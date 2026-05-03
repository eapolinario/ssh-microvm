package firecracker

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ssh-microvm/internal/config"
)

func TestStopNilVMDoesNotPanic(t *testing.T) {
	var vm *VM

	if err := vm.Stop(context.Background(), time.Second); err != nil {
		t.Fatalf("Stop(nil) returned error: %v", err)
	}
}

func TestStopWithoutProcessDoesNotPanic(t *testing.T) {
	logFile, err := os.Create(filepath.Join(t.TempDir(), "firecracker.log"))
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	vm := &VM{logFile: logFile}

	if err := vm.Stop(context.Background(), time.Second); err != nil {
		t.Fatalf("Stop without process returned error: %v", err)
	}
	if _, err := logFile.WriteString("after stop"); err == nil {
		t.Fatalf("Stop did not close the VM log file")
	}
}

func TestStopRejectsNilContext(t *testing.T) {
	vm := &VM{}

	if err := vm.Stop(nil, time.Second); err == nil || !strings.Contains(err.Error(), "context must be set") {
		t.Fatalf("Stop error = %v, want context validation error", err)
	}
}

func TestStopRejectsNonPositiveGracefulTimeout(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start test process: %v", err)
	}
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	vm := &VM{Cmd: cmd}

	err := vm.Stop(context.Background(), 0)
	if err == nil || !strings.Contains(err.Error(), "graceful shutdown timeout must be > 0") {
		t.Fatalf("Stop error = %v, want graceful timeout validation error", err)
	}
	if cmd.ProcessState != nil {
		t.Fatalf("Stop killed or waited for process despite validation failure")
	}
}

func TestStopTreatsGracefulSIGTERMAsSuccess(t *testing.T) {
	logFile, err := os.Create(filepath.Join(t.TempDir(), "firecracker.log"))
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	cmd := exec.Command("sleep", "10")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start test process: %v", err)
	}
	vm := &VM{Cmd: cmd, logFile: logFile}

	if err := vm.Stop(context.Background(), 5*time.Second); err != nil {
		t.Fatalf("Stop returned error for graceful SIGTERM shutdown: %v", err)
	}
	if _, err := logFile.WriteString("after stop"); err == nil {
		t.Fatalf("Stop did not close the VM log file")
	}
}

func TestStopAfterProcessAlreadyWaitedDoesNotReturnWaitError(t *testing.T) {
	cmd := exec.Command("sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start test process: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for test process: %v", err)
	}
	vm := &VM{Cmd: cmd}

	if err := vm.Stop(context.Background(), time.Second); err != nil {
		t.Fatalf("Stop returned error for already-waited process: %v", err)
	}
}

func TestStopReapsProcessAfterGracefulTimeoutKill(t *testing.T) {
	cmd := exec.Command("sh", "-c", "trap '' TERM; echo ready; while :; do :; done")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("create helper stdout pipe: %v", err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper process: %v", err)
	}
	t.Cleanup(func() {
		if cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
	if line, err := bufio.NewReader(stdout).ReadString('\n'); err != nil || line != "ready\n" {
		t.Fatalf("wait for helper readiness = %q, %v", line, err)
	}

	vm := &VM{Cmd: cmd}
	err = vm.Stop(context.Background(), 20*time.Millisecond)
	if err == nil || !strings.Contains(err.Error(), "firecracker shutdown timeout") {
		t.Fatalf("Stop error = %v, want shutdown timeout", err)
	}
	if cmd.ProcessState == nil {
		t.Fatalf("Stop killed the process but did not wait for it")
	}
}

func TestManagerStartRejectsNilDependencies(t *testing.T) {
	tests := []struct {
		name      string
		manager   *Manager
		useNilCtx bool
		wantErr   string
	}{
		{
			name:    "nil manager",
			manager: nil,
			wantErr: "firecracker manager must be set",
		},
		{
			name:    "nil config",
			manager: NewManager(nil),
			wantErr: "config must be set",
		},
		{
			name:      "nil context",
			manager:   NewManager(&config.Config{}),
			useNilCtx: true,
			wantErr:   "context must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.useNilCtx {
				ctx = nil
			}
			vm, err := tt.manager.Start(ctx)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Start error = %v, want containing %q", err, tt.wantErr)
			}
			if vm != nil {
				t.Fatalf("Start returned VM %v, want nil", vm)
			}
		})
	}
}

