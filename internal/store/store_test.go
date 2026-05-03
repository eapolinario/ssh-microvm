package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testAuthorizedKey       = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIPh2vDLbN/0Bu93v5NvdlRQ7WOpknAUgJ0l1ofhOYTpf"
	testKeyFingerprint      = "SHA256:UecLtXI8mKCwPSeFNoPFanZ4gYYgIREcsLQBav+pqAg"
	testOtherAuthorizedKey  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBp154TtDlncesM4JmySe6K3G8KtiAt6PYdPUiWGLFr6"
	testOtherKeyFingerprint = "SHA256:dGAIKjPAvDNRn2eUYFTajUJGNHzLaUHgJsFOfgFzlyI"
)

func testSHA256Fingerprint(suffix byte) string {
	return "SHA256:" + strings.Repeat("A", 42) + string(suffix)
}

func TestEnsureSchemaIsIdempotent(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.EnsureSchema(ctx); err != nil {
		t.Fatalf("second EnsureSchema: %v", err)
	}

	var migrationCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 1")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 1 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 2")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 2: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 2 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 3")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 3: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 3 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 4")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 4: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 4 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 5")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 5: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 5 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 6")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 6: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 6 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 7")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 7: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 7 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 8")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 8: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 8 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 9")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 9: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 9 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 10")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 10: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 10 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 11")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 11: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 11 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 12")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 12: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 12 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 13")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 13: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 13 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 14")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 14: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 14 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 15")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 15: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 15 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 16")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 16: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 16 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 17")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 17: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 17 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 18")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 18: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 18 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 19")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 19: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 19 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 20")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 20: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 20 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 21")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 21: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 21 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 22")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 22: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 22 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 23")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 23: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 23 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 24")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 24: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 24 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 25")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 25: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 25 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 26")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 26: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 26 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 27")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 27: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 27 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 28")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 28: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 28 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 29")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 29: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 29 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 30")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 30: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 30 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 31")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 31: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 31 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 32")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 32: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 32 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 33")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 33: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 33 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 34")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 34: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 34 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 35")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 35: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 35 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 36")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 36: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 36 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 37")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 37: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 37 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 38")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 38: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 38 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 39")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 39: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 39 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 40")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 40: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 40 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 41")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 41: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 41 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 42")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 42: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 42 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 43")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 43: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 43 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 44")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 44: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 44 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 45")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 45: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 45 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 46")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 46: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 46 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 47")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 47: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 47 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 48")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 48: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 48 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 49")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 49: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 49 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 50")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 50: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 50 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 51")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 51: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 51 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 52")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 52: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 52 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 53")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 53: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 53 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 54")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 54: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 54 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 55")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 55: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 55 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 56")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 56: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 56 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 57")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 57: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 57 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 58")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 58: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 58 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 59")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 59: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 59 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 60")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 60: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 60 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 61")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 61: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 61 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 62")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 62: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 62 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 63")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 63: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 63 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 64")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 64: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 64 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 65")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 65: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 65 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 66")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 66: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 66 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 67")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 67: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 67 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 68")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 68: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 68 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 69")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 69: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 69 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 70")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 70: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 70 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 71")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 71: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 71 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 72")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 72: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 72 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 73")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 73: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 73 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 74")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 74: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 74 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 75")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 75: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 75 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 76")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 76: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 76 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 77")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 77: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 77 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 78")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 78: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 78 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 79")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 79: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 79 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 80")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 80: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 80 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 81")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 81: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 81 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 82")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 82: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 82 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 83")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 83: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 83 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 84")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 84: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 84 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 85")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 85: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 85 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 86")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 86: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 86 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 87")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 87: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 87 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 88")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 88: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 88 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 89")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 89: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 89 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 90")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 90: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 90 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 91")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 91: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 91 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 92")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 92: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 92 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 93")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 93: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 93 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 94")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 94: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 94 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 95")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 95: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 95 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 96")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 96: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 96 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 97")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 97: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 97 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 98")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 98: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 98 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 99")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 99: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 99 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 100")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 100: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 100 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 101")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 101: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 101 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 102")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 102: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 102 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 103")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 103: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 103 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 104")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 104: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 104 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 105")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 105: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 105 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 106")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 106: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 106 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 107")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 107: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 107 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 108")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 108: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 108 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 109")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 109: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 109 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 110")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 110: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 110 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 111")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 111: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 111 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 112")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 112: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 112 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 113")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 113: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 113 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 114")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 114: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 114 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 115")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 115: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 115 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 116")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 116: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 116 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 117")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 117: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 117 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 118")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 118: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 118 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 119")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 119: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 119 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 120")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 120: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 120 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 121")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 121: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 121 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 122")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 122: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 122 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 123")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 123: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 123 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 124")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 124: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 124 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 125")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 125: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 125 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 126")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 126: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 126 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 127")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 127: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 127 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 128")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 128: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 128 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 129")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 129: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 129 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 130")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 130: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 130 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 131")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 131: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 131 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 132")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 132: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 132 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 133")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 133: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 133 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 134")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 134: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 134 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 135")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 135: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 135 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 136")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 136: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 136 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 137")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 137: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 137 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 138")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 138: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 138 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 139")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 139: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 139 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 140")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 140: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 140 migration count = %d, want 1", migrationCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = 141")
	if err := row.Scan(&migrationCount); err != nil {
		t.Fatalf("query schema_migrations version 141: %v", err)
	}
	if migrationCount != 1 {
		t.Fatalf("version 141 migration count = %d, want 1", migrationCount)
	}

	for _, table := range []string{"users", "keys", "sessions", "vms", "audit_events"} {
		var count int
		row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?", table)
		if err := row.Scan(&count); err != nil {
			t.Fatalf("query table %s: %v", table, err)
		}
		if count != 1 {
			t.Fatalf("table %s count = %d, want 1", table, count)
		}
	}
}

func TestEnsureSchemaEnforcesUserIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank user ID", userID: " \t "},
		{name: "padded user ID", userID: " user-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.name, now(), now()); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	userID := "user-1"
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", userID, "alice", now(), now()); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank user ID", userID: "\n\t"},
		{name: "padded user ID", userID: "\tuser-2\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE users SET id = ? WHERE id = ?", tt.userID, userID); err == nil {
				t.Fatalf("updated user to %s, want trigger error", tt.name)
			}
		})
	}

	var gotUserID string
	row := st.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = ?", userID)
	if err := row.Scan(&gotUserID); err != nil {
		t.Fatalf("query user ID: %v", err)
	}
	if gotUserID != userID {
		t.Fatalf("user id = %q, want %q", gotUserID, userID)
	}
}

func TestEnsureSchemaEnforcesSessionIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
	}{
		{name: "blank session ID", sessionID: " \t "},
		{name: "padded session ID", sessionID: " session-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name      string
		sessionID string
	}{
		{name: "blank session ID", sessionID: "\n\t"},
		{name: "padded session ID", sessionID: "\tsession-2\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET id = ? WHERE id = ?", tt.sessionID, sessionID); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	var gotSessionID string
	row := st.db.QueryRowContext(ctx, "SELECT id FROM sessions WHERE id = ?", sessionID)
	if err := row.Scan(&gotSessionID); err != nil {
		t.Fatalf("query session ID: %v", err)
	}
	if gotSessionID != sessionID {
		t.Fatalf("session id = %q, want %q", gotSessionID, sessionID)
	}
}

func TestEnsureSchemaEnforcesSessionUserIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank user ID", userID: " \t "},
		{name: "padded user ID", userID: " " + userID + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "bad-"+strings.ReplaceAll(tt.name, " ", "-"), tt.userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank user ID", userID: "\n\t"},
		{name: "padded user ID", userID: "\t" + userID + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET user_id = ? WHERE id = ?", tt.userID, sessionID); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	var gotUserID string
	row := st.db.QueryRowContext(ctx, "SELECT user_id FROM sessions WHERE id = ?", sessionID)
	if err := row.Scan(&gotUserID); err != nil {
		t.Fatalf("query session user_id: %v", err)
	}
	if gotUserID != userID {
		t.Fatalf("session user_id = %q, want %q", gotUserID, userID)
	}
}

func TestEnsureSchemaEnforcesSessionKeyFingerprintValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name           string
		keyFingerprint string
	}{
		{name: "blank key fingerprint", keyFingerprint: " \t "},
		{name: "padded key fingerprint", keyFingerprint: " " + testKeyFingerprint + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "bad-"+strings.ReplaceAll(tt.name, " ", "-"), userID, tt.keyFingerprint, "127.0.0.1:2222", now(), "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name           string
		keyFingerprint string
	}{
		{name: "blank key fingerprint", keyFingerprint: "\n\t"},
		{name: "padded key fingerprint", keyFingerprint: "\t" + testKeyFingerprint + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET key_fingerprint = ? WHERE id = ?", tt.keyFingerprint, sessionID); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	var gotKeyFingerprint string
	row := st.db.QueryRowContext(ctx, "SELECT key_fingerprint FROM sessions WHERE id = ?", sessionID)
	if err := row.Scan(&gotKeyFingerprint); err != nil {
		t.Fatalf("query session key_fingerprint: %v", err)
	}
	if gotKeyFingerprint != testKeyFingerprint {
		t.Fatalf("session key_fingerprint = %q, want %q", gotKeyFingerprint, testKeyFingerprint)
	}
}

func TestEnsureSchemaEnforcesVMIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert session fixture: %v", err)
	}

	for _, tt := range []struct {
		name string
		vmID string
	}{
		{name: "blank VM ID", vmID: " \t "},
		{name: "padded VM ID", vmID: " vm-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, tt.vmID, sessionID, filepath.Join(t.TempDir(), tt.name), 1234, now()); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	vmID := "vm-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, vmID, sessionID, filepath.Join(t.TempDir(), vmID), 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	for _, tt := range []struct {
		name string
		vmID string
	}{
		{name: "blank VM ID", vmID: "\n\t"},
		{name: "padded VM ID", vmID: "\tvm-2\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE vms SET id = ? WHERE id = ?", tt.vmID, vmID); err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
		})
	}

	var gotVMID string
	row := st.db.QueryRowContext(ctx, "SELECT id FROM vms WHERE id = ?", vmID)
	if err := row.Scan(&gotVMID); err != nil {
		t.Fatalf("query VM ID: %v", err)
	}
	if gotVMID != vmID {
		t.Fatalf("VM id = %q, want %q", gotVMID, vmID)
	}
}

func TestEnsureSchemaEnforcesVMSessionIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert session fixture: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
	}{
		{name: "blank session ID", sessionID: " \t "},
		{name: "padded session ID", sessionID: " " + sessionID + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "bad-"+strings.ReplaceAll(tt.name, " ", "-"), tt.sessionID, filepath.Join(t.TempDir(), tt.name), 1234, now())
			if err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "VM session ID must be set and not contain surrounding whitespace") {
				t.Fatalf("insert VM with %s error = %v, want VM session ID trigger error", tt.name, err)
			}
		})
	}

	vmID := "vm-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, vmID, sessionID, filepath.Join(t.TempDir(), vmID), 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	for _, tt := range []struct {
		name      string
		sessionID string
	}{
		{name: "blank session ID", sessionID: "\n\t"},
		{name: "padded session ID", sessionID: "\t" + sessionID + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE vms SET session_id = ? WHERE id = ?", tt.sessionID, vmID)
			if err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "VM session ID must be set and not contain surrounding whitespace") {
				t.Fatalf("update VM to %s error = %v, want VM session ID trigger error", tt.name, err)
			}
		})
	}

	var gotSessionID string
	row := st.db.QueryRowContext(ctx, "SELECT session_id FROM vms WHERE id = ?", vmID)
	if err := row.Scan(&gotSessionID); err != nil {
		t.Fatalf("query VM session_id: %v", err)
	}
	if gotSessionID != sessionID {
		t.Fatalf("VM session_id = %q, want %q", gotSessionID, sessionID)
	}
}

