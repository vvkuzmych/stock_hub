# Email Service 📧

Email notification service for Stock Hub using shared models from `stock_hub_trade/pkg/model`

---

## ✅ Status: WORKING

```bash
# Health check
curl http://localhost:8083/health
# Response: {"status":"ok"}

# Send email (mock mode)
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
# Response: {"status":"sent"}
```

---

## 🎯 Architecture

This service demonstrates **Clean Architecture** with shared domain models:

```
email_service/                    stock_hub_trade/
├── cmd/server/                   ├── pkg/
│   └── main.go                   │   └── model/        ← SHARED ENTITIES
├── internal/                     │       ├── user.go
│   ├── service/                  │       ├── message.go
│   │   └── email_service.go      │       └── stock_order.go
│   └── config/                   │
│       └── config.go             └── PostgreSQL (shared)
└── go.mod → replace stock_hub_trade => ../stock_hub_trade
```

---

## ✉️ Features

- 📨 **Welcome Email**: Sent when user registers
- 📊 **Order Confirmation**: Sent when order is created  
- ❌ **Order Cancellation**: Sent when order is cancelled
- 🔄 **Shared Models**: Uses `stock_hub_trade/pkg/model`
- 🧪 **Mock Mode**: Development without real SMTP

---

## 🚀 Quick Start

### 1. Setup

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# Copy .env.example
cp .env.example .env

# Edit .env (MOCK_EMAIL=true for development)
nano .env
```

### 2. Setup Email Provider

**Option A: Ethereal Email (RECOMMENDED - 100% FREE!)** 🌟

```bash
# Auto-create free test account
./setup_ethereal.sh

# Copy credentials to .env (shown by script)
nano .env
```

**Option B: Mock Mode (for quick dev)**

```bash
# .env
MOCK_EMAIL=true
```

### 3. Run

```bash
make run
```

**Output (Ethereal):**
```
📧 Starting Email Service
SMTP: smtp.ethereal.email:587
✅ Database connected
🌐 Email service running on http://localhost:8083
```

**Output (Mock):**
```
⚠️  MOCK MODE ENABLED - Emails will NOT be sent
```

### 4. Test

```bash
# Welcome email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Expected: {"status":"sent"}
```

### 5. View Email Online

**If using Ethereal:**
- Open https://ethereal.email/messages
- See your email in browser! 📧

**If using Mock Mode:**
- Check terminal logs

---

## ⚙️ Configuration

### Environment Variables (.env)

```bash
# Server
SERVER_PORT=8083

# Database (shared with stock_hub_trade)
DATABASE_DSN=host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable

# Mock Mode (for development)
MOCK_EMAIL=true  # Set to false for real SMTP

# SMTP (only needed if MOCK_EMAIL=false)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=noreply@stockhub.com
FROM_NAME=Stock Hub
```

---

## 📧 Email Modes

### 1. Ethereal Email (Testing) 🌟 RECOMMENDED

```bash
# .env
MOCK_EMAIL=false
SMTP_HOST=smtp.ethereal.email
SMTP_PORT=587
SMTP_USER=auto-generated@ethereal.email
SMTP_PASSWORD=auto-generated-password
```

**Setup:**
```bash
./setup_ethereal.sh  # Auto-creates free account
```

**Behavior:**
- ✅ **Emails sent to online inbox**
- ✅ **View in browser**: https://ethereal.email/messages
- ✅ **100% FREE, no limits**
- ✅ **No registration needed**
- ✅ **Perfect for testing**

### 2. Mock Mode (Development)

```bash
# .env
MOCK_EMAIL=true
```

**Behavior:**
- Emails are **NOT actually sent**
- Email content is **logged to console**
- No SMTP credentials needed
- Fastest for development

**Example log:**
```
📧 [MOCK] Email NOT actually sent
   To:      john@example.com
   Subject: Welcome to Stock Hub!
   ✅ Mock email logged successfully
```

### 3. Real SMTP (Production)

```bash
# Gmail example
MOCK_EMAIL=false
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**Setup Gmail:**
1. Enable 2FA: https://myaccount.google.com/security
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use app password in `.env`

**Behavior:**
- Emails are **sent to real users**
- Production use only

---

## 🧪 Testing

### Health Check

```bash
curl http://localhost:8083/health
# Expected: {"status":"ok"}
```

### Send Welcome Email

```bash
# Prerequisites: user with ID=1 exists in database

# Test
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Expected: {"status":"sent"}
```

### Send Order Confirmation

```bash
# Prerequisites: order with ID=1 exists in database

curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_confirmation", "order_id": 1}'
```

### Send Order Cancellation

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_cancellation", "order_id": 1}'
```

---

## 📡 API Endpoints

### POST /send

Send email notification.

**Request:**
```json
{
  "type": "welcome",        // "welcome" | "order_confirmation" | "order_cancellation"
  "user_id": 1,            // Required for "welcome"
  "order_id": 1            // Required for order emails
}
```

**Response:**
```json
{
  "status": "sent"
}
```

**Errors:**
- `404` - User/Order not found
- `400` - Invalid request or unknown email type
- `500` - Failed to send email

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "ok"
}
```

---

## 🔧 Development

### Build

```bash
make build
```

### Run

```bash
make run
```

### Clean

```bash
make clean
```

### Update Dependencies

```bash
make tidy
```

---

## 🐛 Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues and solutions.

**Quick fixes:**

```bash
# 1. Enable mock mode
echo "MOCK_EMAIL=true" >> .env

# 2. Check database
psql -U postgres -d stock_hub -c "SELECT id, username FROM users LIMIT 5;"

# 3. Check service logs
# (in terminal where service is running)

# 4. Test database connection
psql -U postgres -d stock_hub -c "SELECT 1;"
```

---

## 📚 Documentation

- **Troubleshooting**: `TROUBLESHOOTING.md`
- **Configuration**: `.env.example`
- **API**: This README (API Endpoints section)
- **Architecture**: See monorepo root README

---

## 🔄 Integration with stock_hub_trade

Email service fetches data from `stock_hub_trade` database:

```go
// Fetch user from database
user, err := db.QueryRow("SELECT id, username, created_at FROM users WHERE id = $1", userID)

// Use shared model
import "stock_hub_trade/pkg/model"
emailService.SendWelcomeEmail(user *model.User)
```

**Future: gRPC Integration**

Instead of direct DB access, email service will call stock_hub_trade via gRPC:

```go
// Via gRPC (future)
conn := grpc.Dial("localhost:50051")
client := proto.NewStockHubServiceClient(conn)
user, _ := client.GetUser(ctx, &proto.GetUserRequest{UserId: 1})
```

---

## 📦 Dependencies

```go
require (
	github.com/lib/pq v1.10.9          // PostgreSQL driver
	stock_hub_trade v0.0.0             // Shared models
)

replace stock_hub_trade => ../stock_hub_trade  // Local development
```

---

## 🎓 Learning

This service demonstrates:
- **Clean Architecture** - Separation of concerns
- **Shared Models** - Monorepo with go.work
- **Mock Mode** - Development without external dependencies
- **Repository Pattern** - Data access abstraction

---

## 📝 License

MIT
