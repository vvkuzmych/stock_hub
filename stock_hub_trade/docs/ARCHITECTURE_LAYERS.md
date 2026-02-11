# Архітектурні шари Stock Hub 🏗️

Візуальне представлення архітектури з Repository Pattern.

---

## 🎯 До vs Після

### ❌ Старий підхід (без Repository)

```
┌──────────────────────────────────────────────────────┐
│                    Handler/gRPC                       │
│                                                       │
│  handleCreateOrder()                                  │
└────────────────────┬─────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│                     Service                           │
│  ❌ Mixed responsibilities:                           │
│     • Business logic (validation)                     │
│     • SQL queries                                     │
│     • DB operations                                   │
│                                                       │
│  CreateOrder() {                                      │
│      if price <= 0 { return error }   ← Business     │
│      query := "INSERT INTO..."        ← SQL          │
│      db.QueryRow(query).Scan(&id)     ← DB           │
│  }                                                    │
└────────────────────┬─────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│                   PostgreSQL                          │
└──────────────────────────────────────────────────────┘

Problems:
  🔴 Service залежить від конкретної БД
  🔴 Важко тестувати (потрібна реальна БД)
  🔴 SQL розкидано по всьому Service
  🔴 Порушення SRP
```

---

### ✅ Новий підхід (з Repository)

```
┌──────────────────────────────────────────────────────┐
│                    Handler/gRPC                       │
│                  (Presentation Layer)                 │
│                                                       │
│  handleCreateOrder()                                  │
│  handleGetOrder()                                     │
└────────────────────┬─────────────────────────────────┘
                     │
                     ▼
┌──────────────────────────────────────────────────────┐
│                     Service                           │
│                (Business Logic Layer)                 │
│                                                       │
│  ✅ ONLY business logic:                              │
│     • Validation (price > 0, qty > 0)                │
│     • Business rules (trading hours, limits)         │
│     • Orchestration (координація)                    │
│                                                       │
│  CreateOrder() {                                      │
│      validateInput()             ← Business          │
│      checkTradingHours()         ← Business          │
│      repo.Create(order)          ← Delegation        │
│  }                                                    │
└────────────────────┬─────────────────────────────────┘
                     │ uses interface
                     ▼
┌──────────────────────────────────────────────────────┐
│          Repository Interface (Contract)              │
│                                                       │
│  type StockOrderRepository interface {                │
│      Create(order) error                              │
│      GetByID(id) (*Order, error)                     │
│      Cancel(id, userID) error                        │
│  }                                                    │
└─────────────┬────────────────────┬───────────────────┘
              │                    │
              │ implements         │ implements
              ▼                    ▼
┌─────────────────────┐  ┌────────────────────────────┐
│  PostgreSQL Repo    │  │    Mock Repo (Testing)     │
│  (Production)       │  │                            │
│                     │  │  ✅ No DB needed!          │
│  ✅ ONLY SQL:       │  │  ✅ Fast tests (<1ms)     │
│     • INSERT        │  │  ✅ Easy to write         │
│     • SELECT        │  │                            │
│     • UPDATE        │  │  CreateFunc = func() {    │
│     • DELETE        │  │      // test logic        │
│                     │  │  }                        │
└──────────┬──────────┘  └────────────────────────────┘
           │
           ▼
┌─────────────────────┐
│    PostgreSQL       │
└─────────────────────┘

Benefits:
  ✅ Service незалежний від БД
  ✅ Легко тестувати (mock repository)
  ✅ SQL централізовано в Repository
  ✅ SRP дотримано
  ✅ Легко змінити БД (PostgreSQL → MongoDB)
```

---

## 📊 Відповідальності шарів

### 1. **Handler/gRPC Layer** (Presentation)

```go
// cmd/server/main.go або internal/handler/

Відповідальність:
  ✅ HTTP/gRPC handling
  ✅ Request parsing
  ✅ Response formatting
  ✅ Authentication/Authorization
  
  ❌ NO business logic
  ❌ NO SQL queries
```

**Приклад:**

```go
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
    var req OrderRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Call service
    order, err := service.CreateOrder(req.UserID, req.Symbol, ...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    json.NewEncoder(w).Encode(order)
}
```

---

### 2. **Service Layer** (Business Logic)

```go
// internal/service/

Відповідальність:
  ✅ Business logic (validation, rules)
  ✅ Orchestration (координація між repositories)
  ✅ Transaction management
  ✅ Domain logic
  
  ❌ NO HTTP handling
  ❌ NO SQL queries
  ❌ NO database details
```

**Приклад:**

```go
func (s *Service) CreateOrder(...) (*Order, error) {
    // Business logic
    if price <= 0 {
        return nil, ErrInvalidPrice
    }
    
    // Check trading hours
    if !s.isTradingHours() {
        return nil, ErrMarketClosed
    }
    
    // Delegate to repository
    order := &Order{...}
    if err := s.repo.Create(order); err != nil {
        return nil, err
    }
    
    return order, nil
}
```

