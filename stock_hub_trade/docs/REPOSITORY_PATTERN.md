# Repository Pattern - Відокремлення логіки від БД 🏗️

Цей документ пояснює **чому** і **як** використовувати Repository Pattern для відокремлення бізнес-логіки від доступу до бази даних.

---

## 🎯 Проблема (ДО рефакторингу)

### ❌ Service безпосередньо викликає SQL

```go
// stock_order_service.go (СТАРИЙ підхід)
type StockOrderService struct {
    db     *sql.DB      // ← Пряма залежність від БД!
    driver string
}

func (s *StockOrderService) CreateOrder(...) (*model.StockOrder, error) {
    // Бізнес-логіка + SQL запити в одному місці
    if price <= 0 {
        return nil, fmt.Errorf("invalid price")
    }
    
    // SQL запит прямо в сервісі
    query := `INSERT INTO stock_orders (...) VALUES (...)`
    err := s.db.QueryRow(query+" RETURNING id", ...).Scan(&id)
    
    return &model.StockOrder{...}, nil
}
```

### ⚠️ Проблеми:

1. **Порушення Single Responsibility Principle**
   - Service робить 2 речі: валідацію + SQL

2. **Важко тестувати**
   - Потрібна реальна БД або складний SQL mock
   - Тести повільні

3. **Складно змінити БД**
   - SQL розкидано по всьому сервісу
   - Щоб змінити PostgreSQL → MongoDB = переписати весь Service

4. **Дублювання SQL коду**
   - Однакові запити в різних сервісах

---

## ✅ Рішення: Repository Pattern

### Архітектура

```
┌─────────────────────────────────────────────────────────┐
│                    Clean Architecture                    │
└─────────────────────────────────────────────────────────┘

┌──────────────────┐
│   Handler/gRPC   │  ← Presentation Layer
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     Service      │  ← Business Logic Layer
│                  │     • Validation
│  - ValidatePrice │     • Business Rules
│  - ValidateQty   │     • Orchestration
│  - CreateOrder   │
└────────┬─────────┘
         │ uses
         ▼
┌──────────────────┐
│   Repository     │  ← Data Access Layer
│   (Interface)    │     • SQL queries
│                  │     • DB operations
│  - Create()      │     • Data mapping
│  - GetByID()     │
│  - Cancel()      │
└────────┬─────────┘
         │ implements
         ▼
┌──────────────────┐
│  PostgreSQL Repo │  ← Infrastructure Layer
│  Implementation  │
└──────────────────┘
```

---

## 📁 Нова структура файлів

```
internal/
├── repository/
│   ├── stock_order_repository.go           # Interface + Implementation
│   └── mock_stock_order_repository.go      # Mock for testing
├── service/
│   ├── stock_order_service.go              # OLD (can be removed)
│   ├── stock_order_service_v2.go           # NEW (uses Repository)
│   └── stock_order_service_v2_test.go      # Tests with mock
└── handler/
    └── grpc_handler.go                     # Uses Service
```

---

## 🔧 Реалізація

### 1. Repository Interface

```go
// internal/repository/stock_order_repository.go

// Interface - describes WHAT we need, not HOW
type StockOrderRepository interface {
    Create(order *model.StockOrder) error
    GetByID(id int64) (*model.StockOrder, error)
    GetOpenOrders(symbol string) ([]*model.StockOrder, error)
    GetAllOpenOrders() ([]*model.StockOrder, error)
    GetByUserID(userID int64, limit int) ([]*model.StockOrder, error)
    Cancel(orderID, userID int64) error
}
```

**Чому interface?**
- ✅ Можна мокувати для тестів
- ✅ Можна замінити реалізацію (PostgreSQL → MongoDB)
- ✅ Dependency Inversion Principle

---

### 2. Repository Implementation

```go
// internal/repository/stock_order_repository.go

type stockOrderRepositoryImpl struct {
    db     *sql.DB
    driver string
}

func NewStockOrderRepository(db *sql.DB, driver string) StockOrderRepository {
    return &stockOrderRepositoryImpl{db: db, driver: driver}
}

// Create handles ONLY database operations
func (r *stockOrderRepositoryImpl) Create(order *model.StockOrder) error {
    query := `INSERT INTO stock_orders (...) VALUES (...)`
    query = sqlutil.ConvertPlaceholders(query, r.driver)
    
    err := r.db.QueryRow(query+" RETURNING id", ...).Scan(&order.ID)
    if err != nil {
        return fmt.Errorf("failed to create order: %w", err)
    }
    
    return nil
}
```

