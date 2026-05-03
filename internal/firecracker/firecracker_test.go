package firecracker

import (
	"context"
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
	vm := &VM{}

	if err := vm.Stop(context.Background(), time.Second); err != nil {
		t.Fatalf("Stop without process returned error: %v", err)
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

func TestTapNameForFitsLinuxInterfaceLimit(t *testing.T) {
	got := tapNameFor("tap-prefix-", "abcdef1234567890")
	if len(got) > 15 {
		t.Fatalf("tapNameFor() length = %d, want <= 15 (%q)", len(got), got)
	}
	if got != "tapprefixabcdef" {
		t.Fatalf("tapNameFor() = %q, want %q", got, "tapprefixabcdef")
	}
}
