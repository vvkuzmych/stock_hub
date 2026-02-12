package model

import "time"

// User represents a registered user
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"` // Email address (optional for backwards compatibility)
	CreatedAt time.Time `json:"created_at"`
}
