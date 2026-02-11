# Monorepo Migration Guide 🔄

Міграція до monorepo структури для Stock Hub.

---

## 📦 Що змінилось?

### До (Before)

```
/Users/vkuzm/GolandProjects/
├── stock_hub/              # Trading service
│   ├── pkg/model/         # Shared models
│   └── ...
│
└── email_service/          # Email service (окремо)
    └── ...
```

### Після (After)

```
/Users/vkuzm/GolandProjects/
└── stock_hub/                      # ← MONOREPO
    ├── go.work                     # Go workspace
    ├── README.md                   # Main documentation
    │
    ├── stock_hub_trade/            # Trading service (renamed from stock_hub)
    │   ├── pkg/model/             # Shared models
    │   ├── api/proto/             # gRPC protobuf
    │   └── go.mod                 # module stock_hub_trade
    │
    └── stock_hub_email_service/    # Email service (moved from ../email_service)
        ├── go.mod                  # module email_service
        └── ...                     # replace stock_hub_trade => ../stock_hub_trade
```

---

## 🎯 Переваги Monorepo

### ✅ Організація

- Всі сервіси в одній папці
- Єдина git історія
- Синхронізація версій

### ✅ Розробка

- `go.work` для роботи з декількома модулями
- Shared models (`stock_hub_trade/pkg/model/`)
- Простіше дебажити

### ✅ CI/CD

- Одна конфігурація для всіх сервісів
- Синхронні релізи
- Спільні dependencies

---

## 🔧 Технічні зміни

### 1. Go Module Names

```go
// stock_hub_trade/go.mod
module stock_hub_trade  // Було: module stock_hub

// stock_hub_email_service/go.mod
module email_service
require stock_hub_trade v0.0.0
replace stock_hub_trade => ../stock_hub_trade  // Було: stock_hub => ../stock_hub
```

### 2. Import Paths

**Всі файли оновлені:**

```go
// До
import "stock_hub/pkg/model"
import "stock_hub/internal/service"

// Після
import "stock_hub_trade/pkg/model"
import "stock_hub_trade/internal/service"
```

### 3. Go Workspace

**Файл `go.work`:**

```go
go 1.25

use (
	./stock_hub_trade
	./stock_hub_email_service
)
```

**Команди:**

```bash
# Синхронізувати dependencies
go work sync

# Запустити тести в усіх модулях
go test ./...

# Оновити dependencies
go get -u ./...
```

---

## 📝 Git History

### Збережено через `git mv`

```bash
# Git розпізнав переміщення файлів:
R  CHANGELOG.md -> stock_hub_trade/CHANGELOG.md
R  Makefile -> stock_hub_trade/Makefile
R  cmd/server/main.go -> stock_hub_trade/cmd/server/main.go
# ... і всі інші файли
```

**Історія зберігається!**

```bash
# Подивитись історію файлу після переміщення:
git log --follow stock_hub_trade/cmd/server/main.go
```

---

## 🚀 Як використовувати

### Клонування

```bash
# Клонувати monorepo
git clone <repo> stock_hub
cd stock_hub

# Go workspace вже налаштовано (go.work)
```

### Розробка

```bash
# 1. Запустити trading service
cd stock_hub_trade
make run-postgres

# 2. В іншому терміналі запустити email service
cd ../stock_hub_email_service
make run
```

### Тестування

```bash
# З root directory - всі тести
go test ./...

# Окремо trading service
cd stock_hub_trade && go test ./...

# Окремо email service
cd stock_hub_email_service && go test ./...
```

### Build

```bash
# Trading service
cd stock_hub_trade && go build ./cmd/server/

# Email service
cd stock_hub_email_service && go build ./cmd/server/
```

---

## 🔄 Міграція для розробників

### Якщо ви маєте локальну копію старої структури:

```bash
# 1. Backup старих папок
cp -r /path/to/old/stock_hub /path/to/backup/
cp -r /path/to/old/email_service /path/to/backup/

# 2. Pull нову структуру
cd /path/to/old/stock_hub
git pull origin develop

# 3. Оновити IDE workspace
# GoLand: File → Open → /path/to/stock_hub
# VSCode: Відкрити /path/to/stock_hub

# 4. Синхронізувати dependencies
go work sync
```

### Оновлення import paths (якщо є власні зміни):

```bash
# Automatic find & replace в trading service
cd stock_hub_trade
find . -name "*.go" -type f -exec sed -i '' 's|"stock_hub/|"stock_hub_trade/|g' {} \;

# В email service
cd ../stock_hub_email_service
find . -name "*.go" -type f -exec sed -i '' 's|"stock_hub/|"stock_hub_trade/|g' {} \;
```

---

## 📚 Документація

### Root level (monorepo)
- `README.md` - Загальний опис monorepo
- `MONOREPO_MIGRATION.md` - Цей файл

### Trading Service
- `stock_hub_trade/README.md` - Trading service docs
- `stock_hub_trade/docs/` - Детальна документація
- `stock_hub_trade/QUICKSTART.md` - Швидкий старт

### Email Service
- `stock_hub_email_service/README.md` - Email service docs

---

## 🧪 Перевірка після міграції

### 1. Build

```bash
cd /Users/vkuzm/GolandProjects/stock_hub
cd stock_hub_trade && go build ./cmd/server/
cd ../stock_hub_email_service && go build ./cmd/server/
```

**Очікується:** Успішний build обох сервісів.

### 2. Tests

```bash
cd /Users/vkuzm/GolandProjects/stock_hub
go test ./...
```

**Очікується:** Всі тести проходять.

### 3. Run Services

```bash
# Terminal 1: Trading service
cd stock_hub_trade
make run-postgres

# Terminal 2: Email service
cd stock_hub_email_service
make run
```

**Очікується:** Обидва сервіси запускаються без помилок.

---

## ⚠️ Важливо

### .env файли

**Кожен сервіс має власний `.env`:**

```bash
stock_hub_trade/.env           # Trading service config
stock_hub_email_service/.env   # Email service config
```

### Git

**Структура в git:**

```bash
stock_hub/
├── .git/                      # Єдина git історія
├── stock_hub_trade/
└── stock_hub_email_service/
```

**Не створюйте окремі git repos в subdirectories!**

### IDE

**GoLand/VSCode:**

- Відкривати root папку: `/Users/vkuzm/GolandProjects/stock_hub`
- Go workspace автоматично розпізнається через `go.work`

---

## 🎓 Навчальні матеріали

- [Monorepo Best Practices](https://earthly.dev/blog/golang-monorepo/)
- [Go Workspaces](https://go.dev/doc/tutorial/workspaces)
- [Google's Monorepo](https://research.google/pubs/pub45424/)

---

## 📞 Підтримка

**Питання?**

1. Читайте `README.md` в root
2. Читайте `stock_hub_trade/README.md`
3. Дивіться документацію в `stock_hub_trade/docs/`

---

## ✅ Checklist після міграції

- [ ] `git pull` успішно завершено
- [ ] `go work sync` виконано
- [ ] `go build` працює для обох сервісів
- [ ] `go test ./...` проходить
- [ ] `.env` файли налаштовані
- [ ] IDE workspace оновлено
- [ ] Сервіси запускаються

**Якщо все ✅ - міграція успішна!** 🎉