func TestEnsureSchemaEnforcesSessionVMIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name string
		vmID string
	}{
		{name: "blank VM ID", vmID: " \t "},
		{name: "padded VM ID", vmID: " vm-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status, vm_id)
VALUES(?, ?, ?, ?, ?, ?, ?)`, "bad-"+strings.ReplaceAll(tt.name, " ", "-"), userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active", tt.vmID)
			if err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "session VM ID must be set and not contain surrounding whitespace") {
				t.Fatalf("insert session with %s error = %v, want session VM ID trigger error", tt.name, err)
			}
		})
	}

	sessionID := "session-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	vmID := "vm-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, vmID, sessionID, filepath.Join(t.TempDir(), vmID), 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", vmID, sessionID); err != nil {
		t.Fatalf("attach valid VM: %v", err)
	}

	for _, tt := range []struct {
		name string
		vmID string
	}{
		{name: "blank VM ID", vmID: "\n\t"},
		{name: "padded VM ID", vmID: "\t" + vmID + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", tt.vmID, sessionID)
			if err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "session VM ID must be set and not contain surrounding whitespace") {
				t.Fatalf("update session to %s error = %v, want session VM ID trigger error", tt.name, err)
			}
		})
	}

	var gotVMID sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", sessionID)
	if err := row.Scan(&gotVMID); err != nil {
		t.Fatalf("query session vm_id: %v", err)
	}
	if !gotVMID.Valid || gotVMID.String != vmID {
		t.Fatalf("session vm_id = %v, want %q", gotVMID, vmID)
	}
}

func TestEnsureSchemaEnforcesKeyUserIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID := "user-1"
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", userID, "alice", now(), now()); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}

	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank key user ID", userID: " \t "},
		{name: "padded key user ID", userID: " user-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, tt.userID, testAuthorizedKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "key user ID must be set and not contain surrounding whitespace") {
				t.Fatalf("insert key with %s error = %v, want key user ID trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, userID, testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range []struct {
		name   string
		userID string
	}{
		{name: "blank key user ID", userID: "\n\t"},
		{name: "padded key user ID", userID: "\t" + userID + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET user_id = ? WHERE fingerprint = ?", tt.userID, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "key user ID must be set and not contain surrounding whitespace") {
				t.Fatalf("update key to %s error = %v, want key user ID trigger error", tt.name, err)
			}
		})
	}

	var gotUserID string
	row := st.db.QueryRowContext(ctx, "SELECT user_id FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotUserID); err != nil {
		t.Fatalf("query key user ID: %v", err)
	}
	if gotUserID != userID {
		t.Fatalf("key user_id = %q, want %q", gotUserID, userID)
	}
}

func TestEnsureSchemaEnforcesAuditEventIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name    string
		eventID string
	}{
		{name: "blank audit event ID", eventID: " \t "},
		{name: "padded audit event ID", eventID: " audit-1 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, "test.audit", `{"ok":true}`, now()); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	eventID := "audit-1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, eventID, "test.audit", `{"ok":true}`, now()); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	for _, tt := range []struct {
		name    string
		eventID string
	}{
		{name: "blank audit event ID", eventID: "\n\t"},
		{name: "padded audit event ID", eventID: "\taudit-2\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET id = ? WHERE id = ?", tt.eventID, eventID); err == nil {
				t.Fatalf("updated audit event to %s, want trigger error", tt.name)
			}
		})
	}

	var gotEventID string
	row := st.db.QueryRowContext(ctx, "SELECT id FROM audit_events WHERE id = ?", eventID)
	if err := row.Scan(&gotEventID); err != nil {
		t.Fatalf("query audit event ID: %v", err)
	}
	if gotEventID != eventID {
		t.Fatalf("audit event id = %q, want %q", gotEventID, eventID)
	}
}

func TestEnsureSchemaEnforcesUsernames(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name     string
		userID   string
		username string
	}{
		{name: "blank username", userID: "bad-blank-username", username: " \t "},
		{name: "padded username", userID: "bad-padded-username", username: " alice "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.username, now(), now()); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	username := "alice"
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", username, now(), now()); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	for _, tt := range []struct {
		name     string
		username string
	}{
		{name: "blank username", username: "\n\t"},
		{name: "padded username", username: "\tbob\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE users SET username = ? WHERE id = ?", tt.username, "user-1"); err == nil {
				t.Fatalf("updated user to %s, want trigger error", tt.name)
			}
		})
	}

	var gotUsername string
	row := st.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = ?", "user-1")
	if err := row.Scan(&gotUsername); err != nil {
		t.Fatalf("query username: %v", err)
	}
	if gotUsername != username {
		t.Fatalf("username = %q, want %q", gotUsername, username)
	}
}

func TestEnsureSchemaEnforcesUserCreatedAtValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		userID    string
		createdAt string
	}{
		{name: "blank creation time", userID: "bad-blank-created-at", createdAt: " \t "},
		{name: "padded creation time", userID: "bad-padded-created-at", createdAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.userID, tt.createdAt, now()); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	createdAt := now()
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", createdAt, now()); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	for _, tt := range []struct {
		name      string
		createdAt string
	}{
		{name: "blank creation time", createdAt: "\n\t"},
		{name: "padded creation time", createdAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE users SET created_at = ? WHERE id = ?", tt.createdAt, "user-1"); err == nil {
				t.Fatalf("updated user to %s, want trigger error", tt.name)
			}
		})
	}

	var gotCreatedAt string
	row := st.db.QueryRowContext(ctx, "SELECT created_at FROM users WHERE id = ?", "user-1")
	if err := row.Scan(&gotCreatedAt); err != nil {
		t.Fatalf("query user created_at: %v", err)
	}
	if gotCreatedAt != createdAt {
		t.Fatalf("user created_at = %q, want %q", gotCreatedAt, createdAt)
	}
}

func TestEnsureSchemaEnforcesUserCreatedAtFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		userID    string
		createdAt string
	}{
		{name: "space separated creation time", userID: "bad-space-created-at", createdAt: "2026-05-03 14:03:32"},
		{name: "missing timezone creation time", userID: "bad-missing-timezone-created-at", createdAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.userID, tt.createdAt, now()); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	createdAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", createdAt, now()); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE users SET created_at = ? WHERE id = ?", "2026-05-03 14:03:32", "user-1"); err == nil {
		t.Fatalf("updated user to space separated creation time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE users SET created_at = ? WHERE id = ?", validOffsetTime, "user-1"); err != nil {
		t.Fatalf("updated user to valid offset creation time: %v", err)
	}

	var gotCreatedAt string
	row := st.db.QueryRowContext(ctx, "SELECT created_at FROM users WHERE id = ?", "user-1")
	if err := row.Scan(&gotCreatedAt); err != nil {
		t.Fatalf("query user created_at: %v", err)
	}
	if gotCreatedAt != validOffsetTime {
		t.Fatalf("user created_at = %q, want %q", gotCreatedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesUserLastSeenAtValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name       string
		userID     string
		lastSeenAt string
	}{
		{name: "blank last seen time", userID: "bad-blank-last-seen-at", lastSeenAt: " \t "},
		{name: "padded last seen time", userID: "bad-padded-last-seen-at", lastSeenAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.userID, now(), tt.lastSeenAt); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	lastSeenAt := now()
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), lastSeenAt); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	for _, tt := range []struct {
		name       string
		lastSeenAt string
	}{
		{name: "blank last seen time", lastSeenAt: "\n\t"},
		{name: "padded last seen time", lastSeenAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE users SET last_seen_at = ? WHERE id = ?", tt.lastSeenAt, "user-1"); err == nil {
				t.Fatalf("updated user to %s, want trigger error", tt.name)
			}
		})
	}

	var gotLastSeenAt string
	row := st.db.QueryRowContext(ctx, "SELECT last_seen_at FROM users WHERE id = ?", "user-1")
	if err := row.Scan(&gotLastSeenAt); err != nil {
		t.Fatalf("query user last_seen_at: %v", err)
	}
	if gotLastSeenAt != lastSeenAt {
		t.Fatalf("user last_seen_at = %q, want %q", gotLastSeenAt, lastSeenAt)
	}
}

func TestEnsureSchemaEnforcesUserLastSeenAtFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name       string
		userID     string
		lastSeenAt string
	}{
		{name: "space separated last seen time", userID: "bad-space-last-seen-at", lastSeenAt: "2026-05-03 14:03:32"},
		{name: "missing timezone last seen time", userID: "bad-missing-timezone-last-seen-at", lastSeenAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", tt.userID, tt.userID, now(), tt.lastSeenAt); err == nil {
				t.Fatalf("inserted user with %s, want trigger error", tt.name)
			}
		})
	}

	lastSeenAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), lastSeenAt); err != nil {
		t.Fatalf("insert valid user: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE users SET last_seen_at = ? WHERE id = ?", "2026-05-03 14:03:32", "user-1"); err == nil {
		t.Fatalf("updated user to space separated last seen time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE users SET last_seen_at = ? WHERE id = ?", validOffsetTime, "user-1"); err != nil {
		t.Fatalf("updated user to valid offset last seen time: %v", err)
	}

	var gotLastSeenAt string
	row := st.db.QueryRowContext(ctx, "SELECT last_seen_at FROM users WHERE id = ?", "user-1")
	if err := row.Scan(&gotLastSeenAt); err != nil {
		t.Fatalf("query user last_seen_at: %v", err)
	}
	if gotLastSeenAt != validOffsetTime {
		t.Fatalf("user last_seen_at = %q, want %q", gotLastSeenAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesKeyFingerprintValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
	}{
		{name: "blank fingerprint", fingerprint: " \t "},
		{name: "padded fingerprint", fingerprint: " " + testKeyFingerprint + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, now(), now()); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range []struct {
		name        string
		fingerprint string
	}{
		{name: "blank fingerprint", fingerprint: "\n\t"},
		{name: "padded fingerprint", fingerprint: "\t" + testOtherKeyFingerprint + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE keys SET fingerprint = ? WHERE fingerprint = ?", tt.fingerprint, testKeyFingerprint); err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
		})
	}

	var gotFingerprint string
	row := st.db.QueryRowContext(ctx, "SELECT fingerprint FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotFingerprint); err != nil {
		t.Fatalf("query key fingerprint: %v", err)
	}
	if gotFingerprint != testKeyFingerprint {
		t.Fatalf("key fingerprint = %q, want %q", gotFingerprint, testKeyFingerprint)
	}
}

func TestEnsureSchemaEnforcesKeyFingerprintSHA256Shape(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
	}{
		{name: "missing SHA256 prefix", fingerprint: "MD5:UecLtXI8mKCwPSeFNoPFanZ4gYYgIREcsLQBav+pqAg"},
		{name: "short SHA256 fingerprint", fingerprint: "SHA256:short"},
		{name: "invalid SHA256 character", fingerprint: "SHA256:UecLtXI8mKCwPSeFNoPFanZ4gYYgIREcsLQBav+pqA="},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "key fingerprint must be a SHA256 fingerprint") {
				t.Fatalf("insert key with %s error = %v, want SHA256 fingerprint trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET fingerprint = ? WHERE fingerprint = ?", testOtherKeyFingerprint, testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid SHA256 fingerprint: %v", err)
	}
	for _, tt := range []struct {
		name        string
		fingerprint string
	}{
		{name: "missing SHA256 prefix", fingerprint: "MD5:dGAIKjPAvDNRn2eUYFTajUJGNHzLaUHgJsFOfgFzlyI"},
		{name: "short SHA256 fingerprint", fingerprint: "SHA256:short"},
		{name: "invalid SHA256 character", fingerprint: "SHA256:dGAIKjPAvDNRn2eUYFTajUJGNHzLaUHgJsFOfgFzly="},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET fingerprint = ? WHERE fingerprint = ?", tt.fingerprint, testOtherKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "key fingerprint must be a SHA256 fingerprint") {
				t.Fatalf("update key to %s error = %v, want SHA256 fingerprint trigger error", tt.name, err)
			}
		})
	}

	var gotFingerprint string
	row := st.db.QueryRowContext(ctx, "SELECT fingerprint FROM keys WHERE fingerprint = ?", testOtherKeyFingerprint)
	if err := row.Scan(&gotFingerprint); err != nil {
		t.Fatalf("query key fingerprint: %v", err)
	}
	if gotFingerprint != testOtherKeyFingerprint {
		t.Fatalf("key fingerprint = %q, want %q", gotFingerprint, testOtherKeyFingerprint)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for i, tt := range []struct {
		name      string
		publicKey string
	}{
		{name: "blank public key", publicKey: " \t "},
		{name: "padded public key", publicKey: " " + testAuthorizedKey + "\n"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now()); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range []struct {
		name      string
		publicKey string
	}{
		{name: "blank public key", publicKey: "\n\t"},
		{name: "padded public key", publicKey: "\t" + testOtherAuthorizedKey + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint); err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != testAuthorizedKey {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, testAuthorizedKey)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeySingleLine(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "line-feed separated public keys", publicKey: testAuthorizedKey + "\n" + testOtherAuthorizedKey},
		{name: "carriage-return separated public keys", publicKey: testAuthorizedKey + "\r" + testOtherAuthorizedKey},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must contain exactly one authorized key") {
				t.Fatalf("insert key with %s error = %v, want single authorized key trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must contain exactly one authorized key") {
				t.Fatalf("update key to %s error = %v, want single authorized key trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != testAuthorizedKey {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, testAuthorizedKey)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "single token public key", publicKey: "ssh-ed25519"},
		{name: "plain text public key", publicKey: "not-an-authorized-key"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("insert key with %s error = %v, want authorized key field trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("update key to %s error = %v, want authorized key field trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != testAuthorizedKey {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, testAuthorizedKey)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyType(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "plain text first field", publicKey: "not-an-authorized-key AAAA"},
		{name: "unknown key type", publicKey: "rsa AAAA"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a supported authorized key type") {
				t.Fatalf("insert key with %s error = %v, want authorized key type trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a supported authorized key type") {
				t.Fatalf("update key to %s error = %v, want authorized key type trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != testAuthorizedKey {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, testAuthorizedKey)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyBlob(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "short key blob", publicKey: "ssh-ed25519 AAAA"},
		{name: "non-base64 key blob", publicKey: "ssh-ed25519 !!!!"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("insert key with %s error = %v, want authorized key blob trigger error", tt.name, err)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", strings.Replace(testOtherAuthorizedKey, " ", "\t", 1), testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid tab-separated public key: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("update key to %s error = %v, want authorized key blob trigger error", tt.name, err)
			}
		})
	}

	wantPublicKey := strings.Replace(testOtherAuthorizedKey, " ", "\t", 1)
	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != wantPublicKey {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, wantPublicKey)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyBlobCharacters(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "invalid character after valid prefix", publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5!!!!"},
		{name: "invalid character before comment", publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5!!!! comment"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("insert key with %s error = %v, want authorized key blob character trigger error", tt.name, err)
			}
		})
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	validWithTabComment := strings.Replace(testOtherAuthorizedKey, " ", "\t", 1) + "\trotated@example"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", validWithTabComment, testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid tab-separated public key with comment: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("update key to %s error = %v, want authorized key blob character trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithTabComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithTabComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyBlobBase64Shape(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "blob length not multiple of four", publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5A"},
		{name: "padding before blob end", publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5=AAA"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("insert key with %s error = %v, want authorized key blob base64-shape trigger error", tt.name, err)
			}
		})
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	validWithTabComment := strings.Replace(testOtherAuthorizedKey, " ", "\t", 1) + "\trotated@example"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", validWithTabComment, testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid tab-separated public key with comment: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key must be a valid authorized key") {
				t.Fatalf("update key to %s error = %v, want authorized key blob base64-shape trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithTabComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithTabComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyAuthorizedKeyBlobType(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	ed25519Blob := strings.Fields(testAuthorizedKey)[1]
	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "blob with no encoded key type", publicKey: "ssh-ed25519 AAAAAAAAAAAAAAAA"},
		{name: "blob type does not match first field", publicKey: "ssh-rsa " + ed25519Blob},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key blob must match authorized key type") {
				t.Fatalf("insert key with %s error = %v, want authorized key blob type trigger error", tt.name, err)
			}
		})
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	validWithTabComment := strings.Replace(testOtherAuthorizedKey, " ", "\t", 1) + "\trotated@example"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", validWithTabComment, testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid tab-separated public key with comment: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "public key blob must match authorized key type") {
				t.Fatalf("update key to %s error = %v, want authorized key blob type trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithTabComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithTabComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyEd25519BlobLength(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	tests := []struct {
		name      string
		publicKey string
	}{
		{name: "truncated ed25519 blob", publicKey: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA"},
		{name: "extended ed25519 blob", publicKey: "ssh-ed25519 " + strings.Fields(testAuthorizedKey)[1] + "AAAA"},
	}
	for i, tt := range tests {
		t.Run("insert "+tt.name, func(t *testing.T) {
			fingerprint := testSHA256Fingerprint(byte('A' + i))
			_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, fingerprint, "user-1", tt.publicKey, now(), now())
			if err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "ssh-ed25519 public key blob must be complete") {
				t.Fatalf("insert key with %s error = %v, want ed25519 blob length trigger error", tt.name, err)
			}
		})
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	validWithTabComment := strings.Replace(testOtherAuthorizedKey, " ", "\t", 1) + "\trotated@example"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", validWithTabComment, testKeyFingerprint); err != nil {
		t.Fatalf("update key to valid tab-separated public key with comment: %v", err)
	}
	for _, tt := range tests {
		t.Run("update "+tt.name, func(t *testing.T) {
			_, err := st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", tt.publicKey, testKeyFingerprint)
			if err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
			if !strings.Contains(err.Error(), "ssh-ed25519 public key blob must be complete") {
				t.Fatalf("update key to %s error = %v, want ed25519 blob length trigger error", tt.name, err)
			}
		})
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithTabComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithTabComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	incompleteRSAKey := "ssh-rsa AAAAB3NzaC1yc2E="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", incompleteRSAKey, now(), now())
	if err == nil {
		t.Fatalf("inserted key with incomplete RSA blob, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include exponent and modulus") {
		t.Fatalf("insert key with incomplete RSA blob error = %v, want RSA blob fields trigger error", err)
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", incompleteRSAKey, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to incomplete RSA blob, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include exponent and modulus") {
		t.Fatalf("update key to incomplete RSA blob error = %v, want RSA blob fields trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithoutModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQAB"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithoutModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with RSA blob missing modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include modulus") {
		t.Fatalf("insert key with RSA blob missing modulus error = %v, want RSA blob modulus trigger error", err)
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithoutModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to RSA blob missing modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include modulus") {
		t.Fatalf("update key to RSA blob missing modulus error = %v, want RSA blob modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobModulusLength(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithPartialModulusLength := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAQID"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithPartialModulusLength, now(), now())
	if err == nil {
		t.Fatalf("inserted key with partial RSA modulus length, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include complete modulus length") {
		t.Fatalf("insert key with partial RSA modulus length error = %v, want RSA blob modulus length trigger error", err)
	}

	validWithComment := testAuthorizedKey + " alice@example"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", validWithComment, now(), now()); err != nil {
		t.Fatalf("insert valid key with comment: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithPartialModulusLength, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to partial RSA modulus length, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include complete modulus length") {
		t.Fatalf("update key to partial RSA modulus length error = %v, want RSA blob modulus length trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != validWithComment {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, validWithComment)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobModulusBytes(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithoutModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithoutModulusBytes, now(), now())
	if err == nil {
		t.Fatalf("inserted key with RSA blob missing modulus bytes, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include complete modulus") {
		t.Fatalf("insert key with RSA blob missing modulus bytes error = %v, want RSA blob modulus bytes trigger error", err)
	}

	rsaKeyWithOneModulusByte := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAQE="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithOneModulusByte, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with one modulus byte: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithoutModulusBytes, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to RSA blob missing modulus bytes, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include complete modulus") {
		t.Fatalf("update key to RSA blob missing modulus bytes error = %v, want RSA blob modulus bytes trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithOneModulusByte {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithOneModulusByte)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredModulusBytes(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAgE="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared modulus bytes") {
		t.Fatalf("insert key with short declared RSA modulus error = %v, want RSA declared modulus bytes trigger error", err)
	}

	rsaKeyWithTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAgEC"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared modulus bytes") {
		t.Fatalf("update key to short declared RSA modulus error = %v, want RSA declared modulus bytes trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAwEC"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared three-byte modulus") {
		t.Fatalf("insert key with short declared three-byte RSA modulus error = %v, want RSA declared three-byte modulus trigger error", err)
	}

	rsaKeyWithThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAAwECAw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared three-byte modulus") {
		t.Fatalf("update key to short declared three-byte RSA modulus error = %v, want RSA declared three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABAECAw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared four-byte modulus") {
		t.Fatalf("insert key with short declared four-byte RSA modulus error = %v, want RSA declared four-byte modulus trigger error", err)
	}

	rsaKeyWithFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABAECAwQ="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared four-byte modulus") {
		t.Fatalf("update key to short declared four-byte RSA modulus error = %v, want RSA declared four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABQECAwQ="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared five-byte modulus") {
		t.Fatalf("insert key with short declared five-byte RSA modulus error = %v, want RSA declared five-byte modulus trigger error", err)
	}

	rsaKeyWithFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABQECAwQF"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared five-byte modulus") {
		t.Fatalf("update key to short declared five-byte RSA modulus error = %v, want RSA declared five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABgECAwQF"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared six-byte modulus") {
		t.Fatalf("insert key with short declared six-byte RSA modulus error = %v, want RSA declared six-byte modulus trigger error", err)
	}

	rsaKeyWithSixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABgECAwQFBg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared six-byte modulus") {
		t.Fatalf("update key to short declared six-byte RSA modulus error = %v, want RSA declared six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABwECAwQFBg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seven-byte modulus") {
		t.Fatalf("insert key with short declared seven-byte RSA modulus error = %v, want RSA declared seven-byte modulus trigger error", err)
	}

	rsaKeyWithSevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAABwECAwQFBgc="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seven-byte modulus") {
		t.Fatalf("update key to short declared seven-byte RSA modulus error = %v, want RSA declared seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACAECAwQFBgc="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eight-byte modulus") {
		t.Fatalf("insert key with short declared eight-byte RSA modulus error = %v, want RSA declared eight-byte modulus trigger error", err)
	}

	rsaKeyWithEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACAECAwQFBgcI"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eight-byte modulus") {
		t.Fatalf("update key to short declared eight-byte RSA modulus error = %v, want RSA declared eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACQECAwQFBgcI"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared nine-byte modulus") {
		t.Fatalf("insert key with short declared nine-byte RSA modulus error = %v, want RSA declared nine-byte modulus trigger error", err)
	}

	rsaKeyWithNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACQECAwQFBgcICQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared nine-byte modulus") {
		t.Fatalf("update key to short declared nine-byte RSA modulus error = %v, want RSA declared nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACgECAwQFBgcICQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared ten-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared ten-byte modulus") {
		t.Fatalf("insert key with short declared ten-byte RSA modulus error = %v, want RSA declared ten-byte modulus trigger error", err)
	}

	rsaKeyWithTenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACgECAwQFBgcICQo="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with ten declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared ten-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared ten-byte modulus") {
		t.Fatalf("update key to short declared ten-byte RSA modulus error = %v, want RSA declared ten-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredElevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredElevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACwECAwQFBgcICQo="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredElevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eleven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eleven-byte modulus") {
		t.Fatalf("insert key with short declared eleven-byte RSA modulus error = %v, want RSA declared eleven-byte modulus trigger error", err)
	}

	rsaKeyWithElevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAACwECAwQFBgcICQoL"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithElevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eleven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredElevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eleven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eleven-byte modulus") {
		t.Fatalf("update key to short declared eleven-byte RSA modulus error = %v, want RSA declared eleven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithElevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithElevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwelveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwelveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADAECAwQFBgcICQoL"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwelveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twelve-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twelve-byte modulus") {
		t.Fatalf("insert key with short declared twelve-byte RSA modulus error = %v, want RSA declared twelve-byte modulus trigger error", err)
	}

	rsaKeyWithTwelveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADAECAwQFBgcICQoLDA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwelveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twelve declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwelveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twelve-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twelve-byte modulus") {
		t.Fatalf("update key to short declared twelve-byte RSA modulus error = %v, want RSA declared twelve-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwelveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwelveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADQECAwQFBgcICQoLDA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirteen-byte modulus") {
		t.Fatalf("insert key with short declared thirteen-byte RSA modulus error = %v, want RSA declared thirteen-byte modulus trigger error", err)
	}

	rsaKeyWithThirteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADQECAwQFBgcICQoLDA0="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirteen-byte modulus") {
		t.Fatalf("update key to short declared thirteen-byte RSA modulus error = %v, want RSA declared thirteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFourteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFourteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADgECAwQFBgcICQoLDA0="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFourteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fourteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fourteen-byte modulus") {
		t.Fatalf("insert key with short declared fourteen-byte RSA modulus error = %v, want RSA declared fourteen-byte modulus trigger error", err)
	}

	rsaKeyWithFourteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADgECAwQFBgcICQoLDA0O"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFourteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fourteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFourteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fourteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fourteen-byte modulus") {
		t.Fatalf("update key to short declared fourteen-byte RSA modulus error = %v, want RSA declared fourteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFourteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFourteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFifteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFifteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADwECAwQFBgcICQoLDA0O"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFifteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifteen-byte modulus") {
		t.Fatalf("insert key with short declared fifteen-byte RSA modulus error = %v, want RSA declared fifteen-byte modulus trigger error", err)
	}

	rsaKeyWithFifteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAADwECAwQFBgcICQoLDA0ODw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFifteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFifteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifteen-byte modulus") {
		t.Fatalf("update key to short declared fifteen-byte RSA modulus error = %v, want RSA declared fifteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFifteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFifteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEAECAwQFBgcICQoLDA0ODw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixteen-byte modulus") {
		t.Fatalf("insert key with short declared sixteen-byte RSA modulus error = %v, want RSA declared sixteen-byte modulus trigger error", err)
	}

	rsaKeyWithSixteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEAECAwQFBgcICQoLDA0ODxA="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixteen-byte modulus") {
		t.Fatalf("update key to short declared sixteen-byte RSA modulus error = %v, want RSA declared sixteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventeenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventeenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEQECAwQFBgcICQoLDA0ODxA="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventeenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventeen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventeen-byte modulus") {
		t.Fatalf("insert key with short declared seventeen-byte RSA modulus error = %v, want RSA declared seventeen-byte modulus trigger error", err)
	}

	rsaKeyWithSeventeenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEQECAwQFBgcICQoLDA0ODxAR"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventeenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventeen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventeenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventeen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventeen-byte modulus") {
		t.Fatalf("update key to short declared seventeen-byte RSA modulus error = %v, want RSA declared seventeen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventeenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventeenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredEighteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredEighteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEgECAwQFBgcICQoLDA0ODxAR"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredEighteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eighteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighteen-byte modulus") {
		t.Fatalf("insert key with short declared eighteen-byte RSA modulus error = %v, want RSA declared eighteen-byte modulus trigger error", err)
	}

	rsaKeyWithEighteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEgECAwQFBgcICQoLDA0ODxAREg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithEighteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eighteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredEighteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eighteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighteen-byte modulus") {
		t.Fatalf("update key to short declared eighteen-byte RSA modulus error = %v, want RSA declared eighteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithEighteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithEighteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredNineteenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredNineteenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEwECAwQFBgcICQoLDA0ODxAREg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredNineteenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared nineteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared nineteen-byte modulus") {
		t.Fatalf("insert key with short declared nineteen-byte RSA modulus error = %v, want RSA declared nineteen-byte modulus trigger error", err)
	}

	rsaKeyWithNineteenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAEwECAwQFBgcICQoLDA0ODxAREhM="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithNineteenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with nineteen declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredNineteenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared nineteen-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared nineteen-byte modulus") {
		t.Fatalf("update key to short declared nineteen-byte RSA modulus error = %v, want RSA declared nineteen-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithNineteenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithNineteenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFAECAwQFBgcICQoLDA0ODxAREhM="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-byte modulus") {
		t.Fatalf("insert key with short declared twenty-byte RSA modulus error = %v, want RSA declared twenty-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFAECAwQFBgcICQoLDA0ODxAREhMU"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-byte modulus") {
		t.Fatalf("update key to short declared twenty-byte RSA modulus error = %v, want RSA declared twenty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFQECAwQFBgcICQoLDA0ODxAREhMU"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-one-byte modulus") {
		t.Fatalf("insert key with short declared twenty-one-byte RSA modulus error = %v, want RSA declared twenty-one-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFQECAwQFBgcICQoLDA0ODxAREhMUFQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-one-byte modulus") {
		t.Fatalf("update key to short declared twenty-one-byte RSA modulus error = %v, want RSA declared twenty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFgECAwQFBgcICQoLDA0ODxAREhMUFQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-two-byte modulus") {
		t.Fatalf("insert key with short declared twenty-two-byte RSA modulus error = %v, want RSA declared twenty-two-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFgECAwQFBgcICQoLDA0ODxAREhMUFRY="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-two-byte modulus") {
		t.Fatalf("update key to short declared twenty-two-byte RSA modulus error = %v, want RSA declared twenty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFwECAwQFBgcICQoLDA0ODxAREhMUFRY="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-three-byte modulus") {
		t.Fatalf("insert key with short declared twenty-three-byte RSA modulus error = %v, want RSA declared twenty-three-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAFwECAwQFBgcICQoLDA0ODxAREhMUFRYX"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-three-byte modulus") {
		t.Fatalf("update key to short declared twenty-three-byte RSA modulus error = %v, want RSA declared twenty-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGAECAwQFBgcICQoLDA0ODxAREhMUFRYX"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-four-byte modulus") {
		t.Fatalf("insert key with short declared twenty-four-byte RSA modulus error = %v, want RSA declared twenty-four-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGAECAwQFBgcICQoLDA0ODxAREhMUFRYXGA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-four-byte modulus") {
		t.Fatalf("update key to short declared twenty-four-byte RSA modulus error = %v, want RSA declared twenty-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGQECAwQFBgcICQoLDA0ODxAREhMUFRYXGA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-five-byte modulus") {
		t.Fatalf("insert key with short declared twenty-five-byte RSA modulus error = %v, want RSA declared twenty-five-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBk="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-five-byte modulus") {
		t.Fatalf("update key to short declared twenty-five-byte RSA modulus error = %v, want RSA declared twenty-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBk="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-six-byte modulus") {
		t.Fatalf("insert key with short declared twenty-six-byte RSA modulus error = %v, want RSA declared twenty-six-byte modulus trigger error", err)
	}

	rsaKeyWithTwentySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBka"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-six-byte modulus") {
		t.Fatalf("update key to short declared twenty-six-byte RSA modulus error = %v, want RSA declared twenty-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBka"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-seven-byte modulus") {
		t.Fatalf("insert key with short declared twenty-seven-byte RSA modulus error = %v, want RSA declared twenty-seven-byte modulus trigger error", err)
	}

	rsaKeyWithTwentySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAGwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-seven-byte modulus") {
		t.Fatalf("update key to short declared twenty-seven-byte RSA modulus error = %v, want RSA declared twenty-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-eight-byte modulus") {
		t.Fatalf("insert key with short declared twenty-eight-byte RSA modulus error = %v, want RSA declared twenty-eight-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxw="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-eight-byte modulus") {
		t.Fatalf("update key to short declared twenty-eight-byte RSA modulus error = %v, want RSA declared twenty-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredTwentyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredTwentyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxw="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredTwentyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared twenty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-nine-byte modulus") {
		t.Fatalf("insert key with short declared twenty-nine-byte RSA modulus error = %v, want RSA declared twenty-nine-byte modulus trigger error", err)
	}

	rsaKeyWithTwentyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwd"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithTwentyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with twenty-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredTwentyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared twenty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared twenty-nine-byte modulus") {
		t.Fatalf("update key to short declared twenty-nine-byte RSA modulus error = %v, want RSA declared twenty-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithTwentyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithTwentyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwd"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-byte modulus") {
		t.Fatalf("insert key with short declared thirty-byte RSA modulus error = %v, want RSA declared thirty-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-byte modulus") {
		t.Fatalf("update key to short declared thirty-byte RSA modulus error = %v, want RSA declared thirty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-one-byte modulus") {
		t.Fatalf("insert key with short declared thirty-one-byte RSA modulus error = %v, want RSA declared thirty-one-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAHwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-one-byte modulus") {
		t.Fatalf("update key to short declared thirty-one-byte RSA modulus error = %v, want RSA declared thirty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-two-byte modulus") {
		t.Fatalf("insert key with short declared thirty-two-byte RSA modulus error = %v, want RSA declared thirty-two-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8g"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-two-byte modulus") {
		t.Fatalf("update key to short declared thirty-two-byte RSA modulus error = %v, want RSA declared thirty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8g"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-three-byte modulus") {
		t.Fatalf("insert key with short declared thirty-three-byte RSA modulus error = %v, want RSA declared thirty-three-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gIQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-three-byte modulus") {
		t.Fatalf("update key to short declared thirty-three-byte RSA modulus error = %v, want RSA declared thirty-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gIQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-four-byte modulus") {
		t.Fatalf("insert key with short declared thirty-four-byte RSA modulus error = %v, want RSA declared thirty-four-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISI="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-four-byte modulus") {
		t.Fatalf("update key to short declared thirty-four-byte RSA modulus error = %v, want RSA declared thirty-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISI="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-five-byte modulus") {
		t.Fatalf("insert key with short declared thirty-five-byte RSA modulus error = %v, want RSA declared thirty-five-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAIwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIj"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-five-byte modulus") {
		t.Fatalf("update key to short declared thirty-five-byte RSA modulus error = %v, want RSA declared thirty-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIj"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-six-byte modulus") {
		t.Fatalf("insert key with short declared thirty-six-byte RSA modulus error = %v, want RSA declared thirty-six-byte modulus trigger error", err)
	}

	rsaKeyWithThirtySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-six-byte modulus") {
		t.Fatalf("update key to short declared thirty-six-byte RSA modulus error = %v, want RSA declared thirty-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-seven-byte modulus") {
		t.Fatalf("insert key with short declared thirty-seven-byte RSA modulus error = %v, want RSA declared thirty-seven-byte modulus trigger error", err)
	}

	rsaKeyWithThirtySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCU="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-seven-byte modulus") {
		t.Fatalf("update key to short declared thirty-seven-byte RSA modulus error = %v, want RSA declared thirty-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCU="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-eight-byte modulus") {
		t.Fatalf("insert key with short declared thirty-eight-byte RSA modulus error = %v, want RSA declared thirty-eight-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUm"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-eight-byte modulus") {
		t.Fatalf("update key to short declared thirty-eight-byte RSA modulus error = %v, want RSA declared thirty-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredThirtyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredThirtyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUm"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredThirtyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared thirty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-nine-byte modulus") {
		t.Fatalf("insert key with short declared thirty-nine-byte RSA modulus error = %v, want RSA declared thirty-nine-byte modulus trigger error", err)
	}

	rsaKeyWithThirtyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAJwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithThirtyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with thirty-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredThirtyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared thirty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared thirty-nine-byte modulus") {
		t.Fatalf("update key to short declared thirty-nine-byte RSA modulus error = %v, want RSA declared thirty-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithThirtyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithThirtyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-byte modulus") {
		t.Fatalf("insert key with short declared forty-byte RSA modulus error = %v, want RSA declared forty-byte modulus trigger error", err)
	}

	rsaKeyWithFortyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJyg="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-byte modulus") {
		t.Fatalf("update key to short declared forty-byte RSA modulus error = %v, want RSA declared forty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJyg="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-one-byte modulus") {
		t.Fatalf("insert key with short declared forty-one-byte RSA modulus error = %v, want RSA declared forty-one-byte modulus trigger error", err)
	}

	rsaKeyWithFortyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygp"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-one-byte modulus") {
		t.Fatalf("update key to short declared forty-one-byte RSA modulus error = %v, want RSA declared forty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygp"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-two-byte modulus") {
		t.Fatalf("insert key with short declared forty-two-byte RSA modulus error = %v, want RSA declared forty-two-byte modulus trigger error", err)
	}

	rsaKeyWithFortyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-two-byte modulus") {
		t.Fatalf("update key to short declared forty-two-byte RSA modulus error = %v, want RSA declared forty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-three-byte modulus") {
		t.Fatalf("insert key with short declared forty-three-byte RSA modulus error = %v, want RSA declared forty-three-byte modulus trigger error", err)
	}

	rsaKeyWithFortyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAKwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKis="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-three-byte modulus") {
		t.Fatalf("update key to short declared forty-three-byte RSA modulus error = %v, want RSA declared forty-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKis="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-four-byte modulus") {
		t.Fatalf("insert key with short declared forty-four-byte RSA modulus error = %v, want RSA declared forty-four-byte modulus trigger error", err)
	}

	rsaKeyWithFortyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKiss"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-four-byte modulus") {
		t.Fatalf("update key to short declared forty-four-byte RSA modulus error = %v, want RSA declared forty-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKiss"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-five-byte modulus") {
		t.Fatalf("insert key with short declared forty-five-byte RSA modulus error = %v, want RSA declared forty-five-byte modulus trigger error", err)
	}

	rsaKeyWithFortyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-five-byte modulus") {
		t.Fatalf("update key to short declared forty-five-byte RSA modulus error = %v, want RSA declared forty-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-six-byte modulus") {
		t.Fatalf("insert key with short declared forty-six-byte RSA modulus error = %v, want RSA declared forty-six-byte modulus trigger error", err)
	}

	rsaKeyWithFortySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-six-byte modulus") {
		t.Fatalf("update key to short declared forty-six-byte RSA modulus error = %v, want RSA declared forty-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-seven-byte modulus") {
		t.Fatalf("insert key with short declared forty-seven-byte RSA modulus error = %v, want RSA declared forty-seven-byte modulus trigger error", err)
	}

	rsaKeyWithFortySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAALwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4v"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-seven-byte modulus") {
		t.Fatalf("update key to short declared forty-seven-byte RSA modulus error = %v, want RSA declared forty-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4v"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-eight-byte modulus") {
		t.Fatalf("insert key with short declared forty-eight-byte RSA modulus error = %v, want RSA declared forty-eight-byte modulus trigger error", err)
	}

	rsaKeyWithFortyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-eight-byte modulus") {
		t.Fatalf("update key to short declared forty-eight-byte RSA modulus error = %v, want RSA declared forty-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFortyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFortyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFortyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared forty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-nine-byte modulus") {
		t.Fatalf("insert key with short declared forty-nine-byte RSA modulus error = %v, want RSA declared forty-nine-byte modulus trigger error", err)
	}

	rsaKeyWithFortyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDE="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFortyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with forty-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFortyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared forty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared forty-nine-byte modulus") {
		t.Fatalf("update key to short declared forty-nine-byte RSA modulus error = %v, want RSA declared forty-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFortyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFortyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDE="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-byte modulus") {
		t.Fatalf("insert key with short declared fifty-byte RSA modulus error = %v, want RSA declared fifty-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEy"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-byte modulus") {
		t.Fatalf("update key to short declared fifty-byte RSA modulus error = %v, want RSA declared fifty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEy"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-one-byte modulus") {
		t.Fatalf("insert key with short declared fifty-one-byte RSA modulus error = %v, want RSA declared fifty-one-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAMwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-one-byte modulus") {
		t.Fatalf("update key to short declared fifty-one-byte RSA modulus error = %v, want RSA declared fifty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-two-byte modulus") {
		t.Fatalf("insert key with short declared fifty-two-byte RSA modulus error = %v, want RSA declared fifty-two-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-two-byte modulus") {
		t.Fatalf("update key to short declared fifty-two-byte RSA modulus error = %v, want RSA declared fifty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-three-byte modulus") {
		t.Fatalf("insert key with short declared fifty-three-byte RSA modulus error = %v, want RSA declared fifty-three-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-three-byte modulus") {
		t.Fatalf("update key to short declared fifty-three-byte RSA modulus error = %v, want RSA declared fifty-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-four-byte modulus") {
		t.Fatalf("insert key with short declared fifty-four-byte RSA modulus error = %v, want RSA declared fifty-four-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Ng=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-four-byte modulus") {
		t.Fatalf("update key to short declared fifty-four-byte RSA modulus error = %v, want RSA declared fifty-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Ng=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-five-byte modulus") {
		t.Fatalf("insert key with short declared fifty-five-byte RSA modulus error = %v, want RSA declared fifty-five-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAANwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-five-byte modulus") {
		t.Fatalf("update key to short declared fifty-five-byte RSA modulus error = %v, want RSA declared fifty-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-six-byte modulus") {
		t.Fatalf("insert key with short declared fifty-six-byte RSA modulus error = %v, want RSA declared fifty-six-byte modulus trigger error", err)
	}

	rsaKeyWithFiftySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-six-byte modulus") {
		t.Fatalf("update key to short declared fifty-six-byte RSA modulus error = %v, want RSA declared fifty-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-seven-byte modulus") {
		t.Fatalf("insert key with short declared fifty-seven-byte RSA modulus error = %v, want RSA declared fifty-seven-byte modulus trigger error", err)
	}

	rsaKeyWithFiftySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-seven-byte modulus") {
		t.Fatalf("update key to short declared fifty-seven-byte RSA modulus error = %v, want RSA declared fifty-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-eight-byte modulus") {
		t.Fatalf("insert key with short declared fifty-eight-byte RSA modulus error = %v, want RSA declared fifty-eight-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-eight-byte modulus") {
		t.Fatalf("update key to short declared fifty-eight-byte RSA modulus error = %v, want RSA declared fifty-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredFiftyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredFiftyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredFiftyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared fifty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-nine-byte modulus") {
		t.Fatalf("insert key with short declared fifty-nine-byte RSA modulus error = %v, want RSA declared fifty-nine-byte modulus trigger error", err)
	}

	rsaKeyWithFiftyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAOwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithFiftyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with fifty-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredFiftyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared fifty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared fifty-nine-byte modulus") {
		t.Fatalf("update key to short declared fifty-nine-byte RSA modulus error = %v, want RSA declared fifty-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithFiftyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithFiftyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-byte modulus") {
		t.Fatalf("insert key with short declared sixty-byte RSA modulus error = %v, want RSA declared sixty-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-byte modulus") {
		t.Fatalf("update key to short declared sixty-byte RSA modulus error = %v, want RSA declared sixty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-one-byte modulus") {
		t.Fatalf("insert key with short declared sixty-one-byte RSA modulus error = %v, want RSA declared sixty-one-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-one-byte modulus") {
		t.Fatalf("update key to short declared sixty-one-byte RSA modulus error = %v, want RSA declared sixty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-two-byte modulus") {
		t.Fatalf("insert key with short declared sixty-two-byte RSA modulus error = %v, want RSA declared sixty-two-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-two-byte modulus") {
		t.Fatalf("update key to short declared sixty-two-byte RSA modulus error = %v, want RSA declared sixty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-three-byte modulus") {
		t.Fatalf("insert key with short declared sixty-three-byte RSA modulus error = %v, want RSA declared sixty-three-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAPwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+Pw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-three-byte modulus") {
		t.Fatalf("update key to short declared sixty-three-byte RSA modulus error = %v, want RSA declared sixty-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+Pw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-four-byte modulus") {
		t.Fatalf("insert key with short declared sixty-four-byte RSA modulus error = %v, want RSA declared sixty-four-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0A="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-four-byte modulus") {
		t.Fatalf("update key to short declared sixty-four-byte RSA modulus error = %v, want RSA declared sixty-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0A="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-five-byte modulus") {
		t.Fatalf("insert key with short declared sixty-five-byte RSA modulus error = %v, want RSA declared sixty-five-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BB"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-five-byte modulus") {
		t.Fatalf("update key to short declared sixty-five-byte RSA modulus error = %v, want RSA declared sixty-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BB"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-six-byte modulus") {
		t.Fatalf("insert key with short declared sixty-six-byte RSA modulus error = %v, want RSA declared sixty-six-byte modulus trigger error", err)
	}

	rsaKeyWithSixtySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-six-byte modulus") {
		t.Fatalf("update key to short declared sixty-six-byte RSA modulus error = %v, want RSA declared sixty-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-seven-byte modulus") {
		t.Fatalf("insert key with short declared sixty-seven-byte RSA modulus error = %v, want RSA declared sixty-seven-byte modulus trigger error", err)
	}

	rsaKeyWithSixtySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAQwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkM="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-seven-byte modulus") {
		t.Fatalf("update key to short declared sixty-seven-byte RSA modulus error = %v, want RSA declared sixty-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkM="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-eight-byte modulus") {
		t.Fatalf("insert key with short declared sixty-eight-byte RSA modulus error = %v, want RSA declared sixty-eight-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNE"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-eight-byte modulus") {
		t.Fatalf("update key to short declared sixty-eight-byte RSA modulus error = %v, want RSA declared sixty-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSixtyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSixtyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNE"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSixtyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared sixty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-nine-byte modulus") {
		t.Fatalf("insert key with short declared sixty-nine-byte RSA modulus error = %v, want RSA declared sixty-nine-byte modulus trigger error", err)
	}

	rsaKeyWithSixtyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSixtyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with sixty-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSixtyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared sixty-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared sixty-nine-byte modulus") {
		t.Fatalf("update key to short declared sixty-nine-byte RSA modulus error = %v, want RSA declared sixty-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSixtyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSixtyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-byte modulus") {
		t.Fatalf("insert key with short declared seventy-byte RSA modulus error = %v, want RSA declared seventy-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUY="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-byte modulus") {
		t.Fatalf("update key to short declared seventy-byte RSA modulus error = %v, want RSA declared seventy-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUY="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-one-byte modulus") {
		t.Fatalf("insert key with short declared seventy-one-byte RSA modulus error = %v, want RSA declared seventy-one-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAARwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZH"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-one-byte modulus") {
		t.Fatalf("update key to short declared seventy-one-byte RSA modulus error = %v, want RSA declared seventy-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZH"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-two-byte modulus") {
		t.Fatalf("insert key with short declared seventy-two-byte RSA modulus error = %v, want RSA declared seventy-two-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSA=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-two-byte modulus") {
		t.Fatalf("update key to short declared seventy-two-byte RSA modulus error = %v, want RSA declared seventy-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyThreeByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyThreeByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSA=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyThreeByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-three-byte modulus") {
		t.Fatalf("insert key with short declared seventy-three-byte RSA modulus error = %v, want RSA declared seventy-three-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyThreeModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSEk="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyThreeModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-three declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyThreeByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-three-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-three-byte modulus") {
		t.Fatalf("update key to short declared seventy-three-byte RSA modulus error = %v, want RSA declared seventy-three-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyThreeModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyThreeModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyFourByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyFourByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSEk="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyFourByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-four-byte modulus") {
		t.Fatalf("insert key with short declared seventy-four-byte RSA modulus error = %v, want RSA declared seventy-four-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyFourModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElK"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyFourModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-four declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyFourByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-four-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-four-byte modulus") {
		t.Fatalf("update key to short declared seventy-four-byte RSA modulus error = %v, want RSA declared seventy-four-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyFourModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyFourModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyFiveByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyFiveByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElK"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyFiveByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-five-byte modulus") {
		t.Fatalf("insert key with short declared seventy-five-byte RSA modulus error = %v, want RSA declared seventy-five-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyFiveModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAASwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKSw=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyFiveModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-five declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyFiveByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-five-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-five-byte modulus") {
		t.Fatalf("update key to short declared seventy-five-byte RSA modulus error = %v, want RSA declared seventy-five-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyFiveModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyFiveModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventySixByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventySixByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKSw=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventySixByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-six-byte modulus") {
		t.Fatalf("insert key with short declared seventy-six-byte RSA modulus error = %v, want RSA declared seventy-six-byte modulus trigger error", err)
	}

	rsaKeyWithSeventySixModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0w="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventySixModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-six declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventySixByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-six-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-six-byte modulus") {
		t.Fatalf("update key to short declared seventy-six-byte RSA modulus error = %v, want RSA declared seventy-six-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventySixModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventySixModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventySevenByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventySevenByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0w="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventySevenByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-seven-byte modulus") {
		t.Fatalf("insert key with short declared seventy-seven-byte RSA modulus error = %v, want RSA declared seventy-seven-byte modulus trigger error", err)
	}

	rsaKeyWithSeventySevenModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xN"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventySevenModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-seven declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventySevenByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-seven-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-seven-byte modulus") {
		t.Fatalf("update key to short declared seventy-seven-byte RSA modulus error = %v, want RSA declared seventy-seven-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventySevenModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventySevenModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyEightByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyEightByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xN"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyEightByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-eight-byte modulus") {
		t.Fatalf("insert key with short declared seventy-eight-byte RSA modulus error = %v, want RSA declared seventy-eight-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyEightModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTg=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyEightModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-eight declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyEightByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-eight-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-eight-byte modulus") {
		t.Fatalf("update key to short declared seventy-eight-byte RSA modulus error = %v, want RSA declared seventy-eight-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyEightModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyEightModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredSeventyNineByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredSeventyNineByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTg=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredSeventyNineByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared seventy-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-nine-byte modulus") {
		t.Fatalf("insert key with short declared seventy-nine-byte RSA modulus error = %v, want RSA declared seventy-nine-byte modulus trigger error", err)
	}

	rsaKeyWithSeventyNineModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAATwECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk8="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithSeventyNineModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with seventy-nine declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredSeventyNineByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared seventy-nine-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared seventy-nine-byte modulus") {
		t.Fatalf("update key to short declared seventy-nine-byte RSA modulus error = %v, want RSA declared seventy-nine-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithSeventyNineModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithSeventyNineModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredEightyByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredEightyByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk8="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredEightyByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eighty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-byte modulus") {
		t.Fatalf("insert key with short declared eighty-byte RSA modulus error = %v, want RSA declared eighty-byte modulus trigger error", err)
	}

	rsaKeyWithEightyModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9Q"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithEightyModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eighty declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredEightyByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eighty-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-byte modulus") {
		t.Fatalf("update key to short declared eighty-byte RSA modulus error = %v, want RSA declared eighty-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithEightyModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithEightyModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredEightyOneByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredEightyOneByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9Q"
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredEightyOneByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eighty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-one-byte modulus") {
		t.Fatalf("insert key with short declared eighty-one-byte RSA modulus error = %v, want RSA declared eighty-one-byte modulus trigger error", err)
	}

	rsaKeyWithEightyOneModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUQECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9QUQ=="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithEightyOneModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eighty-one declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredEightyOneByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eighty-one-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-one-byte modulus") {
		t.Fatalf("update key to short declared eighty-one-byte RSA modulus error = %v, want RSA declared eighty-one-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithEightyOneModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithEightyOneModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyPublicKeyRSABlobDeclaredEightyTwoByteModulus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	rsaKeyWithShortDeclaredEightyTwoByteModulus := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9QUQ=="
	_, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testSHA256Fingerprint('A'), "user-1", rsaKeyWithShortDeclaredEightyTwoByteModulus, now(), now())
	if err == nil {
		t.Fatalf("inserted key with short declared eighty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-two-byte modulus") {
		t.Fatalf("insert key with short declared eighty-two-byte RSA modulus error = %v, want RSA declared eighty-two-byte modulus trigger error", err)
	}

	rsaKeyWithEightyTwoModulusBytes := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAUgECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9QUVI="
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", rsaKeyWithEightyTwoModulusBytes, now(), now()); err != nil {
		t.Fatalf("insert valid RSA key with eighty-two declared modulus bytes: %v", err)
	}
	_, err = st.db.ExecContext(ctx, "UPDATE keys SET public_key = ? WHERE fingerprint = ?", rsaKeyWithShortDeclaredEightyTwoByteModulus, testKeyFingerprint)
	if err == nil {
		t.Fatalf("updated key to short declared eighty-two-byte RSA modulus, want trigger error")
	}
	if !strings.Contains(err.Error(), "ssh-rsa public key blob must include declared eighty-two-byte modulus") {
		t.Fatalf("update key to short declared eighty-two-byte RSA modulus error = %v, want RSA declared eighty-two-byte modulus trigger error", err)
	}

	var gotPublicKey string
	row := st.db.QueryRowContext(ctx, "SELECT public_key FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotPublicKey); err != nil {
		t.Fatalf("query key public_key: %v", err)
	}
	if gotPublicKey != rsaKeyWithEightyTwoModulusBytes {
		t.Fatalf("key public_key = %q, want %q", gotPublicKey, rsaKeyWithEightyTwoModulusBytes)
	}
}

func TestEnsureSchemaEnforcesKeyAddedAtValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
		addedAt     string
	}{
		{name: "blank addition time", fingerprint: testSHA256Fingerprint('A'), addedAt: " \t "},
		{name: "padded addition time", fingerprint: testSHA256Fingerprint('B'), addedAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, tt.addedAt, now()); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	addedAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, addedAt, now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range []struct {
		name    string
		addedAt string
	}{
		{name: "blank addition time", addedAt: "\n\t"},
		{name: "padded addition time", addedAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE keys SET added_at = ? WHERE fingerprint = ?", tt.addedAt, testKeyFingerprint); err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
		})
	}

	var gotAddedAt string
	row := st.db.QueryRowContext(ctx, "SELECT added_at FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotAddedAt); err != nil {
		t.Fatalf("query key added_at: %v", err)
	}
	if gotAddedAt != addedAt {
		t.Fatalf("key added_at = %q, want %q", gotAddedAt, addedAt)
	}
}

func TestEnsureSchemaEnforcesKeyAddedAtFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
		addedAt     string
	}{
		{name: "space separated addition time", fingerprint: testSHA256Fingerprint('C'), addedAt: "2026-05-03 14:03:32"},
		{name: "missing timezone addition time", fingerprint: testSHA256Fingerprint('D'), addedAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, tt.addedAt, now()); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	addedAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, addedAt, now()); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET added_at = ? WHERE fingerprint = ?", "2026-05-03 14:03:32", testKeyFingerprint); err == nil {
		t.Fatalf("updated key to space separated addition time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET added_at = ? WHERE fingerprint = ?", validOffsetTime, testKeyFingerprint); err != nil {
		t.Fatalf("updated key to valid offset addition time: %v", err)
	}

	var gotAddedAt string
	row := st.db.QueryRowContext(ctx, "SELECT added_at FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotAddedAt); err != nil {
		t.Fatalf("query key added_at: %v", err)
	}
	if gotAddedAt != validOffsetTime {
		t.Fatalf("key added_at = %q, want %q", gotAddedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesKeyLastSeenAtValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
		lastSeenAt  string
	}{
		{name: "blank last seen time", fingerprint: testSHA256Fingerprint('E'), lastSeenAt: " \t "},
		{name: "padded last seen time", fingerprint: testSHA256Fingerprint('F'), lastSeenAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, now(), tt.lastSeenAt); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	lastSeenAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), lastSeenAt); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	for _, tt := range []struct {
		name       string
		lastSeenAt string
	}{
		{name: "blank last seen time", lastSeenAt: "\n\t"},
		{name: "padded last seen time", lastSeenAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE keys SET last_seen_at = ? WHERE fingerprint = ?", tt.lastSeenAt, testKeyFingerprint); err == nil {
				t.Fatalf("updated key to %s, want trigger error", tt.name)
			}
		})
	}

	var gotLastSeenAt string
	row := st.db.QueryRowContext(ctx, "SELECT last_seen_at FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotLastSeenAt); err != nil {
		t.Fatalf("query key last_seen_at: %v", err)
	}
	if gotLastSeenAt != lastSeenAt {
		t.Fatalf("key last_seen_at = %q, want %q", gotLastSeenAt, lastSeenAt)
	}
}

func TestEnsureSchemaEnforcesKeyLastSeenAtFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert user fixture: %v", err)
	}

	for _, tt := range []struct {
		name        string
		fingerprint string
		lastSeenAt  string
	}{
		{name: "space separated last seen time", fingerprint: testSHA256Fingerprint('G'), lastSeenAt: "2026-05-03 14:03:32"},
		{name: "missing timezone last seen time", fingerprint: testSHA256Fingerprint('H'), lastSeenAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, tt.fingerprint, "user-1", testAuthorizedKey, now(), tt.lastSeenAt); err == nil {
				t.Fatalf("inserted key with %s, want trigger error", tt.name)
			}
		})
	}

	lastSeenAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)`, testKeyFingerprint, "user-1", testAuthorizedKey, now(), lastSeenAt); err != nil {
		t.Fatalf("insert valid key: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET last_seen_at = ? WHERE fingerprint = ?", "2026-05-03 14:03:32", testKeyFingerprint); err == nil {
		t.Fatalf("updated key to space separated last seen time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET last_seen_at = ? WHERE fingerprint = ?", validOffsetTime, testKeyFingerprint); err != nil {
		t.Fatalf("updated key to valid offset last seen time: %v", err)
	}

	var gotLastSeenAt string
	row := st.db.QueryRowContext(ctx, "SELECT last_seen_at FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&gotLastSeenAt); err != nil {
		t.Fatalf("query key last_seen_at: %v", err)
	}
	if gotLastSeenAt != validOffsetTime {
		t.Fatalf("key last_seen_at = %q, want %q", gotLastSeenAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesAuditEventTypeValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		eventID   string
		eventType string
	}{
		{name: "blank event type", eventID: "bad-blank-event-type", eventType: " \t "},
		{name: "padded event type", eventID: "bad-padded-event-type", eventType: " test.audit "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, tt.eventType, `{"ok":true}`, now()); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	eventType := "test.audit"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, "audit-1", eventType, `{"ok":true}`, now()); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	for _, tt := range []struct {
		name      string
		eventType string
	}{
		{name: "blank event type", eventType: "\n\t"},
		{name: "padded event type", eventType: "\ttest.audit\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET event_type = ? WHERE id = ?", tt.eventType, "audit-1"); err == nil {
				t.Fatalf("updated audit event to %s, want trigger error", tt.name)
			}
		})
	}

	var gotEventType string
	row := st.db.QueryRowContext(ctx, "SELECT event_type FROM audit_events WHERE id = ?", "audit-1")
	if err := row.Scan(&gotEventType); err != nil {
		t.Fatalf("query audit event type: %v", err)
	}
	if gotEventType != eventType {
		t.Fatalf("audit event_type = %q, want %q", gotEventType, eventType)
	}
}

func TestEnsureSchemaEnforcesAuditDataValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name    string
		eventID string
		data    string
	}{
		{name: "blank data", eventID: "bad-blank-data", data: " \t "},
		{name: "padded data", eventID: "bad-padded-data", data: ` {"ok":true} `},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, "test.audit", tt.data, now()); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	data := `{"ok":true}`
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, "audit-1", "test.audit", data, now()); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	for _, tt := range []struct {
		name string
		data string
	}{
		{name: "blank data", data: "\n\t"},
		{name: "padded data", data: "\t{\"ok\":true}\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET data_json = ? WHERE id = ?", tt.data, "audit-1"); err == nil {
				t.Fatalf("updated audit event to %s, want trigger error", tt.name)
			}
		})
	}

	var gotData string
	row := st.db.QueryRowContext(ctx, "SELECT data_json FROM audit_events WHERE id = ?", "audit-1")
	if err := row.Scan(&gotData); err != nil {
		t.Fatalf("query audit event data: %v", err)
	}
	if gotData != data {
		t.Fatalf("audit data_json = %q, want %q", gotData, data)
	}
}

