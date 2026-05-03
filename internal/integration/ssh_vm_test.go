package integration_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	_ "modernc.org/sqlite"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/firecracker"
	"ssh-microvm/internal/sshserver"
	"ssh-microvm/internal/store"
)

func TestSSHStartsVM(t *testing.T) {
	if _, err := exec.LookPath("firecracker"); err != nil {
		t.Skip("firecracker not found in PATH")
	}
	if _, err := exec.LookPath("ip"); err != nil {
		t.Skip("ip command not found in PATH")
	}
	if err := exec.Command("sudo", "-n", "true").Run(); err != nil {
		t.Skip("sudo without password is required for tap setup")
	}
	if _, err := os.Stat("/dev/kvm"); err != nil {
		t.Skip("/dev/kvm not available: KVM is required for Firecracker")
	}

	kernel := getenvOr("SSH_MICROVM_KERNEL", "artifacts/vmlinux.bin")
	rootfs := getenvOr("SSH_MICROVM_ROOTFS", "artifacts/ubuntu.ext4")
	if _, err := os.Stat(kernel); err != nil {
		t.Skipf("kernel not found: %s", kernel)
	}
	if _, err := os.Stat(rootfs); err != nil {
		t.Skipf("rootfs not found: %s", rootfs)
	}

	stateDir, err := os.MkdirTemp("", "ssh-microvm-test-")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(stateDir)

	guestKey := getenvOr("SSH_MICROVM_GUEST_KEY", "artifacts/ubuntu.id_rsa")
	if _, err := os.Stat(guestKey); err != nil {
		t.Skipf("guest key not found: %s", guestKey)
	}

	cfg := &config.Config{
		ListenAddr:    "127.0.0.1:0",
		StateDir:      stateDir,
		DBPath:        filepath.Join(stateDir, "db.sqlite"),
		HostKeyPath:   filepath.Join(stateDir, "ssh_host_ed25519"),
		AuthMode:      config.AuthModeAutoEnroll,
		Firecracker:   "firecracker",
		KernelImage:   kernel,
		RootFS:        rootfs,
		BootArgs:      "console=ttyS0 reboot=k panic=1 pci=off",
		VCPUCount:     1,
		MemMiB:        256,
		GracefulStopS: 1,
		GuestUser:     "root",
		GuestKeyPath:  guestKey,
		GuestIP:       "172.16.0.2",
		HostIP:        "172.16.0.1",
		TapPrefix:     "tap",
	}

	st, err := store.New(cfg.DBPath)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}
	defer st.Close()

	if err := st.EnsureSchema(context.Background()); err != nil {
		t.Fatalf("db schema: %v", err)
	}

	manager := firecracker.NewManager(cfg)
	server, err := sshserver.New(cfg, st, manager)
	if err != nil {
		t.Fatalf("server init: %v", err)
	}

	ln, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	addr := ln.Addr().String()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ServeListener(ctx, ln)
	}()

	signer, err := newSigner()
	if err != nil {
		t.Fatalf("signer: %v", err)
	}

	clientCfg := &ssh.ClientConfig{
		User:            "tester",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, clientCfg)
	if err != nil {
		t.Fatalf("ssh dial: %v", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		_ = conn.Close()
		t.Fatalf("ssh session: %v", err)
	}
	_ = session.Shell()
	time.Sleep(200 * time.Millisecond)
	_ = session.Close()
	_ = conn.Close()

	if err := waitForVMRecords(cfg.DBPath, 5*time.Second); err != nil {
		t.Fatalf("vm record check: %v", err)
	}

	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("server did not exit")
	}
}

func newSigner() (ssh.Signer, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(priv)
}

func waitForVMRecords(dbPath string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ok, err := hasClosedSessionAndVM(dbPath)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return context.DeadlineExceeded
}

func hasClosedSessionAndVM(dbPath string) (bool, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var vmID, status string
	var endedAt sql.NullString
	row := db.QueryRow(`SELECT vm_id, status, ended_at FROM sessions LIMIT 1`)
	if err := row.Scan(&vmID, &status, &endedAt); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if vmID == "" || status != "closed" || !endedAt.Valid {
		return false, nil
	}

	var vmEnded sql.NullString
	row = db.QueryRow(`SELECT ended_at FROM vms WHERE id = ?`, vmID)
	if err := row.Scan(&vmEnded); err != nil {
		return false, err
	}
	return vmEnded.Valid, nil
}

func getenvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
