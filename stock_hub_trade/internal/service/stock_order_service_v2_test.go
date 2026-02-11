package service

import (
	"context"
	"fmt"
	"testing"

	"stock_hub_trade/internal/repository"
	"stock_hub_trade/pkg/model"
)

// Test helpers

// setupTest creates a test fixture with mock repository and service
func setupTest(t *testing.T) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
	t.Helper()
	
	mockRepo := repository.NewMockStockOrderRepository()
	service := NewStockOrderServiceV2(mockRepo)
	ctx := context.Background()
	
	return service, mockRepo, ctx
}

// setupTestWithBehavior creates a test fixture with custom repository behavior
func setupTestWithBehavior(t *testing.T, behavior repository.MockBehavior) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
	t.Helper()
	
	mockRepo := repository.NewMockStockOrderRepository()
	mockRepo.Behavior = behavior
	service := NewStockOrderServiceV2(mockRepo)
	ctx := context.Background()
	
	return service, mockRepo, ctx
}

// TestCreateOrderV2_ValidatesPrice tests price validation in service layer
func TestCreateOrderV2_ValidatesPrice(t *testing.T) {
	service, mockRepo, ctx := setupTest(t)

	testCases := []struct {
		name        string
		price       float64
		shouldError bool
	}{
		{"Zero price", 0, true},
		{"Negative price", -100.50, true},
		{"Valid price", 150.25, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, tc.price, 10)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Expected error for price %.2f, but got nil", tc.price)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for price %.2f, but got: %v", tc.price, err)
				}
				// Verify repository was called
				if !mockRepo.Calls.CreateCalled {
					t.Error("Expected repository Create to be called")
				}
			}
		})
	}
}

// TestCreateOrderV2_ValidatesQuantity tests quantity validation
func TestCreateOrderV2_ValidatesQuantity(t *testing.T) {
	service, _, ctx := setupTest(t)

	testCases := []struct {
		name        string
		quantity    int
		shouldError bool
	}{
		{"Zero quantity", 0, true},
		{"Negative quantity", -5, true},
		{"Valid quantity", 100, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, tc.quantity)

			if tc.shouldError && err == nil {
				t.Errorf("Expected error for quantity %d", tc.quantity)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error for quantity %d, got: %v", tc.quantity, err)
			}
		})
	}
}

// TestCreateOrderV2_ValidatesSymbol tests symbol validation
func TestCreateOrderV2_ValidatesSymbol(t *testing.T) {
	service, _, ctx := setupTest(t)

	_, err := service.CreateOrder(ctx, 1, "john", "", model.OrderTypeBid, 150.00, 10)

	if err == nil {
		t.Error("Expected error for empty symbol")
	}
}

// TestCreateOrderV2_ValidatesOrderType tests order type validation
func TestCreateOrderV2_ValidatesOrderType(t *testing.T) {
	service, _, ctx := setupTest(t)

	_, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderType("invalid"), 150.00, 10)

	if err == nil {
		t.Error("Expected error for invalid order type")
	}
}

// TestCreateOrderV2_RepositoryError tests handling of repository errors
func TestCreateOrderV2_RepositoryError(t *testing.T) {
	service, mockRepo, ctx := setupTestWithBehavior(t, repository.MockBehavior{
		CreateFunc: func(ctx context.Context, order *model.StockOrder) error {
			return fmt.Errorf("database connection failed")
		},
	})

	_, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)

	if err == nil {
		t.Error("Expected error when repository fails")
	}

	if !mockRepo.Calls.CreateCalled {
		t.Error("Expected repository Create to be called")
	}
}

// TestCreateOrderV2_Success tests successful order creation
func TestCreateOrderV2_Success(t *testing.T) {
	service, mockRepo, ctx := setupTestWithBehavior(t, repository.MockBehavior{
		CreateFunc: func(ctx context.Context, order *model.StockOrder) error {
			// Simulate database assigning ID
			order.ID = 123
			return nil
		},
	})

	order, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if order == nil {
		t.Fatal("Expected order to be created")
	}

	if order.ID != 123 {
		t.Errorf("Expected order ID to be 123, got: %d", order.ID)
	}

	if order.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got: %s", order.Symbol)
	}

	if !mockRepo.Calls.CreateCalled {
		t.Error("Expected repository Create to be called")
	}
}

// TestGetOrderByID_ValidatesID tests ID validation
func TestGetOrderByID_ValidatesID(t *testing.T) {
	service, mockRepo, ctx := setupTest(t)

	testCases := []struct {
		name        string
		orderID     int64
		shouldError bool
	}{
		{"Zero ID", 0, true},
		{"Negative ID", -1, true},
		{"Valid ID", 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.GetOrderByID(ctx, tc.orderID)

			if tc.shouldError && err == nil {
				t.Errorf("Expected error for order ID %d", tc.orderID)
			}
			if !tc.shouldError && err != nil {
				// Repository returns "not found", which is expected
				if mockRepo.Calls.GetByIDCalled {
					// OK - repository was called
				} else {
					t.Error("Expected repository GetByID to be called")
				}
			}
		})
	}
}

// TestCancelOrder_ValidatesIDs tests ID validation for cancellation
func TestCancelOrder_ValidatesIDs(t *testing.T) {
	service, _, ctx := setupTest(t)

	testCases := []struct {
		name        string
		orderID     int64
		userID      int64
		shouldError bool
	}{
		{"Zero order ID", 0, 1, true},
		{"Negative order ID", -1, 1, true},
		{"Zero user ID", 1, 0, true},
		{"Negative user ID", 1, -1, true},
		{"Valid IDs", 1, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.CancelOrder(ctx, tc.orderID, tc.userID)

			if tc.shouldError && err == nil {
				t.Errorf("Expected error for orderID=%d, userID=%d", tc.orderID, tc.userID)
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

// TestGetOpenOrders_ValidatesSymbol tests symbol validation
func TestGetOpenOrders_ValidatesSymbol(t *testing.T) {
	service, mockRepo, ctx := setupTest(t)

	_, err := service.GetOpenOrders(ctx, "")

	if err == nil {
		t.Error("Expected error for empty symbol")
	}

	if mockRepo.Calls.GetOpenOrdersCalled {
		t.Error("Repository should not be called for invalid input")
	}
}