func TestManagerStartRejectsInvalidConfigBeforeSideEffects(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*config.Config)
		wantErr string
	}{
		{
			name: "blank state dir",
			mutate: func(cfg *config.Config) {
				cfg.StateDir = " \t "
			},
			wantErr: "state dir must be set",
		},
		{
			name: "state dir with surrounding whitespace",
			mutate: func(cfg *config.Config) {
				cfg.StateDir = " " + cfg.StateDir + " "
			},
			wantErr: "state dir must not contain surrounding whitespace",
		},
		{
			name: "blank firecracker binary",
			mutate: func(cfg *config.Config) {
				cfg.Firecracker = " \t "
			},
			wantErr: "firecracker binary must be set",
		},
		{
			name: "firecracker binary with surrounding whitespace",
			mutate: func(cfg *config.Config) {
				cfg.Firecracker = " /bin/firecracker "
			},
			wantErr: "firecracker binary must not contain surrounding whitespace",
		},
		{
			name: "blank kernel image",
			mutate: func(cfg *config.Config) {
				cfg.KernelImage = " \t "
			},
			wantErr: "kernel image must be set",
		},
		{
			name: "kernel image with surrounding whitespace",
			mutate: func(cfg *config.Config) {
				cfg.KernelImage = " /kernel "
			},
			wantErr: "kernel image must not contain surrounding whitespace",
		},
		{
			name: "blank rootfs",
			mutate: func(cfg *config.Config) {
				cfg.RootFS = " \t "
			},
			wantErr: "rootfs must be set",
		},
		{
			name: "rootfs with surrounding whitespace",
			mutate: func(cfg *config.Config) {
				cfg.RootFS = " /rootfs "
			},
			wantErr: "rootfs must not contain surrounding whitespace",
		},
		{
			name: "non-positive vcpus",
			mutate: func(cfg *config.Config) {
				cfg.VCPUCount = 0
			},
			wantErr: "vcpu count must be > 0",
		},
		{
			name: "non-positive memory",
			mutate: func(cfg *config.Config) {
				cfg.MemMiB = 0
			},
			wantErr: "memory must be > 0",
		},
		{
			name: "non-positive graceful shutdown timeout",
			mutate: func(cfg *config.Config) {
				cfg.GracefulStopS = 0
			},
			wantErr: "graceful shutdown timeout must be > 0",
		},
		{
			name: "blank guest ip",
			mutate: func(cfg *config.Config) {
				cfg.GuestIP = " \t "
			},
			wantErr: "guest IP must be set",
		},
		{
			name: "blank host ip",
			mutate: func(cfg *config.Config) {
				cfg.HostIP = " \t "
			},
			wantErr: "host IP must be set",
		},
		{
			name: "invalid guest ip",
			mutate: func(cfg *config.Config) {
				cfg.GuestIP = "not-an-ip"
			},
			wantErr: "guest IP must be a valid IPv4 address",
		},
		{
			name: "invalid host ip",
			mutate: func(cfg *config.Config) {
				cfg.HostIP = "not-an-ip"
			},
			wantErr: "host IP must be a valid IPv4 address",
		},
		{
			name: "same guest and host ip",
			mutate: func(cfg *config.Config) {
				cfg.GuestIP = "172.16.0.2"
				cfg.HostIP = "172.16.0.2"
			},
			wantErr: "guest IP and host IP must be different",
		},
		{
			name: "guest and host ip on different slash24 networks",
			mutate: func(cfg *config.Config) {
				cfg.GuestIP = "172.16.1.2"
				cfg.HostIP = "172.16.0.1"
			},
			wantErr: "guest IP and host IP must be in the same /24 network",
		},
		{
			name: "blank tap prefix",
			mutate: func(cfg *config.Config) {
				cfg.TapPrefix = " \t "
			},
			wantErr: "tap prefix must contain at least one ASCII letter or digit",
		},
		{
			name: "tap prefix without usable characters",
			mutate: func(cfg *config.Config) {
				cfg.TapPrefix = "---:://"
			},
			wantErr: "tap prefix must contain at least one ASCII letter or digit",
		},
		{
			name: "tap prefix with invalid characters",
			mutate: func(cfg *config.Config) {
				cfg.TapPrefix = "tap-bad"
			},
			wantErr: "tap prefix must contain only ASCII letters and digits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workDir := t.TempDir()
			t.Chdir(workDir)
			stateDir := filepath.Join(workDir, "state")
			cfg := validStartConfig(stateDir)
			tt.mutate(cfg)

			vm, err := NewManager(cfg).Start(context.Background())
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Start error = %v, want containing %q", err, tt.wantErr)
			}
			if vm != nil {
				t.Fatalf("Start returned VM %v, want nil", vm)
			}
			entries, readErr := os.ReadDir(workDir)
			if readErr != nil {
				t.Fatalf("read work dir: %v", readErr)
			}
			if len(entries) != 0 {
				t.Fatalf("Start created filesystem entries before validating config: %v", entries)
			}
		})
	}
}