**Відповідальність Repository:**
- ✅ SQL queries
- ✅ Data mapping (SQL → Model)
- ✅ Error handling (DB-specific)
- ❌ NO business logic (валідація, обчислення)

---

### 3. Service (бізнес-логіка)

```go
// internal/service/stock_order_service_v2.go

type StockOrderServiceV2 struct {
    repo repository.StockOrderRepository  // ← Uses Interface!
}

func NewStockOrderServiceV2(repo repository.StockOrderRepository) *StockOrderServiceV2 {
    return &StockOrderServiceV2{repo: repo}
}

// CreateOrder handles ONLY business logic
func (s *StockOrderServiceV2) CreateOrder(...) (*model.StockOrder, error) {
    // Business logic: Validation
    if err := s.validateOrderInput(price, quantity, symbol, orderType); err != nil {
        return nil, err
    }
    
    // Business logic: Additional rules
    // - Check trading hours
    // - Check user balance
    // - Check circuit breakers
    
    // Create order model
    order := &model.StockOrder{
        UserID:    userID,
        Symbol:    symbol,
        Price:     price,
        Quantity:  quantity,
        Status:    "open",
        CreatedAt: time.Now(),
    }
    
    // Delegate to repository
    if err := s.repo.Create(order); err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    
    return order, nil
}
```

**Відповідальність Service:**
- ✅ Business logic (валідація, правила)
- ✅ Orchestration (координація між repositories)
- ✅ Transaction management
- ❌ NO SQL queries

---

### 4. Mock Repository (для тестів)

```go
// internal/repository/mock_stock_order_repository.go

type MockStockOrderRepository struct {
    CreateFunc func(order *model.StockOrder) error
    CreateCalled bool
}

func (m *MockStockOrderRepository) Create(order *model.StockOrder) error {
    m.CreateCalled = true
    if m.CreateFunc != nil {
        return m.CreateFunc(order)
    }
    order.ID = 1  // Default: success
    return nil
}
```

---

## 🧪 Тестування

### ❌ БЕЗ Repository Pattern (складно):

```go
func TestCreateOrder(t *testing.T) {
    // Потрібна реальна БД або складний SQL mock
    db, mock, _ := sqlmock.New()
    defer db.Close()
    
    // Налаштування 10+ SQL expectations
    mock.ExpectExec("CREATE TABLE...")
    mock.ExpectQuery("SELECT...")
    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO...")
    mock.ExpectCommit()
    
    service := NewStockOrderService(db, "postgres")
    order, err := service.CreateOrder(...)
    
    // Перевірка
    if err != nil { t.Error(err) }
    
    // Верифікація всіх SQL моків
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Error(err)
    }
}
```

**Проблеми:**
- 🔴 Потрібна БД або складний SQL mock
- 🔴 Тести повільні (реальна БД)
- 🔴 Тести крихкі (зміна SQL → ламає тести)

---

### ✅ З Repository Pattern (просто):

```go
func TestCreateOrderV2_ValidatesPrice(t *testing.T) {
    // Простий mock - БД не потрібна!
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    
    // Тест бізнес-логіки
    _, err := service.CreateOrder(1, "john", "AAPL", model.OrderTypeBid, -100, 10)
    
    if err == nil {
        t.Error("Expected error for negative price")
    }
}

func TestCreateOrderV2_RepositoryError(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    mockRepo.CreateFunc = func(order *model.StockOrder) error {
        return fmt.Errorf("database connection failed")
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    _, err := service.CreateOrder(1, "john", "AAPL", model.OrderTypeBid, 150, 10)
    
    if err == nil {
        t.Error("Expected error when repository fails")
    }
    
    if !mockRepo.CreateCalled {
        t.Error("Repository should have been called")
    }
}
```

**Переваги:**
- ✅ БД не потрібна
- ✅ Швидкі тести (< 1ms)
- ✅ Легко писати
- ✅ Фокус на бізнес-логіці

