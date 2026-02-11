# Email Service 📧

Email notification service for Stock Hub using shared models from `stock_hub/pkg/model`

---

## 🎯 Architecture

This service demonstrates **Clean Architecture** with shared domain models:

```
email_service/                    stock_hub/
├── cmd/server/                   ├── pkg/
│   └── main.go                   │   └── model/        ← SHARED ENTITIES
├── internal/                     │       ├── user.go
│   ├── service/                  │       ├── message.go
│   │   └── email_service.go      │       └── stock_order.go
│   └── config/                   │
│       └── config.go              └── internal/
└── go.mod                            └── service/

Clean Architecture Layers:
  1. Entities: model.User, model.StockOrder (in stock_hub/pkg/model)
  2. Use Cases: SendWelcomeEmail, SendOrderConfirmation
  3. Interface Adapters: HTTP handlers, DB queries
  4. Frameworks: SMTP, PostgreSQL
```

---

## ✉️ Features

- 📨 **Welcome Email**: Sent when user registers
- 📊 **Order Confirmation**: Sent when order is created
- ❌ **Order Cancellation**: Sent when order is cancelled
- 🔄 **Shared Models**: Uses `stock_hub/pkg/model` (Clean Architecture!)

---

## 🚀 Quick Start

### 1. Setup Environment

```bash
# Copy environment file
cp .env.example .env

# Configure SMTP (Gmail example)
# 1. Enable 2FA in Gmail
# 2. Generate App Password: https://myaccount.google.com/apppasswords
# 3. Update .env with your credentials
```

### 2. Configure go.mod for Local Development

```bash
# For local development, add replace directive to go.mod:
echo "replace stock_hub => ../stock_hub" >> go.mod

# Then tidy
go mod tidy
```

### 3. Run Service

```bash
go run cmd/server/main.go
```

---

## 📡 API Endpoints

### Health Check

```bash
GET /health

Response:
{
  "status": "ok"
}
```

### Send Email

```bash
POST /send
Content-Type: application/json

# Welcome Email
{
  "type": "welcome",
  "user_id": 1
}

# Order Confirmation
{
  "type": "order_confirmation",
  "order_id": 123
}

# Order Cancellation
{
  "type": "order_cancellation",
  "order_id": 123
}

Response:
{
  "status": "sent"
}
```

---

## 🧪 Testing

### Manual Test with curl

```bash
# Welcome email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Order confirmation
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_confirmation", "order_id": 1}'

# Health check
curl http://localhost:8083/health
```

---

## 🔧 Configuration

### Environment Variables

```
SERVER_PORT=8083                  # HTTP server port
SMTP_HOST=smtp.gmail.com          # SMTP server
SMTP_PORT=587                     # SMTP port (587 for STARTTLS)
SMTP_USER=your-email@gmail.com    # SMTP username
SMTP_PASSWORD=your-app-password   # SMTP password (App Password for Gmail)
FROM_EMAIL=noreply@stock-hub.com  # From address
FROM_NAME=Stock Hub               # From name
DATABASE_DSN=host=localhost...    # PostgreSQL connection (shared with stock_hub)
```

### Gmail Setup

1. Enable 2-Factor Authentication
2. Generate App Password:
   - Visit: https://myaccount.google.com/apppasswords
   - App: Mail
   - Device: Other (Custom name)
   - Copy generated password
3. Use App Password in `SMTP_PASSWORD`

---

## 🏗️ Project Structure

```
email_service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point, HTTP handlers
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loader
│   └── service/
│       └── email_service.go     # Email sending logic
├── .env.example                 # Environment template
├── go.mod                       # Module definition
└── README.md                    # This file
```

---

## 🔗 Integration with Stock Hub

### Option 1: HTTP API Call

```go
// In stock_hub/cmd/server/main.go

func notifyEmailService(eventType string, id int64) {
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
		log.Printf("Failed to notify email service: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Printf("Email service error: %d", resp.StatusCode)
	}
}

// Call after user registration
user, err := ctx.userService.RegisterUser(data.Username)
if err == nil {
	go notifyEmailService("welcome", user.ID)
}

// Call after order creation
order, err := ctx.orderService.CreateOrder(...)
if err == nil {
	go notifyEmailService("order_confirmation", order.ID)
}
```

### Option 2: Message Queue (Future)

```
stock_hub → RabbitMQ/Kafka → email_service
            (event stream)
```

---

## 📦 Shared Models

This service uses **shared domain models** from `stock_hub/pkg/model`:

```go
import "stock_hub/pkg/model"

// Available types:
model.User
model.Message  
model.StockOrder
model.OrderType (OrderTypeBid, OrderTypeAsk)
```

This follows **Clean Architecture** principles:
- ✅ Domain entities (models) are shared
- ✅ Use cases (services) are independent
- ✅ Frameworks (SMTP, DB) are isolated

---

## 🔒 Security

- ⚠️ **Never commit `.env`** (contains SMTP credentials)
- ✅ Use **App Passwords** (not your main Gmail password)
- ✅ **TLS/STARTTLS** enabled (port 587)
- ⚠️ **Rate limiting** recommended (prevent spam)
- ✅ **Input validation** on all endpoints

---

## 🚀 Production Deployment

```bash
# Build binary
go build -o email-service cmd/server/main.go

# Run with environment file
./email-service

# Or with Docker
docker build -t email-service .
docker run -p 8083:8083 --env-file .env email-service
```

---

## 📚 Dependencies

```go
// Required:
github.com/lib/pq              // PostgreSQL driver
stock_hub/pkg/model            // Shared domain models

// Standard library:
net/smtp                       // Email sending
database/sql                   // Database access
encoding/json                  // JSON handling
```

---

## 🎯 Future Enhancements

- [ ] Email templates (HTML)
- [ ] Queue-based processing (RabbitMQ/Kafka)
- [ ] Retry mechanism for failed emails
- [ ] Email tracking (open/click rates)
- [ ] Bulk email sending
- [ ] Email scheduling
- [ ] Multiple SMTP providers (SendGrid, Mailgun, AWS SES)

---

Built with ❤️ using Clean Architecture principles
