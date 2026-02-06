-- Down migration: Remove CHECK constraints

ALTER TABLE stock_orders 
DROP CONSTRAINT IF EXISTS check_price_positive;

ALTER TABLE stock_orders 
DROP CONSTRAINT IF EXISTS check_quantity_positive;

ALTER TABLE stock_orders 
DROP CONSTRAINT IF EXISTS check_symbol_not_empty;