func TestEnsureSchemaEnforcesAuditDataJSONFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name    string
		eventID string
		data    string
	}{
		{name: "missing close", eventID: "bad-missing-close-data", data: `{"missing-close":`},
		{name: "plain text", eventID: "bad-plain-text-data", data: `not-json`},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, "test.audit", tt.data, now()); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	data := `{"ok":true}`
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, "audit-1", "test.audit", data, now()); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	for _, tt := range []struct {
		name string
		data string
	}{
		{name: "missing close", data: `{"missing-close":`},
		{name: "plain text", data: `not-json`},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET data_json = ? WHERE id = ?", tt.data, "audit-1"); err == nil {
				t.Fatalf("updated audit event to %s, want trigger error", tt.name)
			}
		})
	}

	validArray := `[{"ok":true}]`
	if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET data_json = ? WHERE id = ?", validArray, "audit-1"); err != nil {
		t.Fatalf("updated audit event to valid JSON array: %v", err)
	}

	var gotData string
	row := st.db.QueryRowContext(ctx, "SELECT data_json FROM audit_events WHERE id = ?", "audit-1")
	if err := row.Scan(&gotData); err != nil {
		t.Fatalf("query audit event data: %v", err)
	}
	if gotData != validArray {
		t.Fatalf("audit data_json = %q, want %q", gotData, validArray)
	}
}

