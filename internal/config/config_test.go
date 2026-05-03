package config

import (
	"flag"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFromArgsAppliesDefaultsAndDerivedPaths(t *testing.T) {
	stateDir := t.TempDir()

	cfg, err := loadFromArgs([]string{
		"--state-dir", stateDir,
		"--kernel", "/images/vmlinux.bin",
		"--rootfs", "/images/rootfs.ext4",
	}, flag.ContinueOnError)
	if err != nil {
		t.Fatalf("loadFromArgs: %v", err)
	}

	if cfg.ListenAddr != ":2222" {
		t.Fatalf("ListenAddr = %q, want :2222", cfg.ListenAddr)
	}
	if cfg.DBPath != filepath.Join(stateDir, "db.sqlite") {
		t.Fatalf("DBPath = %q, want derived state DB path", cfg.DBPath)
	}
	if cfg.HostKeyPath != filepath.Join(stateDir, "ssh_host_ed25519") {
		t.Fatalf("HostKeyPath = %q, want derived state host key path", cfg.HostKeyPath)
	}
	if cfg.AuthMode != AuthModeAutoEnroll {
		t.Fatalf("AuthMode = %q, want %q", cfg.AuthMode, AuthModeAutoEnroll)
	}
	if cfg.Firecracker != "firecracker" {
		t.Fatalf("Firecracker = %q, want firecracker", cfg.Firecracker)
	}
	if cfg.KernelImage != "/images/vmlinux.bin" || cfg.RootFS != "/images/rootfs.ext4" {
		t.Fatalf("kernel/rootfs = %q/%q, want explicit paths", cfg.KernelImage, cfg.RootFS)
	}
	if cfg.VCPUCount != 1 || cfg.MemMiB != 512 || cfg.GracefulStopS != 2 {
		t.Fatalf("resource defaults = vcpu %d mem %d grace %d, want 1/512/2", cfg.VCPUCount, cfg.MemMiB, cfg.GracefulStopS)
	}
	if cfg.GuestUser != "root" || cfg.GuestKeyPath != "artifacts/ubuntu.id_rsa" {
		t.Fatalf("guest defaults = %q/%q, want root/artifacts/ubuntu.id_rsa", cfg.GuestUser, cfg.GuestKeyPath)
	}
	if cfg.GuestIP != "172.16.0.2" || cfg.HostIP != "172.16.0.1" || cfg.TapPrefix != "tap" {
		t.Fatalf("network defaults = guest %q host %q tap %q, want 172.16.0.2/172.16.0.1/tap", cfg.GuestIP, cfg.HostIP, cfg.TapPrefix)
	}
}

func TestLoadFromArgsAppliesExplicitOverrides(t *testing.T) {
	cfg, err := loadFromArgs([]string{
		"--listen", "127.0.0.1:2200",
		"--state-dir", "/state",
		"--db-path", "/db/custom.sqlite",
		"--host-key", "/keys/host",
		"--auth-mode", AuthModeKnownKeys,
		"--firecracker", "/bin/firecracker",
		"--kernel", "/images/kernel",
		"--rootfs", "/images/rootfs",
		"--boot-args", "console=ttyS0 ip=dhcp",
		"--vcpu", "2",
		"--mem", "1024",
		"--grace-stop", "5",
		"--guest-user", "ubuntu",
		"--guest-key", "/keys/guest",
		"--guest-ip", "10.0.0.2",
		"--host-ip", "10.0.0.1",
		"--tap-prefix", "vm",
	}, flag.ContinueOnError)
	if err != nil {
		t.Fatalf("loadFromArgs: %v", err)
	}

	assertConfigValue(t, "ListenAddr", cfg.ListenAddr, "127.0.0.1:2200")
	assertConfigValue(t, "StateDir", cfg.StateDir, "/state")
	assertConfigValue(t, "DBPath", cfg.DBPath, "/db/custom.sqlite")
	assertConfigValue(t, "HostKeyPath", cfg.HostKeyPath, "/keys/host")
	assertConfigValue(t, "AuthMode", cfg.AuthMode, AuthModeKnownKeys)
	assertConfigValue(t, "Firecracker", cfg.Firecracker, "/bin/firecracker")
	assertConfigValue(t, "KernelImage", cfg.KernelImage, "/images/kernel")
	assertConfigValue(t, "RootFS", cfg.RootFS, "/images/rootfs")
	assertConfigValue(t, "BootArgs", cfg.BootArgs, "console=ttyS0 ip=dhcp")
	assertConfigValue(t, "GuestUser", cfg.GuestUser, "ubuntu")
	assertConfigValue(t, "GuestKeyPath", cfg.GuestKeyPath, "/keys/guest")
	assertConfigValue(t, "GuestIP", cfg.GuestIP, "10.0.0.2")
	assertConfigValue(t, "HostIP", cfg.HostIP, "10.0.0.1")
	assertConfigValue(t, "TapPrefix", cfg.TapPrefix, "vm")
	if cfg.VCPUCount != 2 || cfg.MemMiB != 1024 || cfg.GracefulStopS != 5 {
		t.Fatalf("resource overrides = vcpu %d mem %d grace %d, want 2/1024/5", cfg.VCPUCount, cfg.MemMiB, cfg.GracefulStopS)
	}
}

func TestLoadFromArgsDerivesBlankOptionalPaths(t *testing.T) {
	stateDir := t.TempDir()

	cfg, err := loadFromArgs([]string{
		"--state-dir", stateDir,
		"--db-path", " \t ",
		"--host-key", " \t ",
		"--kernel", "/images/vmlinux.bin",
		"--rootfs", "/images/rootfs.ext4",
	}, flag.ContinueOnError)
	if err != nil {
		t.Fatalf("loadFromArgs: %v", err)
	}

	assertConfigValue(t, "DBPath", cfg.DBPath, filepath.Join(stateDir, "db.sqlite"))
	assertConfigValue(t, "HostKeyPath", cfg.HostKeyPath, filepath.Join(stateDir, "ssh_host_ed25519"))
}

func TestLoadFromArgsTrimsStringOverrides(t *testing.T) {
	stateDir := t.TempDir()

	cfg, err := loadFromArgs([]string{
		"--listen", " 127.0.0.1:2200 ",
		"--state-dir", " " + stateDir + " ",
		"--db-path", " " + filepath.Join(stateDir, "custom.sqlite") + " ",
		"--host-key", " " + filepath.Join(stateDir, "host_key") + " ",
		"--auth-mode", " " + AuthModeKnownKeys + " ",
		"--firecracker", " /bin/firecracker ",
		"--kernel", " /images/kernel ",
		"--rootfs", " /images/rootfs ",
		"--guest-user", " ubuntu ",
		"--guest-key", " /keys/guest ",
		"--guest-ip", " 10.0.0.2 ",
		"--host-ip", " 10.0.0.1 ",
		"--tap-prefix", " vm ",
	}, flag.ContinueOnError)
	if err != nil {
		t.Fatalf("loadFromArgs: %v", err)
	}

	assertConfigValue(t, "ListenAddr", cfg.ListenAddr, "127.0.0.1:2200")
	assertConfigValue(t, "StateDir", cfg.StateDir, stateDir)
	assertConfigValue(t, "DBPath", cfg.DBPath, filepath.Join(stateDir, "custom.sqlite"))
	assertConfigValue(t, "HostKeyPath", cfg.HostKeyPath, filepath.Join(stateDir, "host_key"))
	assertConfigValue(t, "AuthMode", cfg.AuthMode, AuthModeKnownKeys)
	assertConfigValue(t, "Firecracker", cfg.Firecracker, "/bin/firecracker")
	assertConfigValue(t, "KernelImage", cfg.KernelImage, "/images/kernel")
	assertConfigValue(t, "RootFS", cfg.RootFS, "/images/rootfs")
	assertConfigValue(t, "GuestUser", cfg.GuestUser, "ubuntu")
	assertConfigValue(t, "GuestKeyPath", cfg.GuestKeyPath, "/keys/guest")
	assertConfigValue(t, "GuestIP", cfg.GuestIP, "10.0.0.2")
	assertConfigValue(t, "HostIP", cfg.HostIP, "10.0.0.1")
	assertConfigValue(t, "TapPrefix", cfg.TapPrefix, "vm")
}

func TestLoadFromArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing kernel",
			args:    []string{"--rootfs", "/rootfs"},
			wantErr: "--kernel is required",
		},
		{
			name:    "missing rootfs",
			args:    []string{"--kernel", "/kernel"},
			wantErr: "--rootfs is required",
		},
		{
			name:    "blank kernel",
			args:    []string{"--kernel", " \t ", "--rootfs", "/rootfs"},
			wantErr: "--kernel is required",
		},
		{
			name:    "blank rootfs",
			args:    []string{"--kernel", "/kernel", "--rootfs", " \t "},
			wantErr: "--rootfs is required",
		},
		{
			name:    "invalid listen address",
			args:    requiredArgs("--listen", "127.0.0.1"),
			wantErr: "--listen must be a valid TCP address: 127.0.0.1",
		},
		{
			name:    "invalid listen port",
			args:    requiredArgs("--listen", "127.0.0.1:not-a-port"),
			wantErr: "--listen port must be valid: not-a-port",
		},
		{
			name:    "invalid listen host",
			args:    requiredArgs("--listen", "999.0.0.1:2222"),
			wantErr: "--listen must resolve to a valid TCP address: 999.0.0.1:2222",
		},
		{
			name:    "missing state dir",
			args:    requiredArgs("--state-dir", ""),
			wantErr: "--state-dir must be set",
		},
		{
			name:    "blank state dir",
			args:    requiredArgs("--state-dir", " \t "),
			wantErr: "--state-dir must be set",
		},
		{
			name:    "invalid auth mode",
			args:    requiredArgs("--auth-mode", "invalid"),
			wantErr: "invalid --auth-mode: invalid",
		},
		{
			name:    "missing firecracker binary",
			args:    requiredArgs("--firecracker", ""),
			wantErr: "--firecracker must be set",
		},
		{
			name:    "blank firecracker binary",
			args:    requiredArgs("--firecracker", " \t "),
			wantErr: "--firecracker must be set",
		},
		{
			name:    "non-positive vcpu",
			args:    requiredArgs("--vcpu", "0"),
			wantErr: "--vcpu must be > 0",
		},
		{
			name:    "non-positive memory",
			args:    requiredArgs("--mem", "0"),
			wantErr: "--mem must be > 0",
		},
		{
			name:    "non-positive graceful stop timeout",
			args:    requiredArgs("--grace-stop", "0"),
			wantErr: "--grace-stop must be > 0",
		},
		{
			name:    "missing guest user",
			args:    requiredArgs("--guest-user", ""),
			wantErr: "--guest-user must be set",
		},
		{
			name:    "blank guest user",
			args:    requiredArgs("--guest-user", " \t "),
			wantErr: "--guest-user must be set",
		},
		{
			name:    "missing guest key",
			args:    requiredArgs("--guest-key", ""),
			wantErr: "--guest-key must be set",
		},
		{
			name:    "blank guest key",
			args:    requiredArgs("--guest-key", " \t "),
			wantErr: "--guest-key must be set",
		},
		{
			name:    "missing guest ip",
			args:    requiredArgs("--guest-ip", ""),
			wantErr: "--guest-ip and --host-ip must be set",
		},
		{
			name:    "missing host ip",
			args:    requiredArgs("--host-ip", ""),
			wantErr: "--guest-ip and --host-ip must be set",
		},
		{
			name:    "blank guest ip",
			args:    requiredArgs("--guest-ip", " \t "),
			wantErr: "--guest-ip and --host-ip must be set",
		},
		{
			name:    "blank host ip",
			args:    requiredArgs("--host-ip", " \t "),
			wantErr: "--guest-ip and --host-ip must be set",
		},
		{
			name:    "invalid guest ip",
			args:    requiredArgs("--guest-ip", "not-an-ip"),
			wantErr: "--guest-ip must be a valid IPv4 address: not-an-ip",
		},
		{
			name:    "invalid host ip",
			args:    requiredArgs("--host-ip", "999.0.0.1"),
			wantErr: "--host-ip must be a valid IPv4 address: 999.0.0.1",
		},
		{
			name:    "ipv6 guest ip",
			args:    requiredArgs("--guest-ip", "2001:db8::2"),
			wantErr: "--guest-ip must be a valid IPv4 address: 2001:db8::2",
		},
		{
			name:    "ipv4 mapped ipv6 guest ip",
			args:    requiredArgs("--guest-ip", "::ffff:172.16.0.2"),
			wantErr: "--guest-ip must be a valid IPv4 address: ::ffff:172.16.0.2",
		},
		{
			name:    "ipv4 mapped ipv6 host ip",
			args:    requiredArgs("--host-ip", "::ffff:172.16.0.1"),
			wantErr: "--host-ip must be a valid IPv4 address: ::ffff:172.16.0.1",
		},
		{
			name:    "same guest and host ip",
			args:    requiredArgs("--guest-ip", "172.16.0.2", "--host-ip", "172.16.0.2"),
			wantErr: "--guest-ip and --host-ip must be different",
		},
		{
			name:    "guest and host ip on different slash24 networks",
			args:    requiredArgs("--guest-ip", "172.16.1.2", "--host-ip", "172.16.0.1"),
			wantErr: "--guest-ip and --host-ip must be in the same /24 network",
		},
		{
			name:    "blank tap prefix",
			args:    requiredArgs("--tap-prefix", " \t "),
			wantErr: "--tap-prefix must contain at least one ASCII letter or digit",
		},
		{
			name:    "tap prefix without usable characters",
			args:    requiredArgs("--tap-prefix", "---:://"),
			wantErr: "--tap-prefix must contain at least one ASCII letter or digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadFromArgs(tt.args, flag.ContinueOnError)
			if err == nil {
				t.Fatalf("loadFromArgs succeeded, want error containing %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("loadFromArgs error = %q, want to contain %q", err, tt.wantErr)
			}
		})
	}
}

func requiredArgs(extra ...string) []string {
	args := []string{"--kernel", "/kernel", "--rootfs", "/rootfs"}
	return append(args, extra...)
}

func assertConfigValue(t *testing.T, name, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("%s = %q, want %q", name, got, want)
	}
}
