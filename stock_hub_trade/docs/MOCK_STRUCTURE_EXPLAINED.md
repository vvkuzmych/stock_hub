# Mock Structure: Чому розділено на Behavior та Calls 🎯

Пояснення структури `MockStockOrderRepository`.

---

## 🔑 Питання

> "Чому ти об'єднав все в один struct, а не розділив для інтерфейсу?"

**Відповідь: Тепер розділено на 2 логічні групи!**

---

## 📊 До vs Після

### ❌ Було (все в купі):

```go
type MockStockOrderRepository struct {
    // Mock functions (для налаштування поведінки)
    CreateFunc             func(order *model.StockOrder) error
    GetByIDFunc            func(id int64) (*model.StockOrder, error)
    CancelFunc             func(orderID, userID int64) error
    
    // Tracking flags (для перевірки викликів)
    CreateCalled           bool
    GetByIDCalled          bool
    CancelCalled           bool
}

Problems:
  🔴 12 полів в одному struct - важко читати
  🔴 Незрозуміло що для чого
  🔴 Немає логічної групування
```

---

### ✅ Стало (розділено на групи):

```go
// Group 1: Behavior (налаштування поведінки)
type MockBehavior struct {
    CreateFunc             func(order *model.StockOrder) error
    GetByIDFunc            func(id int64) (*model.StockOrder, error)
    GetOpenOrdersFunc      func(symbol string) ([]*model.StockOrder, error)
    GetAllOpenOrdersFunc   func() ([]*model.StockOrder, error)
    GetByUserIDFunc        func(userID int64, limit int) ([]*model.StockOrder, error)
    CancelFunc             func(orderID, userID int64) error
}

// Group 2: Calls (відстеження викликів)
type MockCalls struct {
    CreateCalled           bool
    GetByIDCalled          bool
    GetOpenOrdersCalled    bool
    GetAllOpenOrdersCalled bool
    GetByUserIDCalled      bool
    CancelCalled           bool
}

// Main Mock (композиція)
type MockStockOrderRepository struct {
    Behavior MockBehavior  // ← Що робити
    Calls    MockCalls     // ← Що було викликано
}

Benefits:
  ✅ Чітке розділення відповідальностей
  ✅ Легко читати (2 групи замість 12 полів)
  ✅ Зрозуміла структура
  ✅ Легко розширювати
```

---

## 🎯 Використання

### 1. Налаштування поведінки (Behavior)

```go
mockRepo := repository.NewMockStockOrderRepository()

// Симулюємо помилку БД
mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
    return fmt.Errorf("database connection failed")
}

// Симулюємо успішне створення
mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
    order.ID = 123  // Assign ID
    return nil
}

// Симулюємо "order not found"
mockRepo.Behavior.GetByIDFunc = func(id int64) (*model.StockOrder, error) {
    return nil, fmt.Errorf("order not found")
}
```

**Призначення `Behavior`:**
- ✅ Налаштувати **що повертає** mock
- ✅ Симулювати різні сценарії (success, error, not found)
- ✅ Контролювати поведінку в тесті

---

### 2. Перевірка викликів (Calls)

```go
service := NewStockOrderServiceV2(mockRepo)
order, err := service.CreateOrder(...)

// Перевіряємо чи був виклик
if !mockRepo.Calls.CreateCalled {
    t.Error("Expected Create to be called")
}

// Перевіряємо що НЕ був виклик
if mockRepo.Calls.GetByIDCalled {
    t.Error("GetByID should not be called")
}
```

**Призначення `Calls`:**
- ✅ Перевірити **чи був виклик** методу
- ✅ Verify взаємодію між Service і Repository
- ✅ Ensure правильний flow

---

## 📋 Повний приклад тесту

```go
func TestCreateOrderV2_RepositoryError(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    
    // 1. Налаштувати поведінку (Behavior)
    mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
        return fmt.Errorf("database connection failed")
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    
    // 2. Виконати тест
    _, err := service.CreateOrder(1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    
    // 3. Перевірити результат
    if err == nil {
        t.Error("Expected error when repository fails")
    }
    
    // 4. Перевірити виклики (Calls)
    if !mockRepo.Calls.CreateCalled {
        t.Error("Expected repository Create to be called")
    }
}
```

---

## 🏗️ Структура тепер

