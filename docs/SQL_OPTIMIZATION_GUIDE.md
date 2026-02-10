# SQL Оптимізація та Performance Tuning

Детальний гайд по оптимізації SQL запитів та покращенню продуктивності бази даних.

---

## 📖 Зміст

1. [Проблеми продуктивності](#проблеми-продуктивності)
2. [Індексування](#індексування)
3. [Query Optimization](#query-optimization)
4. [N+1 Problem](#n1-problem)
5. [Connection Pooling](#connection-pooling)
6. [Caching Strategies](#caching-strategies)
7. [Partitioning](#partitioning)
8. [Best Practices](#best-practices)

---

## Проблеми продуктивності

### ❌ Повільний запит: Full table scan
```sql
-- Погано: Full table scan
SELECT * FROM stock_orders WHERE symbol = 'AAPL';
```

### ✅ Оптимізовано: З індексом
```sql
-- Створити індекс
CREATE INDEX idx_symbol ON stock_orders(symbol);

-- Тепер швидко
SELECT * FROM stock_orders WHERE symbol = 'AAPL';
```

### Перевірити чи використовується індекс
```sql
-- PostgreSQL
EXPLAIN (ANALYZE, BUFFERS) 
SELECT * FROM stock_orders WHERE symbol = 'AAPL';

-- MySQL
EXPLAIN 
SELECT * FROM stock_orders WHERE symbol = 'AAPL';

-- SQLite
EXPLAIN QUERY PLAN 
SELECT * FROM stock_orders WHERE symbol = 'AAPL';
```

---

## Індексування

### Single Column Index
```sql
CREATE INDEX idx_symbol ON stock_orders(symbol);
CREATE INDEX idx_status ON stock_orders(status);
CREATE INDEX idx_user_id ON stock_orders(user_id);
```

### Composite Index (кращий для multiple conditions)
```sql
-- Порядок важливий! Найселективніша колонка перша
CREATE INDEX idx_symbol_status ON stock_orders(symbol, status);
CREATE INDEX idx_symbol_type_status ON stock_orders(symbol, order_type, status);
```

### Приклад використання composite index
```sql
-- ✅ Використає idx_symbol_status
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open';

-- ✅ Використає idx_symbol_status (partial)
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL';

-- ❌ НЕ використає idx_symbol_status
SELECT * FROM stock_orders 
WHERE status = 'open';
```

### Unique Index
```sql
CREATE UNIQUE INDEX idx_unique_username ON users(username);
```

### Partial Index (PostgreSQL)
```sql
-- Індекс тільки для відкритих ордерів
CREATE INDEX idx_open_orders 
ON stock_orders(symbol, price) 
WHERE status = 'open';
```

### Covering Index (include columns)
```sql
-- PostgreSQL
CREATE INDEX idx_symbol_covering 
ON stock_orders(symbol) 
INCLUDE (price, quantity, status);
```

### Full-Text Search Index
```sql
-- PostgreSQL
CREATE INDEX idx_messages_content_fts 
ON messages USING GIN(to_tsvector('english', content));

-- Query with full-text search
SELECT * FROM messages 
WHERE to_tsvector('english', content) @@ to_tsquery('error & database');

-- MySQL
CREATE FULLTEXT INDEX idx_content_fulltext ON messages(content);
SELECT * FROM messages WHERE MATCH(content) AGAINST('error database');
```

### Коли НЕ потрібні індекси

❌ **Маленькі таблиці** (< 1000 рядків)
```sql
-- Індекс не потрібен для маленької довідкової таблиці
-- CREATE INDEX idx_status ON order_statuses(name); -- не потрібно
```

❌ **Колонки з низькою селективністю**
```sql
-- Погано: тільки 2 можливих значення
CREATE INDEX idx_order_type ON stock_orders(order_type); -- bid/ask

-- Краще: composite index
CREATE INDEX idx_type_status ON stock_orders(order_type, status);
```

❌ **Колонки що часто змінюються**
```sql
-- Погано: часті updates повільнять INSERT/UPDATE
CREATE INDEX idx_updated_at ON stock_orders(updated_at);
```

---

## Query Optimization

### ❌ SELECT * (Погано)
```sql
-- Погано: вибирає всі колонки
SELECT * FROM stock_orders WHERE symbol = 'AAPL';
```

### ✅ SELECT конкретні колонки (Добре)
```sql
-- Добре: тільки потрібні колонки
SELECT id, symbol, price, quantity, status 
FROM stock_orders 
WHERE symbol = 'AAPL';
```

### ❌ Функція в WHERE (ламає індекс)
```sql
-- Погано: не може використати індекс на created_at
SELECT * FROM stock_orders 
WHERE DATE(created_at) = '2024-01-01';
```

### ✅ Range замість функції
```sql
-- Добре: може використати індекс
SELECT * FROM stock_orders 
WHERE created_at >= '2024-01-01' 
  AND created_at < '2024-01-02';
```

### ❌ OR в WHERE (повільно)
```sql
-- Погано: може не використати індекси
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' OR symbol = 'GOOGL' OR symbol = 'MSFT';
```

### ✅ IN замість OR
```sql
-- Добре: краща оптимізація
SELECT * FROM stock_orders 
WHERE symbol IN ('AAPL', 'GOOGL', 'MSFT');
```

### ❌ NOT IN з підзапитом (дуже повільно)
```sql
-- Погано: дуже повільно з великими даними
SELECT * FROM users 
WHERE id NOT IN (SELECT user_id FROM stock_orders);
```

### ✅ LEFT JOIN + IS NULL
```sql
-- Добре: набагато швидше
SELECT u.* 
FROM users u
LEFT JOIN stock_orders so ON u.id = so.user_id
WHERE so.user_id IS NULL;
```

### ❌ Підзапит в SELECT (N+1)
```sql
-- Погано: виконується підзапит для кожного рядка
SELECT 
    u.username,
    (SELECT COUNT(*) FROM stock_orders WHERE user_id = u.id) as order_count
FROM users u;
```

### ✅ JOIN + GROUP BY
```sql
-- Добре: один запит
SELECT 
    u.username,
    COUNT(so.id) as order_count
FROM users u
LEFT JOIN stock_orders so ON u.id = so.user_id
GROUP BY u.id, u.username;
```

### LIMIT для pagination
```sql
-- Добре: обмежити результати
SELECT * FROM stock_orders 
WHERE status = 'open'
ORDER BY created_at DESC
LIMIT 20 OFFSET 0;

-- Ще краще: keyset pagination (для великих offset)
SELECT * FROM stock_orders 
WHERE status = 'open' 
  AND id < 1000  -- останній id з попередньої сторінки
ORDER BY id DESC
LIMIT 20;
```

---

## N+1 Problem

### ❌ N+1 Problem (дуже погано!)
```go
// Погано: 1 запит для users + N запитів для orders
users := GetAllUsers()  // 1 query
for _, user := range users {
    orders := GetOrdersByUserID(user.ID)  // N queries
    // ...
}
```

### ✅ Вирішення: Eager Loading
```go
// Добре: тільки 1 запит
query := `
    SELECT 
        u.id, u.username,
        so.id as order_id, so.symbol, so.price
    FROM users u
    LEFT JOIN stock_orders so ON u.id = so.user_id
`
// Parse results into map[userID][]orders
```

### ✅ Або batch loading
```go
// 1. Отримати users
users := GetAllUsers()  // 1 query

// 2. Отримати всі orders за 1 запит
userIDs := extractIDs(users)
orders := GetOrdersByUserIDs(userIDs)  // 1 query with IN clause

// 3. Згрупувати в пам'яті
ordersByUser := groupByUserID(orders)
```

SQL для batch loading:
```sql
SELECT * FROM stock_orders 
WHERE user_id IN (1, 2, 3, 4, 5, ..., 100);
```

---

## Connection Pooling

### Без пулу (погано)
```go
// Погано: створення нового connection на кожен запит
func GetUser(id int) {
    db, err := sql.Open("sqlite3", "stock_hub.db")
    defer db.Close()
    // query...
}
```

### З пулом (добре)
```go
// Добре: reuse connections
var db *sql.DB

func init() {
    db, _ = sql.Open("sqlite3", "stock_hub.db")
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
}

func GetUser(id int) {
    // використовує connection з пулу
    row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
}
```

### Optimal connection pool settings

**PostgreSQL:**
```go
db.SetMaxOpenConns(25)      // max connections
db.SetMaxIdleConns(5)       // idle connections в пулі
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(10 * time.Minute)
```

**MySQL:**
```go
db.SetMaxOpenConns(50)      // MySQL can handle more
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(10 * time.Minute)
```

**SQLite:**
```go
db.SetMaxOpenConns(1)       // SQLite is single-writer
db.SetMaxIdleConns(1)
```

---

## Caching Strategies

### 1. Query Result Caching (Application Level)
```go
var orderCache = make(map[string][]Order)
var cacheMutex sync.RWMutex

func GetOpenOrdersForSymbol(symbol string) []Order {
    cacheKey := fmt.Sprintf("orders:%s:open", symbol)
    
    // Check cache
    cacheMutex.RLock()
    if cached, ok := orderCache[cacheKey]; ok {
        cacheMutex.RUnlock()
        return cached
    }
    cacheMutex.RUnlock()
    
    // Query DB
    orders := queryDB(symbol)
    
    // Store in cache
    cacheMutex.Lock()
    orderCache[cacheKey] = orders
    cacheMutex.Unlock()
    
    return orders
}
```

### 2. Materialized View (Database Level)
```sql
-- PostgreSQL
CREATE MATERIALIZED VIEW order_book_summary AS
SELECT 
    symbol,
    order_type,
    COUNT(*) as order_count,
    SUM(quantity) as total_quantity,
    AVG(price) as avg_price
FROM stock_orders
WHERE status = 'open'
GROUP BY symbol, order_type;

-- Оновити view
REFRESH MATERIALIZED VIEW order_book_summary;

-- Query materialized view (дуже швидко!)
SELECT * FROM order_book_summary WHERE symbol = 'AAPL';
```

### 3. Redis для Order Book
```go
// Cache order book в Redis
func GetOrderBook(symbol string) OrderBook {
    // Try Redis first
    cached, err := redisClient.Get(ctx, "orderbook:"+symbol).Result()
    if err == nil {
        return parseOrderBook(cached)
    }
    
    // Query DB
    orderBook := queryOrderBookFromDB(symbol)
    
    // Cache в Redis на 1 секунду
    redisClient.Set(ctx, "orderbook:"+symbol, serialize(orderBook), 1*time.Second)
    
    return orderBook
}
```

---

## Partitioning

### Partitioning by Range (PostgreSQL)
```sql
-- Головна таблиця
CREATE TABLE stock_orders (
    id SERIAL,
    user_id INTEGER NOT NULL,
    symbol TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    ...
) PARTITION BY RANGE (created_at);

-- Партиції по місяцях
CREATE TABLE stock_orders_2024_01 PARTITION OF stock_orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE stock_orders_2024_02 PARTITION OF stock_orders
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- Query автоматично використає потрібну партицію
SELECT * FROM stock_orders 
WHERE created_at >= '2024-01-15' AND created_at < '2024-01-20';
```

### Partitioning by List
```sql
CREATE TABLE stock_orders (
    ...
) PARTITION BY LIST (order_type);

CREATE TABLE stock_orders_bid PARTITION OF stock_orders
    FOR VALUES IN ('bid');

CREATE TABLE stock_orders_ask PARTITION OF stock_orders
    FOR VALUES IN ('ask');
```

---

## Best Practices

### 1. Використовуйте транзакції для consistency
```sql
BEGIN TRANSACTION;

-- Match bid and ask
UPDATE stock_orders SET status = 'filled' WHERE id = 1;
UPDATE stock_orders SET status = 'filled' WHERE id = 2;
INSERT INTO trades (bid_id, ask_id, price, quantity) VALUES (1, 2, 150.00, 10);

COMMIT;
```

### 2. Batch INSERT замість single inserts
```sql
-- ❌ Погано: N запитів
INSERT INTO messages (content, client_id) VALUES ('msg1', 'client1');
INSERT INTO messages (content, client_id) VALUES ('msg2', 'client2');
INSERT INTO messages (content, client_id) VALUES ('msg3', 'client3');

-- ✅ Добре: 1 запит
INSERT INTO messages (content, client_id) VALUES 
    ('msg1', 'client1'),
    ('msg2', 'client2'),
    ('msg3', 'client3');
```

### 3. UPSERT для conditional insert/update
```sql
-- PostgreSQL
INSERT INTO users (username, created_at) 
VALUES ('john', NOW())
ON CONFLICT (username) DO UPDATE 
SET created_at = EXCLUDED.created_at;

-- MySQL
INSERT INTO users (username, created_at) 
VALUES ('john', NOW())
ON DUPLICATE KEY UPDATE created_at = VALUES(created_at);

-- SQLite
INSERT OR REPLACE INTO users (username, created_at) 
VALUES ('john', datetime('now'));
```

### 4. Prepared Statements (prevent SQL injection + performance)
```go
// Добре: prepared statement з параметрами
stmt, _ := db.Prepare("SELECT * FROM users WHERE username = ?")
defer stmt.Close()
row := stmt.QueryRow("john_trader")
```

### 5. Read Replicas для scaling reads
```go
// Write до master
_, err := masterDB.Exec("INSERT INTO stock_orders ...")

// Read з replica
rows, err := replicaDB.Query("SELECT * FROM stock_orders WHERE symbol = ?", "AAPL")
```

### 6. Архівування старих даних
```sql
-- Переміщувати старі filled orders в архів
INSERT INTO stock_orders_archive 
SELECT * FROM stock_orders 
WHERE status = 'filled' 
  AND created_at < datetime('now', '-90 days');

DELETE FROM stock_orders 
WHERE status = 'filled' 
  AND created_at < datetime('now', '-90 days');
```

### 7. Моніторинг повільних запитів
```sql
-- PostgreSQL: enable slow query log
ALTER SYSTEM SET log_min_duration_statement = 1000; -- 1 second
SELECT pg_reload_conf();

-- MySQL: slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;
```

### 8. Connection timeout
```go
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(10 * time.Minute)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

row := db.QueryRowContext(ctx, "SELECT ...")
```

---

## Performance Checklist

### Before deploying
- [ ] Всі часто використовувані WHERE колонки мають індекси
- [ ] EXPLAIN ANALYZE для критичних запитів
- [ ] Немає SELECT * в production коді
- [ ] Використовується connection pooling
- [ ] Prepared statements для параметризованих запитів
- [ ] Транзакції для критичних операцій
- [ ] Batch inserts замість single inserts

### Regular maintenance
- [ ] Моніторинг повільних запитів
- [ ] VACUUM/OPTIMIZE таблиць
- [ ] Перевірка невикористовуваних індексів
- [ ] Архівування старих даних
- [ ] Backup бази даних

### Monitoring metrics
- [ ] Query latency (p50, p95, p99)
- [ ] Connection pool usage
- [ ] Cache hit ratio
- [ ] Index usage statistics
- [ ] Database size growth
- [ ] Slow query count

---

## Common Anti-Patterns

### ❌ 1. Fetching all data then filtering in code
```go
// Погано
orders := GetAllOrders()  // SELECT * FROM stock_orders
filtered := filterBySymbol(orders, "AAPL")
```

### ✅ Filter in database
```go
// Добре
orders := GetOrdersBySymbol("AAPL")  // SELECT * WHERE symbol = ?
```

### ❌ 2. Multiple queries for related data (N+1)
```go
// Погано
for _, order := range orders {
    user := GetUser(order.UserID)  // N queries!
}
```

### ✅ Join or batch fetch
```sql
-- Добре
SELECT o.*, u.username 
FROM stock_orders o 
JOIN users u ON o.user_id = u.id;
```

### ❌ 3. Updating in loop
```go
// Погано: N updates
for _, orderID := range orderIDs {
    db.Exec("UPDATE stock_orders SET status = ? WHERE id = ?", "filled", orderID)
}
```

### ✅ Bulk update
```sql
-- Добре: 1 update
UPDATE stock_orders 
SET status = 'filled' 
WHERE id IN (1, 2, 3, 4, 5);
```

---

**Створено:** 2026-02-05  
**Версія:** 1.0  
**Тема:** SQL Optimization and Performance Tuning