func TestBuildBootArgsAddsGuestIPConfiguration(t *testing.T) {
	cfg := &config.Config{
		BootArgs: "console=ttyS0 reboot=k panic=1 pci=off",
		GuestIP:  "172.16.0.2",
		HostIP:   "172.16.0.1",
	}

	got := buildBootArgs(cfg)
	want := "console=ttyS0 reboot=k panic=1 pci=off ip=172.16.0.2::172.16.0.1:255.255.255.0::eth0:off"
	if got != want {
		t.Fatalf("buildBootArgs() = %q, want %q", got, want)
	}
}

func TestBuildBootArgsPreservesExistingIPConfiguration(t *testing.T) {
	cfg := &config.Config{
		BootArgs: "console=ttyS0 ip=dhcp",
		GuestIP:  "172.16.0.2",
		HostIP:   "172.16.0.1",
	}

	got := buildBootArgs(cfg)
	if got != cfg.BootArgs {
		t.Fatalf("buildBootArgs() = %q, want %q", got, cfg.BootArgs)
	}
}

func TestBuildBootArgsDoesNotTreatEmbeddedIPAsIPConfiguration(t *testing.T) {
	cfg := &config.Config{
		BootArgs: "console=ttyS0 fooip=bar",
		GuestIP:  "172.16.0.2",
		HostIP:   "172.16.0.1",
	}

	got := buildBootArgs(cfg)
	want := "console=ttyS0 fooip=bar ip=172.16.0.2::172.16.0.1:255.255.255.0::eth0:off"
	if got != want {
		t.Fatalf("buildBootArgs() = %q, want %q", got, want)
	}
}

func TestBuildBootArgsWithoutBaseArgsDoesNotAddLeadingWhitespace(t *testing.T) {
	cfg := &config.Config{
		BootArgs: " \t ",
		GuestIP:  "172.16.0.2",
		HostIP:   "172.16.0.1",
	}

	got := buildBootArgs(cfg)
	want := "ip=172.16.0.2::172.16.0.1:255.255.255.0::eth0:off"
	if got != want {
		t.Fatalf("buildBootArgs() = %q, want %q", got, want)
	}
}

func TestTapNameForFitsLinuxInterfaceLimit(t *testing.T) {
	got := tapNameFor("tap-prefix-", "abcdef123456")
	if len(got) > 15 {
		t.Fatalf("tapNameFor() length = %d, want <= 15 (%q)", len(got), got)
	}
	if got != "tapabcdef123456" {
		t.Fatalf("tapNameFor() = %q, want %q", got, "tapabcdef123456")
	}
}

func TestTapNameForKeepsVMIDWithLongPrefix(t *testing.T) {
	first := tapNameFor("very-long-prefix", "abcdef123456")
	second := tapNameFor("very-long-prefix", "123456abcdef")

	if first == second {
		t.Fatalf("tapNameFor() generated colliding names for distinct VM IDs: %q", first)
	}
	if !strings.HasSuffix(first, "abcdef123456") {
		t.Fatalf("tapNameFor() = %q, want VM ID suffix to avoid tap collisions", first)
	}
}

func TestTapNameForSanitizesInvalidCharacters(t *testing.T) {
	got := tapNameFor("tap/: bad-", "ab-cd/12:34")

	if got != "tapbadabcd1234" {
		t.Fatalf("tapNameFor() = %q, want sanitized name %q", got, "tapbadabcd1234")
	}
	if strings.ContainsAny(got, "/: -") {
		t.Fatalf("tapNameFor() = %q, want no invalid tap name characters", got)
	}
}

