package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(path string) (*Store, error) {
	if isBlank(path) {
		return nil, errors.New("database path must be set")
	}
	if path != strings.TrimSpace(path) {
		return nil, errors.New("database path must not contain surrounding whitespace")
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode = WAL;"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enable WAL: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return errors.New("store must be set")
	}
	return s.db.Close()
}

func (s *Store) EnsureSchema(ctx context.Context) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if err := s.ensureMigrationsTable(ctx); err != nil {
		return err
	}
	for _, m := range migrations {
		applied, err := s.hasMigration(ctx, m.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := s.applyMigration(ctx, m); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ensureMigrationsTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
	version INTEGER PRIMARY KEY,
	applied_at TEXT NOT NULL
);`)
	return err
}

func (s *Store) hasMigration(ctx context.Context, version int) (bool, error) {
	row := s.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM schema_migrations WHERE version = ?", version)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) applyMigration(ctx context.Context, m migration) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, m.sql); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES(?, ?)", m.version, now()); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) HasKey(ctx context.Context, fingerprint string) (bool, error) {
	if err := s.validate(ctx); err != nil {
		return false, err
	}
	if isBlank(fingerprint) {
		return false, errors.New("key fingerprint must be set")
	}
	if hasSurroundingWhitespace(fingerprint) {
		return false, errors.New("key fingerprint must not contain surrounding whitespace")
	}
	row := s.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM keys WHERE fingerprint = ?", fingerprint)
	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) EnsureUserAndKey(ctx context.Context, username, fingerprint, publicKey string) (string, error) {
	if err := s.validate(ctx); err != nil {
		return "", err
	}
	if isBlank(username) {
		return "", errors.New("username must be set")
	}
	if hasSurroundingWhitespace(username) {
		return "", errors.New("username must not contain surrounding whitespace")
	}
	if isBlank(fingerprint) {
		return "", errors.New("key fingerprint must be set")
	}
	if hasSurroundingWhitespace(fingerprint) {
		return "", errors.New("key fingerprint must not contain surrounding whitespace")
	}
	if isBlank(publicKey) {
		return "", errors.New("public key must be set")
	}
	if hasSurroundingWhitespace(publicKey) {
		return "", errors.New("public key must not contain surrounding whitespace")
	}
	key, err := parseAuthorizedKey(publicKey)
	if err != nil {
		return "", err
	}
	if got := ssh.FingerprintSHA256(key); got != fingerprint {
		return "", errors.New("key fingerprint must match public key")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	userID, err := s.ensureUserTx(ctx, tx, username)
	if err != nil {
		_ = tx.Rollback()
		return "", err
	}
	if err := s.upsertKeyTx(ctx, tx, userID, fingerprint, publicKey); err != nil {
		_ = tx.Rollback()
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Store) ensureUserTx(ctx context.Context, tx *sql.Tx, username string) (string, error) {
	row := tx.QueryRowContext(ctx, "SELECT id FROM users WHERE username = ?", username)
	var id string
	if err := row.Scan(&id); err == nil {
		if _, err := tx.ExecContext(ctx, "UPDATE users SET last_seen_at = ? WHERE id = ?", now(), id); err != nil {
			return "", err
		}
		return id, nil
	} else if err != sql.ErrNoRows {
		return "", err
	}
	id = newID()
	if _, err := tx.ExecContext(ctx, "INSERT INTO users(id, username, created_at, last_seen_at) VALUES(?, ?, ?, ?)", id, username, now(), now()); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) upsertKeyTx(ctx context.Context, tx *sql.Tx, userID, fingerprint, publicKey string) error {
	result, err := tx.ExecContext(ctx, `INSERT INTO keys(fingerprint, user_id, public_key, added_at, last_seen_at)
VALUES(?, ?, ?, ?, ?)
ON CONFLICT(fingerprint) DO UPDATE SET
	public_key=excluded.public_key,
	last_seen_at=excluded.last_seen_at
WHERE keys.user_id = excluded.user_id`, fingerprint, userID, publicKey, now(), now())
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("key fingerprint is already enrolled for another user")
	}
	return nil
}

func (s *Store) CreateSession(ctx context.Context, session Session) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(session.ID) {
		return errors.New("session ID must be set")
	}
	if hasSurroundingWhitespace(session.ID) {
		return errors.New("session ID must not contain surrounding whitespace")
	}
	if isBlank(session.UserID) {
		return errors.New("session user ID must be set")
	}
	if hasSurroundingWhitespace(session.UserID) {
		return errors.New("session user ID must not contain surrounding whitespace")
	}
	if isBlank(session.KeyFingerprint) {
		return errors.New("session key fingerprint must be set")
	}
	if hasSurroundingWhitespace(session.KeyFingerprint) {
		return errors.New("session key fingerprint must not contain surrounding whitespace")
	}
	if isBlank(session.RemoteAddr) {
		return errors.New("session remote address must be set")
	}
	if hasSurroundingWhitespace(session.RemoteAddr) {
		return errors.New("session remote address must not contain surrounding whitespace")
	}
	if err := validateTCPAddr("session remote address", session.RemoteAddr); err != nil {
		return err
	}
	if isBlank(session.StartedAt) {
		return errors.New("session start time must be set")
	}
	if hasSurroundingWhitespace(session.StartedAt) {
		return errors.New("session start time must not contain surrounding whitespace")
	}
	if err := validateTimestamp("session start time", session.StartedAt); err != nil {
		return err
	}
	if isBlank(session.Status) {
		return errors.New("session status must be set")
	}
	if hasSurroundingWhitespace(session.Status) {
		return errors.New("session status must not contain surrounding whitespace")
	}
	if session.Status != "active" {
		return errors.New("session status must be active")
	}
	return execOne(ctx, s.db, `INSERT INTO sessions(id, user_id, key_fingerprint, remote_addr, started_at, status)
SELECT ?, ?, ?, ?, ?, ?
WHERE EXISTS (SELECT 1 FROM keys WHERE fingerprint = ? AND user_id = ?)`, session.ID, session.UserID, session.KeyFingerprint, session.RemoteAddr, session.StartedAt, session.Status, session.KeyFingerprint, session.UserID)
}

func (s *Store) EndSession(ctx context.Context, sessionID, status string) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(sessionID) {
		return errors.New("session ID must be set")
	}
	if hasSurroundingWhitespace(sessionID) {
		return errors.New("session ID must not contain surrounding whitespace")
	}
	if isBlank(status) {
		return errors.New("session status must be set")
	}
	if hasSurroundingWhitespace(status) {
		return errors.New("session status must not contain surrounding whitespace")
	}
	if status != "closed" && status != "vm_failed" {
		return errors.New("session end status must be closed or vm_failed")
	}
	return execOne(ctx, s.db, `UPDATE sessions SET ended_at = ?, status = ? WHERE id = ? AND status = 'active' AND ended_at IS NULL`, now(), status, sessionID)
}

func (s *Store) AttachVM(ctx context.Context, sessionID, vmID string) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(sessionID) {
		return errors.New("session ID must be set")
	}
	if hasSurroundingWhitespace(sessionID) {
		return errors.New("session ID must not contain surrounding whitespace")
	}
	if isBlank(vmID) {
		return errors.New("VM ID must be set")
	}
	if hasSurroundingWhitespace(vmID) {
		return errors.New("VM ID must not contain surrounding whitespace")
	}
	return execOne(ctx, s.db, `UPDATE sessions SET vm_id = ? WHERE id = ? AND status = 'active' AND ended_at IS NULL AND vm_id IS NULL AND EXISTS (SELECT 1 FROM vms WHERE id = ? AND session_id = ?)`, vmID, sessionID, vmID, sessionID)
}

func (s *Store) CreateVM(ctx context.Context, vm VM) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(vm.ID) {
		return errors.New("VM ID must be set")
	}
	if hasSurroundingWhitespace(vm.ID) {
		return errors.New("VM ID must not contain surrounding whitespace")
	}
	if isBlank(vm.SessionID) {
		return errors.New("VM session ID must be set")
	}
	if hasSurroundingWhitespace(vm.SessionID) {
		return errors.New("VM session ID must not contain surrounding whitespace")
	}
	if isBlank(vm.StateDir) {
		return errors.New("VM state directory must be set")
	}
	if hasSurroundingWhitespace(vm.StateDir) {
		return errors.New("VM state directory must not contain surrounding whitespace")
	}
	if isBlank(vm.StartedAt) {
		return errors.New("VM start time must be set")
	}
	if hasSurroundingWhitespace(vm.StartedAt) {
		return errors.New("VM start time must not contain surrounding whitespace")
	}
	if err := validateTimestamp("VM start time", vm.StartedAt); err != nil {
		return err
	}
	if vm.FCPid <= 0 {
		return errors.New("VM Firecracker PID must be > 0")
	}
	return execOne(ctx, s.db, `INSERT INTO vms(id, session_id, state_dir, fc_pid, started_at)
SELECT ?, ?, ?, ?, ?
WHERE EXISTS (SELECT 1 FROM sessions WHERE id = ? AND status = 'active' AND ended_at IS NULL AND vm_id IS NULL)
AND NOT EXISTS (SELECT 1 FROM vms WHERE session_id = ?)`, vm.ID, vm.SessionID, vm.StateDir, vm.FCPid, vm.StartedAt, vm.SessionID, vm.SessionID)
}

func (s *Store) EndVM(ctx context.Context, vmID string, exitStatus int) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(vmID) {
		return errors.New("VM ID must be set")
	}
	if hasSurroundingWhitespace(vmID) {
		return errors.New("VM ID must not contain surrounding whitespace")
	}
	if exitStatus < 0 {
		return errors.New("VM exit status must be >= 0")
	}
	return execOne(ctx, s.db, `UPDATE vms SET ended_at = ?, exit_status = ? WHERE id = ? AND ended_at IS NULL`, now(), exitStatus, vmID)
}

func (s *Store) Audit(ctx context.Context, eventType, dataJSON string) error {
	if err := s.validate(ctx); err != nil {
		return err
	}
	if isBlank(eventType) {
		return errors.New("audit event type must be set")
	}
	if hasSurroundingWhitespace(eventType) {
		return errors.New("audit event type must not contain surrounding whitespace")
	}
	if isBlank(dataJSON) {
		return errors.New("audit data must be set")
	}
	if hasSurroundingWhitespace(dataJSON) {
		return errors.New("audit data must not contain surrounding whitespace")
	}
	if !json.Valid([]byte(dataJSON)) {
		return errors.New("audit data must be valid JSON")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO audit_events(id, event_type, data_json, created_at)
VALUES(?, ?, ?, ?)`, newID(), eventType, dataJSON, now())
	return err
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func isBlank(value string) bool {
	return strings.TrimSpace(value) == ""
}