func TestEnsureSchemaEnforcesAuditCreatedAtValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		eventID   string
		createdAt string
	}{
		{name: "blank creation time", eventID: "bad-blank-created-at", createdAt: " \t "},
		{name: "padded creation time", eventID: "bad-padded-created-at", createdAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, "test.audit", `{"ok":true}`, tt.createdAt); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	createdAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, "audit-1", "test.audit", `{"ok":true}`, createdAt); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	for _, tt := range []struct {
		name      string
		createdAt string
	}{
		{name: "blank creation time", createdAt: "\n\t"},
		{name: "padded creation time", createdAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET created_at = ? WHERE id = ?", tt.createdAt, "audit-1"); err == nil {
				t.Fatalf("updated audit event to %s, want trigger error", tt.name)
			}
		})
	}

	var gotCreatedAt string
	row := st.db.QueryRowContext(ctx, "SELECT created_at FROM audit_events WHERE id = ?", "audit-1")
	if err := row.Scan(&gotCreatedAt); err != nil {
		t.Fatalf("query audit event created_at: %v", err)
	}
	if gotCreatedAt != createdAt {
		t.Fatalf("audit created_at = %q, want %q", gotCreatedAt, createdAt)
	}
}

func TestEnsureSchemaEnforcesAuditCreatedAtFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		eventID   string
		createdAt string
	}{
		{name: "space separated creation time", eventID: "bad-space-created-at", createdAt: "2026-05-03 14:03:32"},
		{name: "missing timezone creation time", eventID: "bad-missing-timezone-created-at", createdAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, tt.eventID, "test.audit", `{"ok":true}`, tt.createdAt); err == nil {
				t.Fatalf("inserted audit event with %s, want trigger error", tt.name)
			}
		})
	}

	createdAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, "audit-1", "test.audit", `{"ok":true}`, createdAt); err != nil {
		t.Fatalf("insert valid audit event: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET created_at = ? WHERE id = ?", "2026-05-03 14:03:32", "audit-1"); err == nil {
		t.Fatalf("updated audit event to space separated creation time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE audit_events SET created_at = ? WHERE id = ?", validOffsetTime, "audit-1"); err != nil {
		t.Fatalf("updated audit event to valid offset creation time: %v", err)
	}

	var gotCreatedAt string
	row := st.db.QueryRowContext(ctx, "SELECT created_at FROM audit_events WHERE id = ?", "audit-1")
	if err := row.Scan(&gotCreatedAt); err != nil {
		t.Fatalf("query audit event created_at: %v", err)
	}
	if gotCreatedAt != validOffsetTime {
		t.Fatalf("audit created_at = %q, want %q", gotCreatedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesVMStartTimeValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name      string
		vmID      string
		startedAt string
	}{
		{name: "blank start time", vmID: "bad-blank-start-time", startedAt: " \t "},
		{name: "padded start time", vmID: "bad-padded-start-time", startedAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, tt.vmID, session.ID, filepath.Join(t.TempDir(), tt.vmID), 1234, tt.startedAt); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	startedAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, startedAt); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	for _, tt := range []struct {
		name      string
		startedAt string
	}{
		{name: "blank start time", startedAt: "\n\t"},
		{name: "padded start time", startedAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE vms SET started_at = ? WHERE id = ?", tt.startedAt, "vm-1"); err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
		})
	}

	var gotStartedAt string
	row := st.db.QueryRowContext(ctx, "SELECT started_at FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotStartedAt); err != nil {
		t.Fatalf("query VM started_at: %v", err)
	}
	if gotStartedAt != startedAt {
		t.Fatalf("VM started_at = %q, want %q", gotStartedAt, startedAt)
	}
}

func TestEnsureSchemaEnforcesVMStartTimeFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name      string
		vmID      string
		startedAt string
	}{
		{name: "space separated start time", vmID: "bad-space-start-time", startedAt: "2026-05-03 14:03:32"},
		{name: "missing timezone start time", vmID: "bad-missing-timezone-start-time", startedAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, tt.vmID, session.ID, filepath.Join(t.TempDir(), tt.vmID), 1234, tt.startedAt); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	startedAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, startedAt); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET started_at = ? WHERE id = ?", "2026-05-03 14:03:32", "vm-1"); err == nil {
		t.Fatalf("updated VM to space separated start time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET started_at = ? WHERE id = ?", validOffsetTime, "vm-1"); err != nil {
		t.Fatalf("updated VM to valid offset start time: %v", err)
	}

	var gotStartedAt string
	row := st.db.QueryRowContext(ctx, "SELECT started_at FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotStartedAt); err != nil {
		t.Fatalf("query VM started_at: %v", err)
	}
	if gotStartedAt != validOffsetTime {
		t.Fatalf("VM started_at = %q, want %q", gotStartedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesVMEndTimeValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name    string
		vmID    string
		endedAt string
	}{
		{name: "blank end time", vmID: "bad-blank-end-time", endedAt: " \t "},
		{name: "padded end time", vmID: "bad-padded-end-time", endedAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at, ended_at, exit_status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, tt.vmID, session.ID, filepath.Join(t.TempDir(), tt.vmID), 1234, now(), tt.endedAt, 0); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid active VM: %v", err)
	}
	for _, tt := range []struct {
		name    string
		endedAt string
	}{
		{name: "blank end time", endedAt: "\n\t"},
		{name: "padded end time", endedAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", tt.endedAt, 0, "vm-1"); err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
		})
	}

	endedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", endedAt, 0, "vm-1"); err != nil {
		t.Fatalf("complete VM with valid end time: %v", err)
	}
	var gotEndedAt sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT ended_at FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotEndedAt); err != nil {
		t.Fatalf("query VM ended_at: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != endedAt {
		t.Fatalf("VM ended_at = %v, want %q", gotEndedAt, endedAt)
	}
}

