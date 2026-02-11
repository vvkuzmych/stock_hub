# Mailhog Setup - Локальний Email Testing 📬

**Mailhog** - найпростіший спосіб тестувати email локально!

---

## 🌟 Переваги Mailhog

✅ **Повністю локальний** - не потрібен інтернет  
✅ **Web UI** - http://localhost:8025 (просто відкрий браузер!)  
✅ **Docker** - запускається однією командою  
✅ **Не потрібні credentials** - просто працює  
✅ **Бачиш email негайно** - без реєстрацій та API  
✅ **100% безкоштовно**

vs Ethereal:
- ❌ Ethereal потребує інтернет
- ✅ Mailhog працює офлайн
- ✅ Mailhog - просто localhost:8025

---

## 🚀 Quick Start (2 хвилини)

### 1. Запустити Mailhog

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# Запустити Docker container
docker-compose up -d

# Перевірити що працює
docker ps | grep mailhog
```

**Очікується:**
```
stock_hub_mailhog   Up   0.0.0.0:1025->1025/tcp, 0.0.0.0:8025->8025/tcp
```

---

### 2. Налаштувати Email Service

```bash
# Відкрити .env
nano .env

# Змінити на Mailhog:
MOCK_EMAIL=false
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASSWORD=
FROM_EMAIL=noreply@stockhub.com
FROM_NAME=Stock Hub
DATABASE_DSN=host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable
SERVER_PORT=8083
```

**⚠️ ВАЖЛИВО:**
- `SMTP_HOST=localhost` (не 127.0.0.1, не smtp.*)
- `SMTP_PORT=1025` (Mailhog SMTP port)
- `SMTP_USER` і `SMTP_PASSWORD` - **ПУСТІ** (не потрібні!)

---

### 3. Restart Email Service

```bash
# Зупинити старий
pkill -f "email.*server"

# Запустити новий
make run
```

**Логи:**
```
📧 Starting Email Service
Server Port: 8083
SMTP: localhost:1025
✅ Database connected
🌐 Email service running on http://localhost:8083
```

---

### 4. Відправити Test Email

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
```

**Response:**
```json
{"status":"sent"}
```

---

### 5. 👀 Побачити Email в Браузері!

**Відкрити:** http://localhost:8025

**Побачиш:**
```
┌─────────────────────────────────────────┐
│  Mailhog                                │
├─────────────────────────────────────────┤
│  📧 Welcome to Stock Hub!               │
│  From: Stock Hub <noreply@stockhub.com>│
│  To: test1@example.com                  │
│  Date: Just now                         │
│                                         │
│  [View Message]                         │
└─────────────────────────────────────────┘
```

**✅ ЦЕ ВСЕ!** Просто відкрив браузер і побачив email!

---

## 📊 Mailhog Web UI Features

### Main Page (http://localhost:8025)

- **📬 Inbox** - всі email в одному місці
- **🔍 Search** - пошук по темі, відправнику, отримувачу
- **🗑️ Delete** - видалити окремі або всі email
- **📥 Download** - завантажити .eml файл

### Email View

- **📧 Preview** - як виглядає email
- **📄 HTML** - HTML версія
- **📝 Plain Text** - текстова версія
- **🔧 Source** - raw MIME
- **📎 Attachments** - файли (якщо є)

---

## 🎯 Повний Test Flow

### Terminal 1: Start Mailhog

```bash
cd stock_hub_email_service
docker-compose up

# Логи:
[HTTP] Binding to address: 0.0.0.0:8025
[SMTP] Binding to address: 0.0.0.0:1025
```

### Terminal 2: Start Email Service

```bash
cd stock_hub_email_service
make run

# Логи:
📧 Starting Email Service
SMTP: localhost:1025
✅ Database connected
```

### Terminal 3: Send Email

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Response:
{"status":"sent"}
```

### Browser: View Email

**Open:** http://localhost:8025

**✅ Email миттєво з'явиться в inbox!**

---

## 🔧 Commands

### Start Mailhog

```bash
# Background mode
docker-compose up -d

# Foreground (see logs)
docker-compose up

