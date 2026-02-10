# SQL Advanced Examples - Складні Запити

Колекція складних та корисних SQL запитів з поясненнями.

---

## 📖 Зміст

1. [Advanced JOINs](#advanced-joins)
2. [Window Functions](#window-functions)
3. [Common Table Expressions (CTE)](#common-table-expressions-cte)
4. [Recursive Queries](#recursive-queries)
5. [Pivot Tables](#pivot-tables)
6. [Running Totals](#running-totals)
7. [Gap Detection](#gap-detection)
8. [Real-World Scenarios](#real-world-scenarios)

---

## Advanced JOINs

### Self-Join: Знайти користувачів з однаковими іменами
```sql
SELECT 
    u1.id as user1_id,
    u1.username,
    u2.id as user2_id
FROM users u1
JOIN users u2 ON u1.username = u2.username AND u1.id < u2.id;
```

### Multiple JOINs: Повна інформація про ордер
```sql
SELECT 
    u.username,
    so.symbol,
    so.order_type,
    so.price,
    so.quantity,
    so.status,
    m.content as last_message
FROM stock_orders so
JOIN users u ON so.user_id = u.id
LEFT JOIN messages m ON m.client_id = CAST(u.id AS TEXT)
WHERE so.symbol = 'AAPL'
ORDER BY so.created_at DESC;
```

### Cross JOIN: Генерація всіх комбінацій
```sql
-- Приклад: всі можливі пари користувачів для matching
SELECT 
    u1.username as buyer,
    u2.username as seller
FROM users u1
CROSS JOIN users u2
WHERE u1.id != u2.id;
```

### JOIN з підзапитом
```sql
SELECT 
    u.username,
    recent_orders.order_count,
    recent_orders.total_volume
FROM users u
JOIN (
    SELECT 
        user_id,
        COUNT(*) as order_count,
        SUM(quantity) as total_volume
    FROM stock_orders
    WHERE created_at >= datetime('now', '-7 days')
    GROUP BY user_id
) recent_orders ON u.id = recent_orders.user_id
WHERE recent_orders.order_count > 10;
```

### LATERAL JOIN (PostgreSQL) - для correlated subqueries
```sql
SELECT 
    u.username,
    latest.symbol,
    latest.price,
    latest.created_at
FROM users u
CROSS JOIN LATERAL (
    SELECT symbol, price, created_at
    FROM stock_orders
    WHERE user_id = u.id
    ORDER BY created_at DESC
    LIMIT 3
) latest;
```

---

## Window Functions

### ROW_NUMBER: Нумерація ордерів по користувачу
```sql
SELECT 
    username,
    symbol,
    price,
    created_at,
    ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at DESC) as order_rank
FROM stock_orders
WHERE status = 'filled';
```

### RANK vs DENSE_RANK
```sql
SELECT 
    symbol,
    price,
    quantity,
    -- RANK: 1, 2, 2, 4 (пропускає номери)
    RANK() OVER (ORDER BY price DESC) as rank,
    -- DENSE_RANK: 1, 2, 2, 3 (без пропусків)
    DENSE_RANK() OVER (ORDER BY price DESC) as dense_rank
FROM stock_orders
WHERE symbol = 'AAPL' AND status = 'filled';
```

### LAG та LEAD: Порівняння з попереднім/наступним
```sql
SELECT 
    symbol,
    price,
    created_at,
    LAG(price) OVER (PARTITION BY symbol ORDER BY created_at) as prev_price,
    LEAD(price) OVER (PARTITION BY symbol ORDER BY created_at) as next_price,
    price - LAG(price) OVER (PARTITION BY symbol ORDER BY created_at) as price_change
FROM stock_orders
WHERE symbol = 'AAPL'
ORDER BY created_at;
```

### FIRST_VALUE та LAST_VALUE
```sql
SELECT 
    symbol,
    price,
    created_at,
    FIRST_VALUE(price) OVER (
        PARTITION BY symbol 
        ORDER BY created_at 
        ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
    ) as first_price,
    LAST_VALUE(price) OVER (
        PARTITION BY symbol 
        ORDER BY created_at 
        ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
    ) as last_price
FROM stock_orders
WHERE status = 'filled';
```

### Moving Average (ковзне середнє)
```sql
SELECT 
    symbol,
    price,
    created_at,
    AVG(price) OVER (
        PARTITION BY symbol 
        ORDER BY created_at 
        ROWS BETWEEN 4 PRECEDING AND CURRENT ROW
    ) as moving_avg_5
FROM stock_orders
WHERE symbol = 'AAPL'
ORDER BY created_at;
```

### Percentile (PERCENT_RANK)
```sql
SELECT 
    username,
    total_orders,
    PERCENT_RANK() OVER (ORDER BY total_orders) as percentile
FROM (
    SELECT 
        u.username,
        COUNT(so.id) as total_orders
    FROM users u
    LEFT JOIN stock_orders so ON u.id = so.user_id
    GROUP BY u.id, u.username
) user_stats;
```

---

## Common Table Expressions (CTE)

### Простий CTE
```sql
WITH active_traders AS (
    SELECT user_id, COUNT(*) as order_count
    FROM stock_orders
    WHERE created_at >= datetime('now', '-30 days')
    GROUP BY user_id
    HAVING COUNT(*) > 10
)
SELECT 
    u.username,
    at.order_count
FROM active_traders at
JOIN users u ON at.user_id = u.id
ORDER BY at.order_count DESC;
```

### Multiple CTEs
```sql
WITH 
bid_orders AS (
    SELECT symbol, AVG(price) as avg_bid_price
    FROM stock_orders
    WHERE order_type = 'bid' AND status = 'open'
    GROUP BY symbol
),
ask_orders AS (
    SELECT symbol, AVG(price) as avg_ask_price
    FROM stock_orders
    WHERE order_type = 'ask' AND status = 'open'
    GROUP BY symbol
)
SELECT 
    b.symbol,
    b.avg_bid_price,
    a.avg_ask_price,
    (a.avg_ask_price - b.avg_bid_price) as spread
FROM bid_orders b
JOIN ask_orders a ON b.symbol = a.symbol
ORDER BY spread DESC;
```

### Nested CTEs
```sql
WITH order_stats AS (
    SELECT 
        user_id,
        COUNT(*) as total_orders,
        SUM(CASE WHEN status = 'filled' THEN 1 ELSE 0 END) as filled_orders
    FROM stock_orders
    GROUP BY user_id
),
top_traders AS (
    SELECT *
    FROM order_stats
    WHERE total_orders > 50
)
SELECT 
    u.username,
    tt.total_orders,
    tt.filled_orders,
    ROUND(tt.filled_orders * 100.0 / tt.total_orders, 2) as fill_rate
FROM top_traders tt
JOIN users u ON tt.user_id = u.id
ORDER BY fill_rate DESC;
```

---

## Recursive Queries

### Recursive CTE: Generate series (dates)
```sql
-- PostgreSQL
WITH RECURSIVE dates AS (
    SELECT DATE '2024-01-01' as date
    UNION ALL
    SELECT date + INTERVAL '1 day'
    FROM dates
    WHERE date < DATE '2024-01-31'
)
SELECT 
    d.date,
    COUNT(so.id) as orders_count
FROM dates d
LEFT JOIN stock_orders so ON DATE(so.created_at) = d.date
GROUP BY d.date
ORDER BY d.date;
```

### Hierarchical data (якби була така структура)
```sql
-- Припустимо, є таблиця categories з parent_id
WITH RECURSIVE category_tree AS (
    -- Anchor: root categories
    SELECT id, name, parent_id, 1 as level
    FROM categories
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Recursive: child categories
    SELECT c.id, c.name, c.parent_id, ct.level + 1
    FROM categories c
    JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT * FROM category_tree ORDER BY level, name;
```

### Number sequence
```sql
WITH RECURSIVE numbers AS (
    SELECT 1 as n
    UNION ALL
    SELECT n + 1
    FROM numbers
    WHERE n < 100
)
SELECT * FROM numbers;
```

---

## Pivot Tables

### Динамічна pivot table (order types по символах)
```sql
SELECT 
    symbol,
    SUM(CASE WHEN order_type = 'bid' THEN quantity ELSE 0 END) as bid_volume,
    SUM(CASE WHEN order_type = 'ask' THEN quantity ELSE 0 END) as ask_volume,
    SUM(CASE WHEN order_type = 'bid' THEN 1 ELSE 0 END) as bid_count,
    SUM(CASE WHEN order_type = 'ask' THEN 1 ELSE 0 END) as ask_count
FROM stock_orders
WHERE status = 'open'
GROUP BY symbol;
```

### Статуси по користувачах
```sql
SELECT 
    username,
    SUM(CASE WHEN status = 'open' THEN 1 ELSE 0 END) as open_orders,
    SUM(CASE WHEN status = 'filled' THEN 1 ELSE 0 END) as filled_orders,
    SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) as cancelled_orders
FROM stock_orders so
JOIN users u ON so.user_id = u.id
GROUP BY u.id, u.username;
```

### Pivot по датах (orders per day)
```sql
SELECT 
    symbol,
    SUM(CASE WHEN DATE(created_at) = DATE('now') THEN 1 ELSE 0 END) as today,
    SUM(CASE WHEN DATE(created_at) = DATE('now', '-1 day') THEN 1 ELSE 0 END) as yesterday,
    SUM(CASE WHEN DATE(created_at) >= DATE('now', '-7 days') THEN 1 ELSE 0 END) as last_7_days
FROM stock_orders
GROUP BY symbol;
```

---

## Running Totals

### Running sum (накопичувальна сума)
```sql
SELECT 
    symbol,
    price,
    quantity,
    created_at,
    SUM(quantity) OVER (
        PARTITION BY symbol 
        ORDER BY created_at 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as cumulative_volume
FROM stock_orders
WHERE symbol = 'AAPL' AND status = 'filled'
ORDER BY created_at;
```

### Running average
```sql
SELECT 
    symbol,
    price,
    created_at,
    AVG(price) OVER (
        PARTITION BY symbol 
        ORDER BY created_at 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as cumulative_avg_price
FROM stock_orders
WHERE symbol = 'AAPL'
ORDER BY created_at;
```

### Running count
```sql
SELECT 
    username,
    symbol,
    created_at,
    ROW_NUMBER() OVER (
        PARTITION BY user_id 
        ORDER BY created_at
    ) as order_number
FROM stock_orders so
JOIN users u ON so.user_id = u.id
ORDER BY user_id, created_at;
```

---

## Gap Detection

### Знайти пропуски в ID sequence
```sql
WITH numbered AS (
    SELECT 
        id,
        ROW_NUMBER() OVER (ORDER BY id) as rn
    FROM stock_orders
)
SELECT 
    id + 1 as gap_start,
    (SELECT MIN(n2.id) - 1 
     FROM numbered n2 
     WHERE n2.id > n1.id) as gap_end
FROM numbered n1
WHERE NOT EXISTS (
    SELECT 1 FROM numbered n2 
    WHERE n2.id = n1.id + 1
)
AND id < (SELECT MAX(id) FROM stock_orders);
```

### Знайти дні без активності
```sql
WITH RECURSIVE all_dates AS (
    SELECT DATE('2024-01-01') as date
    UNION ALL
    SELECT date(date, '+1 day')
    FROM all_dates
    WHERE date < DATE('2024-01-31')
)
SELECT ad.date
FROM all_dates ad
LEFT JOIN stock_orders so ON DATE(so.created_at) = ad.date
WHERE so.id IS NULL;
```

---

## Real-World Scenarios

### 1. Order Matching Algorithm (простий приклад)
```sql
-- Знайти matching bid/ask orders
SELECT 
    b.id as bid_id,
    a.id as ask_id,
    b.symbol,
    b.price as bid_price,
    a.price as ask_price,
    MIN(b.quantity, a.quantity) as match_quantity
FROM stock_orders b
JOIN stock_orders a ON b.symbol = a.symbol
WHERE b.order_type = 'bid' 
  AND a.order_type = 'ask'
  AND b.status = 'open'
  AND a.status = 'open'
  AND b.price >= a.price  -- Bid price >= Ask price = match!
ORDER BY 
    b.symbol,
    a.price ASC,  -- Найкраща ask ціна
    b.price DESC  -- Найкраща bid ціна
LIMIT 10;
```

### 2. Top Movers (найбільша зміна ціни)
```sql
WITH price_changes AS (
    SELECT 
        symbol,
        FIRST_VALUE(price) OVER (
            PARTITION BY symbol 
            ORDER BY created_at
        ) as first_price,
        LAST_VALUE(price) OVER (
            PARTITION BY symbol 
            ORDER BY created_at
            ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
        ) as last_price
    FROM stock_orders
    WHERE created_at >= datetime('now', '-1 day')
      AND status = 'filled'
)
SELECT DISTINCT
    symbol,
    first_price,
    last_price,
    (last_price - first_price) as price_change,
    ROUND(((last_price - first_price) / first_price) * 100, 2) as percent_change
FROM price_changes
ORDER BY ABS(percent_change) DESC
LIMIT 10;
```

### 3. Volume Weighted Average Price (VWAP)
```sql
SELECT 
    symbol,
    SUM(price * quantity) / SUM(quantity) as vwap,
    SUM(quantity) as total_volume
FROM stock_orders
WHERE status = 'filled'
  AND created_at >= datetime('now', '-1 day')
GROUP BY symbol
ORDER BY total_volume DESC;
```

### 4. User Trading Statistics
```sql
WITH user_stats AS (
    SELECT 
        u.username,
        COUNT(*) as total_orders,
        SUM(CASE WHEN so.status = 'filled' THEN 1 ELSE 0 END) as filled_orders,
        SUM(CASE WHEN so.order_type = 'bid' THEN 1 ELSE 0 END) as buy_orders,
        SUM(CASE WHEN so.order_type = 'ask' THEN 1 ELSE 0 END) as sell_orders,
        AVG(so.price * so.quantity) as avg_order_value,
        MAX(so.created_at) as last_trade_date
    FROM users u
    LEFT JOIN stock_orders so ON u.id = so.user_id
    GROUP BY u.id, u.username
)
SELECT 
    *,
    ROUND(filled_orders * 100.0 / NULLIF(total_orders, 0), 2) as fill_rate,
    CASE 
        WHEN buy_orders > sell_orders THEN 'Buyer'
        WHEN sell_orders > buy_orders THEN 'Seller'
        ELSE 'Balanced'
    END as trader_type
FROM user_stats
WHERE total_orders > 0
ORDER BY total_orders DESC;
```

### 5. Market Depth (Order Book)
```sql
WITH order_book AS (
    SELECT 
        symbol,
        order_type,
        price,
        SUM(quantity) as total_quantity,
        COUNT(*) as order_count
    FROM stock_orders
    WHERE status = 'open'
    GROUP BY symbol, order_type, price
)
SELECT 
    symbol,
    order_type,
    price,
    total_quantity,
    order_count,
    SUM(total_quantity) OVER (
        PARTITION BY symbol, order_type 
        ORDER BY 
            CASE WHEN order_type = 'bid' THEN -price ELSE price END
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) as cumulative_quantity
FROM order_book
WHERE symbol = 'AAPL'
ORDER BY 
    symbol,
    CASE WHEN order_type = 'bid' THEN 0 ELSE 1 END,
    CASE WHEN order_type = 'bid' THEN -price ELSE price END;
```

### 6. Daily Trading Summary
```sql
SELECT 
    DATE(created_at) as trade_date,
    symbol,
    COUNT(*) as total_trades,
    SUM(quantity) as volume,
    MIN(price) as low,
    MAX(price) as high,
    FIRST_VALUE(price) OVER (
        PARTITION BY DATE(created_at), symbol 
        ORDER BY created_at
    ) as open_price,
    LAST_VALUE(price) OVER (
        PARTITION BY DATE(created_at), symbol 
        ORDER BY created_at
        ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
    ) as close_price
FROM stock_orders
WHERE status = 'filled'
  AND created_at >= datetime('now', '-7 days')
GROUP BY DATE(created_at), symbol, id, price, created_at
ORDER BY trade_date DESC, symbol;
```

### 7. User Cohort Analysis
```sql
WITH user_cohorts AS (
    SELECT 
        u.id,
        u.username,
        DATE(u.created_at) as signup_date,
        MIN(DATE(so.created_at)) as first_trade_date
    FROM users u
    LEFT JOIN stock_orders so ON u.id = so.user_id
    GROUP BY u.id, u.username, u.created_at
)
SELECT 
    signup_date,
    COUNT(*) as users_signed_up,
    COUNT(first_trade_date) as users_who_traded,
    ROUND(COUNT(first_trade_date) * 100.0 / COUNT(*), 2) as activation_rate,
    ROUND(AVG(julianday(first_trade_date) - julianday(signup_date)), 1) as avg_days_to_first_trade
FROM user_cohorts
GROUP BY signup_date
ORDER BY signup_date DESC;
```

### 8. Outlier Detection (аномальні ордери)
```sql
WITH price_stats AS (
    SELECT 
        symbol,
        AVG(price) as avg_price,
        STDEV(price) as stddev_price  -- PostgreSQL: stddev_pop()
    FROM stock_orders
    WHERE status = 'filled'
      AND created_at >= datetime('now', '-7 days')
    GROUP BY symbol
)
SELECT 
    so.id,
    so.symbol,
    so.price,
    ps.avg_price,
    ABS(so.price - ps.avg_price) / ps.stddev_price as z_score
FROM stock_orders so
JOIN price_stats ps ON so.symbol = ps.symbol
WHERE ABS(so.price - ps.avg_price) / ps.stddev_price > 3  -- 3 standard deviations
ORDER BY z_score DESC;
```

---

## Performance Tips для складних запитів

### 1. Матеріалізуйте підзапити якщо вони використовуються кілька разів
```sql
-- Погано: підзапит виконується двічі
SELECT * FROM orders 
WHERE price > (SELECT AVG(price) FROM orders WHERE symbol = 'AAPL')
  AND price < (SELECT AVG(price) FROM orders WHERE symbol = 'AAPL') * 1.1;

-- Добре: обчислити один раз
WITH avg_price AS (
    SELECT AVG(price) as avg FROM orders WHERE symbol = 'AAPL'
)
SELECT * FROM orders, avg_price
WHERE price > avg_price.avg AND price < avg_price.avg * 1.1;
```

### 2. Обмежуйте window functions якщо можливо
```sql
-- Швидше: обмежити ROWS
AVG(price) OVER (ORDER BY created_at ROWS BETWEEN 10 PRECEDING AND CURRENT ROW)

-- Повільніше: необмежене вікно
AVG(price) OVER (ORDER BY created_at)
```

### 3. Використовуйте EXISTS замість IN для великих списків
```sql
-- Повільніше
WHERE user_id IN (SELECT user_id FROM active_users)

-- Швидше
WHERE EXISTS (SELECT 1 FROM active_users WHERE user_id = orders.user_id)
```

---

**Створено:** 2026-02-05  
**Версія:** 1.0  
**Тема:** Advanced SQL Queries and Real-World Examples
