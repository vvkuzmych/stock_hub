# User Registration with Email - Test Guide 📧

Тепер користувачі **ОБОВ'ЯЗКОВО** мають вказувати реальний email при реєстрації!

---

## 🎯 Що Змінилось?

### До ❌
```json
{
  "type": "register",
  "data": {
    "username": "john"
  }
}
```
**Email:** `john@example.com` (фейковий!)

### Після ✅
```json
{
  "type": "register",
  "data": {
    "username": "john",
    "email": "john@real-email.com"
  }
}
```
**Email:** Реальний email користувача!

---

## 🧪 Тестування

### 1. Реєстрація Нового Користувача

```bash
# WebSocket (через wscat або websocat)
wscat -c ws://localhost:8082/ws

# Send:
{
  "type": "register",
  "data": {
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

**Успішна відповідь:**
```json
{
  "type": "registration_success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

---

### 2. Валідація Email

#### ❌ Без Email

```json
{
  "type": "register",
  "data": {
    "username": "testuser"
  }
}
```

**Помилка:**
```json
{
  "type": "error",
  "message": "Email is required"
}
```

---

#### ❌ Невалідний Email

```json
{
  "type": "register",
  "data": {
    "username": "testuser",
    "email": "invalid-email"
  }
}
```

**Помилка:**
```json
{
  "type": "error",
  "message": "invalid email format"
}
```

---

#### ❌ Email Вже Використовується

```json
{
  "type": "register",
  "data": {
    "username": "newuser",
    "email": "test@example.com"
  }
}
```

**Помилка:**
```json
{
  "type": "error",
  "message": "email already in use"
}
```

---

### 3. Перевірка в БД

```bash
psql -U postgres -d stock_hub -c "SELECT id, username, email FROM users ORDER BY id DESC LIMIT 5;"
```

**Очікується:**
```
 id |  username  |        email        
----+------------+---------------------
  5 | testuser   | test@example.com
  4 | john_doe   | john@example.com
  3 | test2      | NULL
  2 | еуіе3      | NULL
  1 | test1      | test1@example.com
```

---

### 4. Тест Email Service

**Крок 1: Зареєструвати користувача з реальним email**
```bash
# Via WebSocket
{
  "type": "register",
  "data": {
    "username": "newtester",
    "email": "newtester@example.com"
  }
}
```

**Крок 2: Отримати user_id з відповіді**
```json
{
  "type": "registration_success",
  "data": {
    "user_id": 6,  ← Use this!
    "username": "newtester",
    "email": "newtester@example.com"
  }
}
```

**Крок 3: Відправити Welcome Email**
```bash
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 6}'
```

**Крок 4: Перевірити в Mailhog**
```
http://localhost:8025
```

**Побачиш email на:** `newtester@example.com` ✅

---

## 📋 Email Validation Rules

### ✅ Валідні Email

```
user@example.com
john.doe@company.co.uk
test+tag@gmail.com
admin@localhost
support@mail-server.org
```

### ❌ Невалідні Email

```
invalid                  (без @)
@example.com            (без username)
user@                   (без domain)
user@@example.com       (подвійний @)
user @example.com       (пробіл)
```

---

## 🔧 API Endpoints

### WebSocket Registration

**Endpoint:** `ws://localhost:8082/ws`

**Message Type:** `register`

**Payload:**
```json
{
  "type": "register",
  "data": {
    "username": "string (required)",
    "email": "string (required, valid email format)"
  }
}
```

**Success Response:**
```json
{
  "type": "registration_success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

**Error Responses:**
```json
// Missing username
{"type": "error", "message": "Username is required"}

// Missing email
{"type": "error", "message": "Email is required"}

// Invalid email
{"type": "error", "message": "invalid email format"}

// Duplicate username
{"type": "error", "message": "user already exists"}

// Duplicate email
{"type": "error", "message": "email already in use"}
```

---

## 🧪 Complete Test Flow

```bash
# 1. Запустити PostgreSQL
cd stock_hub_trade
docker-compose up -d

# 2. Запустити Trading Service
make run-postgres

# 3. Запустити Mailhog
cd ../stock_hub_email_service
docker-compose up -d

# 4. Запустити Email Service
export MOCK_EMAIL=false SMTP_HOST=localhost SMTP_PORT=1025
make run

# 5. Реєстрація користувача (через WebSocket client)
# Або використати curl для HTTP endpoint якщо є

# 6. Відправити Welcome Email
curl -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": NEW_USER_ID}'

# 7. Перевірити email
open http://localhost:8025
```

---

## 📊 Database Schema

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  email TEXT,  -- Now required for new users!
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_username ON users(username);
CREATE INDEX idx_user_email ON users(email);
```

---

## 🎯 Migration Summary

### Migration 005: Add Email Field

**File:** `migrations/005_add_user_email.up.sql`

```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT;
CREATE INDEX IF NOT EXISTS idx_user_email ON users(email);
```

**Already applied!** ✅

---

## 🐛 Troubleshooting

### Email Field is NULL for Old Users

**Причина:** Існуючі користувачі були створені до міграції.

**Рішення:**
```sql
-- Update old users manually
UPDATE users SET email = username || '@example.com' WHERE email IS NULL;

-- Or delete test users
DELETE FROM users WHERE email IS NULL;
```

---

### Registration Fails with "invalid email format"

**Check:**
```bash
# Email має бути валідним:
test@example.com ✅
invalid-email ❌
```

---

## ✅ Success Checklist

- [ ] Migration 005 applied (email field exists)
- [ ] Trading service запущений
- [ ] Email service з Mailhog запущений
- [ ] Реєстрація з email працює
- [ ] Email валідація працює
- [ ] Welcome email приходить в Mailhog на правильний email

**Все працює!** 🎉

---

**Updated:** 2026-02-11  
**Status:** ✅ Email Required for Registration
