# Stock Hub - Monorepo 🚀

Мікросервісна архітектура для trading платформи.

---

## 📁 Структура

```
stock_hub/
├── go.work                          # Go workspace (monorepo)
├── README.md                        # Цей файл
│
├── stock_hub_trade/                 # 📈 Trading Service
│   ├── cmd/server/                  # Main HTTP/WebSocket/gRPC сервер
│   ├── internal/                    # Business logic
│   │   ├── service/                 # Domain services
│   │   ├── repository/              # Data access layer
│   │   ├── grpc/                    # gRPC server implementation
│   │   ├── config/                  # Configuration
│   │   └── handler/                 # HTTP handlers
│   ├── pkg/                         # 🔥 SHARED DOMAIN MODELS
│   │   ├── model/                   # User, StockOrder, Message
│   │   ├── websocket/               # WebSocket utilities
│   │   ├── migrate/                 # Database migrations
│   │   └── sqlutil/                 # SQL helpers
│   ├── api/proto/                   # gRPC Protobuf definitions
│   ├── migrations/                  # SQL migration files
│   ├── docs/                        # Documentation
│   ├── go.mod
│   └── README.md
│
└── stock_hub_email_service/         # 📧 Email Service
    ├── cmd/server/                  # Main HTTP сервер
    ├── internal/
    │   ├── service/                 # Email sending logic
    │   └── config/                  # Configuration
    ├── go.mod                       # replace stock_hub_trade => ../stock_hub_trade
    └── README.md
```

---

## 🎯 Сервіси

### 1️⃣ Stock Hub Trade (порт 8082, 50051)

**Технології:**
- HTTP/REST API
- WebSocket (real-time trading)
- gRPC (inter-service communication)
- PostgreSQL

**Endpoints:**
```bash
# HTTP
http://localhost:8082/register
http://localhost:8082/ws

# gRPC
grpc://localhost:50051
```

**RPC методи:**
- `GetUser(user_id)` → User
- `GetOrder(order_id)` → StockOrder
- `GetOrders(user_id)` → []StockOrder
- `GetMessages(limit)` → []Message

---

### 2️⃣ Email Service (порт 50052)

**Технології:**
- HTTP API
- gRPC client (викликає stock_hub_trade)
- SMTP (email sending)

**Features:**
- Welcome emails
- Order confirmations
- Order cancellations

---

## 🚀 Швидкий старт

### Prerequisites

```bash
# Встановити protoc
brew install protobuf

# Встановити Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Запуск

```bash
# 1. Запустити PostgreSQL
cd stock_hub_trade
docker-compose up -d

# 2. Згенерувати proto файли (якщо змінювали)
make proto-gen

# 3. Запустити міграції
make migrate-up

# 4. Запустити trading service
make run-postgres

# 5. В іншому терміналі запустити email service
cd ../stock_hub_email_service
make run
```

---

## 🔧 Розробка

### Go Workspace

Проект використовує `go.work` для роботи з monorepo:

```bash
# Оновити dependencies
go work sync

# Запустити тести в усіх модулях
go test ./...

# Запустити тести в конкретному сервісі
cd stock_hub_trade && go test ./...
```

### Shared Models

Обидва сервіси використовують спільні моделі з `stock_hub_trade/pkg/model/`:

```go
// stock_hub_email_service/internal/service/email_service.go
import "stock_hub_trade/pkg/model"  // ← Shared models!

func (s *EmailService) SendWelcomeEmail(user *model.User) error {
    // ...
}
```

### gRPC Communication

Email service викликає trading service через gRPC:

```go
// Email service → Trading service
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := proto.NewStockHubServiceClient(conn)
user, _ := client.GetUser(ctx, &proto.GetUserRequest{UserId: 1})
```

---

## 📚 Документація

### Trading Service
- `stock_hub_trade/README.md` - Повний опис
- `stock_hub_trade/QUICKSTART.md` - Швидкий старт
- `stock_hub_trade/docs/GRPC_SETUP.md` - gRPC setup
- `stock_hub_trade/docs/REPOSITORY_PATTERN.md` - Архітектура
- `stock_hub_trade/docs/REST_VS_GRPC_VS_WEBSOCKET.md` - Порівняння протоколів

### Email Service
- `stock_hub_email_service/README.md` - Повний опис

---

## 🧪 Тестування

```bash
# Всі тести
go test ./...

# Trading service
cd stock_hub_trade
go test ./internal/service/...
go test ./internal/repository/...
go test ./pkg/migrate/...

# Email service
cd stock_hub_email_service
go test ./...
```

---

## 🏗️ Архітектура

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│         Frameworks & Drivers            │
│  HTTP, gRPC, PostgreSQL, SMTP, WebSocket│
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│        Interface Adapters               │
│  Handlers, Repositories, Converters     │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│           Use Cases                     │
│  Services (CreateOrder, SendEmail)      │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│            Entities                     │
│  pkg/model (User, StockOrder, Message)  │
└─────────────────────────────────────────┘
```

### Microservices Communication

```
User → stock_hub_trade:8082/register
         ↓
      [Create user in DB]
         ↓
      gRPC call → stock_hub_email_service:50052
         ↓
      gRPC client → stock_hub_trade:50051/GetUser
         ↓
      [Send welcome email via SMTP]
```

---

## 🔐 Environment Variables

### Trading Service (.env)
```bash
DB_TYPE=postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=stock_hub
SERVER_PORT=8082
GRPC_PORT=50051
WS_PATH=/ws
```

### Email Service (.env)
```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
EMAIL_FROM=noreply@stockhub.com
DATABASE_DSN=postgres://postgres:postgres@localhost:5432/stock_hub?sslmode=disable
SERVER_PORT=50052
```

---

## 📦 Dependencies

### Trading Service
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/lib/pq` - PostgreSQL driver
- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protobuf support
- `github.com/DATA-DOG/go-sqlmock` - SQL mocking for tests

### Email Service
- `github.com/lib/pq` - PostgreSQL driver (for fetching data)
- `stock_hub_trade/pkg/model` - Shared domain models

---

## 🎓 Навчальні матеріали

- [Clean Architecture vs MVVM](stock_hub_trade/docs/MICROSERVICES_ARCHITECTURE.md)
- [Repository Pattern](stock_hub_trade/docs/REPOSITORY_PATTERN.md)
- [Go Interfaces Explained](stock_hub_trade/docs/GO_INTERFACES_EXPLAINED.md)
- [Table-Driven Tests](stock_hub_trade/docs/TABLE_DRIVEN_TESTS.md)
- [gRPC Integration](stock_hub_trade/docs/GRPC_QUICKSTART.md)

---

## 🤝 Contribution

1. Створити feature branch
2. Написати тести
3. Оновити документацію
4. Створити PR

---

## 📝 License

MIT
