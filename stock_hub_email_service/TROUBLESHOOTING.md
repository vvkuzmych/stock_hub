# Email Service - Troubleshooting 🔧

Common issues and how to fix them.

---

## ❌ Error: "Failed to send email"

### Possible Causes

#### 1. User/Order not found in database

**Symptom:**
```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Response: User not found (404) або Failed to send email (500)
```

**Check:**
```bash
# Перевірити чи є user з ID=1
psql -U postgres -d stock_hub -c "SELECT * FROM users WHERE id = 1;"
```

**Fix:**
```bash
# Створити тестового користувача
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade

# Через HTTP API
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test_user"}'
```

---

#### 2. SMTP not configured

**Symptom:**
```
Failed to send welcome email: failed to send email: dial tcp: lookup smtp.gmail.com: no such host
```

**Check:**
```bash
# Перевірити .env файл
cat stock_hub_email_service/.env

# Має бути:
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

**Fix для Development:**

**Варіант 1: Mock Mode (рекомендовано для розробки)**

```bash
# В .env додати:
MOCK_EMAIL=true
```

Оновлений код буде просто логувати email замість відправки.

**Варіант 2: Реальний Gmail**

1. Увійти в Google Account
2. Enable 2-Factor Authentication
3. Generate App Password:
   - https://myaccount.google.com/apppasswords
4. Використати App Password в `.env`:

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
FROM_EMAIL=your-email@gmail.com
```

**Варіант 3: Mailtrap (для тестування)**

```bash
# Зареєструватись на https://mailtrap.io (безкоштовно)
SMTP_HOST=smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your-mailtrap-user
SMTP_PASSWORD=your-mailtrap-password
FROM_EMAIL=test@stockhub.com
```

---

#### 3. Database connection failed

**Symptom:**
```
Failed to connect to database: dial tcp [::1]:5432: connect: connection refused
```

**Check:**
```bash
# Перевірити чи запущений PostgreSQL
docker ps | grep postgres

# Або перевірити локальний
lsof -i :5432
```

**Fix:**
```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
docker-compose up -d
```

---

## 🔍 Debug Mode

### Enable detailed logging

```bash
# В main.go додати більше логів
# Або запустити з debug:
cd stock_hub_email_service
go run ./cmd/server/ 2>&1 | tee email_service.log
```

### Test database connection

```bash
# Test query
psql -U postgres -d stock_hub -c "SELECT id, username FROM users LIMIT 5;"
```

### Test SMTP connection

```go
// Створити test file: test_smtp.go
package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func main() {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", user, pass, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	message := []byte("From: test@test.com\r\n" +
		"To: test@test.com\r\n" +
		"Subject: Test\r\n" +
		"\r\n" +
		"This is a test.\r\n")

	err := smtp.SendMail(addr, auth, user, []string{"test@test.com"}, message)
	if err != nil {
		fmt.Printf("❌ SMTP Error: %v\n", err)
		return
	}
	fmt.Println("✅ SMTP OK")
}
```

```bash
# Run test
go run test_smtp.go
```

---

## 🧪 Testing Endpoints

### 1. Health Check

```bash
curl http://localhost:8083/health
# Expected: {"status":"ok"}
```

### 2. Welcome Email (requires user_id=1)

```bash
# First, create user
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe"}'

# Then send email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Expected: {"status":"sent"}
```

### 3. Order Confirmation (requires order_id)

```bash
# First, create order via WebSocket or stock_hub_trade API
# Then:
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_confirmation", "order_id": 1}'
```

---

## 📊 Logs Location

```bash
# Email service logs
docker logs <email-service-container>

# Or if running directly:
# Output in terminal where you ran `make run`
```

---

## 🚀 Quick Fix Steps

### Step 1: Check PostgreSQL

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
docker-compose ps
# Має бути running
```

### Step 2: Check Database

```bash
psql -U postgres -d stock_hub -c "SELECT COUNT(*) FROM users;"
# Має повернути число > 0
```

### Step 3: Enable Mock Mode

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# Create .env if not exists
cp .env.example .env

# Edit .env
echo "MOCK_EMAIL=true" >> .env

# Restart service
# Ctrl+C у терміналі де запущено
make run
```

### Step 4: Test Again

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
```

---

## 🔄 Start Fresh

```bash
# 1. Stop all services
docker-compose down

# 2. Start PostgreSQL
cd stock_hub_trade
docker-compose up -d

# 3. Run migrations
make migrate-up

# 4. Create test user
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test"}'

# 5. Enable mock mode
cd ../stock_hub_email_service
cp .env.example .env
echo "MOCK_EMAIL=true" >> .env

# 6. Start email service
make run

# 7. Test
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
```

---

## 📞 Still not working?

1. Check logs in terminal
2. Verify DATABASE_DSN in .env
3. Verify user exists: `SELECT * FROM users;`
4. Try mock mode: `MOCK_EMAIL=true`
5. Check SMTP credentials if using real SMTP

---

## ✅ Expected Behavior

### With MOCK_EMAIL=true

```bash
# Terminal output:
[MOCK] Would send email to: test@example.com
[MOCK] Subject: Welcome to Stock Hub!
[MOCK] Body:
Hello test,

Welcome to Stock Hub! Your account has been created successfully.
...
```

### With real SMTP

```bash
# Terminal output:
✅ Email sent to test@example.com

# Email received in inbox
```
