package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"

	"github.com/TheRealSibasishBehera/syncsh/internal/parser"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed schema.sql
	Schema string

	ErrEntryNotFound = errors.New("history entry not found")
	ErrDuplicateHash = errors.New("duplicate hash - entry already exists")
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// InitSchema creates the database tables if they don't exist
func (s *Store) InitSchema(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, Schema)
	return err
}

// CreateEntry inserts a new history entry into the database
func (s *Store) CreateEntry(ctx context.Context, entry *parser.HistoryEntry) error {
	// Generate hash if not provided
	if entry.Hash == "" {
		entry.Hash = generateHash(entry.Timestamp, entry.MachineID, entry.Command)
	}

	query := `INSERT INTO history_entries (timestamp, machine_id, command, duration, exit_code, hash) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query,
		entry.Timestamp, entry.MachineID, entry.Command,
		entry.Duration, entry.ExitCode, entry.Hash)

	if err != nil {
		// Check for unique constraint violation on hash
		if isUniqueConstraintError(err) {
			return ErrDuplicateHash
		}
		return fmt.Errorf("insert history entry: %w", err)
	}

	return nil
}

// GetEntryByHash retrieves a history entry by its hash
func (s *Store) GetEntryByHash(ctx context.Context, hash string) (parser.HistoryEntry, error) {
	query := `SELECT id, timestamp, machine_id, command, duration, exit_code, hash 
	          FROM history_entries WHERE hash = ?`

	row := s.db.QueryRowContext(ctx, query, hash)

	var entry parser.HistoryEntry
	err := row.Scan(&entry.ID, &entry.Timestamp, &entry.MachineID,
		&entry.Command, &entry.Duration, &entry.ExitCode, &entry.Hash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return parser.HistoryEntry{}, ErrEntryNotFound
		}
		return parser.HistoryEntry{}, fmt.Errorf("scan history entry: %w", err)
	}

	return entry, nil
}

// ListEntries retrieves history entries with optional filtering
func (s *Store) ListEntries(ctx context.Context, machineID string, since int64, limit int) ([]parser.HistoryEntry, error) {
	query := `SELECT id, timestamp, machine_id, command, duration, exit_code, hash 
	          FROM history_entries WHERE 1=1`
	args := []interface{}{}

	if machineID != "" {
		query += " AND machine_id = ?"
		args = append(args, machineID)
	}

	if since > 0 {
		query += " AND timestamp > ?"
		args = append(args, since)
	}

	query += " ORDER BY timestamp DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query history entries: %w", err)
	}
	defer rows.Close()

	var entries []parser.HistoryEntry
	for rows.Next() {
		var entry parser.HistoryEntry
		err := rows.Scan(&entry.ID, &entry.Timestamp, &entry.MachineID,
			&entry.Command, &entry.Duration, &entry.ExitCode, &entry.Hash)
		if err != nil {
			return nil, fmt.Errorf("scan history entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetLastSyncTimestamp retrieves the last sync timestamp for a machine
func (s *Store) GetLastSyncTimestamp(ctx context.Context, machineID string) (int64, error) {
	query := `SELECT last_sync_timestamp FROM sync_state WHERE machine_id = ?`

	var timestamp int64
	err := s.db.QueryRowContext(ctx, query, machineID).Scan(&timestamp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // No previous sync
		}
		return 0, fmt.Errorf("get last sync timestamp: %w", err)
	}

	return timestamp, nil
}

// UpdateLastSyncTimestamp updates the last sync timestamp for a machine
func (s *Store) UpdateLastSyncTimestamp(ctx context.Context, machineID string, timestamp int64) error {
	query := `INSERT OR REPLACE INTO sync_state (machine_id, last_sync_timestamp) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, query, machineID, timestamp)
	if err != nil {
		return fmt.Errorf("update last sync timestamp: %w", err)
	}

	return nil
}

// DeleteEntry removes a history entry by hash
func (s *Store) DeleteEntry(ctx context.Context, hash string) error {
	query := `DELETE FROM history_entries WHERE hash = ?`

	result, err := s.db.ExecContext(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("delete history entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrEntryNotFound
	}

	return nil
}

// generateHash creates a unique hash for a history entry
func generateHash(timestamp int64, machineID, command string) string {
	data := strconv.FormatInt(timestamp, 10) + machineID + command
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// isUniqueConstraintError checks if error is due to unique constraint violation
func isUniqueConstraintError(err error) bool {
	return err != nil && (
	// SQLite unique constraint error messages
	contains(err.Error(), "UNIQUE constraint failed") ||
		contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || hasSubstring(s, substr))
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
