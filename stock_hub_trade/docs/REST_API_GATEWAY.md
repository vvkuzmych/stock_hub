# REST API Gateway - REST → gRPC 🚪

## 📐 Архітектура

```
Browser/Client
    ↓ (REST/JSON)
HTTP Server (:8082)
    ↓
API Gateway (converts REST → gRPC)
    ↓ (gRPC/Protobuf)
gRPC Server (:50051)
    ↓ (business logic)
Database
```

---

## ✅ Переваги

```
✅ REST для браузерів (JSON)
✅ gRPC для бізнес-логіки
✅ Єдине джерело правди (protobuf)
✅ Автоматична конвертація JSON ↔ Protobuf
✅ Типізація через protobuf
✅ Легко підтримувати
```

---

## 🔌 Available REST Endpoints

### Base URL
```
http://localhost:8082/api/v1/
```

---

### 1. Get User

**Request:**
```bash
GET /api/v1/users/:id
```

**Example:**
```bash
curl http://localhost:8082/api/v1/users/1
```

**Response:**
```json
{
  "id": "1",
  "username": "test1",
  "email": "test1@stockhub.com",
  "created_at": "2026-02-05T20:04:21.631883Z"
}
```

---

### 2. Get Order

**Request:**
```bash
GET /api/v1/orders/:id
```

**Example:**
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
  "status": "open",
  "created_at": "2026-02-11T..."
}
```

---

### 3. Get Orders by User

**Request:**
```bash
GET /api/v1/orders?user_id=:id&limit=:limit
```

**Example:**
```bash
curl "http://localhost:8082/api/v1/orders?user_id=1&limit=10"
```

**Response:**
```json
[
  {
    "id": "1",
    "user_id": "1",
    "username": "test1",
    "symbol": "AAPL",
    "order_type": "ORDER_TYPE_BID",
    "price": 150.5,
    "quantity": 10,
    "status": "open",
    "created_at": "..."
  },
  {
    "id": "2",
    "symbol": "TSLA",
    ...
  }
]
```

---

### 4. Get All Orders

**Request:**
```bash
GET /api/v1/orders
GET /api/v1/orders?symbol=AAPL
```

**Example:**
```bash
# All orders
curl http://localhost:8082/api/v1/orders

# Filter by symbol
curl http://localhost:8082/api/v1/orders?symbol=AAPL
```

**Response:**
```json
[
  {
    "id": "1",
    "user_id": "1",
    "username": "test1",
    "symbol": "AAPL",
    ...
  },
  {
    "id": "2",
    "user_id": "2",
    "username": "test2",
    "symbol": "TSLA",
    ...
  }
]
```

---

### 5. Get Messages

**Request:**
```bash
GET /api/v1/messages?limit=:limit
```

**Example:**
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
  },
  {
    "id": "2",
    "content": "How are you?",
    "client_id": "user456",
    "created_at": "..."
  }
]
```

---

## 🔄 How It Works

### Example Flow: Get User

**1. Client sends REST:**
```bash
curl http://localhost:8082/api/v1/users/1
```

**2. API Gateway receives:**
```go
// api_gateway.go
func (g *APIGateway) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := extractID(r.URL.Path)
    
    // Convert REST → gRPC
    resp, err := g.grpcClient.GetUser(r.Context(), &proto.GetUserRequest{
        UserId: userID,
    })
    
    // Convert Protobuf → JSON
    json.NewEncoder(w).Encode(resp.User)
}
```

**3. Calls internal gRPC:**
```go
grpcClient.GetUser(ctx, &proto.GetUserRequest{UserId: 1})
```

**4. gRPC server processes:**
```go
// internal/grpc/server.go
func (s *Server) GetUser(ctx, req) {
    user, err := s.userService.GetUserByID(req.UserId)
    return &proto.GetUserResponse{User: ToProtoUser(user)}, nil
}
```

**5. Returns protobuf → Gateway converts to JSON**

---

## 🛡️ Error Handling

Gateway автоматично конвертує gRPC errors → HTTP status:

