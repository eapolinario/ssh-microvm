package util

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestRandomHexReturnsRequestedByteLengthAsHex(t *testing.T) {
	got, err := RandomHex(8)
	if err != nil {
		t.Fatalf("RandomHex: %v", err)
	}
	if len(got) != 16 {
		t.Fatalf("RandomHex(8) length = %d, want 16 hex characters", len(got))
	}
	if _, err := hex.DecodeString(got); err != nil {
		t.Fatalf("RandomHex(8) returned non-hex string %q: %v", got, err)
	}
}

func TestRandomHexAllowsZeroLength(t *testing.T) {
	got, err := RandomHex(0)
	if err != nil {
		t.Fatalf("RandomHex(0): %v", err)
	}
	if got != "" {
		t.Fatalf("RandomHex(0) = %q, want empty string", got)
	}
}

func TestRandomHexRejectsNegativeLength(t *testing.T) {
	got, err := RandomHex(-1)
	if err == nil {
		t.Fatalf("RandomHex(-1) = %q, want error", got)
	}
	if !strings.Contains(err.Error(), "non-negative") {
		t.Fatalf("RandomHex(-1) error = %q, want non-negative length error", err)
	}
}
