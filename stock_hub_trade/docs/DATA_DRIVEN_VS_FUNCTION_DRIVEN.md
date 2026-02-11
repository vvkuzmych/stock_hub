# Data-Driven vs Function-Driven Mock Setup 🎯

Порівняння двох підходів до організації mock setup в тестах.

---

## ❌ Неможливо: Блок коду як поле

### Те що ви хочете (НЕ працює):

```go
testCases := []struct {
    name      string
    mockSetup func(mock sqlmock.Sqlmock, order *model.StockOrder) {  // ← Синтаксична помилка!
        mock.ExpectQuery(...).
            WithArgs(...).
            WillReturnRows(...)
    }
}
```

**Причина:** В Go не можна мати тіло функції прямо в struct field definition.

---

## ✅ Варіант 1: Function-Driven (поточний)

```go
testCases := []struct {
    name      string
    order     *model.StockOrder
    mockSetup func(sqlmock.Sqlmock, *model.StockOrder)  // ← Type
}{
    {
        name: "Success",
        order: &model.StockOrder{...},
        mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {  // ← Implementation
            mock.ExpectQuery(...).
                WithArgs(order.UserID, order.Username, ...).
                WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
        },
    },
}
```

**Переваги:**
- ✅ Гнучкість - можна будь-яка логіка
- ✅ Кастомізація - різні setups для різних кейсів

**Недоліки:**
- ❌ Багатослівно
- ❌ Багато коду в кожному кейсі

---

## ✅ Варіант 2: Data-Driven (альтернатива)

Замість функції використовуємо дані:

```go
testCases := []struct {
    name        string
    order       *model.StockOrder
    returnID    int64        // ← Дані замість функції
    returnError error        // ← Дані замість функції
    expectError bool
}{
    {
        name: "Success",
        order: &model.StockOrder{
            UserID:    1,
            Username:  "john",
            Symbol:    "AAPL",
            OrderType: model.OrderTypeBid,
            Price:     150.00,
            Quantity:  10,
            Status:    "open",
        },
        returnID:    123,      // ← Просто число!
        returnError: nil,
        expectError: false,
    },
    {
        name:        "Database error",
        order:       &model.StockOrder{...},
        returnError: errors.New("connection lost"),  // ← Просто error!
        expectError: true,
    },
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        repo, mock, cleanup := setupRepository(t)
        defer cleanup()

        // Generic mock setup based on data
        if tc.returnError != nil {
            mock.ExpectQuery(`INSERT INTO stock_orders ...`).
                WithArgs(tc.order.UserID, tc.order.Username, ...).
                WillReturnError(tc.returnError)
        } else {
            mock.ExpectQuery(`INSERT INTO stock_orders ...`).
                WithArgs(tc.order.UserID, tc.order.Username, ...).
                WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tc.returnID))
        }

        ctx := context.Background()
        err := repo.Create(ctx, tc.order)

        if tc.expectError {
            assertError(t, err)
        } else {
            assertNoError(t, err)
        }
    })
}
```

**Переваги:**
- ✅ Менше коду
- ✅ Дані в одному місці
- ✅ Немає функцій в test cases

**Недоліки:**
- ❌ Менш гнучко для складних кейсів
- ❌ Mock setup в циклі (а не в test case)

---

## 🎯 Гібридний підхід (найкращий)

Комбінуємо обидва - helper functions + простий data:

```go
// Helper для стандартних mock setups
func expectCreateSuccess(mock sqlmock.Sqlmock, order *model.StockOrder, returnID int64) {
    mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
        WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnID))
}

func expectCreateError(mock sqlmock.Sqlmock, order *model.StockOrder, err error) {
    mock.ExpectQuery(`INSERT INTO stock_orders`).
        WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
        WillReturnError(err)
}

// Використання:
testCases := []struct {
    name        string
    order       *model.StockOrder
    mockSetup   func(sqlmock.Sqlmock, *model.StockOrder)
    expectError bool
}{
    {
        name:  "Success",
        order: &model.StockOrder{...},
        mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
            expectCreateSuccess(mock, order, 123)  // ← Одна лінія!
        },
    },
    {
        name:  "Database error",
        order: &model.StockOrder{...},
        mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
            expectCreateError(mock, order, errors.New("connection lost"))  // ← Одна лінія!
        },
    },
}
```

