package model

import "time"

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeBid OrderType = "bid" // Buy order
	OrderTypeAsk OrderType = "ask" // Sell order
)

// StockOrder represents a bid or ask order
type StockOrder struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Symbol    string    `json:"symbol"`
	OrderType OrderType `json:"order_type"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"` // "open", "filled", "cancelled"
	CreatedAt time.Time `json:"created_at"`
}
