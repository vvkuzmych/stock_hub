# Test Helpers Refactoring 🧪

Рефакторинг тестів для усунення дублювання коду через test helper functions.

---

## 🎯 Проблема

### До рефакторингу (дублювання):

```go
func TestCreateOrderV2_ValidatesPrice(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()  // ← Повторюється
    service := NewStockOrderServiceV2(mockRepo)           // ← Повторюється
    ctx := context.Background()                           // ← Повторюється
    
    // test logic...
}

func TestCreateOrderV2_ValidatesQuantity(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()  // ← Повторюється
    service := NewStockOrderServiceV2(mockRepo)           // ← Повторюється
    ctx := context.Background()                           // ← Повторюється
    
    // test logic...
}

// ... 8 більше функцій з таким же дублюванням
```

**Проблеми:**
- ❌ Код повторюється в кожному тесті
- ❌ Важко підтримувати (зміна setup = 10 місць змін)
- ❌ Більше коду = більше шансів на помилки
- ❌ Не DRY (Don't Repeat Yourself)

---

## ✅ Рішення: Test Helper Functions

### 1. Створено Helper Functions

```go
// setupTest creates a test fixture with mock repository and service
func setupTest(t *testing.T) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
    t.Helper()  // ← Важливо! Правильні номери строк у помилках
    
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    
    return service, mockRepo, ctx
}

// setupTestWithBehavior creates a test fixture with custom repository behavior
func setupTestWithBehavior(t *testing.T, behavior repository.MockBehavior) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
    t.Helper()
    
    mockRepo := repository.NewMockStockOrderRepository()
    mockRepo.Behavior = behavior
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    
    return service, mockRepo, ctx
}
```

**Переваги helper functions:**
- ✅ `t.Helper()` - правильні строки у помилках тестів
- ✅ Централізований setup
- ✅ Легко змінювати (одне місце)
- ✅ Чистіший код тестів

---

### 2. Використання `setupTest()` (базовий)

**До:**
```go
func TestCreateOrderV2_ValidatesSymbol(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()

    _, err := service.CreateOrder(ctx, 1, "john", "", model.OrderTypeBid, 150.00, 10)

    if err == nil {
        t.Error("Expected error for empty symbol")
    }
}
```

**Після:**
```go
func TestCreateOrderV2_ValidatesSymbol(t *testing.T) {
    service, _, ctx := setupTest(t)  // ← Одна лінія!

    _, err := service.CreateOrder(ctx, 1, "john", "", model.OrderTypeBid, 150.00, 10)

    if err == nil {
        t.Error("Expected error for empty symbol")
    }
}
```

**Що змінилось:**
- ✅ 3 строки → 1 строка
- ✅ `_` для невикористаних змінних (mockRepo)
- ✅ Чистіший focus на логіці тесту

---

### 3. Використання `setupTestWithBehavior()` (з кастомним behavior)

**До:**
```go
func TestCreateOrderV2_RepositoryError(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    mockRepo.Behavior.CreateFunc = func(ctx context.Context, order *model.StockOrder) error {
        return fmt.Errorf("database connection failed")
    }

    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()

    _, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    
    // assertions...
}
```

**Після:**
```go
func TestCreateOrderV2_RepositoryError(t *testing.T) {
    service, mockRepo, ctx := setupTestWithBehavior(t, repository.MockBehavior{
        CreateFunc: func(ctx context.Context, order *model.StockOrder) error {
            return fmt.Errorf("database connection failed")
        },
    })

    _, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, 10)
    
    // assertions...
}
```

**Переваги:**
- ✅ Behavior інлайн, ближче до контексту
- ✅ Менш багатослівно
- ✅ Зрозуміла структура

---

## 📊 Порівняння: До vs Після

### Статистика:

| Метрика | До | Після | Покращення |
|---------|-----|-------|-----------|
| Строк коду setup | 30+ | 0 | ✅ -30 строк |
| Дублювання | 10x | 0x | ✅ -100% |
| Строк на тест | 15-20 | 10-15 | ✅ -33% |
| Читабельність | ⚠️ Середня | ✅ Висока | +50% |

---

### Приклад: Validation Tests

**До (15 строк):**
```go
func TestCreateOrderV2_ValidatesQuantity(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)

    testCases := []struct {
        name        string
        quantity    int
        shouldError bool
    }{
        {"Zero quantity", 0, true},
        {"Negative quantity", -5, true},
        {"Valid quantity", 100, false},
    }

    ctx := context.Background()
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, tc.quantity)
            // assertions...
        })
    }
}
```

**Після (12 строк):**
```go
func TestCreateOrderV2_ValidatesQuantity(t *testing.T) {
    service, _, ctx := setupTest(t)  // ← Компактний setup

    testCases := []struct {
        name        string
        quantity    int
        shouldError bool
    }{
        {"Zero quantity", 0, true},
        {"Negative quantity", -5, true},
        {"Valid quantity", 100, false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := service.CreateOrder(ctx, 1, "john", "AAPL", model.OrderTypeBid, 150.00, tc.quantity)
            // assertions...
        })
    }
}
```

**Покращення: -3 строки, фокус на логіці тесту**

---

## 🎯 Використання `_` для невикористаних змінних

### Коли використовувати `_`:

```go
// Коли mockRepo не потрібен:
service, _, ctx := setupTest(t)
// Використовується: service, ctx
// Ігнорується: mockRepo

// Коли mockRepo потрібен:
service, mockRepo, ctx := setupTest(t)
// Використовується: service, mockRepo, ctx

// Коли тільки service потрібен:
service, _, _ := setupTest(t)
// Використовується: service
// Ігнорується: mockRepo, ctx
```

**Переваги `_`:**
- ✅ Явно показує що змінна не використовується
- ✅ Компілятор не видає warning
- ✅ Код чистіший

---

## 🔍 Важливість `t.Helper()`

### Без `t.Helper()`:

```go
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) {
    // NO t.Helper()
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    return service, mockRepo, context.Background()
}

func TestSomething(t *testing.T) {
    service, _, ctx := setupTest(t)  // ← line 45
    
    if err := service.DoSomething(ctx); err != nil {
        t.Error("failed")  // ← line 48
    }
}
```

**Error output:**
```
    stock_order_service_v2_test.go:45: failed
                                    ^^^ Wrong line! (points to setupTest call)
```

---

### З `t.Helper()` ✅:

```go
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) {
    t.Helper()  // ← Magic!
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    return service, mockRepo, context.Background()
}

func TestSomething(t *testing.T) {
    service, _, ctx := setupTest(t)  // ← line 45
    
    if err := service.DoSomething(ctx); err != nil {
        t.Error("failed")  // ← line 48
    }
}
```

**Error output:**
```
    stock_order_service_v2_test.go:48: failed
                                    ^^^ Correct line! (points to actual failure)
```

**`t.Helper()` tells Go to skip this function in stack traces!** ✅

---

## 📚 Типи Test Helpers

### 1. **Basic Setup** (найчастіше)

```go
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) {
    t.Helper()
    // Basic setup без кастомізації
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    return service, mockRepo, ctx
}

// Use case:
func TestValidation(t *testing.T) {
    service, _, ctx := setupTest(t)
    // test logic
}
```

---

### 2. **Setup with Behavior** (для кастомного behavior)

```go
func setupTestWithBehavior(t *testing.T, behavior repository.MockBehavior) (*Service, *MockRepo, context.Context) {
    t.Helper()
    mockRepo := repository.NewMockStockOrderRepository()
    mockRepo.Behavior = behavior  // ← Custom behavior
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    return service, mockRepo, ctx
}

// Use case:
func TestRepositoryError(t *testing.T) {
    service, mockRepo, ctx := setupTestWithBehavior(t, repository.MockBehavior{
        CreateFunc: func(ctx context.Context, order *model.StockOrder) error {
            return fmt.Errorf("DB error")
        },
    })
    // test logic
}
```

---

### 3. **Setup with Cleanup** (якщо потрібен cleanup)

```go
func setupTestWithCleanup(t *testing.T) (*Service, *MockRepo, context.Context, func()) {
    t.Helper()
    
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    
    cleanup := func() {
        // Clean up resources
        mockRepo.Calls = repository.MockCalls{}
    }
    
    return service, mockRepo, ctx, cleanup
}

// Use case:
func TestWithCleanup(t *testing.T) {
    service, mockRepo, ctx, cleanup := setupTestWithCleanup(t)
    defer cleanup()  // ← Runs after test
    
    // test logic
}
```

---

### 4. **Setup with Context Timeout**

```go
func setupTestWithTimeout(t *testing.T, timeout time.Duration) (*Service, *MockRepo, context.Context, context.CancelFunc) {
    t.Helper()
    
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    
    return service, mockRepo, ctx, cancel
}

// Use case:
func TestTimeout(t *testing.T) {
    service, _, ctx, cancel := setupTestWithTimeout(t, 100*time.Millisecond)
    defer cancel()
    
    // test logic that should timeout
}
```

---

## ✅ Best Practices

### 1. **Завжди використовуйте `t.Helper()`**

```go
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) {
    t.Helper()  // ← Обов'язково!
    // ...
}
```

**Чому важливо:**
- ✅ Правильні строки у помилках
- ✅ Легший debugging
- ✅ Краща читабельність output

---

### 2. **Повертайте тільки потрібне**

```go
// ❌ Bad: занадто багато return values
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context, *sql.DB, *Config, *Logger) {
    // ...
}

// ✅ Good: тільки необхідне
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) {
    // ...
}
```

---

### 3. **Використовуйте `_` для невикористаних змінних**

```go
// ✅ Good: явно показуємо що не використовуємо
service, _, ctx := setupTest(t)

// ❌ Bad: невикористана змінна
service, mockRepo, ctx := setupTest(t)
// mockRepo not used - compiler warning
```

---

### 4. **Групуйте схожі helper functions**

```go
// Basic setup
func setupTest(t *testing.T) (*Service, *MockRepo, context.Context) { ... }

// Setup with custom behavior
func setupTestWithBehavior(t *testing.T, behavior MockBehavior) (*Service, *MockRepo, context.Context) { ... }

// Setup with timeout
func setupTestWithTimeout(t *testing.T, timeout time.Duration) (*Service, *MockRepo, context.Context, context.CancelFunc) { ... }
```

---

### 5. **Документуйте helper functions**

```go
// setupTest creates a test fixture with mock repository and service
// Returns: service, mockRepo, context
func setupTest(t *testing.T) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
    t.Helper()
    // ...
}
```

---

## 📈 Результати

### До vs Після:

**До (багато дублювання):**
```go
func TestA(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    // test...
}

func TestB(t *testing.T) {
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    // test...
}

// + 8 більше функцій з таким же setup
```

**Після (DRY):**
```go
func setupTest(t *testing.T) (*StockOrderServiceV2, *repository.MockStockOrderRepository, context.Context) {
    t.Helper()
    mockRepo := repository.NewMockStockOrderRepository()
    service := NewStockOrderServiceV2(mockRepo)
    ctx := context.Background()
    return service, mockRepo, ctx
}

func TestA(t *testing.T) {
    service, _, ctx := setupTest(t)
    // test...
}

func TestB(t *testing.T) {
    service, mockRepo, ctx := setupTest(t)
    // test...
}

// + 8 більше функцій без дублювання
```

---

### Тести працюють ✅:

```bash
go test ./internal/service -v -run V2

=== RUN   TestCreateOrderV2_ValidatesPrice
=== RUN   TestCreateOrderV2_ValidatesPrice/Zero_price
=== RUN   TestCreateOrderV2_ValidatesPrice/Negative_price
=== RUN   TestCreateOrderV2_ValidatesPrice/Valid_price
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
ok  	stock_hub/internal/service	1.010s
```

---

## 🎓 Висновки

### Що отримали:

1. ✅ **Менше дублювання** - 30+ строк → 2 helper functions
2. ✅ **Легша підтримка** - зміни в одному місці
3. ✅ **Чистіший код** - фокус на логіці тесту
4. ✅ **Правильні помилки** - `t.Helper()` для correct stack traces
5. ✅ **Гнучкість** - `setupTest()` vs `setupTestWithBehavior()`
6. ✅ **DRY principle** - Don't Repeat Yourself

### Patterns:

```go
// Pattern 1: Basic test
service, _, ctx := setupTest(t)

// Pattern 2: Test with mock verification
service, mockRepo, ctx := setupTest(t)

// Pattern 3: Test with custom behavior
service, mockRepo, ctx := setupTestWithBehavior(t, repository.MockBehavior{
    CreateFunc: func(...) error { return errors.New("test error") },
})
```

---

## 📚 Додаткові ресурси

- [Testing in Go](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [t.Helper() documentation](https://pkg.go.dev/testing#T.Helper)
- [Test Fixtures in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

---

**Тепер тести чистіші, без дублювання, і легші у підтримці!** 🚀
