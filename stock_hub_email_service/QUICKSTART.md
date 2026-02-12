# Email Service - Quick Start ⚡

Швидкий старт для email service з Mailtrap інтеграцією.

---

## ✅ Email Service працює!

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Response: {"status":"sent"}
```

---

## 🚀 3 режими роботи

### 1️⃣ Mock Mode (Розробка)

**Найшвидший** - просто логи в консолі

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# .env
MOCK_EMAIL=true

# Run
make run
```

**Лог:**
```
📧 [MOCK] Email NOT actually sent
   To:      test1@example.com
   Subject: Welcome to Stock Hub!
   ✅ Mock email logged successfully
```

---

### 2️⃣ Mailtrap (Тестування) 🌟 RECOMMENDED

**Онлайн inbox** - бачиш email в браузері!

```bash
# 1. Sign up: https://mailtrap.io/ (безкоштовно)

# 2. Copy SMTP credentials:
#    sandbox.smtp.mailtrap.io:2525

# 3. Edit .env
nano .env

# Add:
MOCK_EMAIL=false
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your-mailtrap-username
SMTP_PASSWORD=your-mailtrap-password

# 4. Restart
make run

# 5. Send email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# 6. Check https://mailtrap.io/inboxes
#    👀 Побачиш email онлайн!
```

**Переваги:**
- ✅ Бачиш як виглядає email
- ✅ HTML preview
- ✅ SPAM score
- ✅ Безпечно (не йде реальним користувачам)

---

### 3️⃣ Real SMTP (Production)

Gmail, SendGrid, AWS SES для production.

```bash
# Gmail (потрібен App Password)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-16-char-app-password
```

---

## 🎯 Швидкий тест

### 1. Setup

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service
cp .env.example .env
nano .env  # Set MOCK_EMAIL=true
```

### 2. Start Service

```bash
make run
```

### 3. Create User (якщо немає)

```bash
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john_doe", "email": "john@example.com"}'
```

### 4. Send Email

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Expected: {"status":"sent"}
```

### 5. Check Logs

```
📧 [MOCK] Email NOT actually sent
   To:      john@example.com
   Subject: Welcome to Stock Hub!
   Body:
   Hello john_doe,
   Welcome to Stock Hub! ...
```

---

## 📧 Email Types

### Welcome Email

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
```

### Order Confirmation

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_confirmation", "order_id": 1}'
```

### Order Cancellation

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_cancellation", "order_id": 1}'
```

---

## 🔧 Troubleshooting

### "Failed to send email"

**Причина:** User/Order не існує в БД.

**Рішення:**
```bash
# Check users
psql -U postgres -d stock_hub -c "SELECT id, username, email FROM users;"

# Create test user
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "email": "test@example.com"}'
```

### PostgreSQL not running

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
docker-compose up -d
```

### Детальна допомога

See: `TROUBLESHOOTING.md`

---

## 📚 Документація

- **Mailtrap Setup**: `MAILTRAP_SETUP.md`
- **Troubleshooting**: `TROUBLESHOOTING.md`
- **Full README**: `README.md`

---

## ⚡ TL;DR - Найшвидше

```bash
# 1. Start
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service
cp .env.example .env
echo "MOCK_EMAIL=true" >> .env
make run

# 2. Test
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# 3. ✅ Working!
{"status":"sent"}
```

---

## 🌟 Next Steps

1. **Mailtrap**: Sign up → Copy credentials → Update .env
2. **HTML Templates**: Beautiful email design
3. **gRPC**: Integration with stock_hub_trade gRPC
4. **Queue**: Async email sending with RabbitMQ/Redis

---

**Готово до роботи!** 🚀
