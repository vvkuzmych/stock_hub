# Stock Hub - Real-time Trading Platform 📊

Stock Hub - це real-time платформа для торгівлі акціями з підтримкою WebSocket з'єднань та PostgreSQL базою даних.

## 🚀 Швидкий старт

### Попередні вимоги

- Go 1.25+
- PostgreSQL 14+
- Docker (опціонально, для легкого запуску PostgreSQL)

### Встановлення

1. **Клонуйте репозиторій**
```bash
cd /Users/vkuzm/GolandProjects/stock_hub
```

2. **Встановіть залежності**
```bash
make install-deps
```

3. **Запустіть PostgreSQL**
```bash
# Варіант 1: Docker
make docker-up

# Варіант 2: Локальна установка (потрібно встановити PostgreSQL окремо)
# Створіть базу даних вручну:
psql -U postgres -c "CREATE DATABASE stock_hub;"
```

4. **Налаштуйте середовище**
```bash
# Скопіюйте приклад .env
cp .env.example .env

# Відредагуйте .env при необхідності (за замовчуванням налаштовано для Docker)
```

5. **Застосуйте міграції**
```bash
make migrate-up
```

6. **Запустіть сервер**
```bash
make run
```

Сервер буде доступний за адресою: http://localhost:8082

## 📦 Доступні команди

```bash
make help           # Показати всі доступні команди
make run            # Запустити сервер (PostgreSQL)
make migrate-up     # Застосувати міграції
make migrate-down   # Відкотити останню міграцію
make docker-up      # Запустити PostgreSQL в Docker
make docker-down    # Зупинити PostgreSQL
make docker-logs    # Показати логи PostgreSQL
make psql           # Підключитись до PostgreSQL
make db-reset       # Повністю очистити БД та застосувати міграції
make install-deps   # Встановити залежності
make test           # Запустити тести
make build          # Зібрати binary
make clean          # Очистити тимчасові файли
```

## 🗄️ База даних

Проект використовує **PostgreSQL** як основну базу даних.

### Структура таблиць

- **messages** - Повідомлення користувачів
- **users** - Користувачі платформи
- **stock_orders** - Заявки на купівлю/продаж акцій

### Міграції

Міграції знаходяться в директорії `migrations/`:

```
migrations/
├── 001_create_messages_table.up.sql
├── 001_create_messages_table.down.sql
├── 002_create_users_table.up.sql
├── 002_create_users_table.down.sql
├── 003_create_stock_orders_table.up.sql
└── 003_create_stock_orders_table.down.sql
```

## 🔧 Конфігурація

Всі налаштування можна змінити через змінні середовища (`.env` файл):

```env
# Server
SERVER_PORT=8082
WS_PATH=/ws

# Database (завжди PostgreSQL)
DB_TYPE=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=stock_hub
POSTGRES_SSLMODE=disable
```

## 🌐 API та WebSocket

### WebSocket Endpoint
```
ws://localhost:8082/ws
```

### Типи повідомлень

1. **Реєстрація користувача**
```json
{
  "type": "register",
  "data": {
    "username": "john_doe"
  }
}
```

2. **Відправка повідомлення**
```json
{
  "type": "message",
  "data": "Hello, world!"
}
```

3. **Створення заявки**
```json
{
  "type": "order",
  "data": {
    "symbol": "AAPL",
    "order_type": "bid",
    "price": 150.50,
    "quantity": 100
  }
}
```

4. **Скасування заявки**
```json
{
  "type": "cancel_order",
  "data": {
    "order_id": 123
  }
}
```

## 📚 Документація

### SQL Документація

В директорії `docs/` є детальна SQL документація:

- [`SQL_TOP_70_QUERIES.md`](./docs/SQL_TOP_70_QUERIES.md) - ТОП-70 найважливіших SQL запитів
- [`SQL_STOCK_HUB_QUERIES.md`](./docs/SQL_STOCK_HUB_QUERIES.md) - Запити для цього проекту
- [`SQL_ADMIN_MONITORING.md`](./docs/SQL_ADMIN_MONITORING.md) - Адміністрування та моніторинг
- [`SQL_OPTIMIZATION_GUIDE.md`](./docs/SQL_OPTIMIZATION_GUIDE.md) - Оптимізація запитів
- [`SQL_ADVANCED_EXAMPLES.md`](./docs/SQL_ADVANCED_EXAMPLES.md) - Просунуті приклади
- [`SQL_README.md`](./docs/SQL_README.md) - Центральний README для всієї SQL документації

### Налаштування та швидкий старт

- [QUICKSTART.md](./QUICKSTART.md) - Швидкий старт за 3 хвилини
- [POSTGRES_SETUP.md](./docs/POSTGRES_SETUP.md) - Детальне налаштування PostgreSQL
- [CHANGELOG.md](./CHANGELOG.md) - Історія змін

## 🛠️ Розробка

### Структура проекту

```
stock_hub/
├── cmd/
│   ├── server/          # Головний сервер
│   └── migrate/         # CLI для міграцій
├── internal/
│   ├── config/          # Конфігурація
│   ├── model/           # Моделі даних
│   └── service/         # Бізнес логіка
├── pkg/
│   ├── migrate/         # Система міграцій
│   ├── sqlutil/         # SQL утиліти
│   └── websocket/       # WebSocket hub
├── migrations/          # SQL міграції
└── static/              # Статичні файли (HTML, CSS, JS)
```

### Запуск тестів

```bash
make test
```

### Збірка бінарника

```bash
make build
# Бінарник буде в bin/stock_hub
./bin/stock_hub
```

## 🐳 Docker

Для легкого запуску PostgreSQL використовуйте Docker Compose:

```bash
# Запустити PostgreSQL
make docker-up

# Зупинити PostgreSQL
make docker-down

# Переглянути логи
make docker-logs
```

## 🔍 Налагодження

### Підключення до PostgreSQL

```bash
# З Makefile
make psql

# Або напряму
psql -U postgres -d stock_hub
```

### Перезавантаження бази даних

```bash
# Повне очищення та повторне застосування міграцій
make db-reset
```

### Перевірка таблиць

```sql
\dt                          # Список таблиць
SELECT * FROM messages;      # Всі повідомлення
SELECT * FROM users;         # Всі користувачі
SELECT * FROM stock_orders;  # Всі заявки
```

## 📝 Файли конфігурації

- [docker-compose.yml](./docker-compose.yml) - Docker конфігурація для PostgreSQL
- [.env.example](./.env.example) - Приклад змінних середовища
- [Makefile](./Makefile) - Команди для розробки

## 🤝 Внесок

Внески вітаються! Будь ласка, створюйте Pull Request або Issue.

## 📄 Ліцензія

MIT License

---

Створено з ❤️ для навчання та практики Go + PostgreSQL + WebSockets
