# Go Interfaces: Чому interface визначається один раз 🎯

Пояснення для питання: "Чому я не бачу `type StockOrderRepository interface` в `mock_stock_order_repository.go`?"

---

## 🔑 Відповідь

**Interface визначається ОДИН РАЗ**, а потім різні типи його **неявно реалізують (implicitly)**.

```
stock_order_repository.go:        mock_stock_order_repository.go:

┌──────────────────────────┐     ┌──────────────────────────┐
│ type Repository          │     │ type Mock struct {       │
│   interface {            │     │   CreateFunc func(...)   │
│                          │     │ }                        │
│   Create(...) error      │     │                          │
│   GetByID(...) error     │     │ func (m *Mock) Create    │
│   Cancel(...) error      │     │ func (m *Mock) GetByID   │
│ }                        │     │ func (m *Mock) Cancel    │
│                          │     │                          │
│ ← Interface (contract)   │     │ ← Implementation         │
└──────────────────────────┘     └──────────────────────────┘
                                           │
                                           │ implements
                                           │ (implicitly!)
                                           ▼
                              ✅ Mock реалізує Repository
```

---

## 📁 Структура файлів

### Файл 1: Interface (контракт)

```go
// internal/repository/stock_order_repository.go

package repository

// Interface визначається ТУТ (один раз!)
type StockOrderRepository interface {
    Create(order *model.StockOrder) error
    GetByID(id int64) (*model.StockOrder, error)
    Cancel(orderID, userID int64) error
}

// PostgreSQL implementation
type stockOrderRepositoryImpl struct {
    db *sql.DB
}

// ✅ Compile-time check
var _ StockOrderRepository = (*stockOrderRepositoryImpl)(nil)

func (r *stockOrderRepositoryImpl) Create(order *model.StockOrder) error {
    // PostgreSQL specific code
}
```

---

### Файл 2: Mock (реалізація для тестів)

```go
// internal/repository/mock_stock_order_repository.go

package repository

// Interface НЕ повторюється! Він вже визначений в stock_order_repository.go
// Mock просто реалізує методи цього interface

// Mock implementation
type MockStockOrderRepository struct {
    CreateFunc func(order *model.StockOrder) error
}

// ✅ Compile-time check що Mock реалізує interface
var _ StockOrderRepository = (*MockStockOrderRepository)(nil)

func (m *MockStockOrderRepository) Create(order *model.StockOrder) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(order)
    }
    return nil
}
```

---

## 🎓 Go Interfaces - Implicit Implementation

### Відмінність від інших мов

#### Java/C# (Explicit):

```java
// Java - EXPLICIT implementation
interface Repository {
    void create(Order order);
}

class PostgresRepo implements Repository {  // ← explicit "implements"
    public void create(Order order) { ... }
}

class MockRepo implements Repository {      // ← explicit "implements"
    public void create(Order order) { ... }
}
```

#### Go (Implicit):

```go
// Go - IMPLICIT implementation
type Repository interface {
    Create(order *Order) error
}

// PostgreSQL - автоматично реалізує interface якщо має всі методи
type PostgresRepo struct { db *sql.DB }
func (r *PostgresRepo) Create(order *Order) error { ... }  // ← NO "implements" keyword!

// Mock - автоматично реалізує interface якщо має всі методи
type MockRepo struct { CreateFunc func(...) error }
func (m *MockRepo) Create(order *Order) error { ... }     // ← NO "implements" keyword!
```

**У Go немає ключового слова `implements`!**

---

## ✅ Compile-time Check

Щоб **гарантувати** що тип реалізує interface, використовуємо:

```go
// ✅ Це перевірка на етапі компіляції
var _ StockOrderRepository = (*MockStockOrderRepository)(nil)

// Що це означає:
// 1. Створюємо nil pointer типу MockStockOrderRepository
// 2. Присвоюємо його змінній типу StockOrderRepository
// 3. Якщо Mock НЕ реалізує всі методи → compile error!
```

### Приклад помилки:

```go
type MockRepo struct {}

// Забули реалізувати Cancel()!
func (m *MockRepo) Create(order *Order) error { ... }
func (m *MockRepo) GetByID(id int64) (*Order, error) { ... }

// Compile-time check виявить помилку:
var _ StockOrderRepository = (*MockRepo)(nil)

// Error:
// cannot use (*MockRepo)(nil) (type *MockRepo) as type StockOrderRepository
// in assignment: *MockRepo does not implement StockOrderRepository
// (missing Cancel method)
```

---

## 📊 Де що визначається

```
internal/repository/
├── stock_order_repository.go
│   │
│   ├── type StockOrderRepository interface {   ← INTERFACE (один раз!)
│   │       Create()
│   │       GetByID()
│   │       Cancel()
│   │   }
│   │
│   ├── type stockOrderRepositoryImpl struct   ← PostgreSQL implementation
│   │
│   └── var _ StockOrderRepository = (*stockOrderRepositoryImpl)(nil)  ← Check
│
└── mock_stock_order_repository.go
    │
    ├── // НЕ дублюємо interface!
    │
    ├── type MockStockOrderRepository struct   ← Mock implementation
    │
    └── var _ StockOrderRepository = (*MockStockOrderRepository)(nil)  ← Check
```

---

## 🎯 Як це працює

### 1. Interface визначений в основному файлі

```go
// stock_order_repository.go

type StockOrderRepository interface {
    Create(order *model.StockOrder) error
    GetByID(id int64) (*model.StockOrder, error)
    Cancel(orderID, userID int64) error
}
```

### 2. PostgreSQL реалізація (implicit)

