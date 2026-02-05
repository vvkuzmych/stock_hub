-- Down migration: Drop stock_orders table

DROP INDEX IF EXISTS idx_status;
DROP INDEX IF EXISTS idx_order_type;
DROP INDEX IF EXISTS idx_symbol;
DROP INDEX IF EXISTS idx_user_id;
DROP TABLE IF EXISTS stock_orders;
