# PostgreSQL Setup для Stock Hub

Інструкція по налаштуванню PostgreSQL замість SQLite.

---

## 📦 Крок 1: Встановити PostgreSQL драйвер

```bash
cd stock_hub
go get github.com/lib/pq
go mod tidy
```

---

## 🐘 Крок 2: Встановити PostgreSQL (якщо ще не встановлено)

### macOS (Homebrew):
```bash
brew install postgresql@15
brew services start postgresql@15
```

### Ubuntu/Debian:
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### Docker (рекомендовано для розробки):
```bash
docker run --name stock_hub_postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=stock_hub \
  -p 5432:5432 \
  -d postgres:15
```

---

## 🗄️ Крок 3: Створити базу даних

### Варіант A: Через psql
```bash
# Підключитись до PostgreSQL
psql -U postgres

# Створити базу даних
CREATE DATABASE stock_hub;

# Створити користувача (опціонально)
CREATE USER stock_hub_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE stock_hub TO stock_hub_user;

# Вийти
\q
```

### Варіант B: Через Docker
База даних вже створена (`POSTGRES_DB=stock_hub`)

---

## ⚙️ Крок 4: Налаштувати змінні оточення

Створити файл `.env` в корені проекту:

```bash
# Database Type
DB_TYPE=postgres

# PostgreSQL Connection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=stock_hub
POSTGRES_SSLMODE=disable

# Server
SERVER_PORT=8082
```

---

## 🚀 Крок 5: Запустити міграції

### Автоматично (при старті сервера):
```bash
go run cmd/server/main.go
```

### Або вручну:
```bash
go run cmd/migrate/main.go -command=up
```

---

## 📊 Крок 6: Перевірити підключення

```bash
# Підключитись до БД
psql -U postgres -d stock_hub

# Перевірити таблиці
\dt

# Повинні побачити:
#  public | messages      | table | postgres
#  public | stock_orders  | table | postgres
#  public | users         | table | postgres
```

---

## 🔄 Перемикання між SQLite та PostgreSQL

### Використовувати SQLite (за замовчуванням):
```bash
export DB_TYPE=sqlite
go run cmd/server/main.go
```

### Використовувати PostgreSQL:
```bash
export DB_TYPE=postgres
go run cmd/server/main.go
```

---

## 🛠️ Корисні команди

### PostgreSQL CLI
```bash
# Підключитись
psql -U postgres -d stock_hub

# Список баз даних
\l

# Список таблиць
\dt

# Опис таблиці
\d stock_orders

# Вийти
\q
```

### SQL Запити
```sql
-- Список користувачів
SELECT * FROM users;

-- Відкриті ордери
SELECT * FROM stock_orders WHERE status = 'open';

-- Order book для AAPL
SELECT 
    order_type, 
    price, 
    SUM(quantity) as total_quantity
FROM stock_orders
WHERE symbol = 'AAPL' AND status = 'open'
GROUP BY order_type, price
ORDER BY 
    CASE WHEN order_type = 'bid' THEN -price ELSE price END;
```

---

## 🐳 Docker Compose (рекомендовано)

Створити `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: stock_hub_postgres
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: stock_hub
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

Запустити:
```bash
docker-compose up -d
```

Зупинити:
```bash
docker-compose down
```

---

## 📈 Порівняння: SQLite vs PostgreSQL

| Параметр | SQLite | PostgreSQL |
|----------|--------|------------|
| **Тип** | Embedded (файлова) | Server-based |
| **Продуктивність** | Відмінно для < 1GB | Відмінно для будь-якого розміру |
| **Конкурентність** | Single writer | Multiple writers |
| **Транзакції** | ✅ ACID | ✅ ACID |
| **Складність** | Нульова | Потребує сервер |
| **Use Case** | Розробка, mobile | Production, enterprise |

---

## 🔍 Troubleshooting

### Помилка: "connection refused"
```bash
# Перевірити чи запущено PostgreSQL
brew services list  # macOS
sudo systemctl status postgresql  # Linux
docker ps  # Docker
```

### Помилка: "database does not exist"
```bash
psql -U postgres -c "CREATE DATABASE stock_hub;"
```

### Помилка: "password authentication failed"
```bash
# Перевірити credentials в .env
cat .env | grep POSTGRES
```

---

## 🎯 Production рекомендації

1. **Connection Pooling:**
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

2. **SSL Mode:**
```bash
POSTGRES_SSLMODE=require  # Для production!
```

3. **Backup:**
```bash
pg_dump -U postgres stock_hub > backup.sql
```

4. **Restore:**
```bash
psql -U postgres stock_hub < backup.sql
```

---

**Створено:** 2026-02-05  
**Версія:** 1.0  
**Database:** PostgreSQL 15+
