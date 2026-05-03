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
}
