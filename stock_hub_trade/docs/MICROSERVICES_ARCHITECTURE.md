# Microservices Architecture з Clean Architecture 🏗️

Цей документ описує архітектуру мікросервісів `stock_hub` та `email_service` з використанням принципів **Clean Architecture**.

---

## 🎯 Архітектурний Огляд

```
┌─────────────────────────────────────────────────────────────────┐
│                     MICROSERVICES ECOSYSTEM                      │
└─────────────────────────────────────────────────────────────────┘

                    ┌──────────────────────┐
                    │   SHARED MODELS      │
                    │  (Domain Entities)   │
                    │                      │
                    │  stock_hub/pkg/      │
                    │    └── model/        │
                    │        ├── user.go   │
                    │        ├── message.go│
                    │        └── stock_    │
                    │            order.go  │
                    └──────────┬───────────┘
                               │
            ┌──────────────────┴───────────────────┐
            │                                      │
            ▼                                      ▼
┌────────────────────────┐             ┌────────────────────────┐
│    STOCK_HUB SERVICE   │◄────────────┤  EMAIL_SERVICE         │
│                        │   HTTP API  │                        │
│  Port: 8082            │   or Queue  │  Port: 8083            │
│                        │             │                        │
│  ┌──────────────────┐  │             │  ┌──────────────────┐  │
│  │   WebSocket Hub  │  │             │  │  Email Sender    │  │
│  │   User Service   │  │             │  │  SMTP Client     │  │
│  │   Order Service  │  │             │  │  Template Engine │  │
│  └──────────────────┘  │             │  └──────────────────┘  │
│           │            │             │           │            │
│           ▼            │             │           ▼            │
│  ┌──────────────────┐  │             │  ┌──────────────────┐  │
│  │   PostgreSQL     │  │             │  │   PostgreSQL     │  │
│  │   (stock_hub)    │◄─┼─────────────┼──│   (stock_hub)    │  │
│  └──────────────────┘  │   Shared DB │  └──────────────────┘  │
└────────────────────────┘             └────────────────────────┘
```

---

## 📦 Clean Architecture Layers

### 1. **Entities (Domain Models)** - `stock_hub/pkg/model`

Це **ядро системи**, незалежне від будь-яких фреймворків або баз даних.

```go
// stock_hub/pkg/model/user.go
package model

import "time"

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}
```

```go
// stock_hub/pkg/model/stock_order.go
package model

type OrderType string

const (
    OrderTypeBid OrderType = "bid"
    OrderTypeAsk OrderType = "ask"
)

type StockOrder struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Username  string    `json:"username"`
    Symbol    string    `json:"symbol"`
    OrderType OrderType `json:"order_type"`
    Price     float64   `json:"price"`
    Quantity  int       `json:"quantity"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Чому в `pkg/model`?**

✅ Моделі використовуються **кількома сервісами** (`stock_hub`, `email_service`)  
✅ Це **Domain Entities** - ядро бізнес-логіки  
✅ **Zero dependencies** - не залежать від баз даних, фреймворків, транспортів  
✅ Можуть бути опубліковані як **окремий Go модуль** для інших команд

---

### 2. **Use Cases (Business Logic)** - `internal/service`

Кожен сервіс має свої **Use Cases**, які використовують Domain Entities.

#### Stock Hub Use Cases

```go
// stock_hub/internal/service/user_service.go
package service

import "stock_hub/pkg/model"

type UserService struct {
    db *sql.DB
}

func (s *UserService) RegisterUser(username string) (*model.User, error) {
    // Business logic: register user
}
```

```go
// stock_hub/internal/service/stock_order_service.go
package service

import "stock_hub/pkg/model"

type StockOrderService struct {
    db *sql.DB
}

func (s *StockOrderService) CreateOrder(...) (*model.StockOrder, error) {
    // Business logic: validate, create order
}
```

#### Email Service Use Cases

```go
// email_service/internal/service/email_service.go
package service

import "stock_hub/pkg/model" // ← Import shared models

type EmailService struct {
    smtpHost string
    // ...
}

func (s *EmailService) SendWelcomeEmail(user *model.User) error {
    // Business logic: compose and send email
}

func (s *EmailService) SendOrderConfirmation(order *model.StockOrder) error {
    // Business logic: compose and send order email
}
```

---

### 3. **Interface Adapters (Controllers, Gateways)** - `cmd/server`

Це HTTP handlers, WebSocket handlers, DB queries - все, що адаптує Use Cases до конкретних технологій.

#### Stock Hub Adapter

