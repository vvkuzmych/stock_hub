# gRPC в Stock Hub - Стислий Опис 🚀

## 📁 Структура

```
api/proto/
├── model.proto          # Моделі: User, StockOrder, Message
├── stock_hub.proto      # RPC методи: GetUser, GetOrder, GetMessages
└── *.pb.go             # Згенеровані Go файли

internal/grpc/
├── server.go           # Імплементація gRPC сервера
└── converter.go        # Конвертери між model ↔ proto
```

## 🔧 Використання

### 1️⃣ Запуск

```bash
# Генерація proto файлів
make proto-gen

# Запуск (HTTP:8082 + gRPC:50051 одночасно)
make run-postgres
```

### 2️⃣ Конфігурація

```bash
# .env
GRPC_PORT=50051        # gRPC сервер
SERVER_PORT=8082       # HTTP/WebSocket
```

### 3️⃣ RPC Методи

```protobuf
service StockHubService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc GetOrder(GetOrderRequest) returns (GetOrderResponse);
  rpc GetOrders(GetOrdersRequest) returns (GetOrdersResponse);
  rpc GetMessages(GetMessagesRequest) returns (GetMessagesResponse);
}
```

### 4️⃣ Архітектура

```
cmd/server/main.go
  ├─> HTTP Server (8082)    → handlers + WebSocket
  └─> gRPC Server (50051)   → internal/grpc/server.go
        └─> service layer   → UserService, StockOrderService, MessageService
              └─> repository → PostgreSQL
```

### 5️⃣ Тестування

```bash
# Використовуємо grpcurl
grpcurl -plaintext \
  -d '{"user_id": 1}' \
  localhost:50051 \
  stock_hub.StockHubService/GetUser

# Або curl для HTTP
curl http://localhost:8082/...
```

## 🎯 Навіщо?

- **Мікросервіси** - інші сервіси (email_service) можуть викликати методи через gRPC
- **Типобезпека** - protobuf замість JSON
- **Швидкість** - HTTP/2, бінарний протокол
- **Контракт** - .proto файли = документація API

## 📚 Детальніше

- `GRPC_SETUP.md` - повна інструкція
- `GRPC_QUICKSTART.md` - швидкий старт
- `REST_VS_GRPC_VS_WEBSOCKET.md` - порівняння
