package firecracker

import (
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