func hasSurroundingWhitespace(value string) bool {
	return value != strings.TrimSpace(value)
}

func validateTimestamp(label, value string) error {
	if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
		return fmt.Errorf("%s must be a valid RFC3339 timestamp", label)
	}
	return nil
}

func validateTCPAddr(label, value string) error {
	_, port, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("%s must be a valid TCP address", label)
	}
	if port == "" {
		return fmt.Errorf("%s port must be set", label)
	}
	if _, err := net.LookupPort("tcp", port); err != nil {
		return fmt.Errorf("%s port must be valid", label)
	}
	if _, err := net.ResolveTCPAddr("tcp", value); err != nil {
		return fmt.Errorf("%s must resolve to a valid TCP address", label)
	}
	return nil
}

func parseAuthorizedKey(value string) (ssh.PublicKey, error) {
	key, _, _, rest, err := ssh.ParseAuthorizedKey([]byte(value))
	if err != nil {
		return nil, errors.New("public key must be a valid authorized key")
	}
	if len(rest) != 0 {
		return nil, errors.New("public key must contain exactly one authorized key")
	}
	return key, nil
}

func (s *Store) validate(ctx context.Context) error {
	if s == nil || s.db == nil {
		return errors.New("store must be set")
	}
	if ctx == nil {
		return errors.New("context must be set")
	}
	return nil
}

type execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

func execOne(ctx context.Context, db execer, query string, args ...any) error {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