func TestTapNameForDefaultsEmptySanitizedPrefix(t *testing.T) {
	got := tapNameFor("---:://", "abcdef123456")

	if got != "tapabcdef123456" {
		t.Fatalf("tapNameFor() = %q, want default tap prefix after sanitization", got)
	}
}

func TestRandomMACIsDeterministicAndLocallyAdministered(t *testing.T) {
	got := randomMAC("vm-seed")
	if got != randomMAC("vm-seed") {
		t.Fatalf("randomMAC() was not deterministic for the same seed")
	}
	if got == randomMAC("other-seed") {
		t.Fatalf("randomMAC() returned the same address for different seeds: %q", got)
	}
	if !strings.HasPrefix(got, "02:") {
		t.Fatalf("randomMAC() = %q, want locally administered unicast prefix 02", got)
	}
}

func TestTapCommandHelpersRejectInvalidState(t *testing.T) {
	tests := []struct {
		name    string
		run     func() error
		wantErr string
	}{
		{
			name: "setupTap nil context",
			run: func() error {
				return setupTap(nil, "tap0", "172.16.0.1")
			},
			wantErr: "context must be set",
		},
		{
			name: "setupTap blank tap name",
			run: func() error {
				return setupTap(context.Background(), " \t ", "172.16.0.1")
			},
			wantErr: "tap name is empty",
		},
		{
			name: "setupTap invalid tap name characters",
			run: func() error {
				return setupTap(context.Background(), "tap bad", "172.16.0.1")
			},
			wantErr: "tap name must contain only ASCII letters and digits",
		},
		{
			name: "setupTap overlong tap name",
			run: func() error {
				return setupTap(context.Background(), "tap0123456789012", "172.16.0.1")
			},
			wantErr: "tap name must be <= 15 characters",
		},
		{
			name: "setupTap blank host IP",
			run: func() error {
				return setupTap(context.Background(), "tap0", " \t ")
			},
			wantErr: "host IP is empty",
		},
		{
			name: "setupTap invalid host IP",
			run: func() error {
				return setupTap(context.Background(), "tap0", "not-an-ip")
			},
			wantErr: "host IP must be a valid IPv4 address",
		},
		{
			name: "teardownTap nil context",
			run: func() error {
				return teardownTap(nil, "tap0")
			},
			wantErr: "context must be set",
		},
		{
			name: "teardownTap blank tap name",
			run: func() error {
				return teardownTap(context.Background(), " \t ")
			},
			wantErr: "tap name is empty",
		},
		{
			name: "teardownTap invalid tap name characters",
			run: func() error {
				return teardownTap(context.Background(), "tap/bad")
			},
			wantErr: "tap name must contain only ASCII letters and digits",
		},
		{
			name: "runCmd nil context",
			run: func() error {
				return runCmd(nil, "true")
			},
			wantErr: "context must be set",
		},
		{
			name: "runCmd blank command name",
			run: func() error {
				return runCmd(context.Background(), " \t ")
			},
			wantErr: "command name is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestWaitForSocket(t *testing.T) {
	socketPath := t.TempDir() + "/firecracker.sock"

	if err := waitForSocket(socketPath, 20*time.Millisecond); err == nil {
		t.Fatalf("waitForSocket() succeeded for missing socket")
	}

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen unix socket: %v", err)
	}
	defer ln.Close()

	if err := waitForSocket(socketPath, time.Second); err != nil {
		t.Fatalf("waitForSocket() existing socket returned error: %v", err)
	}
}

func TestWaitForSocketRejectsNonSocketPath(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "firecracker.sock")
	if err := os.WriteFile(socketPath, []byte("not a socket"), 0o600); err != nil {
		t.Fatalf("write placeholder socket path: %v", err)
	}

	if err := waitForSocket(socketPath, 20*time.Millisecond); err == nil {
		t.Fatalf("waitForSocket() succeeded for non-socket path")
	}
}

func TestWaitForSocketHonorsShortTimeout(t *testing.T) {
	socketPath := filepath.Join(t.TempDir(), "firecracker.sock")

	start := time.Now()
	if err := waitForSocket(socketPath, 20*time.Millisecond); err == nil {
		t.Fatalf("waitForSocket() succeeded for missing socket")
	}
	if elapsed := time.Since(start); elapsed > 45*time.Millisecond {
		t.Fatalf("waitForSocket() elapsed = %v, want short timeout honored", elapsed)
	}
}

