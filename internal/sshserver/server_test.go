package sshserver

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/firecracker"
	"ssh-microvm/internal/store"
)

func TestServeListenerReturnsOnContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := newReadyTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ServeListener(ctx, ln)
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

func TestServeRejectsNilDependencies(t *testing.T) {
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	st := newTestStore(t)
	manager := firecracker.NewManager(cfg)
	signer := newTestSigner(t)

	tests := []struct {
		name    string
		server  *Server
		ctx     context.Context
		wantErr string
	}{
		{
			name:    "nil server",
			ctx:     context.Background(),
			wantErr: "server must be set",
		},
		{
			name:    "nil config",
			server:  &Server{},
			ctx:     context.Background(),
			wantErr: "config must be set",
		},
		{
			name:    "nil context",
			server:  &Server{cfg: cfg, store: st, manager: manager, hostSigner: signer},
			wantErr: "context must be set",
		},
		{
			name:    "nil store",
			server:  &Server{cfg: cfg, manager: manager, hostSigner: signer},
			ctx:     context.Background(),
			wantErr: "store must be set",
		},
		{
			name:    "nil manager",
			server:  &Server{cfg: cfg, store: st, hostSigner: signer},
			ctx:     context.Background(),
			wantErr: "firecracker manager must be set",
		},
		{
			name:    "nil host signer",
			server:  &Server{cfg: cfg, store: st, manager: manager},
			ctx:     context.Background(),
			wantErr: "host signer must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Serve(tt.ctx)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Serve error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestServeListenerRejectsNilDependencies(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	cfg := &config.Config{}
	st := newTestStore(t)
	manager := firecracker.NewManager(cfg)
	signer := newTestSigner(t)

	tests := []struct {
		name    string
		server  *Server
		ctx     context.Context
		ln      net.Listener
		wantErr string
	}{
		{
			name:    "nil server",
			ctx:     context.Background(),
			ln:      ln,
			wantErr: "server must be set",
		},
		{
			name:    "nil context",
			server:  &Server{},
			ln:      ln,
			wantErr: "context must be set",
		},
		{
			name:    "nil listener",
			server:  &Server{cfg: cfg, store: st, manager: manager, hostSigner: signer},
			ctx:     context.Background(),
			wantErr: "listener must be set",
		},
		{
			name:    "nil config",
			server:  &Server{store: st, manager: manager, hostSigner: signer},
			ctx:     context.Background(),
			ln:      ln,
			wantErr: "config must be set",
		},
		{
			name:    "nil store",
			server:  &Server{cfg: cfg, manager: manager, hostSigner: signer},
			ctx:     context.Background(),
			ln:      ln,
			wantErr: "store must be set",
		},
		{
			name:    "nil manager",
			server:  &Server{cfg: cfg, store: st, hostSigner: signer},
			ctx:     context.Background(),
			ln:      ln,
			wantErr: "firecracker manager must be set",
		},
		{
			name:    "nil host signer",
			server:  &Server{cfg: cfg, store: st, manager: manager},
			ctx:     context.Background(),
			ln:      ln,
			wantErr: "host signer must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.ServeListener(tt.ctx, tt.ln)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("ServeListener error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestHandleConnRejectsInvalidState(t *testing.T) {
	readyServer := newReadyTestServer(t)
	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	st := newTestStore(t)
	manager := firecracker.NewManager(cfg)
	signer := newTestSigner(t)

	tests := []struct {
		name       string
		server     *Server
		ctx        context.Context
		conn       net.Conn
		wantClosed bool
	}{
		{
			name:       "nil server",
			ctx:        context.Background(),
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:       "nil context",
			server:     &Server{cfg: cfg, store: st, manager: manager, hostSigner: signer},
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:       "nil config",
			server:     &Server{store: st, manager: manager, hostSigner: signer},
			ctx:        context.Background(),
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:       "nil store",
			server:     &Server{cfg: cfg, manager: manager, hostSigner: signer},
			ctx:        context.Background(),
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:       "nil manager",
			server:     &Server{cfg: cfg, store: st, hostSigner: signer},
			ctx:        context.Background(),
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:       "nil host signer",
			server:     &Server{cfg: cfg, store: st, manager: manager},
			ctx:        context.Background(),
			conn:       &testNetConn{},
			wantClosed: true,
		},
		{
			name:   "nil network connection",
			server: readyServer,
			ctx:    context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.server.handleConn(tt.ctx, tt.conn)

			testConn, ok := tt.conn.(*testNetConn)
			if tt.wantClosed {
				if !ok {
					t.Fatalf("test connection has type %T, want *testNetConn", tt.conn)
				}
				if !testConn.closed {
					t.Fatalf("handleConn did not close invalid connection")
				}
			} else if ok && testConn.closed {
				t.Fatalf("handleConn closed connection unexpectedly")
			}
		})
	}
}

func TestNewRejectsNilDependencies(t *testing.T) {
	cfg := &config.Config{HostKeyPath: filepath.Join(t.TempDir(), "ssh_host_ed25519")}
	st := newTestStore(t)
	manager := firecracker.NewManager(cfg)

	tests := []struct {
		name    string
		cfg     *config.Config
		store   *store.Store
		manager *firecracker.Manager
		wantErr string
	}{
		{
			name:    "nil config",
			store:   st,
			manager: manager,
			wantErr: "config must be set",
		},
		{
			name:    "nil store",
			cfg:     cfg,
			manager: manager,
			wantErr: "store must be set",
		},
		{
			name:    "nil manager",
			cfg:     cfg,
			store:   st,
			wantErr: "firecracker manager must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := New(tt.cfg, tt.store, tt.manager)
			if err == nil {
				t.Fatalf("New succeeded, want error containing %q", tt.wantErr)
			}
			if server != nil {
				t.Fatalf("New returned server %#v, want nil", server)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("New error = %q, want to contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestStopVMRejectsInvalidState(t *testing.T) {
	validServer := &Server{cfg: &config.Config{GracefulStopS: 1}}
	validVM := &firecracker.VM{}

	tests := []struct {
		name    string
		server  *Server
		vm      *firecracker.VM
		wantErr string
	}{
		{
			name:    "nil server",
			vm:      validVM,
			wantErr: "server must be set",
		},
		{
			name:    "nil config",
			server:  &Server{},
			vm:      validVM,
			wantErr: "config must be set",
		},
		{
			name:    "nil VM",
			server:  validServer,
			wantErr: "vm not available",
		},
		{
			name:    "non-positive graceful shutdown timeout",
			server:  &Server{cfg: &config.Config{}},
			vm:      validVM,
			wantErr: "graceful shutdown timeout must be > 0",
		},
		{
			name:    "invalid tap name",
			server:  validServer,
			vm:      &firecracker.VM{TapName: "tap/bad"},
			wantErr: "tap name must contain only ASCII letters and digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.stopVM(tt.vm)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("stopVM error = %v, want containing %q", err, tt.wantErr)
			}
		})
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

func TestLoadOrCreateHostKeyRejectsPathWithSurroundingWhitespaceBeforeSideEffects(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	_, err := loadOrCreateHostKey(" " + filepath.Join(workDir, "ssh_host_ed25519") + " ")
	if err == nil {
		t.Fatalf("loadOrCreateHostKey accepted path with surrounding whitespace")
	}
	if !strings.Contains(err.Error(), "host key path must not contain surrounding whitespace") {
		t.Fatalf("loadOrCreateHostKey error = %q, want surrounding whitespace error", err)
	}
	entries, readErr := os.ReadDir(workDir)
	if readErr != nil {
		t.Fatalf("read work dir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("loadOrCreateHostKey created filesystem entries before validating path: %v", entries)
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

func TestPublicKeyCallbackRejectsInvalidState(t *testing.T) {
	key := newTestSigner(t).PublicKey()

	tests := []struct {
		name    string
		server  *Server
		key     ssh.PublicKey
		wantErr string
	}{
		{
			name:    "nil server",
			key:     key,
			wantErr: "server must be set",
		},
		{
			name:    "nil config",
			server:  &Server{},
			key:     key,
			wantErr: "config must be set",
		},
		{
			name:    "nil public key",
			server:  &Server{cfg: &config.Config{AuthMode: config.AuthModeAutoEnroll}},
			wantErr: "public key must be set",
		},
		{
			name:    "known keys mode nil store",
			server:  &Server{cfg: &config.Config{AuthMode: config.AuthModeKnownKeys}},
			key:     key,
			wantErr: "store must be set",
		},
		{
			name:    "invalid auth mode",
			server:  &Server{cfg: &config.Config{AuthMode: "bogus"}},
			key:     key,
			wantErr: "invalid auth mode: bogus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perms, err := tt.server.publicKeyCallback(nil, tt.key)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("publicKeyCallback error = %v, want containing %q", err, tt.wantErr)
			}
			if perms != nil {
				t.Fatalf("publicKeyCallback permissions = %#v, want nil", perms)
			}
		})
	}
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
	for _, tt := range []struct {
		name    string
		term    string
		columns uint32
		rows    uint32
	}{
		{name: "blank term", term: " \t ", columns: 120, rows: 40},
		{name: "term with surrounding whitespace", term: " xterm-256color ", columns: 120, rows: 40},
		{name: "zero columns", term: "xterm-256color", rows: 40},
		{name: "zero rows", term: "xterm-256color", columns: 120},
	} {
		payload := ssh.Marshal(struct {
			Term          string
			Columns, Rows uint32
			Width, Height uint32
			TerminalModes []byte
		}{
			Term:    tt.term,
			Columns: tt.columns,
			Rows:    tt.rows,
		})
		if _, ok := parsePtyRequest(payload); ok {
			t.Fatalf("parsePtyRequest accepted invalid %s payload", tt.name)
		}
	}
	for _, tt := range []struct {
		name    string
		columns uint32
		rows    uint32
	}{
		{name: "zero columns", rows: 33},
		{name: "zero rows", columns: 100},
	} {
		payload := ssh.Marshal(struct {
			Columns, Rows uint32
			Width, Height uint32
		}{
			Columns: tt.columns,
			Rows:    tt.rows,
		})
		if _, ok := parseWindowChange(payload); ok {
			t.Fatalf("parseWindowChange accepted invalid %s payload", tt.name)
		}
	}
	for _, command := range []string{"", " \t ", " echo hi "} {
		payload := ssh.Marshal(struct {
			Command string
		}{Command: command})
		if _, ok := parseExecRequest(payload); ok {
			t.Fatalf("parseExecRequest accepted invalid command %q", command)
		}
	}
}

func TestProxyToGuestRejectsInvalidState(t *testing.T) {
	validCfg := &config.Config{
		GuestUser:    "root",
		GuestKeyPath: "/keys/guest",
	}
	validServer := &Server{cfg: validCfg}
	validChannel := &testSSHChannel{}
	validWinCh := make(chan windowChange)
	validVM := &firecracker.VM{GuestIP: "127.0.0.1"}

	tests := []struct {
		name    string
		server  *Server
		channel ssh.Channel
		ptyReq  *ptyRequest
		winCh   <-chan windowChange
		shell   bool
		execCmd string
		vm      *firecracker.VM
		wantErr string
	}{
		{
			name:    "nil server",
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "server must be set",
		},
		{
			name:    "nil config",
			server:  &Server{},
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "config must be set",
		},
		{
			name:    "nil channel",
			server:  validServer,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "ssh channel must be set",
		},
		{
			name:    "nil VM",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			wantErr: "vm not available",
		},
		{
			name:    "blank pty terminal",
			server:  validServer,
			channel: validChannel,
			ptyReq:  &ptyRequest{Term: " \t ", Width: 80, Height: 24},
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "pty terminal must be set",
		},
		{
			name:    "pty terminal with surrounding whitespace",
			server:  validServer,
			channel: validChannel,
			ptyReq:  &ptyRequest{Term: " xterm ", Width: 80, Height: 24},
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "pty terminal must not contain surrounding whitespace",
		},
		{
			name:    "zero pty width",
			server:  validServer,
			channel: validChannel,
			ptyReq:  &ptyRequest{Term: "xterm", Height: 24},
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "pty dimensions must be positive",
		},
		{
			name:    "zero pty height",
			server:  validServer,
			channel: validChannel,
			ptyReq:  &ptyRequest{Term: "xterm", Width: 80},
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "pty dimensions must be positive",
		},
		{
			name:    "nil window change channel",
			server:  validServer,
			channel: validChannel,
			shell:   true,
			vm:      validVM,
			wantErr: "window change channel must be set",
		},
		{
			name:    "blank exec command",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			execCmd: " \t ",
			vm:      validVM,
			wantErr: "exec command must be set",
		},
		{
			name:    "exec command with surrounding whitespace",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			execCmd: " echo hi ",
			vm:      validVM,
			wantErr: "exec command must not contain surrounding whitespace",
		},
		{
			name:    "exec command with shell session",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			execCmd: "echo hi",
			vm:      validVM,
			wantErr: "exec command cannot be set for shell sessions",
		},
		{
			name:    "blank guest IP",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      &firecracker.VM{GuestIP: " \t "},
			wantErr: "guest IP must be set",
		},
		{
			name:    "invalid guest IP",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      &firecracker.VM{GuestIP: "not-an-ip"},
			wantErr: "guest IP must be a valid IPv4 address",
		},
		{
			name:    "non-canonical guest IP",
			server:  validServer,
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      &firecracker.VM{GuestIP: " 127.0.0.1 "},
			wantErr: "guest IP must be a valid IPv4 address",
		},
		{
			name:    "blank guest user",
			server:  &Server{cfg: &config.Config{GuestUser: " \t ", GuestKeyPath: "/keys/guest"}},
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "guest user must be set",
		},
		{
			name:    "guest user with surrounding whitespace",
			server:  &Server{cfg: &config.Config{GuestUser: " root ", GuestKeyPath: "/keys/guest"}},
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "guest user must not contain surrounding whitespace",
		},
		{
			name:    "blank guest key",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: " \t "}},
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "guest key path must be set",
		},
		{
			name:    "guest key with surrounding whitespace",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: " /keys/guest "}},
			channel: validChannel,
			winCh:   validWinCh,
			shell:   true,
			vm:      validVM,
			wantErr: "guest key path must not contain surrounding whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.proxyToGuest(tt.channel, tt.ptyReq, tt.winCh, tt.shell, tt.execCmd, tt.vm)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("proxyToGuest error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestHandleChannelsRejectsInvalidState(t *testing.T) {
	validVM := &firecracker.VM{GuestIP: "127.0.0.1"}
	nilNewChannels := make(chan ssh.NewChannel, 1)
	nilNewChannels <- nil
	close(nilNewChannels)

	tests := []struct {
		name     string
		server   *Server
		channels <-chan ssh.NewChannel
		vm       *firecracker.VM
	}{
		{
			name:     "nil server",
			channels: make(chan ssh.NewChannel),
			vm:       validVM,
		},
		{
			name:   "nil channel stream",
			server: &Server{},
			vm:     validVM,
		},
		{
			name:     "nil VM",
			server:   &Server{},
			channels: make(chan ssh.NewChannel),
		},
		{
			name:     "nil new channel",
			server:   &Server{},
			channels: nilNewChannels,
			vm:       validVM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				tt.server.handleChannels(tt.channels, tt.vm)
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("handleChannels did not return for invalid state")
			}
		})
	}
}

func TestHandleSessionRejectsInvalidState(t *testing.T) {
	validRequests := make(chan *ssh.Request)
	close(validRequests)
	validVM := &firecracker.VM{GuestIP: "127.0.0.1"}

	tests := []struct {
		name       string
		server     *Server
		channel    *testSSHChannel
		requests   <-chan *ssh.Request
		vm         *firecracker.VM
		wantClosed bool
	}{
		{
			name:     "nil channel",
			server:   &Server{cfg: &config.Config{}},
			requests: validRequests,
			vm:       validVM,
		},
		{
			name:       "nil server",
			channel:    &testSSHChannel{},
			requests:   validRequests,
			vm:         validVM,
			wantClosed: true,
		},
		{
			name:       "nil config",
			server:     &Server{},
			channel:    &testSSHChannel{},
			requests:   validRequests,
			vm:         validVM,
			wantClosed: true,
		},
		{
			name:       "nil requests",
			server:     &Server{cfg: &config.Config{}},
			channel:    &testSSHChannel{},
			vm:         validVM,
			wantClosed: true,
		},
		{
			name:       "nil VM",
			server:     &Server{cfg: &config.Config{}},
			channel:    &testSSHChannel{},
			requests:   validRequests,
			wantClosed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})
			go func() {
				tt.server.handleSession(tt.channel, tt.requests, tt.vm)
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("handleSession did not return for invalid state")
			}

			if tt.channel != nil && tt.channel.closed != tt.wantClosed {
				t.Fatalf("channel closed = %t, want %t", tt.channel.closed, tt.wantClosed)
			}
		})
	}
}

func TestHandleSessionSkipsNilRequests(t *testing.T) {
	requests := make(chan *ssh.Request, 1)
	requests <- nil
	close(requests)
	channel := &testSSHChannel{}
	server := &Server{cfg: &config.Config{}}

	done := make(chan struct{})
	go func() {
		server.handleSession(channel, requests, &firecracker.VM{GuestIP: "127.0.0.1"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("handleSession did not return after nil request stream closed")
	}

	if !channel.closed {
		t.Fatalf("handleSession did not close channel")
	}
}

func TestDialGuestRejectsInvalidState(t *testing.T) {
	tests := []struct {
		name    string
		server  *Server
		guestIP string
		wantErr string
	}{
		{
			name:    "nil server",
			guestIP: "127.0.0.1",
			wantErr: "server must be set",
		},
		{
			name:    "nil config",
			server:  &Server{},
			guestIP: "127.0.0.1",
			wantErr: "config must be set",
		},
		{
			name:    "blank guest IP",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: "/keys/guest"}},
			guestIP: " \t ",
			wantErr: "guest IP must be set",
		},
		{
			name:    "invalid guest IP",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: "/keys/guest"}},
			guestIP: "not-an-ip",
			wantErr: "guest IP must be a valid IPv4 address",
		},
		{
			name:    "non-canonical guest IP",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: "/keys/guest"}},
			guestIP: " 127.0.0.1 ",
			wantErr: "guest IP must be a valid IPv4 address",
		},
		{
			name:    "blank guest user",
			server:  &Server{cfg: &config.Config{GuestUser: " \t ", GuestKeyPath: "/keys/guest"}},
			guestIP: "127.0.0.1",
			wantErr: "guest user must be set",
		},
		{
			name:    "guest user with surrounding whitespace",
			server:  &Server{cfg: &config.Config{GuestUser: " root ", GuestKeyPath: "/keys/guest"}},
			guestIP: "127.0.0.1",
			wantErr: "guest user must not contain surrounding whitespace",
		},
		{
			name:    "blank guest key",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: " \t "}},
			guestIP: "127.0.0.1",
			wantErr: "guest key path must be set",
		},
		{
			name:    "guest key with surrounding whitespace",
			server:  &Server{cfg: &config.Config{GuestUser: "root", GuestKeyPath: " /keys/guest "}},
			guestIP: "127.0.0.1",
			wantErr: "guest key path must not contain surrounding whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := tt.server.dialGuest(tt.guestIP)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("dialGuest error = %v, want containing %q", err, tt.wantErr)
			}
			if client != nil {
				t.Fatalf("dialGuest client = %#v, want nil", client)
			}
		})
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

