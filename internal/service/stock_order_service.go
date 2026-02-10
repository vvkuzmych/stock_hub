package service

import (
	"database/sql"
	"fmt"
	"time"

	"stock_hub/internal/model"
	"stock_hub/pkg/sqlutil"

	_ "github.com/lib/pq"
)

// StockOrderService handles stock order operations
type StockOrderService struct {
	db     *sql.DB
	driver string
}

// NewStockOrderService creates a new stock order service
func NewStockOrderService(db *sql.DB) *StockOrderService {
	return &StockOrderService{
		db:     db,
		driver: "postgres", // PostgreSQL only
	}
}

// NewStockOrderServiceWithDriver creates a new stock order service with explicit driver
func NewStockOrderServiceWithDriver(db *sql.DB, driver string) *StockOrderService {
	return &StockOrderService{
		db:     db,
		driver: driver,
	}
}

// CreateOrder creates a new bid or ask order
func (s *StockOrderService) CreateOrder(userID int64, username, symbol string, orderType model.OrderType, price float64, quantity int) (*model.StockOrder, error) {
	// Validate input parameters
	if price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0, got: %.2f", price)
	}
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be greater than 0, got: %d", quantity)
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol cannot be empty")
	}
	if orderType != model.OrderTypeBid && orderType != model.OrderTypeAsk {
		return nil, fmt.Errorf("invalid order type: %s", orderType)
	}

	query := `INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity, status, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, 'open', ?)`
	query = sqlutil.ConvertPlaceholders(query, s.driver)
	
	now := time.Now()
	
	// PostgreSQL: use RETURNING to get ID
	var id int64
	err := s.db.QueryRow(query+" RETURNING id", userID, username, symbol, string(orderType), price, quantity, now).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	return &model.StockOrder{
		ID:        id,
		UserID:    userID,
		Username:  username,
		Symbol:    symbol,
		OrderType: orderType,
		Price:     price,
		Quantity:  quantity,
		Status:    "open",
		CreatedAt: now,
	}, nil
}

// GetOpenOrders retrieves all open orders for a symbol
func (s *StockOrderService) GetOpenOrders(symbol string) ([]*model.StockOrder, error) {
	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE symbol = ? AND status = 'open' 
	          ORDER BY 
	            CASE WHEN order_type = 'bid' THEN price END DESC,
	            CASE WHEN order_type = 'ask' THEN price END ASC,
	            created_at ASC`
	query = sqlutil.ConvertPlaceholders(query, s.driver)
	
	rows, err := s.db.Query(query, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []*model.StockOrder
	for rows.Next() {
		var order model.StockOrder
		var orderTypeStr string
		err := rows.Scan(&order.ID, &order.UserID, &order.Username, &order.Symbol, 
			&orderTypeStr, &order.Price, &order.Quantity, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		order.OrderType = model.OrderType(orderTypeStr)
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// GetAllOpenOrders retrieves all open orders
func (s *StockOrderService) GetAllOpenOrders() ([]*model.StockOrder, error) {
	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE status = 'open' 
	          ORDER BY symbol, 
	            CASE WHEN order_type = 'bid' THEN price END DESC,
	            CASE WHEN order_type = 'ask' THEN price END ASC,
	            created_at ASC`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []*model.StockOrder
	for rows.Next() {
		var order model.StockOrder
		var orderTypeStr string
		err := rows.Scan(&order.ID, &order.UserID, &order.Username, &order.Symbol, 
			&orderTypeStr, &order.Price, &order.Quantity, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		order.OrderType = model.OrderType(orderTypeStr)
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// CancelOrder cancels an order
func (s *StockOrderService) CancelOrder(orderID, userID int64) error {
	query := `UPDATE stock_orders SET status = 'cancelled' WHERE id = ? AND user_id = ? AND status = 'open'`
	query = sqlutil.ConvertPlaceholders(query, s.driver)
	result, err := s.db.Exec(query, orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("order not found or already cancelled")
	}
	
	return nil
}
