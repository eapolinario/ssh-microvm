package sshserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/store"
)

func TestServeListenerReturnsOnContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- (&Server{}).ServeListener(ctx, ln)
	}()

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("ServeListener returned error after context cancellation: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatalf("ServeListener did not return after context cancellation")
	}
}

func TestLoadOrCreateHostKeyCreatesAndReusesKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ssh_host_ed25519")

	signer, err := loadOrCreateHostKey(path)
	if err != nil {
		t.Fatalf("loadOrCreateHostKey create: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat host key: %v", err)
	}
	if got, want := info.Mode().Perm(), os.FileMode(0o600); got != want {
		t.Fatalf("host key mode = %v, want %v", got, want)
	}

	reloaded, err := loadOrCreateHostKey(path)
	if err != nil {
		t.Fatalf("loadOrCreateHostKey reload: %v", err)
	}
	if got, want := ssh.FingerprintSHA256(reloaded.PublicKey()), ssh.FingerprintSHA256(signer.PublicKey()); got != want {
		t.Fatalf("reloaded key fingerprint = %q, want %q", got, want)
	}
}

func TestLoadOrCreateHostKeyRejectsOpenPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ssh_host_ed25519")
	if _, err := loadOrCreateHostKey(path); err != nil {
		t.Fatalf("loadOrCreateHostKey create: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("chmod host key: %v", err)
	}

	_, err := loadOrCreateHostKey(path)
	if err == nil {
		t.Fatalf("loadOrCreateHostKey accepted host key with open permissions")
	}
	if !strings.Contains(err.Error(), "permissions too open") {
		t.Fatalf("loadOrCreateHostKey error = %q, want permissions error", err)
	}
}

func TestLoadOrCreateHostKeyRejectsWritableParentDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "open")
	if err := os.Mkdir(dir, 0o777); err != nil {
		t.Fatalf("mkdir open host key dir: %v", err)
	}
	if err := os.Chmod(dir, 0o777); err != nil {
		t.Fatalf("chmod open host key dir: %v", err)
	}

	_, err := loadOrCreateHostKey(filepath.Join(dir, "ssh_host_ed25519"))
	if err == nil {
		t.Fatalf("loadOrCreateHostKey accepted a group/world-writable parent directory")
	}
	if !strings.Contains(err.Error(), "permissions too open") {
		t.Fatalf("loadOrCreateHostKey error = %q, want permissions error", err)
	}
}

func TestLoadOrCreateHostKeyRejectsBlankPath(t *testing.T) {
	_, err := loadOrCreateHostKey(" \t ")
	if err == nil {
		t.Fatalf("loadOrCreateHostKey accepted a blank path")
	}
	if !strings.Contains(err.Error(), "host key path must be set") {
		t.Fatalf("loadOrCreateHostKey error = %q, want blank path error", err)
	}
}

func TestPublicKeyCallbackAuthModes(t *testing.T) {
	signer := newTestSigner(t)
	key := signer.PublicKey()
	fingerprint := ssh.FingerprintSHA256(key)
	authorizedKey := string(ssh.MarshalAuthorizedKey(key))

	t.Run("auto enroll mode accepts unknown key", func(t *testing.T) {
		server := &Server{
			cfg:   &config.Config{AuthMode: config.AuthModeAutoEnroll},
			store: newTestStore(t),
		}

		perms, err := server.publicKeyCallback(nil, key)
		if err != nil {
			t.Fatalf("publicKeyCallback: %v", err)
		}
		assertPermissionExtension(t, perms, "pubkey-fp", fingerprint)
		assertPermissionExtension(t, perms, "pubkey", authorizedKey)
	})

	t.Run("known keys mode rejects unknown key", func(t *testing.T) {
		server := &Server{
			cfg:   &config.Config{AuthMode: config.AuthModeKnownKeys},
			store: newTestStore(t),
		}

		if _, err := server.publicKeyCallback(nil, key); err == nil {
			t.Fatalf("publicKeyCallback accepted unknown key in known-keys mode")
		}
	})

	t.Run("known keys mode accepts enrolled key", func(t *testing.T) {
		st := newTestStore(t)
		if _, err := st.EnsureUserAndKey(context.Background(), "alice", fingerprint, authorizedKey); err != nil {
			t.Fatalf("EnsureUserAndKey: %v", err)
		}
		server := &Server{
			cfg:   &config.Config{AuthMode: config.AuthModeKnownKeys},
			store: st,
		}

		perms, err := server.publicKeyCallback(nil, key)
		if err != nil {
			t.Fatalf("publicKeyCallback: %v", err)
		}
		assertPermissionExtension(t, perms, "pubkey-fp", fingerprint)
		assertPermissionExtension(t, perms, "pubkey", authorizedKey)
	})
}

