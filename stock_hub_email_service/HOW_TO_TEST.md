# Як Перевірити Email - Повна Інструкція 🧪

---

## 🚀 Швидкий Тест (3 хвилини)

```bash
# 1. Перейди в папку email service
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# 2. Створи Ethereal акаунт (автоматично)
./setup_ethereal.sh

# 3. Скопіюй credentials в .env (вручну з output вище)
nano .env

# 4. Запусти service
make run

# 5. Відправ email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# 6. Відкрий браузер і введи свій username
open https://ethereal.email/
```

---

## 📧 Як Побачити Email в Ethereal

### Спосіб 1: Web Interface (НАЙПРОСТІШЕ)

1. **Відкрити** https://ethereal.email/

2. **Прокрутити вниз** до розділу "Messages"

3. **Ввести свій username** (з .env файлу):
   ```
   Приклад: n5ghpyjy3xfvsjj5@ethereal.email
   ```

4. **Натиснути "Show messages"**

5. **Побачиш всі email!** 📬

**Screenshot (як має виглядати):**
```
┌─────────────────────────────────────────┐
│  Ethereal Email                         │
├─────────────────────────────────────────┤
│  Your Messages                          │
│                                         │
│  From: Stock Hub                        │
│  To: test1@example.com                  │
│  Subject: Welcome to Stock Hub!         │
│  Date: 2026-02-11 12:51                 │
│                                         │
│  [View Message]                         │
└─────────────────────────────────────────┘
```

---

### Спосіб 2: Direct Message URL

Якщо маєш messageId, можна відкрити напряму:

```
https://ethereal.email/message/MESSAGE_ID
```

**Як отримати messageId?** - Оновимо service щоб повертав його!

---

## 🔍 Перевірка Чи Email Відправився

### Метод 1: Логи Service

**В терміналі де запущено `make run`:**

```bash
# ✅ Успішно:
2026/02/11 12:51:50 ✅ Email sent to test1@example.com

# ❌ Помилка:
2026/02/11 12:51:50 Failed to send email: dial tcp: connection refused
```

---

### Метод 2: cURL Response

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# ✅ Успішно:
{"status":"sent"}

# ❌ Помилка:
Failed to send email
# або
User not found
```

---

### Метод 3: Health Check

```bash
# Перевірити що service працює
curl http://localhost:8083/health

# Очікується:
{"status":"ok"}
```

---

## 🧪 Повний Test Flow

### Крок 1: Підготовка

```bash
# Перевір що PostgreSQL запущений
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
docker-compose ps

# Якщо НЕ запущений:
docker-compose up -d

# Перевір що є користувачі
psql -U postgres -d stock_hub -c "SELECT id, username, email FROM users LIMIT 3;"
```

**Очікується:**
```
 id | username |       email       
----+----------+-------------------
  1 | test1    | test1@example.com
  2 | test2    | test2@example.com
```

**Якщо немає користувачів:**
```bash
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "email": "test@example.com"}'
```

---

### Крок 2: Налаштування Email Service

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# Створити Ethereal акаунт
./setup_ethereal.sh

# Результат:
# ✅ Ethereal Email account created!
# Username: abc123@ethereal.email
# Password: xyz789
# Inbox URL: https://ethereal.email
```

**Скопіювати credentials:**
```bash
nano .env

# Додати:
MOCK_EMAIL=false
SMTP_HOST=smtp.ethereal.email
SMTP_PORT=587
SMTP_USER=abc123@ethereal.email
SMTP_PASSWORD=xyz789
FROM_EMAIL=noreply@stockhub.com
FROM_NAME=Stock Hub
DATABASE_DSN=host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable
SERVER_PORT=8083
```

---

### Крок 3: Запуск Service

```bash
# Зупинити старий service (якщо працює)
pkill -f "email"

# Запустити новий
make run
```

**Очікується в логах:**
```
📧 Starting Email Service
Server Port: 8083
SMTP: smtp.ethereal.email:587
✅ Database connected
🌐 Email service running on http://localhost:8083
```

