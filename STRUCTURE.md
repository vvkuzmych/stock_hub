# Stock Hub Monorepo - Structure 🏗️

Детальна структура monorepo проекту.

---

## 📁 Root Level

```
/Users/vkuzm/GolandProjects/stock_hub/
├── .git/                       # Git repository
├── .gitignore                  # Root gitignore
├── .env.example                # Example environment variables
├── README.md                   # Main documentation
├── QUICKSTART.md               # Quick start guide
├── MONOREPO_MIGRATION.md       # Migration details
├── STRUCTURE.md                # This file
├── go.work                     # Go workspace configuration
├── go.work.sum                 # Go workspace checksums
│
├── stock_hub_trade/            # Trading Service
└── stock_hub_email_service/    # Email Service
```

---

## 📈 Trading Service (`stock_hub_trade/`)

```
stock_hub_trade/
├── .env                        # Local environment (gitignored)
├── .env.example                # Example environment
├── .gitignore                  # Service-specific gitignore
├── go.mod                      # module stock_hub_trade
├── go.sum                      # Dependencies checksums
├── Makefile                    # Build & run commands
├── README.md                   # Service documentation
├── QUICKSTART.md               # Quick start
├── CHANGELOG.md                # Version history
├── MIGRATION_GUIDE.md          # PostgreSQL migration
├── REFACTORING_SUMMARY.md      # Architecture refactoring
├── docker-compose.yml          # PostgreSQL setup
│
├── cmd/                        # Executables
│   ├── server/                 # Main HTTP/WS/gRPC server
│   │   └── main.go
│   └── migrate/                # Migration CLI
│       └── main.go
│
├── internal/                   # Private application code
│   ├── config/                 # Configuration
│   │   ├── config.go
│   │   └── config_test.go
│   ├── grpc/                   # gRPC server implementation
│   │   ├── server.go           # StockHubService implementation
│   │   └── converter.go        # Proto ↔ Model converters
│   ├── repository/             # Data access layer (Repository Pattern)
│   │   ├── stock_order_repository.go
│   │   ├── stock_order_repository_test.go
│   │   ├── stock_order_repository_test_helpers.go
│   │   └── mock_stock_order_repository.go
│   └── service/                # Business logic layer
│       ├── user_service.go
│       ├── message_service.go
│       ├── stock_order_service.go          # V1 (direct DB)
│       ├── stock_order_service_test.go
│       ├── stock_order_service_v2.go       # V2 (Repository Pattern)
│       └── stock_order_service_v2_test.go
│
├── pkg/                        # 🔥 SHARED PUBLIC CODE
│   ├── model/                  # Domain models (shared with email_service)
│   │   ├── user.go             # User entity
│   │   ├── stock_order.go      # StockOrder entity
│   │   └── message.go          # Message entity
│   ├── websocket/              # WebSocket utilities
│   │   ├── hub.go
│   │   └── hub_test.go
│   ├── migrate/                # Database migration utilities
│   │   ├── migrate.go
│   │   ├── migrate_test.go
│   │   ├── migrate_db_test.go
│   │   ├── migrate_test_helpers.go
│   │   └── README_TESTING.md
│   └── sqlutil/                # SQL utilities
│       └── placeholder.go      # Placeholder conversion (? vs $1)
│
├── api/                        # API definitions
│   └── proto/                  # gRPC Protobuf definitions
│       ├── model.proto         # Domain models in protobuf
│       ├── stock_hub.proto     # Service definition
│       ├── generate.sh         # Code generation script
│       ├── model.pb.go         # Generated Go code
│       ├── stock_hub.pb.go     # Generated Go code
│       └── stock_hub_grpc.pb.go # Generated gRPC code
│
├── migrations/                 # SQL migration files
│   ├── 001_create_messages_table.up.sql
│   ├── 001_create_messages_table.down.sql
│   ├── 002_create_users_table.up.sql
│   ├── 002_create_users_table.down.sql
│   ├── 003_create_stock_orders_table.up.sql
│   ├── 003_create_stock_orders_table.down.sql
│   ├── 004_add_order_value_constraints.up.sql
│   └── 004_add_order_value_constraints.down.sql
│
├── docs/                       # Documentation
│   ├── ARCHITECTURE_LAYERS.md
│   ├── CONTEXT_DBMODEL_REFACTORING.md
│   ├── DATA_DRIVEN_VS_FUNCTION_DRIVEN.md
│   ├── GO_INTERFACES_EXPLAINED.md
│   ├── GRPC_SETUP.md
│   ├── GRPC_QUICKSTART.md
│   ├── GRPC_USAGE_SHORT.md
│   ├── MICROSERVICES_ARCHITECTURE.md
│   ├── MOCK_EMBEDDING_EXAMPLE.md
│   ├── MOCK_STRUCTURE_EXPLAINED.md
│   ├── POSTGRES_SETUP.md
│   ├── REPOSITORY_PATTERN.md
│   ├── REST_VS_GRPC_VS_WEBSOCKET.md
│   ├── SQL_README.md
│   ├── SQL_TOP_70_QUERIES.md
│   ├── SQL_STOCK_HUB_QUERIES.md
│   ├── SQL_ADMIN_MONITORING.md
│   ├── SQL_ADVANCED_EXAMPLES.md
│   ├── SQL_OPTIMIZATION_GUIDE.md
│   ├── TABLE_DRIVEN_TESTS.md
│   ├── TABLE_DRIVEN_WITH_HELPERS.md
│   └── TEST_HELPERS_REFACTORING.md
│
├── web/                        # Frontend (React + Vite)
│   ├── src/
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   ├── App.css
│   │   └── index.css
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   └── README.md
│
├── static/                     # Built frontend assets
│   ├── index.html
│   └── assets/
│       ├── index-*.js
│       └── index-*.css
│
└── scripts/                    # Utility scripts
```

