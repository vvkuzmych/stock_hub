# Changelog

## [2.0.0] - 2026-02-06

### 🔄 BREAKING CHANGES
- **Перехід на PostgreSQL** - проект тепер використовує тільки PostgreSQL
- Видалено підтримку SQLite повністю
- Змінено структуру команд у Makefile

### ✨ Додано
- Повна підтримка PostgreSQL як єдиної БД
- Оновлено всі сервіси для роботи з PostgreSQL
- Спрощено структуру проекту (один тип БД замість двох)
- Оновлено документацію для PostgreSQL-only підходу

### 🗑️ Видалено
- Підтримка SQLite (`modernc.org/sqlite`)
- Старі файли конфігурації для SQLite
- Подвійна логіка для SQLite/PostgreSQL в сервісах
- `main_postgres.go` (тепер просто `main.go`)
- `message_service_postgres.go` (тепер просто `message_service.go`)
- `config_postgres.go` (тепер просто `config.go`)
- `migrations_postgres/` (тепер просто `migrations/`)

### 🔧 Змінено
- Спрощено всі сервіси - використовується тільки PostgreSQL синтаксис
- Оновлено Makefile - видалено команди для SQLite
- `make run` тепер завжди запускає PostgreSQL версію
- `make migrate-up` завжди працює з PostgreSQL
- Оновлено `go.mod` - залишено тільки PostgreSQL драйвер
- Змінено логіку placeholder'ів в міграторі - завжди `$1, $2...`

### 🐛 Виправлено
- Видалено конфлікти між SQLite та PostgreSQL синтаксисом
- Виправлено всі `LastInsertId()` на `RETURNING id` (PostgreSQL стандарт)
- Видалено непотрібні перевірки типу БД

### 📚 Документація
- Оновлено README.md під PostgreSQL-only
- Створено CHANGELOG.md
- Оновлено всю SQL документацію

---

## [1.0.0] - 2026-02-05

### ✨ Початковий реліз
- Базова функціональність WebSocket сервера
- Підтримка SQLite та PostgreSQL
- Система міграцій
- Real-time торгівля акціями
- Чат між користувачами
- Створення та скасування заявок
