package store

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	maxQueryResults    = 1000
	defaultQueryLimit  = 100
	maxFactSize        = 1024 * 1024 // 1MB
	maxMessageSize     = 64 * 1024   // 64KB
	maxTagLength       = 100
	maxTagCount        = 50
)

// SQLiteStore implements the Store interface using SQLite
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLite-backed store
func NewSQLiteStore(dataDir string) (*SQLiteStore, error) {
	dbPath := filepath.Join(dataDir, "claude_peepee.db")

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{db: db}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

func (s *SQLiteStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS facts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		tags TEXT,
		source_dir TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_facts_source_dir ON facts(source_dir);
	CREATE INDEX IF NOT EXISTS idx_facts_created_at ON facts(created_at);

	CREATE VIRTUAL TABLE IF NOT EXISTS facts_fts USING fts5(
		content,
		tags,
		content=facts,
		content_rowid=id
	);

	CREATE TRIGGER IF NOT EXISTS facts_ai AFTER INSERT ON facts BEGIN
		INSERT INTO facts_fts(rowid, content, tags) VALUES (new.id, new.content, new.tags);
	END;

	CREATE TRIGGER IF NOT EXISTS facts_ad AFTER DELETE ON facts BEGIN
		INSERT INTO facts_fts(facts_fts, rowid, content, tags) VALUES('delete', old.id, old.content, old.tags);
	END;

	CREATE TRIGGER IF NOT EXISTS facts_au AFTER UPDATE ON facts BEGIN
		INSERT INTO facts_fts(facts_fts, rowid, content, tags) VALUES('delete', old.id, old.content, old.tags);
		INSERT INTO facts_fts(rowid, content, tags) VALUES (new.id, new.content, new.tags);
	END;

	CREATE TABLE IF NOT EXISTS instances (
		id TEXT PRIMARY KEY,
		pid INTEGER NOT NULL,
		working_dir TEXT NOT NULL,
		started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_heartbeat DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_instance TEXT NOT NULL,
		to_instance TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		read_at DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_messages_to_instance ON messages(to_instance);
	`

	_, err := s.db.Exec(schema)
	return err
}

// sanitizeFTSQuery escapes special characters for FTS5 queries
func sanitizeFTSQuery(query string) string {
	// Escape double quotes by doubling them and wrap in quotes
	escaped := strings.ReplaceAll(query, `"`, `""`)
	return `"` + escaped + `"`
}

// AddFact stores a new fact
func (s *SQLiteStore) AddFact(content string, tags []string, sourceDir string) (int64, error) {
	if len(content) > maxFactSize {
		return 0, fmt.Errorf("fact content exceeds maximum size of %d bytes", maxFactSize)
	}
	if len(tags) > maxTagCount {
		return 0, fmt.Errorf("too many tags (max %d)", maxTagCount)
	}
	for _, tag := range tags {
		if len(tag) > maxTagLength {
			return 0, fmt.Errorf("tag exceeds maximum length of %d characters", maxTagLength)
		}
	}

	tagsStr := strings.Join(tags, ",")
	result, err := s.db.Exec(
		"INSERT INTO facts (content, tags, source_dir) VALUES (?, ?, ?)",
		content, tagsStr, sourceDir,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert fact: %w", err)
	}

	return result.LastInsertId()
}

// QueryFacts searches for facts matching the query
func (s *SQLiteStore) QueryFacts(query string, tags []string, sourceDir string, limit int) ([]Fact, error) {
	if limit <= 0 || limit > maxQueryResults {
		limit = defaultQueryLimit
	}

	var args []interface{}
	var conditions []string

	baseQuery := `SELECT f.id, f.content, f.tags, f.source_dir, f.created_at, f.updated_at FROM facts f`

	if query != "" {
		baseQuery += ` JOIN facts_fts ON f.id = facts_fts.rowid`
		conditions = append(conditions, `facts_fts MATCH ?`)
		args = append(args, sanitizeFTSQuery(query))
	}

	if sourceDir != "" {
		conditions = append(conditions, `f.source_dir = ?`)
		args = append(args, sourceDir)
	}

	if len(tags) > 0 {
		tagConditions := make([]string, len(tags))
		for i, tag := range tags {
			tagConditions[i] = `f.tags LIKE ?`
			args = append(args, "%"+tag+"%")
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " OR ")+")")
	}

	if len(conditions) > 0 {
		baseQuery += ` WHERE ` + strings.Join(conditions, " AND ")
	}

	baseQuery += ` ORDER BY f.created_at DESC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query facts: %w", err)
	}
	defer rows.Close()

	var facts []Fact
	for rows.Next() {
		var f Fact
		var tagsStr sql.NullString
		var sourceDir sql.NullString

		if err := rows.Scan(&f.ID, &f.Content, &tagsStr, &sourceDir, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan fact: %w", err)
		}

		if tagsStr.Valid && tagsStr.String != "" {
			f.Tags = strings.Split(tagsStr.String, ",")
		}
		if sourceDir.Valid {
			f.SourceDir = sourceDir.String
		}

		facts = append(facts, f)
	}

	return facts, nil
}

// GetFact retrieves a fact by ID
func (s *SQLiteStore) GetFact(id int64) (*Fact, error) {
	var f Fact
	var tagsStr sql.NullString
	var sourceDir sql.NullString

	err := s.db.QueryRow(
		`SELECT id, content, tags, source_dir, created_at, updated_at FROM facts WHERE id = ?`,
		id,
	).Scan(&f.ID, &f.Content, &tagsStr, &sourceDir, &f.CreatedAt, &f.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fact: %w", err)
	}

	if tagsStr.Valid && tagsStr.String != "" {
		f.Tags = strings.Split(tagsStr.String, ",")
	}
	if sourceDir.Valid {
		f.SourceDir = sourceDir.String
	}

	return &f, nil
}

// DeleteFact removes a fact by ID
func (s *SQLiteStore) DeleteFact(id int64) error {
	_, err := s.db.Exec("DELETE FROM facts WHERE id = ?", id)
	return err
}

// CountFacts returns the total and local fact counts
func (s *SQLiteStore) CountFacts(sourceDir string) (total int, local int, err error) {
	err = s.db.QueryRow("SELECT COUNT(*) FROM facts").Scan(&total)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count total facts: %w", err)
	}

	if sourceDir != "" {
		err = s.db.QueryRow("SELECT COUNT(*) FROM facts WHERE source_dir = ?", sourceDir).Scan(&local)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to count local facts: %w", err)
		}
	}

	return total, local, nil
}

// RegisterInstance adds a new instance to the store
func (s *SQLiteStore) RegisterInstance(id string, pid int, workingDir string) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO instances (id, pid, working_dir, started_at, last_heartbeat)
		 VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		id, pid, workingDir,
	)
	return err
}

// UpdateHeartbeat updates the last heartbeat time for an instance
func (s *SQLiteStore) UpdateHeartbeat(id string) error {
	_, err := s.db.Exec(
		`UPDATE instances SET last_heartbeat = CURRENT_TIMESTAMP WHERE id = ?`,
		id,
	)
	return err
}

// UnregisterInstance removes an instance from the store
func (s *SQLiteStore) UnregisterInstance(id string) error {
	_, err := s.db.Exec("DELETE FROM instances WHERE id = ?", id)
	return err
}

// GetInstances returns all registered instances
func (s *SQLiteStore) GetInstances() ([]Instance, error) {
	rows, err := s.db.Query(
		`SELECT id, pid, working_dir, started_at, last_heartbeat FROM instances ORDER BY started_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query instances: %w", err)
	}
	defer rows.Close()

	var instances []Instance
	for rows.Next() {
		var inst Instance
		if err := rows.Scan(&inst.ID, &inst.PID, &inst.WorkingDir, &inst.StartedAt, &inst.LastHeartbeat); err != nil {
			return nil, fmt.Errorf("failed to scan instance: %w", err)
		}
		instances = append(instances, inst)
	}

	return instances, nil
}

