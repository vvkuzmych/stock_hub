# Context + DBModel Refactoring 🔄

Рефакторинг Repository Pattern для використання `context.Context` та публічного `DBModel` struct (як у `usual_store`).

---

## 🎯 Що змінилось?

### 1. **Додано Context у всі методи** ✅

**Було:**
```go
type StockOrderRepository interface {
    Create(order *model.StockOrder) error
    GetByID(id int64) (*model.StockOrder, error)
    Cancel(orderID, userID int64) error
}
```

**Стало:**
```go
type StockOrderRepository interface {
    Create(ctx context.Context, order *model.StockOrder) error
    GetByID(ctx context.Context, id int64) (*model.StockOrder, error)
    Cancel(ctx context.Context, orderID, userID int64) error
}
```

---

### 2. **Публічний struct DBModel** ✅

**Було (приватний):**
```go
type stockOrderRepositoryImpl struct {
    db     *sql.DB
    driver string
}

func NewStockOrderRepository(db *sql.DB, driver string) StockOrderRepository {
    return &stockOrderRepositoryImpl{
        db:     db,
        driver: driver,
    }
}
```

**Стало (публічний):**
```go
type StockOrderDBModel struct {
    DB     *sql.DB
    Driver string
}

func NewStockOrderDBModel(db *sql.DB, driver string) *StockOrderDBModel {
    return &StockOrderDBModel{
        DB:     db,
        Driver: driver,
    }
}
```

---

### 3. **Context у всіх методах Repository**

**Було:**
```go
func (r *stockOrderRepositoryImpl) Create(order *model.StockOrder) error {
    result, err := r.db.Exec(query, ...)
    // ...
}
```

**Стало:**
```go
func (m *StockOrderDBModel) Create(ctx context.Context, order *model.StockOrder) error {
    result, err := m.DB.ExecContext(ctx, query, ...)
    // ...
}
```

**Методи з Context:**
- ✅ `QueryRowContext()` замість `QueryRow()`
- ✅ `ExecContext()` замість `Exec()`
- ✅ `QueryContext()` замість `Query()`

---

### 4. **Context у Service Layer**

**Було:**
```go
func (s *StockOrderServiceV2) CreateOrder(
    userID int64,
    username, symbol string,
    orderType model.OrderType,
    price float64,
    quantity int,
) (*model.StockOrder, error) {
    // ...
    err := s.repo.Create(order)
    return order, err
}
```

**Стало:**
```go
func (s *StockOrderServiceV2) CreateOrder(
    ctx context.Context,
    userID int64,
    username, symbol string,
    orderType model.OrderType,
    price float64,
    quantity int,
) (*model.StockOrder, error) {
    // ...
    err := s.repo.Create(ctx, order)
    return order, err
}
```

---

### 5. **Context у Тестах**

**Було:**
```go
func TestCreateOrder(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    
    order, err := service.CreateOrder(1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    
    // assertions...
}
```

**Стало:**
```go
func TestCreateOrder(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    
    order, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    
    // assertions...
}
```

---

### 6. **Context у Mock Repository**

**Було:**
```go
type MockBehavior struct {
    CreateFunc func(order *model.StockOrder) error
    GetByIDFunc func(id int64) (*model.StockOrder, error)
}

func (m *MockStockOrderRepository) Create(order *model.StockOrder) error {
    m.Calls.CreateCalled = true
    if m.Behavior.CreateFunc != nil {
        return m.Behavior.CreateFunc(order)
    }
    return nil
}
```

**Стало:**
```go
type MockBehavior struct {
    CreateFunc func(ctx context.Context, order *model.StockOrder) error
    GetByIDFunc func(ctx context.Context, id int64) (*model.StockOrder, error)
}

func (m *MockStockOrderRepository) Create(ctx context.Context, order *model.StockOrder) error {
    m.Calls.CreateCalled = true
    if m.Behavior.CreateFunc != nil {
        return m.Behavior.CreateFunc(ctx, order)
    }
    return nil
}
```

---

## 🎯 Чому це важливо?

### 1. **Context для Cancellation** ⏱️

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

order, err := service.CreateOrder(ctx, ...)
// Query автоматично скасується через 5 секунд!
```

**Переваги:**
- ✅ Timeout для довгих запитів
- ✅ Cancellation для HTTP requests
- ✅ Deadlines для критичних операцій

---

### 2. **Context для Values** 🔑

```go
// Передача request ID для логування
ctx := context.WithValue(r.Context(), "requestID", uuid.New())

order, err := service.CreateOrder(ctx, ...)

