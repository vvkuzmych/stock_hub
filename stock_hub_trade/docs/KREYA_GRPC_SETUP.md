# Kreya - gRPC Testing для Stock Hub 🚀

**Kreya** - GUI клієнт для тестування gRPC (як Postman, але для gRPC).

---

## 📥 Встановлення Kreya

### macOS

```bash
brew install --cask kreya
```

Або завантажити з: https://kreya.app/downloads/

---

## 🔧 Налаштування Stock Hub gRPC

### 1. Запустити Stock Hub Trade

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade

# Запустити PostgreSQL
docker-compose up -d

# Запустити міграції
make migrate-up

# Запустити сервер (HTTP:8082 + gRPC:50051)
make run-postgres
```

**Логи:**
```
🚀 Starting gRPC server on port 50051
✅ HTTP server on port 8082
```

---

### 2. Відкрити Kreya

1. **Launch Kreya app**
2. **Create New Project** → "Stock Hub gRPC"

---

### 3. Імпортувати Proto Files

#### Варіант 1: Import Proto Files (рекомендовано)

1. **В Kreya:** 
   - File → Import → **Proto files**

2. **Вибрати файли:**
   ```
   /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto/model.proto
   /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto/stock_hub.proto
   ```

3. **Import settings:**
   - Import path: `/Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade/api/proto`
   - Package: `stock_hub`

4. **Click Import**

#### Варіант 2: Server Reflection (якщо не працює Варіант 1)

1. **В Kreya:**
   - New → **gRPC Request**
   - Method: Click "Select method"

2. **Add Server:**
   - Host: `localhost`
   - Port: `50051`
   - **Enable "Server Reflection"**

3. **Kreya автоматично завантажить методи!**

---

## 🎯 Створення Environment

### 1. В Kreya: Settings → Environments

```yaml
name: Local
variables:
  - host: localhost
  - port: 50051
  - protocol: grpc (plaintext)
```

### 2. В Kreya: Settings → Certificates

**Для локальної розробки:**
- ✅ Enable "Allow insecure connections"
- ✅ Disable "TLS/SSL"

---

## 📡 Тестування Endpoints

### GetUser - Отримати користувача

**1. Create Request:**
- Method: `stock_hub.StockHubService/GetUser`
- Host: `localhost:50051`
- Protocol: `gRPC (plaintext)`

**2. Request Body (JSON):**
```json
{
  "user_id": 1
}
```

**3. Send Request**

**4. Expected Response:**
```json
{
  "user": {
    "id": 1,
    "username": "test1",
    "email": "test1@example.com",
    "created_at": "2026-02-05T20:04:21.631883Z"
  }
}
```

---

### GetOrder - Отримати ордер

**Request:**
```json
{
  "order_id": 1
}
```

**Response:**
```json
{
  "order": {
    "id": 1,
    "user_id": 1,
    "username": "test1",
    "symbol": "AAPL",
    "order_type": "BID",
    "price": 150.50,
    "quantity": 10,
    "status": "open",
    "created_at": "..."
  }
}
```

---

### GetOrders - Всі ордери користувача

**Request:**
```json
{
  "user_id": 1,
  "limit": 10
}
```

**Response:**
```json
{
  "orders": [
    {
      "id": 1,
      "symbol": "AAPL",
      "order_type": "BID",
      ...
    },
    {
      "id": 2,
      "symbol": "TSLA",
      ...
    }
  ]
}
```

---

### GetMessages - Останні повідомлення

**Request:**
```json
{
  "limit": 20
}
```

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "user_id": 1,
      "username": "test1",
      "content": "Hello!",
      "created_at": "..."
    }
  ]
}
```

---

## 🎨 Kreya UI - Швидкий Огляд

### Структура проекту в Kreya:

```
Stock Hub gRPC/
├── Environments/
│   └── Local (localhost:50051)
│
├── Requests/
│   ├── GetUser
│   ├── GetOrder
│   ├── GetOrders
│   ├── GetOrdersByUser
│   ├── GetOpenOrders
│   ├── GetAllOpenOrders
│   └── GetMessages
│
└── Proto Files/
    ├── model.proto
    └── stock_hub.proto
```

### Request Template:

```
┌─────────────────────────────────────────┐
│ GetUser                           [Send]│
├─────────────────────────────────────────┤
│ Method: stock_hub.StockHubService/GetUser
│ Target: localhost:50051
│
│ Request Body (JSON):
│ {
│   "user_id": 1
│ }
│
│ Response:
│ {
│   "user": {
│     "id": 1,
│     "username": "test1",
│     ...
│   }
│ }
└─────────────────────────────────────────┘
```

---

## 🧪 Тестові Дані