# Check status
docker ps | grep mailhog
```

### Stop Mailhog

```bash
docker-compose down
```

### Restart Mailhog

```bash
docker-compose restart
```

### View Mailhog Logs

```bash
docker-compose logs -f mailhog
```

### Clear All Emails

**Web UI:** http://localhost:8025 → "Delete all messages"

**API:**
```bash
curl -X DELETE http://localhost:8025/api/v1/messages
```

---

## 📡 Mailhog API

### Get All Messages

```bash
curl http://localhost:8025/api/v2/messages
```

### Get Single Message

```bash
curl http://localhost:8025/api/v1/messages/MESSAGE_ID
```

### Delete Message

```bash
curl -X DELETE http://localhost:8025/api/v1/messages/MESSAGE_ID
```

### Delete All Messages

```bash
curl -X DELETE http://localhost:8025/api/v1/messages
```

---

## 🐛 Troubleshooting

### Port 8025 already in use

```bash
# Check what's using port
lsof -i :8025

# Kill it
kill -9 PID

# Or change Mailhog port in docker-compose.yml:
ports:
  - "8026:8025"  # Use 8026 instead
```

### Email не приходить

**Check 1: Mailhog running?**
```bash
docker ps | grep mailhog
# Should show: Up
```

**Check 2: Email service connected?**
```bash
# In email service logs:
✅ Email sent to ...
```

**Check 3: Correct SMTP settings?**
```bash
cat .env | grep SMTP
# Should be:
SMTP_HOST=localhost
SMTP_PORT=1025
```

### Web UI not loading

**Check:**
```bash
# Try different URL
http://localhost:8025
http://127.0.0.1:8025

# Check docker logs
docker-compose logs mailhog
```

---

## 🆚 Порівняння

| Feature | Mailhog | Ethereal | Mailtrap |
|---------|---------|----------|----------|
| Ціна | ✅ Free | ✅ Free | 💰 Limited |
| Setup | ✅ 1 команда | ⚠️ API call | ⚠️ Реєстрація |
| Офлайн | ✅ Так | ❌ Ні | ❌ Ні |
| Web UI | ✅ localhost:8025 | ⚠️ ethereal.email | ✅ mailtrap.io |
| Credentials | ✅ Не потрібні | ⚠️ Потрібні | ⚠️ Потрібні |
| Перегляд | ✅ Instant | ⚠️ Manual | ✅ Auto |

**Переможець: Mailhog** 🏆 - найпростіше для локальної розробки!

---

## 💡 Use Cases

### Development (Mailhog) ✅

```bash
# Швидка розробка, не потрібен інтернет
docker-compose up -d
make run
# Open http://localhost:8025
```

### Demo/Testing (Ethereal)

```bash
# Показати email клієнту онлайн
./setup_ethereal.sh
# Share https://ethereal.email/ link
```

### Production (Real SMTP)

```bash
# Gmail, SendGrid, AWS SES
SMTP_HOST=smtp.gmail.com
```

---

## 📦 Makefile Integration

Додамо команди в Makefile:

```makefile
.PHONY: mailhog-start mailhog-stop mailhog-logs mailhog-clean

mailhog-start: ## Start Mailhog
	@echo "🚀 Starting Mailhog..."
	docker-compose up -d
	@echo "✅ Mailhog running on http://localhost:8025"

mailhog-stop: ## Stop Mailhog
	@echo "⏹️  Stopping Mailhog..."
	docker-compose down

mailhog-logs: ## View Mailhog logs
	docker-compose logs -f mailhog

mailhog-clean: ## Clear all emails
	@echo "🗑️  Clearing all emails..."
	curl -X DELETE http://localhost:8025/api/v1/messages
	@echo "✅ Done"
```

**Usage:**
```bash
make mailhog-start   # Start
make mailhog-stop    # Stop
make mailhog-logs    # Logs
make mailhog-clean   # Clear emails
```

---

## ✅ Quick Test

```bash
# 1. Start Mailhog
docker-compose up -d

# 2. Update .env
cat > .env << EOF
MOCK_EMAIL=false
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USER=
SMTP_PASSWORD=
FROM_EMAIL=noreply@stockhub.com
FROM_NAME=Stock Hub
DATABASE_DSN=host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable
SERVER_PORT=8083
EOF

# 3. Start service
make run

# 4. Send email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# 5. Open browser
open http://localhost:8025

# ✅ Побачиш email!
```

---

## 🎉 Done!

**Mailhog працює!** Тепер просто:

1. `docker-compose up -d` - запустити
2. `make run` - запустити email service
3. Відправити email
4. `http://localhost:8025` - побачити!

**Найпростіше рішення!** 🚀

---

## 🔗 Links

- **Mailhog GitHub**: https://github.com/mailhog/MailHog
- **Docker Image**: https://hub.docker.com/r/mailhog/mailhog/
- **API Docs**: https://github.com/mailhog/MailHog/blob/master/docs/APIv2.md
