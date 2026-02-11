# Go Embedding: Прибрати Behavior/Calls з викликів 🎯

Два підходи до структури Mock.

---

## Варіант 1: Nested (поточний) ✅

```go
type MockBehavior struct {
    CreateFunc   func(order *model.StockOrder) error
    GetByIDFunc  func(id int64) (*model.StockOrder, error)
}

type MockCalls struct {
    CreateCalled  bool
    GetByIDCalled bool
}

type MockStockOrderRepository struct {
    Behavior MockBehavior  // ← Явне поле
    Calls    MockCalls     // ← Явне поле
}

// Використання:
mockRepo.Behavior.CreateFunc = func(...) error { return nil }
if mockRepo.Calls.CreateCalled { ... }
```

**Переваги:**
- ✅ Явна групування (видно що `Behavior` vs `Calls`)
- ✅ Зрозуміло призначення кожного поля
- ✅ Легко читати код

**Недоліки:**
- ❌ Довші виклики (`mockRepo.Behavior.CreateFunc`)
- ❌ Додатковий рівень вкладеності

---

## Варіант 2: Embedded (альтернатива) 🚀

```go
type MockBehavior struct {
    CreateFunc   func(order *model.StockOrder) error
    GetByIDFunc  func(id int64) (*model.StockOrder, error)
}

type MockCalls struct {
    CreateCalled  bool
    GetByIDCalled bool
}

type MockStockOrderRepository struct {
    MockBehavior  // ← Embedded (анонімне поле)
    MockCalls     // ← Embedded (анонімне поле)
}

// Використання:
mockRepo.CreateFunc = func(...) error { return nil }  // ← Без .Behavior!
if mockRepo.CreateCalled { ... }                      // ← Без .Calls!
```

**Переваги:**
- ✅ Коротші виклики
- ✅ Менше коду
- ✅ Прямий доступ до полів

**Недоліки:**
- ❌ Втрата явної групування
- ❌ Менш зрозуміло що `CreateFunc` vs `CreateCalled`
- ❌ Можливі конфлікти імен (якщо в `MockBehavior` і `MockCalls` однакові імена)

---

## 🔍 Як працює Embedding?

### Embedding в Go:

```go
type A struct {
    Field1 string
}

type B struct {
    Field2 int
}

type C struct {
    A  // ← Embedded (promoted fields)
    B  // ← Embedded (promoted fields)
}

c := C{}
c.Field1 = "hello"  // ← Можна без c.A.Field1
c.Field2 = 42       // ← Можна без c.B.Field2

// Але також можна:
c.A.Field1 = "hello"  // ← Явний доступ через A
c.B.Field2 = 42       // ← Явний доступ через B
```

**Field promotion**: поля вбудованих типів "піднімаються" на рівень вище.

---

## 📊 Порівняння підходів

### Nested (Variant 1):

```go
type Mock struct {
    Behavior MockBehavior  // Explicit field
    Calls    MockCalls     // Explicit field
}

// Usage:
mockRepo.Behavior.CreateFunc = func() { ... }
mockRepo.Calls.CreateCalled = true

Pros:
  ✅ Clear intent (Behavior vs Calls)
  ✅ Better documentation
  ✅ No name conflicts

Cons:
  ❌ Longer calls
  ❌ More typing
```

---

### Embedded (Variant 2):

```go
type Mock struct {
    MockBehavior  // Embedded
    MockCalls     // Embedded
}

// Usage:
mockRepo.CreateFunc = func() { ... }
mockRepo.CreateCalled = true

Pros:
  ✅ Shorter calls
  ✅ Less typing
  ✅ Direct access

Cons:
  ❌ Less clear intent
  ❌ Harder to distinguish Behavior vs Calls
  ❌ Potential name conflicts
```

---

## 🎯 Приклад реалізації (Embedded)

```go
// internal/repository/mock_stock_order_repository.go

package repository

import (
	"fmt"
	"stock_hub/pkg/model"
)

// MockBehavior defines custom behavior functions
type MockBehavior struct {
	CreateFunc             func(order *model.StockOrder) error
	GetByIDFunc            func(id int64) (*model.StockOrder, error)
	GetOpenOrdersFunc      func(symbol string) ([]*model.StockOrder, error)
	GetAllOpenOrdersFunc   func() ([]*model.StockOrder, error)
	GetByUserIDFunc        func(userID int64, limit int) ([]*model.StockOrder, error)
	CancelFunc             func(orderID, userID int64) error
}

// MockCalls tracks method calls
type MockCalls struct {
	CreateCalled           bool
	GetByIDCalled          bool
	GetOpenOrdersCalled    bool
	GetAllOpenOrdersCalled bool
	GetByUserIDCalled      bool
	CancelCalled           bool
}

// MockStockOrderRepository with embedded types
type MockStockOrderRepository struct {
	MockBehavior  // ← Embedded (no field name)
	MockCalls     // ← Embedded (no field name)
}

// Compile-time check
var _ StockOrderRepository = (*MockStockOrderRepository)(nil)

// NewMockStockOrderRepository creates new mock
func NewMockStockOrderRepository() *MockStockOrderRepository {
	return &MockStockOrderRepository{
		MockBehavior: MockBehavior{},
		MockCalls:    MockCalls{},
	}
}

// Create implements interface
func (m *MockStockOrderRepository) Create(order *model.StockOrder) error {
	m.CreateCalled = true  // ← Direct access!
	if m.CreateFunc != nil {  // ← Direct access!
		return m.CreateFunc(order)
	}
	order.ID = 1
	return nil
}

// GetByID implements interface
func (m *MockStockOrderRepository) GetByID(id int64) (*model.StockOrder, error) {
	m.GetByIDCalled = true  // ← Direct access!
	if m.GetByIDFunc != nil {  // ← Direct access!
		return m.GetByIDFunc(id)
	}
	return nil, fmt.Errorf("order not found")
}

// ... інші методи
```

