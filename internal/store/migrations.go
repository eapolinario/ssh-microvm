package store

type migration struct {
	version int
	sql     string
}

var migrations = []migration{
	{
		version: 1,
		sql: `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version INTEGER PRIMARY KEY,
	applied_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	created_at TEXT NOT NULL,
	last_seen_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS keys (
	fingerprint TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	public_key TEXT NOT NULL,
	added_at TEXT NOT NULL,
	last_seen_at TEXT NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	key_fingerprint TEXT NOT NULL,
	remote_addr TEXT NOT NULL,
	started_at TEXT NOT NULL,
	ended_at TEXT,
	status TEXT NOT NULL,
	vm_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY(key_fingerprint) REFERENCES keys(fingerprint)
);

CREATE TABLE IF NOT EXISTS vms (
	id TEXT PRIMARY KEY,
	session_id TEXT NOT NULL,
	state_dir TEXT NOT NULL,
	fc_pid INTEGER,
	started_at TEXT NOT NULL,
	ended_at TEXT,
	exit_status INTEGER,
	FOREIGN KEY(session_id) REFERENCES sessions(id)
);

CREATE TABLE IF NOT EXISTS audit_events (
	id TEXT PRIMARY KEY,
	event_type TEXT NOT NULL,
	data_json TEXT NOT NULL,
	created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_keys_user_id ON keys(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_vms_session_id ON vms(session_id);
`,
	},
	{
		version: 2,
		sql: `
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username ON users(username);
`,
	},
	{
		version: 3,
		sql: `
CREATE UNIQUE INDEX IF NOT EXISTS idx_vms_session_id_unique ON vms(session_id);
`,
	},
	{
		version: 4,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_vm_id_insert_valid
BEFORE INSERT ON sessions
WHEN NEW.vm_id IS NOT NULL
AND NOT EXISTS (
	SELECT 1 FROM vms WHERE id = NEW.vm_id AND session_id = NEW.id
)
BEGIN
	SELECT RAISE(ABORT, 'session vm_id must reference a VM for the same session');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_vm_id_update_valid
BEFORE UPDATE OF id, vm_id ON sessions
WHEN NEW.vm_id IS NOT NULL
AND NOT EXISTS (
	SELECT 1 FROM vms WHERE id = NEW.vm_id AND session_id = NEW.id
)
BEGIN
	SELECT RAISE(ABORT, 'session vm_id must reference a VM for the same session');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_identity_update_preserves_session_links
BEFORE UPDATE OF id, session_id ON vms
WHEN EXISTS (
	SELECT 1 FROM sessions
	WHERE vm_id = OLD.id
	AND (NEW.id != OLD.id OR id != NEW.session_id)
)
BEGIN
	SELECT RAISE(ABORT, 'referenced VM identity must not break session link');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_delete_preserves_session_links
BEFORE DELETE ON vms
WHEN EXISTS (
	SELECT 1 FROM sessions WHERE vm_id = OLD.id
)
BEGIN
	SELECT RAISE(ABORT, 'referenced VM must not be deleted');
END;
`,
	},
	{
		version: 5,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_key_owner_insert_valid
BEFORE INSERT ON sessions
WHEN NOT EXISTS (
	SELECT 1 FROM keys
	WHERE fingerprint = NEW.key_fingerprint
	AND user_id = NEW.user_id
)
BEGIN
	SELECT RAISE(ABORT, 'session key must belong to session user');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_key_owner_update_valid
BEFORE UPDATE OF user_id, key_fingerprint ON sessions
WHEN NOT EXISTS (
	SELECT 1 FROM keys
	WHERE fingerprint = NEW.key_fingerprint
	AND user_id = NEW.user_id
)
BEGIN
	SELECT RAISE(ABORT, 'session key must belong to session user');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_owner_update_preserves_sessions
BEFORE UPDATE OF fingerprint, user_id ON keys
WHEN EXISTS (
	SELECT 1 FROM sessions
	WHERE key_fingerprint = OLD.fingerprint
	AND (NEW.fingerprint != OLD.fingerprint OR user_id != NEW.user_id)
)
BEGIN
	SELECT RAISE(ABORT, 'referenced key ownership must not break session link');
END;
`,
	},
	{
		version: 6,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_status_insert_valid
BEFORE INSERT ON sessions
WHEN NEW.status NOT IN ('active', 'closed', 'vm_failed')
BEGIN
	SELECT RAISE(ABORT, 'session status must be active, closed, or vm_failed');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_status_update_valid
BEFORE UPDATE OF status ON sessions
WHEN NEW.status NOT IN ('active', 'closed', 'vm_failed')
BEGIN
	SELECT RAISE(ABORT, 'session status must be active, closed, or vm_failed');
END;
`,
	},
	{
		version: 7,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_exit_status_insert_valid
BEFORE INSERT ON vms
WHEN NEW.exit_status IS NOT NULL AND NEW.exit_status < 0
BEGIN
	SELECT RAISE(ABORT, 'VM exit status must be >= 0');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_exit_status_update_valid
BEFORE UPDATE OF exit_status ON vms
WHEN NEW.exit_status IS NOT NULL AND NEW.exit_status < 0
BEGIN
	SELECT RAISE(ABORT, 'VM exit status must be >= 0');
END;
`,
	},
	{
		version: 8,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_completion_insert_consistent
BEFORE INSERT ON vms
WHEN (NEW.ended_at IS NULL AND NEW.exit_status IS NOT NULL)
OR (NEW.ended_at IS NOT NULL AND NEW.exit_status IS NULL)
BEGIN
	SELECT RAISE(ABORT, 'VM ended_at and exit_status must be set together');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_completion_update_consistent
BEFORE UPDATE OF ended_at, exit_status ON vms
WHEN (NEW.ended_at IS NULL AND NEW.exit_status IS NOT NULL)
OR (NEW.ended_at IS NOT NULL AND NEW.exit_status IS NULL)
BEGIN
	SELECT RAISE(ABORT, 'VM ended_at and exit_status must be set together');
END;
`,
	},
	{
		version: 9,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_completion_insert_consistent
BEFORE INSERT ON sessions
WHEN (NEW.status = 'active' AND NEW.ended_at IS NOT NULL)
OR (NEW.status IN ('closed', 'vm_failed') AND NEW.ended_at IS NULL)
BEGIN
	SELECT RAISE(ABORT, 'session ended_at must match terminal status');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_completion_update_consistent
BEFORE UPDATE OF status, ended_at ON sessions
WHEN (NEW.status = 'active' AND NEW.ended_at IS NOT NULL)
OR (NEW.status IN ('closed', 'vm_failed') AND NEW.ended_at IS NULL)
BEGIN
	SELECT RAISE(ABORT, 'session ended_at must match terminal status');
END;
`,
	},
	{
		version: 10,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_terminal_update_requires_no_active_vms
BEFORE UPDATE OF id, status, ended_at ON sessions
WHEN NEW.status IN ('closed', 'vm_failed')
AND EXISTS (
	SELECT 1 FROM vms
	WHERE session_id = NEW.id
	AND ended_at IS NULL
)
BEGIN
	SELECT RAISE(ABORT, 'session cannot end while VM is active');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_insert_requires_active_session
BEFORE INSERT ON vms
WHEN NOT EXISTS (
	SELECT 1 FROM sessions
	WHERE id = NEW.session_id
	AND status = 'active'
	AND ended_at IS NULL
)
BEGIN
	SELECT RAISE(ABORT, 'VM session must be active');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_session_update_requires_active_session
BEFORE UPDATE OF session_id ON vms
WHEN NOT EXISTS (
	SELECT 1 FROM sessions
	WHERE id = NEW.session_id
	AND status = 'active'
	AND ended_at IS NULL
)
BEGIN
	SELECT RAISE(ABORT, 'VM session must be active');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_active_update_requires_active_session
BEFORE UPDATE OF ended_at ON vms
WHEN NEW.ended_at IS NULL
AND NOT EXISTS (
	SELECT 1 FROM sessions
	WHERE id = NEW.session_id
	AND status = 'active'
	AND ended_at IS NULL
)
BEGIN
	SELECT RAISE(ABORT, 'active VM session must be active');
END;
`,
	},
}
