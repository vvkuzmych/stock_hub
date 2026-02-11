package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"stock_hub_trade/pkg/model"

	"github.com/DATA-DOG/go-sqlmock"
)

// Test helpers

// setupMockDB creates a mock database and sqlmock instance for testing
func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

// setupRepository creates a repository with mock database
func setupRepository(t *testing.T) (*StockOrderDBModel, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, cleanup := setupMockDB(t)
	repo := NewStockOrderDBModel(db, "postgres")

	return repo, mock, cleanup
}

// assertNoError checks that no error occurred
func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// assertError checks that error occurred
func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// assertMockExpectations checks that all mock expectations were met
func assertMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled mock expectations: %v", err)
	}
}

// orderColumns is now defined in stock_order_repository_test_helpers.go for reusability

// Test Create

func TestCreate(t *testing.T) {
	testCases := []struct {
		name          string
		order         *model.StockOrder
		mockSetup     func(sqlmock.Sqlmock, *model.StockOrder)
		expectedID    int64
		expectError   bool
		validateOrder func(*testing.T, *model.StockOrder)
	}{
		{
			name: "Success",
			order: &model.StockOrder{
				UserID:    1,
				Username:  "john",
				Symbol:    "AAPL",
				OrderType: model.OrderTypeBid,
				Price:     150.00,
				Quantity:  10,
				Status:    "open",
				CreatedAt: time.Now(),
			},
			mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
				mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
					WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
			},
			expectedID:  123,
			expectError: false,
			validateOrder: func(t *testing.T, order *model.StockOrder) {
				if order.ID != 123 {
					t.Errorf("Expected order ID to be 123, got: %d", order.ID)
				}
			},
		},
		{
			name: "Database error",
			order: &model.StockOrder{
				UserID:    1,
				Username:  "john",
				Symbol:    "AAPL",
				OrderType: model.OrderTypeBid,
				Price:     150.00,
				Quantity:  10,
				Status:    "open",
				CreatedAt: time.Now(),
			},
			mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
				mock.ExpectQuery(`INSERT INTO stock_orders`).
					WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
					WillReturnError(errors.New("connection lost"))
			},
			expectError: true,
		},
		{
			name: "Context cancelled",
			order: &model.StockOrder{
				UserID:    1,
				Username:  "john",
				Symbol:    "AAPL",
				OrderType: model.OrderTypeBid,
				Price:     150.00,
				Quantity:  10,
				Status:    "open",
				CreatedAt: time.Now(),
			},
			mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
				mock.ExpectQuery(`INSERT INTO stock_orders`).
					WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
					WillReturnError(context.Canceled)
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock, tc.order)

			ctx := context.Background()
			err := repo.Create(ctx, tc.order)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				if tc.validateOrder != nil {
					tc.validateOrder(t, tc.order)
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test GetByID

func TestGetByID(t *testing.T) {
	expectedTime := time.Now()

	testCases := []struct {
		name          string
		orderID       int64
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectNil     bool
		validateOrder func(*testing.T, *model.StockOrder)
	}{
		{
			name:    "Success",
			orderID: 123,
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "username", "symbol", "order_type", "price", "quantity", "status", "created_at"}).
					AddRow(123, 1, "john", "AAPL", "bid", 150.00, 10, "open", expectedTime)
				mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE id = \$1`).
					WithArgs(int64(123)).
					WillReturnRows(rows)
			},
			expectError: false,
			expectNil:   false,
			validateOrder: func(t *testing.T, order *model.StockOrder) {
				if order.ID != 123 {
					t.Errorf("Expected ID 123, got: %d", order.ID)
				}
				if order.Username != "john" {
					t.Errorf("Expected username 'john', got: %s", order.Username)
				}
				if order.Symbol != "AAPL" {
					t.Errorf("Expected symbol 'AAPL', got: %s", order.Symbol)
				}
				if order.OrderType != model.OrderTypeBid {
					t.Errorf("Expected order type 'bid', got: %s", order.OrderType)
				}
				if order.Price != 150.00 {
					t.Errorf("Expected price 150.00, got: %.2f", order.Price)
				}
				if order.Quantity != 10 {
					t.Errorf("Expected quantity 10, got: %d", order.Quantity)
				}
				if order.Status != "open" {
					t.Errorf("Expected status 'open', got: %s", order.Status)
				}
			},
		},
		{
			name:    "Not found",
			orderID: 999,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE id = \$1`).
					WithArgs(int64(999)).
					WillReturnError(sql.ErrNoRows)
			},
			expectError: true,
			expectNil:   true,
		},
		{
			name:    "Database error",
			orderID: 123,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE id = \$1`).
					WithArgs(int64(123)).
					WillReturnError(errors.New("connection timeout"))
			},
			expectError: true,
			expectNil:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock)

			ctx := context.Background()
			order, err := repo.GetByID(ctx, tc.orderID)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}

			if tc.expectNil {
				if order != nil {
					t.Error("Expected nil order")
				}
			} else {
				if order == nil {
					t.Fatal("Expected order, got nil")
				}
				if tc.validateOrder != nil {
					tc.validateOrder(t, order)
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test GetOpenOrders

func TestGetOpenOrders(t *testing.T) {
	const queryPattern = `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE symbol = \$1 AND status = 'open'`

	testCases := []struct {
		name           string
		symbol         string
		mockSetup      func(sqlmock.Sqlmock, string)
		expectError    bool
		expectedCount  int
		validateOrders func(*testing.T, []*model.StockOrder)
	}{
		{
			name:   "Success with multiple orders",
			symbol: "AAPL",
			mockSetup: func(mock sqlmock.Sqlmock, symbol string) {
				rows := sqlmock.NewRows(orderColumns()).
					AddRow(1, int64(1), "john", "AAPL", "bid", 150.00, 10, "open", time.Now()).
					AddRow(2, int64(2), "jane", "AAPL", "ask", 155.00, 5, "open", time.Now())
				mock.ExpectQuery(queryPattern).
					WithArgs(symbol).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 2,
			validateOrders: func(t *testing.T, orders []*model.StockOrder) {
				if orders[0].Symbol != "AAPL" {
					t.Errorf("Expected symbol 'AAPL', got: %s", orders[0].Symbol)
				}
			},
		},
		{
			name:   "Empty result",
			symbol: "TSLA",
			mockSetup: func(mock sqlmock.Sqlmock, symbol string) {
				rows := sqlmock.NewRows(orderColumns())
				mock.ExpectQuery(queryPattern).
					WithArgs(symbol).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:   "Database error",
			symbol: "AAPL",
			mockSetup: func(mock sqlmock.Sqlmock, symbol string) {
				mock.ExpectQuery(queryPattern).
					WithArgs(symbol).
					WillReturnError(errors.New("database error"))
			},
			expectError:   true,
			expectedCount: -1,
		},
		{
			name:   "Scan error",
			symbol: "AAPL",
			mockSetup: func(mock sqlmock.Sqlmock, symbol string) {
				rows := sqlmock.NewRows(orderColumns()).
					AddRow(1, "invalid", "john", "AAPL", "bid", 150.00, 10, "open", time.Now())
				mock.ExpectQuery(queryPattern).
					WithArgs(symbol).
					WillReturnRows(rows)
			},
			expectError:   true,
			expectedCount: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock, tc.symbol)

			ctx := context.Background()
			orders, err := repo.GetOpenOrders(ctx, tc.symbol)

			if tc.expectError {
				assertError(t, err)
				if orders != nil {
					t.Error("Expected nil orders on error")
				}
			} else {
				assertNoError(t, err)
				if len(orders) != tc.expectedCount {
					t.Errorf("Expected %d orders, got: %d", tc.expectedCount, len(orders))
				}
				if tc.validateOrders != nil && len(orders) > 0 {
					tc.validateOrders(t, orders)
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test GetAllOpenOrders

func TestGetAllOpenOrders(t *testing.T) {
	testCases := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectError   bool
		expectedCount int
	}{
		{
			name: "Success with multiple orders",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "username", "symbol", "order_type", "price", "quantity", "status", "created_at"}).
					AddRow(1, 1, "john", "AAPL", "bid", 150.00, 10, "open", time.Now()).
					AddRow(2, 2, "jane", "TSLA", "ask", 800.00, 3, "open", time.Now()).
					AddRow(3, 3, "bob", "GOOGL", "bid", 2800.00, 2, "open", time.Now())
				mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE status = 'open'`).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 3,
		},
		{
			name: "Empty result",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "username", "symbol", "order_type", "price", "quantity", "status", "created_at"})
				mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE status = 'open'`).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock)

			ctx := context.Background()
			orders, err := repo.GetAllOpenOrders(ctx)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				if len(orders) != tc.expectedCount {
					t.Errorf("Expected %d orders, got: %d", tc.expectedCount, len(orders))
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test GetByUserID

func TestGetByUserID(t *testing.T) {
	const queryPattern = `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE user_id = \$1 ORDER BY created_at DESC LIMIT \$2`

	testCases := []struct {
		name           string
		userID         int64
		limit          int
		mockSetup      func(sqlmock.Sqlmock, int64, int)
		expectError    bool
		expectedCount  int
		validateOrders func(*testing.T, []*model.StockOrder)
	}{
		{
			name:   "Success with custom limit",
			userID: 1,
			limit:  100,
			mockSetup: func(mock sqlmock.Sqlmock, userID int64, limit int) {
				rows := sqlmock.NewRows(orderColumns()).
					AddRow(1, int64(1), "john", "AAPL", "bid", 150.00, 10, "open", time.Now()).
					AddRow(2, int64(1), "john", "TSLA", "ask", 800.00, 5, "cancelled", time.Now())
				mock.ExpectQuery(queryPattern).
					WithArgs(userID, limit).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 2,
			validateOrders: func(t *testing.T, orders []*model.StockOrder) {
				if orders[0].UserID != 1 {
					t.Errorf("Expected UserID 1, got: %d", orders[0].UserID)
				}
			},
		},
		{
			name:   "Default limit when zero",
			userID: 1,
			limit:  0,
			mockSetup: func(mock sqlmock.Sqlmock, userID int64, limit int) {
				rows := sqlmock.NewRows(orderColumns())
				mock.ExpectQuery(queryPattern).
					WithArgs(userID, 100).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:   "Custom limit 50",
			userID: 1,
			limit:  50,
			mockSetup: func(mock sqlmock.Sqlmock, userID int64, limit int) {
				rows := sqlmock.NewRows(orderColumns()).
					AddRow(1, int64(1), "john", "AAPL", "bid", 150.00, 10, "open", time.Now())
				mock.ExpectQuery(queryPattern).
					WithArgs(userID, limit).
					WillReturnRows(rows)
			},
			expectError:   false,
			expectedCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			// Setup mock
			effectiveLimit := tc.limit
			if effectiveLimit <= 0 {
				effectiveLimit = 100
			}
			tc.mockSetup(mock, tc.userID, effectiveLimit)

			ctx := context.Background()
			orders, err := repo.GetByUserID(ctx, tc.userID, tc.limit)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				if len(orders) != tc.expectedCount {
					t.Errorf("Expected %d orders, got: %d", tc.expectedCount, len(orders))
				}
				if tc.validateOrders != nil && len(orders) > 0 {
					tc.validateOrders(t, orders)
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test Cancel

func TestCancel(t *testing.T) {
	testCases := []struct {
		name        string
		orderID     int64
		userID      int64
		mockSetup   func(sqlmock.Sqlmock)
		expectError bool
	}{
		{
			name:    "Success",
			orderID: 123,
			userID:  1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
					WithArgs(int64(123), int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
			},
			expectError: false,
		},
		{
			name:    "Order not found",
			orderID: 999,
			userID:  1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
					WithArgs(int64(999), int64(1)).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError: true,
		},
		{
			name:    "Wrong user",
			orderID: 123,
			userID:  2,
			mockSetup: func(mock sqlmock.Sqlmock) {
				// User 2 tries to cancel order belonging to User 1
				mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
					WithArgs(int64(123), int64(2)).
					WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
			},
			expectError: true,
		},
		{
			name:    "Database error",
			orderID: 123,
			userID:  1,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
					WithArgs(int64(123), int64(1)).
					WillReturnError(errors.New("connection timeout"))
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock)

			ctx := context.Background()
			err := repo.Cancel(ctx, tc.orderID, tc.userID)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}

			assertMockExpectations(t, mock)
		})
	}
}

// Test Interface Compliance

func TestInterfaceCompliance(t *testing.T) {
	// Compile-time check
	var _ StockOrderRepository = (*StockOrderDBModel)(nil)
}