func TestEnsureSchemaEnforcesVMEndTimeFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name    string
		vmID    string
		endedAt string
	}{
		{name: "space separated end time", vmID: "bad-space-end-time", endedAt: "2026-05-03 14:03:32"},
		{name: "missing timezone end time", vmID: "bad-missing-timezone-end-time", endedAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at, ended_at, exit_status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, tt.vmID, session.ID, filepath.Join(t.TempDir(), tt.vmID), 1234, now(), tt.endedAt, 0); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid active VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", "2026-05-03 14:03:32", 0, "vm-1"); err == nil {
		t.Fatalf("updated VM to space separated end time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", validOffsetTime, 0, "vm-1"); err != nil {
		t.Fatalf("updated VM to valid offset end time: %v", err)
	}

	var gotEndedAt sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT ended_at FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotEndedAt); err != nil {
		t.Fatalf("query VM ended_at: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != validOffsetTime {
		t.Fatalf("VM ended_at = %v, want %q", gotEndedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesSessionEndTimeValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
		endedAt   string
	}{
		{name: "blank end time", sessionID: "bad-blank-end-time", endedAt: " \t "},
		{name: "padded end time", sessionID: "bad-padded-end-time", endedAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, ended_at, status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), tt.endedAt, "closed"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid active session: %v", err)
	}
	for _, tt := range []struct {
		name    string
		endedAt string
	}{
		{name: "blank end time", endedAt: "\n\t"},
		{name: "padded end time", endedAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", tt.endedAt, "closed", "session-1"); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	endedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", endedAt, "closed", "session-1"); err != nil {
		t.Fatalf("complete session with valid end time: %v", err)
	}

	var gotEndedAt sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT ended_at FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotEndedAt); err != nil {
		t.Fatalf("query session ended_at: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != endedAt {
		t.Fatalf("session ended_at = %v, want %q", gotEndedAt, endedAt)
	}
}

func TestEnsureSchemaEnforcesSessionEndTimeFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
		endedAt   string
	}{
		{name: "space separated end time", sessionID: "bad-space-end-time", endedAt: "2026-05-03 14:03:32"},
		{name: "missing timezone end time", sessionID: "bad-missing-timezone-end-time", endedAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, ended_at, status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", now(), tt.endedAt, "closed"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	startedAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", startedAt, "active"); err != nil {
		t.Fatalf("insert valid active session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", "2026-05-03 14:03:32", "closed", "session-1"); err == nil {
		t.Fatalf("updated session to space separated end time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", validOffsetTime, "closed", "session-1"); err != nil {
		t.Fatalf("updated session to valid offset end time: %v", err)
	}

	var gotEndedAt sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT ended_at FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotEndedAt); err != nil {
		t.Fatalf("query session ended_at: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != validOffsetTime {
		t.Fatalf("session ended_at = %v, want %q", gotEndedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesSessionStartTimeValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
		startedAt string
	}{
		{name: "blank start time", sessionID: "bad-blank-start-time", startedAt: " \t "},
		{name: "padded start time", sessionID: "bad-padded-start-time", startedAt: " " + now() + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", tt.startedAt, "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	startedAt := now()
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", startedAt, "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name      string
		startedAt string
	}{
		{name: "blank start time", startedAt: "\n\t"},
		{name: "padded start time", startedAt: "\t" + now() + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET started_at = ? WHERE id = ?", tt.startedAt, "session-1"); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	var gotStartedAt string
	row := st.db.QueryRowContext(ctx, "SELECT started_at FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotStartedAt); err != nil {
		t.Fatalf("query session started_at: %v", err)
	}
	if gotStartedAt != startedAt {
		t.Fatalf("session started_at = %q, want %q", gotStartedAt, startedAt)
	}
}

func TestEnsureSchemaEnforcesSessionStartTimeFormat(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name      string
		sessionID string
		startedAt string
	}{
		{name: "space separated start time", sessionID: "bad-space-start-time", startedAt: "2026-05-03 14:03:32"},
		{name: "missing timezone start time", sessionID: "bad-missing-timezone-start-time", startedAt: "2026-05-03T14:03:32"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, "127.0.0.1:2222", tt.startedAt, "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	startedAt := "2026-05-03T14:03:32Z"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", startedAt, "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET started_at = ? WHERE id = ?", "2026-05-03 14:03:32", "session-1"); err == nil {
		t.Fatalf("updated session to space separated start time, want trigger error")
	}

	validOffsetTime := "2026-05-03T14:03:32-04:00"
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET started_at = ? WHERE id = ?", validOffsetTime, "session-1"); err != nil {
		t.Fatalf("updated session to valid offset start time: %v", err)
	}

	var gotStartedAt string
	row := st.db.QueryRowContext(ctx, "SELECT started_at FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotStartedAt); err != nil {
		t.Fatalf("query session started_at: %v", err)
	}
	if gotStartedAt != validOffsetTime {
		t.Fatalf("session started_at = %q, want %q", gotStartedAt, validOffsetTime)
	}
}

func TestEnsureSchemaEnforcesSessionRemoteAddressValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name       string
		sessionID  string
		remoteAddr string
	}{
		{name: "blank remote address", sessionID: "bad-blank-remote-address", remoteAddr: " \t "},
		{name: "padded remote address", sessionID: "bad-padded-remote-address", remoteAddr: " 127.0.0.1:2222 "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, tt.remoteAddr, now(), "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	remoteAddr := "127.0.0.1:2222"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, remoteAddr, now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name       string
		remoteAddr string
	}{
		{name: "blank remote address", remoteAddr: "\n\t"},
		{name: "padded remote address", remoteAddr: "\t127.0.0.1:2223\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET remote_addr = ? WHERE id = ?", tt.remoteAddr, "session-1"); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	var gotRemoteAddr string
	row := st.db.QueryRowContext(ctx, "SELECT remote_addr FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotRemoteAddr); err != nil {
		t.Fatalf("query session remote_addr: %v", err)
	}
	if gotRemoteAddr != remoteAddr {
		t.Fatalf("session remote_addr = %q, want %q", gotRemoteAddr, remoteAddr)
	}
}

func TestEnsureSchemaEnforcesSessionRemoteAddressTCPValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	for _, tt := range []struct {
		name       string
		sessionID  string
		remoteAddr string
	}{
		{name: "missing port", sessionID: "bad-missing-port", remoteAddr: "127.0.0.1"},
		{name: "missing host", sessionID: "bad-missing-host", remoteAddr: ":2222"},
		{name: "empty port", sessionID: "bad-empty-port", remoteAddr: "127.0.0.1:"},
		{name: "non-numeric port", sessionID: "bad-nonnumeric-port", remoteAddr: "127.0.0.1:ssh"},
		{name: "zero port", sessionID: "bad-zero-port", remoteAddr: "127.0.0.1:0"},
		{name: "out of range port", sessionID: "bad-out-of-range-port", remoteAddr: "127.0.0.1:65536"},
		{name: "unbracketed IPv6 address", sessionID: "bad-unbracketed-ipv6", remoteAddr: "::1:2222"},
		{name: "bracketed IPv6 missing port", sessionID: "bad-bracketed-ipv6-missing-port", remoteAddr: "[::1]:"},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, tt.sessionID, userID, testKeyFingerprint, tt.remoteAddr, now(), "active"); err == nil {
				t.Fatalf("inserted session with %s, want trigger error", tt.name)
			}
		})
	}

	remoteAddr := "127.0.0.1:2222"
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, remoteAddr, now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	for _, tt := range []struct {
		name       string
		remoteAddr string
	}{
		{name: "missing port", remoteAddr: "127.0.0.1"},
		{name: "missing host", remoteAddr: ":2222"},
		{name: "non-numeric port", remoteAddr: "127.0.0.1:ssh"},
		{name: "zero port", remoteAddr: "127.0.0.1:0"},
		{name: "out of range port", remoteAddr: "127.0.0.1:65536"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET remote_addr = ? WHERE id = ?", tt.remoteAddr, "session-1"); err == nil {
				t.Fatalf("updated session to %s, want trigger error", tt.name)
			}
		})
	}

	validIPv6Addr := "[::1]:2223"
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET remote_addr = ? WHERE id = ?", validIPv6Addr, "session-1"); err != nil {
		t.Fatalf("updated session to valid bracketed IPv6 remote address: %v", err)
	}

	var gotRemoteAddr string
	row := st.db.QueryRowContext(ctx, "SELECT remote_addr FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotRemoteAddr); err != nil {
		t.Fatalf("query session remote_addr: %v", err)
	}
	if gotRemoteAddr != validIPv6Addr {
		t.Fatalf("session remote_addr = %q, want %q", gotRemoteAddr, validIPv6Addr)
	}
}

func TestEnsureSchemaEnforcesVMStateDirectoryValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name     string
		stateDir string
	}{
		{name: "blank state directory", stateDir: " \t "},
		{name: "padded state directory", stateDir: " " + filepath.Join(t.TempDir(), "bad-vm") + " "},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "bad-"+tt.name, session.ID, tt.stateDir, 1234, now()); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	stateDir := filepath.Join(t.TempDir(), "vm-1")
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, stateDir, 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	for _, tt := range []struct {
		name     string
		stateDir string
	}{
		{name: "blank state directory", stateDir: "\n\t"},
		{name: "padded state directory", stateDir: "\t" + filepath.Join(t.TempDir(), "bad-vm") + "\n"},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE vms SET state_dir = ? WHERE id = ?", tt.stateDir, "vm-1"); err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
		})
	}

	var gotStateDir string
	row := st.db.QueryRowContext(ctx, "SELECT state_dir FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotStateDir); err != nil {
		t.Fatalf("query VM state_dir: %v", err)
	}
	if gotStateDir != stateDir {
		t.Fatalf("VM state_dir = %q, want %q", gotStateDir, stateDir)
	}
}

func TestEnsureSchemaEnforcesVMFirecrackerPIDValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	for _, tt := range []struct {
		name  string
		fcPID any
	}{
		{name: "null PID", fcPID: nil},
		{name: "zero PID", fcPID: 0},
		{name: "negative PID", fcPID: -1},
	} {
		t.Run("insert "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "bad-"+tt.name, session.ID, filepath.Join(t.TempDir(), tt.name), tt.fcPID, now()); err == nil {
				t.Fatalf("inserted VM with %s, want trigger error", tt.name)
			}
		})
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	for _, tt := range []struct {
		name  string
		fcPID any
	}{
		{name: "null PID", fcPID: nil},
		{name: "zero PID", fcPID: 0},
		{name: "negative PID", fcPID: -1},
	} {
		t.Run("update "+tt.name, func(t *testing.T) {
			if _, err := st.db.ExecContext(ctx, "UPDATE vms SET fc_pid = ? WHERE id = ?", tt.fcPID, "vm-1"); err == nil {
				t.Fatalf("updated VM to %s, want trigger error", tt.name)
			}
		})
	}

	var fcPID int
	row := st.db.QueryRowContext(ctx, "SELECT fc_pid FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&fcPID); err != nil {
		t.Fatalf("query VM fc_pid: %v", err)
	}
	if fcPID != 1234 {
		t.Fatalf("VM fc_pid = %d, want 1234", fcPID)
	}
}

