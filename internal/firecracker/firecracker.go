package firecracker

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"ssh-microvm/internal/config"
	"ssh-microvm/internal/util"
)

type Manager struct {
	cfg *config.Config
}

type VM struct {
	ID       string
	StateDir string
	APISock  string
	Cmd      *exec.Cmd
	TapName  string
	GuestIP  string
	logFile  *os.File
}

const firecrackerAPITimeout = 5 * time.Second

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) Start(ctx context.Context) (*VM, error) {
	if m == nil {
		return nil, errors.New("firecracker manager must be set")
	}
	if m.cfg == nil {
		return nil, errors.New("config must be set")
	}
	if ctx == nil {
		return nil, errors.New("context must be set")
	}
	if err := validateStartConfig(m.cfg); err != nil {
		return nil, err
	}
	id, err := util.RandomHex(6)
	if err != nil {
		return nil, err
	}
	vmDir := filepath.Join(m.cfg.StateDir, "vm-"+id)
	if err := os.MkdirAll(vmDir, 0o750); err != nil {
		return nil, err
	}

	apiSock := filepath.Join(vmDir, "firecracker.sock")
	logPath := filepath.Join(vmDir, "firecracker.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return nil, err
	}

	tapName := tapNameFor(m.cfg.TapPrefix, id)
	if err := setupTap(ctx, tapName, m.cfg.HostIP); err != nil {
		_ = logFile.Close()
		return nil, err
	}

	cmd := exec.CommandContext(ctx, m.cfg.Firecracker, "--api-sock", apiSock)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		_ = teardownTap(context.Background(), tapName)
		_ = logFile.Close()
		return nil, err
	}

	vm := &VM{
		ID:       id,
		StateDir: vmDir,
		APISock:  apiSock,
		Cmd:      cmd,
		TapName:  tapName,
		GuestIP:  m.cfg.GuestIP,
		logFile:  logFile,
	}

	if err := waitForSocket(apiSock, 2*time.Second); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}

	client := newUnixClient(apiSock)
	if err := putJSON(client, "/machine-config", map[string]any{
		"vcpu_count":   m.cfg.VCPUCount,
		"mem_size_mib": m.cfg.MemMiB,
		"smt":          false,
	}); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}
	if err := putJSON(client, "/network-interfaces/eth0", map[string]any{
		"iface_id":      "eth0",
		"host_dev_name": tapName,
		"guest_mac":     randomMAC(id),
	}); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}
	if err := putJSON(client, "/boot-source", map[string]any{
		"kernel_image_path": m.cfg.KernelImage,
		"boot_args":         buildBootArgs(m.cfg),
	}); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}
	if err := putJSON(client, "/drives/rootfs", map[string]any{
		"drive_id":       "rootfs",
		"path_on_host":   m.cfg.RootFS,
		"is_root_device": true,
		"is_read_only":   false,
	}); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}
	if err := putJSON(client, "/actions", map[string]any{
		"action_type": "InstanceStart",
	}); err != nil {
		_ = vm.Stop(context.Background(), 1*time.Second)
		return nil, err
	}

	return vm, nil
}

