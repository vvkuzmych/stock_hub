package service

import (
	"database/sql"
	"fmt"
	"time"

	"stock_hub/internal/model"
	"stock_hub/pkg/migrate"

	_ "modernc.org/sqlite"
)

// MessageService handles message persistence
type MessageService struct {
	db *sql.DB
}

// NewMessageService creates a new message service
func NewMessageService(dbPath string) (*MessageService, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	service := &MessageService{db: db}

	// Run migrations
	migrator := migrate.NewMigrator(db, "migrations")
	if err := migrator.Up(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return service, nil
}

// SaveMessage saves a message to the database
func (s *MessageService) SaveMessage(content string, clientID string) error {
	query := `INSERT INTO messages (content, client_id, created_at) VALUES (?, ?, ?)`
	_, err := s.db.Exec(query, content, clientID, time.Now())
	return err
}

// GetMessages retrieves messages from the database
func (s *MessageService) GetMessages(limit int) ([]*model.Message, error) {
	query := `SELECT id, content, client_id, created_at FROM messages ORDER BY created_at DESC LIMIT ?`
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
func (s *MessageService) GetMessagesByClient(clientID string, limit int) ([]*model.Message, error) {
	query := `SELECT id, content, client_id, created_at FROM messages WHERE client_id = ? ORDER BY created_at DESC LIMIT ?`
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
func (s *MessageService) GetDB() *sql.DB {
	return s.db
}

// Close closes the database connection
func (s *MessageService) Close() error {
	return s.db.Close()
}