---

## 📧 Email Service (`stock_hub_email_service/`)

```
stock_hub_email_service/
├── .env                        # Local environment (gitignored)
├── .env.example                # Example environment
├── go.mod                      # module email_service
├── go.sum                      # Dependencies
├── Makefile                    # Build & run commands
├── README.md                   # Service documentation
│
├── cmd/                        # Executables
│   └── server/                 # Main HTTP server
│       └── main.go
│
└── internal/                   # Private application code
    ├── config/                 # Configuration
    │   └── config.go
    └── service/                # Business logic
        └── email_service.go    # Email sending logic
```

---

## 🔗 Inter-Service Communication

### Shared Models

```go
// Email service imports models from trading service
// stock_hub_email_service/go.mod
require stock_hub_trade v0.0.0
replace stock_hub_trade => ../stock_hub_trade

// stock_hub_email_service/internal/service/email_service.go
import "stock_hub_trade/pkg/model"

func (s *EmailService) SendWelcomeEmail(user *model.User) error {
    // ↑ Uses shared User model from stock_hub_trade
}
```

### gRPC Communication

```
Email Service → gRPC Client
    ↓
Trading Service:50051 (gRPC Server)
    ↓
Service Layer → Repository → PostgreSQL
```

---

## 🌳 Module Dependencies

```
go.work
  │
  ├── stock_hub_trade/
  │   └── depends on: gorilla/websocket, lib/pq, grpc, protobuf, sqlmock
  │
  └── stock_hub_email_service/
      └── depends on: stock_hub_trade (local replace), lib/pq
```

---

## 🚀 Ports

| Service | Protocol | Port | Description |
|---------|----------|------|-------------|
| Trading | HTTP     | 8082 | REST API & WebSocket |
| Trading | gRPC     | 50051| Inter-service communication |
| Email   | HTTP     | 50052| Trigger emails |
| PostgreSQL | TCP   | 5432 | Database |

---

## 📦 Build Artifacts (gitignored)

```
stock_hub_trade/
├── server                      # Built binary
├── stock_hub.db                # SQLite DB (legacy, not used)
└── .idea/                      # IDE files

stock_hub_email_service/
└── server                      # Built binary
```

---

## 🧪 Test Structure

### Trading Service Tests

```
internal/service/
├── stock_order_service_test.go         # V1 tests (direct DB)
└── stock_order_service_v2_test.go      # V2 tests (with mocks)

internal/repository/
├── stock_order_repository_test.go       # Repository tests (sqlmock)
└── stock_order_repository_test_helpers.go

pkg/migrate/
├── migrate_test.go                      # Unit tests (in-memory)
├── migrate_db_test.go                   # Integration tests (sqlmock)
└── migrate_test_helpers.go

pkg/websocket/
└── hub_test.go                          # WebSocket hub tests
```

### Email Service Tests

```
internal/service/
└── email_service_test.go (TODO)
```

---

## 🎯 Key Files

### Configuration

- `stock_hub_trade/.env` - Trading service config
- `stock_hub_email_service/.env` - Email service config
- `go.work` - Workspace configuration

### Documentation

- `README.md` - Main overview
- `QUICKSTART.md` - Quick start
- `MONOREPO_MIGRATION.md` - Migration guide
- `STRUCTURE.md` - This file

### Entry Points

- `stock_hub_trade/cmd/server/main.go` - Trading server
- `stock_hub_email_service/cmd/server/main.go` - Email server
- `stock_hub_trade/cmd/migrate/main.go` - Migration CLI

---

## 🔐 Ignored Files (.gitignore)

### Root

```
.env
*.db
.DS_Store
go.work.sum
```

### Services

```
bin/
*.exe
.idea/
server (binary)
```

---

## 📊 Lines of Code (approx)

| Component | Files | Lines |
|-----------|-------|-------|
| Trading Service (Go) | ~30 | ~3500 |
| Email Service (Go) | ~4 | ~200 |
| Tests | ~15 | ~2000 |
| Documentation | ~25 | ~5000 |
| Frontend (React) | ~5 | ~500 |
| **Total** | **~80** | **~11000** |

---

## 🏆 Architecture Highlights

### Clean Architecture

```
Frameworks (HTTP, gRPC, PostgreSQL)
    ↓
Interface Adapters (Handlers, Repositories, Converters)
    ↓
Use Cases (Services)
    ↓
Entities (pkg/model)
```

### Design Patterns

- **Repository Pattern** (`internal/repository/`)
- **Dependency Injection** (services via constructors)
- **Table-Driven Tests** (all test files)
- **Mock Objects** (for unit testing)
- **Builder Pattern** (test helpers)

---

## 📚 Further Reading

- [Go Workspaces](https://go.dev/doc/tutorial/workspaces)
- [Clean Architecture](stock_hub_trade/docs/MICROSERVICES_ARCHITECTURE.md)
- [Repository Pattern](stock_hub_trade/docs/REPOSITORY_PATTERN.md)
- [gRPC Setup](stock_hub_trade/docs/GRPC_SETUP.md)
