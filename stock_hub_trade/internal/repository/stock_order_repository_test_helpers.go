package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"stock_hub_trade/pkg/model"
)

// orderColumns returns the standard column names for stock_orders table
func orderColumns() []string {
	return []string{"id", "user_id", "username", "symbol", "order_type", "price", "quantity", "status", "created_at"}
}

// Mock expectation helpers for cleaner test code

// expectCreateSuccess sets up mock for successful Create operation
func expectCreateSuccess(mock sqlmock.Sqlmock, order *model.StockOrder, returnID int64) {
	mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
		WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnID))
}

// expectCreateError sets up mock for failed Create operation
func expectCreateError(mock sqlmock.Sqlmock, order *model.StockOrder, err error) {
	mock.ExpectQuery(`INSERT INTO stock_orders`).
		WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
		WillReturnError(err)
}

// expectGetByIDSuccess sets up mock for successful GetByID operation
func expectGetByIDSuccess(mock sqlmock.Sqlmock, orderID int64, order *model.StockOrder) {
	rows := sqlmock.NewRows(orderColumns()).
		AddRow(order.ID, order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, order.CreatedAt)
	mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnRows(rows)
}

// expectGetByIDError sets up mock for failed GetByID operation
func expectGetByIDError(mock sqlmock.Sqlmock, orderID int64, err error) {
	mock.ExpectQuery(`SELECT id, user_id, username, symbol, order_type, price, quantity, status, created_at FROM stock_orders WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnError(err)
}

// expectCancelSuccess sets up mock for successful Cancel operation
func expectCancelSuccess(mock sqlmock.Sqlmock, orderID, userID int64) {
	mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
		WithArgs(orderID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

// expectCancelNotFound sets up mock for Cancel when order not found
func expectCancelNotFound(mock sqlmock.Sqlmock, orderID, userID int64) {
	mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
		WithArgs(orderID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
}

// expectCancelError sets up mock for failed Cancel operation
func expectCancelError(mock sqlmock.Sqlmock, orderID, userID int64, err error) {
	mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
		WithArgs(orderID, userID).
		WillReturnError(err)
}