---

## 📝 Використання в тестах

### Nested (поточний):

```go
func TestCreateOrder(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    
    // Setup behavior
    mockRepo.Behavior.CreateFunc = func(order *model.StockOrder) error {
        order.ID = 123
        return nil
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    order, err := service.CreateOrder(...)
    
    // Verify
    if !mockRepo.Calls.CreateCalled {
        t.Error("Expected Create to be called")
    }
}
```

---

### Embedded (альтернатива):

```go
func TestCreateOrder(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    
    // Setup behavior (shorter!)
    mockRepo.CreateFunc = func(order *model.StockOrder) error {
        order.ID = 123
        return nil
    }
    
    service := NewStockOrderServiceV2(mockRepo)
    order, err := service.CreateOrder(...)
    
    // Verify (shorter!)
    if !mockRepo.CreateCalled {
        t.Error("Expected Create to be called")
    }
}
```

**Різниця: без `.Behavior` і `.Calls`!**

---

## 🤔 Який підхід обрати?

### Обирайте **Nested** (поточний) якщо:

- ✅ Команда велика (багато розробників)
- ✅ Важлива читабельність та явність
- ✅ Хочете чітке розділення Behavior vs Calls
- ✅ Проект довгостроковий (maintainability)

**Рекомендація: для production коду**

---

### Обирайте **Embedded** якщо:

- ✅ Команда мала (solo або 2-3 розробники)
- ✅ Хочете менше коду
- ✅ Пріоритет - швидкість написання тестів
- ✅ Всі розуміють структуру mock

**Рекомендація: для особистих проектів або прототипів**

---

## 🎭 Гібридний підхід

Можна зробити обидва способи доступу:

```go
type MockStockOrderRepository struct {
    MockBehavior  // Embedded - прямий доступ
    MockCalls     // Embedded - прямий доступ
}

// Прямий доступ (через embedding):
mockRepo.CreateFunc = func() { ... }
mockRepo.CreateCalled = true

// Явний доступ (через тип):
mockRepo.MockBehavior.CreateFunc = func() { ... }
mockRepo.MockCalls.CreateCalled = true

// ОБА працюють одночасно!
```

---

## 📚 Реальні приклади

### testify/mock (популярна бібліотека):

```go
type Mock struct {
    ExpectedCalls []*Call  // ← Явне поле
    Calls         []Call   // ← Явне поле
    mutex         sync.Mutex
}

// Використання:
mock.ExpectedCalls = append(...)
mock.Calls = []Call{}

// Явний підхід (як у нас зараз)
```

---

### gomock (офіційний mock від Go team):

```go
type MockController struct {
    mu           sync.Mutex
    expectedCalls *callSet
    calls        []*Call
}

// Теж явний підхід
```

---

### sqlmock (для SQL):

```go
type sqlmock struct {
    expectations []expectation  // ← Явне поле
    ordered      bool
}

// Явний підхід
```

**Висновок: більшість Go mock бібліотек використовують явні поля (Nested)!**

---

## ⚖️ Підсумок

| Критерій | Nested (поточний) | Embedded (альтернатива) |
|----------|------------------|------------------------|
| Читабельність | ✅ Відмінна | ⚠️ Добра |
| Довжина коду | ⚠️ Довше | ✅ Коротше |
| Явність | ✅ Дуже явно | ❌ Менш явно |
| Maintainability | ✅ Легко | ⚠️ Середньо |
| Популярність | ✅ Стандарт | ❌ Рідко |
| Для великих команд | ✅ Так | ❌ Ні |
| Для solo проектів | ✅ Так | ✅ Так |

---

## 💡 Моя рекомендація

### Залишити Nested (поточний підхід) ✅

**Причини:**

1. ✅ **Явність** - зрозуміло що `Behavior` (налаштування) vs `Calls` (верифікація)
2. ✅ **Best Practice** - так роблять testify, gomock, sqlmock
3. ✅ **Maintainability** - легко підтримувати в майбутньому
4. ✅ **Team-friendly** - зрозуміло для нових розробників
5. ✅ **Documentation** - self-documenting code

**Незначний недолік:**
- ❌ Довші виклики (`mockRepo.Behavior.CreateFunc`)

**Але це не проблема:**
- IDE autocomplete допомагає
- Явність важливіша за 10 символів коду
- Production код читають частіше ніж пишуть

---

## 🔄 Якщо все ж таки хочете Embedded

Переваги:
- Коротший код
- Менше typing

Недоліки:
- Менш явно
- Важче розуміти новим розробникам

**Вибір за вами!** Обидва підходи валідні в Go.

---

## 📖 Висновок

### Поточний (Nested):
```go
type Mock struct {
    Behavior MockBehavior  // Explicit
    Calls    MockCalls     // Explicit
}

mockRepo.Behavior.CreateFunc = ...
mockRepo.Calls.CreateCalled
```

**Рекомендація: ✅ Залишити для production**

---

### Альтернатива (Embedded):
```go
type Mock struct {
    MockBehavior  // Embedded
    MockCalls     // Embedded
}

mockRepo.CreateFunc = ...
mockRepo.CreateCalled
```

**Використання: ⚠️ Для особистих проектів**

---

**Мій вибір: залишити Nested підхід** для кращої читабельності та maintainability! 🚀