func TestEnsureSchemaEnforcesVMExitStatusValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at, exit_status)
VALUES(?, ?, ?, ?, ?, ?)`, "bad-exit-status-vm", session.ID, filepath.Join(t.TempDir(), "bad-vm"), 1234, now(), -1); err == nil {
		t.Fatalf("inserted VM with negative exit status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET exit_status = ? WHERE id = ?", -1, "vm-1"); err == nil {
		t.Fatalf("updated VM to negative exit status, want trigger error")
	}

	var exitStatus sql.NullInt64
	row := st.db.QueryRowContext(ctx, "SELECT exit_status FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&exitStatus); err != nil {
		t.Fatalf("query VM exit_status: %v", err)
	}
	if exitStatus.Valid {
		t.Fatalf("VM exit_status = %d, want NULL", exitStatus.Int64)
	}
}

func TestEnsureSchemaEnforcesVMCompletionConsistency(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at, ended_at)
VALUES(?, ?, ?, ?, ?, ?)`, "ended-without-exit-vm", session.ID, filepath.Join(t.TempDir(), "ended-without-exit-vm"), 1234, now(), now()); err == nil {
		t.Fatalf("inserted ended VM without exit status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at, exit_status)
VALUES(?, ?, ?, ?, ?, ?)`, "exit-without-ended-vm", session.ID, filepath.Join(t.TempDir(), "exit-without-ended-vm"), 1234, now(), 0); err == nil {
		t.Fatalf("inserted active VM with exit status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid active VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ? WHERE id = ?", now(), "vm-1"); err == nil {
		t.Fatalf("updated VM ended_at without exit status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET exit_status = ? WHERE id = ?", 0, "vm-1"); err == nil {
		t.Fatalf("updated VM exit_status without ended_at, want trigger error")
	}

	endedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", endedAt, 0, "vm-1"); err != nil {
		t.Fatalf("complete VM consistently: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = NULL WHERE id = ?", "vm-1"); err == nil {
		t.Fatalf("cleared ended_at without exit status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET exit_status = NULL WHERE id = ?", "vm-1"); err == nil {
		t.Fatalf("cleared exit_status without ended_at, want trigger error")
	}

	var gotEndedAt sql.NullString
	var gotExitStatus sql.NullInt64
	row := st.db.QueryRowContext(ctx, "SELECT ended_at, exit_status FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotEndedAt, &gotExitStatus); err != nil {
		t.Fatalf("query VM completion: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != endedAt {
		t.Fatalf("VM ended_at = %v, want %q", gotEndedAt, endedAt)
	}
	if !gotExitStatus.Valid || gotExitStatus.Int64 != 0 {
		t.Fatalf("VM exit_status = %v, want 0", gotExitStatus)
	}
}

func TestEnsureSchemaEnforcesSessionCompletionConsistency(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, ended_at, status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, "active-ended-session", userID, testKeyFingerprint, "127.0.0.1:2222", now(), now(), "active"); err == nil {
		t.Fatalf("inserted active session with ended_at, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "closed-without-ended-session", userID, testKeyFingerprint, "127.0.0.1:2222", now(), "closed"); err == nil {
		t.Fatalf("inserted closed session without ended_at, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "failed-without-ended-session", userID, testKeyFingerprint, "127.0.0.1:2223", now(), "vm_failed"); err == nil {
		t.Fatalf("inserted vm_failed session without ended_at, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2224", now(), "active"); err != nil {
		t.Fatalf("insert valid active session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ? WHERE id = ?", now(), "session-1"); err == nil {
		t.Fatalf("updated active session ended_at without terminal status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET status = ? WHERE id = ?", "closed", "session-1"); err == nil {
		t.Fatalf("updated session to terminal status without ended_at, want trigger error")
	}

	endedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", endedAt, "closed", "session-1"); err != nil {
		t.Fatalf("complete session consistently: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = NULL WHERE id = ?", "session-1"); err == nil {
		t.Fatalf("cleared ended_at for terminal session, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET status = ? WHERE id = ?", "active", "session-1"); err == nil {
		t.Fatalf("reopened session without clearing ended_at, want trigger error")
	}

	var gotEndedAt sql.NullString
	var gotStatus string
	row := st.db.QueryRowContext(ctx, "SELECT ended_at, status FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotEndedAt, &gotStatus); err != nil {
		t.Fatalf("query session completion: %v", err)
	}
	if !gotEndedAt.Valid || gotEndedAt.String != endedAt {
		t.Fatalf("session ended_at = %v, want %q", gotEndedAt, endedAt)
	}
	if gotStatus != "closed" {
		t.Fatalf("session status = %q, want closed", gotStatus)
	}
}

func TestEnsureSchemaEnforcesSessionVMCompletionOrdering(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid active session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", "session-1", filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert valid active VM: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", now(), "closed", "session-1"); err == nil {
		t.Fatalf("ended session with active VM, want trigger error")
	}

	vmEndedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ?", vmEndedAt, 0, "vm-1"); err != nil {
		t.Fatalf("complete VM consistently: %v", err)
	}
	sessionEndedAt := now()
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET ended_at = ?, status = ? WHERE id = ?", sessionEndedAt, "closed", "session-1"); err != nil {
		t.Fatalf("complete session after VM ended: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET ended_at = NULL, exit_status = NULL WHERE id = ?", "vm-1"); err == nil {
		t.Fatalf("reopened VM for terminal session, want trigger error")
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, ended_at, status)
VALUES(?, ?, ?, ?, ?, ?, ?)`, "session-2", userID, testKeyFingerprint, "127.0.0.1:2223", now(), now(), "vm_failed"); err != nil {
		t.Fatalf("insert valid terminal session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-2", "session-2", filepath.Join(t.TempDir(), "vm-2"), 1235, now()); err == nil {
		t.Fatalf("inserted VM for terminal session, want trigger error")
	}

	var (
		gotSessionStatus string
		gotSessionEnd    sql.NullString
		gotVMEnd         sql.NullString
		gotExitStatus    sql.NullInt64
		vmCount          int
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&gotSessionStatus, &gotSessionEnd); err != nil {
		t.Fatalf("query session-1: %v", err)
	}
	if gotSessionStatus != "closed" || !gotSessionEnd.Valid || gotSessionEnd.String != sessionEndedAt {
		t.Fatalf("session-1 completion = (%q, %v), want (closed, %q)", gotSessionStatus, gotSessionEnd, sessionEndedAt)
	}
	row = st.db.QueryRowContext(ctx, "SELECT ended_at, exit_status FROM vms WHERE id = ?", "vm-1")
	if err := row.Scan(&gotVMEnd, &gotExitStatus); err != nil {
		t.Fatalf("query vm-1: %v", err)
	}
	if !gotVMEnd.Valid || gotVMEnd.String != vmEndedAt {
		t.Fatalf("vm-1 ended_at = %v, want %q", gotVMEnd, vmEndedAt)
	}
	if !gotExitStatus.Valid || gotExitStatus.Int64 != 0 {
		t.Fatalf("vm-1 exit_status = %v, want 0", gotExitStatus)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms WHERE session_id = ?", "session-2")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query session-2 VMs: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("terminal session VM count = %d, want 0", vmCount)
	}
}

func TestEnsureSchemaEnforcesSessionStatusValues(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "bad-status-session", userID, testKeyFingerprint, "127.0.0.1:2222", now(), "paused"); err == nil {
		t.Fatalf("inserted session with invalid status, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", userID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET status = ? WHERE id = ?", "paused", "session-1"); err == nil {
		t.Fatalf("updated session to invalid status, want trigger error")
	}

	var status string
	row := st.db.QueryRowContext(ctx, "SELECT status FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&status); err != nil {
		t.Fatalf("query session status: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
}

func TestEnsureSchemaEnforcesSessionKeyOwnership(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	aliceID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey alice: %v", err)
	}
	bobID, err := st.EnsureUserAndKey(ctx, "bob", testOtherKeyFingerprint, testOtherAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey bob: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "cross-key-session", aliceID, testOtherKeyFingerprint, "127.0.0.1:2222", now(), "active"); err == nil {
		t.Fatalf("inserted session with another user's key, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
VALUES(?, ?, ?, ?, ?, ?)`, "session-1", aliceID, testKeyFingerprint, "127.0.0.1:2222", now(), "active"); err != nil {
		t.Fatalf("insert valid session: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET key_fingerprint = ? WHERE id = ?", testOtherKeyFingerprint, "session-1"); err == nil {
		t.Fatalf("updated session to another user's key, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE keys SET user_id = ? WHERE fingerprint = ?", bobID, testKeyFingerprint); err == nil {
		t.Fatalf("reassigned referenced key to another user, want trigger error")
	}

	var sessionUserID, sessionKeyFingerprint string
	row := st.db.QueryRowContext(ctx, "SELECT user_id, key_fingerprint FROM sessions WHERE id = ?", "session-1")
	if err := row.Scan(&sessionUserID, &sessionKeyFingerprint); err != nil {
		t.Fatalf("query session: %v", err)
	}
	if sessionUserID != aliceID || sessionKeyFingerprint != testKeyFingerprint {
		t.Fatalf("session ownership = (%q, %q), want (%q, %q)", sessionUserID, sessionKeyFingerprint, aliceID, testKeyFingerprint)
	}

	var keyUserID string
	row = st.db.QueryRowContext(ctx, "SELECT user_id FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&keyUserID); err != nil {
		t.Fatalf("query key: %v", err)
	}
	if keyUserID != aliceID {
		t.Fatalf("key user_id = %q, want %q", keyUserID, aliceID)
	}
}

func TestEnsureSchemaEnforcesUniqueUsernames(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-1", "alice", now(), now()); err != nil {
		t.Fatalf("insert first user: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", "user-2", "alice", now(), now()); err == nil {
		t.Fatalf("inserted duplicate username, want unique constraint error")
	}
}

func TestEnsureSchemaEnforcesOneVMPerSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-1", session.ID, filepath.Join(t.TempDir(), "vm-1"), 1234, now()); err != nil {
		t.Fatalf("insert first VM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-2", session.ID, filepath.Join(t.TempDir(), "vm-2"), 1235, now()); err == nil {
		t.Fatalf("inserted duplicate VM session, want unique constraint error")
	}
}

func TestEnsureSchemaEnforcesSessionVMReference(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	sessions := []Session{
		{
			ID:             "session-1",
			UserID:         userID,
			KeyFingerprint: testKeyFingerprint,
			RemoteAddr:     "127.0.0.1:2222",
			StartedAt:      now(),
			Status:         "active",
		},
		{
			ID:             "session-2",
			UserID:         userID,
			KeyFingerprint: testKeyFingerprint,
			RemoteAddr:     "127.0.0.1:2223",
			StartedAt:      now(),
			Status:         "active",
		},
	}
	for _, session := range sessions {
		if err := st.CreateSession(ctx, session); err != nil {
			t.Fatalf("CreateSession %s: %v", session.ID, err)
		}
	}
	if _, err := st.db.ExecContext(ctx, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
VALUES(?, ?, ?, ?, ?)`, "vm-2", sessions[1].ID, filepath.Join(t.TempDir(), "vm-2"), 1234, now()); err != nil {
		t.Fatalf("insert VM fixture: %v", err)
	}

	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", "missing-vm", sessions[0].ID); err == nil {
		t.Fatalf("updated session with dangling vm_id, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", "vm-2", sessions[0].ID); err == nil {
		t.Fatalf("updated session with cross-session vm_id, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", "vm-2", sessions[1].ID); err != nil {
		t.Fatalf("update same-session vm_id: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET session_id = ? WHERE id = ?", sessions[0].ID, "vm-2"); err == nil {
		t.Fatalf("reassigned linked VM to another session, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE vms SET id = ? WHERE id = ?", "vm-renamed", "vm-2"); err == nil {
		t.Fatalf("renamed linked VM, want trigger error")
	}
	if _, err := st.db.ExecContext(ctx, "DELETE FROM vms WHERE id = ?", "vm-2"); err == nil {
		t.Fatalf("deleted linked VM, want trigger error")
	}

	for _, session := range sessions {
		var attachedVM sql.NullString
		row := st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
		if err := row.Scan(&attachedVM); err != nil {
			t.Fatalf("query session %s vm_id: %v", session.ID, err)
		}
		if session.ID == sessions[0].ID {
			if attachedVM.Valid {
				t.Fatalf("rejected vm_id update mutated session %s to %q", session.ID, attachedVM.String)
			}
			continue
		}
		if !attachedVM.Valid || attachedVM.String != "vm-2" {
			t.Fatalf("session %s vm_id = %v, want vm-2", session.ID, attachedVM)
		}
	}
}

func TestStoreUsesSingleConnectionForConnectionScopedPragmas(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if maxOpen := st.db.Stats().MaxOpenConnections; maxOpen != 1 {
		t.Fatalf("MaxOpenConnections = %d, want 1", maxOpen)
	}

	var foreignKeys int
	row := st.db.QueryRowContext(ctx, "PRAGMA foreign_keys")
	if err := row.Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys pragma: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("PRAGMA foreign_keys = %d, want 1", foreignKeys)
	}
}

func TestNewRejectsBlankPath(t *testing.T) {
	for _, path := range []string{"", " \t\n"} {
		st, err := New(path)
		if err == nil {
			if st != nil {
				_ = st.Close()
			}
			t.Fatalf("New(%q) error = nil, want validation error", path)
		}
		if st != nil {
			t.Fatalf("New(%q) store = %#v, want nil", path, st)
		}
	}
}

func TestNewRejectsPathWithSurroundingWhitespaceBeforeSideEffects(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	st, err := New(" test.sqlite ")
	if err == nil {
		if st != nil {
			_ = st.Close()
		}
		t.Fatalf("New accepted database path with surrounding whitespace")
	}
	if st != nil {
		t.Fatalf("New store = %#v, want nil", st)
	}
	if err.Error() != "database path must not contain surrounding whitespace" {
		t.Fatalf("New error = %q, want surrounding whitespace validation error", err)
	}
	entries, readErr := os.ReadDir(workDir)
	if readErr != nil {
		t.Fatalf("read work dir: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("New created filesystem entries before validating path: %v", entries)
	}
}

func TestStoreAPIsRejectNilStore(t *testing.T) {
	ctx := context.Background()
	validSession := Session{
		ID:             "session-1",
		UserID:         "user-1",
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	validVM := VM{
		ID:        "vm-1",
		SessionID: validSession.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	tests := []struct {
		name string
		run  func(*Store) error
	}{
		{
			name: "Close",
			run: func(st *Store) error {
				return st.Close()
			},
		},
		{
			name: "EnsureSchema",
			run: func(st *Store) error {
				return st.EnsureSchema(ctx)
			},
		},
		{
			name: "HasKey",
			run: func(st *Store) error {
				_, err := st.HasKey(ctx, testKeyFingerprint)
				return err
			},
		},
		{
			name: "EnsureUserAndKey",
			run: func(st *Store) error {
				_, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
				return err
			},
		},
		{
			name: "CreateSession",
			run: func(st *Store) error {
				return st.CreateSession(ctx, validSession)
			},
		},
		{
			name: "EndSession",
			run: func(st *Store) error {
				return st.EndSession(ctx, validSession.ID, "closed")
			},
		},
		{
			name: "AttachVM",
			run: func(st *Store) error {
				return st.AttachVM(ctx, validSession.ID, validVM.ID)
			},
		},
		{
			name: "CreateVM",
			run: func(st *Store) error {
				return st.CreateVM(ctx, validVM)
			},
		},
		{
			name: "EndVM",
			run: func(st *Store) error {
				return st.EndVM(ctx, validVM.ID, 0)
			},
		},
		{
			name: "Audit",
			run: func(st *Store) error {
				return st.Audit(ctx, "test.audit", `{"ok":true}`)
			},
		},
	}

	stores := []struct {
		name string
		st   *Store
	}{
		{name: "nil receiver", st: nil},
		{name: "nil database", st: &Store{}},
	}
	for _, storeCase := range stores {
		for _, tt := range tests {
			t.Run(storeCase.name+" "+tt.name, func(t *testing.T) {
				st := storeCase.st
				if err := tt.run(st); err == nil {
					t.Fatalf("%s accepted %s", tt.name, storeCase.name)
				}
			})
		}
	}
}

func TestStoreAPIsRejectNilContext(t *testing.T) {
	st := newTestStore(t)
	validSession := Session{
		ID:             "session-1",
		UserID:         "user-1",
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	validVM := VM{
		ID:        "vm-1",
		SessionID: validSession.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "EnsureSchema",
			run: func() error {
				return st.EnsureSchema(nil)
			},
		},
		{
			name: "HasKey",
			run: func() error {
				_, err := st.HasKey(nil, testKeyFingerprint)
				return err
			},
		},
		{
			name: "EnsureUserAndKey",
			run: func() error {
				_, err := st.EnsureUserAndKey(nil, "alice", testKeyFingerprint, testAuthorizedKey)
				return err
			},
		},
		{
			name: "CreateSession",
			run: func() error {
				return st.CreateSession(nil, validSession)
			},
		},
		{
			name: "EndSession",
			run: func() error {
				return st.EndSession(nil, validSession.ID, "closed")
			},
		},
		{
			name: "AttachVM",
			run: func() error {
				return st.AttachVM(nil, validSession.ID, validVM.ID)
			},
		},
		{
			name: "CreateVM",
			run: func() error {
				return st.CreateVM(nil, validVM)
			},
		},
		{
			name: "EndVM",
			run: func() error {
				return st.EndVM(nil, validVM.ID, 0)
			},
		},
		{
			name: "Audit",
			run: func() error {
				return st.Audit(nil, "test.audit", `{"ok":true}`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Fatalf("%s accepted nil context", tt.name)
			}
		})
	}
}

func TestStoreUserSessionAndVMLifecycle(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	if userID == "" {
		t.Fatalf("EnsureUserAndKey returned empty user ID")
	}

	userIDAgain, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("second EnsureUserAndKey: %v", err)
	}
	if userIDAgain != userID {
		t.Fatalf("second EnsureUserAndKey user ID = %q, want %q", userIDAgain, userID)
	}

	hasKey, err := st.HasKey(ctx, testKeyFingerprint)
	if err != nil {
		t.Fatalf("HasKey: %v", err)
	}
	if !hasKey {
		t.Fatalf("HasKey returned false for enrolled key")
	}

	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if err := st.AttachVM(ctx, session.ID, vm.ID); err != nil {
		t.Fatalf("AttachVM: %v", err)
	}
	if err := st.EndVM(ctx, vm.ID, 7); err != nil {
		t.Fatalf("EndVM: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession: %v", err)
	}
	if err := st.Audit(ctx, "session.closed", `{"session_id":"session-1"}`); err != nil {
		t.Fatalf("Audit: %v", err)
	}

	var (
		status     string
		attachedVM string
		sessionEnd sql.NullString
		vmEnd      sql.NullString
		exitStatus int
	)
	row := st.db.QueryRowContext(ctx, `
SELECT s.status, s.vm_id, s.ended_at, v.ended_at, v.exit_status
FROM sessions s
JOIN vms v ON v.id = s.vm_id
WHERE s.id = ?`, session.ID)
	if err := row.Scan(&status, &attachedVM, &sessionEnd, &vmEnd, &exitStatus); err != nil {
		t.Fatalf("query lifecycle record: %v", err)
	}
	if status != "closed" {
		t.Fatalf("session status = %q, want closed", status)
	}
	if attachedVM != vm.ID {
		t.Fatalf("attached VM = %q, want %q", attachedVM, vm.ID)
	}
	if !sessionEnd.Valid {
		t.Fatalf("session ended_at was not set")
	}
	if !vmEnd.Valid {
		t.Fatalf("vm ended_at was not set")
	}
	if exitStatus != 7 {
		t.Fatalf("vm exit_status = %d, want 7", exitStatus)
	}

	var auditCount int
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE event_type = ?", "session.closed")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query audit_events: %v", err)
	}
	if auditCount != 1 {
		t.Fatalf("audit event count = %d, want 1", auditCount)
	}
}

func TestCreateSessionRequiresExistingUserAndKey(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	err := st.CreateSession(ctx, Session{
		ID:             "session-with-missing-key",
		UserID:         "missing-user",
		KeyFingerprint: "missing-key",
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	})
	if err == nil {
		t.Fatalf("CreateSession with missing user/key succeeded, want error")
	}
}

func TestCreateSessionRequiresKeyOwnedByUser(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	aliceID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey alice: %v", err)
	}
	if _, err := st.EnsureUserAndKey(ctx, "bob", testOtherKeyFingerprint, testOtherAuthorizedKey); err != nil {
		t.Fatalf("EnsureUserAndKey bob: %v", err)
	}

	err = st.CreateSession(ctx, Session{
		ID:             "cross-key-session",
		UserID:         aliceID,
		KeyFingerprint: testOtherKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	})
	if err != sql.ErrNoRows {
		t.Fatalf("CreateSession cross-user key error = %v, want sql.ErrNoRows", err)
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("cross-user key CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestCreateSessionRejectsBlankFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	valid := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	tests := []struct {
		name   string
		mutate func(*Session)
	}{
		{
			name: "blank ID",
			mutate: func(session *Session) {
				session.ID = " \t "
			},
		},
		{
			name: "blank user ID",
			mutate: func(session *Session) {
				session.UserID = " \t "
			},
		},
		{
			name: "blank key fingerprint",
			mutate: func(session *Session) {
				session.KeyFingerprint = " \t "
			},
		},
		{
			name: "blank remote address",
			mutate: func(session *Session) {
				session.RemoteAddr = " \t "
			},
		},
		{
			name: "blank start time",
			mutate: func(session *Session) {
				session.StartedAt = " \t "
			},
		},
		{
			name: "blank status",
			mutate: func(session *Session) {
				session.Status = " \t "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := valid
			tt.mutate(&session)
			if err := st.CreateSession(ctx, session); err == nil {
				t.Fatalf("CreateSession accepted %s", tt.name)
			}
		})
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("blank CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestCreateSessionRejectsWhitespacePaddedFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	valid := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	tests := []struct {
		name   string
		mutate func(*Session)
	}{
		{
			name: "padded ID",
			mutate: func(session *Session) {
				session.ID = " " + session.ID + " "
			},
		},
		{
			name: "padded user ID",
			mutate: func(session *Session) {
				session.UserID = " " + session.UserID + " "
			},
		},
		{
			name: "padded key fingerprint",
			mutate: func(session *Session) {
				session.KeyFingerprint = " " + session.KeyFingerprint + " "
			},
		},
		{
			name: "padded remote address",
			mutate: func(session *Session) {
				session.RemoteAddr = " " + session.RemoteAddr + " "
			},
		},
		{
			name: "padded start time",
			mutate: func(session *Session) {
				session.StartedAt = " " + session.StartedAt + " "
			},
		},
		{
			name: "padded status",
			mutate: func(session *Session) {
				session.Status = " " + session.Status + " "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := valid
			tt.mutate(&session)
			if err := st.CreateSession(ctx, session); err == nil {
				t.Fatalf("CreateSession accepted %s", tt.name)
			} else if err == sql.ErrNoRows {
				t.Fatalf("CreateSession returned sql.ErrNoRows for %s, want validation error", tt.name)
			}
		})
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("padded CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestCreateSessionRejectsInvalidStartTime(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      "2026-05-03 14:03:32",
		Status:         "active",
	}
	err = st.CreateSession(ctx, session)
	if err == nil {
		t.Fatalf("CreateSession accepted invalid start time")
	}
	if !strings.Contains(err.Error(), "session start time must be a valid RFC3339 timestamp") {
		t.Fatalf("CreateSession error = %q, want RFC3339 timestamp validation error", err)
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("invalid start time CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestCreateSessionRejectsInvalidRemoteAddr(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	valid := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	tests := []struct {
		name       string
		remoteAddr string
		wantErr    string
	}{
		{
			name:       "missing port",
			remoteAddr: "127.0.0.1",
			wantErr:    "session remote address must be a valid TCP address",
		},
		{
			name:       "missing host",
			remoteAddr: ":2222",
			wantErr:    "session remote address host must be set",
		},
		{
			name:       "invalid port",
			remoteAddr: "127.0.0.1:not-a-port",
			wantErr:    "session remote address port must be valid",
		},
		{
			name:       "zero port",
			remoteAddr: "127.0.0.1:0",
			wantErr:    "session remote address port must be > 0",
		},
		{
			name:       "unresolvable address",
			remoteAddr: "999.0.0.1:2222",
			wantErr:    "session remote address must resolve to a valid TCP address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := valid
			session.ID = tt.name
			session.RemoteAddr = tt.remoteAddr
			err := st.CreateSession(ctx, session)
			if err == nil {
				t.Fatalf("CreateSession accepted invalid remote address %q", tt.remoteAddr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("CreateSession error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("invalid remote address CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestCreateSessionRejectsUnsupportedStatus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "closed",
	}
	err = st.CreateSession(ctx, session)
	if err == nil {
		t.Fatalf("CreateSession accepted unsupported status")
	}
	if !strings.Contains(err.Error(), "session status must be active") {
		t.Fatalf("CreateSession error = %q, want initial status validation error", err)
	}

	var sessionCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sessions")
	if err := row.Scan(&sessionCount); err != nil {
		t.Fatalf("query sessions: %v", err)
	}
	if sessionCount != 0 {
		t.Fatalf("unsupported status CreateSession inserted sessions=%d, want 0", sessionCount)
	}
}

func TestEnsureUserAndKeyRejectsBlankInputs(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	tests := []struct {
		name        string
		username    string
		fingerprint string
		publicKey   string
	}{
		{
			name:        "blank username",
			username:    " \t ",
			fingerprint: testKeyFingerprint,
			publicKey:   testAuthorizedKey,
		},
		{
			name:        "blank fingerprint",
			username:    "alice",
			fingerprint: " \t ",
			publicKey:   testAuthorizedKey,
		},
		{
			name:        "blank public key",
			username:    "alice",
			fingerprint: testKeyFingerprint,
			publicKey:   " \t ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := st.EnsureUserAndKey(ctx, tt.username, tt.fingerprint, tt.publicKey); err == nil {
				t.Fatalf("EnsureUserAndKey accepted %s", tt.name)
			}
		})
	}

	var userCount, keyCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err != nil {
		t.Fatalf("query users: %v", err)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM keys")
	if err := row.Scan(&keyCount); err != nil {
		t.Fatalf("query keys: %v", err)
	}
	if userCount != 0 || keyCount != 0 {
		t.Fatalf("blank EnsureUserAndKey inserted users=%d keys=%d, want 0/0", userCount, keyCount)
	}
}

func TestEnsureUserAndKeyRejectsWhitespacePaddedMetadata(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	tests := []struct {
		name        string
		username    string
		fingerprint string
		publicKey   string
	}{
		{
			name:        "padded username",
			username:    " alice ",
			fingerprint: testKeyFingerprint,
		},
		{
			name:        "padded fingerprint",
			username:    "alice",
			fingerprint: " " + testKeyFingerprint + " ",
		},
		{
			name:        "padded public key",
			username:    "alice",
			fingerprint: testKeyFingerprint,
			publicKey:   " " + testAuthorizedKey + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicKey := tt.publicKey
			if publicKey == "" {
				publicKey = testAuthorizedKey
			}
			if _, err := st.EnsureUserAndKey(ctx, tt.username, tt.fingerprint, publicKey); err == nil {
				t.Fatalf("EnsureUserAndKey accepted %s", tt.name)
			}
		})
	}

	var userCount, keyCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err != nil {
		t.Fatalf("query users: %v", err)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM keys")
	if err := row.Scan(&keyCount); err != nil {
		t.Fatalf("query keys: %v", err)
	}
	if userCount != 0 || keyCount != 0 {
		t.Fatalf("padded EnsureUserAndKey inserted users=%d keys=%d, want 0/0", userCount, keyCount)
	}
}

func TestEnsureUserAndKeyRejectsInvalidPublicKey(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	tests := []struct {
		name      string
		publicKey string
		wantErr   string
	}{
		{
			name:      "malformed public key",
			publicKey: "ssh-ed25519 AAAA alice",
			wantErr:   "public key must be a valid authorized key",
		},
		{
			name:      "multiple public keys",
			publicKey: testAuthorizedKey + "\n" + testAuthorizedKey,
			wantErr:   "public key must contain exactly one authorized key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, tt.publicKey)
			if err == nil {
				t.Fatalf("EnsureUserAndKey accepted %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("EnsureUserAndKey error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}

	var userCount, keyCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err != nil {
		t.Fatalf("query users: %v", err)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM keys")
	if err := row.Scan(&keyCount); err != nil {
		t.Fatalf("query keys: %v", err)
	}
	if userCount != 0 || keyCount != 0 {
		t.Fatalf("invalid public key inserted users=%d keys=%d, want 0/0", userCount, keyCount)
	}
}

func TestEnsureUserAndKeyRejectsFingerprintMismatch(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	_, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:mismatch", testAuthorizedKey)
	if err == nil {
		t.Fatalf("EnsureUserAndKey accepted mismatched fingerprint")
	}
	if !strings.Contains(err.Error(), "key fingerprint must match public key") {
		t.Fatalf("EnsureUserAndKey error = %q, want fingerprint mismatch", err)
	}

	var userCount, keyCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err != nil {
		t.Fatalf("query users: %v", err)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM keys")
	if err := row.Scan(&keyCount); err != nil {
		t.Fatalf("query keys: %v", err)
	}
	if userCount != 0 || keyCount != 0 {
		t.Fatalf("fingerprint mismatch inserted users=%d keys=%d, want 0/0", userCount, keyCount)
	}
}

func TestEnsureUserAndKeyRejectsExistingFingerprintForDifferentUser(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	aliceID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey alice: %v", err)
	}

	bobID, err := st.EnsureUserAndKey(ctx, "bob", testKeyFingerprint, testAuthorizedKey)
	if err == nil {
		t.Fatalf("EnsureUserAndKey reassigned existing fingerprint to bob, bobID=%q", bobID)
	}
	if !strings.Contains(err.Error(), "key fingerprint is already enrolled for another user") {
		t.Fatalf("EnsureUserAndKey error = %q, want existing fingerprint ownership error", err)
	}

	var (
		keyUserID string
		userCount int
	)
	row := st.db.QueryRowContext(ctx, "SELECT user_id FROM keys WHERE fingerprint = ?", testKeyFingerprint)
	if err := row.Scan(&keyUserID); err != nil {
		t.Fatalf("query key owner: %v", err)
	}
	if keyUserID != aliceID {
		t.Fatalf("key user_id = %q, want original alice user ID %q", keyUserID, aliceID)
	}
	row = st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	if err := row.Scan(&userCount); err != nil {
		t.Fatalf("query users: %v", err)
	}
	if userCount != 1 {
		t.Fatalf("rejected cross-user key enrollment persisted users=%d, want 1", userCount)
	}
}

func TestHasKeyRejectsBlankFingerprint(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	hasKey, err := st.HasKey(ctx, " \t ")
	if err == nil {
		t.Fatalf("HasKey accepted a blank fingerprint")
	}
	if hasKey {
		t.Fatalf("HasKey returned true for a blank fingerprint")
	}
}

func TestHasKeyRejectsWhitespacePaddedFingerprint(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if _, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey); err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	hasKey, err := st.HasKey(ctx, " "+testKeyFingerprint+" ")
	if err == nil {
		t.Fatalf("HasKey accepted a whitespace-padded fingerprint")
	}
	if hasKey {
		t.Fatalf("HasKey returned true for a whitespace-padded fingerprint")
	}
}

func TestLifecycleUpdatesRequireExistingRecords(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.EndSession(ctx, "missing-session", "closed"); err != sql.ErrNoRows {
		t.Fatalf("EndSession missing record error = %v, want sql.ErrNoRows", err)
	}
	if err := st.EndVM(ctx, "missing-vm", 0); err != sql.ErrNoRows {
		t.Fatalf("EndVM missing record error = %v, want sql.ErrNoRows", err)
	}
	if err := st.AttachVM(ctx, "missing-session", "missing-vm"); err != sql.ErrNoRows {
		t.Fatalf("AttachVM missing records error = %v, want sql.ErrNoRows", err)
	}
}

func TestLifecycleUpdatesRejectBlankInputs(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "EndSession blank session ID",
			run: func() error {
				return st.EndSession(ctx, " \t ", "closed")
			},
		},
		{
			name: "EndSession blank status",
			run: func() error {
				return st.EndSession(ctx, session.ID, " \t ")
			},
		},
		{
			name: "AttachVM blank session ID",
			run: func() error {
				return st.AttachVM(ctx, " \t ", vm.ID)
			},
		},
		{
			name: "AttachVM blank VM ID",
			run: func() error {
				return st.AttachVM(ctx, session.ID, " \t ")
			},
		},
		{
			name: "EndVM blank VM ID",
			run: func() error {
				return st.EndVM(ctx, " \t ", 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Fatalf("%s succeeded", tt.name)
			} else if err == sql.ErrNoRows {
				t.Fatalf("%s returned sql.ErrNoRows, want validation error", tt.name)
			}
		})
	}

	var (
		status     string
		attachedVM sql.NullString
		sessionEnd sql.NullString
		vmEnd      sql.NullString
	)
	row := st.db.QueryRowContext(ctx, `
SELECT s.status, s.vm_id, s.ended_at, v.ended_at
FROM sessions s
JOIN vms v ON v.session_id = s.id
WHERE s.id = ?`, session.ID)
	if err := row.Scan(&status, &attachedVM, &sessionEnd, &vmEnd); err != nil {
		t.Fatalf("query lifecycle records: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if attachedVM.Valid {
		t.Fatalf("session vm_id = %q, want NULL", attachedVM.String)
	}
	if sessionEnd.Valid {
		t.Fatalf("session ended_at = %q, want NULL", sessionEnd.String)
	}
	if vmEnd.Valid {
		t.Fatalf("vm ended_at = %q, want NULL", vmEnd.String)
	}
}

func TestLifecycleUpdatesRejectWhitespacePaddedInputs(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}

	tests := []struct {
		name string
		run  func() error
	}{
		{
			name: "EndSession padded session ID",
			run: func() error {
				return st.EndSession(ctx, " "+session.ID+" ", "closed")
			},
		},
		{
			name: "EndSession padded status",
			run: func() error {
				return st.EndSession(ctx, session.ID, " closed ")
			},
		},
		{
			name: "AttachVM padded session ID",
			run: func() error {
				return st.AttachVM(ctx, " "+session.ID+" ", vm.ID)
			},
		},
		{
			name: "AttachVM padded VM ID",
			run: func() error {
				return st.AttachVM(ctx, session.ID, " "+vm.ID+" ")
			},
		},
		{
			name: "EndVM padded VM ID",
			run: func() error {
				return st.EndVM(ctx, " "+vm.ID+" ", 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Fatalf("%s succeeded", tt.name)
			} else if err == sql.ErrNoRows {
				t.Fatalf("%s returned sql.ErrNoRows, want validation error", tt.name)
			}
		})
	}

	var (
		status     string
		attachedVM sql.NullString
		sessionEnd sql.NullString
		vmEnd      sql.NullString
	)
	row := st.db.QueryRowContext(ctx, `
SELECT s.status, s.vm_id, s.ended_at, v.ended_at
FROM sessions s
JOIN vms v ON v.session_id = s.id
WHERE s.id = ?`, session.ID)
	if err := row.Scan(&status, &attachedVM, &sessionEnd, &vmEnd); err != nil {
		t.Fatalf("query lifecycle records: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if attachedVM.Valid {
		t.Fatalf("session vm_id = %q, want NULL", attachedVM.String)
	}
	if sessionEnd.Valid {
		t.Fatalf("session ended_at = %q, want NULL", sessionEnd.String)
	}
	if vmEnd.Valid {
		t.Fatalf("vm ended_at = %q, want NULL", vmEnd.String)
	}
}

func TestEndVMRejectsNegativeExitStatus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}

	err = st.EndVM(ctx, vm.ID, -1)
	if err == nil {
		t.Fatalf("EndVM accepted a negative exit status")
	}
	if !strings.Contains(err.Error(), "VM exit status must be >= 0") {
		t.Fatalf("EndVM error = %q, want exit status validation error", err)
	}

	var (
		endedAt    sql.NullString
		exitStatus sql.NullInt64
	)
	row := st.db.QueryRowContext(ctx, "SELECT ended_at, exit_status FROM vms WHERE id = ?", vm.ID)
	if err := row.Scan(&endedAt, &exitStatus); err != nil {
		t.Fatalf("query VM lifecycle state: %v", err)
	}
	if endedAt.Valid {
		t.Fatalf("EndVM set ended_at for invalid exit status: %q", endedAt.String)
	}
	if exitStatus.Valid {
		t.Fatalf("EndVM set exit_status for invalid exit status: %d", exitStatus.Int64)
	}
}

func TestEndVMRejectsAlreadyEndedVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if err := st.EndVM(ctx, vm.ID, 7); err != nil {
		t.Fatalf("EndVM: %v", err)
	}

	err = st.EndVM(ctx, vm.ID, 9)
	if err != sql.ErrNoRows {
		t.Fatalf("second EndVM error = %v, want sql.ErrNoRows", err)
	}

	var (
		endedAt    string
		exitStatus int
	)
	row := st.db.QueryRowContext(ctx, "SELECT ended_at, exit_status FROM vms WHERE id = ?", vm.ID)
	if err := row.Scan(&endedAt, &exitStatus); err != nil {
		t.Fatalf("query VM lifecycle state: %v", err)
	}
	if endedAt == "" {
		t.Fatalf("VM ended_at was not set by first EndVM")
	}
	if exitStatus != 7 {
		t.Fatalf("VM exit_status = %d, want 7", exitStatus)
	}
}

func TestEndSessionRejectsUnsupportedTerminalStatus(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	err = st.EndSession(ctx, session.ID, "active")
	if err == nil {
		t.Fatalf("EndSession accepted unsupported terminal status")
	}
	if !strings.Contains(err.Error(), "session end status must be closed or vm_failed") {
		t.Fatalf("EndSession error = %q, want terminal status validation error", err)
	}

	var (
		status  string
		endedAt sql.NullString
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&status, &endedAt); err != nil {
		t.Fatalf("query session lifecycle state: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if endedAt.Valid {
		t.Fatalf("EndSession set ended_at for invalid terminal status: %q", endedAt.String)
	}
}

func TestEndSessionRejectsAlreadyEndedSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession: %v", err)
	}

	err = st.EndSession(ctx, session.ID, "vm_failed")
	if err != sql.ErrNoRows {
		t.Fatalf("second EndSession error = %v, want sql.ErrNoRows", err)
	}

	var (
		status  string
		endedAt string
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&status, &endedAt); err != nil {
		t.Fatalf("query session lifecycle state: %v", err)
	}
	if status != "closed" {
		t.Fatalf("session status = %q, want closed", status)
	}
	if endedAt == "" {
		t.Fatalf("session ended_at was not set by first EndSession")
	}
}

func TestEndSessionRejectsClosedStatusWithActiveAttachedVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if err := st.AttachVM(ctx, session.ID, vm.ID); err != nil {
		t.Fatalf("AttachVM: %v", err)
	}

	err = st.EndSession(ctx, session.ID, "closed")
	if err != sql.ErrNoRows {
		t.Fatalf("EndSession with active attached VM error = %v, want sql.ErrNoRows", err)
	}

	var (
		status     string
		endedAt    sql.NullString
		attachedVM string
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at, vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&status, &endedAt, &attachedVM); err != nil {
		t.Fatalf("query session lifecycle state: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if endedAt.Valid {
		t.Fatalf("EndSession set ended_at while attached VM was active: %q", endedAt.String)
	}
	if attachedVM != vm.ID {
		t.Fatalf("session vm_id = %q, want attached VM %q", attachedVM, vm.ID)
	}

	if err := st.EndVM(ctx, vm.ID, 0); err != nil {
		t.Fatalf("EndVM: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession after VM ended: %v", err)
	}
}

func TestEndSessionRejectsVMFailedStatusWithActiveAttachedVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if err := st.AttachVM(ctx, session.ID, vm.ID); err != nil {
		t.Fatalf("AttachVM: %v", err)
	}

	err = st.EndSession(ctx, session.ID, "vm_failed")
	if err != sql.ErrNoRows {
		t.Fatalf("EndSession vm_failed with active attached VM error = %v, want sql.ErrNoRows", err)
	}

	var (
		status     string
		endedAt    sql.NullString
		attachedVM string
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at, vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&status, &endedAt, &attachedVM); err != nil {
		t.Fatalf("query session lifecycle state: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if endedAt.Valid {
		t.Fatalf("EndSession set ended_at while attached VM was active: %q", endedAt.String)
	}
	if attachedVM != vm.ID {
		t.Fatalf("session vm_id = %q, want attached VM %q", attachedVM, vm.ID)
	}

	if err := st.EndVM(ctx, vm.ID, 1); err != nil {
		t.Fatalf("EndVM: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "vm_failed"); err != nil {
		t.Fatalf("EndSession after VM ended: %v", err)
	}
}

func TestEndSessionRejectsActiveUnattachedVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}

	err = st.EndSession(ctx, session.ID, "closed")
	if err != sql.ErrNoRows {
		t.Fatalf("EndSession with active unattached VM error = %v, want sql.ErrNoRows", err)
	}

	var (
		status     string
		endedAt    sql.NullString
		attachedVM sql.NullString
		vmEnd      sql.NullString
	)
	row := st.db.QueryRowContext(ctx, `
SELECT s.status, s.ended_at, s.vm_id, v.ended_at
FROM sessions s
JOIN vms v ON v.session_id = s.id
WHERE s.id = ?`, session.ID)
	if err := row.Scan(&status, &endedAt, &attachedVM, &vmEnd); err != nil {
		t.Fatalf("query lifecycle state: %v", err)
	}
	if status != "active" {
		t.Fatalf("session status = %q, want active", status)
	}
	if endedAt.Valid {
		t.Fatalf("EndSession set ended_at while unattached VM was active: %q", endedAt.String)
	}
	if attachedVM.Valid {
		t.Fatalf("session vm_id = %q, want NULL", attachedVM.String)
	}
	if vmEnd.Valid {
		t.Fatalf("VM ended_at = %q, want NULL", vmEnd.String)
	}

	if err := st.EndVM(ctx, vm.ID, 0); err != nil {
		t.Fatalf("EndVM: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession after VM ended: %v", err)
	}
}

func TestAttachVMRequiresExistingVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	if err := st.AttachVM(ctx, session.ID, "missing-vm"); err != sql.ErrNoRows {
		t.Fatalf("AttachVM missing VM error = %v, want sql.ErrNoRows", err)
	}

	var attachedVM sql.NullString
	row := st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&attachedVM); err != nil {
		t.Fatalf("query session vm_id: %v", err)
	}
	if attachedVM.Valid {
		t.Fatalf("AttachVM set vm_id to %q for missing VM", attachedVM.String)
	}
}

func TestAttachVMRequiresVMForSameSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	sessions := []Session{
		{
			ID:             "session-1",
			UserID:         userID,
			KeyFingerprint: testKeyFingerprint,
			RemoteAddr:     "127.0.0.1:2222",
			StartedAt:      now(),
			Status:         "active",
		},
		{
			ID:             "session-2",
			UserID:         userID,
			KeyFingerprint: testKeyFingerprint,
			RemoteAddr:     "127.0.0.1:2223",
			StartedAt:      now(),
			Status:         "active",
		},
	}
	for _, session := range sessions {
		if err := st.CreateSession(ctx, session); err != nil {
			t.Fatalf("CreateSession %s: %v", session.ID, err)
		}
	}
	vm := VM{
		ID:        "vm-2",
		SessionID: sessions[1].ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-2"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}

	if err := st.AttachVM(ctx, sessions[0].ID, vm.ID); err != sql.ErrNoRows {
		t.Fatalf("AttachVM cross-session VM error = %v, want sql.ErrNoRows", err)
	}

	for _, session := range sessions {
		var attachedVM sql.NullString
		row := st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
		if err := row.Scan(&attachedVM); err != nil {
			t.Fatalf("query session %s vm_id: %v", session.ID, err)
		}
		if attachedVM.Valid {
			t.Fatalf("session %s vm_id = %q, want NULL", session.ID, attachedVM.String)
		}
	}
}

func TestAttachVMRejectsAlreadyAttachedSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if _, err := st.db.ExecContext(ctx, "UPDATE sessions SET vm_id = ? WHERE id = ?", vm.ID, session.ID); err != nil {
		t.Fatalf("mark session attached fixture: %v", err)
	}

	err = st.AttachVM(ctx, session.ID, vm.ID)
	if err != sql.ErrNoRows {
		t.Fatalf("AttachVM already-attached session error = %v, want sql.ErrNoRows", err)
	}

	var attachedVM string
	row := st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&attachedVM); err != nil {
		t.Fatalf("query session vm_id: %v", err)
	}
	if attachedVM != vm.ID {
		t.Fatalf("session vm_id = %q, want original VM %q", attachedVM, vm.ID)
	}
}

func TestCreateVMRejectsSecondVMForSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM first VM: %v", err)
	}

	otherVM := VM{
		ID:        "vm-2",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-2"),
		FCPid:     1235,
		StartedAt: now(),
	}
	err = st.CreateVM(ctx, otherVM)
	if err != sql.ErrNoRows {
		t.Fatalf("CreateVM second VM error = %v, want sql.ErrNoRows", err)
	}

	var (
		vmCount    int
		attachedVM sql.NullString
	)
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms WHERE session_id = ?", session.ID)
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query session VMs: %v", err)
	}
	if vmCount != 1 {
		t.Fatalf("CreateVM inserted VMs for session=%d, want 1", vmCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&attachedVM); err != nil {
		t.Fatalf("query attached VM: %v", err)
	}
	if attachedVM.Valid {
		t.Fatalf("CreateVM attached VM unexpectedly to %q", attachedVM.String)
	}
}

func TestAttachVMRejectsEndedSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM: %v", err)
	}
	if err := st.EndVM(ctx, vm.ID, 0); err != nil {
		t.Fatalf("EndVM: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession: %v", err)
	}

	err = st.AttachVM(ctx, session.ID, vm.ID)
	if err != sql.ErrNoRows {
		t.Fatalf("AttachVM ended session error = %v, want sql.ErrNoRows", err)
	}

	var (
		status     string
		endedAt    sql.NullString
		attachedVM sql.NullString
	)
	row := st.db.QueryRowContext(ctx, "SELECT status, ended_at, vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&status, &endedAt, &attachedVM); err != nil {
		t.Fatalf("query session lifecycle state: %v", err)
	}
	if status != "closed" {
		t.Fatalf("session status = %q, want closed", status)
	}
	if !endedAt.Valid {
		t.Fatalf("session ended_at was not set by EndSession")
	}
	if attachedVM.Valid {
		t.Fatalf("AttachVM set vm_id on ended session to %q", attachedVM.String)
	}
}

func TestCreateVMRejectsEndedSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	if err := st.EndSession(ctx, session.ID, "closed"); err != nil {
		t.Fatalf("EndSession: %v", err)
	}

	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	err = st.CreateVM(ctx, vm)
	if err != sql.ErrNoRows {
		t.Fatalf("CreateVM ended session error = %v, want sql.ErrNoRows", err)
	}

	var vmCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("CreateVM inserted VMs for ended session=%d, want 0", vmCount)
	}
}

func TestCreateVMRejectsAlreadyAttachedSession(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}
	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	if err := st.CreateVM(ctx, vm); err != nil {
		t.Fatalf("CreateVM first VM: %v", err)
	}
	if err := st.AttachVM(ctx, session.ID, vm.ID); err != nil {
		t.Fatalf("AttachVM: %v", err)
	}

	otherVM := VM{
		ID:        "vm-2",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-2"),
		FCPid:     1235,
		StartedAt: now(),
	}
	err = st.CreateVM(ctx, otherVM)
	if err != sql.ErrNoRows {
		t.Fatalf("CreateVM already-attached session error = %v, want sql.ErrNoRows", err)
	}

	var (
		vmCount    int
		attachedVM string
	)
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 1 {
		t.Fatalf("CreateVM inserted VMs for already-attached session=%d, want 1", vmCount)
	}
	row = st.db.QueryRowContext(ctx, "SELECT vm_id FROM sessions WHERE id = ?", session.ID)
	if err := row.Scan(&attachedVM); err != nil {
		t.Fatalf("query attached VM: %v", err)
	}
	if attachedVM != vm.ID {
		t.Fatalf("session vm_id = %q, want original VM %q", attachedVM, vm.ID)
	}
}

func TestCreateVMRejectsBlankFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	valid := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	tests := []struct {
		name   string
		mutate func(*VM)
	}{
		{
			name: "blank ID",
			mutate: func(vm *VM) {
				vm.ID = " \t "
			},
		},
		{
			name: "blank session ID",
			mutate: func(vm *VM) {
				vm.SessionID = " \t "
			},
		},
		{
			name: "blank state directory",
			mutate: func(vm *VM) {
				vm.StateDir = " \t "
			},
		},
		{
			name: "blank start time",
			mutate: func(vm *VM) {
				vm.StartedAt = " \t "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := valid
			tt.mutate(&vm)
			if err := st.CreateVM(ctx, vm); err == nil {
				t.Fatalf("CreateVM accepted %s", tt.name)
			}
		})
	}

	var vmCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("blank CreateVM inserted VMs=%d, want 0", vmCount)
	}
}

func TestCreateVMRejectsWhitespacePaddedFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	valid := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	tests := []struct {
		name   string
		mutate func(*VM)
	}{
		{
			name: "padded ID",
			mutate: func(vm *VM) {
				vm.ID = " " + vm.ID + " "
			},
		},
		{
			name: "padded session ID",
			mutate: func(vm *VM) {
				vm.SessionID = " " + vm.SessionID + " "
			},
		},
		{
			name: "padded state directory",
			mutate: func(vm *VM) {
				vm.StateDir = " " + vm.StateDir + " "
			},
		},
		{
			name: "padded start time",
			mutate: func(vm *VM) {
				vm.StartedAt = " " + vm.StartedAt + " "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := valid
			tt.mutate(&vm)
			if err := st.CreateVM(ctx, vm); err == nil {
				t.Fatalf("CreateVM accepted %s", tt.name)
			} else if err == sql.ErrNoRows {
				t.Fatalf("CreateVM returned sql.ErrNoRows for %s, want validation error", tt.name)
			}
		})
	}

	var vmCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("whitespace-padded CreateVM inserted VMs=%d, want 0", vmCount)
	}
}

func TestCreateVMRejectsInvalidStartTime(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	vm := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: "2026-05-03 14:03:32",
	}
	err = st.CreateVM(ctx, vm)
	if err == nil {
		t.Fatalf("CreateVM accepted invalid start time")
	}
	if !strings.Contains(err.Error(), "VM start time must be a valid RFC3339 timestamp") {
		t.Fatalf("CreateVM error = %q, want RFC3339 timestamp validation error", err)
	}

	var vmCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("invalid start time CreateVM inserted VMs=%d, want 0", vmCount)
	}
}

func TestCreateVMRejectsNonPositiveFirecrackerPID(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", testKeyFingerprint, testAuthorizedKey)
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: testKeyFingerprint,
		RemoteAddr:     "127.0.0.1:2222",
		StartedAt:      now(),
		Status:         "active",
	}
	if err := st.CreateSession(ctx, session); err != nil {
		t.Fatalf("CreateSession: %v", err)
	}

	valid := VM{
		ID:        "vm-1",
		SessionID: session.ID,
		StateDir:  filepath.Join(t.TempDir(), "vm-1"),
		FCPid:     1234,
		StartedAt: now(),
	}
	tests := []struct {
		name  string
		fcPid int
	}{
		{name: "zero PID", fcPid: 0},
		{name: "negative PID", fcPid: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vm := valid
			vm.ID = tt.name
			vm.FCPid = tt.fcPid
			err := st.CreateVM(ctx, vm)
			if err == nil {
				t.Fatalf("CreateVM accepted %s", tt.name)
			}
			if !strings.Contains(err.Error(), "VM Firecracker PID must be > 0") {
				t.Fatalf("CreateVM error = %q, want Firecracker PID validation error", err)
			}
		})
	}

	var vmCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM vms")
	if err := row.Scan(&vmCount); err != nil {
		t.Fatalf("query vms: %v", err)
	}
	if vmCount != 0 {
		t.Fatalf("non-positive PID CreateVM inserted VMs=%d, want 0", vmCount)
	}
}

func TestAuditRequiresValidJSON(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.Audit(ctx, "bad.audit", `{"missing-close":`); err == nil {
		t.Fatalf("Audit accepted invalid JSON data")
	}

	var auditCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events WHERE event_type = ?", "bad.audit")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query invalid audit count: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("invalid audit event count = %d, want 0", auditCount)
	}

	if err := st.Audit(ctx, "good.audit", `{"ok":true}`); err != nil {
		t.Fatalf("Audit valid JSON: %v", err)
	}
}

func TestAuditRejectsBlankData(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	for _, data := range []string{"", " \t "} {
		t.Run(data, func(t *testing.T) {
			err := st.Audit(ctx, "blank.audit", data)
			if err == nil {
				t.Fatalf("Audit accepted blank data")
			}
			if !strings.Contains(err.Error(), "audit data must be set") {
				t.Fatalf("Audit error = %q, want required data validation error", err)
			}
		})
	}

	var auditCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query audit_events: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("blank audit data event count = %d, want 0", auditCount)
	}
}

func TestAuditRejectsWhitespacePaddedData(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	err := st.Audit(ctx, "padded.audit", ` {"ok":true} `)
	if err == nil {
		t.Fatalf("Audit accepted whitespace-padded data")
	}
	if !strings.Contains(err.Error(), "audit data must not contain surrounding whitespace") {
		t.Fatalf("Audit error = %q, want surrounding whitespace validation error", err)
	}

	var auditCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query audit_events: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("whitespace-padded audit data event count = %d, want 0", auditCount)
	}
}

func TestAuditRejectsBlankEventType(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.Audit(ctx, " \t ", `{"ok":true}`); err == nil {
		t.Fatalf("Audit accepted blank event type")
	}

	var auditCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query audit_events: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("blank audit event count = %d, want 0", auditCount)
	}
}

func TestAuditRejectsWhitespacePaddedEventType(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	if err := st.Audit(ctx, " test.audit ", `{"ok":true}`); err == nil {
		t.Fatalf("Audit accepted whitespace-padded event type")
	}

	var auditCount int
	row := st.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM audit_events")
	if err := row.Scan(&auditCount); err != nil {
		t.Fatalf("query audit_events: %v", err)
	}
	if auditCount != 0 {
		t.Fatalf("whitespace-padded audit event count = %d, want 0", auditCount)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()

	st, err := New(filepath.Join(t.TempDir(), "test.sqlite"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if err := st.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	if err := st.EnsureSchema(context.Background()); err != nil {
		t.Fatalf("EnsureSchema: %v", err)
	}
	return st
}
