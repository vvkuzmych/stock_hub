# SQL Адміністративні та Моніторингові Запити

Запити для адміністрування, моніторингу та оптимізації баз даних.

---

## 📖 Зміст

1. [PostgreSQL Admin](#postgresql-admin)
2. [MySQL/MariaDB Admin](#mysqlmariadb-admin)
3. [SQLite Admin](#sqlite-admin)
4. [Моніторинг Performance](#моніторинг-performance)
5. [Backup та Recovery](#backup-та-recovery)
6. [Security та Permissions](#security-та-permissions)

---

## PostgreSQL Admin

### Перевірити версію
```sql
SELECT version();
```

### Список всіх баз даних
```sql
SELECT datname, pg_size_pretty(pg_database_size(datname)) as size
FROM pg_database
ORDER BY pg_database_size(datname) DESC;
```

### Список всіх таблиць
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Активні з'єднання
```sql
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    query,
    query_start
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY query_start;
```

### Завершити конкретне з'єднання
```sql
SELECT pg_terminate_backend(12345);  -- pid
```

### Завершити всі з'єднання до БД
```sql
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'target_database'
  AND pid <> pg_backend_pid();
```

### Найповільніші запити
```sql
SELECT 
    calls,
    total_exec_time,
    mean_exec_time,
    query
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Статистика по індексах
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan ASC;  -- Індекси, що рідко використовуються
```

### Невикористовувані індекси
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;
```

### Cache hit ratio
```sql
SELECT 
    sum(heap_blks_read) as heap_read,
    sum(heap_blks_hit)  as heap_hit,
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) as ratio
FROM pg_statio_user_tables;
```

### Bloat в таблицях
```sql
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size,
    n_dead_tup,
    n_live_tup,
    ROUND(n_dead_tup * 100.0 / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_tuple_percent
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
```

### Vacuum та Analyze
```sql
-- Manual vacuum
VACUUM VERBOSE ANALYZE table_name;

-- Full vacuum (locks table)
VACUUM FULL table_name;

-- Analyze only
ANALYZE table_name;
```

### Перевірити replication lag
```sql
SELECT 
    client_addr,
    state,
    sync_state,
    pg_wal_lsn_diff(pg_current_wal_lsn(), replay_lsn) as lag_bytes
FROM pg_stat_replication;
```

---

## MySQL/MariaDB Admin

### Перевірити версію
```sql
SELECT VERSION();
```

### Список баз даних з розміром
```sql
SELECT 
    table_schema as 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) as 'Size (MB)'
FROM information_schema.tables
GROUP BY table_schema
ORDER BY SUM(data_length + index_length) DESC;
```

### Список таблиць з розміром
```sql
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) as 'Size (MB)',
    table_rows,
    ROUND(((data_length) / 1024 / 1024), 2) as 'Data Size (MB)',
    ROUND(((index_length) / 1024 / 1024), 2) as 'Index Size (MB)'
FROM information_schema.tables
WHERE table_schema = 'your_database'
ORDER BY (data_length + index_length) DESC;
```

### Активні процеси
```sql
SHOW FULL PROCESSLIST;
```

### Завершити процес
```sql
KILL 12345;  -- process id
```

### Статистика по запитах (MySQL 8.0+)
```sql
SELECT 
    digest_text as query,
    count_star as exec_count,
    avg_timer_wait / 1000000000000 as avg_time_sec,
    sum_timer_wait / 1000000000000 as total_time_sec
FROM performance_schema.events_statements_summary_by_digest
ORDER BY sum_timer_wait DESC
LIMIT 10;
```

### Статистика по індексах
```sql
SELECT 
    table_schema,
    table_name,
    index_name,
    cardinality
FROM information_schema.statistics
WHERE table_schema = 'your_database'
ORDER BY cardinality DESC;
```

### Невикористовувані індекси
```sql
SELECT 
    object_schema,
    object_name,
    index_name
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE index_name IS NOT NULL
  AND count_star = 0
  AND object_schema != 'mysql'
ORDER BY object_schema, object_name;
```

### InnoDB статус
```sql
SHOW ENGINE INNODB STATUS\G
```

### Перевірити реплікацію
```sql
SHOW SLAVE STATUS\G
```

### Buffer pool usage
```sql
SELECT 
    (data_size / 1024 / 1024) as data_mb,
    (pages_data * page_size / 1024 / 1024) as pages_mb
FROM information_schema.innodb_buffer_pool_stats;
```

### Оптимізувати таблицю
```sql
OPTIMIZE TABLE table_name;
```

### Аналізувати таблицю
```sql
ANALYZE TABLE table_name;
```

---

## SQLite Admin

### Перевірити версію
```sql
SELECT sqlite_version();
```

### Розмір бази даних
```sql
SELECT 
    page_count * page_size / 1024 / 1024 as size_mb 
FROM pragma_page_count(), pragma_page_size();
```

### Список таблиць
```sql
SELECT name, type 
FROM sqlite_master 
WHERE type = 'table'
ORDER BY name;
```

### Інформація про таблицю
```sql
PRAGMA table_info(table_name);
```

### Список індексів
```sql
SELECT name, tbl_name, sql 
FROM sqlite_master 
WHERE type = 'index'
ORDER BY tbl_name, name;
```

### Статистика таблиць
```sql
SELECT 
    name,
    (SELECT COUNT(*) FROM pragma_table_info(name)) as column_count
FROM sqlite_master 
WHERE type = 'table';
```

### Перевірити цілісність
```sql
PRAGMA integrity_check;
```

### Швидка перевірка цілісності
```sql
PRAGMA quick_check;
```

### Оптимізувати базу (VACUUM)
```sql
VACUUM;
```

### Аналізувати для оптимізатора запитів
```sql
ANALYZE;
```

### Foreign keys статус
```sql
PRAGMA foreign_keys;
```

### Увімкнути foreign keys
```sql
PRAGMA foreign_keys = ON;
```

### Journal mode
```sql
PRAGMA journal_mode;
```

### Змінити на WAL mode (краща продуктивність)
```sql
PRAGMA journal_mode = WAL;
```

### Auto-vacuum статус
```sql
PRAGMA auto_vacuum;
```

### Compile options
```sql
PRAGMA compile_options;
```

---

## Моніторинг Performance

### Explain запиту (PostgreSQL)
```sql
EXPLAIN (ANALYZE, BUFFERS, VERBOSE)
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open';
```

### Explain запиту (MySQL)
```sql
EXPLAIN FORMAT=JSON
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open';
```

### Explain запиту (SQLite)
```sql
EXPLAIN QUERY PLAN
SELECT * FROM stock_orders 
WHERE symbol = 'AAPL' AND status = 'open';
```

### Час виконання запиту (PostgreSQL)
```sql
\timing on
SELECT COUNT(*) FROM stock_orders;
```

### Lock monitoring (PostgreSQL)
```sql
SELECT 
    pg_class.relname,
    pg_locks.locktype,
    pg_locks.mode,
    pg_stat_activity.state,
    pg_stat_activity.query
FROM pg_locks
JOIN pg_class ON pg_locks.relation = pg_class.oid
JOIN pg_stat_activity ON pg_locks.pid = pg_stat_activity.pid
WHERE pg_locks.granted = false;
```

### Deadlock detection (MySQL)
```sql
SHOW ENGINE INNODB STATUS\G
-- Look for "LATEST DETECTED DEADLOCK" section
```

### Transaction isolation level
```sql
-- PostgreSQL
SHOW transaction_isolation;

-- MySQL
SELECT @@transaction_isolation;
```

---

## Backup та Recovery

### PostgreSQL Backup
```bash
# Повний dump бази
pg_dump -U username -d database_name > backup.sql

# Compressed backup
pg_dump -U username -d database_name | gzip > backup.sql.gz

# Custom format (recommended)
pg_dump -U username -Fc database_name > backup.dump

# Відновлення
psql -U username -d database_name < backup.sql
pg_restore -U username -d database_name backup.dump
```

### MySQL Backup
```bash
# Повний dump
mysqldump -u username -p database_name > backup.sql

# Compressed
mysqldump -u username -p database_name | gzip > backup.sql.gz

# Конкретна таблиця
mysqldump -u username -p database_name table_name > table_backup.sql

# Відновлення
mysql -u username -p database_name < backup.sql
```

### SQLite Backup
```bash
# Командна лінія
sqlite3 stock_hub.db ".backup backup.db"

# Або copy файлу
cp stock_hub.db stock_hub_backup.db
```

### Backup через SQL (SQLite)
```sql
-- Attach backup database
ATTACH DATABASE 'backup.db' AS backup;

-- Copy data
CREATE TABLE backup.users AS SELECT * FROM main.users;
CREATE TABLE backup.messages AS SELECT * FROM main.messages;
CREATE TABLE backup.stock_orders AS SELECT * FROM main.stock_orders;

DETACH DATABASE backup;
```

---

## Security та Permissions

### Створити користувача (PostgreSQL)
```sql
CREATE USER trader WITH PASSWORD 'secure_password';
```

### Надати права
```sql
-- PostgreSQL
GRANT SELECT, INSERT, UPDATE ON stock_orders TO trader;
GRANT ALL PRIVILEGES ON DATABASE stock_hub TO admin_user;

-- MySQL
GRANT SELECT, INSERT, UPDATE ON stock_hub.stock_orders TO 'trader'@'localhost';
GRANT ALL PRIVILEGES ON stock_hub.* TO 'admin'@'localhost';
FLUSH PRIVILEGES;
```

### Відібрати права
```sql
-- PostgreSQL
REVOKE INSERT ON stock_orders FROM trader;

-- MySQL
REVOKE INSERT ON stock_hub.stock_orders FROM 'trader'@'localhost';
```

### Список користувачів (PostgreSQL)
```sql
SELECT usename, usesuper, usecreatedb 
FROM pg_user;
```

### Список користувачів (MySQL)
```sql
SELECT user, host FROM mysql.user;
```

### Перевірити права користувача
```sql
-- PostgreSQL
\du username

-- MySQL
SHOW GRANTS FOR 'username'@'localhost';
```

### Змінити пароль
```sql
-- PostgreSQL
ALTER USER username WITH PASSWORD 'new_password';

-- MySQL
ALTER USER 'username'@'localhost' IDENTIFIED BY 'new_password';
```

---

## Корисні налаштування

### PostgreSQL config
```sql
-- Показати всі налаштування
SHOW ALL;

-- Конкретне налаштування
SHOW max_connections;
SHOW shared_buffers;
SHOW work_mem;

-- Змінити (session level)
SET work_mem = '256MB';
```

### MySQL config
```sql
-- Показати всі змінні
SHOW VARIABLES;

-- Конкретна змінна
SHOW VARIABLES LIKE 'max_connections';

-- Змінити
SET GLOBAL max_connections = 200;
```

### SQLite pragmas
```sql
-- Performance settings
PRAGMA cache_size = 10000;
PRAGMA temp_store = MEMORY;
PRAGMA synchronous = NORMAL;
PRAGMA journal_mode = WAL;

-- Security
PRAGMA foreign_keys = ON;
```

---

## Maintenance Scripts

### Очистка старих даних (приклад)
```sql
-- PostgreSQL/MySQL
DELETE FROM logs 
WHERE created_at < NOW() - INTERVAL '30 days';

-- SQLite
DELETE FROM messages 
WHERE created_at < datetime('now', '-30 days');
```

### Регулярна оптимізація
```sql
-- PostgreSQL (автоматично через autovacuum, але можна manual)
VACUUM ANALYZE;

-- MySQL
OPTIMIZE TABLE stock_orders;

-- SQLite
VACUUM;
ANALYZE;
```

### Архівування даних
```sql
-- Створити архівну таблицю
CREATE TABLE stock_orders_archive AS 
SELECT * FROM stock_orders 
WHERE status = 'filled' 
  AND created_at < datetime('now', '-1 year');

-- Видалити з основної
DELETE FROM stock_orders 
WHERE status = 'filled' 
  AND created_at < datetime('now', '-1 year');
```

---

## Health Check Checklist

### Daily
- [ ] Перевірити активні з'єднання
- [ ] Моніторити повільні запити
- [ ] Перевірити розмір бази даних
- [ ] Перевірити помилки в логах

### Weekly
- [ ] Перевірити невикористовувані індекси
- [ ] Аналізувати найповільніші запити
- [ ] Перевірити bloat в таблицях
- [ ] Backup бази даних

### Monthly
- [ ] Повна оптимізація (VACUUM/OPTIMIZE)
- [ ] Очистка старих даних
- [ ] Перевірка реплікації (якщо є)
- [ ] Аудит користувачів та прав доступу

---

**Створено:** 2026-02-05  
**Версія:** 1.0  
**Підтримувані БД:** PostgreSQL, MySQL/MariaDB, SQLite
