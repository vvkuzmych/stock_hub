# Stock Hub Monorepo - Quick Start ⚡

Швидкий старт для роботи з monorepo.

---

## 🚀 Перше використання

### 1. Клонування (якщо ще не клонували)

```bash
git clone <your-repo-url> stock_hub
cd stock_hub
```

### 2. Перевірка структури

```bash
tree -L 2 -I '.git|node_modules|.idea'
```

**Очікується:**
```
stock_hub/
├── stock_hub_trade/        # Trading service
├── stock_hub_email_service/ # Email service
└── go.work                 # Go workspace
```

### 3. Синхронізація dependencies

```bash
go work sync
```

---

## 🔧 Налаштування

### Trading Service

```bash
cd stock_hub_trade

# 1. Скопіювати .env.example
cp .env.example .env

# 2. Відредагувати .env (опціонально)
nano .env

# 3. Запустити PostgreSQL
docker-compose up -d

# 4. Запустити міграції
make migrate-up
```

### Email Service

```bash
cd stock_hub_email_service

# 1. Скопіювати .env.example
cp .env.example .env

# 2. Відредагувати SMTP налаштування
nano .env
```

---

## ▶️ Запуск

### Варіант 1: Окремі термінали (рекомендовано)

```bash
# Terminal 1: Trading service
cd stock_hub_trade
make run-postgres

# Terminal 2: Email service
cd stock_hub_email_service
make run
```

### Варіант 2: Один термінал (background)

```bash
cd stock_hub_trade
make run-postgres &

cd ../stock_hub_email_service
make run
```

---

## 🧪 Тестування

### Всі тести

```bash
# З root directory
cd /Users/vkuzm/GolandProjects/stock_hub
go test ./...
```

### Окремі сервіси

```bash
# Trading service
cd stock_hub_trade
go test ./internal/service/...
go test ./internal/repository/...

# Email service
cd stock_hub_email_service
go test ./...
```

---

## 🔍 Перевірка

### Trading Service

```bash
# HTTP API
curl http://localhost:8082/

# WebSocket (через websocat або browser)
websocat ws://localhost:8082/ws

# gRPC
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"user_id": 1}' localhost:50051 stock_hub.StockHubService/GetUser
```

### Email Service

```bash
# HTTP API
curl http://localhost:50052/health
```

---

## 📦 Build

```bash
# Trading service
cd stock_hub_trade
go build -o bin/trade ./cmd/server/

# Email service
cd stock_hub_email_service
go build -o bin/email ./cmd/server/

# Запустити бінарі
./stock_hub_trade/bin/trade
./stock_hub_email_service/bin/email
```

---

## 🔄 Оновлення Dependencies

```bash
# З root directory
go work sync

# Або окремо для кожного сервісу
cd stock_hub_trade && go mod tidy
cd ../stock_hub_email_service && go mod tidy
```

---

## 🐛 Дебаг

### Перевірити імпорти

```bash
# Trading service
cd stock_hub_trade
go list -m all | grep stock_hub

# Email service
cd stock_hub_email_service
go list -m all | grep stock_hub
```

### Перевірити workspace

```bash
go work edit -json
```

### Очистити cache

```bash
go clean -modcache
go work sync
```

---

## 📚 Документація

- **Root README**: `README.md` - Загальний опис
- **Migration Guide**: `MONOREPO_MIGRATION.md` - Деталі міграції
- **Trading Service**: `stock_hub_trade/README.md`
- **Email Service**: `stock_hub_email_service/README.md`
- **Architecture**: `stock_hub_trade/docs/`

---

## 🎯 Типові завдання

### Додати нову feature в trading service

```bash
cd stock_hub_trade
git checkout -b feature/my-feature
# Зробити зміни
go test ./...
git add .
git commit -m "feat: add my feature"
```

### Оновити shared models

```bash
# Models в stock_hub_trade/pkg/model/
cd stock_hub_trade/pkg/model
# Змінити user.go, stock_order.go, message.go

# Автоматично доступно в email_service через go.work!
cd ../../stock_hub_email_service
go test ./...  # Переконатись що все працює
```

### Згенерувати proto файли

```bash
cd stock_hub_trade
make proto-gen

# Або вручну
cd api/proto && ./generate.sh
```

---

## ⚠️ Часті помилки

### "package stock_hub/... not found"

**Причина:** Старі import paths.

**Рішення:**
```bash
find . -name "*.go" -exec sed -i '' 's|"stock_hub/|"stock_hub_trade/|g' {} \;
```

### "no required module provides package..."

**Причина:** Dependencies не синхронізовані.

**Рішення:**
```bash
go work sync
cd stock_hub_trade && go mod tidy
cd ../stock_hub_email_service && go mod tidy
```

### PostgreSQL не запускається

**Причина:** Порт 5432 зайнятий.

**Рішення:**
```bash
# Знайти процес
lsof -i :5432

# Або змінити порт в docker-compose.yml
```

---

## ✅ Checklist

Перед початком роботи:

- [ ] `git pull` успішно
- [ ] `go work sync` виконано
- [ ] `.env` файли налаштовані
- [ ] PostgreSQL запущений
- [ ] Міграції виконані (`make migrate-up`)
- [ ] Обидва сервіси запускаються
- [ ] Тести проходять (`go test ./...`)

**Готово до роботи!** 🎉
