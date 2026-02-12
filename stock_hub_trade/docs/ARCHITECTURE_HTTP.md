# HTTP Architecture - Clean Structure 🏗️

## 📁 Структура

```
internal/
├── handler/
│   ├── http/
│   │   ├── handlers.go    ← HTTP handlers (business logic)
│   │   └── router.go      ← Routes configuration
│   └── websocket/
│       ├── handler.go     ← WebSocket handlers
│       ├── types.go       ← Request/Response DTOs
│       └── errors.go      ← Domain errors
│
├── middleware/
│   ├── manager.go         ← Middleware manager
│   ├── auth.go            ← Authentication
│   ├── logging.go         ← Request logging
│   ├── recovery.go        ← Panic recovery
│   ├── cors.go            ← CORS headers
│   └── rate_limit.go      ← Rate limiting
│
cmd/server/main.go          ← Only initialization
```

---

## 🎯 Separation of Concerns

### 1. Handlers (Business Logic)

**`internal/handler/http/handlers.go`:**
```go
type Handler struct {
    // Dependencies
}

func (h *Handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
    // Handle request
}
```

**Відповідальність:**
- Обробка HTTP requests
- Виклик services
- Формування responses

---

### 2. Router (Route Configuration)

**`internal/handler/http/router.go`:**
```go
type Router struct {
    mux        *http.ServeMux
    httpH      *Handler
    wsH        *wsHandler.Handler
    middleware *middleware.Manager
}

func (r *Router) Setup() http.Handler {
    // Configure routes with middleware
    r.mux.HandleFunc("/ws", r.middleware.Chain(
        r.wsH.HandleConnection,
        r.middleware.Logging,
        r.middleware.Recovery,
    ))
    
    return r.middleware.CORS(r.mux)
}
```

**Відповідальність:**
- Mapping URLs → Handlers
- Застосування middleware
- Налаштування маршрутів

---

### 3. Middleware (Cross-Cutting Concerns)

**`internal/middleware/`:**

#### Auth Middleware
```go
func (m *Manager) Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Validate token
        // Extract user info
        // Add to context
        next(w, r.WithContext(ctx))
    }
}
```

#### Logging Middleware
```go
func (m *Manager) Logging(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("[HTTP] %s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}
```

---

## 🔄 Request Flow

```
HTTP Request
    ↓
CORS Middleware (global)
    ↓
Logging Middleware
    ↓
Recovery Middleware
    ↓
Auth Middleware (optional)
    ↓
Handler (business logic)
    ↓
Service Layer
    ↓
Repository Layer
    ↓
Database
```

---

## 📋 Available Routes

### Public Routes

```
GET  /                → SPA index.html
GET  /assets/*        → Static files
GET  /health          → Health check
GET  /ws              → WebSocket connection
```

### Protected Routes (with Auth)

To add protected routes:

```go
// In router.go
r.mux.HandleFunc("/api/profile", r.middleware.Chain(
    r.apiH.GetProfile,
    r.middleware.Auth,        ← Add auth middleware
    r.middleware.Logging,
    r.middleware.Recovery,
))
```

---

## 🛡️ Middleware Usage

### 1. Logging

**Автоматично логує всі requests:**
```
[HTTP] GET /health 200 2.5ms 127.0.0.1:50123
[HTTP] POST /ws 101 1.2ms 127.0.0.1:50124
```

### 2. Recovery

**Ловить panics і повертає 500:**
```go
// If handler panics:
panic("something went wrong")

// Middleware recovers:
❌ PANIC: something went wrong
[stack trace...]
→ Returns: {"error":"Internal server error"}
```

### 3. CORS

**Додає CORS headers:**
```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
```

### 4. Auth

**Validates Bearer token:**
```bash
# Request:
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8082/api/profile

# Handler can access user:
userID, ok := middleware.GetUserID(r.Context())
```

### 5. Rate Limiting

**Limits requests per IP:**
```go
// In main.go:
rateLimiter := middleware.NewRateLimiter(100, time.Minute)

// In router:
r.mux.HandleFunc("/api/data", r.middleware.Chain(
    r.apiH.GetData,
    r.middleware.RateLimit(rateLimiter),  // 100 req/min
))
```

---

## ✅ Benefits

### 1. Separation of Concerns

```
✅ Handlers  → Business logic
✅ Router    → URL mapping
✅ Middleware → Cross-cutting concerns
```

### 2. Easy Testing

```go
// Test handler without router
func TestServeIndex(t *testing.T) {
    h := NewHandler()
    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    
    h.ServeIndex(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### 3. Reusable Middleware

```go
// Use same middleware for different routes
protected := r.middleware.Chain(
    handler,
    r.middleware.Auth,
    r.middleware.Logging,
    r.middleware.Recovery,
)
```

### 4. Clean Main

```go
// main.go - only initialization
func main() {
    // Create handlers
    httpH := httpHandler.NewHandler()
    wsH := wsHandler.NewHandler(...)
    
    // Create middleware
    mw := middleware.NewManager()
    
    // Setup router
    router := httpHandler.NewRouter(httpH, wsH, mw)
    
    // Start server
    http.ListenAndServe(":8082", router.Handler())
}
```

---

## 🔧 How to Add New Endpoint

### Step 1: Add Handler

**`internal/handler/http/handlers.go`:**
```go
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Get users from service
    users, err := h.userService.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Return JSON
    json.NewEncoder(w).Encode(users)
}
```

### Step 2: Add Route

**`internal/handler/http/router.go`:**
```go
r.mux.HandleFunc("/api/users", r.middleware.Chain(
    r.httpH.GetUsers,
    r.middleware.Auth,      // Protected
    r.middleware.Logging,
    r.middleware.Recovery,
))
```

### Step 3: Done!

```bash
curl -H "Authorization: Bearer TOKEN" http://localhost:8082/api/users
```

---

## 🧪 Testing

```bash
# Start server
make run

# Test health check
curl http://localhost:8082/health
# → {"status":"ok"}

# Test with logging
curl http://localhost:8082/
# → Logs: [HTTP] GET / 200 5ms ...

# Test rate limiting (if enabled)
for i in {1..200}; do curl http://localhost:8082/api/data; done
# → First 100: OK, Next 100: 429 Rate limit exceeded
```

---

## 📚 Related Docs

- [Clean Architecture](./ARCHITECTURE_LAYERS.md)
- [Testing Guide](./TEST_HELPERS_REFACTORING.md)
- [WebSocket Handler](../internal/handler/websocket/handler.go)

---

**Perfect Clean Architecture!** 🎉