```go
// stock_hub/cmd/server/main.go
package main

import (
    "stock_hub/internal/service"
    "stock_hub/pkg/model"
)

func handleMessage(ctx *ServerContext, client *ws.Client, msg *WSMessage) {
    switch msg.Type {
    case "register":
        // Call use case
        user, err := ctx.userService.RegisterUser(data.Username)
        if err == nil {
            // Notify email service
            notifyEmailService("welcome", user.ID)
        }
    
    case "order":
        // Call use case
        order, err := ctx.orderService.CreateOrder(...)
        if err == nil {
            // Notify email service
            notifyEmailService("order_confirmation", order.ID)
        }
    }
}
```

#### Email Service Adapter

```go
// email_service/cmd/server/main.go
package main

import (
    "email_service/internal/service"
    "stock_hub/pkg/model"
)

func sendEmailHandler(emailService *service.EmailService, db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        switch req.Type {
        case "welcome":
            user, _ := fetchUser(db, req.UserID)
            emailService.SendWelcomeEmail(user)
        
        case "order_confirmation":
            order, _ := fetchOrder(db, req.OrderID)
            emailService.SendOrderConfirmation(order)
        }
    }
}
```

---

### 4. **Frameworks & Drivers (Infrastructure)** - бази даних, SMTP, WebSocket

- PostgreSQL
- SMTP (Gmail, SendGrid)
- WebSocket (Gorilla)
- HTTP (net/http)

---

## 🔄 Communication Flow

### Scenario 1: User Registration

```
┌─────────────┐                    ┌─────────────┐                    ┌─────────────┐
│   Client    │                    │ Stock Hub   │                    │   Email     │
│  (Browser)  │                    │   Service   │                    │   Service   │
└──────┬──────┘                    └──────┬──────┘                    └──────┬──────┘
       │                                  │                                  │
       │ 1. Register (WS)                 │                                  │
       ├─────────────────────────────────►│                                  │
       │                                  │                                  │
       │                                  │ 2. Save User to DB               │
       │                                  ├──────────┐                       │
       │                                  │          │                       │
       │                                  │◄─────────┘                       │
       │                                  │                                  │
       │                                  │ 3. POST /send (welcome)          │
       │                                  ├─────────────────────────────────►│
       │                                  │                                  │
       │                                  │                                  │ 4. Fetch User
       │                                  │                                  ├────────┐
       │                                  │                                  │        │
       │                                  │                                  │◄───────┘
       │                                  │                                  │
       │                                  │                                  │ 5. Send Email
       │                                  │                                  ├────────┐
       │                                  │                                  │        │
       │                                  │                                  │◄───────┘
       │                                  │                                  │
       │                                  │          6. OK                   │
       │                                  │◄─────────────────────────────────┤
       │                                  │                                  │
       │  7. Registered (WS)              │                                  │
       │◄─────────────────────────────────┤                                  │
       │                                  │                                  │
```

### Scenario 2: Order Creation

```
┌─────────────┐                    ┌─────────────┐                    ┌─────────────┐
│   Client    │                    │ Stock Hub   │                    │   Email     │
│  (Browser)  │                    │   Service   │                    │   Service   │
└──────┬──────┘                    └──────┬──────┘                    └──────┬──────┘
       │                                  │                                  │
       │ 1. Create Order (WS)             │                                  │
       ├─────────────────────────────────►│                                  │
       │                                  │                                  │
       │                                  │ 2. Validate Order                │
       │                                  ├──────────┐                       │
       │                                  │          │                       │
       │                                  │◄─────────┘                       │
       │                                  │                                  │
       │                                  │ 3. Save Order to DB              │
       │                                  ├──────────┐                       │
       │                                  │          │                       │
       │                                  │◄─────────┘                       │
       │                                  │                                  │
       │                                  │ 4. POST /send (order_confirmation)│
       │                                  ├─────────────────────────────────►│
       │                                  │                                  │
       │                                  │                                  │ 5. Fetch Order
       │                                  │                                  ├────────┐
       │                                  │                                  │        │
       │                                  │                                  │◄───────┘
       │                                  │                                  │
       │                                  │                                  │ 6. Send Email
       │                                  │                                  ├────────┐
       │                                  │                                  │        │
       │                                  │                                  │◄───────┘
       │                                  │                                  │
       │                                  │          7. OK                   │
       │                                  │◄─────────────────────────────────┤
       │                                  │                                  │
       │  8. Order Created (WS)           │                                  │
       │◄─────────────────────────────────┤                                  │
       │                                  │                                  │
```

---

## 🛠️ Integration Implementation

### Option 1: HTTP API (Current)

#### Stock Hub → Email Service

