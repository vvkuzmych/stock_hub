package service

import (
	"database/sql"
	"fmt"
	"time"

	"stock_hub/internal/model"
	"stock_hub/pkg/migrate"

	_ "github.com/lib/pq"  // PostgreSQL driver
	_ "modernc.org/sqlite" // SQLite driver
)

// MessageServicePostgres handles message persistence with PostgreSQL support
type MessageServicePostgres struct {
	db     *sql.DB
	driver string
}

// NewMessageServicePostgres creates a new message service with database type support
func NewMessageServicePostgres(driver, dsn, migrationsDir string) (*MessageServicePostgres, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	service := &MessageServicePostgres{
		db:     db,
		driver: driver,
	}

	// Run migrations
	migrator := migrate.NewMigrator(db, migrationsDir)
	if err := migrator.Up(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return service, nil
}

// SaveMessage saves a message to the database
func (s *MessageServicePostgres) SaveMessage(content string, clientID string) error {
	query := `INSERT INTO messages (content, client_id, created_at) VALUES ($1, $2, $3)`
	if s.driver == "sqlite" {
		query = `INSERT INTO messages (content, client_id, created_at) VALUES (?, ?, ?)`
	}
	_, err := s.db.Exec(query, content, clientID, time.Now())
	return err
}

// GetMessages retrieves messages from the database
func (s *MessageServicePostgres) GetMessages(limit int) ([]*model.Message, error) {
	query := `SELECT id, content, client_id, created_at FROM messages ORDER BY created_at DESC LIMIT $1`
	if s.driver == "sqlite" {
		query = `SELECT id, content, client_id, created_at FROM messages ORDER BY created_at DESC LIMIT ?`
	}

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.ID, &msg.Content, &msg.ClientID, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, rows.Err()
}

// GetMessagesByClient retrieves messages for a specific client
func (s *MessageServicePostgres) GetMessagesByClient(clientID string, limit int) ([]*model.Message, error) {
	query := `SELECT id, content, client_id, created_at FROM messages WHERE client_id = $1 ORDER BY created_at DESC LIMIT $2`
	if s.driver == "sqlite" {
		query = `SELECT id, content, client_id, created_at FROM messages WHERE client_id = ? ORDER BY created_at DESC LIMIT ?`
	}

	rows, err := s.db.Query(query, clientID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		var msg model.Message
		err := rows.Scan(&msg.ID, &msg.Content, &msg.ClientID, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, rows.Err()
}

// GetDB returns the database connection
func (s *MessageServicePostgres) GetDB() *sql.DB {
	return s.db
}

// Close closes the database connection
func (s *MessageServicePostgres) Close() error {
	return s.db.Close()
}
