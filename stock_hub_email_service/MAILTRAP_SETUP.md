# Mailtrap Setup - Тестування Email онлайн 📧

**Mailtrap** - це sandbox для тестування email. Всі email "відправляються", але приходять в Mailtrap inbox, а не реальним користувачам.

---

## 🎯 Навіщо Mailtrap?

### ❌ Проблема з Mock Mode

```bash
# Mock mode - тільки логи в консолі
MOCK_EMAIL=true
```

**Недоліки:**
- Не бачиш як виглядає email насправді
- Не перевіриш HTML/CSS форматування
- Не тестуєш SMTP підключення

### ✅ Переваги Mailtrap

```bash
# Mailtrap - реальні email в sandbox
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
```

**Переваги:**
- ✅ **Онлайн inbox** - бачиш всі відправлені email
- ✅ **HTML preview** - перевіряєш форматування
- ✅ **SPAM тест** - перевіряєш чи не потрапить в SPAM
- ✅ **Тестування SMTP** - реальне підключення
- ✅ **Безпечно** - email НЕ йдуть реальним користувачам

---

## 🚀 Quick Setup

### 1. Реєстрація

1. Перейти на https://mailtrap.io/
2. Зареєструватись (безкоштовно!)
3. Email Sandbox → Inboxes → **My Inbox**

### 2. Отримати SMTP Credentials

В Mailtrap inbox знайдеш:

```
SMTP Settings:
  Host: sandbox.smtp.mailtrap.io
  Port: 2525 (або 25, 465, 587)
  Username: your-mailtrap-username
  Password: your-mailtrap-password
```

### 3. Налаштувати .env

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_email_service

# Edit .env
nano .env
```

**Додати:**

```bash
# Server
SERVER_PORT=8083

# Database
DATABASE_DSN=host=localhost port=5432 user=postgres password=postgres dbname=stock_hub sslmode=disable

# Mailtrap SMTP (для тестування)
MOCK_EMAIL=false
SMTP_HOST=sandbox.smtp.mailtrap.io
SMTP_PORT=2525
SMTP_USER=your-mailtrap-username
SMTP_PASSWORD=your-mailtrap-password
FROM_EMAIL=noreply@stockhub.com
FROM_NAME=Stock Hub
```

### 4. Restart Service

```bash
# Stop old service
pkill -f "email_service"

# Start with Mailtrap
make run
```

**Лог:**

```
📧 Starting Email Service
Server Port: 8083
SMTP: sandbox.smtp.mailtrap.io:2525
✅ Database connected
🌐 Email service running on http://localhost:8083
```

---

## 🧪 Тестування

### 1. Відправити Welcome Email

```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}'
```

**Очікується:**

```json
{"status":"sent"}
```

### 2. Перевірити в Mailtrap

1. Відкрити https://mailtrap.io/inboxes
2. Вибрати **My Inbox**
3. Побачиш новий email!

**Screenshot:**

```
┌─────────────────────────────────────────┐
│ From: Stock Hub <noreply@stockhub.com> │
│ To: test1@example.com                  │
│ Subject: Welcome to Stock Hub!          │
│                                         │
│ Hello test1,                            │
│                                         │
│ Welcome to Stock Hub! Your account...  │
└─────────────────────────────────────────┘
```

---

## 📊 Mailtrap Features

### Email Preview

- **HTML** - як виглядає email
- **TEXT** - text версія
- **Raw** - source код
- **Headers** - SMTP headers

### Analysis

- **SPAM Score** - чи потрапить в SPAM
- **Blacklist Check** - чи IP в blacklist
- **Validation** - правильність форматування

### Testing

- **Forward** - переслати на реальний email
- **HTML Check** - валідація HTML
- **Links Check** - перевірка посилань

---

## 🔄 Порівняння режимів

| Режим | Швидкість | Реалізм | Онлайн перегляд | SMTP тест |
|-------|-----------|---------|-----------------|-----------|
| **Mock Mode** | ⚡ Instant | ❌ Низький | ❌ Ні | ❌ Ні |
| **Mailtrap** | ⚡ Fast | ✅ Високий | ✅ Так | ✅ Так |
| **Real SMTP** | 🐌 Slow | ✅ Повний | ✅ Так | ✅ Так |

---

## 🎯 Коли що використовувати?

### 🧪 Development (Mock Mode)

```bash
MOCK_EMAIL=true
```

**Використовуй коли:**
- Швидка розробка
- Не потрібно бачити email
- Тестуєш логіку (не зовнішній вигляд)

---

### 📧 Testing (Mailtrap)

```bash
MOCK_EMAIL=false
SMTP_HOST=sandbox.smtp.mailtrap.io
```

**Використовуй коли:**
- Перевіряєш як виглядає email
- Тестуєш SMTP підключення
- Дебажиш email форматування
- Показуєш demo

---

### 🚀 Production (Real SMTP)

```bash
SMTP_HOST=smtp.gmail.com  # або SendGrid, AWS SES, etc.
```

**Використовуй коли:**
- Production environment
- Реальні користувачі

---

## 📚 Додаткові фічі Mailtrap

### 1. Multiple Inboxes

Створи окремі inboxes для:
- Development
- Staging
- QA Testing

### 2. Email Forwarding

Forward email на реальний email для перевірки.

### 3. API Access

```bash
# Get emails via API
curl -u "api-token:" https://mailtrap.io/api/v1/inboxes/123/messages
```

### 4. Webhooks

Отримуй notifications коли email приходить.

---

## 🐛 Troubleshooting

### Email не приходить в Mailtrap

**1. Check credentials**

```bash
# Verify .env
cat .env | grep SMTP
```

**2. Check Mailtrap status**

https://status.mailtrap.io/

**3. Check logs**

```bash
# In terminal where service runs
# Look for:
✅ Email sent to test@example.com
```

### SMTP Connection Error

```
Error: dial tcp: lookup sandbox.smtp.mailtrap.io: no such host
```

**Fix:**

```bash
# Check internet connection
ping sandbox.smtp.mailtrap.io

# Verify SMTP_HOST in .env
echo $SMTP_HOST
```

---

## 💡 Pro Tips

### 1. Use .env files per environment

```bash
.env.development   # Mock mode
.env.testing       # Mailtrap
.env.production    # Real SMTP
```

### 2. Add email to users on registration

```bash
# Update registration endpoint to accept email
curl -X POST http://localhost:8082/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com"}'
```

### 3. HTML Email Templates

Later: use HTML templates instead of plain text.

---

## 🔗 Links

- **Mailtrap**: https://mailtrap.io/
- **Docs**: https://help.mailtrap.io/
- **API**: https://api-docs.mailtrap.io/
- **Status**: https://status.mailtrap.io/

---

## ✅ Quick Test Checklist

- [ ] Зареєструватись на Mailtrap.io
- [ ] Скопіювати SMTP credentials
- [ ] Оновити `.env` з Mailtrap settings
- [ ] Restart email service
- [ ] Відправити test email
- [ ] Перевірити Mailtrap inbox
- [ ] 🎉 Побачити email онлайн!

---

**Готово!** Тепер можеш бачити всі email онлайн в Mailtrap! 🚀
