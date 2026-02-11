# TOP 70 Найважливіших SQL Запитів

Цей документ містить найважливіші SQL запити для роботи з базами даних.

---

## 📖 Зміст

### 🗄️ [Список всіх баз даних](#список-всіх-баз-даних)

### 🗄️ [Типи Баз Даних](#типи-баз-даних)
- Реляційні бази (PostgreSQL, MySQL, SQLite, SQL Server, Oracle)
- NoSQL бази (MongoDB, Redis, Cassandra, Elasticsearch, DynamoDB)
- NewSQL (CockroachDB, Google Spanner)
- Time-Series (InfluxDB, TimescaleDB)
- Graph (Neo4j)
- Як вибрати базу даних

### 📚 Базові SELECT запити (1-5)
1. [Вибірка всіх даних з таблиці](#1-вибірка-всіх-даних-з-таблиці)
2. [Вибірка конкретних колонок](#2-вибірка-конкретних-колонок)
3. [Вибірка з умовою WHERE](#3-вибірка-з-умовою-where)
4. [Сортування ORDER BY](#4-сортування-order-by)
5. [Групування GROUP BY](#5-групування-group-by)

### 🔢 Агрегатні функції (6-10)
6. [COUNT - підрахунок записів](#6-count---підрахунок-записів)
7. [SUM - сума значень](#7-sum---сума-значень)
8. [AVG - середнє значення](#8-avg---середнє-значення)
9. [MIN - мінімальне значення](#9-min---мінімальне-значення)
10. [MAX - максимальне значення](#10-max---максимальне-значення)

### 🔗 JOIN операції (11-15)
11. [INNER JOIN](#11-inner-join)
12. [LEFT JOIN](#12-left-join)
13. [RIGHT JOIN](#13-right-join)
14. [FULL OUTER JOIN](#14-full-outer-join)
15. [CROSS JOIN](#15-cross-join)

### 🎯 Фільтрація та пошук (16-20)
16. [LIKE - пошук за шаблоном](#16-like---пошук-за-шаблоном)
17. [IN - список значень](#17-in---список-значень)
18. [BETWEEN - діапазон значень](#18-between---діапазон-значень)
19. [IS NULL - перевірка на NULL](#19-is-null---перевірка-на-null)
20. [IS NOT NULL](#20-is-not-null)

### 📊 Підзапити (Subqueries) (21-25)
21. [Підзапит в WHERE](#21-підзапит-в-where)
22. [Підзапит в SELECT](#22-підзапит-в-select)
23. [EXISTS](#23-exists)
24. [NOT EXISTS](#24-not-exists)
25. [IN з підзапитом](#25-in-з-підзапитом)

### 🛠️ DDL - створення та модифікація (26-31)
26. [CREATE DATABASE](#26-create-database)
27. [CREATE TABLE](#27-create-table)
28. [ALTER TABLE - додати колонку](#28-alter-table---додати-колонку)
29. [ALTER TABLE - змінити колонку](#29-alter-table---змінити-колонку)
30. [DROP TABLE](#30-drop-table)
31. [TRUNCATE - очистити таблицю](#31-truncate---очистити-таблицю)

### ➕ DML - маніпуляція даними (32-35)
32. [INSERT - один запис](#32-insert---один-запис)
33. [INSERT - багато записів](#33-insert---багато-записів)
34. [UPDATE - оновлення](#34-update---оновлення)
35. [DELETE - видалення](#35-delete---видалення)

### 🔑 Індекси та ключі (36-40)
36. [CREATE INDEX](#36-create-index)
37. [CREATE UNIQUE INDEX](#37-create-unique-index)
38. [DROP INDEX](#38-drop-index)
39. [PRIMARY KEY constraint](#39-primary-key-constraint)
40. [FOREIGN KEY constraint](#40-foreign-key-constraint)

### 👁️ VIEW - представлення (41-44)
41. [CREATE VIEW](#41-create-view)
42. [SELECT з VIEW](#42-select-з-view)
43. [UPDATE VIEW](#43-update-view)
44. [DROP VIEW](#44-drop-view)

### 🔄 UNION операції (45-46)
45. [UNION - об'єднання без дублікатів](#45-union---обєднання-без-дублікатів)
46. [UNION ALL - об'єднання з дублікатами](#46-union-all---обєднання-з-дублікатами)

### 📈 Аналітичні функції (47-52)
47. [ROW_NUMBER()](#47-row_number)
48. [RANK()](#48-rank)
49. [DENSE_RANK()](#49-dense_rank)
50. [LAG - попереднє значення](#50-lag---попереднє-значення)
51. [LEAD - наступне значення](#51-lead---наступне-значення)
52. [PARTITION BY](#52-partition-by)

### 🎲 Додаткові функції (53-56)
53. [CASE WHEN](#53-case-when)
54. [COALESCE - перше не-NULL значення](#54-coalesce---перше-не-null-значення)
55. [NULLIF - повертає NULL якщо рівні](#55-nullif---повертає-null-якщо-рівні)
56. [CAST - приведення типів](#56-cast---приведення-типів)

### 📅 Робота з датами (57-60)
57. [CURRENT_TIMESTAMP](#57-current_timestamp)
58. [DATE_TRUNC](#58-date_trunc)
59. [EXTRACT](#59-extract)
60. [AGE](#60-age)

### 🔐 Транзакції та безпека (61-63)
61. [BEGIN TRANSACTION](#61-begin-transaction)
62. [ROLLBACK](#62-rollback)
63. [SAVEPOINT](#63-savepoint)

### 📊 CTE (Common Table Expressions) (64-65)
64. [WITH - простий CTE](#64-with---простий-cte)
65. [Recursive CTE](#65-recursive-cte)

### 🔍 Пошук та текст (66-70)
66. [UPPER / LOWER](#66-upper--lower)
67. [CONCAT](#67-concat)
68. [SUBSTRING](#68-substring)
69. [LENGTH](#69-length)
70. [DISTINCT - унікальні значення](#70-distinct---унікальні-значення)

---

## Список всіх баз даних

```sql
-- PostgreSQL
SELECT datname FROM pg_database;

-- MySQL
SHOW DATABASES;

-- SQLite
SELECT name FROM sqlite_master WHERE type='table';

-- SQL Server
SELECT name FROM sys.databases;
```

---

## 🗄️ Типи Баз Даних

### Реляційні бази даних (SQL/RDBMS)

#### 1. **PostgreSQL**
- 🔥 **Open Source**: Повністю безкоштовна
- ⚡ **Продуктивність**: Відмінна для складних запитів
- 🎯 **Use Cases**: Enterprise додатки, аналітика, геопросторові дані
- 📊 **Особливості**: 
  - ACID compliance
  - JSON/JSONB підтримка
  - Full-text search
  - Materialized views
  - Partitioning
- 🌐 **Компанії**: Instagram, Spotify, Netflix, Apple

#### 2. **MySQL / MariaDB**
- 🔥 **Популярність**: Найпопулярніша open-source БД
- ⚡ **Швидкість**: Відмінно для read-heavy операцій
- 🎯 **Use Cases**: Web додатки, WordPress, e-commerce
- 📊 **Особливості**:
  - Проста у налаштуванні
  - Чудова реплікація
  - InnoDB storage engine
- 🌐 **Компанії**: Facebook, Twitter, YouTube, Uber
- 📌 **MariaDB**: Fork MySQL з додатковими features

#### 3. **SQLite**
- 💾 **Embedded**: Файлова БД, не потребує сервера
- 🎯 **Use Cases**: Mobile apps, embedded systems, прототипи
- 📊 **Особливості**:
  - Нульова конфігурація
  - Один файл = вся БД
  - Відмінно для < 1GB даних
- 🌐 **Використання**: Android, iOS apps, браузери

#### 4. **Microsoft SQL Server**
- 💼 **Enterprise**: Комерційна БД від Microsoft
- 🎯 **Use Cases**: Windows-based додатки, .NET
- 📊 **Особливості**:
  - T-SQL (Transact-SQL)
  - Integration з Azure
  - Business Intelligence tools
- 🌐 **Компанії**: Banks, insurance, Microsoft ecosystem

#### 5. **Oracle Database**
- 👑 **Enterprise-grade**: Найпотужніша commercial БД
- 🎯 **Use Cases**: Banking, телеком, великі корпорації
- 📊 **Особливості**:
  - RAC (Real Application Clusters)
  - Advanced security
  - Partitioning та compression
- 🌐 **Компанії**: Banks, SAP, Oracle Cloud

---

### NoSQL бази даних (Non-relational)

#### 6. **MongoDB** (Document Store)
- 📄 **Документи**: JSON-подібні документи (BSON)
- 🎯 **Use Cases**: Content management, IoT, real-time analytics
- 📊 **Особливості**:
  - Flexible schema
  - Horizontal scaling (sharding)
  - Aggregation framework
- 🌐 **Компанії**: eBay, MetLife, Adobe

#### 7. **Redis** (Key-Value Store)
- ⚡ **In-Memory**: Надшвидка БД в пам'яті
- 🎯 **Use Cases**: Caching, session store, pub/sub, queues
- 📊 **Особливості**:
  - Millisecond latency
  - Data structures (strings, lists, sets, hashes)
  - Persistence options
- 🌐 **Компанії**: Twitter, GitHub, Snapchat, StackOverflow

#### 8. **Cassandra** (Wide Column Store)
- 📊 **Distributed**: Розподілена БД без single point of failure
- 🎯 **Use Cases**: Time-series data, IoT, messaging
- 📊 **Особливості**:
  - Linear scalability
  - Multi-datacenter replication
  - High availability
- 🌐 **Компанії**: Netflix, Apple, Instagram

#### 9. **Elasticsearch** (Search Engine)
- 🔍 **Full-Text Search**: Спеціалізована БД для пошуку
- 🎯 **Use Cases**: Log analytics, site search, APM
- 📊 **Особливості**:
  - Real-time indexing
  - RESTful API
  - ELK Stack (Elasticsearch, Logstash, Kibana)
- 🌐 **Компанії**: GitHub, Netflix, LinkedIn

#### 10. **DynamoDB** (Key-Value + Document)
- ☁️ **AWS Managed**: Serverless NoSQL від Amazon
- 🎯 **Use Cases**: Gaming, mobile apps, serverless
- 📊 **Особливості**:
  - Auto-scaling
  - Single-digit millisecond latency
  - Serverless
- 🌐 **Компанії**: Amazon, Lyft, Samsung

---

### NewSQL (Hybrid)

#### 11. **CockroachDB**
- 🪳 **Distributed SQL**: PostgreSQL-compatible
- 🎯 **Use Cases**: Global apps, multi-region
- 📊 **Особливості**:
  - Auto-sharding
  - Geo-partitioning
  - ACID transactions

#### 12. **Google Spanner**
- 🌍 **Global**: Перша глобально розподілена БД
- 🎯 **Use Cases**: Global financial systems
- 📊 **Особливості**:
  - True time API
  - Multi-region ACID

---

### Time-Series Databases

#### 13. **InfluxDB**
- 📈 **Time-Series**: Оптимізована для часових даних
- 🎯 **Use Cases**: Metrics, monitoring, IoT sensors
- 📊 **Особливості**:
  - High write throughput
  - Data retention policies
  - InfluxQL (SQL-like)

#### 14. **TimescaleDB**
- ⏰ **PostgreSQL Extension**: Time-series на базі PostgreSQL
- 🎯 **Use Cases**: IoT, financial data, monitoring
- 📊 **Особливості**:
  - Full SQL support
  - Compression
  - Continuous aggregates

---

### Graph Databases

#### 15. **Neo4j**
- 🕸️ **Graph**: Для зв'язаних даних
- 🎯 **Use Cases**: Social networks, fraud detection, recommendations
- 📊 **Особливості**:
  - Cypher query language
  - Native graph storage
  - ACID transactions
- 🌐 **Компанії**: LinkedIn, eBay, Walmart

---

## 🎯 Як вибрати базу даних?

### За типом даних:

| Тип даних | Рекомендована БД |
|-----------|------------------|
| Структуровані (табличні) | PostgreSQL, MySQL |
| Документи (JSON) | MongoDB, DynamoDB |
| Key-Value (швидкий cache) | Redis, Memcached |
| Time-Series | InfluxDB, TimescaleDB |
| Graph (соціальні зв'язки) | Neo4j, Amazon Neptune |
| Full-Text Search | Elasticsearch, Algolia |

### За навантаженням:

| Навантаження | Рекомендована БД |
|--------------|------------------|
| Read-heavy | MySQL, Redis (cache) |
| Write-heavy | Cassandra, MongoDB |
| Analytical (OLAP) | PostgreSQL, ClickHouse |
| Transactional (OLTP) | PostgreSQL, MySQL |
| Real-time | Redis, DynamoDB |

### За масштабом:

| Розмір | Рекомендована БД |
|--------|------------------|
| < 1GB (mobile, prototypes) | SQLite |
| 1GB - 100GB (small-medium apps) | PostgreSQL, MySQL |
| 100GB - 1TB (large apps) | PostgreSQL (partitioning), MongoDB |
| 1TB+ (big data) | Cassandra, CockroachDB, Spanner |

---

## 📚 Базові SELECT запити

### 1 Вибірка всіх даних з таблиці
```sql
SELECT * FROM users;
```

### 2 Вибірка конкретних колонок
```sql
SELECT first_name, last_name, email FROM users;
```

### 3 Вибірка з умовою WHERE
```sql
SELECT * FROM products WHERE price > 100;
```

### 4 Сортування ORDER BY
```sql
SELECT * FROM orders ORDER BY created_at DESC;
```

### 5 Групування GROUP BY
```sql
SELECT category, COUNT(*) as count 
FROM products 
GROUP BY category;
```

---

## 🔢 Агрегатні функції

### 6 COUNT - підрахунок записів
```sql
SELECT COUNT(*) FROM users WHERE status = 'active';
```

### 7 SUM - сума значень
```sql
SELECT SUM(amount) as total_sales FROM orders;
```

### 8 AVG - середнє значення
```sql
SELECT AVG(price) as average_price FROM products;
```

### 9 MIN - мінімальне значення
```sql
SELECT MIN(price) FROM products;
```

### 10 MAX - максимальне значення
```sql
SELECT MAX(salary) FROM employees;
```

---

## 🔗 JOIN операції

### 11 INNER JOIN
```sql
SELECT u.name, o.order_id 
FROM users u 
INNER JOIN orders o ON u.id = o.user_id;
```

### 12 LEFT JOIN
```sql
SELECT u.name, o.order_id 
FROM users u 
LEFT JOIN orders o ON u.id = o.user_id;
```

### 13 RIGHT JOIN
```sql
SELECT u.name, o.order_id 
FROM users u 
RIGHT JOIN orders o ON u.id = o.user_id;
```

### 14 FULL OUTER JOIN
```sql
SELECT u.name, o.order_id 
FROM users u 
FULL OUTER JOIN orders o ON u.id = o.user_id;
```

### 15 CROSS JOIN
```sql
SELECT p.name, c.color 
FROM products p 
CROSS JOIN colors c;
```

---

## 🎯 Фільтрація та пошук

### 16 LIKE - пошук за шаблоном
```sql
SELECT * FROM users WHERE email LIKE '%@gmail.com';
```

### 17 IN - список значень
```sql
SELECT * FROM products WHERE category IN ('Electronics', 'Books', 'Toys');
```

### 18 BETWEEN - діапазон значень
```sql
SELECT * FROM orders WHERE created_at BETWEEN '2024-01-01' AND '2024-12-31';
```

### 19 IS NULL - перевірка на NULL
```sql
SELECT * FROM users WHERE phone_number IS NULL;
```

### 20 IS NOT NULL
```sql
SELECT * FROM users WHERE email IS NOT NULL;
```

---

## 📊 Підзапити (Subqueries)

### 21 Підзапит в WHERE
```sql
SELECT * FROM employees 
WHERE salary > (SELECT AVG(salary) FROM employees);
```

### 22 Підзапит в SELECT
```sql
SELECT name, (SELECT COUNT(*) FROM orders WHERE user_id = u.id) as order_count
FROM users u;
```

### 23 EXISTS
```sql
SELECT * FROM users u 
WHERE EXISTS (SELECT 1 FROM orders WHERE user_id = u.id);
```

### 24 NOT EXISTS
```sql
SELECT * FROM products p 
WHERE NOT EXISTS (SELECT 1 FROM order_items WHERE product_id = p.id);
```

### 25 IN з підзапитом
```sql
SELECT * FROM users 
WHERE id IN (SELECT user_id FROM orders WHERE total > 1000);
```

---

## 🛠️ DDL - створення та модифікація структури

### 26 CREATE DATABASE
```sql
CREATE DATABASE my_shop;
```

### 27 CREATE TABLE
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 28 ALTER TABLE - додати колонку
```sql
ALTER TABLE users ADD COLUMN phone VARCHAR(20);
```

### 29 ALTER TABLE - змінити колонку
```sql
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(150);
```

### 30 DROP TABLE
```sql
DROP TABLE IF EXISTS temp_table;
```

### 31 TRUNCATE - очистити таблицю
```sql
TRUNCATE TABLE logs;
```

---

## ➕ DML - маніпуляція даними

### 32 INSERT - один запис
```sql
INSERT INTO users (username, email) 
VALUES ('john_doe', 'john@example.com');
```

### 33 INSERT - багато записів
```sql
INSERT INTO products (name, price) 
VALUES 
    ('Laptop', 999.99),
    ('Mouse', 29.99),
    ('Keyboard', 79.99);
```

### 34 UPDATE - оновлення
```sql
UPDATE users 
SET status = 'active' 
WHERE last_login > NOW() - INTERVAL '30 days';
```

### 35 DELETE - видалення
```sql
DELETE FROM orders 
WHERE status = 'cancelled' AND created_at < NOW() - INTERVAL '1 year';
```

---

## 🔑 Індекси та ключі

### 36 CREATE INDEX
```sql
CREATE INDEX idx_users_email ON users(email);
```

### 37 CREATE UNIQUE INDEX
```sql
CREATE UNIQUE INDEX idx_users_username ON users(username);
```

### 38 DROP INDEX
```sql
DROP INDEX IF EXISTS idx_users_email;
```

### 39 PRIMARY KEY constraint
```sql
ALTER TABLE orders ADD PRIMARY KEY (id);
```

### 40 FOREIGN KEY constraint
```sql
ALTER TABLE orders 
ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) REFERENCES users(id);
```

---

## 👁️ VIEW - представлення

### 41 CREATE VIEW
```sql
CREATE VIEW active_users AS
SELECT * FROM users WHERE status = 'active';
```

### 42 SELECT з VIEW
```sql
SELECT * FROM active_users;
```

### 43 UPDATE VIEW
```sql
CREATE OR REPLACE VIEW user_stats AS
SELECT 
    u.id,
    u.username,
    COUNT(o.id) as order_count
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.id, u.username;
```

### 44 DROP VIEW
```sql
DROP VIEW IF EXISTS active_users;
```

---

## 🔄 UNION операції

### 45 UNION - об'єднання без дублікатів
```sql
SELECT email FROM customers
UNION
SELECT email FROM suppliers;
```

### 46 UNION ALL - об'єднання з дублікатами
```sql
SELECT name FROM products_2023
UNION ALL
SELECT name FROM products_2024;
```

---

## 📈 Аналітичні функції (Window Functions)

### 47 ROW_NUMBER()
```sql
SELECT 
    name,
    salary,
    ROW_NUMBER() OVER (ORDER BY salary DESC) as rank
FROM employees;
```

### 48 RANK()
```sql
SELECT 
    name,
    score,
    RANK() OVER (ORDER BY score DESC) as rank
FROM students;
```

### 49 DENSE_RANK()
```sql
SELECT 
    product_name,
    sales,
    DENSE_RANK() OVER (ORDER BY sales DESC) as rank
FROM product_sales;
```

### 50 LAG - попереднє значення
```sql
SELECT 
    date,
    revenue,
    LAG(revenue) OVER (ORDER BY date) as prev_revenue
FROM daily_sales;
```

### 51 LEAD - наступне значення
```sql
SELECT 
    date,
    revenue,
    LEAD(revenue) OVER (ORDER BY date) as next_revenue
FROM daily_sales;
```

### 52 PARTITION BY
```sql
SELECT 
    department,
    name,
    salary,
    AVG(salary) OVER (PARTITION BY department) as dept_avg
FROM employees;
```

---

## 🎲 Додаткові функції

### 53 CASE WHEN
```sql
SELECT 
    name,
    age,
    CASE 
        WHEN age < 18 THEN 'Minor'
        WHEN age < 65 THEN 'Adult'
        ELSE 'Senior'
    END as age_group
FROM users;
```

### 54 COALESCE - перше не-NULL значення
```sql
SELECT 
    name,
    COALESCE(phone, mobile, 'No contact') as contact
FROM users;
```

### 55 NULLIF - повертає NULL якщо рівні
```sql
SELECT 
    product,
    NULLIF(discount_price, regular_price) as actual_discount
FROM products;
```

### 56 CAST - приведення типів
```sql
SELECT CAST(price AS INTEGER) FROM products;
```

---

## 📅 Робота з датами

### 57 CURRENT_TIMESTAMP
```sql
SELECT CURRENT_TIMESTAMP;
```

### 58 DATE_TRUNC
```sql
SELECT DATE_TRUNC('month', created_at) as month, COUNT(*) 
FROM orders 
GROUP BY month;
```

### 59 EXTRACT
```sql
SELECT EXTRACT(YEAR FROM created_at) as year 
FROM orders;
```

### 60 AGE
```sql
SELECT name, AGE(birth_date) as age FROM users;
```

---

## 🔐 Транзакції та безпека

### 61 BEGIN TRANSACTION
```sql
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;
```

### 62 ROLLBACK
```sql
BEGIN;
DELETE FROM important_data;
ROLLBACK; -- Скасувати зміни
```

### 63 SAVEPOINT
```sql
BEGIN;
UPDATE users SET status = 'inactive';
SAVEPOINT sp1;
DELETE FROM users WHERE status = 'inactive';
ROLLBACK TO sp1; -- Повернутись до savepoint
COMMIT;
```

---

## 📊 CTE (Common Table Expressions)

### 64 WITH - простий CTE
```sql
WITH high_earners AS (
    SELECT * FROM employees WHERE salary > 100000
)
SELECT department, COUNT(*) 
FROM high_earners 
GROUP BY department;
```

### 65 Recursive CTE
```sql
WITH RECURSIVE employee_hierarchy AS (
    SELECT id, name, manager_id, 1 as level
    FROM employees WHERE manager_id IS NULL
    
    UNION ALL
    
    SELECT e.id, e.name, e.manager_id, eh.level + 1
    FROM employees e
    JOIN employee_hierarchy eh ON e.manager_id = eh.id
)
SELECT * FROM employee_hierarchy;
```

---

## 🔍 Пошук та текст

### 66 UPPER / LOWER
```sql
SELECT UPPER(name) as name_upper, LOWER(email) as email_lower 
FROM users;
```

### 67 CONCAT
```sql
SELECT CONCAT(first_name, ' ', last_name) as full_name 
FROM users;
```

### 68 SUBSTRING
```sql
SELECT SUBSTRING(email FROM 1 FOR POSITION('@' IN email) - 1) as username 
FROM users;
```

### 69 LENGTH
```sql
SELECT name, LENGTH(name) as name_length 
FROM products 
WHERE LENGTH(name) > 20;
```

### 70 DISTINCT - унікальні значення
```sql
SELECT DISTINCT category FROM products;
```

---

## 🎯 Бонусні корисні запити

### TOP N записів
```sql
SELECT * FROM products ORDER BY sales DESC LIMIT 10;
```

### Пагінація
```sql
SELECT * FROM users 
ORDER BY created_at DESC 
LIMIT 20 OFFSET 40; -- Сторінка 3, по 20 записів
```

### Перевірка існування таблиці
```sql
SELECT EXISTS (
    SELECT FROM information_schema.tables 
    WHERE table_schema = 'public' 
    AND table_name = 'users'
);
```

### Копіювання структури таблиці
```sql
CREATE TABLE users_backup (LIKE users INCLUDING ALL);
```

### Вставка з SELECT
```sql
INSERT INTO archive_orders 
SELECT * FROM orders 
WHERE created_at < NOW() - INTERVAL '1 year';
```

---

## 📝 Корисні поради

1. **Використовуйте індекси** для колонок, що часто використовуються в WHERE, JOIN, ORDER BY
2. **Уникайте SELECT *** у продакшн коді - вказуйте тільки потрібні колонки
3. **Використовуйте EXPLAIN ANALYZE** для оптимізації запитів
4. **Обов'язково** використовуйте транзакції для критичних операцій
5. **Завжди** робіть backup перед масовими UPDATE/DELETE

---

## 🔗 Джерела

- PostgreSQL Documentation
- MySQL Documentation  
- SQL Standard (ISO/IEC 9075)
- Натхнення: [ByteScout Blog](https://bytescout.com/blog/20-important-sql-queries.html)

---

**Створено:** 2026-01-27  
**Оновлено:** 2026-02-05  
**Версія:** 1.2  
**Автор:** Developer Guide for Stock Hub Project  
**Зміни v1.2:** Додано SQL запити для отримання списку баз даних (PostgreSQL, MySQL, SQLite, SQL Server, Oracle)  
**Зміни v1.1:** Додано детальний огляд 15 типів баз даних та рекомендації по вибору