// В repository можна отримати requestID:
requestID := ctx.Value("requestID")
log.Printf("Request %s: Creating order...", requestID)
```

**Use cases:**
- ✅ Request tracing (Jaeger, Zipkin)
- ✅ Authentication tokens
- ✅ User ID для audit logs
- ✅ Correlation IDs

---

### 3. **Context для Graceful Shutdown** 🛑

```go
func (s *Server) Shutdown() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Всі активні запити завершаться або будуть скасовані через 30s
    s.httpServer.Shutdown(ctx)
}
```

---

### 4. **Публічний DBModel** 🔓

**Чому публічний?**

```go
// ❌ Було (приватний):
repo := repository.NewStockOrderRepository(db, driver)
// Не можна звернутись до repo.db або repo.driver

// ✅ Стало (публічний):
repo := repository.NewStockOrderDBModel(db, driver)
// Можна repo.DB або repo.Driver (якщо потрібно)
```

**Переваги:**
- ✅ Flexibility - можна перевірити `repo.DB` у тестах
- ✅ Debugging - легше debug через публічні поля
- ✅ Extensions - легше розширювати функціонал
- ✅ **Consistency з usual_store** 🎯

---

## 📊 Порівняння: stock_hub vs usual_store

### usual_store (reference):

```go
type TokenRepository interface {
    InsertToken(ctx context.Context, token *models.Token, user models.User) error
    GetUserForToken(ctx context.Context, token string) (*models.User, error)
}

type DBModel struct {
    DB *sql.DB
}

func NewDBModel(db *sql.DB) *DBModel {
    return &DBModel{DB: db}
}

func (m *DBModel) InsertToken(ctx context.Context, token *models.Token, user models.User) error {
    stmt := `DELETE FROM tokens WHERE user_id = $1`
    _, err := m.DB.ExecContext(ctx, stmt, user.ID)
    // ...
}
```

---

### stock_hub (тепер аналогічно):

```go
type StockOrderRepository interface {
    Create(ctx context.Context, order *model.StockOrder) error
    GetByID(ctx context.Context, id int64) (*model.StockOrder, error)
}

type StockOrderDBModel struct {
    DB     *sql.DB
    Driver string
}

func NewStockOrderDBModel(db *sql.DB, driver string) *StockOrderDBModel {
    return &StockOrderDBModel{
        DB:     db,
        Driver: driver,
    }
}

func (m *StockOrderDBModel) Create(ctx context.Context, order *model.StockOrder) error {
    query := `INSERT INTO stock_orders (...) VALUES (...)`
    err := m.DB.QueryRowContext(ctx, query+" RETURNING id", ...).Scan(&order.ID)
    // ...
}
```

**Тепер обидва проекти використовують однаковий підхід!** ✅

---

## 🔄 Міграція: До vs Після

### Handler виклик (майбутнє):

**До:**
```go
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    order, err := h.service.CreateOrder(userID, username, symbol, orderType, price, qty)
    // ...
}
```

**Після:**
```go
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // ← Context з HTTP request!
    order, err := h.service.CreateOrder(ctx, userID, username, symbol, orderType, price, qty)
    // ...
}
```

**Переваги:**
- ✅ Automatic cancellation коли client закриває з'єднання
- ✅ Request timeout propagation
- ✅ Distributed tracing support

---

## 📚 Best Practices з Context

### 1. **Context як перший параметр** ✅

```go
// ✅ Good
func CreateOrder(ctx context.Context, userID int64, ...) error

// ❌ Bad
func CreateOrder(userID int64, ctx context.Context, ...) error
```

---

### 2. **Не зберігати Context у struct** ❌

```go
// ❌ Bad
type Service struct {
    ctx context.Context
}

// ✅ Good
type Service struct {
    repo Repository
}

