# 🚀 Швидкий старт Stock Hub

## За 3 хвилини до запуску!

### 1️⃣ Запустіть PostgreSQL в Docker
```bash
make docker-up
```

### 2️⃣ Застосуйте міграції
```bash
make migrate-up
```

### 3️⃣ Запустіть сервер
```bash
make run
```

**Готово!** 🎉

Відкрийте в браузері: http://localhost:8082

---

## Корисні команди

```bash
make help          # Показати всі команди
make psql          # Підключитись до PostgreSQL
make db-reset      # Очистити БД та застосувати міграції заново
make docker-down   # Зупинити PostgreSQL
make docker-logs   # Логи PostgreSQL
```

---

## Що робити якщо щось не працює?

### Порт 8082 вже зайнятий
```bash
# Знайдіть процес що використовує порт
lsof -ti:8082 | xargs kill -9
```

### PostgreSQL не запускається
```bash
# Зупиніть контейнер і запустіть знову
make docker-down
make docker-up
```

### Міграції не застосовуються
```bash
# Скиньте БД повністю
make db-reset
```

### База даних не існує
```bash
# Створіть вручну
psql -U postgres -c "CREATE DATABASE stock_hub;"
```

---

## Налаштування

Скопіюйте `.env.example` в `.env` та відредагуйте при необхідності:

```bash
cp .env.example .env
```

Стандартна конфігурація:
- **Порт**: 8082
- **PostgreSQL**: localhost:5432
- **База**: stock_hub
- **Користувач**: postgres
- **Пароль**: postgres

---

**Більше інформації**: [README.md](./README.md)