| gRPC Code | HTTP Status | Example |
|-----------|-------------|---------|
| `OK` | 200 OK | Success |
| `InvalidArgument` | 400 Bad Request | Invalid input |
| `NotFound` | 404 Not Found | User/Order not found |
| `AlreadyExists` | 409 Conflict | Duplicate entry |
| `PermissionDenied` | 403 Forbidden | No access |
| `Unauthenticated` | 401 Unauthorized | Not logged in |
| `Unavailable` | 503 Service Unavailable | Server down |
| `DeadlineExceeded` | 504 Gateway Timeout | Timeout |

**Example Error Response:**
```json
{
  "error": "user not found",
  "status": 404
}
```

---

## 📝 Code Structure

```
internal/handler/http/
├── grpc_client.go    ← Internal gRPC client
├── api_gateway.go    ← REST → gRPC converter
├── handlers.go       ← Static files handler
└── router.go         ← Routes configuration
```

---

## 🚀 Usage Examples

### JavaScript (Fetch)

```javascript
// Get user
const response = await fetch('http://localhost:8082/api/v1/users/1');
const user = await response.json();
console.log(user.username); // "test1"

// Get orders by user
const orders = await fetch('http://localhost:8082/api/v1/orders?user_id=1')
  .then(r => r.json());
console.log(orders); // Array of orders
```

---

### Python (requests)

```python
import requests

# Get user
response = requests.get('http://localhost:8082/api/v1/users/1')
user = response.json()
print(user['username'])  # "test1"

# Get messages
messages = requests.get('http://localhost:8082/api/v1/messages?limit=10').json()
for msg in messages:
    print(msg['content'])
```

---

### curl

```bash
# Get user
curl http://localhost:8082/api/v1/users/1

# Get orders (with pretty print)
curl http://localhost:8082/api/v1/orders | jq .

# Get orders by user
curl "http://localhost:8082/api/v1/orders?user_id=1&limit=5" | jq .

# Get messages
curl "http://localhost:8082/api/v1/messages?limit=20" | jq .
```

---

## 🧪 Testing

```bash
# Start server
make run

# Test REST endpoints
curl http://localhost:8082/api/v1/users/1
curl http://localhost:8082/api/v1/orders
curl http://localhost:8082/api/v1/messages

# Test with invalid ID (should return 404)
curl http://localhost:8082/api/v1/users/99999

# Test health check
curl http://localhost:8082/health
```

---

## 💡 Benefits vs Direct REST

### Without Gateway (Direct REST):

```go
// Duplicate business logic in HTTP handler
func GetUser(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Validate
    // Query database  ← Duplicated logic
    // Format response
}
```

### With Gateway (REST → gRPC):

```go
// Single business logic in gRPC
func (s *Server) GetUser(ctx, req) {
    // All logic here (validation, DB, etc.)
}

// Gateway just converts
func (g *Gateway) GetUser(w, r) {
    resp := g.grpcClient.GetUser(...)  ← Reuses gRPC
    json.Encode(resp)
}
```

**✅ No duplicate code!**

---

## 🔧 Adding New Endpoint

### Step 1: Already in proto!

Your gRPC method is the source of truth.

### Step 2: Add gateway method

```go
// api_gateway.go
func (g *APIGateway) NewMethod(w http.ResponseWriter, r *http.Request) {
    // Parse REST params
    // Call gRPC
    resp, err := g.grpcClient.NewMethod(ctx, &proto.Request{...})
    // Return JSON
    g.jsonResponse(w, resp, 200)
}
```

### Step 3: Add route

```go
// router.go
r.mux.HandleFunc("/api/v1/new-endpoint", apiGateway.NewMethod)
```

**Done!** 🎉

---

## 🎯 Summary

```
REST Client (browser/curl)
    ↓ JSON
API Gateway
    ↓ Protobuf
gRPC Business Logic
    ↓
Database

✅ Best of both worlds!
```

---

## 📚 Related Docs

- [gRPC Setup](./KREYA_GRPC_SETUP.md)
- [HTTP Architecture](./ARCHITECTURE_HTTP.md)
- [REST vs gRPC](./REST_VS_GRPC_VS_WEBSOCKET.md)

**Тепер маєш REST і gRPC одночасно!** 🚀
