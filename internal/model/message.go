package model

import "time"

// Message represents a communication record between clients
type Message struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	ClientID  string    `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
}