```go
// stock_hub/cmd/server/main.go

func notifyEmailService(eventType string, id int64) error {
    payload := map[string]interface{}{
        "type": eventType,
    }
    
    if eventType == "welcome" {
        payload["user_id"] = id
    } else {
        payload["order_id"] = id
    }
    
    jsonData, _ := json.Marshal(payload)
    resp, err := http.Post("http://localhost:8083/send", 
        "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to notify email service: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("email service error: %d", resp.StatusCode)
    }
    
    return nil
}

// After user registration
user, err := ctx.userService.RegisterUser(data.Username)
if err == nil {
    go notifyEmailService("welcome", user.ID)
}

// After order creation
order, err := ctx.orderService.CreateOrder(...)
if err == nil {
    go notifyEmailService("order_confirmation", order.ID)
}
```

---

### Option 2: Message Queue (Future - Recommended for Production)

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│ Stock Hub   │─event──►│  RabbitMQ/  │◄──poll──│   Email     │
│   Service   │         │   Kafka     │         │   Service   │
└─────────────┘         └─────────────┘         └─────────────┘
                            │
                            │ Topics:
                            │  - user.registered
                            │  - order.created
                            │  - order.cancelled
                            │
```

#### Producer (Stock Hub)

```go
// stock_hub/internal/event/publisher.go
package event

import "stock_hub/pkg/model"

type Publisher interface {
    PublishUserRegistered(user *model.User) error
    PublishOrderCreated(order *model.StockOrder) error
    PublishOrderCancelled(order *model.StockOrder) error
}

// RabbitMQ implementation
type RabbitMQPublisher struct {
    conn *amqp.Connection
}

func (p *RabbitMQPublisher) PublishUserRegistered(user *model.User) error {
    payload, _ := json.Marshal(user)
    return p.publish("user.registered", payload)
}
```

#### Consumer (Email Service)

```go
// email_service/internal/event/consumer.go
package event

import "stock_hub/pkg/model"

type Consumer struct {
    emailService *service.EmailService
}

func (c *Consumer) Start() {
    // Subscribe to topics
    c.subscribe("user.registered", c.handleUserRegistered)
    c.subscribe("order.created", c.handleOrderCreated)
}

func (c *Consumer) handleUserRegistered(payload []byte) error {
    var user model.User
    json.Unmarshal(payload, &user)
    return c.emailService.SendWelcomeEmail(&user)
}
```

---

## 📁 Project Structure

```
/Users/vkuzm/GolandProjects/
│
├── stock_hub/                          # Main trading service
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                 # HTTP + WebSocket handlers
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go               # Configuration
│   │   └── service/
│   │       ├── user_service.go         # User business logic
│   │       ├── stock_order_service.go  # Order business logic
│   │       └── message_service.go      # Message business logic
│   ├── pkg/
│   │   ├── model/                      # ← SHARED DOMAIN ENTITIES
│   │   │   ├── user.go
│   │   │   ├── message.go
│   │   │   └── stock_order.go
│   │   ├── websocket/
│   │   │   └── hub.go
│   │   └── sqlutil/
│   │       └── placeholder.go
│   ├── migrations/
│   │   ├── 001_create_users_table.up.sql
│   │   ├── 002_create_messages_table.up.sql
│   │   └── 003_create_stock_orders_table.up.sql
│   ├── go.mod
│   └── README.md
│
└── email_service/                      # Email notification service
    ├── cmd/
    │   └── server/
    │       └── main.go                 # HTTP handlers
    ├── internal/
    │   ├── config/
    │   │   └── config.go               # Configuration
    │   └── service/
    │       └── email_service.go        # Email business logic
    ├── go.mod                          # replace stock_hub => ../stock_hub
    ├── .env.example
    ├── Makefile
    └── README.md

Clean Architecture Mapping:
  - Entities: stock_hub/pkg/model/
  - Use Cases: */internal/service/
  - Interface Adapters: */cmd/server/
  - Frameworks: PostgreSQL, SMTP, WebSocket
```

---

## 🔑 Key Design Principles

### 1. **Dependency Rule**

```
Entities ← Use Cases ← Interface Adapters ← Frameworks
   ↑          ↑              ↑                   ↑
   │          │              │                   │
No dependencies outward      │                   │
                             │                   │
                  Dependencies point inward
