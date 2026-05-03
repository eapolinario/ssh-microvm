package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
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

func TestStoreUserSessionAndVMLifecycle(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}
	if userID == "" {
		t.Fatalf("EnsureUserAndKey returned empty user ID")
	}

	userIDAgain, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("second EnsureUserAndKey: %v", err)
	}
	if userIDAgain != userID {
		t.Fatalf("second EnsureUserAndKey user ID = %q, want %q", userIDAgain, userID)
	}

	hasKey, err := st.HasKey(ctx, "SHA256:test")
	if err != nil {
		t.Fatalf("HasKey: %v", err)
	}
	if !hasKey {
		t.Fatalf("HasKey returned false for enrolled key")
	}

	session := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: "SHA256:test",
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
		t.Fatalf("CreateSession with missing user/key succeeded, want foreign key error")
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