**Якщо бачиш:**
```
⚠️  MOCK MODE ENABLED
```
**→ Перевір що `MOCK_EMAIL=false` в .env!**

---

### Крок 4: Відправка Email

```bash
# Welcome email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Очікується:
{"status":"sent"}
```

**В логах service має з'явитись:**
```
2026/02/11 12:51:50 ✅ Email sent to test1@example.com
```

---

### Крок 5: Перегляд Email

#### A. Web Interface

1. **Відкрити:** https://ethereal.email/
2. **Scroll down** до "Messages"
3. **Ввести username:** `abc123@ethereal.email` (твій з .env)
4. **Натиснути "Show messages"**
5. **✅ Побачиш email!**

#### B. Alternative - Nodemailer Dashboard

Спробуй:
- https://ethereal.email/messages
- https://ethereal.email/login

---

## ❌ Troubleshooting

### Problem 1: "Failed to send email"

**Причина:** User не існує в БД

**Рішення:**
```bash
# Перевірити
psql -U postgres -d stock_hub -c "SELECT * FROM users WHERE id = 1;"

# Якщо немає - створити
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test1", "email": "test1@example.com"}'
```

---

### Problem 2: "Connection refused"

**Причина:** Service не запущений

**Рішення:**
```bash
# Перевірити чи працює
lsof -i :8083

# Якщо НЕ працює - запустити
cd stock_hub_email_service
make run
```

---

### Problem 3: Email не з'являється в Ethereal

**Причини:**
1. SMTP credentials невірні
2. MOCK_EMAIL=true (не відправляє реально)
3. Email відправився, але треба чекати

**Рішення:**
```bash
# 1. Перевірити .env
cat .env | grep SMTP

# 2. Перевірити що NOT mock mode
cat .env | grep MOCK_EMAIL
# Має бути: MOCK_EMAIL=false

# 3. Check service logs
# В терміналі де запущено service - має бути:
✅ Email sent to ...
```

---

### Problem 4: Ethereal показує "No messages"

**Можливі причини:**
1. Невірний username
2. Email відправлений на інший акаунт
3. Email expired (24h TTL)

**Рішення:**
```bash
# Подивитись який username в .env
cat .env | grep SMTP_USER

# Використати ЦЕЙ username в Ethereal web
```

---

## 🎯 Quick Commands Cheat Sheet

```bash
# 1. Setup
cd stock_hub_email_service
./setup_ethereal.sh
nano .env  # Copy credentials

# 2. Start
make run

# 3. Test
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# 4. View
# Browser → https://ethereal.email/
# Enter username from .env

# 5. Debug
lsof -i :8083                    # Service running?
cat .env | grep SMTP             # Credentials OK?
psql -U postgres -d stock_hub -c "SELECT * FROM users;"  # Users exist?
```

---

## 📱 Test All Email Types

```bash
# Welcome email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'

# Order confirmation (потрібен order_id)
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_confirmation", "order_id": 1}'

# Order cancellation
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "order_cancellation", "order_id": 1}'
```

---

## ✅ Expected Results

### In Terminal (service logs):
```
2026/02/11 12:51:50 ✅ Email sent to test1@example.com
```

### In cURL response:
```json
{"status":"sent"}
```

### In Browser (https://ethereal.email/):
```
From: Stock Hub <noreply@stockhub.com>
To: test1@example.com
Subject: Welcome to Stock Hub!

Hello test1,

Welcome to Stock Hub! Your account has been created successfully.
...
```

---

## 🎉 Success Checklist

- [ ] Ethereal акаунт створено (`./setup_ethereal.sh`)
- [ ] Credentials в `.env` (не в git!)
- [ ] Service запущений (`make run`)
- [ ] PostgreSQL працює
- [ ] User існує в БД
- [ ] Email відправлений (`{"status":"sent"}`)
- [ ] Email видно в https://ethereal.email/

**Якщо всі ✅ - працює!** 🚀

---

**Last Updated:** 2026-02-11  
**Status:** Tested & Working
