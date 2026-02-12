# Kreya Quick Start - 3 хвилини ⚡

## ✅ У ТЕБЕ ВЖЕ Є ГОТОВІ ПРОЕКТИ!

Kreya проекти вже налаштовані в:
- `traders1/traders1.krproj`
- `traders2/traders2.krproj`

**Просто відкрий їх!**

---

## 🚀 Швидкий старт (використай існуючий проект)

### Крок 1: Відкрити проект

**В Kreya:**
```
File → Open Project
→ Вибрати: /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/traders1/traders1.krproj
```

**✅ Всі requests вже створені:**
- GetUser
- GetOrder
- GetOrders
- GetOrdersByUser
- GetMessages

---

### Крок 2: Запустити сервер

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
make run
```

**Очікуємо:**
```
🚀 gRPC server listening on :50051
```

---

### Крок 3: Send Request

**В Kreya:**
1. **Відкрити:** `stock_hub/StockHubService/GetUser.krop`
2. **Request body вже є:** `{"userId": "275"}`
3. **Натисни Send!**

**Expected Response:**
```json
{
  "user": {
    "id": 275,
    "username": "...",
    "email": "...@...",
    "created_at": "..."
  }
}
```

---

## 🔥 Якщо Send Button Disabled

### Причина: Неправильний endpoint format!

**Перевір налаштування в `directory.krpref`:**

**✅ ПРАВИЛЬНО:**
```json
{
  "settings": [
    {
      "options": {
        "grpc": {
          "endpoint": "localhost:50051",
          "mode": "grpc"
        }
      }
    }
  ]
}
```

**❌ НЕПРАВИЛЬНО:**
- `"endpoint": "grpc://localhost:50051"` ← НЕ додавай `grpc://`
- `"endpoint": "http://localhost:50051"` ← НЕ додавай `http://`

**Kreya очікує просто:** `localhost:50051` без префіксів!

---

## 📝 Створити НОВИЙ проект (якщо потрібно)

### Крок 1: New Project

```
File → New Project
Name: Stock Hub
```

### Крок 2: Import Proto

1. **Import → Proto files**
2. **Add files:**
   ```
   /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto/model.proto
   /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto/stock_hub.proto
   ```
3. **Import path:** `/Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade`

### Крок 3: Configure Endpoint

**Створи файл `directory.krpref` в корені проекту:**

```json
{
  "settings": [
    {
      "options": {
        "grpc": {
          "endpoint": "localhost:50051",
          "mode": "grpc"
        }
      }
    }
  ]
}
```

**⚠️ БЕЗ `grpc://` або `http://` префіксів!**

### Крок 4: Create Request

1. **New Request**
2. **Method:** `stock_hub.StockHubService/GetUser`
3. **Body:** `{"user_id": 1}`
4. **Send!**

---

## 📋 Proto Files Location

```
/Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto/
├── model.proto          # User, StockOrder, Message моделі
└── stock_hub.proto      # Service methods
```

**Copy path для Kreya:**
```
/Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto
```

---

## 🧪 Test Requests

### GetUser
```json
{"user_id": 1}
```

### GetOrder
```json
{"order_id": 1}
```

### GetOrdersByUser
```json
{"user_id": 1, "limit": 10}
```

### GetMessages
```json
{"limit": 20}
```

---

## ⚙️ Connection Settings

```
Server Address: localhost:50051
Protocol: gRPC
Security: None (plaintext)
Reflection: Auto (optional)
```

---

## 🐛 Troubleshooting

### "Send button disabled" або "endpoint not supported"

**Проблема:** Неправильний формат endpoint!

**Рішення:**
1. Відкрий `directory.krpref` в корені Kreya проекту
2. Перевір що є:
   ```json
   "endpoint": "localhost:50051"
   ```
3. **НЕ повинно бути:**
   - ❌ `grpc://localhost:50051`
   - ❌ `http://localhost:50051`
   - ❌ Будь-яких префіксів!

### "Connection refused"

```bash
# Check server
lsof -i :50051

# If not running:
make run
```

### "Method not found"

**Server reflection працює?**
```bash
grpcurl -plaintext localhost:50051 list
```

**Expected:**
```
stock_hub.StockHubService
```

**Якщо ні - перезапусти сервер:**
```bash
pkill -f "cmd/server/main.go"
make run
```

---

## 📚 Full Guide

See: `docs/KREYA_GRPC_SETUP.md`

---

## 🎯 Alternative: grpcui Web UI

**Якщо Kreya не працює:**

```bash
# Install
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest

# Start (відкриється браузер з GUI)
make grpc-ui
```

**Або вручну:**
```bash
grpcui -plaintext localhost:50051
```

---

## ✅ Success Checklist

- [ ] Використав існуючий проект (`traders1`) АБО створив новий
- [ ] Endpoint = `localhost:50051` (БЕЗ префіксів!)
- [ ] gRPC server запущений `:50051`
- [ ] Send button активний
- [ ] Response отриманий

**Готово!** 🎉

---

## 🔑 Key Points

1. **Endpoint формат:** `localhost:50051` (БЕЗ `grpc://`!)
2. **Mode:** `"grpc"` в `directory.krpref`
3. **Використовуй готові проекти:** `traders1` або `traders2`
4. **Якщо не працює:** спробуй `make grpc-ui` (Web UI)

---

## 🔗 Links

- **Existing projects:** `traders1/`, `traders2/`
- **Full guide:** `docs/KREYA_GRPC_SETUP.md`
- **Alternative:** `make grpc-ui` (Web UI)
- **Server:** `localhost:50051`