```go
// stock_order_repository.go

type stockOrderRepositoryImpl struct {
    db *sql.DB
}

// Реалізуємо всі 3 методи = автоматично реалізуємо interface!
func (r *stockOrderRepositoryImpl) Create(...) error { /* SQL */ }
func (r *stockOrderRepositoryImpl) GetByID(...) (*Order, error) { /* SQL */ }
func (r *stockOrderRepositoryImpl) Cancel(...) error { /* SQL */ }

// Перевірка
var _ StockOrderRepository = (*stockOrderRepositoryImpl)(nil) ✅
```

### 3. Mock реалізація (implicit)

```go
// mock_stock_order_repository.go

type MockStockOrderRepository struct {
    CreateFunc func(...) error
}

// Реалізуємо всі 3 методи = автоматично реалізуємо interface!
func (m *MockStockOrderRepository) Create(...) error { /* Mock */ }
func (m *MockStockOrderRepository) GetByID(...) (*Order, error) { /* Mock */ }
func (m *MockStockOrderRepository) Cancel(...) error { /* Mock */ }

// Перевірка
var _ StockOrderRepository = (*MockStockOrderRepository)(nil) ✅
```

---

## 🔄 Використання

### Production (PostgreSQL):

```go
// Повертає interface, а не конкретний тип
repo := repository.NewStockOrderRepository(db, "postgres")
//       ↓ returns StockOrderRepository interface
service := service.NewStockOrderServiceV2(repo)
```

### Testing (Mock):

```go
// Також повертає interface (той самий!)
mockRepo := repository.NewMockStockOrderRepository()
//           ↓ returns StockOrderRepository interface (same!)
service := service.NewStockOrderServiceV2(mockRepo)
```

### Service не знає різниці:

```go
type StockOrderServiceV2 struct {
    repo repository.StockOrderRepository  // ← Interface!
    // Service не знає чи це PostgreSQL, MongoDB, або Mock
    // Service працює тільки з методами interface
}

func (s *Service) CreateOrder(...) {
    // Викликає Create() з interface
    // Під капотом це може бути:
    //   - PostgreSQL.Create() (production)
    //   - Mock.Create() (testing)
    //   - MongoDB.Create() (якщо створимо)
    return s.repo.Create(order)
}
```

---

## 💡 Чому так зроблено в Go?

### Duck Typing

> "If it walks like a duck and quacks like a duck, it's a duck"

```go
// Interface в Go - це "duck typing"
type Duck interface {
    Quack() string
    Walk() string
}

// Якщо тип має методи Quack() та Walk() → він Duck!
type RealDuck struct {}
func (r *RealDuck) Quack() string { return "Quack!" }
func (r *RealDuck) Walk() string { return "Waddle!" }

type Robot struct {}
func (r *Robot) Quack() string { return "Beep!" }
func (r *Robot) Walk() string { return "Roll!" }

// Обидва реалізують Duck interface!
var _ Duck = (*RealDuck)(nil)  ✅
var _ Duck = (*Robot)(nil)     ✅
```

---

## 🎯 Best Practice

### ✅ Правильно:

```go
// Define interface once
// stock_order_repository.go
type Repository interface {
    Create(...) error
}

// Implementation 1
type PostgresRepo struct {}
func (r *PostgresRepo) Create(...) error { /* SQL */ }
var _ Repository = (*PostgresRepo)(nil)

// Implementation 2 (in separate file)
// mock_stock_order_repository.go
type MockRepo struct {}
func (m *MockRepo) Create(...) error { /* Mock */ }
var _ Repository = (*MockRepo)(nil)
```

### ❌ НЕправильно:

```go
// DON'T repeat interface in each file!

// stock_order_repository.go
type Repository interface { Create() }

// mock_stock_order_repository.go
type Repository interface { Create() }  // ❌ DUPLICATE!
```

---

## 📚 Структура (фінальна)

```
internal/repository/
│
├── stock_order_repository.go
│   ├── type StockOrderRepository interface { ... }  ← INTERFACE (ONE TIME)
│   ├── type stockOrderRepositoryImpl struct { ... } ← PostgreSQL impl
│   ├── var _ StockOrderRepository = (*stockOrderRepositoryImpl)(nil)
│   └── func (r *stockOrderRepositoryImpl) Create() { SQL }
│
└── mock_stock_order_repository.go
    ├── // Interface вже визначений в stock_order_repository.go
    ├── type MockStockOrderRepository struct { ... }  ← Mock impl
    ├── var _ StockOrderRepository = (*MockStockOrderRepository)(nil)
    └── func (m *MockStockOrderRepository) Create() { Mock }

Both implement the SAME interface!
```

---

## ✅ Тепер в обох файлах є перевірка:

```go
// stock_order_repository.go
var _ StockOrderRepository = (*stockOrderRepositoryImpl)(nil)  ✅

// mock_stock_order_repository.go
var _ StockOrderRepository = (*MockStockOrderRepository)(nil)   ✅
```

**Це гарантує що обидві реалізації дотримуються контракту interface!**

---

## 🎓 Висновок

### Чому interface в одному файлі?

1. ✅ **DRY** - не дублюємо визначення
2. ✅ **Single Source of Truth** - одне місце для змін
3. ✅ **Go convention** - implicit implementation
4. ✅ **Compile-time safety** - перевірка з `var _`

### Як переконатися що все працює?

```bash
go build ./internal/repository

# Якщо Mock НЕ реалізує всі методи → compile error!
# Якщо компілюється → Mock реалізує interface ✅
```

**Це стандартна практика в Go!** 🚀
