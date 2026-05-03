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