```

**Правило:** Dependencies завжди спрямовані **всередину** (до Domain Entities).

- ✅ `email_service` імпортує `stock_hub/pkg/model` (Use Case → Entity)
- ✅ `service` імпортує `model` (Use Case → Entity)
- ❌ `model` **НЕ** імпортує `service` (Entity НЕ залежить від Use Case)
- ❌ `model` **НЕ** імпортує `database/sql` (Entity НЕ залежить від DB)

---

### 2. **Shared Entities в pkg/**

Чому моделі в `pkg/` замість `internal/`?

| Аспект          | `internal/model`          | `pkg/model`               |
|-----------------|---------------------------|---------------------------|
| **Видимість**   | Тільки для цього проекту  | Публічна (інші сервіси)   |
| **Призначення** | Приватна реалізація       | Shared Domain Entities    |
| **Go Import**   | `stock_hub/internal/model`| `stock_hub/pkg/model`     |
| **Архітектура** | Monolith                  | Microservices             |

**Коли використовувати `pkg/`?**

✅ Моделі використовуються **кількома сервісами**  
✅ Domain Entities, які описують **бізнес-логіку**  
✅ Zero dependencies на фреймворки/БД  
✅ Можливість публікації як **окремий модуль**

**Коли використовувати `internal/`?**

✅ Приватна реалізація, специфічна для одного сервісу  
✅ Monolithic додаток без планів на мікросервіси  
✅ Внутрішні утиліти, які не варто експортувати

---

### 3. **go.mod replace для Local Development**

```go
// email_service/go.mod
module email_service

go 1.25

require (
    stock_hub v0.0.0
)

// Local development: use local stock_hub
replace stock_hub => ../stock_hub
```

**Production:** Publish `stock_hub/pkg/model` as separate module:

```bash
# Option 1: Separate module
cd stock_hub/pkg/model
go mod init github.com/yourusername/stock-hub-models
git tag v1.0.0
git push --tags

# email_service/go.mod
require github.com/yourusername/stock-hub-models v1.0.0
```

```bash
# Option 2: Go workspace (local development)
go work init
go work use ./stock_hub
go work use ./email_service
```

---

## 🚀 Running Services

### Terminal 1: PostgreSQL

```bash
cd stock_hub
make docker-up
```

### Terminal 2: Stock Hub

```bash
cd stock_hub
make migrate-up
make run

# Output:
# 🚀 Starting Stock Hub with postgres database
# 📊 DSN: host=localhost port=5432 user=postgres password=*** dbname=stock_hub sslmode=disable
# ✅ Database connected and migrations applied
# 🌐 Web interface: http://localhost:8082
# 🔌 WebSocket endpoint: ws://localhost:8082/ws
```

### Terminal 3: Email Service

```bash
cd email_service
cp .env.example .env
# Edit .env with SMTP credentials
make run

# Output:
# 📧 Starting Email Service
# Server Port: 8083
# SMTP: smtp.gmail.com:587
# ✅ Database connected
# 🌐 Email service running on http://localhost:8083
```

---

## 🧪 Testing

### Test User Registration with Email

```bash
# 1. Register user (WebSocket to Stock Hub)
wscat -c ws://localhost:8082/ws
> {"type": "register", "data": {"username": "john_doe"}}

# 2. Check email (should receive welcome email)

# 3. Verify email service was called
curl http://localhost:8083/health
```

### Test Order Creation with Email

```bash
# 1. Create order (WebSocket to Stock Hub)
> {"type": "order", "data": {"symbol": "AAPL", "order_type": "bid", "price": 150.00, "quantity": 10}}

# 2. Check email (should receive order confirmation)
```

---

## 📊 Benefits of This Architecture

### ✅ Clean Architecture

- **Testable**: Domain logic не залежить від БД/фреймворків
- **Independent**: Entities можна змінювати незалежно від UI/DB
- **Flexible**: Легко замінити PostgreSQL на MongoDB, SMTP на SendGrid

### ✅ Microservices

- **Scalability**: `email_service` можна скейлити незалежно
- **Resilience**: Якщо email service падає, trading продовжує працювати
- **Technology Freedom**: Кожен сервіс може використовувати різні технології

### ✅ Shared Models

- **Consistency**: Одні й ті ж Domain Entities в усіх сервісах
- **Type Safety**: Go compiler перевіряє типи між сервісами
- **No Duplication**: DRY принцип для моделей

---

## 🔮 Future Enhancements

- [ ] **Event-Driven Architecture**: RabbitMQ/Kafka для асинхронної комунікації
- [ ] **API Gateway**: Єдина точка входу для всіх сервісів
- [ ] **Service Discovery**: Consul/Eureka для динамічного знаходження сервісів
- [ ] **Circuit Breaker**: Resilience patterns (Hystrix, resilience4go)
- [ ] **Distributed Tracing**: OpenTelemetry, Jaeger
- [ ] **Monitoring**: Prometheus + Grafana
- [ ] **Authentication Service**: JWT-based auth microservice
- [ ] **Notification Service**: Push notifications (Firebase, WebPush)
- [ ] **Analytics Service**: Real-time trading analytics

---

## 📚 Resources

- [Clean Architecture (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Microservices Patterns](https://microservices.io/patterns/microservices.html)
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html)

---

Built with ❤️ using **Clean Architecture** principles
