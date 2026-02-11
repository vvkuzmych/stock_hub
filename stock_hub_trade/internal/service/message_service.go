package service

import (
	"database/sql"
	"fmt"
	"time"

	"stock_hub_trade/pkg/migrate"
	"stock_hub_trade/pkg/model"
	"stock_hub_trade/pkg/sqlutil"

	_ "github.com/lib/pq"
)

// MessageService handles message persistence
type MessageService struct {
	db     *sql.DB
	driver string
}

// NewMessageService creates a new message service with database type support
func NewMessageService(driver, dsn, migrationsDir string) (*MessageService, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	service := &MessageService{
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
func (s *MessageService) SaveMessage(content string, clientID string) error {
	query := `INSERT INTO messages (content, client_id, created_at) VALUES (?, ?, ?)`
	query = sqlutil.ConvertPlaceholders(query, s.driver)
	_, err := s.db.Exec(query, content, clientID, time.Now())
	return err
}

// GetMessages retrieves messages from the database
func (s *MessageService) GetMessages(limit int) ([]*model.Message, error) {
	query := `SELECT id, content, client_id, created_at FROM messages ORDER BY created_at DESC LIMIT ?`
	query = sqlutil.ConvertPlaceholders(query, s.driver)

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
	query = sqlutil.ConvertPlaceholders(query, s.driver)

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
