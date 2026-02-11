# Migration Guide: internal/model → pkg/model 🔄

This document describes the migration of domain models from `internal/model` to `pkg/model` to enable sharing between microservices.

---

## 📋 Summary of Changes

### What Changed?

1. **Models moved**: `internal/model/` → `pkg/model/`
2. **Imports updated**: All files now import `stock_hub/pkg/model`
3. **New service created**: `email_service` can now import shared models

### Files Modified

#### Models Moved

- ✅ `internal/model/user.go` → `pkg/model/user.go`
- ✅ `internal/model/message.go` → `pkg/model/message.go`
- ✅ `internal/model/stock_order.go` → `pkg/model/stock_order.go`

#### Import Updates

1. `cmd/server/main.go`
   ```diff
   - import "stock_hub/internal/model"
   + import "stock_hub/pkg/model"
   ```

2. `internal/service/message_service.go`
   ```diff
   - import "stock_hub/internal/model"
   + import "stock_hub/pkg/model"
   ```

3. `internal/service/user_service.go`
   ```diff
   - import "stock_hub/internal/model"
   + import "stock_hub/pkg/model"
   ```

4. `internal/service/stock_order_service.go`
   ```diff
   - import "stock_hub/internal/model"
   + import "stock_hub/pkg/model"
   ```

5. `internal/service/stock_order_service_test.go`
   ```diff
   - import "stock_hub/internal/model"
   + import "stock_hub/pkg/model"
   ```

---

## 🎯 Why This Change?

### Before (Monolith)

```
stock_hub/
├── internal/
│   ├── model/           ← Private to stock_hub only
│   │   ├── user.go
│   │   ├── message.go
│   │   └── stock_order.go
│   └── service/
│       └── user_service.go

❌ Other services cannot import internal/model
❌ Code duplication if new service needs same models
```

### After (Microservices)

```
stock_hub/
├── pkg/
│   └── model/           ← Shared across all services
│       ├── user.go
│       ├── message.go
│       └── stock_order.go
└── internal/
    └── service/
        └── user_service.go

email_service/
└── internal/
    └── service/
        └── email_service.go
            ↓
        import "stock_hub/pkg/model"  ✅

✅ Other services can import pkg/model
✅ No code duplication
✅ Single source of truth for domain entities
```

---

## 🏗️ Clean Architecture Alignment

This change aligns with **Clean Architecture** principles:

### Layer Structure

```
┌─────────────────────────────────────────────────────┐
│  Frameworks & Drivers (PostgreSQL, SMTP, HTTP)      │
│  ┌───────────────────────────────────────────────┐  │
│  │  Interface Adapters (HTTP handlers, DB)       │  │
│  │  ┌─────────────────────────────────────────┐  │  │
│  │  │  Use Cases (Services)                    │  │  │
│  │  │  ┌───────────────────────────────────┐  │  │  │
│  │  │  │  Entities (Domain Models)         │  │  │  │
│  │  │  │                                   │  │  │  │
│  │  │  │  pkg/model/                       │  │  │  │
│  │  │  │    ├── user.go                    │  │  │  │
│  │  │  │    ├── message.go                 │  │  │  │
│  │  │  │    └── stock_order.go             │  │  │  │
│  │  │  │                                   │  │  │  │
│  │  │  └───────────────────────────────────┘  │  │  │
│  │  │                                         │  │  │
│  │  └─────────────────────────────────────────┘  │  │
│  │                                               │  │
│  └───────────────────────────────────────────────┘  │
│                                                     │
└─────────────────────────────────────────────────────┘

Dependencies point INWARD: Frameworks → Adapters → Use Cases → Entities
```

### Benefits

1. **Dependency Rule**: Inner layers (Entities) don't depend on outer layers
2. **Testability**: Domain logic can be tested without DB/frameworks
3. **Flexibility**: Easy to swap implementations (PostgreSQL → MongoDB)
4. **Reusability**: Domain models shared across microservices

---

## 🚀 Verification

### Build & Test

```bash
# Rebuild stock_hub
cd /Users/vkuzm/GolandProjects/stock_hub
go mod tidy
go build ./...
go test ./...

# Build email_service
cd /Users/vkuzm/GolandProjects/email_service
go mod tidy
go build ./...
```

### Expected Output

```
✅ All packages compile without errors
✅ All tests pass
✅ email_service can import stock_hub/pkg/model
```

---

## 🔄 Rollback (If Needed)

To rollback this change:

```bash
cd /Users/vkuzm/GolandProjects/stock_hub

# Move models back
mkdir -p internal/model
mv pkg/model/*.go internal/model/

# Update imports (reverse)
find . -name "*.go" -exec sed -i '' 's|stock_hub/pkg/model|stock_hub/internal/model|g' {} +

# Rebuild
go mod tidy
go build ./...
```

---

## 📦 Future: Publishing Shared Models

### Option 1: Separate Module (Recommended for Production)

```bash
cd /Users/vkuzm/GolandProjects/stock_hub/pkg/model
go mod init github.com/yourusername/stock-hub-models
git init
git add .
git commit -m "Initial commit"
git tag v1.0.0
git push origin main --tags
```

Then in `email_service/go.mod`:

```go
require github.com/yourusername/stock-hub-models v1.0.0
```

### Option 2: Go Workspace (Local Development)

```bash
cd /Users/vkuzm/GolandProjects
go work init
go work use ./stock_hub
go work use ./email_service
```

Then in `email_service/go.mod`:

```go
// No replace directive needed with workspace
require stock_hub v0.0.0
```

---

## 🎓 Learning Resources

- **Clean Architecture**: [Uncle Bob's Blog](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **Go Project Layout**: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- **Microservices**: [microservices.io](https://microservices.io/)
- **Domain-Driven Design**: [Martin Fowler](https://martinfowler.com/tags/domain%20driven%20design.html)

---

## ✅ Migration Completed

- ✅ Models moved to `pkg/model`
- ✅ All imports updated
- ✅ `email_service` created
- ✅ All tests passing
- ✅ Documentation updated

**Date**: 2026-02-05  
**Status**: ✅ Complete
