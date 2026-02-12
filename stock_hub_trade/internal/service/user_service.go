package service

import (
	"database/sql"
	"fmt"
	"regexp"

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
func (s *UserService) RegisterUser(username, email string) (*model.User, error) {
	// Validate username
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	// Validate email
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Simple email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Check if user exists
	var user model.User
	var emailNull sql.NullString
	query := sqlutil.ConvertPlaceholders("SELECT id, username, email, created_at FROM users WHERE username = ?", s.driver)
	err := s.db.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &emailNull, &user.CreatedAt)

	if err == nil {
		// User exists, return it
		if emailNull.Valid {
			user.Email = emailNull.String
		}
		return &user, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	// Create new user (PostgreSQL: use RETURNING to get ID)
	insertQuery := sqlutil.ConvertPlaceholders("INSERT INTO users (username, email) VALUES (?, ?)", s.driver)
	err = s.db.QueryRow(insertQuery+" RETURNING id, email", username, email).Scan(&user.ID, &emailNull)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.Username = username
	if emailNull.Valid {
		user.Email = emailNull.String
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	var emailNull sql.NullString
	query := sqlutil.ConvertPlaceholders("SELECT id, username, email, created_at FROM users WHERE id = ?", s.driver)
	err := s.db.QueryRow(query, id).
		Scan(&user.ID, &user.Username, &emailNull, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if emailNull.Valid {
		user.Email = emailNull.String
	}
	return &user, nil
}
