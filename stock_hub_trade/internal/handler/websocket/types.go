package websocket

import "encoding/json"

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type     string          `json:"type"` // "register", "message", "order", "cancel_order"
	Data     json.RawMessage `json:"data"`
	Username string          `json:"username,omitempty"`
	UserID   int64           `json:"user_id,omitempty"`
}

// RegisterData for user registration
type RegisterData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// OrderData for stock orders
type OrderData struct {
	Symbol    string  `json:"symbol"`
	OrderType string  `json:"order_type"` // "bid" or "ask"
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// CancelOrderData for cancelling orders
type CancelOrderData struct {
	OrderID int64 `json:"order_id"`
}
