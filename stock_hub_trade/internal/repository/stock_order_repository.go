package repository

import (
	"context"
	"database/sql"
	"fmt"

	"stock_hub_trade/pkg/model"
	"stock_hub_trade/pkg/sqlutil"
)

// StockOrderRepository defines the interface for stock order data access
type StockOrderRepository interface {
	Create(ctx context.Context, order *model.StockOrder) error
	GetByID(ctx context.Context, id int64) (*model.StockOrder, error)
	GetOpenOrders(ctx context.Context, symbol string) ([]*model.StockOrder, error)
	GetAllOpenOrders(ctx context.Context) ([]*model.StockOrder, error)
	GetByUserID(ctx context.Context, userID int64, limit int) ([]*model.StockOrder, error)
	Cancel(ctx context.Context, orderID, userID int64) error
}

// StockOrderDBModel is the PostgreSQL implementation of StockOrderRepository
type StockOrderDBModel struct {
	DB     *sql.DB
	Driver string
}

// Compile-time check that StockOrderDBModel implements StockOrderRepository interface
var _ StockOrderRepository = (*StockOrderDBModel)(nil)

// NewStockOrderDBModel creates a new stock order repository
func NewStockOrderDBModel(db *sql.DB, driver string) *StockOrderDBModel {
	return &StockOrderDBModel{
		DB:     db,
		Driver: driver,
	}
}

// Create inserts a new stock order into the database
func (m *StockOrderDBModel) Create(ctx context.Context, order *model.StockOrder) error {
	query := `INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity, status, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	query = sqlutil.ConvertPlaceholders(query, m.Driver)

	// PostgreSQL: use RETURNING to get ID
	err := m.DB.QueryRowContext(ctx, query+" RETURNING id",
		order.UserID,
		order.Username,
		order.Symbol,
		string(order.OrderType),
		order.Price,
		order.Quantity,
		order.Status,
		order.CreatedAt,
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetByID retrieves an order by its ID
func (m *StockOrderDBModel) GetByID(ctx context.Context, id int64) (*model.StockOrder, error) {
	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE id = ?`
	query = sqlutil.ConvertPlaceholders(query, m.Driver)

	var order model.StockOrder
	var orderTypeStr string
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Username,
		&order.Symbol,
		&orderTypeStr,
		&order.Price,
		&order.Quantity,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.OrderType = model.OrderType(orderTypeStr)
	return &order, nil
}

// GetOpenOrders retrieves all open orders for a specific symbol
func (m *StockOrderDBModel) GetOpenOrders(ctx context.Context, symbol string) ([]*model.StockOrder, error) {
	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE symbol = ? AND status = 'open' 
	          ORDER BY 
	            CASE WHEN order_type = 'bid' THEN price END DESC,
	            CASE WHEN order_type = 'ask' THEN price END ASC,
	            created_at ASC`
	query = sqlutil.ConvertPlaceholders(query, m.Driver)

	rows, err := m.DB.QueryContext(ctx, query, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	return m.scanOrders(rows)
}

// GetAllOpenOrders retrieves all open orders
func (m *StockOrderDBModel) GetAllOpenOrders(ctx context.Context) ([]*model.StockOrder, error) {
	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE status = 'open' 
	          ORDER BY symbol, 
	            CASE WHEN order_type = 'bid' THEN price END DESC,
	            CASE WHEN order_type = 'ask' THEN price END ASC,
	            created_at ASC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	return m.scanOrders(rows)
}

// GetByUserID retrieves orders for a specific user
func (m *StockOrderDBModel) GetByUserID(ctx context.Context, userID int64, limit int) ([]*model.StockOrder, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at 
	          FROM stock_orders 
	          WHERE user_id = ? 
	          ORDER BY created_at DESC 
	          LIMIT ?`
	query = sqlutil.ConvertPlaceholders(query, m.Driver)

	rows, err := m.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	return m.scanOrders(rows)
}

// Cancel cancels an order
func (m *StockOrderDBModel) Cancel(ctx context.Context, orderID, userID int64) error {
	query := `UPDATE stock_orders SET status = 'cancelled' WHERE id = ? AND user_id = ? AND status = 'open'`
	query = sqlutil.ConvertPlaceholders(query, m.Driver)

	result, err := m.DB.ExecContext(ctx, query, orderID, userID)
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

// scanOrders is a helper function to scan multiple orders from rows
func (m *StockOrderDBModel) scanOrders(rows *sql.Rows) ([]*model.StockOrder, error) {
	var orders []*model.StockOrder
	for rows.Next() {
		var order model.StockOrder
		var orderTypeStr string
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Username,
			&order.Symbol,
			&orderTypeStr,
			&order.Price,
			&order.Quantity,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		order.OrderType = model.OrderType(orderTypeStr)
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}