// GetInstance retrieves an instance by ID
func (s *SQLiteStore) GetInstance(id string) (*Instance, error) {
	var inst Instance
	err := s.db.QueryRow(
		`SELECT id, pid, working_dir, started_at, last_heartbeat FROM instances WHERE id = ?`,
		id,
	).Scan(&inst.ID, &inst.PID, &inst.WorkingDir, &inst.StartedAt, &inst.LastHeartbeat)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return &inst, nil
}

// CleanupStaleInstances removes instances that haven't sent a heartbeat recently
func (s *SQLiteStore) CleanupStaleInstances(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)
	_, err := s.db.Exec(
		`DELETE FROM instances WHERE last_heartbeat < ?`,
		cutoff,
	)
	return err
}

// SendMessage sends a message to another instance
func (s *SQLiteStore) SendMessage(fromInstance, toInstance, content string) (int64, error) {
	if len(content) > maxMessageSize {
		return 0, fmt.Errorf("message content exceeds maximum size of %d bytes", maxMessageSize)
	}

	result, err := s.db.Exec(
		`INSERT INTO messages (from_instance, to_instance, content) VALUES (?, ?, ?)`,
		fromInstance, toInstance, content,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	return result.LastInsertId()
}

// GetMessages retrieves messages for an instance
func (s *SQLiteStore) GetMessages(instanceID string, unreadOnly bool) ([]Message, error) {
	query := `SELECT id, from_instance, to_instance, content, created_at, read_at
	          FROM messages WHERE to_instance = ?`
	if unreadOnly {
		query += ` AND read_at IS NULL`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.db.Query(query, instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var readAt sql.NullTime

		if err := rows.Scan(&msg.ID, &msg.FromInstance, &msg.ToInstance, &msg.Content, &msg.CreatedAt, &readAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		if readAt.Valid {
			msg.ReadAt = &readAt.Time
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// MarkMessageRead marks a message as read
func (s *SQLiteStore) MarkMessageRead(id int64) error {
	_, err := s.db.Exec(
		`UPDATE messages SET read_at = CURRENT_TIMESTAMP WHERE id = ?`,
		id,
	)
	return err
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
