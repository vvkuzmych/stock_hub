package service

import (
	"strings"
	"testing"

	"stock_hub_trade/pkg/model"
)

func TestCreateOrder_ValidatesPrice(t *testing.T) {
	service := &StockOrderService{
		driver: "postgres",
	}

	testCases := []struct {
		name        string
		price       float64
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Negative price",
			price:       -10.50,
			shouldError: true,
			errorMsg:    "price must be greater than 0",
		},
		{
			name:        "Zero price",
			price:       0.0,
			shouldError: true,
			errorMsg:    "price must be greater than 0",
		},
		{
			name:        "Very small negative price",
			price:       -0.01,
			shouldError: true,
			errorMsg:    "price must be greater than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(
				1,
				"testuser",
				"AAPL",
				model.OrderTypeBid,
				tc.price,
				100,
			)

			if err == nil {
				t.Errorf("Expected error for price %.2f, but got nil", tc.price)
			} else if !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
			}
		})
	}
}

func TestCreateOrder_ValidatesQuantity(t *testing.T) {
	service := &StockOrderService{
		driver: "postgres",
	}

	testCases := []struct {
		name        string
		quantity    int
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Negative quantity",
			quantity:    -100,
			shouldError: true,
			errorMsg:    "quantity must be greater than 0",
		},
		{
			name:        "Zero quantity",
			quantity:    0,
			shouldError: true,
			errorMsg:    "quantity must be greater than 0",
		},
		{
			name:        "Negative one quantity",
			quantity:    -1,
			shouldError: true,
			errorMsg:    "quantity must be greater than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(
				1,
				"testuser",
				"AAPL",
				model.OrderTypeBid,
				150.50,
				tc.quantity,
			)

			if err == nil {
				t.Errorf("Expected error for quantity %d, but got nil", tc.quantity)
			} else if !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
			}
		})
	}
}

func TestCreateOrder_ValidatesSymbol(t *testing.T) {
	service := &StockOrderService{
		driver: "postgres",
	}

	testCases := []struct {
		name        string
		symbol      string
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Empty symbol",
			symbol:      "",
			shouldError: true,
			errorMsg:    "symbol cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(
				1,
				"testuser",
				tc.symbol,
				model.OrderTypeBid,
				150.50,
				100,
			)

			if err == nil {
				t.Errorf("Expected error for symbol '%s', but got nil", tc.symbol)
			} else if !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
			}
		})
	}
}

func TestCreateOrder_ValidatesOrderType(t *testing.T) {
	service := &StockOrderService{
		driver: "postgres",
	}

	testCases := []struct {
		name        string
		orderType   model.OrderType
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Invalid order type",
			orderType:   model.OrderType("invalid"),
			shouldError: true,
			errorMsg:    "invalid order type",
		},
		{
			name:        "Empty order type",
			orderType:   model.OrderType(""),
			shouldError: true,
			errorMsg:    "invalid order type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(
				1,
				"testuser",
				"AAPL",
				tc.orderType,
				150.50,
				100,
			)

			if err == nil {
				t.Errorf("Expected error for order type '%s', but got nil", tc.orderType)
			} else if !strings.Contains(err.Error(), tc.errorMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errorMsg, err)
			}
		})
	}
}

func TestCreateOrder_CombinedValidation(t *testing.T) {
	service := &StockOrderService{
		driver: "postgres",
	}

	// Test multiple invalid inputs at once
	testCases := []struct {
		name      string
		userID    int64
		username  string
		symbol    string
		orderType model.OrderType
		price     float64
		quantity  int
		expectErr string
	}{
		{
			name:      "All invalid",
			userID:    1,
			username:  "test",
			symbol:    "",
			orderType: model.OrderType("invalid"),
			price:     -10.0,
			quantity:  -100,
			expectErr: "price must be greater than 0", // First validation to fail
		},
		{
			name:      "Invalid price and quantity",
			userID:    1,
			username:  "test",
			symbol:    "AAPL",
			orderType: model.OrderTypeBid,
			price:     0,
			quantity:  0,
			expectErr: "price must be greater than 0",
		},
		{
			name:      "Valid price, invalid quantity",
			userID:    1,
			username:  "test",
			symbol:    "AAPL",
			orderType: model.OrderTypeBid,
			price:     100.0,
			quantity:  -50,
			expectErr: "quantity must be greater than 0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateOrder(
				tc.userID,
				tc.username,
				tc.symbol,
				tc.orderType,
				tc.price,
				tc.quantity,
			)

			if err == nil {
				t.Error("Expected validation error but got nil")
			} else if !strings.Contains(err.Error(), tc.expectErr) {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectErr, err)
			}
		})
	}
}