func (s *Service) CreateOrder(ctx context.Context, ...) error {
    return s.repo.Create(ctx, order)
}
```

---

### 3. **Context у тестах: Background()** ✅

```go
func TestCreateOrder(t *testing.T) {
    ctx := context.Background()  // ← Для тестів
    
    order, err := service.CreateOrder(ctx, ...)
    // ...
}
```

---

### 4. **Context з Timeout для критичних операцій** ⏱️

```go
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // Parent context з HTTP request
    ctx := r.Context()
    
    // Додаємо timeout для DB операції
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    order, err := h.service.CreateOrder(ctx, ...)
    // ...
}
```

---

## 🧪 Приклади використання

### Приклад 1: Timeout

```go
func main() {
    db, _ := sql.Open("postgres", dsn)
    repo := repository.NewStockOrderDBModel(db, "postgres")
    service := service.NewStockOrderServiceV2(repo)
    
    // Context з timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    order, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            log.Println("Operation timed out!")
        }
    }
}
```

---

### Приклад 2: Request Tracing

```go
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // Додаємо request ID до context
    requestID := uuid.New().String()
    ctx := context.WithValue(r.Context(), "requestID", requestID)
    
    log.Printf("[%s] Creating order...", requestID)
    
    order, err := h.service.CreateOrder(ctx, ...)
    
    // В repository можна логувати з requestID:
    // log.Printf("[%s] Executing query...", ctx.Value("requestID"))
}
```

---

### Приклад 3: Cancellation

```go
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Якщо client закриє з'єднання, ctx буде cancelled
    orders, err := h.service.GetAllOpenOrders(ctx)
    
    if ctx.Err() == context.Canceled {
        log.Println("Client disconnected, query cancelled")
        return
    }
    
    json.NewEncoder(w).Encode(orders)
}
```

---

## ✅ Результати тестів

```bash
go test ./internal/service -v -run V2

=== RUN   TestCreateOrderV2_ValidatesPrice
=== RUN   TestCreateOrderV2_ValidatesPrice/Zero_price
=== RUN   TestCreateOrderV2_ValidatesPrice/Negative_price
=== RUN   TestCreateOrderV2_ValidatesPrice/Valid_price
--- PASS: TestCreateOrderV2_ValidatesPrice (0.00s)
=== RUN   TestCreateOrderV2_ValidatesQuantity
=== RUN   TestCreateOrderV2_ValidatesQuantity/Zero_quantity
=== RUN   TestCreateOrderV2_ValidatesQuantity/Negative_quantity
=== RUN   TestCreateOrderV2_ValidatesQuantity/Valid_quantity
--- PASS: TestCreateOrderV2_ValidatesQuantity (0.00s)
=== RUN   TestCreateOrderV2_ValidatesSymbol
--- PASS: TestCreateOrderV2_ValidatesSymbol (0.00s)
=== RUN   TestCreateOrderV2_ValidatesOrderType
--- PASS: TestCreateOrderV2_ValidatesOrderType (0.00s)
=== RUN   TestCreateOrderV2_RepositoryError
--- PASS: TestCreateOrderV2_RepositoryError (0.00s)
=== RUN   TestCreateOrderV2_Success
--- PASS: TestCreateOrderV2_Success (0.00s)

PASS
ok  	stock_hub/internal/service	0.862s
```

**Всі тести пройшли!** ✅

---

## 📈 Змінені файли

### 1. Repository Layer:

- ✅ `/internal/repository/stock_order_repository.go`
  - Інтерфейс з `context.Context`
  - `StockOrderDBModel` (публічний)
  - `NewStockOrderDBModel()` повертає `*StockOrderDBModel`
  - Всі методи з `Context` та `*Context()` функції

- ✅ `/internal/repository/mock_stock_order_repository.go`
  - `MockBehavior` з `context.Context`
  - Методи mock з `context.Context`

---

### 2. Service Layer:

- ✅ `/internal/service/stock_order_service_v2.go`
  - Всі методи з `context.Context` як перший параметр
  - Передача `ctx` у `repo` виклики

---

### 3. Tests:

- ✅ `/internal/service/stock_order_service_v2_test.go`
  - `ctx := context.Background()` у всіх тестах
  - Передача `ctx` у виклики сервісу
  - Mock functions з `context.Context`

---

## 🎯 Висновки

### Що отримали:

1. ✅ **Context Support** - timeout, cancellation, values
2. ✅ **DBModel Pattern** - як у `usual_store`
3. ✅ **Public Fields** - `DB`, `Driver`
4. ✅ **Consistency** - однаковий підхід в обох проектах
5. ✅ **Best Practices** - standard Go idioms
6. ✅ **Testability** - `context.Background()` у тестах
7. ✅ **Production Ready** - готово для HTTP handlers

---

### Наступні кроки:

- [ ] Оновити HTTP handlers для передачі `r.Context()`
- [ ] Додати timeout middleware
- [ ] Додати request tracing (Jaeger/Zipkin)
- [ ] Додати metrics (Prometheus)
- [ ] Додати graceful shutdown

---

## 📚 Додаткові ресурси

- [Go Context Package](https://pkg.go.dev/context)
- [Context Best Practices](https://go.dev/blog/context)
- [Database/SQL Context Methods](https://pkg.go.dev/database/sql#DB.QueryContext)
- [Effective Go - Context](https://go.dev/doc/effective_go#context)

---

**Тепер `stock_hub` repository layer відповідає стилю `usual_store`!** 🚀