**Переваги:**
- ✅ Короткий mockSetup (1 лінія)
- ✅ Гнучкість (можна комбінувати helpers)
- ✅ Зрозуміло (helper name каже що робить)

---

## 📚 Реальний приклад

### Створимо helpers:

```go
// internal/repository/stock_order_repository_test_helpers.go

package repository

import (
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"stock_hub/pkg/model"
)

// expectCreateSuccess sets up mock for successful Create
func expectCreateSuccess(mock sqlmock.Sqlmock, order *model.StockOrder, returnID int64) {
	mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
		WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnID))
}

// expectCreateError sets up mock for Create that fails
func expectCreateError(mock sqlmock.Sqlmock, order *model.StockOrder, err error) {
	mock.ExpectQuery(`INSERT INTO stock_orders`).
		WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
		WillReturnError(err)
}

// expectGetByIDSuccess sets up mock for successful GetByID
func expectGetByIDSuccess(mock sqlmock.Sqlmock, orderID int64, order *model.StockOrder) {
	rows := sqlmock.NewRows(orderColumns()).
		AddRow(order.ID, order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, order.CreatedAt)
	mock.ExpectQuery(`SELECT ... WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnRows(rows)
}

// expectGetByIDError sets up mock for GetByID that fails
func expectGetByIDError(mock sqlmock.Sqlmock, orderID int64, err error) {
	mock.ExpectQuery(`SELECT ... WHERE id = \$1`).
		WithArgs(orderID).
		WillReturnError(err)
}

// expectGetOpenOrdersSuccess sets up mock for successful GetOpenOrders
func expectGetOpenOrdersSuccess(mock sqlmock.Sqlmock, symbol string, orders ...*model.StockOrder) {
	rows := sqlmock.NewRows(orderColumns())
	for _, order := range orders {
		rows.AddRow(order.ID, order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, order.CreatedAt)
	}
	mock.ExpectQuery(`SELECT ... WHERE symbol = \$1 AND status = 'open'`).
		WithArgs(symbol).
		WillReturnRows(rows)
}

// expectCancelSuccess sets up mock for successful Cancel
func expectCancelSuccess(mock sqlmock.Sqlmock, orderID, userID int64) {
	mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
		WithArgs(orderID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

// expectCancelNotFound sets up mock for Cancel when order not found
func expectCancelNotFound(mock sqlmock.Sqlmock, orderID, userID int64) {
	mock.ExpectExec(`UPDATE stock_orders SET status = 'cancelled' WHERE id = \$1 AND user_id = \$2 AND status = 'open'`).
		WithArgs(orderID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
}
```

---

### Використання helpers:

```go
func TestCreate(t *testing.T) {
	testCases := []struct {
		name        string
		order       *model.StockOrder
		mockSetup   func(sqlmock.Sqlmock, *model.StockOrder)
		expectError bool
	}{
		{
			name: "Success",
			order: &model.StockOrder{
				UserID:    1,
				Username:  "john",
				Symbol:    "AAPL",
				OrderType: model.OrderTypeBid,
				Price:     150.00,
				Quantity:  10,
				Status:    "open",
				CreatedAt: time.Now(),
			},
			mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
				expectCreateSuccess(mock, order, 123)  // ← 1 лінія!
			},
			expectError: false,
		},
		{
			name:  "Database error",
			order: &model.StockOrder{...},
			mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
				expectCreateError(mock, order, errors.New("connection lost"))  // ← 1 лінія!
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			tc.mockSetup(mock, tc.order)

			ctx := context.Background()
			err := repo.Create(ctx, tc.order)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
			}

			assertMockExpectations(t, mock)
		})
	}
}
```

---

## 🔄 Ще компактніше: Data-Only

Якщо хочете зовсім прибрати `mockSetup`, можна так:

```go
func TestCreate(t *testing.T) {
	testCases := []struct {
		name        string
		order       *model.StockOrder
		returnID    int64   // ← Дані
		returnError error   // ← Дані
		expectError bool
	}{
		{
			name: "Success",
			order: &model.StockOrder{
				UserID:    1,
				Username:  "john",
				Symbol:    "AAPL",
				OrderType: model.OrderTypeBid,
				Price:     150.00,
				Quantity:  10,
				Status:    "open",
			},
			returnID:    123,  // ← Просто число
			returnError: nil,
			expectError: false,
		},
		{
			name:        "Database error",
			order:       &model.StockOrder{...},
			returnError: errors.New("connection lost"),  // ← Просто error
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock, cleanup := setupRepository(t)
			defer cleanup()

			// Generic setup based on data
			if tc.returnError != nil {
				expectCreateError(mock, tc.order, tc.returnError)
			} else {
				expectCreateSuccess(mock, tc.order, tc.returnID)
			}

			ctx := context.Background()
			err := repo.Create(ctx, tc.order)

			if tc.expectError {
				assertError(t, err)
			} else {
				assertNoError(t, err)
				if tc.order.ID != tc.returnID {
					t.Errorf("Expected ID %d, got: %d", tc.returnID, tc.order.ID)
				}
			}

			assertMockExpectations(t, mock)
		})
	}
}
```

**Переваги:**
- ✅ Тільки дані в test cases
- ✅ Жодних функцій
- ✅ Найкомпактніше

**Недоліки:**
- ❌ Менш гнучко для складних кейсів
- ❌ Mock setup в циклі (а не біля даних)

---

## 📊 Порівняння всіх підходів

### 1. Inline Functions (багато коду):

```go
{
    name: "Success",
    mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
        mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
            WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))
    },
}