```
MockStockOrderRepository
├── Behavior (MockBehavior)
│   ├── CreateFunc              ← Що повертати при Create()
│   ├── GetByIDFunc             ← Що повертати при GetByID()
│   ├── GetOpenOrdersFunc       ← Що повертати при GetOpenOrders()
│   ├── GetAllOpenOrdersFunc    ← Що повертати при GetAllOpenOrders()
│   ├── GetByUserIDFunc         ← Що повертати при GetByUserID()
│   └── CancelFunc              ← Що повертати при Cancel()
│
└── Calls (MockCalls)
    ├── CreateCalled            ← Чи був виклик Create()?
    ├── GetByIDCalled           ← Чи був виклик GetByID()?
    ├── GetOpenOrdersCalled     ← Чи був виклик GetOpenOrders()?
    ├── GetAllOpenOrdersCalled  ← Чи був виклик GetAllOpenOrders()?
    ├── GetByUserIDCalled       ← Чи був виклик GetByUserID()?
    └── CancelCalled            ← Чи був виклик Cancel()?
```

---

## 💡 Переваги розділення

### 1. **Separation of Concerns** ✅

```
Behavior: WHAT to return (налаштування)
Calls:    WAS it called (верифікація)
```

### 2. **Читабельність** ✅

```go
// ❌ Було
mockRepo.CreateFunc = ...
mockRepo.CreateCalled

// ✅ Стало
mockRepo.Behavior.CreateFunc = ...  // Налаштування
mockRepo.Calls.CreateCalled         // Перевірка

Зрозуміліше що для чого!
```

### 3. **Розширюваність** ✅

Легко додати нові групи:

```go
type MockStockOrderRepository struct {
    Behavior MockBehavior
    Calls    MockCalls
    Stats    MockStats    // ← Нова група (лічильники)
    Errors   MockErrors   // ← Нова група (помилки)
}

type MockStats struct {
    CreateCount   int
    GetByIDCount  int
}
```

### 4. **Ізоляція** ✅

Кожна група має свою відповідальність:

```go
// Налаштувати поведінку
mockRepo.Behavior = MockBehavior{
    CreateFunc: func(...) error { return nil },
}

// Перевірити виклики
if mockRepo.Calls.CreateCalled {
    // assertions
}

// Очистити виклики для нового тесту
mockRepo.Calls = MockCalls{}
```

---

## 📈 Порівняння підходів

### Варіант 1: Flat Structure (❌ Було)

```go
type Mock struct {
    CreateFunc   func() error
    CreateCalled bool
    GetByIDFunc  func() error
    GetByIDCalled bool
    // ... 12 полів в одному рівні
}

mockRepo.CreateFunc = ...
mockRepo.CreateCalled

Problems:
  🔴 Важко читати
  🔴 Немає групування
  🔴 Складно розширювати
```

---

### Варіант 2: Nested Structure (✅ Стало)

```go
type MockBehavior struct {
    CreateFunc  func() error
    GetByIDFunc func() error
}

type MockCalls struct {
    CreateCalled  bool
    GetByIDCalled bool
}

type Mock struct {
    Behavior MockBehavior
    Calls    MockCalls
}

mockRepo.Behavior.CreateFunc = ...
mockRepo.Calls.CreateCalled

Benefits:
  ✅ Чітка структура
  ✅ Логічне групування
  ✅ Легко розширювати
```

---

### Варіант 3: Separate Structs (альтернатива)

```go
type MockBehavior struct {
    Create  func() error
    GetByID func() error
}

type MockTracker struct {
    CreateCalls  int
    GetByIDCalls int
}

type Mock struct {
    behavior MockBehavior
    tracker  MockTracker
}

// Private fields + getter methods
func (m *Mock) SetBehavior(b MockBehavior) { m.behavior = b }
func (m *Mock) GetCalls() MockTracker { return m.tracker }
```

**Наш підхід (Варіант 2) - найпростіший і найпоширеніший!**

---

## 🧪 Приклади використання

### Тест 1: Симуляція помилки БД

```go
func TestRepositoryError(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    
    // Налаштування поведінки
    mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
        return fmt.Errorf("connection timeout")
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    _, err := service.CreateOrder(...)
    
    // Перевірка результату
    if err == nil {
        t.Error("Expected error")
    }
    
    // Перевірка викликів
    if !mockRepo.Calls.CreateCalled {
        t.Error("Expected Create to be called")
    }
}
```

---

### Тест 2: Симуляція успіху

```go
func TestCreateSuccess(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    
    // Налаштування: повертаємо ID
    mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
        order.ID = 999
        return nil
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    order, err := service.CreateOrder(...)
    
    if err != nil {
        t.Error("Expected no error")
    }
    
    if order.ID != 999 {
        t.Errorf("Expected ID 999, got %d", order.ID)
    }
    
    if !mockRepo.Calls.CreateCalled {
        t.Error("Expected Create to be called")
    }
}
```

---

### Тест 3: Перевірка що метод НЕ викликався