---

## 📊 Порівняння

| Аспект | БЕЗ Repository | З Repository |
|--------|----------------|--------------|
| **Архітектура** | ❌ Monolith | ✅ Layered |
| **Тестування** | 🔴 Складно | 🟢 Легко |
| **Швидкість тестів** | 🔴 Повільно (100-500ms) | 🟢 Швидко (<1ms) |
| **Залежності** | ❌ Потрібна БД | ✅ Mock |
| **Зміна БД** | 🔴 Важко | 🟢 Легко |
| **Дублювання SQL** | ❌ Так | ✅ Централізовано |
| **SRP** | ❌ Порушено | ✅ Дотримано |

---

## 🎯 Переваги Repository Pattern

### 1. **Separation of Concerns** ✅

```
Service:     Business Logic (validation, rules)
Repository:  Data Access (SQL, DB operations)
```

### 2. **Легке тестування** ✅

```go
// Без реальної БД!
mockRepo := repository.NewMockStockOrderRepository()
service := NewStockOrderServiceV2(mockRepo)
```

### 3. **Гнучкість** ✅

Легко замінити PostgreSQL → MongoDB:

```go
// Тільки створити нову реалізацію
type MongoStockOrderRepository struct {
    client *mongo.Client
}

func (r *MongoStockOrderRepository) Create(order *model.StockOrder) error {
    // MongoDB specific code
}

// Service залишається без змін!
service := NewStockOrderServiceV2(NewMongoStockOrderRepository(client))
```

### 4. **Централізація SQL** ✅

Всі SQL запити в одному місці:

```
repository/stock_order_repository.go  ← ALL SQL queries here
```

### 5. **Легка підтримка** ✅

```go
// Змінити SQL запит? Тільки в Repository
// Service не потребує змін
```

---

## 🚀 Міграція з старого коду

### Крок 1: Створити Repository

```bash
# Створити interface + implementation
touch internal/repository/stock_order_repository.go
touch internal/repository/mock_stock_order_repository.go
```

### Крок 2: Створити новий Service

```go
// Старий Service залишається для backward compatibility
// stock_order_service.go

// Новий Service з Repository
// stock_order_service_v2.go
type StockOrderServiceV2 struct {
    repo repository.StockOrderRepository
}
```

### Крок 3: Оновити тести

```go
// Старі тести з SQL mock залишаються
// stock_order_service_test.go

// Нові тести з Repository mock
// stock_order_service_v2_test.go
```

### Крок 4: Поступово мігрувати handlers

```go
// Old
service := service.NewStockOrderService(db, "postgres")

// New
repo := repository.NewStockOrderRepository(db, "postgres")
service := service.NewStockOrderServiceV2(repo)
```

---

## 📈 Результати

### Тести

```bash
go test ./internal/service -v -run TestCreateOrderV2

=== RUN   TestCreateOrderV2_ValidatesPrice
--- PASS: TestCreateOrderV2_ValidatesPrice (0.00s)
=== RUN   TestCreateOrderV2_ValidatesQuantity
--- PASS: TestCreateOrderV2_ValidatesQuantity (0.00s)
=== RUN   TestCreateOrderV2_ValidatesSymbol
--- PASS: TestCreateOrderV2_ValidatesSymbol (0.00s)
=== RUN   TestCreateOrderV2_RepositoryError
--- PASS: TestCreateOrderV2_RepositoryError (0.00s)
=== RUN   TestCreateOrderV2_Success
--- PASS: TestCreateOrderV2_Success (0.00s)

PASS
ok  	stock_hub/internal/service	0.715s
```

**6 тестів, < 1 секунда, без БД!** ⚡

---

## 🎓 Висновок

**Repository Pattern** дозволяє:

✅ **Clean Architecture** - чітке розділення шарів  
✅ **Easy Testing** - моки замість реальної БД  
✅ **Maintainability** - легко змінювати SQL  
✅ **Flexibility** - легко змінювати БД (PostgreSQL → MongoDB)  
✅ **SRP** - кожен клас має одну відповідальність  
✅ **DRY** - SQL не дублюється

**Це industry standard для Go backend додатків!** 🚀

---

## 📚 Додаткові ресурси

- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
