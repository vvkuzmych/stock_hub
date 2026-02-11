# SQL Documentation - Stock Hub Project

Повна колекція SQL документації для проекту Stock Hub.

---

## 📚 Доступні документи

### 1. 📖 [SQL_TOP_70_QUERIES.md](SQL_TOP_70_QUERIES.md)
**Базова SQL довідка + Огляд Баз Даних**

🗄️ **Нове:** Огляд 15 типів баз даних:
- **Реляційні:** PostgreSQL, MySQL, SQLite, SQL Server, Oracle
- **NoSQL:** MongoDB, Redis, Cassandra, Elasticsearch, DynamoDB
- **NewSQL:** CockroachDB, Google Spanner
- **Time-Series:** InfluxDB, TimescaleDB
- **Graph:** Neo4j
- **Рекомендації:** Як вибрати БД за типом даних, навантаженням та масштабом

📝 **70 найважливіших SQL запитів:**
- ✅ SELECT, INSERT, UPDATE, DELETE
- ✅ JOIN операції (INNER, LEFT, RIGHT, FULL OUTER)
- ✅ Агрегатні функції (COUNT, SUM, AVG, MIN, MAX)
- ✅ Підзапити та CTE
- ✅ Window Functions (ROW_NUMBER, RANK, LAG, LEAD)
- ✅ DDL команди (CREATE, ALTER, DROP)
- ✅ Індекси та Foreign Keys
- ✅ Транзакції та Views

**Для кого:** Початківці та середній рівень  
**Час на вивчення:** 2-3 години  
**Версія:** 1.1 (оновлено 2026-02-05)

---

### 2. 🎯 [SQL_STOCK_HUB_QUERIES.md](SQL_STOCK_HUB_QUERIES.md)
**Практичні запити для Stock Hub проекту**

SQL запити специфічно для таблиць Stock Hub:
- 👤 **Users**: створення, пошук, статистика
- 💬 **Messages**: логування, пошук, архівування
- 📈 **Stock Orders**: BID/ASK ордери, Order Book, matching
- 📊 **Аналітика**: топ символів, обсяг торгів, активність користувачів
- 📄 **Звітність**: денні/місячні звіти, топ трейдери
- 🔧 **Оптимізація**: індекси, Views, транзакції

**Для кого:** Розробники Stock Hub  
**База даних:** SQLite (stock_hub.db)

---

### 3. 🛠️ [SQL_ADMIN_MONITORING.md](SQL_ADMIN_MONITORING.md)
**Адміністрування та моніторинг баз даних**

Запити для адміністрування, моніторингу та troubleshooting:
- 🐘 **PostgreSQL Admin**: версія, розмір БД, активні з'єднання, найповільніші запити
- 🐬 **MySQL/MariaDB Admin**: процеси, InnoDB status, реплікація
- 🗄️ **SQLite Admin**: розмір, цілісність, оптимізація (VACUUM)
- 📊 **Performance Monitoring**: EXPLAIN, locks, cache hit ratio
- 💾 **Backup та Recovery**: pg_dump, mysqldump, SQLite backup
- 🔐 **Security**: користувачі, права доступу, permissions

**Для кого:** DevOps, Database Administrators  
**Підтримка:** PostgreSQL, MySQL, SQLite

---

### 4. ⚡ [SQL_OPTIMIZATION_GUIDE.md](SQL_OPTIMIZATION_GUIDE.md)
**Оптимізація продуктивності SQL запитів**

Детальний гайд по оптимізації та performance tuning:
- 🚫 **Anti-patterns**: що НЕ робити
- 🎯 **Індексування**: single, composite, partial, covering indexes
- 🔍 **Query Optimization**: SELECT *, WHERE optimization, JOIN vs IN
- 🔄 **N+1 Problem**: eager loading, batch loading
- 🏊 **Connection Pooling**: налаштування для різних БД
- 💾 **Caching**: application-level, materialized views, Redis
- 📂 **Partitioning**: range, list partitioning
- ✅ **Best Practices**: транзакції, prepared statements, архівування

**Для кого:** Середній та advanced рівень  
**Результат:** 10-100x швидші запити

---

### 5. 🎓 [SQL_ADVANCED_EXAMPLES.md](SQL_ADVANCED_EXAMPLES.md)
**Складні SQL запити та real-world приклади**

