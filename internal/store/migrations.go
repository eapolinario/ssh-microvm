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
	{
		version: 11,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_fc_pid_insert_valid
BEFORE INSERT ON vms
WHEN NEW.fc_pid IS NULL OR NEW.fc_pid <= 0
BEGIN
	SELECT RAISE(ABORT, 'VM Firecracker PID must be > 0');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_fc_pid_update_valid
BEFORE UPDATE OF fc_pid ON vms
WHEN NEW.fc_pid IS NULL OR NEW.fc_pid <= 0
BEGIN
	SELECT RAISE(ABORT, 'VM Firecracker PID must be > 0');
END;
`,
	},
	{
		version: 12,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_state_dir_insert_valid
BEFORE INSERT ON vms
WHEN trim(NEW.state_dir, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.state_dir != trim(NEW.state_dir, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM state directory must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_state_dir_update_valid
BEFORE UPDATE OF state_dir ON vms
WHEN trim(NEW.state_dir, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.state_dir != trim(NEW.state_dir, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM state directory must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 13,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_started_at_insert_valid
BEFORE INSERT ON vms
WHEN trim(NEW.started_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.started_at != trim(NEW.started_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM start time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_started_at_update_valid
BEFORE UPDATE OF started_at ON vms
WHEN trim(NEW.started_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.started_at != trim(NEW.started_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM start time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 14,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_started_at_insert_rfc3339
BEFORE INSERT ON vms
WHEN julianday(NEW.started_at) IS NULL
OR instr(NEW.started_at, 'T') != 11
OR (
	substr(NEW.started_at, -1) != 'Z'
	AND substr(NEW.started_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'VM start time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_started_at_update_rfc3339
BEFORE UPDATE OF started_at ON vms
WHEN julianday(NEW.started_at) IS NULL
OR instr(NEW.started_at, 'T') != 11
OR (
	substr(NEW.started_at, -1) != 'Z'
	AND substr(NEW.started_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'VM start time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 15,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_started_at_insert_valid
BEFORE INSERT ON sessions
WHEN trim(NEW.started_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.started_at != trim(NEW.started_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session start time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_started_at_update_valid
BEFORE UPDATE OF started_at ON sessions
WHEN trim(NEW.started_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.started_at != trim(NEW.started_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session start time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 16,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_started_at_insert_rfc3339
BEFORE INSERT ON sessions
WHEN julianday(NEW.started_at) IS NULL
OR instr(NEW.started_at, 'T') != 11
OR (
	substr(NEW.started_at, -1) != 'Z'
	AND substr(NEW.started_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'session start time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_started_at_update_rfc3339
BEFORE UPDATE OF started_at ON sessions
WHEN julianday(NEW.started_at) IS NULL
OR instr(NEW.started_at, 'T') != 11
OR (
	substr(NEW.started_at, -1) != 'Z'
	AND substr(NEW.started_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'session start time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 17,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_remote_addr_insert_valid
BEFORE INSERT ON sessions
WHEN trim(NEW.remote_addr, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.remote_addr != trim(NEW.remote_addr, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session remote address must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_remote_addr_update_valid
BEFORE UPDATE OF remote_addr ON sessions
WHEN trim(NEW.remote_addr, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.remote_addr != trim(NEW.remote_addr, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session remote address must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 18,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_remote_addr_insert_tcp
BEFORE INSERT ON sessions
WHEN (
	substr(NEW.remote_addr, 1, 1) = '['
	AND (
		instr(NEW.remote_addr, ']:') <= 2
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) = ''
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) GLOB '*[^0-9]*'
		OR CAST(substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) AS INTEGER) NOT BETWEEN 1 AND 65535
	)
)
OR (
	substr(NEW.remote_addr, 1, 1) != '['
	AND (
		instr(NEW.remote_addr, ':') <= 1
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) = ''
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) GLOB '*[^0-9]*'
		OR CAST(substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) AS INTEGER) NOT BETWEEN 1 AND 65535
	)
)
BEGIN
	SELECT RAISE(ABORT, 'session remote address must include a host and valid TCP port');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_remote_addr_update_tcp
BEFORE UPDATE OF remote_addr ON sessions
WHEN (
	substr(NEW.remote_addr, 1, 1) = '['
	AND (
		instr(NEW.remote_addr, ']:') <= 2
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) = ''
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) GLOB '*[^0-9]*'
		OR CAST(substr(NEW.remote_addr, instr(NEW.remote_addr, ']:') + 2) AS INTEGER) NOT BETWEEN 1 AND 65535
	)
)
OR (
	substr(NEW.remote_addr, 1, 1) != '['
	AND (
		instr(NEW.remote_addr, ':') <= 1
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) = ''
		OR substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) GLOB '*[^0-9]*'
		OR CAST(substr(NEW.remote_addr, instr(NEW.remote_addr, ':') + 1) AS INTEGER) NOT BETWEEN 1 AND 65535
	)
)
BEGIN
	SELECT RAISE(ABORT, 'session remote address must include a host and valid TCP port');
END;
`,
	},
	{
		version: 19,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_ended_at_insert_valid
BEFORE INSERT ON vms
WHEN NEW.ended_at IS NOT NULL
AND (
	trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.ended_at != trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'VM end time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_ended_at_update_valid
BEFORE UPDATE OF ended_at ON vms
WHEN NEW.ended_at IS NOT NULL
AND (
	trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.ended_at != trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'VM end time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 20,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_ended_at_insert_rfc3339
BEFORE INSERT ON vms
WHEN NEW.ended_at IS NOT NULL
AND (
	julianday(NEW.ended_at) IS NULL
	OR instr(NEW.ended_at, 'T') != 11
	OR (
		substr(NEW.ended_at, -1) != 'Z'
		AND substr(NEW.ended_at, -6, 1) NOT IN ('+', '-')
	)
)
BEGIN
	SELECT RAISE(ABORT, 'VM end time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_ended_at_update_rfc3339
BEFORE UPDATE OF ended_at ON vms
WHEN NEW.ended_at IS NOT NULL
AND (
	julianday(NEW.ended_at) IS NULL
	OR instr(NEW.ended_at, 'T') != 11
	OR (
		substr(NEW.ended_at, -1) != 'Z'
		AND substr(NEW.ended_at, -6, 1) NOT IN ('+', '-')
	)
)
BEGIN
	SELECT RAISE(ABORT, 'VM end time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 21,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_ended_at_insert_valid
BEFORE INSERT ON sessions
WHEN NEW.ended_at IS NOT NULL
AND (
	trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.ended_at != trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'session end time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_ended_at_update_valid
BEFORE UPDATE OF ended_at ON sessions
WHEN NEW.ended_at IS NOT NULL
AND (
	trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.ended_at != trim(NEW.ended_at, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'session end time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 22,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_ended_at_insert_rfc3339
BEFORE INSERT ON sessions
WHEN NEW.ended_at IS NOT NULL
AND (
	julianday(NEW.ended_at) IS NULL
	OR instr(NEW.ended_at, 'T') != 11
	OR (
		substr(NEW.ended_at, -1) != 'Z'
		AND substr(NEW.ended_at, -6, 1) NOT IN ('+', '-')
	)
)
BEGIN
	SELECT RAISE(ABORT, 'session end time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_ended_at_update_rfc3339
BEFORE UPDATE OF ended_at ON sessions
WHEN NEW.ended_at IS NOT NULL
AND (
	julianday(NEW.ended_at) IS NULL
	OR instr(NEW.ended_at, 'T') != 11
	OR (
		substr(NEW.ended_at, -1) != 'Z'
		AND substr(NEW.ended_at, -6, 1) NOT IN ('+', '-')
	)
)
BEGIN
	SELECT RAISE(ABORT, 'session end time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 23,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_event_type_insert_valid
BEFORE INSERT ON audit_events
WHEN trim(NEW.event_type, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.event_type != trim(NEW.event_type, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event type must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_event_type_update_valid
BEFORE UPDATE OF event_type ON audit_events
WHEN trim(NEW.event_type, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.event_type != trim(NEW.event_type, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event type must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 24,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_data_json_insert_valid
BEFORE INSERT ON audit_events
WHEN trim(NEW.data_json, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.data_json != trim(NEW.data_json, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit data must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_data_json_update_valid
BEFORE UPDATE OF data_json ON audit_events
WHEN trim(NEW.data_json, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.data_json != trim(NEW.data_json, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit data must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 25,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_data_json_insert_json
BEFORE INSERT ON audit_events
WHEN json_valid(NEW.data_json) != 1
BEGIN
	SELECT RAISE(ABORT, 'audit data must be valid JSON');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_data_json_update_json
BEFORE UPDATE OF data_json ON audit_events
WHEN json_valid(NEW.data_json) != 1
BEGIN
	SELECT RAISE(ABORT, 'audit data must be valid JSON');
END;
`,
	},
	{
		version: 26,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_created_at_insert_valid
BEFORE INSERT ON audit_events
WHEN trim(NEW.created_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.created_at != trim(NEW.created_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event creation time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_created_at_update_valid
BEFORE UPDATE OF created_at ON audit_events
WHEN trim(NEW.created_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.created_at != trim(NEW.created_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event creation time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 27,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_created_at_insert_rfc3339
BEFORE INSERT ON audit_events
WHEN julianday(NEW.created_at) IS NULL
OR instr(NEW.created_at, 'T') != 11
OR (
	substr(NEW.created_at, -1) != 'Z'
	AND substr(NEW.created_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'audit event creation time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_created_at_update_rfc3339
BEFORE UPDATE OF created_at ON audit_events
WHEN julianday(NEW.created_at) IS NULL
OR instr(NEW.created_at, 'T') != 11
OR (
	substr(NEW.created_at, -1) != 'Z'
	AND substr(NEW.created_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'audit event creation time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 28,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_username_insert_valid
BEFORE INSERT ON users
WHEN trim(NEW.username, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.username != trim(NEW.username, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'username must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_username_update_valid
BEFORE UPDATE OF username ON users
WHEN trim(NEW.username, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.username != trim(NEW.username, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'username must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 29,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_created_at_insert_valid
BEFORE INSERT ON users
WHEN trim(NEW.created_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.created_at != trim(NEW.created_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user creation time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_created_at_update_valid
BEFORE UPDATE OF created_at ON users
WHEN trim(NEW.created_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.created_at != trim(NEW.created_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user creation time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 30,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_created_at_insert_rfc3339
BEFORE INSERT ON users
WHEN julianday(NEW.created_at) IS NULL
OR instr(NEW.created_at, 'T') != 11
OR (
	substr(NEW.created_at, -1) != 'Z'
	AND substr(NEW.created_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'user creation time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_created_at_update_rfc3339
BEFORE UPDATE OF created_at ON users
WHEN julianday(NEW.created_at) IS NULL
OR instr(NEW.created_at, 'T') != 11
OR (
	substr(NEW.created_at, -1) != 'Z'
	AND substr(NEW.created_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'user creation time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 31,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_last_seen_at_insert_valid
BEFORE INSERT ON users
WHEN trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.last_seen_at != trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user last seen time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_last_seen_at_update_valid
BEFORE UPDATE OF last_seen_at ON users
WHEN trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.last_seen_at != trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user last seen time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 32,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_last_seen_at_insert_rfc3339
BEFORE INSERT ON users
WHEN julianday(NEW.last_seen_at) IS NULL
OR instr(NEW.last_seen_at, 'T') != 11
OR (
	substr(NEW.last_seen_at, -1) != 'Z'
	AND substr(NEW.last_seen_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'user last seen time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_last_seen_at_update_rfc3339
BEFORE UPDATE OF last_seen_at ON users
WHEN julianday(NEW.last_seen_at) IS NULL
OR instr(NEW.last_seen_at, 'T') != 11
OR (
	substr(NEW.last_seen_at, -1) != 'Z'
	AND substr(NEW.last_seen_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'user last seen time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 33,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_fingerprint_insert_valid
BEFORE INSERT ON keys
WHEN trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.fingerprint != trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key fingerprint must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_fingerprint_update_valid
BEFORE UPDATE OF fingerprint ON keys
WHEN trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.fingerprint != trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key fingerprint must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 34,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_valid
BEFORE INSERT ON keys
WHEN trim(NEW.public_key, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.public_key != trim(NEW.public_key, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'public key must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_valid
BEFORE UPDATE OF public_key ON keys
WHEN trim(NEW.public_key, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.public_key != trim(NEW.public_key, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'public key must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 35,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_added_at_insert_valid
BEFORE INSERT ON keys
WHEN trim(NEW.added_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.added_at != trim(NEW.added_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key addition time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_added_at_update_valid
BEFORE UPDATE OF added_at ON keys
WHEN trim(NEW.added_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.added_at != trim(NEW.added_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key addition time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 36,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_added_at_insert_rfc3339
BEFORE INSERT ON keys
WHEN julianday(NEW.added_at) IS NULL
OR instr(NEW.added_at, 'T') != 11
OR (
	substr(NEW.added_at, -1) != 'Z'
	AND substr(NEW.added_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'key addition time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_added_at_update_rfc3339
BEFORE UPDATE OF added_at ON keys
WHEN julianday(NEW.added_at) IS NULL
OR instr(NEW.added_at, 'T') != 11
OR (
	substr(NEW.added_at, -1) != 'Z'
	AND substr(NEW.added_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'key addition time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 37,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_last_seen_at_insert_valid
BEFORE INSERT ON keys
WHEN trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.last_seen_at != trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key last seen time must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_last_seen_at_update_valid
BEFORE UPDATE OF last_seen_at ON keys
WHEN trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.last_seen_at != trim(NEW.last_seen_at, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key last seen time must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 38,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_last_seen_at_insert_rfc3339
BEFORE INSERT ON keys
WHEN julianday(NEW.last_seen_at) IS NULL
OR instr(NEW.last_seen_at, 'T') != 11
OR (
	substr(NEW.last_seen_at, -1) != 'Z'
	AND substr(NEW.last_seen_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'key last seen time must be a valid RFC3339 timestamp');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_last_seen_at_update_rfc3339
BEFORE UPDATE OF last_seen_at ON keys
WHEN julianday(NEW.last_seen_at) IS NULL
OR instr(NEW.last_seen_at, 'T') != 11
OR (
	substr(NEW.last_seen_at, -1) != 'Z'
	AND substr(NEW.last_seen_at, -6, 1) NOT IN ('+', '-')
)
BEGIN
	SELECT RAISE(ABORT, 'key last seen time must be a valid RFC3339 timestamp');
END;
`,
	},
	{
		version: 39,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_users_id_insert_valid
BEFORE INSERT ON users
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_users_id_update_valid
BEFORE UPDATE OF id ON users
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'user ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 40,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_id_insert_valid
BEFORE INSERT ON sessions
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_id_update_valid
BEFORE UPDATE OF id ON sessions
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 41,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_id_insert_valid
BEFORE INSERT ON vms
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_id_update_valid
BEFORE UPDATE OF id ON vms
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 42,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_audit_events_id_insert_valid
BEFORE INSERT ON audit_events
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_audit_events_id_update_valid
BEFORE UPDATE OF id ON audit_events
WHEN trim(NEW.id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.id != trim(NEW.id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'audit event ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 43,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_user_id_insert_valid
BEFORE INSERT ON sessions
WHEN trim(NEW.user_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.user_id != trim(NEW.user_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session user ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_user_id_update_valid
BEFORE UPDATE OF user_id ON sessions
WHEN trim(NEW.user_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.user_id != trim(NEW.user_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session user ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 44,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_key_fingerprint_insert_valid
BEFORE INSERT ON sessions
WHEN trim(NEW.key_fingerprint, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.key_fingerprint != trim(NEW.key_fingerprint, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session key fingerprint must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_key_fingerprint_update_valid
BEFORE UPDATE OF key_fingerprint ON sessions
WHEN trim(NEW.key_fingerprint, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.key_fingerprint != trim(NEW.key_fingerprint, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'session key fingerprint must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 45,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_vms_session_id_insert_valid
BEFORE INSERT ON vms
WHEN trim(NEW.session_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.session_id != trim(NEW.session_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM session ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_vms_session_id_update_valid
BEFORE UPDATE OF session_id ON vms
WHEN trim(NEW.session_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.session_id != trim(NEW.session_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'VM session ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 46,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_sessions_vm_id_insert_valid_value
BEFORE INSERT ON sessions
WHEN NEW.vm_id IS NOT NULL
AND (
	trim(NEW.vm_id, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.vm_id != trim(NEW.vm_id, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'session VM ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_sessions_vm_id_update_valid_value
BEFORE UPDATE OF vm_id ON sessions
WHEN NEW.vm_id IS NOT NULL
AND (
	trim(NEW.vm_id, char(9, 10, 11, 12, 13, 32)) = ''
	OR NEW.vm_id != trim(NEW.vm_id, char(9, 10, 11, 12, 13, 32))
)
BEGIN
	SELECT RAISE(ABORT, 'session VM ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 47,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_user_id_insert_valid
BEFORE INSERT ON keys
WHEN trim(NEW.user_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.user_id != trim(NEW.user_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key user ID must be set and not contain surrounding whitespace');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_user_id_update_valid
BEFORE UPDATE OF user_id ON keys
WHEN trim(NEW.user_id, char(9, 10, 11, 12, 13, 32)) = ''
OR NEW.user_id != trim(NEW.user_id, char(9, 10, 11, 12, 13, 32))
BEGIN
	SELECT RAISE(ABORT, 'key user ID must be set and not contain surrounding whitespace');
END;
`,
	},
	{
		version: 48,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_single_line
BEFORE INSERT ON keys
WHEN instr(NEW.public_key, char(10)) > 0
OR instr(NEW.public_key, char(13)) > 0
BEGIN
	SELECT RAISE(ABORT, 'public key must contain exactly one authorized key');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_single_line
BEFORE UPDATE OF public_key ON keys
WHEN instr(NEW.public_key, char(10)) > 0
OR instr(NEW.public_key, char(13)) > 0
BEGIN
	SELECT RAISE(ABORT, 'public key must contain exactly one authorized key');
END;
`,
	},
	{
		version: 49,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_fields
BEFORE INSERT ON keys
WHEN instr(NEW.public_key, ' ') = 0
AND instr(NEW.public_key, char(9)) = 0
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_fields
BEFORE UPDATE OF public_key ON keys
WHEN instr(NEW.public_key, ' ') = 0
AND instr(NEW.public_key, char(9)) = 0
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;
`,
	},
	{
		version: 50,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_type
BEFORE INSERT ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND (
	CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
	ELSE NEW.public_key
	END
) NOT IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
BEGIN
	SELECT RAISE(ABORT, 'public key must be a supported authorized key type');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_type
BEFORE UPDATE OF public_key ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND (
	CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
	ELSE NEW.public_key
	END
) NOT IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
BEGIN
	SELECT RAISE(ABORT, 'public key must be a supported authorized key type');
END;
`,
	},
	{
		version: 51,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_blob
BEFORE INSERT ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
WHEN instr(NEW.public_key, char(9)) > 0
THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
ELSE NEW.public_key
END) IN (
'ssh-rsa',
'ssh-rsa-cert-v01@openssh.com',
'ssh-dss',
'ssh-dss-cert-v01@openssh.com',
'ssh-ed25519',
'ssh-ed25519-cert-v01@openssh.com',
'ecdsa-sha2-nistp256',
'ecdsa-sha2-nistp256-cert-v01@openssh.com',
'ecdsa-sha2-nistp384',
'ecdsa-sha2-nistp384-cert-v01@openssh.com',
'ecdsa-sha2-nistp521',
'ecdsa-sha2-nistp521-cert-v01@openssh.com',
'sk-ssh-ed25519@openssh.com',
'sk-ssh-ed25519-cert-v01@openssh.com',
'sk-ecdsa-sha2-nistp256@openssh.com',
'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (
ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)) = ''
OR length(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32))) < 16
OR ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)) GLOB '[^A-Za-z0-9+/=]*'
)
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_blob
BEFORE UPDATE OF public_key ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
WHEN instr(NEW.public_key, char(9)) > 0
THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
ELSE NEW.public_key
END) IN (
'ssh-rsa',
'ssh-rsa-cert-v01@openssh.com',
'ssh-dss',
'ssh-dss-cert-v01@openssh.com',
'ssh-ed25519',
'ssh-ed25519-cert-v01@openssh.com',
'ecdsa-sha2-nistp256',
'ecdsa-sha2-nistp256-cert-v01@openssh.com',
'ecdsa-sha2-nistp384',
'ecdsa-sha2-nistp384-cert-v01@openssh.com',
'ecdsa-sha2-nistp521',
'ecdsa-sha2-nistp521-cert-v01@openssh.com',
'sk-ssh-ed25519@openssh.com',
'sk-ssh-ed25519-cert-v01@openssh.com',
'sk-ecdsa-sha2-nistp256@openssh.com',
'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (
ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)) = ''
OR length(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32))) < 16
OR ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)) GLOB '[^A-Za-z0-9+/=]*'
)
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;
`,
	},
	{
		version: 52,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_blob_chars
BEFORE INSERT ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
WHEN instr(NEW.public_key, char(9)) > 0
THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
ELSE NEW.public_key
END) IN (
'ssh-rsa',
'ssh-rsa-cert-v01@openssh.com',
'ssh-dss',
'ssh-dss-cert-v01@openssh.com',
'ssh-ed25519',
'ssh-ed25519-cert-v01@openssh.com',
'ecdsa-sha2-nistp256',
'ecdsa-sha2-nistp256-cert-v01@openssh.com',
'ecdsa-sha2-nistp384',
'ecdsa-sha2-nistp384-cert-v01@openssh.com',
'ecdsa-sha2-nistp521',
'ecdsa-sha2-nistp521-cert-v01@openssh.com',
'sk-ssh-ed25519@openssh.com',
'sk-ssh-ed25519-cert-v01@openssh.com',
'sk-ecdsa-sha2-nistp256@openssh.com',
'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (CASE
WHEN instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') > 0
AND (instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) = 0
OR instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') < instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)))
THEN substr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), 1, instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') - 1)
WHEN instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) > 0
THEN substr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), 1, instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) - 1)
ELSE ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32))
END) GLOB '*[^A-Za-z0-9+/=]*'
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_blob_chars
BEFORE UPDATE OF public_key ON keys
WHEN (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN substr(NEW.public_key, 1, instr(NEW.public_key, ' ') - 1)
WHEN instr(NEW.public_key, char(9)) > 0
THEN substr(NEW.public_key, 1, instr(NEW.public_key, char(9)) - 1)
ELSE NEW.public_key
END) IN (
'ssh-rsa',
'ssh-rsa-cert-v01@openssh.com',
'ssh-dss',
'ssh-dss-cert-v01@openssh.com',
'ssh-ed25519',
'ssh-ed25519-cert-v01@openssh.com',
'ecdsa-sha2-nistp256',
'ecdsa-sha2-nistp256-cert-v01@openssh.com',
'ecdsa-sha2-nistp384',
'ecdsa-sha2-nistp384-cert-v01@openssh.com',
'ecdsa-sha2-nistp521',
'ecdsa-sha2-nistp521-cert-v01@openssh.com',
'sk-ssh-ed25519@openssh.com',
'sk-ssh-ed25519-cert-v01@openssh.com',
'sk-ecdsa-sha2-nistp256@openssh.com',
'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (CASE
WHEN instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') > 0
AND (instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) = 0
OR instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') < instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)))
THEN substr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), 1, instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), ' ') - 1)
WHEN instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) > 0
THEN substr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), 1, instr(ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32)), char(9)) - 1)
ELSE ltrim(substr(NEW.public_key, (CASE
WHEN instr(NEW.public_key, ' ') > 0
AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
THEN instr(NEW.public_key, ' ')
WHEN instr(NEW.public_key, char(9)) > 0
THEN instr(NEW.public_key, char(9))
ELSE 0
END) + 1), char(9, 32))
END) GLOB '*[^A-Za-z0-9+/=]*'
BEGIN
	SELECT RAISE(ABORT, 'public key must be a valid authorized key');
END;
`,
	},
	{
		version: 53,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_blob_base64_shape
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'public key must be a valid authorized key')
FROM blob
WHERE (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND key_type IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (
	length(value) % 4 != 0
	OR value GLOB '*=[A-Za-z0-9+/]*'
	OR value GLOB '*===*'
);
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_blob_base64_shape
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'public key must be a valid authorized key')
FROM blob
WHERE (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND key_type IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND (
	length(value) % 4 != 0
	OR value GLOB '*=[A-Za-z0-9+/]*'
	OR value GLOB '*===*'
);
END;
`,
	},
	{
		version: 54,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_authorized_key_blob_type
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'public key blob must match authorized key type')
FROM blob
WHERE (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND key_type IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND length(value) >= 16
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND value NOT GLOB CASE key_type
	WHEN 'ssh-rsa' THEN 'AAAAB3NzaC1yc2*'
	WHEN 'ssh-rsa-cert-v01@openssh.com' THEN 'AAAAHHNzaC1yc2EtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ssh-dss' THEN 'AAAAB3NzaC1kc3*'
	WHEN 'ssh-dss-cert-v01@openssh.com' THEN 'AAAAHHNzaC1kc3MtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ssh-ed25519' THEN 'AAAAC3NzaC1lZDI1NTE5*'
	WHEN 'ssh-ed25519-cert-v01@openssh.com' THEN 'AAAAIHNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29t*'
	WHEN 'ecdsa-sha2-nistp256' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHAyNT*'
	WHEN 'ecdsa-sha2-nistp256-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHAyNTYtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ecdsa-sha2-nistp384' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHAzOD*'
	WHEN 'ecdsa-sha2-nistp384-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHAzODQtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ecdsa-sha2-nistp521' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHA1Mj*'
	WHEN 'ecdsa-sha2-nistp521-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHA1MjEtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'sk-ssh-ed25519@openssh.com' THEN 'AAAAGnNrLXNzaC1lZDI1NTE5QG9wZW5zc2guY29t*'
	WHEN 'sk-ssh-ed25519-cert-v01@openssh.com' THEN 'AAAAI3NrLXNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29t*'
	WHEN 'sk-ecdsa-sha2-nistp256@openssh.com' THEN 'AAAAInNrLWVjZHNhLXNoYTItbmlzdHAyNTZAb3BlbnNzaC5jb2*'
	WHEN 'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com' THEN 'AAAAK3NrLWVjZHNhLXNoYTItbmlzdHAyNTYtY2VydC12MDFAb3BlbnNzaC5jb2*'
END;
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_authorized_key_blob_type
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'public key blob must match authorized key type')
FROM blob
WHERE (instr(NEW.public_key, ' ') > 0 OR instr(NEW.public_key, char(9)) > 0)
AND instr(NEW.public_key, char(10)) = 0
AND instr(NEW.public_key, char(13)) = 0
AND key_type IN (
	'ssh-rsa',
	'ssh-rsa-cert-v01@openssh.com',
	'ssh-dss',
	'ssh-dss-cert-v01@openssh.com',
	'ssh-ed25519',
	'ssh-ed25519-cert-v01@openssh.com',
	'ecdsa-sha2-nistp256',
	'ecdsa-sha2-nistp256-cert-v01@openssh.com',
	'ecdsa-sha2-nistp384',
	'ecdsa-sha2-nistp384-cert-v01@openssh.com',
	'ecdsa-sha2-nistp521',
	'ecdsa-sha2-nistp521-cert-v01@openssh.com',
	'sk-ssh-ed25519@openssh.com',
	'sk-ssh-ed25519-cert-v01@openssh.com',
	'sk-ecdsa-sha2-nistp256@openssh.com',
	'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com'
)
AND length(value) >= 16
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND value NOT GLOB CASE key_type
	WHEN 'ssh-rsa' THEN 'AAAAB3NzaC1yc2*'
	WHEN 'ssh-rsa-cert-v01@openssh.com' THEN 'AAAAHHNzaC1yc2EtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ssh-dss' THEN 'AAAAB3NzaC1kc3*'
	WHEN 'ssh-dss-cert-v01@openssh.com' THEN 'AAAAHHNzaC1kc3MtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ssh-ed25519' THEN 'AAAAC3NzaC1lZDI1NTE5*'
	WHEN 'ssh-ed25519-cert-v01@openssh.com' THEN 'AAAAIHNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29t*'
	WHEN 'ecdsa-sha2-nistp256' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHAyNT*'
	WHEN 'ecdsa-sha2-nistp256-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHAyNTYtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ecdsa-sha2-nistp384' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHAzOD*'
	WHEN 'ecdsa-sha2-nistp384-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHAzODQtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'ecdsa-sha2-nistp521' THEN 'AAAAE2VjZHNhLXNoYTItbmlzdHA1Mj*'
	WHEN 'ecdsa-sha2-nistp521-cert-v01@openssh.com' THEN 'AAAAKGVjZHNhLXNoYTItbmlzdHA1MjEtY2VydC12MDFAb3BlbnNzaC5jb2*'
	WHEN 'sk-ssh-ed25519@openssh.com' THEN 'AAAAGnNrLXNzaC1lZDI1NTE5QG9wZW5zc2guY29t*'
	WHEN 'sk-ssh-ed25519-cert-v01@openssh.com' THEN 'AAAAI3NrLXNzaC1lZDI1NTE5LWNlcnQtdjAxQG9wZW5zc2guY29t*'
	WHEN 'sk-ecdsa-sha2-nistp256@openssh.com' THEN 'AAAAInNrLWVjZHNhLXNoYTItbmlzdHAyNTZAb3BlbnNzaC5jb2*'
	WHEN 'sk-ecdsa-sha2-nistp256-cert-v01@openssh.com' THEN 'AAAAK3NrLWVjZHNhLXNoYTItbmlzdHAyNTYtY2VydC12MDFAb3BlbnNzaC5jb2*'
END;
END;
`,
	},
	{
		version: 55,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_fingerprint_insert_sha256
BEFORE INSERT ON keys
WHEN trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32)) != ''
AND NEW.fingerprint = trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32))
AND (
	length(NEW.fingerprint) != 50
	OR substr(NEW.fingerprint, 1, 7) != 'SHA256:'
	OR substr(NEW.fingerprint, 8) GLOB '*[^A-Za-z0-9+/]*'
)
BEGIN
	SELECT RAISE(ABORT, 'key fingerprint must be a SHA256 fingerprint');
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_fingerprint_update_sha256
BEFORE UPDATE OF fingerprint ON keys
WHEN trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32)) != ''
AND NEW.fingerprint = trim(NEW.fingerprint, char(9, 10, 11, 12, 13, 32))
AND (
	length(NEW.fingerprint) != 50
	OR substr(NEW.fingerprint, 1, 7) != 'SHA256:'
	OR substr(NEW.fingerprint, 8) GLOB '*[^A-Za-z0-9+/]*'
)
BEGIN
	SELECT RAISE(ABORT, 'key fingerprint must be a SHA256 fingerprint');
END;
`,
	},
	{
		version: 56,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_ed25519_blob_length
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-ed25519 public key blob must be complete')
FROM blob
WHERE key_type = 'ssh-ed25519'
AND value GLOB 'AAAAC3NzaC1lZDI1NTE5*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) != 68;
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_ed25519_blob_length
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-ed25519 public key blob must be complete')
FROM blob
WHERE key_type = 'ssh-ed25519'
AND value GLOB 'AAAAC3NzaC1lZDI1NTE5*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) != 68;
END;
`,
	},
	{
		version: 57,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_fields
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include exponent and modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) <= 16;
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_fields
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include exponent and modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) <= 16;
END;
`,
	},
	{
		version: 58,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_modulus
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) > 16
AND length(value) <= 24;
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_modulus
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) > 16
AND length(value) <= 24;
END;
`,
	},
	{
		version: 59,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_modulus_length
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include complete modulus length')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) > 24
AND length(value) < 32;
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_modulus_length
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include complete modulus length')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) > 24
AND length(value) < 32;
END;
`,
	},
	{
		version: 60,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_modulus_bytes
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include complete modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQAB*=='
AND substr(value, 25) != 'AAAAAA==';
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_modulus_bytes
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include complete modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQAB*=='
AND substr(value, 25) != 'AAAAAA==';
END;
`,
	},
	{
		version: 61,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_declared_modulus_bytes
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include declared modulus bytes')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQABAAAAAg?=';
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_declared_modulus_bytes
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include declared modulus bytes')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQABAAAAAg?=';
END;
`,
	},
	{
		version: 62,
		sql: `
CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_insert_rsa_blob_declared_three_modulus_bytes
BEFORE INSERT ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include declared three-byte modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQABAAAAAw??';
END;

CREATE TRIGGER IF NOT EXISTS trg_keys_public_key_update_rsa_blob_declared_three_modulus_bytes
BEFORE UPDATE OF public_key ON keys
BEGIN
WITH parsed AS (
	SELECT CASE
	WHEN instr(NEW.public_key, ' ') > 0
	AND (instr(NEW.public_key, char(9)) = 0 OR instr(NEW.public_key, ' ') < instr(NEW.public_key, char(9)))
	THEN instr(NEW.public_key, ' ')
	WHEN instr(NEW.public_key, char(9)) > 0
	THEN instr(NEW.public_key, char(9))
	ELSE 0
	END AS first_separator
),
fields AS (
	SELECT CASE
	WHEN first_separator > 0 THEN substr(NEW.public_key, 1, first_separator - 1)
	ELSE NEW.public_key
	END AS key_type,
	ltrim(substr(NEW.public_key, first_separator + 1), char(9, 32)) AS after_type
	FROM parsed
),
blob AS (
	SELECT key_type, CASE
	WHEN instr(after_type, ' ') > 0
	AND (instr(after_type, char(9)) = 0 OR instr(after_type, ' ') < instr(after_type, char(9)))
	THEN substr(after_type, 1, instr(after_type, ' ') - 1)
	WHEN instr(after_type, char(9)) > 0
	THEN substr(after_type, 1, instr(after_type, char(9)) - 1)
	ELSE after_type
	END AS value
	FROM fields
)
SELECT RAISE(ABORT, 'ssh-rsa public key blob must include declared three-byte modulus')
FROM blob
WHERE key_type = 'ssh-rsa'
AND value GLOB 'AAAAB3NzaC1yc2*'
AND value NOT GLOB '*[^A-Za-z0-9+/=]*'
AND length(value) % 4 = 0
AND value NOT GLOB '*=[A-Za-z0-9+/]*'
AND value NOT GLOB '*===*'
AND length(value) = 32
AND value GLOB 'AAAAB3NzaC1yc2EAAAADAQABAAAAAw??';
END;
`,
	},
}
