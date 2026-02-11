package service

import (
	"context"
	"fmt"
	"time"

	"stock_hub_trade/internal/repository"
	"stock_hub_trade/pkg/model"
)

// StockOrderServiceV2 handles stock order business logic (separated from data access)
type StockOrderServiceV2 struct {
	repo repository.StockOrderRepository
}

// NewStockOrderServiceV2 creates a new stock order service with repository pattern
func NewStockOrderServiceV2(repo repository.StockOrderRepository) *StockOrderServiceV2 {
	return &StockOrderServiceV2{
		repo: repo,
	}
}

// CreateOrder creates a new bid or ask order with business logic validation
func (s *StockOrderServiceV2) CreateOrder(
	ctx context.Context,
	userID int64,
	username, symbol string,
	orderType model.OrderType,
	price float64,
	quantity int,
) (*model.StockOrder, error) {
	// Business logic: Validate input parameters
	if err := s.validateOrderInput(price, quantity, symbol, orderType); err != nil {
		return nil, err
	}

	// Business logic: Additional validations (can be extended)
	// For example: check user balance, trading hours, circuit breakers, etc.

	// Create order model
	order := &model.StockOrder{
		UserID:    userID,
		Username:  username,
		Symbol:    symbol,
		OrderType: orderType,
		Price:     price,
		Quantity:  quantity,
		Status:    "open",
		CreatedAt: time.Now(),
	}

	// Data access: Repository handles DB interaction
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID
func (s *StockOrderServiceV2) GetOrderByID(ctx context.Context, orderID int64) (*model.StockOrder, error) {
	if orderID <= 0 {
		return nil, fmt.Errorf("invalid order ID")
	}

	return s.repo.GetByID(ctx, orderID)
}

// GetOpenOrders retrieves all open orders for a symbol
func (s *StockOrderServiceV2) GetOpenOrders(ctx context.Context, symbol string) ([]*model.StockOrder, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}

	return s.repo.GetOpenOrders(ctx, symbol)
}

// GetAllOpenOrders retrieves all open orders
func (s *StockOrderServiceV2) GetAllOpenOrders(ctx context.Context) ([]*model.StockOrder, error) {
	return s.repo.GetAllOpenOrders(ctx)
}

// GetOrdersByUserID retrieves orders for a specific user
func (s *StockOrderServiceV2) GetOrdersByUserID(ctx context.Context, userID int64, limit int) ([]*model.StockOrder, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	return s.repo.GetByUserID(ctx, userID, limit)
}

// CancelOrder cancels an order
func (s *StockOrderServiceV2) CancelOrder(ctx context.Context, orderID, userID int64) error {
	// Business logic: Validate IDs
	if orderID <= 0 {
		return fmt.Errorf("invalid order ID")
	}
	if userID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// Business logic: Additional checks can be added here
	// For example: check if order is still open, within cancellation window, etc.

	return s.repo.Cancel(ctx, orderID, userID)
}

// validateOrderInput validates order creation parameters
func (s *StockOrderServiceV2) validateOrderInput(
	price float64,
	quantity int,
	symbol string,
	orderType model.OrderType,
) error {
	if price <= 0 {
		return fmt.Errorf("price must be greater than 0, got: %.2f", price)
	}
	if quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0, got: %d", quantity)
	}
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if orderType != model.OrderTypeBid && orderType != model.OrderTypeAsk {
		return fmt.Errorf("invalid order type: %s", orderType)
	}

	// Additional business rules can be added here:
	// - Max order size
	// - Price limits (circuit breakers)
	// - Trading hours check
	// - Symbol validation against allowed list

	return nil
}
