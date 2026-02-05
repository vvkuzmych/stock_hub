# Stock Hub - Практичні SQL Запити

Колекція SQL запитів специфічно для проекту Stock Hub.

---

## 📖 Зміст

1. [Запити для Users](#запити-для-users)
2. [Запити для Messages](#запити-для-messages)
3. [Запити для Stock Orders](#запити-для-stock-orders)
4. [Аналітичні запити](#аналітичні-запити)
5. [Запити для звітності](#запити-для-звітності)
6. [Оптимізація та моніторинг](#оптимізація-та-моніторинг)

---

## Запити для Users

### Створити нового користувача
```sql
INSERT INTO users (username) 
VALUES ('john_trader');
```

### Отримати всіх користувачів
```sql
SELECT id, username, created_at 
FROM users 
ORDER BY created_at DESC;
```

### Знайти користувача за username
```sql
SELECT * FROM users 
WHERE username = 'john_trader';
```

### Перевірити чи існує користувач
```sql
SELECT EXISTS(
    SELECT 1 FROM users 
    WHERE username = 'john_trader'
) as user_exists;
```

### Отримати кількість користувачів
```sql
SELECT COUNT(*) as total_users FROM users;
```

### Отримати нових користувачів за останній день
```sql
SELECT * FROM users 
WHERE created_at >= datetime('now', '-1 day')
ORDER BY created_at DESC;
```

### Отримати користувачів з найбільшою активністю
```sql
SELECT 
    u.id,
    u.username,
    COUNT(so.id) as total_orders,
    SUM(CASE WHEN so.status = 'filled' THEN 1 ELSE 0 END) as filled_orders
FROM users u
LEFT JOIN stock_orders so ON u.id = so.user_id
GROUP BY u.id, u.username
ORDER BY total_orders DESC
LIMIT 10;
```

---

## Запити для Messages

### Створити повідомлення
```sql
INSERT INTO messages (content, client_id) 
VALUES ('User connected', 'client_123');
```

### Отримати всі повідомлення
```sql
SELECT * FROM messages 
ORDER BY created_at DESC 
LIMIT 100;
```

### Отримати повідомлення конкретного клієнта
```sql
SELECT * FROM messages 
WHERE client_id = 'client_123'
ORDER BY created_at DESC;
```

### Підрахувати повідомлення за останню годину
```sql
SELECT COUNT(*) as messages_last_hour 
FROM messages 
WHERE created_at >= datetime('now', '-1 hour');
```

### Отримати активність по клієнтах
```sql
SELECT 
    client_id,
    COUNT(*) as message_count,
    MIN(created_at) as first_message,
    MAX(created_at) as last_message
FROM messages
GROUP BY client_id
ORDER BY message_count DESC;
```

### Очистити старі повідомлення (> 30 днів)
```sql
DELETE FROM messages 
WHERE created_at < datetime('now', '-30 days');
```

### Пошук повідомлень за змістом
```sql
SELECT * FROM messages 
WHERE content LIKE '%error%'
ORDER BY created_at DESC;
```

---

## Запити для Stock Orders

### Створити BID ордер
```sql
INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity, status) 
VALUES (1, 'john_trader', 'AAPL', 'bid', 150.50, 10, 'open');
```

### Створити ASK ордер
```sql
INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity, status) 
VALUES (2, 'jane_investor', 'AAPL', 'ask', 151.00, 5, 'open');
```

### Отримати всі відкриті ордери
```sql
SELECT * FROM stock_orders 
WHERE status = 'open'
ORDER BY created_at DESC;
```

### Отримати ордери конкретного користувача
```sql
SELECT * FROM stock_orders 
WHERE user_id = 1
ORDER BY created_at DESC;
```

### Отримати ордери по символу
```sql
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL'
ORDER BY created_at DESC;
```

### Отримати відкриті BID ордери для символу
```sql
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' 
  AND order_type = 'bid' 
  AND status = 'open'
ORDER BY price DESC, created_at ASC;
```

### Отримати відкриті ASK ордери для символу
```sql
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' 
  AND order_type = 'ask' 
  AND status = 'open'
ORDER BY price ASC, created_at ASC;
```

### Знайти найкращу BID ціну
```sql
SELECT MAX(price) as best_bid_price 
FROM stock_orders 
WHERE symbol = 'AAPL' 
  AND order_type = 'bid' 
  AND status = 'open';
```

### Знайти найкращу ASK ціну
```sql
SELECT MIN(price) as best_ask_price 
FROM stock_orders 
WHERE symbol = 'AAPL' 
  AND order_type = 'ask' 
  AND status = 'open';
```

### Отримати Order Book (BID + ASK)
```sql
SELECT 
    order_type,
    price,
    SUM(quantity) as total_quantity,
    COUNT(*) as order_count
FROM stock_orders
WHERE symbol = 'AAPL' AND status = 'open'
GROUP BY order_type, price
ORDER BY 
    CASE WHEN order_type = 'bid' THEN -price ELSE price END;
```

### Оновити статус ордера
```sql
UPDATE stock_orders 
SET status = 'filled' 
WHERE id = 1;
```

### Скасувати ордер
```sql
UPDATE stock_orders 
SET status = 'cancelled' 
WHERE id = 1;
```

### Підрахувати ордери за статусом
```sql
SELECT 
    status,
    COUNT(*) as count
FROM stock_orders
GROUP BY status;
```

---

## Аналітичні запити

### Топ найактивніших символів
```sql
SELECT 
    symbol,
    COUNT(*) as total_orders,
    SUM(CASE WHEN order_type = 'bid' THEN 1 ELSE 0 END) as bid_count,
    SUM(CASE WHEN order_type = 'ask' THEN 1 ELSE 0 END) as ask_count,
    SUM(CASE WHEN status = 'filled' THEN 1 ELSE 0 END) as filled_count
FROM stock_orders
GROUP BY symbol
ORDER BY total_orders DESC
LIMIT 10;
```

### Середня ціна за символом
```sql
SELECT 
    symbol,
    order_type,
    AVG(price) as avg_price,
    MIN(price) as min_price,
    MAX(price) as max_price
FROM stock_orders
WHERE status = 'filled'
GROUP BY symbol, order_type;
```

### Обсяг торгів за символом
```sql
SELECT 
    symbol,
    SUM(quantity) as total_volume,
    SUM(price * quantity) as total_value
FROM stock_orders
WHERE status = 'filled'
GROUP BY symbol
ORDER BY total_value DESC;
```

### Співвідношення BID/ASK
```sql
SELECT 
    symbol,
    SUM(CASE WHEN order_type = 'bid' THEN quantity ELSE 0 END) as bid_volume,
    SUM(CASE WHEN order_type = 'ask' THEN quantity ELSE 0 END) as ask_volume,
    CAST(SUM(CASE WHEN order_type = 'bid' THEN quantity ELSE 0 END) AS REAL) / 
    NULLIF(SUM(CASE WHEN order_type = 'ask' THEN quantity ELSE 0 END), 0) as bid_ask_ratio
FROM stock_orders
WHERE status = 'open'
GROUP BY symbol;
```

### Активність користувачів по днях
```sql
SELECT 
    DATE(created_at) as trade_date,
    COUNT(DISTINCT user_id) as active_users,
    COUNT(*) as total_orders,
    SUM(CASE WHEN status = 'filled' THEN 1 ELSE 0 END) as filled_orders
FROM stock_orders
GROUP BY DATE(created_at)
ORDER BY trade_date DESC;
```

### Найбільші ордери
```sql
SELECT 
    so.*,
    (so.price * so.quantity) as order_value
FROM stock_orders so
WHERE status = 'filled'
ORDER BY order_value DESC
LIMIT 10;
```

### Час життя ордера (від створення до виконання)
```sql
SELECT 
    id,
    symbol,
    order_type,
    price,
    quantity,
    created_at,
    -- Note: You'd need updated_at column for this
    -- julianday(updated_at) - julianday(created_at) as days_to_fill
    'N/A' as days_to_fill
FROM stock_orders
WHERE status = 'filled'
ORDER BY created_at DESC
LIMIT 10;
```

---

## Запити для звітності

### Денний звіт торгів
```sql
SELECT 
    DATE(created_at) as trade_date,
    symbol,
    COUNT(*) as total_orders,
    SUM(CASE WHEN status = 'filled' THEN 1 ELSE 0 END) as filled_orders,
    SUM(quantity) as total_quantity,
    AVG(price) as avg_price,
    MIN(price) as min_price,
    MAX(price) as max_price
FROM stock_orders
WHERE DATE(created_at) = DATE('now')
GROUP BY DATE(created_at), symbol
ORDER BY symbol;
```

### Місячний звіт по користувачах
```sql
SELECT 
    u.username,
    COUNT(so.id) as total_orders,
    SUM(CASE WHEN so.status = 'filled' THEN 1 ELSE 0 END) as filled_orders,
    SUM(CASE WHEN so.status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_orders,
    SUM(so.quantity) as total_volume
FROM users u
LEFT JOIN stock_orders so ON u.id = so.user_id 
    AND strftime('%Y-%m', so.created_at) = strftime('%Y-%m', 'now')
GROUP BY u.id, u.username
HAVING COUNT(so.id) > 0
ORDER BY total_orders DESC;
```

### Звіт по статусах ордерів за період
```sql
SELECT 
    status,
    COUNT(*) as count,
    SUM(quantity) as total_quantity,
    SUM(price * quantity) as total_value
FROM stock_orders
WHERE created_at >= datetime('now', '-7 days')
GROUP BY status;
```

### Top користувачів за прибутковістю (припустимо)
```sql
SELECT 
    u.username,
    COUNT(so.id) as filled_orders,
    SUM(so.quantity) as total_volume,
    AVG(so.price) as avg_price
FROM users u
INNER JOIN stock_orders so ON u.id = so.user_id
WHERE so.status = 'filled'
GROUP BY u.id, u.username
ORDER BY filled_orders DESC
LIMIT 10;
```

---

## Оптимізація та моніторинг

### Перевірити розмір бази даних (SQLite)
```sql
SELECT page_count * page_size as size_bytes 
FROM pragma_page_count(), pragma_page_size();
```

### Статистика по таблицях
```sql
SELECT 
    'users' as table_name, COUNT(*) as row_count FROM users
UNION ALL
SELECT 
    'messages' as table_name, COUNT(*) as row_count FROM messages
UNION ALL
SELECT 
    'stock_orders' as table_name, COUNT(*) as row_count FROM stock_orders;
```

### Перевірити індекси
```sql
SELECT name, tbl_name, sql 
FROM sqlite_master 
WHERE type = 'index';
```

### Аналіз запиту (EXPLAIN QUERY PLAN)
```sql
EXPLAIN QUERY PLAN
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open';
```

### Перевірити foreign keys
```sql
PRAGMA foreign_keys;
```

### Увімкнути foreign keys (якщо вимкнені)
```sql
PRAGMA foreign_keys = ON;
```

### Оптимізувати базу даних
```sql
VACUUM;
```

### Аналізувати для оптимізації запитів
```sql
ANALYZE;
```

### Перевірити цілісність бази даних
```sql
PRAGMA integrity_check;
```

---

## Корисні Views

### View для активних ордерів
```sql
CREATE VIEW IF NOT EXISTS active_orders AS
SELECT 
    so.*,
    u.username as user_name,
    (so.price * so.quantity) as order_value
FROM stock_orders so
JOIN users u ON so.user_id = u.id
WHERE so.status = 'open';
```

### View для Order Book
```sql
CREATE VIEW IF NOT EXISTS order_book AS
SELECT 
    symbol,
    order_type,
    price,
    SUM(quantity) as total_quantity,
    COUNT(*) as order_count
FROM stock_orders
WHERE status = 'open'
GROUP BY symbol, order_type, price;
```

### View для статистики користувачів
```sql
CREATE VIEW IF NOT EXISTS user_statistics AS
SELECT 
    u.id,
    u.username,
    u.created_at,
    COUNT(so.id) as total_orders,
    SUM(CASE WHEN so.status = 'filled' THEN 1 ELSE 0 END) as filled_orders,
    SUM(CASE WHEN so.status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_orders,
    SUM(CASE WHEN so.status = 'open' THEN 1 ELSE 0 END) as open_orders
FROM users u
LEFT JOIN stock_orders so ON u.id = so.user_id
GROUP BY u.id, u.username, u.created_at;
```

### Використання View
```sql
SELECT * FROM active_orders WHERE symbol = 'AAPL';
SELECT * FROM order_book ORDER BY symbol, order_type, price;
SELECT * FROM user_statistics ORDER BY total_orders DESC;
```

---

## Транзакції

### Створити ордер з перевіркою користувача
```sql
BEGIN TRANSACTION;

-- Перевірити чи існує користувач
SELECT id FROM users WHERE username = 'john_trader';

-- Якщо існує, створити ордер
INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity)
VALUES (1, 'john_trader', 'AAPL', 'bid', 150.00, 10);

COMMIT;
```

### Виконати matching bid/ask
```sql
BEGIN TRANSACTION;

-- Знайти matching orders
-- Mark bid as filled
UPDATE stock_orders SET status = 'filled' WHERE id = 1;

-- Mark ask as filled
UPDATE stock_orders SET status = 'filled' WHERE id = 2;

COMMIT;
```

### Скасувати всі відкриті ордери користувача
```sql
BEGIN TRANSACTION;

UPDATE stock_orders 
SET status = 'cancelled' 
WHERE user_id = 1 AND status = 'open';

COMMIT;
```

---

## Міграції даних

### Додати колонку updated_at до stock_orders
```sql
ALTER TABLE stock_orders 
ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;
```

### Створити тригер для автоматичного оновлення updated_at
```sql
CREATE TRIGGER IF NOT EXISTS update_stock_orders_timestamp 
AFTER UPDATE ON stock_orders
BEGIN
    UPDATE stock_orders 
    SET updated_at = CURRENT_TIMESTAMP 
    WHERE id = NEW.id;
END;
```

### Backup даних в тимчасову таблицю
```sql
CREATE TABLE stock_orders_backup AS 
SELECT * FROM stock_orders 
WHERE status = 'filled' 
  AND created_at < datetime('now', '-30 days');
```

---

## Performance Tips

### 1. Використовуйте індекси для часто використовуваних фільтрів
```sql
CREATE INDEX IF NOT EXISTS idx_symbol_status 
ON stock_orders(symbol, status);
```

### 2. Composite index для order book queries
```sql
CREATE INDEX IF NOT EXISTS idx_symbol_type_status_price 
ON stock_orders(symbol, order_type, status, price);
```

### 3. Для пошуку повідомлень
```sql
CREATE INDEX IF NOT EXISTS idx_messages_content 
ON messages(content);
```

### 4. Регулярна очистка старих даних
```sql
-- Видаляти старі cancelled ордери
DELETE FROM stock_orders 
WHERE status = 'cancelled' 
  AND created_at < datetime('now', '-90 days');

-- Видаляти старі повідомлення
DELETE FROM messages 
WHERE created_at < datetime('now', '-30 days');
```

---

**Створено:** 2026-02-05  
**Проект:** Stock Hub  
**База даних:** SQLite
