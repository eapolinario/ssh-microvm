package config

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"ssh-microvm/internal/util"
)

const (
	AuthModeAutoEnroll = "auto-enroll"
	AuthModeKnownKeys  = "known-keys"
)

type Config struct {
	ListenAddr    string
	StateDir      string
	DBPath        string
	HostKeyPath   string
	AuthMode      string
	Firecracker   string
	KernelImage   string
	RootFS        string
	BootArgs      string
	VCPUCount     int
	MemMiB        int
	GracefulStopS int
	GuestUser     string
	GuestKeyPath  string
	GuestIP       string
	HostIP        string
	TapPrefix     string
}

func Load() (*Config, error) {
	id, err := util.RandomHex(6)
	if err != nil {
		return nil, err
	}

	defaultStateDir := fmt.Sprintf("/tmp/ssh-microvm-%s", id)

	cfg := &Config{}
	flag.StringVar(&cfg.ListenAddr, "listen", ":2222", "SSH listen address")
	flag.StringVar(&cfg.StateDir, "state-dir", defaultStateDir, "State directory for sockets, logs, and DB")
	flag.StringVar(&cfg.DBPath, "db-path", "", "SQLite database path (default: <state-dir>/db.sqlite)")
	flag.StringVar(&cfg.HostKeyPath, "host-key", "", "SSH host private key path (default: <state-dir>/ssh_host_ed25519)")
	flag.StringVar(&cfg.AuthMode, "auth-mode", AuthModeAutoEnroll, "Auth mode: auto-enroll or known-keys")
	flag.StringVar(&cfg.Firecracker, "firecracker", "firecracker", "Firecracker binary path")
	flag.StringVar(&cfg.KernelImage, "kernel", "", "Kernel image path")
	flag.StringVar(&cfg.RootFS, "rootfs", "", "Root filesystem image path (ext4)")
	flag.StringVar(&cfg.BootArgs, "boot-args", "console=ttyS0 reboot=k panic=1 pci=off", "Kernel boot args")
	flag.IntVar(&cfg.VCPUCount, "vcpu", 1, "VM vCPU count")
	flag.IntVar(&cfg.MemMiB, "mem", 512, "VM memory in MiB")
	flag.IntVar(&cfg.GracefulStopS, "grace-stop", 2, "Seconds to wait before force-killing Firecracker")
	flag.StringVar(&cfg.GuestUser, "guest-user", "root", "Guest SSH user")
	flag.StringVar(&cfg.GuestKeyPath, "guest-key", "artifacts/ubuntu.id_rsa", "Guest SSH private key path")
	flag.StringVar(&cfg.GuestIP, "guest-ip", "172.16.0.2", "Guest IP address")
	flag.StringVar(&cfg.HostIP, "host-ip", "172.16.0.1", "Host IP address for tap device")
	flag.StringVar(&cfg.TapPrefix, "tap-prefix", "tap", "Tap device name prefix")
	flag.Parse()

	if cfg.DBPath == "" {
		cfg.DBPath = filepath.Join(cfg.StateDir, "db.sqlite")
	}
	if cfg.HostKeyPath == "" {
		cfg.HostKeyPath = filepath.Join(cfg.StateDir, "ssh_host_ed25519")
	}
	if cfg.KernelImage == "" {
		return nil, errors.New("--kernel is required")
	}
	if cfg.RootFS == "" {
		return nil, errors.New("--rootfs is required")
	}
	if cfg.AuthMode != AuthModeAutoEnroll && cfg.AuthMode != AuthModeKnownKeys {
		return nil, fmt.Errorf("invalid --auth-mode: %s", cfg.AuthMode)
	}
	if cfg.VCPUCount <= 0 {
		return nil, errors.New("--vcpu must be > 0")
	}
	if cfg.MemMiB <= 0 {
		return nil, errors.New("--mem must be > 0")
	}
	if cfg.GuestUser == "" {
		return nil, errors.New("--guest-user must be set")
	}
	if cfg.GuestKeyPath == "" {
		return nil, errors.New("--guest-key must be set")
	}
	if cfg.GuestIP == "" || cfg.HostIP == "" {
		return nil, errors.New("--guest-ip and --host-ip must be set")
	}

	return cfg, nil
}