func TestWaitForSocketRejectsInvalidState(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		timeout time.Duration
		wantErr string
	}{
		{
			name:    "blank socket path",
			path:    " \t ",
			timeout: time.Second,
			wantErr: "api socket path is empty",
		},
		{
			name:    "zero timeout",
			path:    filepath.Join(t.TempDir(), "firecracker.sock"),
			timeout: 0,
			wantErr: "api socket timeout must be positive",
		},
		{
			name:    "negative timeout",
			path:    filepath.Join(t.TempDir(), "firecracker.sock"),
			timeout: -time.Second,
			wantErr: "api socket timeout must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := waitForSocket(tt.path, tt.timeout)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestPutJSONOverUnixSocket(t *testing.T) {
	socketPath := t.TempDir() + "/firecracker.sock"
	requests := make(chan map[string]any, 1)
	server := newUnixHTTPServer(t, socketPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/machine-config" {
			t.Errorf("path = %s, want /machine-config", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", got)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		requests <- payload
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	err := putJSON(newUnixClient(socketPath), "/machine-config", map[string]any{
		"vcpu_count": 2,
		"smt":        false,
	})
	if err != nil {
		t.Fatalf("putJSON() returned error: %v", err)
	}

	got := <-requests
	if got["vcpu_count"] != float64(2) {
		t.Fatalf("vcpu_count = %#v, want 2", got["vcpu_count"])
	}
	if got["smt"] != false {
		t.Fatalf("smt = %#v, want false", got["smt"])
	}
}

func TestNewUnixClientSetsRequestTimeout(t *testing.T) {
	client := newUnixClient(filepath.Join(t.TempDir(), "firecracker.sock"))

	if client.Timeout != firecrackerAPITimeout {
		t.Fatalf("newUnixClient() timeout = %v, want %v", client.Timeout, firecrackerAPITimeout)
	}
}

func TestPutJSONRejectsInvalidState(t *testing.T) {
	tests := []struct {
		name    string
		client  *http.Client
		path    string
		wantErr string
	}{
		{
			name:    "nil client",
			client:  nil,
			path:    "/machine-config",
			wantErr: "http client must be set",
		},
		{
			name:    "blank path",
			client:  http.DefaultClient,
			path:    " \t ",
			wantErr: "firecracker api path is empty",
		},
		{
			name:    "relative path",
			client:  http.DefaultClient,
			path:    "machine-config",
			wantErr: "firecracker api path must start with /",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := putJSON(tt.client, tt.path, map[string]any{"ok": true})
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestPutJSONReturnsAPIErrorBody(t *testing.T) {
	socketPath := t.TempDir() + "/firecracker.sock"
	server := newUnixHTTPServer(t, socketPath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad config", http.StatusBadRequest)
	}))
	defer server.Close()

	err := putJSON(newUnixClient(socketPath), "/machine-config", map[string]any{"bad": true})
	if err == nil {
		t.Fatalf("putJSON() succeeded, want API error")
	}
	if !strings.Contains(err.Error(), "/machine-config") ||
		!strings.Contains(err.Error(), "400 Bad Request") ||
		!strings.Contains(err.Error(), "bad config") {
		t.Fatalf("putJSON() error = %q, want path, status, and response body", err)
	}
}

func TestPutJSONReturnsMarshalError(t *testing.T) {
	err := putJSON(http.DefaultClient, "/machine-config", map[string]any{"bad": make(chan int)})
	if err == nil {
		t.Fatalf("putJSON() succeeded, want marshal error")
	}
}

func newUnixHTTPServer(t *testing.T, socketPath string, handler http.Handler) *http.Server {
	t.Helper()

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen unix socket: %v", err)
	}

	server := &http.Server{Handler: handler}
	go func() {
		err := server.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("unix HTTP server: %v", err)
		}
	}()
	t.Cleanup(func() {
		_ = server.Close()
	})
	return server
}

func validStartConfig(stateDir string) *config.Config {
	return &config.Config{
		StateDir:      stateDir,
		Firecracker:   "firecracker",
		KernelImage:   "/kernel",
		RootFS:        "/rootfs",
		VCPUCount:     1,
		MemMiB:        512,
		GracefulStopS: 2,
		GuestIP:       "172.16.0.2",
		HostIP:        "172.16.0.1",
		TapPrefix:     "tap",
	}
}
