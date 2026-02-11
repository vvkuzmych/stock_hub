-- Up migration: Add CHECK constraints for price and quantity validation

-- Add CHECK constraint to ensure price is positive
ALTER TABLE stock_orders 
ADD CONSTRAINT check_price_positive 
CHECK (price > 0);

-- Add CHECK constraint to ensure quantity is positive
ALTER TABLE stock_orders 
ADD CONSTRAINT check_quantity_positive 
CHECK (quantity > 0);

-- Add CHECK constraint to ensure symbol is not empty
ALTER TABLE stock_orders 
ADD CONSTRAINT check_symbol_not_empty 
CHECK (symbol != '');
