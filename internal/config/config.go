package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

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
	return loadFromArgs(os.Args[1:], flag.ExitOnError)
}

func loadFromArgs(args []string, errorHandling flag.ErrorHandling) (*Config, error) {
	id, err := util.RandomHex(6)
	if err != nil {
		return nil, err
	}

	defaultStateDir := fmt.Sprintf("/tmp/ssh-microvm-%s", id)

	cfg := &Config{}
	fs := flag.NewFlagSet("ssh-microvm", errorHandling)
	if errorHandling == flag.ContinueOnError {
		fs.SetOutput(io.Discard)
	}
	fs.StringVar(&cfg.ListenAddr, "listen", ":2222", "SSH listen address")
	fs.StringVar(&cfg.StateDir, "state-dir", defaultStateDir, "State directory for sockets, logs, and DB")
	fs.StringVar(&cfg.DBPath, "db-path", "", "SQLite database path (default: <state-dir>/db.sqlite)")
	fs.StringVar(&cfg.HostKeyPath, "host-key", "", "SSH host private key path (default: <state-dir>/ssh_host_ed25519)")
	fs.StringVar(&cfg.AuthMode, "auth-mode", AuthModeAutoEnroll, "Auth mode: auto-enroll or known-keys")
	fs.StringVar(&cfg.Firecracker, "firecracker", "firecracker", "Firecracker binary path")
	fs.StringVar(&cfg.KernelImage, "kernel", "", "Kernel image path")
	fs.StringVar(&cfg.RootFS, "rootfs", "", "Root filesystem image path (ext4)")
	fs.StringVar(&cfg.BootArgs, "boot-args", "console=ttyS0 reboot=k panic=1 pci=off", "Kernel boot args")
	fs.IntVar(&cfg.VCPUCount, "vcpu", 1, "VM vCPU count")
	fs.IntVar(&cfg.MemMiB, "mem", 512, "VM memory in MiB")
	fs.IntVar(&cfg.GracefulStopS, "grace-stop", 2, "Seconds to wait before force-killing Firecracker")
	fs.StringVar(&cfg.GuestUser, "guest-user", "root", "Guest SSH user")
	fs.StringVar(&cfg.GuestKeyPath, "guest-key", "artifacts/ubuntu.id_rsa", "Guest SSH private key path")
	fs.StringVar(&cfg.GuestIP, "guest-ip", "172.16.0.2", "Guest IP address")
	fs.StringVar(&cfg.HostIP, "host-ip", "172.16.0.1", "Host IP address for tap device")
	fs.StringVar(&cfg.TapPrefix, "tap-prefix", "tap", "Tap device name prefix")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if isBlank(cfg.StateDir) {
		return nil, errors.New("--state-dir must be set")
	}
	if hasSurroundingWhitespace(cfg.StateDir) {
		return nil, errors.New("--state-dir must not contain surrounding whitespace")
	}
	if isBlank(cfg.DBPath) {
		cfg.DBPath = filepath.Join(cfg.StateDir, "db.sqlite")
	} else if hasSurroundingWhitespace(cfg.DBPath) {
		return nil, errors.New("--db-path must not contain surrounding whitespace")
	}
	if isBlank(cfg.HostKeyPath) {
		cfg.HostKeyPath = filepath.Join(cfg.StateDir, "ssh_host_ed25519")
	} else if hasSurroundingWhitespace(cfg.HostKeyPath) {
		return nil, errors.New("--host-key must not contain surrounding whitespace")
	}
	if hasSurroundingWhitespace(cfg.ListenAddr) {
		return nil, errors.New("--listen must not contain surrounding whitespace")
	}
	if err := validateListenAddr(cfg.ListenAddr); err != nil {
		return nil, err
	}
	if isBlank(cfg.KernelImage) {
		return nil, errors.New("--kernel is required")
	}
	if hasSurroundingWhitespace(cfg.KernelImage) {
		return nil, errors.New("--kernel must not contain surrounding whitespace")
	}
	if isBlank(cfg.RootFS) {
		return nil, errors.New("--rootfs is required")
	}
	if hasSurroundingWhitespace(cfg.RootFS) {
		return nil, errors.New("--rootfs must not contain surrounding whitespace")
	}
	if isBlank(cfg.AuthMode) {
		return nil, errors.New("--auth-mode must be set")
	}
	if hasSurroundingWhitespace(cfg.AuthMode) {
		return nil, errors.New("--auth-mode must not contain surrounding whitespace")
	}
	if cfg.AuthMode != AuthModeAutoEnroll && cfg.AuthMode != AuthModeKnownKeys {
		return nil, fmt.Errorf("invalid --auth-mode: %s", cfg.AuthMode)
	}
	if !isBlank(cfg.BootArgs) && hasSurroundingWhitespace(cfg.BootArgs) {
		return nil, errors.New("--boot-args must not contain surrounding whitespace")
	}
	if isBlank(cfg.Firecracker) {
		return nil, errors.New("--firecracker must be set")
	}
	if hasSurroundingWhitespace(cfg.Firecracker) {
		return nil, errors.New("--firecracker must not contain surrounding whitespace")
	}
	if cfg.VCPUCount <= 0 {
		return nil, errors.New("--vcpu must be > 0")
	}
	if cfg.MemMiB <= 0 {
		return nil, errors.New("--mem must be > 0")
	}
	if cfg.GracefulStopS <= 0 {
		return nil, errors.New("--grace-stop must be > 0")
	}
	if isBlank(cfg.GuestUser) {
		return nil, errors.New("--guest-user must be set")
	}
	if hasSurroundingWhitespace(cfg.GuestUser) {
		return nil, errors.New("--guest-user must not contain surrounding whitespace")
	}
	if isBlank(cfg.GuestKeyPath) {
		return nil, errors.New("--guest-key must be set")
	}
	if hasSurroundingWhitespace(cfg.GuestKeyPath) {
		return nil, errors.New("--guest-key must not contain surrounding whitespace")
	}
	if isBlank(cfg.GuestIP) || isBlank(cfg.HostIP) {
		return nil, errors.New("--guest-ip and --host-ip must be set")
	}
	if hasSurroundingWhitespace(cfg.GuestIP) {
		return nil, errors.New("--guest-ip must not contain surrounding whitespace")
	}
	if hasSurroundingWhitespace(cfg.HostIP) {
		return nil, errors.New("--host-ip must not contain surrounding whitespace")
	}
	if !isIPv4(cfg.GuestIP) {
		return nil, fmt.Errorf("--guest-ip must be a valid IPv4 address: %s", cfg.GuestIP)
	}
	if !isIPv4(cfg.HostIP) {
		return nil, fmt.Errorf("--host-ip must be a valid IPv4 address: %s", cfg.HostIP)
	}
	if cfg.GuestIP == cfg.HostIP {
		return nil, errors.New("--guest-ip and --host-ip must be different")
	}
	if !sameIPv4Slash24(cfg.GuestIP, cfg.HostIP) {
		return nil, errors.New("--guest-ip and --host-ip must be in the same /24 network")
	}
	if !isBlank(cfg.TapPrefix) && hasSurroundingWhitespace(cfg.TapPrefix) {
		return nil, errors.New("--tap-prefix must not contain surrounding whitespace")
	}
	if err := validateTapPrefix(cfg.TapPrefix); err != nil {
		return nil, err
	}

	return cfg, nil
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func hasSurroundingWhitespace(value string) bool {
	return value != strings.TrimSpace(value)
}

func isIPv4(value string) bool {
	ip := net.ParseIP(value)
	ipv4 := ip.To4()
	return ipv4 != nil && value == ipv4.String()
}

func sameIPv4Slash24(a, b string) bool {
	ipA := net.ParseIP(a).To4()
	ipB := net.ParseIP(b).To4()
	if ipA == nil || ipB == nil {
		return false
	}
	return ipA[0] == ipB[0] && ipA[1] == ipB[1] && ipA[2] == ipB[2]
}

func validateTapPrefix(value string) error {
	hasUsableChar := false
	hasInvalidChar := false
	for _, r := range value {
		if isASCIIAlphaNumeric(r) {
			hasUsableChar = true
			continue
		}
		hasInvalidChar = true
	}
	if !hasUsableChar {
		return errors.New("--tap-prefix must contain at least one ASCII letter or digit")
	}
	if hasInvalidChar {
		return errors.New("--tap-prefix must contain only ASCII letters and digits")
	}
	return nil
}

func isASCIIAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func validateListenAddr(value string) error {
	if isBlank(value) {
		return errors.New("--listen must be set")
	}
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("--listen must be a valid TCP address: %s", value)
	}
	if port == "" {
		return errors.New("--listen port must be set")
	}
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("--listen port must be valid: %s", port)
	}
	if _, err := net.ResolveTCPAddr("tcp", value); err != nil {
		return fmt.Errorf("--listen must resolve to a valid TCP address: %s", value)
	}
	return nil
}