### Create Test User (через WebSocket)

```bash
# Use websocat or browser WebSocket client
websocat ws://localhost:8082/ws

# Send:
{"type": "register", "data": {"username": "john", "email": "john@test.com"}}
```

### Create Test Order

```bash
# Via WebSocket:
{"type": "order", "data": {"symbol": "AAPL", "order_type": "bid", "price": 150.5, "quantity": 10}}
```

### Тепер тестуй в Kreya!

---

## 📚 Всі Доступні Methods

```protobuf
service StockHubService {
  // User operations
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  
  // Order operations
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc GetOrders(GetOrdersRequest) returns (GetOrdersResponse);
  rpc GetOrdersByUser(GetOrdersByUserRequest) returns (GetOrdersByUserResponse);
  rpc GetOpenOrders(GetOpenOrdersRequest) returns (GetOpenOrdersResponse);
  rpc GetAllOpenOrders(GetAllOpenOrdersRequest) returns (GetAllOpenOrdersResponse);
  
  // Message operations
  rpc GetMessages(GetMessagesRequest) returns (GetMessagesResponse);
}
```

---

## 🔍 Debug в Kreya

### Enable Request/Response Logging

1. **Settings → Logging**
2. **Enable:** "Log all requests"
3. **Enable:** "Log all responses"

### View Metadata

**В Kreya request:**
- Tab: **Metadata**
- Add custom headers:
  ```
  authorization: Bearer YOUR_TOKEN
  x-user-id: 123
  ```

### Performance Metrics

Kreya показує:
- **Request time** (ms)
- **Response size** (bytes)
- **Status code**

---

## 🆚 Kreya vs Alternatives

| Tool | Proto Import | Reflection | UI | Price |
|------|--------------|------------|-----|-------|
| **Kreya** | ✅ Yes | ✅ Yes | 🎨 Beautiful | Free |
| **Postman** | ✅ Yes | ✅ Yes | 🎨 Good | Free/Paid |
| **BloomRPC** | ✅ Yes | ✅ Yes | ⚠️ Basic | Free |
| **grpcurl** | ❌ CLI | ✅ Yes | ❌ CLI | Free |

**Переможець: Kreya** для GUI тестування! 🏆

---

## 💡 Pro Tips

### 1. Save Collections

**В Kreya:**
- File → Save Collection
- Зберігай requests для команди

### 2. Variables

**Use Kreya variables:**
```json
{
  "user_id": {{userId}}
}
```

**In Environment:**
```yaml
userId: 1
```

### 3. Chain Requests

**Scenario:**
1. GetUser → Save `user.id`
2. GetOrders → Use `{{user.id}}`

### 4. Export for CI/CD

Kreya може експортувати requests як:
- **grpcurl commands**
- **Go code**
- **Python code**

---

## 🐛 Troubleshooting

### "Failed to connect"

**Check:**
```bash
# gRPC server running?
lsof -i :50051

# If not:
cd stock_hub_trade && make run-postgres
```

### "Unknown method"

**Fix:**
```bash
# Regenerate proto files
cd stock_hub_trade
make proto-gen

# Reimport in Kreya
```

### "Invalid request"

**Check:**
- Request format matches proto definition
- All required fields present
- Correct data types (int64 vs string)

---

## 🚀 Quick Start Checklist

- [ ] Kreya встановлено
- [ ] Stock Hub Trade запущений (port 50051)
- [ ] Proto files імпортовані в Kreya
- [ ] Environment налаштований (localhost:50051)
- [ ] Test request sent успішно
- [ ] ✅ Побачив response!

---

## 📝 Example: Full Test Flow

### 1. Setup Kreya

```
New Project → "Stock Hub"
Import Proto Files → model.proto, stock_hub.proto
Add Environment → Local (localhost:50051)
```

### 2. Create Requests

```
GetUser:
  Method: stock_hub.StockHubService/GetUser
  Body: {"user_id": 1}
  
GetOrders:
  Method: stock_hub.StockHubService/GetOrdersByUser
  Body: {"user_id": 1, "limit": 10}
```

### 3. Send & Verify

```
Click "Send" → See response in Kreya
✅ Status: OK
✅ Response time: 5ms
✅ Data received
```

---

## 🔗 Links

- **Kreya Website**: https://kreya.app/
- **Kreya Docs**: https://kreya.app/docs/
- **gRPC Reflection**: https://github.com/grpc/grpc/blob/master/doc/server-reflection.md

---

## 🎓 Video Tutorial

Kreya має відео туторіали на сайті для:
- Importing proto files
- Creating requests
- Using variables
- Server reflection

---

**Готово!** Тепер можеш тестувати gRPC через красивий UI! 🎨