---

### 3. **Repository Layer** (Data Access)

```go
// internal/repository/

Відповідальність:
  ✅ SQL queries
  ✅ Database operations
  ✅ Data mapping (SQL ↔ Model)
  ✅ Transaction handling
  
  ❌ NO business logic
  ❌ NO validation (beyond DB constraints)
```

**Приклад:**

```go
func (r *Repository) Create(order *Order) error {
    query := `INSERT INTO stock_orders 
              (user_id, symbol, price, quantity) 
              VALUES ($1, $2, $3, $4) 
              RETURNING id`
    
    err := r.db.QueryRow(query, 
        order.UserID, 
        order.Symbol, 
        order.Price, 
        order.Quantity,
    ).Scan(&order.ID)
    
    return err
}
```

---

### 4. **Model Layer** (Domain Entities)

```go
// pkg/model/

Відповідальність:
  ✅ Domain entities (data structures)
  ✅ Business rules (simple methods on struct)
  
  ❌ NO database code
  ❌ NO HTTP code
  ❌ NO external dependencies
```

**Приклад:**

```go
type StockOrder struct {
    ID        int64
    UserID    int64
    Symbol    string
    OrderType OrderType
    Price     float64
    Quantity  int
    Status    string
    CreatedAt time.Time
}

// Simple business method on entity
func (o *StockOrder) TotalValue() float64 {
    return o.Price * float64(o.Quantity)
}
```

---

## 🔄 Data Flow

### Створення ордера

```
1. Request arrives
   ↓
2. Handler parses JSON
   {
     "symbol": "AAPL",
     "price": 150.00,
     "quantity": 10
   }
   ↓
3. Handler calls Service
   service.CreateOrder(1, "john", "AAPL", "bid", 150.00, 10)
   ↓
4. Service validates
   ✓ price > 0
   ✓ quantity > 0
   ✓ symbol not empty
   ✓ trading hours
   ↓
5. Service creates model
   order := &model.StockOrder{
       UserID: 1,
       Symbol: "AAPL",
       Price: 150.00,
       ...
   }
   ↓
6. Service calls Repository
   repo.Create(order)
   ↓
7. Repository executes SQL
   INSERT INTO stock_orders (...) VALUES (...)
   ↓
8. Database returns ID
   RETURNING id → 123
   ↓
9. Repository sets ID
   order.ID = 123
   ↓
10. Service returns order
    return order, nil
    ↓
11. Handler formats response
    {
      "id": 123,
      "symbol": "AAPL",
      "price": 150.00,
      "status": "open"
    }
```

---

## 🧪 Тестування різних шарів

### Handler Tests (Integration)

```go
func TestHandleCreateOrder_Integration(t *testing.T) {
    // Use real service + real repository + test DB
    db := setupTestDB(t)
    repo := repository.NewStockOrderRepository(db, "postgres")
    service := service.NewStockOrderServiceV2(repo)
    
    req := httptest.NewRequest("POST", "/orders", body)
    w := httptest.NewRecorder()
    
    handleCreateOrder(w, req, service)
    
    // Assert response
}
```

### Service Tests (Unit with Mock)

```go
func TestCreateOrder_ValidatesPrice(t *testing.T) {
    // Mock repository - NO DB needed!
    mockRepo := repository.NewMockStockOrderRepository()
    service := service.NewStockOrderServiceV2(mockRepo)
    
    _, err := service.CreateOrder(1, "john", "AAPL", "bid", -100, 10)
    
    if err == nil {
        t.Error("Expected error for negative price")
    }
}
```

### Repository Tests (Unit with SQL Mock)

```go
func TestCreate_InsertsOrder(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()
    
    mock.ExpectQuery("INSERT INTO stock_orders").
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
    
    repo := repository.NewStockOrderRepository(db, "postgres")
    order := &model.StockOrder{...}
    
    err := repo.Create(order)
    
    if err != nil {
        t.Error(err)
    }
    if order.ID != 123 {
        t.Errorf("Expected ID 123, got %d", order.ID)
    }
}
```

---

## 📈 Переваги кожного шару

| Шар | Переваги | Недоліки без розділення |
|-----|----------|------------------------|
| **Handler** | HTTP logic isolated | Складно переключитися gRPC → REST |
| **Service** | Business rules in one place | Логіка розкидана |
| **Repository** | SQL centralized | SQL дублюється |
| **Model** | Type-safe domain | Mixing with DB code |

---

## 🎯 Висновок

**3-tier Architecture з Repository Pattern:**

```
Handler → Service → Repository → Database
   ↓         ↓          ↓
 HTTP    Business    Data
Logic     Logic     Access
```

**Кожен шар має чітку відповідальність!**

✅ **Clean Architecture**  
✅ **Easy to test**  
✅ **Easy to maintain**  
✅ **Easy to extend**

**This is the industry standard!** 🚀
