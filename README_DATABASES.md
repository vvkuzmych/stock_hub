# Stock Hub - Database Options

Stock Hub тепер підтримує **2 типи баз даних**: SQLite та PostgreSQL.

---

## 🗄️ Підтримувані бази даних

### 1. **SQLite** (за замовчуванням)
✅ **Переваги:**
- Нульова конфігурація
- Один файл = вся БД
- Ідеально для розробки та прототипів
- Не потребує сервера

❌ **Обмеження:**
- Single writer (один процес пише одночасно)
- Обмежена масштабованість
- Не підходить для high-load production

📁 **Файл БД:** `stock_hub.db`

---

### 2. **PostgreSQL** (нова підтримка)
✅ **Переваги:**
- Production-ready
- Multiple concurrent writers
- Відмінна масштабованість
- Підтримка складних запитів
- ACID транзакції
- Активна спільнота

❌ **Вимоги:**
- Потребує PostgreSQL сервер
- Більш складне налаштування

🐘 **Версія:** PostgreSQL 15+

---

## 🚀 Quick Start

### Варіант A: SQLite (простий)
```bash
# Просто запусти!
go run cmd/server/main.go
```

### Варіант B: PostgreSQL (з Docker)
```bash
# 1. Запустити PostgreSQL
make docker-up

# 2. Встановити драйвер
make install-deps

# 3. Створити .env
cp .env.example .env

# 4. Запустити сервер
make run-postgres
```

---

## 📊 Порівняння

| Параметр | SQLite | PostgreSQL |
|----------|--------|------------|
| **Setup** | ✅ Миттєво | ⚠️ Потребує налаштування |
| **Performance (reads)** | ✅ Відмінно | ✅ Відмінно |
| **Performance (writes)** | ⚠️ Single writer | ✅ Concurrent |
| **Розмір БД** | < 1GB рекомендовано | Без обмежень |
| **Production** | ❌ Не рекомендовано | ✅ Так |
| **Backup** | Копія файлу | pg_dump |
| **Monitoring** | Базовий | Повний |

---

## 📁 Структура файлів

```
stock_hub/
├── migrations/              # SQLite міграції
├── migrations_postgres/     # PostgreSQL міграції
├── cmd/
│   ├── server/
│   │   ├── main.go         # SQLite версія
│   │   └── main_postgres.go # PostgreSQL версія
│   └── migrate/
│       └── main.go         # Міграції
├── internal/
│   ├── config/
│   │   ├── config.go       # SQLite config
│   │   └── config_postgres.go # PostgreSQL config
│   └── service/
│       ├── message_service.go    # SQLite
│       └── message_service_postgres.go # PostgreSQL
├── docker-compose.yml      # PostgreSQL Docker
├── .env.example           # Приклад конфігурації
├── Makefile              # Команди
├── POSTGRES_SETUP.md     # Детальна інструкція
└── QUICK_START_POSTGRES.md # Швидкий старт
```

---

## ⚙️ Конфігурація

### SQLite (за замовчуванням)
```bash
# Не потребує конфігурації!
# Або встанови шлях до файлу:
export DATABASE_PATH=my_database.db
```

### PostgreSQL
```bash
# Створи .env файл:
DB_TYPE=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=stock_hub
POSTGRES_SSLMODE=disable
```

---

## 🔄 Міграції

### SQLite міграції
```
migrations/
├── 001_create_messages_table.up.sql
├── 002_create_users_table.up.sql
└── 003_create_stock_orders_table.up.sql
```

### PostgreSQL міграції
```
migrations_postgres/
├── 001_create_messages_table.up.sql    # SERIAL замість AUTOINCREMENT
├── 002_create_users_table.up.sql
└── 003_create_stock_orders_table.up.sql # DECIMAL замість REAL
```

**Основні відмінності:**
- SQLite: `INTEGER PRIMARY KEY AUTOINCREMENT`
- PostgreSQL: `SERIAL PRIMARY KEY`
- SQLite: `DATETIME`
- PostgreSQL: `TIMESTAMP`
- SQLite: `REAL`
- PostgreSQL: `DECIMAL(10, 2)`

---

## 🎯 Рекомендації

### Використовуй SQLite для:
- ✅ Локальної розробки
- ✅ Прототипів та MVP
- ✅ Mobile/Desktop додатків
- ✅ Embedded систем
- ✅ Testing

### Використовуй PostgreSQL для:
- ✅ Production environment
- ✅ High-load додатків
- ✅ Multiple concurrent users
- ✅ Enterprise систем
- ✅ Cloud deployments

---

## 🛠️ Makefile команди

```bash
make help              # Показати всі команди
make run               # SQLite версія
make run-postgres      # PostgreSQL версія
make docker-up         # Запустити PostgreSQL
make docker-down       # Зупинити PostgreSQL
make setup-postgres    # Повне налаштування PostgreSQL
make migrate-up        # Застосувати міграції
make psql              # Підключитись до PostgreSQL
```

---

## 📚 Документація

- **SQLite документація:** https://www.sqlite.org/docs.html
- **PostgreSQL документація:** https://www.postgresql.org/docs/
- **Детальна інструкція PostgreSQL:** [POSTGRES_SETUP.md](docs/POSTGRES_SETUP.md)
- **Швидкий старт:** [QUICK_START_POSTGRES.md](docs/QUICK_START_POSTGRES.md)
- **SQL запити:** [SQL_STOCK_HUB_QUERIES.md](docs/SQL_STOCK_HUB_QUERIES.md)

---

## ❓ FAQ

### Як перемкнутись з SQLite на PostgreSQL?
```bash
# 1. Запусти PostgreSQL
make docker-up

# 2. Створи .env
cp .env.example .env

# 3. Запусти PostgreSQL версію
make run-postgres
```

### Чи можна використовувати обидві бази одночасно?
Ні, вибери одну:
- `make run` для SQLite
- `make run-postgres` для PostgreSQL

### Як експортувати дані з SQLite в PostgreSQL?
```bash
# 1. Export з SQLite
sqlite3 stock_hub.db .dump > export.sql

# 2. Адаптувати для PostgreSQL (замінити синтаксис)
# 3. Import в PostgreSQL
psql -U postgres -d stock_hub < export.sql
```

### Де зберігаються дані?
- **SQLite:** `stock_hub.db` файл в корені проекту
- **PostgreSQL:** В PostgreSQL сервері (Docker volume або локально)

---

**Створено:** 2026-02-05  
**Версія:** 1.0  
**Підтримка:** SQLite + PostgreSQL
