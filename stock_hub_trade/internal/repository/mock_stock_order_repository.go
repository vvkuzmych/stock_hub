package repository

import (
	"context"
	"fmt"
	"stock_hub_trade/pkg/model"
)

// MockBehavior defines mock functions for customizing behavior in tests
type MockBehavior struct {
	CreateFunc             func(ctx context.Context, order *model.StockOrder) error
	GetByIDFunc            func(ctx context.Context, id int64) (*model.StockOrder, error)
	GetOpenOrdersFunc      func(ctx context.Context, symbol string) ([]*model.StockOrder, error)
	GetAllOpenOrdersFunc   func(ctx context.Context) ([]*model.StockOrder, error)
	GetByUserIDFunc        func(ctx context.Context, userID int64, limit int) ([]*model.StockOrder, error)
	CancelFunc             func(ctx context.Context, orderID, userID int64) error
}

// MockCalls tracks which methods were called for verification in tests
type MockCalls struct {
	CreateCalled           bool
	GetByIDCalled          bool
	GetOpenOrdersCalled    bool
	GetAllOpenOrdersCalled bool
	GetByUserIDCalled      bool
	CancelCalled           bool
}

// MockStockOrderRepository is a mock implementation of StockOrderRepository interface
// Used for testing services without a real database
type MockStockOrderRepository struct {
	Behavior MockBehavior  // Customize behavior
	Calls    MockCalls     // Track calls for verification
}

// Compile-time check that MockStockOrderRepository implements StockOrderRepository interface
var _ StockOrderRepository = (*MockStockOrderRepository)(nil)

// NewMockStockOrderRepository creates a new mock repository
func NewMockStockOrderRepository() *MockStockOrderRepository {
	return &MockStockOrderRepository{}
}

func (m *MockStockOrderRepository) Create(ctx context.Context, order *model.StockOrder) error {
	m.Calls.CreateCalled = true
	if m.Behavior.CreateFunc != nil {
		return m.Behavior.CreateFunc(ctx, order)
	}
	// Default: assign ID and return success
	order.ID = 1
	return nil
}

func (m *MockStockOrderRepository) GetByID(ctx context.Context, id int64) (*model.StockOrder, error) {
	m.Calls.GetByIDCalled = true
	if m.Behavior.GetByIDFunc != nil {
		return m.Behavior.GetByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("order not found")
}

func (m *MockStockOrderRepository) GetOpenOrders(ctx context.Context, symbol string) ([]*model.StockOrder, error) {
	m.Calls.GetOpenOrdersCalled = true
	if m.Behavior.GetOpenOrdersFunc != nil {
		return m.Behavior.GetOpenOrdersFunc(ctx, symbol)
	}
	return []*model.StockOrder{}, nil
}

func (m *MockStockOrderRepository) GetAllOpenOrders(ctx context.Context) ([]*model.StockOrder, error) {
	m.Calls.GetAllOpenOrdersCalled = true
	if m.Behavior.GetAllOpenOrdersFunc != nil {
		return m.Behavior.GetAllOpenOrdersFunc(ctx)
	}
	return []*model.StockOrder{}, nil
}

func (m *MockStockOrderRepository) GetByUserID(ctx context.Context, userID int64, limit int) ([]*model.StockOrder, error) {
	m.Calls.GetByUserIDCalled = true
	if m.Behavior.GetByUserIDFunc != nil {
		return m.Behavior.GetByUserIDFunc(ctx, userID, limit)
	}
	return []*model.StockOrder{}, nil
}

func (m *MockStockOrderRepository) Cancel(ctx context.Context, orderID, userID int64) error {
	m.Calls.CancelCalled = true
	if m.Behavior.CancelFunc != nil {
		return m.Behavior.CancelFunc(ctx, orderID, userID)
	}
	return nil
}