Просунуті техніки та практичні сценарії:
- 🔗 **Advanced JOINs**: self-join, lateral join, multiple joins
- 📊 **Window Functions**: LAG/LEAD, moving average, percentile
- 🔁 **Recursive CTEs**: ієрархії, генерація sequences
- 📈 **Pivot Tables**: динамічна аналітика
- 📉 **Running Totals**: накопичувальні суми, середні
- 🎯 **Real-World Scenarios**: 
  - Order matching algorithm
  - VWAP calculation
  - Market depth (Order Book)
  - Cohort analysis
  - Outlier detection

**Для кого:** Advanced level  
**Приклади:** Всі з проекту Stock Hub

---

## 🚀 Quick Start

### Для початківців
1. Почни з [SQL_TOP_70_QUERIES.md](SQL_TOP_70_QUERIES.md)
2. Практикуй основні запити на [SQL Fiddle](http://sqlfiddle.com/)
3. Перейди до [SQL_STOCK_HUB_QUERIES.md](SQL_STOCK_HUB_QUERIES.md)

### Для розробників Stock Hub
1. Ознайомся зі структурою БД: `migrations/`
2. Вивчи [SQL_STOCK_HUB_QUERIES.md](SQL_STOCK_HUB_QUERIES.md)
3. Використовуй [SQL_OPTIMIZATION_GUIDE.md](SQL_OPTIMIZATION_GUIDE.md) для оптимізації

### Для DevOps/DBA
1. Налаштуй моніторинг з [SQL_ADMIN_MONITORING.md](SQL_ADMIN_MONITORING.md)
2. Оптимізуй запити з [SQL_OPTIMIZATION_GUIDE.md](SQL_OPTIMIZATION_GUIDE.md)
3. Налаштуй backup стратегію

---

## 📊 Структура Stock Hub БД

### Таблиці

#### `users`
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### `messages`
```sql
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    client_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### `stock_orders`
```sql
CREATE TABLE stock_orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    username TEXT NOT NULL,
    symbol TEXT NOT NULL,
    order_type TEXT CHECK(order_type IN ('bid', 'ask')),
    price REAL NOT NULL,
    quantity INTEGER NOT NULL,
    status TEXT DEFAULT 'open' CHECK(status IN ('open', 'filled', 'cancelled')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## 🎯 Практичні сценарії

### Сценарій 1: Новий користувач реєструється
```sql
-- 1. Створити користувача
INSERT INTO users (username) VALUES ('new_trader');

-- 2. Перевірити чи створено
SELECT * FROM users WHERE username = 'new_trader';
```

### Сценарій 2: Розмістити BID ордер
```sql
-- 1. Створити BID ордер
INSERT INTO stock_orders (user_id, username, symbol, order_type, price, quantity)
VALUES (1, 'john_trader', 'AAPL', 'bid', 150.00, 10);

-- 2. Перевірити Order Book
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open'
ORDER BY 
    CASE WHEN order_type = 'bid' THEN -price ELSE price END;
```

### Сценарій 3: Match bid та ask
```sql
BEGIN TRANSACTION;

-- Знайти matching orders
SELECT 
    b.id as bid_id, a.id as ask_id, b.price as bid_price, a.price as ask_price
FROM stock_orders b
JOIN stock_orders a ON b.symbol = a.symbol
WHERE b.order_type = 'bid' 
  AND a.order_type = 'ask'
  AND b.price >= a.price
  AND b.status = 'open' 
  AND a.status = 'open'
LIMIT 1;

-- Виконати matching
UPDATE stock_orders SET status = 'filled' WHERE id IN (bid_id, ask_id);

COMMIT;
```

### Сценарій 4: Аналітика за день
```sql
-- Статистика торгів за сьогодні
SELECT 
    symbol,
    COUNT(*) as trades,
    SUM(quantity) as volume,
    AVG(price) as avg_price,
    MIN(price) as low,
    MAX(price) as high
FROM stock_orders
WHERE DATE(created_at) = DATE('now')
  AND status = 'filled'
GROUP BY symbol
ORDER BY volume DESC;
```

---

## 🛠️ Корисні інструменти

### SQLite CLI
```bash
# Відкрити базу даних
sqlite3 stock_hub.db

# Показати таблиці
.tables

# Показати схему
.schema stock_orders

# Експорт в CSV
.mode csv
.output orders.csv
SELECT * FROM stock_orders;
.output stdout

# Включити foreign keys
PRAGMA foreign_keys = ON;

# Перевірити цілісність
PRAGMA integrity_check;
```

### GUI Tools
- **DB Browser for SQLite**: https://sqlitebrowser.org/
- **DBeaver**: https://dbeaver.io/
- **TablePlus**: https://tableplus.com/

---

## 📖 Корисні ресурси

### Онлайн SQL Практика
- [SQLBolt](https://sqlbolt.com/) - Interactive SQL tutorial
- [SQL Fiddle](http://sqlfiddle.com/) - Test SQL online
- [LeetCode Database](https://leetcode.com/problemset/database/) - SQL challenges

### Документація
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)

### Книги
- "SQL Performance Explained" by Markus Winand
- "High Performance MySQL" by Baron Schwartz
- "PostgreSQL: Up and Running" by Regina Obe

---

## 🎓 Learning Path

### Beginner (0-3 місяці)
- [ ] Вивчити основи SELECT, WHERE, ORDER BY
- [ ] Опанувати JOIN операції
- [ ] Зрозуміти агрегатні функції
- [ ] Навчитися використовувати підзапити

**Ресурси:** SQL_TOP_70_QUERIES.md

### Intermediate (3-6 місяців)
- [ ] Window Functions
- [ ] Common Table Expressions (CTE)
- [ ] Індексування
- [ ] Transactions

**Ресурси:** SQL_STOCK_HUB_QUERIES.md, SQL_OPTIMIZATION_GUIDE.md

### Advanced (6+ місяців)
- [ ] Query optimization
- [ ] Recursive queries
- [ ] Partitioning
- [ ] Replication та sharding

**Ресурси:** SQL_ADVANCED_EXAMPLES.md, SQL_ADMIN_MONITORING.md

---

## 🚨 Common Mistakes

### 1. Відсутність індексів
```sql
-- ❌ Погано: Full table scan
SELECT * FROM stock_orders WHERE symbol = 'AAPL';

-- ✅ Добре: Створити індекс
CREATE INDEX idx_symbol ON stock_orders(symbol);
```

### 2. SELECT *
```sql
-- ❌ Погано
SELECT * FROM stock_orders;

-- ✅ Добре
SELECT id, symbol, price, quantity FROM stock_orders;
```

### 3. N+1 запитів
```go
// ❌ Погано
for _, user := range users {
    orders := GetOrdersByUserID(user.ID)  // N queries!
}

// ✅ Добре
orders := GetAllOrdersWithUsers()  // 1 query with JOIN
```

### 4. Відсутність LIMIT
```sql
-- ❌ Погано: може повернути мільйони рядків
SELECT * FROM stock_orders;

-- ✅ Добре
SELECT * FROM stock_orders LIMIT 100;
```

---

## 📞 Support

Маєш питання? Перевір:
1. Відповідний SQL документ вище
2. SQLite документацію: https://www.sqlite.org/
3. Stack Overflow: https://stackoverflow.com/questions/tagged/sql

---

## 🔄 Changelog

- **2026-02-05 v1.1**: Оновлено SQL_TOP_70_QUERIES.md
  - ✨ Додано детальний огляд 15 типів баз даних
  - 📊 Порівняльні таблиці по вибору БД
  - 🎯 Рекомендації за типом даних, навантаженням та масштабом
  
- **2026-02-05 v1.0**: Initial documentation created
  - SQL_TOP_70_QUERIES.md (70 базових SQL запитів)
  - SQL_STOCK_HUB_QUERIES.md (практичні запити для Stock Hub)
  - SQL_ADMIN_MONITORING.md (адміністрування та моніторинг)
  - SQL_OPTIMIZATION_GUIDE.md (оптимізація продуктивності)
  - SQL_ADVANCED_EXAMPLES.md (складні запити та сценарії)

---

**Останнє оновлення:** 2026-02-05  
**Версія:** 1.0  
**Проект:** Stock Hub  
**База даних:** SQLite
