# REST vs gRPC vs WebSocket: Вибір методу комунікації 🔄

Based on: [Medium Article by Alex Botha](https://medium.com/@alexbotha_18115/restful-apis-vs-grpc-choosing-the-best-communication-method-for-real-time-data-updates-9a9dfc0cc947)

---

## 📊 Порівняння

| Feature | REST | WebSocket | gRPC |
|---------|------|-----------|------|
| **Protocol** | HTTP/1.1 (зазвичай) | WS/WSS | HTTP/2 |
| **Model** | Request/Response | Bi-directional | Bi-directional Streaming |
| **Data Format** | JSON (text) | JSON/Text/Binary | Protobuf (binary) |
| **Real-time** | ❌ (polling) або ⚠️ (WebSocket) | ✅ Native | ✅ Built-in |
| **Type Safety** | ❌ Runtime | ⚠️ Manual | ✅ Strong (.proto) |
| **Learning Curve** | 🟢 Easy | 🟡 Medium | 🔴 Steep |
| **Efficiency** | 🔴 Low (JSON) | 🟡 Medium | 🟢 High (binary) |
| **Browser Support** | ✅ Native | ✅ Native | ⚠️ gRPC-Web |
| **Best For** | Simple CRUD | Real-time web apps | Microservices |

---

## 🎯 Use Cases

### REST - Simple Applications ✅

**Коли використовувати:**
- Simple CRUD operations (Create, Read, Update, Delete)
- User profile updates
- Product listings
- Reviews, comments
- Public APIs (easy to consume)

**Приклад:**
```
Client: "GET /users/123"
Server: { "id": 123, "name": "John" }

Client: "POST /orders" + { product_id: 5 }
Server: { "order_id": 789, "status": "created" }
```

**Обмеження для Real-time:**
- ❌ **Polling**: Client запитує кожну секунду → inefficient
- ⚠️ **WebSocket**: можна, але потребує додаткової настройки

---

### gRPC - Microservices & Real-time ✅

**Коли використовувати:**
- **Microservices** communication (service-to-service)
- **Real-time data** updates (trading, IoT, security cameras)
- **High performance** requirements (low latency, high throughput)
- **Cross-language** systems (Go ↔ Python ↔ Java)
- **Streaming** (audio, video, logs, metrics)

**Приклад (Home Security System):**
```
Camera (Server) → gRPC Stream → Tablet (Client)

Camera detects motion:
  → Instant push to tablet (no polling!)
  → Binary data (Protobuf) = faster
  → Bi-directional: tablet can control camera
```

**Ключові переваги:**
1. **Bi-directional Streaming** - both client and server send data continuously
2. **Protobuf** - binary format, smaller size, faster than JSON
3. **Multiplexing** - multiple requests over single connection (HTTP/2)
4. **Strongly Typed Schema** (.proto file) - compile-time validation
5. **Cross-language** - one .proto file for Go, Python, Java, C++, etc.

**Приклад .proto:**
```protobuf
message User {
  int64 id = 1;
  string username = 2;
  google.protobuf.Timestamp created_at = 3;
}

service StockHubService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc StreamOrders(StreamRequest) returns (stream StockOrder); // Streaming!
}
```

---

### WebSocket - Real-time Web Apps ⚠️

**Коли використовувати:**
- **Browser-based** real-time apps (chat, notifications)
- **Simple real-time** without microservices complexity
- When you need **native browser support**

**Приклад:**
```
Client ↔ WebSocket ↔ Server

Chat app:
  User types message → Server → Broadcast to all users
  Persistent connection (no repeated handshakes)
```

**Переваги:**
- ✅ Native browser support
- ✅ Bi-directional communication
- ✅ Simpler than gRPC for web clients

**Недоліки порівняно з gRPC:**
- ❌ **Manual management** - you handle connection, reconnection, data flow
- ❌ **No built-in type safety** - need manual validation
- ❌ **JSON overhead** - larger payloads than Protobuf
- ❌ **No multiplexing** - each WebSocket = separate connection

---

## 🏗️ Architectural Patterns

### REST: Monolith or Simple API

```
┌─────────────┐         HTTP/REST         ┌─────────────┐
│   Browser   │ ────────────────────────► │   Server    │
│  (Client)   │ ◄──────────────────────── │  (Backend)  │
└─────────────┘         JSON              └─────────────┘

Use case: E-commerce website, blog, admin panel
```

---

### WebSocket: Real-time Web

```
┌─────────────┐      WebSocket (WS)       ┌─────────────┐
│   Browser   │ ◄════════════════════════► │  Server     │
│             │    Persistent Connection   │             │
└─────────────┘                            └─────────────┘

Use case: Chat app (WhatsApp, Messenger), live notifications
```

---

### gRPC: Microservices Architecture

```
┌──────────────┐     gRPC Stream      ┌──────────────┐
│  Service A   │ ◄══════════════════► │  Service B   │
│  (Go)        │   Protobuf Binary    │  (Python)    │
└──────┬───────┘                      └──────┬───────┘
       │                                     │
       │              gRPC                   │
       ▼                                     ▼
┌──────────────┐                    ┌──────────────┐
│  Service C   │                    │  Service D   │
│  (Java)      │                    │  (C++)       │
└──────────────┘                    └──────────────┘

Use case: Stock trading platform, IoT systems, video streaming
```

---

## 🎓 Key Insights from Article

### 1. REST Limitations for Real-time

**Polling (❌ Inefficient):**
```
Every 1 second:
  Client: "Any updates?"
  Server: "No"
  Client: "Any updates?"
  Server: "No"
  Client: "Any updates?"
  Server: "Yes! Motion detected"
  
Problems:
  - Wasted bandwidth (99% "No" responses)
  - Delayed alerts (up to 1 second lag)
  - Missed events (if happens between polls)
```

**WebSocket (⚠️ Better, but manual):**
```
Client connects → Server pushes updates
  
Problems:
  - Manual connection management
  - No built-in type safety
  - JSON overhead
```

---

### 2. gRPC Advantages

#### A. Protobuf Efficiency

```
JSON (REST):
{
  "id": 123,
  "username": "john_doe",
  "created_at": "2026-02-05T10:30:00Z"
}
Size: ~80 bytes (text)

Protobuf (gRPC):
[binary data]
Size: ~20 bytes (binary)

Result: 4x smaller! = faster network transfer
```

#### B. HTTP/2 Multiplexing

```
HTTP/1.1 (REST):
  Request 1 → [wait] → Response 1
  Request 2 → [wait] → Response 2
  Request 3 → [wait] → Response 3

HTTP/2 (gRPC):
  Request 1 ──┐
  Request 2 ──┼──→ Single Connection ──→ Response 1, 2, 3 (parallel)
  Request 3 ──┘

Result: Lower latency, fewer connections
```

#### C. Strongly Typed Schema

```protobuf
// .proto file (shared between services)
message User {
  int64 id = 1;           // MUST be int64
  string username = 2;    // MUST be string
  Timestamp created_at = 3; // MUST be Timestamp
}

✅ Compile-time validation
✅ Auto-generated code (Go, Python, Java)
✅ Cross-language compatibility
❌ Can't send wrong data type (unlike JSON)
```

---

### 3. gRPC vs WebSocket

| Aspect | gRPC | WebSocket |
|--------|------|-----------|
| **Bi-directional** | ✅ Built-in | ✅ Built-in |
| **Learning Curve** | 🔴 Steep (Protobuf) | 🟡 Medium (Manual) |
| **Type Safety** | ✅ .proto file | ⚠️ Manual validation |
| **Cross-language** | ✅ .proto = translator | ⚠️ Manual serialization |
| **Browser Support** | ❌ Need gRPC-Web | ✅ Native |
| **Microservices** | ✅ Ideal | ⚠️ Possible |
| **Data Format** | Protobuf (binary) | JSON/Text (larger) |

---

## 🚀 Real-world Example: Stock Hub

### Architecture Decision

```
┌─────────────────────────────────────────────────────────┐
│  Stock Hub System                                       │
└─────────────────────────────────────────────────────────┘

Frontend (Browser):
  - WebSocket for real-time order updates (native support)
  - REST for user profile, settings (simple CRUD)

Backend Services:
  - gRPC between stock_hub ↔ email_service (efficient)
  - gRPC between stock_hub ↔ analytics_service (streaming)
  - gRPC between stock_hub ↔ matching_engine (low latency)
```

### Why Not Use Only One?

**Why not only REST?**
- ❌ Too slow for real-time order matching
- ❌ Inefficient for microservice communication (JSON overhead)
- ❌ No streaming support

**Why not only WebSocket?**
- ❌ Poor for microservices (no type safety)
- ❌ Manual connection management complexity
- ❌ Larger payloads (JSON vs Protobuf)

**Why not only gRPC?**
- ❌ Limited browser support (need gRPC-Web proxy)
- ❌ Overkill for simple CRUD operations

**✅ Best Approach: Hybrid**
- **gRPC** for backend microservices (service-to-service)
- **WebSocket** for browser real-time updates (frontend)
- **REST** for simple CRUD, public APIs

---

## 📈 Performance Comparison

### Latency (Lower is Better)

```
Simple Request (100ms):
  REST:       100ms (JSON parsing + HTTP/1.1)
  WebSocket:   80ms (JSON parsing, persistent connection)
  gRPC:        50ms (Protobuf binary + HTTP/2 multiplexing)

1000 Parallel Requests:
  REST:      5000ms (connection overhead)
  WebSocket: 3000ms (single connection)
  gRPC:      1500ms (multiplexing)
```

### Bandwidth (Lower is Better)

```
1000 User objects:
  REST (JSON):     800 KB
  WebSocket (JSON): 800 KB
  gRPC (Protobuf):  200 KB (4x smaller!)
```

---

## 🎯 Decision Matrix

### Use REST when:
- ✅ Building public APIs (easy for third-party developers)
- ✅ Simple CRUD operations (users, products, reviews)
- ✅ Need wide browser compatibility
- ✅ Team has no Protobuf experience
- ✅ Real-time not critical (profile updates, settings)

### Use WebSocket when:
- ✅ Browser-based real-time (chat, notifications, live feeds)
- ✅ Simple bi-directional communication
- ✅ Native browser support needed
- ✅ Don't need microservices complexity
- ❌ NOT for microservices (use gRPC instead)

### Use gRPC when:
- ✅ Microservices architecture (service-to-service)
- ✅ Real-time data streaming (IoT, trading, logs)
- ✅ High performance requirements (low latency)
- ✅ Cross-language services (Go ↔ Python ↔ Java)
- ✅ Type safety critical (financial, healthcare)
- ✅ Bandwidth optimization needed
- ❌ NOT for browser clients (use gRPC-Web or WebSocket)

---

## 🏁 Conclusion (from Article)

> "gRPC is a fantastic choice for **complex applications** with **multiple services** that require **real-time communication**. Its built-in features like bi-directional streaming, strong typing via .proto files, and efficient cross-platform compatibility make it a great tool for **microservices**."

> "While there is a **steep learning curve**, the **long-term benefits**, especially in terms of **scalability**, make it worthwhile."

> "For simpler applications with fewer real-time communication needs, **WebSockets** can also be a great option, but for **larger, more complex systems**, gRPC will certainly ensure you don't miss important real-time updates."

---

## 🔑 Key Takeaway

**Don't choose ONE, choose the RIGHT ONE for each use case:**

```
┌───────────────────────────────────────────────┐
│          Modern Application Stack             │
├───────────────────────────────────────────────┤
│  Frontend (Browser)                           │
│    - REST for CRUD                            │
│    - WebSocket for real-time UI updates      │
├───────────────────────────────────────────────┤
│  Backend (Microservices)                      │
│    - gRPC for service-to-service              │
│    - gRPC Streaming for real-time data        │
├───────────────────────────────────────────────┤
│  Public API                                   │
│    - REST for third-party integrations        │
└───────────────────────────────────────────────┘
```

**Example: Stock Hub**
- 🌐 **Browser → Server**: WebSocket (real-time order updates)
- 🔄 **stock_hub ↔ email_service**: gRPC (efficient, typed)
- 🔄 **stock_hub ↔ analytics_service**: gRPC Streaming (real-time metrics)
- 🔓 **Public API**: REST (easy for partners to integrate)

---

**Source**: [RESTful APIs vs gRPC by Alex Botha](https://medium.com/@alexbotha_18115/restful-apis-vs-grpc-choosing-the-best-communication-method-for-real-time-data-updates-9a9dfc0cc947)
