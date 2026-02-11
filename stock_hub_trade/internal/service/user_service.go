package service

import (
	"database/sql"
	"fmt"

	"stock_hub_trade/pkg/model"
	"stock_hub_trade/pkg/sqlutil"

	_ "github.com/lib/pq"
)

// UserService handles user operations
type UserService struct {
	db     *sql.DB
	driver string
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db:     db,
		driver: "postgres", // PostgreSQL only
	}
}

// NewUserServiceWithDriver creates a new user service with explicit driver
func NewUserServiceWithDriver(db *sql.DB, driver string) *UserService {
	return &UserService{
		db:     db,
		driver: driver,
	}
}

// RegisterUser creates a new user or returns existing one
func (s *UserService) RegisterUser(username string) (*model.User, error) {
	// Check if user exists
	var user model.User
	query := sqlutil.ConvertPlaceholders("SELECT id, username, created_at FROM users WHERE username = ?", s.driver)
	err := s.db.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err == nil {
		// User exists, return it
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	// Create new user (PostgreSQL: use RETURNING to get ID)
	insertQuery := sqlutil.ConvertPlaceholders("INSERT INTO users (username) VALUES (?)", s.driver)
	err = s.db.QueryRow(insertQuery+" RETURNING id", username).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.Username = username
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	query := sqlutil.ConvertPlaceholders("SELECT id, username, created_at FROM users WHERE id = ?", s.driver)
	err := s.db.QueryRow(query, id).
		Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