```go
func TestValidationFailsBeforeRepository(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    
    // Невалідний price (service має відхилити до виклику repository)
    _, err := service.CreateOrder(1, "john", "AAPL", model.OrderTypeBid, -100, 10)
    
    if err == nil {
        t.Error("Expected validation error")
    }
    
    // Перевіряємо що repository НЕ викликався
    if mockRepo.Calls.CreateCalled {
        t.Error("Repository should NOT be called for invalid input")
    }
}
```

---

## 🎯 Чому саме така структура?

### 1. **Clear Intent** (зрозумілий намір)

```go
// ✅ Зрозуміло - це налаштування
mockRepo.Behavior.CreateFunc = func() { ... }

// ✅ Зрозуміло - це перевірка
if mockRepo.Calls.CreateCalled { ... }

// ❌ Було незрозуміло
mockRepo.CreateFunc = ...
mockRepo.CreateCalled
```

### 2. **Maintainability** (легка підтримка)

```go
// Легко додати нові групи
type Mock struct {
    Behavior MockBehavior
    Calls    MockCalls
    Stats    MockStats    // ← Нова група
    Errors   []error      // ← Нова група
}
```

### 3. **Testability** (зручність тестування)

```go
// Очистити tracking між sub-tests
t.Run("test1", func(t *testing.T) {
    mockRepo.Calls = MockCalls{}  // Reset
    // test...
})

t.Run("test2", func(t *testing.T) {
    mockRepo.Calls = MockCalls{}  // Reset
    // test...
})
```

---

## 📚 Аналоги в інших мовах

### Python (unittest.mock)

```python
from unittest.mock import Mock

mock_repo = Mock()

# Behavior
mock_repo.create.return_value = Order(id=123)
mock_repo.create.side_effect = Exception("DB error")

# Calls
mock_repo.create.assert_called_once()
mock_repo.create.assert_not_called()
```

### Java (Mockito)

```java
Repository mockRepo = mock(Repository.class);

// Behavior
when(mockRepo.create(any())).thenReturn(order);
when(mockRepo.create(any())).thenThrow(new SQLException());

// Calls
verify(mockRepo).create(any());
verify(mockRepo, never()).getById(anyLong());
```

### Go (наш підхід)

```go
mockRepo := repository.NewMockStockOrderRepository()

// Behavior
mockRepo.Behavior.CreateFunc = func(order *Order) error {
    order.ID = 123
    return nil
}

// Calls
if !mockRepo.Calls.CreateCalled {
    t.Error("Expected call")
}
```

**Наш підхід схожий на інші мови, але Go-idiomatic!**

---

## 🔄 Фінальна структура

```go
// internal/repository/mock_stock_order_repository.go

type MockBehavior struct {
    CreateFunc   func(order *model.StockOrder) error
    GetByIDFunc  func(id int64) (*model.StockOrder, error)
    CancelFunc   func(orderID, userID int64) error
    // ... інші методи
}

type MockCalls struct {
    CreateCalled  bool
    GetByIDCalled bool
    CancelCalled  bool
    // ... інші флаги
}

type MockStockOrderRepository struct {
    Behavior MockBehavior
    Calls    MockCalls
}

// Methods
func (m *MockStockOrderRepository) Create(order *model.StockOrder) error {
    m.Calls.CreateCalled = true         // ← Track call
    if m.Behavior.CreateFunc != nil {   // ← Use custom behavior
        return m.Behavior.CreateFunc(order)
    }
    order.ID = 1  // Default behavior
    return nil
}
```

---

## ✅ Результати

```bash
go test ./internal/service -v -run V2

=== RUN   TestCreateOrderV2_ValidatesPrice
--- PASS: TestCreateOrderV2_ValidatesPrice (0.00s)
=== RUN   TestCreateOrderV2_ValidatesQuantity
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
ok  	stock_hub/internal/service	0.753s
```

**Всі тести пройшли з новою структурою!** ✅

---

## 🎓 Висновок

### Чому розділили на `Behavior` і `Calls`?

1. ✅ **Clear Separation** - налаштування vs перевірка
2. ✅ **Readable** - зрозуміло що для чого
3. ✅ **Maintainable** - легко додавати нові групи
4. ✅ **Testable** - легко reset між тестами
5. ✅ **Best Practice** - схоже на Mockito/unittest.mock

### Типи структур:

```
MockBehavior: ЯК поводиться mock (налаштування)
MockCalls:    ЩО було викликано (верифікація)
Mock:         Композиція обох (головний тип)
```

**Це стандартний підхід для Go mocks!** 🚀

---

## 📚 Додаткові ресурси

- [Testing with Mocks in Go](https://go.dev/blog/test-fixtures)
- [Testify Mock Package](https://github.com/stretchr/testify#mock-package)
- [Go Mock Best Practices](https://github.com/golang/mock)
