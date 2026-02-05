# Quick Start - PostgreSQL для Stock Hub

Швидкий старт за 5 хвилин! 🚀

---

## 🎯 Варіант 1: Docker (Найпростіший)

### Крок 1: Запустити PostgreSQL
```bash
cd /Users/vkuzm/GolandProjects/stock_hub
docker-compose up -d
```

### Крок 2: Встановити драйвер
```bash
go get github.com/lib/pq
go mod tidy
```

### Крок 3: Налаштувати .env
```bash
cp .env.example .env
# Залишити все як є (вже налаштовано для Docker)
```

### Крок 4: Запустити сервер
```bash
DB_TYPE=postgres go run cmd/server/main_postgres.go
```

**Готово!** Відкрий http://localhost:8082 🎉

---

## ⚡ Варіант 2: Локальний PostgreSQL (macOS)

### Крок 1: Встановити PostgreSQL
```bash
brew install postgresql@15
brew services start postgresql@15
```

### Крок 2: Створити базу даних
```bash
createdb stock_hub
```

### Крок 3: Встановити драйвер
```bash
cd /Users/vkuzm/GolandProjects/stock_hub
go get github.com/lib/pq
go mod tidy
```

### Крок 4: Налаштувати .env
```bash
cp .env.example .env
# Відредагувати якщо потрібно
```

### Крок 5: Запустити
```bash
DB_TYPE=postgres go run cmd/server/main_postgres.go
```

---

## 🔄 Перемикання між SQLite та PostgreSQL

### SQLite (оригінальна версія):
```bash
go run cmd/server/main.go
```

### PostgreSQL (нова версія):
```bash
DB_TYPE=postgres go run cmd/server/main_postgres.go
```

---

## 🧪 Тестування

### 1. Перевірити з'єднання
```bash
psql -U postgres -d stock_hub -c "SELECT 1;"
```

### 2. Перевірити таблиці
```bash
psql -U postgres -d stock_hub -c "\dt"
```

Повинні побачити:
- `messages`
- `users`
- `stock_orders`
- `schema_migrations`

### 3. Перевірити дані
```bash
psql -U postgres -d stock_hub -c "SELECT * FROM users;"
```

---

## 📊 Корисні команди

### Docker
```bash
# Запустити PostgreSQL
docker-compose up -d

# Зупинити
docker-compose down

# Логи
docker-compose logs -f postgres

# Видалити всі дані (обережно!)
docker-compose down -v
```

### PostgreSQL CLI
```bash
# Підключитись
psql -U postgres -d stock_hub

# Список таблиць
\dt

# Вийти
\q
```

### Міграції
```bash
# Застосувати міграції
go run cmd/migrate/main.go -command=up

# Відкотити останню
go run cmd/migrate/main.go -command=down
```

---

## ❓ Troubleshooting

### PostgreSQL не запущено?
```bash
# Docker
docker-compose ps

# Локально (macOS)
brew services list

# Запустити
brew services start postgresql@15
```

### База даних не існує?
```bash
createdb stock_hub
# або
psql -U postgres -c "CREATE DATABASE stock_hub;"
```

### Connection refused?
Перевір `.env`:
```bash
cat .env | grep POSTGRES
```

---

## 🎉 Перевірка роботи

1. Відкрий http://localhost:8082
2. Зареєструй користувача
3. Створи BID ордер для AAPL на $150.00
4. Створи ASK ордер для AAPL на $151.00
5. Перевір в базі:
```bash
psql -U postgres -d stock_hub -c "SELECT * FROM stock_orders;"
```

---

**Готово!** Тепер твій Stock Hub працює на PostgreSQL! 🐘✨
