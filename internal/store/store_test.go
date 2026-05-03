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