func validateStartConfig(cfg *config.Config) error {
	if strings.TrimSpace(cfg.StateDir) == "" {
		return errors.New("state dir must be set")
	}
	if strings.TrimSpace(cfg.Firecracker) == "" {
		return errors.New("firecracker binary must be set")
	}
	if strings.TrimSpace(cfg.KernelImage) == "" {
		return errors.New("kernel image must be set")
	}
	if strings.TrimSpace(cfg.RootFS) == "" {
		return errors.New("rootfs must be set")
	}
	if cfg.VCPUCount <= 0 {
		return errors.New("vcpu count must be > 0")
	}
	if cfg.MemMiB <= 0 {
		return errors.New("memory must be > 0")
	}
	if strings.TrimSpace(cfg.GuestIP) == "" {
		return errors.New("guest IP must be set")
	}
	if strings.TrimSpace(cfg.HostIP) == "" {
		return errors.New("host IP must be set")
	}
	if !isIPv4(cfg.GuestIP) {
		return fmt.Errorf("guest IP must be a valid IPv4 address: %s", cfg.GuestIP)
	}
	if !isIPv4(cfg.HostIP) {
		return fmt.Errorf("host IP must be a valid IPv4 address: %s", cfg.HostIP)
	}
	if cfg.GuestIP == cfg.HostIP {
		return errors.New("guest IP and host IP must be different")
	}
	if !sameIPv4Slash24(cfg.GuestIP, cfg.HostIP) {
		return errors.New("guest IP and host IP must be in the same /24 network")
	}
	sanitizedTapPrefix := sanitizeTapNamePart(cfg.TapPrefix)
	if sanitizedTapPrefix == "" {
		return errors.New("tap prefix must contain at least one ASCII letter or digit")
	}
	if cfg.TapPrefix != sanitizedTapPrefix {
		return errors.New("tap prefix must contain only ASCII letters and digits")
	}
	return nil
}

func (v *VM) Stop(ctx context.Context, graceful time.Duration) error {
	if v == nil {
		return nil
	}
	if ctx == nil {
		return errors.New("context must be set")
	}
	if graceful <= 0 {
		return errors.New("graceful shutdown timeout must be > 0")
	}
	defer v.closeLog()
	if v.Cmd == nil || v.Cmd.Process == nil {
		_ = teardownTap(context.Background(), v.TapName)
		return nil
	}
	if v.Cmd.ProcessState != nil {
		_ = teardownTap(context.Background(), v.TapName)
		return nil
	}
	_ = v.Cmd.Process.Signal(syscall.SIGTERM)

	done := make(chan error, 1)
	go func() {
		done <- v.Cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		_ = v.Cmd.Process.Kill()
		<-done
		_ = teardownTap(context.Background(), v.TapName)
		return ctx.Err()
	case <-time.After(graceful):
		_ = v.Cmd.Process.Kill()
		<-done
		_ = teardownTap(context.Background(), v.TapName)
		return errors.New("firecracker shutdown timeout")
	case err := <-done:
		_ = teardownTap(context.Background(), v.TapName)
		if wasSignal(err, syscall.SIGTERM) {
			return nil
		}
		return err
	}
}

func (v *VM) closeLog() {
	if v.logFile == nil {
		return
	}
	_ = v.logFile.Close()
	v.logFile = nil
}

func wasSignal(err error, signal syscall.Signal) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	return ok && status.Signaled() && status.Signal() == signal
}

func buildBootArgs(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	bootArgs := strings.TrimSpace(cfg.BootArgs)
	if hasIPArg(bootArgs) {
		return bootArgs
	}
	ipArg := fmt.Sprintf("ip=%s::%s:255.255.255.0::eth0:off", cfg.GuestIP, cfg.HostIP)
	if bootArgs == "" {
		return ipArg
	}
	return fmt.Sprintf("%s %s", bootArgs, ipArg)
}

func hasIPArg(args string) bool {
	for _, field := range strings.Fields(args) {
		if strings.HasPrefix(field, "ip=") {
			return true
		}
	}
	return false
}

func randomMAC(seed string) string {
	hash := sha1.Sum([]byte(seed))
	return fmt.Sprintf("02:%02x:%02x:%02x:%02x:%02x", hash[0], hash[1], hash[2], hash[3], hash[4])
}

