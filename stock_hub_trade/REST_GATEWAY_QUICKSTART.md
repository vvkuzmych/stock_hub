# REST API Gateway - Quick Start 🚀

## 🎯 Що Це?

**REST Gateway** - конвертує REST запити в gRPC:

```
Browser → REST/JSON → Gateway → gRPC/Protobuf → Business Logic
```

**Переваги:**
- ✅ REST для браузерів
- ✅ gRPC для бізнес-логіки
- ✅ Один source of truth (protobuf)
- ✅ Не дублюємо код

---

## 🚀 Запуск

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade
make run
```

**Сервер запустить:**
- 🌐 HTTP: `localhost:8082`
- 🚪 REST API: `localhost:8082/api/v1/`
- 🔌 gRPC: `localhost:50051`
- 📡 WebSocket: `ws://localhost:8082/ws`

---

## 📋 REST Endpoints

### Base URL
```
http://localhost:8082/api/v1/
```

---

### 1️⃣ Get User

```bash
curl http://localhost:8082/api/v1/users/1
```

**Response:**
```json
{
  "id": "1",
  "username": "test1",
  "email": "test1@stockhub.com",
  "created_at": "2026-02-05T20:04:21Z"
}
```

---

### 2️⃣ Get Order

```bash
curl http://localhost:8082/api/v1/orders/1
```

**Response:**
```json
{
  "id": "1",
  "user_id": "1",
  "username": "test1",
  "symbol": "AAPL",
  "order_type": "ORDER_TYPE_BID",
  "price": 150.5,
  "quantity": 10,
  "status": "open"
}
```

---

### 3️⃣ Get Orders by User

```bash
curl "http://localhost:8082/api/v1/orders?user_id=1&limit=10"
```

**Response:**
```json
[
  {"id": "1", "symbol": "AAPL", ...},
  {"id": "2", "symbol": "TSLA", ...}
]
```

---

### 4️⃣ Get All Orders

```bash
# All orders
curl http://localhost:8082/api/v1/orders

# Filter by symbol
curl http://localhost:8082/api/v1/orders?symbol=AAPL
```

---

### 5️⃣ Get Messages

```bash
curl "http://localhost:8082/api/v1/messages?limit=20"
```

**Response:**
```json
[
  {
    "id": "1",
    "content": "Hello!",
    "client_id": "user123",
    "created_at": "..."
  }
]
```

---

## 🧪 Test All Endpoints

```bash
# Get user
curl http://localhost:8082/api/v1/users/1 | jq .

# Get order
curl http://localhost:8082/api/v1/orders/1 | jq .

# Get orders by user
curl "http://localhost:8082/api/v1/orders?user_id=1" | jq .

# Get all orders
curl http://localhost:8082/api/v1/orders | jq .

# Filter by symbol
curl "http://localhost:8082/api/v1/orders?symbol=AAPL" | jq .

# Get messages
curl "http://localhost:8082/api/v1/messages?limit=10" | jq .

# Health check
curl http://localhost:8082/health
```

---

## 🔄 How It Works

```
1. REST Request
   ↓
   curl http://localhost:8082/api/v1/users/1

2. API Gateway (converts)
   ↓
   gRPC: GetUser(user_id=1)

3. gRPC Handler (business logic)
   ↓
   userService.GetUserByID(1)

4. Database
   ↓
   SELECT * FROM users WHERE id = 1

5. Response (converts back)
   ↓
   JSON: {"id":1,"username":"test1",...}
```

---

## 💻 JavaScript Example

```javascript
// Get user
async function getUser(userId) {
    const response = await fetch(`http://localhost:8082/api/v1/users/${userId}`);
    return await response.json();
}

// Get orders
async function getUserOrders(userId, limit = 10) {
    const url = `http://localhost:8082/api/v1/orders?user_id=${userId}&limit=${limit}`;
    const response = await fetch(url);
    return await response.json();
}

// Usage
const user = await getUser(1);
console.log(user.username); // "test1"

const orders = await getUserOrders(1, 5);
console.log(orders.length); // 5 orders
```

---

## 🐍 Python Example

```python
import requests

BASE_URL = "http://localhost:8082/api/v1"

# Get user
def get_user(user_id):
    response = requests.get(f"{BASE_URL}/users/{user_id}")
    return response.json()

# Get orders
def get_user_orders(user_id, limit=10):
    response = requests.get(
        f"{BASE_URL}/orders",
        params={"user_id": user_id, "limit": limit}
    )
    return response.json()

# Usage
user = get_user(1)
print(user['username'])  # "test1"

orders = get_user_orders(1, 5)
print(len(orders))  # 5
```

---

## ❌ Error Handling

**Gateway автоматично конвертує gRPC errors:**

```bash
# User not found
curl http://localhost:8082/api/v1/users/99999

# Response:
{
  "error": "user not found",
  "status": 404
}
```

**Error codes:**
- `400` - Bad Request (invalid input)
- `404` - Not Found (user/order doesn't exist)
- `500` - Internal Server Error
- `503` - Service Unavailable

---

## 🔍 Compare: REST vs gRPC

### Same Operation, Different Protocols:

**REST:**
```bash
curl http://localhost:8082/api/v1/users/1
```

**gRPC:**
```bash
grpcurl -plaintext \
  -d '{"user_id": 1}' \
  localhost:50051 \
  stock_hub.StockHubService/GetUser
```

**✅ Обидва використовують ОДНУ бізнес-логіку!**

---

## 📚 Full Documentation

Детальна документація: **[REST_API_GATEWAY.md](docs/REST_API_GATEWAY.md)**

---

## ✅ Summary

```
✅ REST endpoints працюють
✅ Внутрішньо використовують gRPC
✅ Не дублюємо бізнес-логіку
✅ Type-safe через protobuf
✅ Easy to use для браузерів
```

**Best of both worlds!** 🎉
