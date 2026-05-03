package store

import (
	"context"
	"database/sql"
	"os"
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
		KeyFingerprint: "SHA256:test",
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
				_, err := st.HasKey(ctx, "SHA256:test")
				return err
			},
		},
		{
			name: "EnsureUserAndKey",
			run: func(st *Store) error {
				_, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
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
		KeyFingerprint: "SHA256:test",
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
				_, err := st.HasKey(nil, "SHA256:test")
				return err
			},
		},
		{
			name: "EnsureUserAndKey",
			run: func() error {
				_, err := st.EnsureUserAndKey(nil, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
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

func TestCreateSessionRejectsBlankFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	valid := Session{
		ID:             "session-1",
		UserID:         userID,
		KeyFingerprint: "SHA256:test",
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
			fingerprint: "SHA256:test",
			publicKey:   "ssh-ed25519 AAAA alice",
		},
		{
			name:        "blank fingerprint",
			username:    "alice",
			fingerprint: " \t ",
			publicKey:   "ssh-ed25519 AAAA alice",
		},
		{
			name:        "blank public key",
			username:    "alice",
			fingerprint: "SHA256:test",
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
	}{
		{
			name:        "padded username",
			username:    " alice ",
			fingerprint: "SHA256:test",
		},
		{
			name:        "padded fingerprint",
			username:    "alice",
			fingerprint: " SHA256:test ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := st.EnsureUserAndKey(ctx, tt.username, tt.fingerprint, "ssh-ed25519 AAAA alice"); err == nil {
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

	if _, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice"); err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
	}

	hasKey, err := st.HasKey(ctx, " SHA256:test ")
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

	userID, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
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

func TestAttachVMRequiresExistingVM(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
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

func TestCreateVMRejectsBlankFields(t *testing.T) {
	st := newTestStore(t)
	ctx := context.Background()

	userID, err := st.EnsureUserAndKey(ctx, "alice", "SHA256:test", "ssh-ed25519 AAAA alice")
	if err != nil {
		t.Fatalf("EnsureUserAndKey: %v", err)
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