func setupTap(ctx context.Context, tapName, hostIP string) error {
	if ctx == nil {
		return errors.New("context must be set")
	}
	if err := validateTapName(tapName); err != nil {
		return err
	}
	hostIP = strings.TrimSpace(hostIP)
	if hostIP == "" {
		return errors.New("host IP is empty")
	}
	if !isIPv4(hostIP) {
		return fmt.Errorf("host IP must be a valid IPv4 address: %s", hostIP)
	}
	if err := runCmd(ctx, "sudo", "ip", "tuntap", "add", "dev", tapName, "mode", "tap"); err != nil {
		return err
	}
	if err := runCmd(ctx, "sudo", "ip", "addr", "add", hostIP+"/24", "dev", tapName); err != nil {
		_ = runCmd(context.Background(), "sudo", "ip", "link", "del", tapName)
		return err
	}
	if err := runCmd(ctx, "sudo", "ip", "link", "set", tapName, "up"); err != nil {
		_ = runCmd(context.Background(), "sudo", "ip", "link", "del", tapName)
		return err
	}
	return nil
}

func teardownTap(ctx context.Context, tapName string) error {
	if ctx == nil {
		return errors.New("context must be set")
	}
	if err := validateTapName(tapName); err != nil {
		return err
	}
	return runCmd(ctx, "sudo", "ip", "link", "del", tapName)
}

func validateTapName(tapName string) error {
	if strings.TrimSpace(tapName) == "" {
		return errors.New("tap name is empty")
	}
	if len(tapName) > 15 {
		return errors.New("tap name must be <= 15 characters")
	}
	if tapName != sanitizeTapNamePart(tapName) {
		return errors.New("tap name must contain only ASCII letters and digits")
	}
	return nil
}

func runCmd(ctx context.Context, name string, args ...string) error {
	if ctx == nil {
		return errors.New("context must be set")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("command name is empty")
	}
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %v failed: %w: %s", name, args, err, string(out))
	}
	return nil
}

func tapNameFor(prefix, id string) string {
	prefix = sanitizeTapNamePart(prefix)
	if prefix == "" {
		prefix = "tap"
	}
	id = sanitizeTapNamePart(id)
	name := prefix + id
	if len(name) > 15 {
		if len(id) >= 15 {
			return id[:15]
		}
		return prefix[:15-len(id)] + id
	}
	return name
}

func sanitizeTapNamePart(value string) string {
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		}
	}
	return b.String()
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

func waitForSocket(path string, timeout time.Duration) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("api socket path is empty")
	}
	if timeout <= 0 {
		return errors.New("api socket timeout must be positive")
	}

	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		info, err := os.Stat(path)
		if err != nil {
			lastErr = err
		} else if info.Mode()&os.ModeSocket == 0 {
			lastErr = fmt.Errorf("%s exists but is not a unix socket", path)
		} else {
			dialTimeout := time.Until(deadline)
			if dialTimeout > 50*time.Millisecond {
				dialTimeout = 50 * time.Millisecond
			}
			conn, err := net.DialTimeout("unix", path, dialTimeout)
			if err == nil {
				_ = conn.Close()
				return nil
			}
			lastErr = err
		}
		sleepFor := time.Until(deadline)
		if sleepFor > 50*time.Millisecond {
			sleepFor = 50 * time.Millisecond
		}
		if sleepFor > 0 {
			time.Sleep(sleepFor)
		}
	}
	if lastErr != nil {
		return fmt.Errorf("timeout waiting for api socket %s: %w", path, lastErr)
	}
	return fmt.Errorf("timeout waiting for api socket: %s", path)
}

func newUnixClient(sock string) *http.Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", sock)
		},
	}
	return &http.Client{Transport: transport, Timeout: firecrackerAPITimeout}
}

func putJSON(client *http.Client, path string, payload any) error {
	if client == nil {
		return errors.New("http client must be set")
	}
	if strings.TrimSpace(path) == "" {
		return errors.New("firecracker api path is empty")
	}
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("firecracker api path must start with /: %s", path)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", "http://unix"+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body [512]byte
		n, _ := resp.Body.Read(body[:])
		return fmt.Errorf("firecracker api %s failed: %s: %s", path, resp.Status, bytes.TrimSpace(body[:n]))
	}
	return nil
}
