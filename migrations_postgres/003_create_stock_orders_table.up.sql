-- Up migration: Create stock_orders table

CREATE TABLE IF NOT EXISTS stock_orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    symbol TEXT NOT NULL,
    order_type TEXT NOT NULL CHECK(order_type IN ('bid', 'ask')),
    price DECIMAL(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'open' CHECK(status IN ('open', 'filled', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_id ON stock_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_symbol ON stock_orders(symbol);
CREATE INDEX IF NOT EXISTS idx_order_type ON stock_orders(order_type);
CREATE INDEX IF NOT EXISTS idx_status ON stock_orders(status);
CREATE INDEX IF NOT EXISTS idx_symbol_status ON stock_orders(symbol, status);