func TestParseSSHRequestPayloads(t *testing.T) {
	ptyPayload := ssh.Marshal(struct {
		Term          string
		Columns, Rows uint32
		Width, Height uint32
		TerminalModes []byte
	}{
		Term:          "xterm-256color",
		Columns:       120,
		Rows:          40,
		Width:         800,
		Height:        600,
		TerminalModes: nil,
	})
	pty, ok := parsePtyRequest(ptyPayload)
	if !ok {
		t.Fatalf("parsePtyRequest returned !ok")
	}
	if pty.Term != "xterm-256color" || pty.Width != 120 || pty.Height != 40 {
		t.Fatalf("parsePtyRequest = %+v, want term xterm-256color width 120 height 40", pty)
	}

	windowPayload := ssh.Marshal(struct {
		Columns, Rows uint32
		Width, Height uint32
	}{
		Columns: 100,
		Rows:    33,
		Width:   640,
		Height:  480,
	})
	win, ok := parseWindowChange(windowPayload)
	if !ok {
		t.Fatalf("parseWindowChange returned !ok")
	}
	if win.Width != 100 || win.Height != 33 {
		t.Fatalf("parseWindowChange = %+v, want width 100 height 33", win)
	}

	execPayload := ssh.Marshal(struct {
		Command string
	}{Command: "echo hi"})
	cmd, ok := parseExecRequest(execPayload)
	if !ok {
		t.Fatalf("parseExecRequest returned !ok")
	}
	if cmd != "echo hi" {
		t.Fatalf("parseExecRequest = %q, want %q", cmd, "echo hi")
	}
}

func TestParseSSHRequestPayloadsRejectInvalidData(t *testing.T) {
	invalid := []byte{0, 1, 2}

	if _, ok := parsePtyRequest(invalid); ok {
		t.Fatalf("parsePtyRequest accepted invalid payload")
	}
	if _, ok := parseWindowChange(invalid); ok {
		t.Fatalf("parseWindowChange accepted invalid payload")
	}
	if _, ok := parseExecRequest(invalid); ok {
		t.Fatalf("parseExecRequest accepted invalid payload")
	}
}

func TestWaitForPort(t *testing.T) {
	t.Run("ready listener", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		t.Cleanup(func() {
			_ = ln.Close()
		})

		if err := waitForPort(ln.Addr().String(), time.Second); err != nil {
			t.Fatalf("waitForPort ready listener: %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("listen: %v", err)
		}
		addr := ln.Addr().String()
		if err := ln.Close(); err != nil {
			t.Fatalf("close listener: %v", err)
		}

		err = waitForPort(addr, 20*time.Millisecond)
		if err == nil {
			t.Fatalf("waitForPort succeeded for closed listener")
		}
		if !strings.Contains(err.Error(), "timeout waiting for "+addr) {
			t.Fatalf("waitForPort error = %q, want timeout for %s", err, addr)
		}
	})
}

func TestWaitForPortCapsDialTimeoutToRemainingDeadline(t *testing.T) {
	var maxDialTimeout time.Duration
	err := waitForPortWithDial("203.0.113.1:22", 20*time.Millisecond, func(_ string, timeout time.Duration) (net.Conn, error) {
		if timeout > maxDialTimeout {
			maxDialTimeout = timeout
		}
		time.Sleep(timeout)
		return nil, errors.New("dial timeout")
	})
	if err == nil {
		t.Fatalf("waitForPortWithDial succeeded, want timeout")
	}
	if maxDialTimeout > 20*time.Millisecond {
		t.Fatalf("dial timeout = %v, want capped to remaining deadline", maxDialTimeout)
	}
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()

	st, err := store.New(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("store.Close: %v", err)
		}
	})

	if err := st.EnsureSchema(context.Background()); err != nil {
		t.Fatalf("EnsureSchema: %v", err)
	}
	return st
}

func newTestSigner(t *testing.T) ssh.Signer {
	t.Helper()

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatalf("NewSignerFromKey: %v", err)
	}
	return signer
}

func assertPermissionExtension(t *testing.T, perms *ssh.Permissions, key, want string) {
	t.Helper()

	if perms == nil {
		t.Fatalf("permissions are nil")
	}
	got, ok := perms.Extensions[key]
	if !ok {
		t.Fatalf("permissions missing extension %q", key)
	}
	if got != want {
		t.Fatalf("permissions extension %q = %q, want %q", key, got, want)
	}
}
