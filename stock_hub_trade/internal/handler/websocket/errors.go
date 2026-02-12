package websocket

import "errors"

var (
	ErrInvalidOrderType = errors.New("invalid order type")
	ErrInvalidPrice     = errors.New("price must be greater than 0")
	ErrInvalidQuantity  = errors.New("quantity must be greater than 0")
	ErrEmptySymbol      = errors.New("symbol cannot be empty")
)