Lines: 5
Readability: ⚠️ Medium
Flexibility: ✅ High
```

---

### 2. Helper Functions (збалансовано):

```go
{
    name: "Success",
    mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
        expectCreateSuccess(mock, order, 123)  // ← 1 лінія
    },
}

Lines: 1
Readability: ✅ High
Flexibility: ✅ High
```

---

### 3. Data-Driven (мінімум коду):

```go
{
    name:     "Success",
    order:    &model.StockOrder{...},
    returnID: 123,  // ← Тільки дані
}

Lines: 1 (but mock setup in loop)
Readability: ✅ Very High
Flexibility: ⚠️ Medium
```

---

## 💡 Рекомендація: Helper Functions

Створимо helper functions для спрощення:

```go
// stock_order_repository_test_helpers.go

// expectCreateSuccess sets up mock for successful Create
func expectCreateSuccess(mock sqlmock.Sqlmock, order *model.StockOrder, returnID int64) {
    mock.ExpectQuery(`INSERT INTO stock_orders \(user_id, username, symbol, order_type, price, quantity, status, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) RETURNING id`).
        WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(returnID))
}

// expectCreateError sets up mock for failed Create
func expectCreateError(mock sqlmock.Sqlmock, order *model.StockOrder, err error) {
    mock.ExpectQuery(`INSERT INTO stock_orders`).
        WithArgs(order.UserID, order.Username, order.Symbol, string(order.OrderType), order.Price, order.Quantity, order.Status, sqlmock.AnyArg()).
        WillReturnError(err)
}
```

**Використання:**
```go
{
    name: "Success",
    mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
        expectCreateSuccess(mock, order, 123)  // ← Clean!
    },
}
```

---

## 🎯 Фінальний варіант

Найкраща комбінація:

1. ✅ **Helper functions** для типових setups
2. ✅ **Параметризований mockSetup** для гнучкості
3. ✅ **orderColumns()** helper
4. ✅ **Query constants**

```go
func TestCreate(t *testing.T) {
    testCases := []struct {
        name        string
        order       *model.StockOrder
        mockSetup   func(sqlmock.Sqlmock, *model.StockOrder)
        expectError bool
    }{
        {
            name:  "Success",
            order: &model.StockOrder{...},
            mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
                expectCreateSuccess(mock, order, 123)  // ← Helper
            },
        },
        {
            name:  "Error",
            order: &model.StockOrder{...},
            mockSetup: func(mock sqlmock.Sqlmock, order *model.StockOrder) {
                expectCreateError(mock, order, errors.New("failed"))  // ← Helper
            },
        },
    }
}
```

---

## ⚖️ Висновок

| Підхід | Lines/case | Гнучкість | Читабельність |
|--------|-----------|-----------|---------------|
| Inline Functions | 5-7 | ✅ High | ⚠️ Medium |
| Helper Functions | 1-2 | ✅ High | ✅ High |
| Data-Driven | 1 | ⚠️ Medium | ✅ High |

**Рекомендація: Helper Functions** - найкращий баланс! ✅