func TestWaitForPortRejectsInvalidState(t *testing.T) {
	validDial := func(string, time.Duration) (net.Conn, error) {
		return nil, errors.New("dial should not be called for invalid state")
	}

	tests := []struct {
		name    string
		addr    string
		timeout time.Duration
		dial    func(string, time.Duration) (net.Conn, error)
		wantErr string
	}{
		{
			name:    "blank address",
			addr:    " \t ",
			timeout: time.Second,
			dial:    validDial,
			wantErr: "guest port address must be set",
		},
		{
			name:    "address with surrounding whitespace",
			addr:    " 127.0.0.1:22 ",
			timeout: time.Second,
			dial:    validDial,
			wantErr: "guest port address must not contain surrounding whitespace",
		},
		{
			name:    "malformed address",
			addr:    "127.0.0.1",
			timeout: time.Second,
			dial:    validDial,
			wantErr: "guest port address must be a valid TCP address",
		},
		{
			name:    "invalid port",
			addr:    "127.0.0.1:not-a-port",
			timeout: time.Second,
			dial:    validDial,
			wantErr: "guest port address port must be valid",
		},
		{
			name:    "zero timeout",
			addr:    "127.0.0.1:22",
			timeout: 0,
			dial:    validDial,
			wantErr: "guest port timeout must be positive",
		},
		{
			name:    "negative timeout",
			addr:    "127.0.0.1:22",
			timeout: -time.Second,
			dial:    validDial,
			wantErr: "guest port timeout must be positive",
		},
		{
			name:    "nil dial function",
			addr:    "127.0.0.1:22",
			timeout: time.Second,
			wantErr: "guest port dial function must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := waitForPortWithDial(tt.addr, tt.timeout, tt.dial)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
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

func newReadyTestServer(t *testing.T) *Server {
	t.Helper()

	cfg := &config.Config{ListenAddr: "127.0.0.1:0"}
	return &Server{
		cfg:        cfg,
		store:      newTestStore(t),
		manager:    firecracker.NewManager(cfg),
		hostSigner: newTestSigner(t),
	}
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

type testSSHChannel struct {
	bytes.Buffer
	closed bool
}

func (c *testSSHChannel) Close() error {
	c.closed = true
	return nil
}

func (c *testSSHChannel) CloseWrite() error {
	return nil
}

func (c *testSSHChannel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return false, nil
}

func (c *testSSHChannel) Stderr() io.ReadWriter {
	return &bytes.Buffer{}
}

type testNetConn struct {
	closed bool
}

func (c *testNetConn) Read(_ []byte) (int, error) {
	return 0, io.EOF
}

func (c *testNetConn) Write(p []byte) (int, error) {
	return len(p), nil
}

func (c *testNetConn) Close() error {
	c.closed = true
	return nil
}

func (c *testNetConn) LocalAddr() net.Addr {
	return testAddr("local")
}

func (c *testNetConn) RemoteAddr() net.Addr {
	return testAddr("remote")
}

func (c *testNetConn) SetDeadline(_ time.Time) error {
	return nil
}

func (c *testNetConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (c *testNetConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

type testAddr string

func (a testAddr) Network() string {
	return string(a)
}

func (a testAddr) String() string {
	return string(a)
}
