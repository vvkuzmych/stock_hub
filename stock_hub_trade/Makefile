.PHONY: help run migrate-up migrate-down docker-up docker-down clean

help: ## Показати допомогу
	@echo "Stock Hub - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Запустити сервер (PostgreSQL)
	@echo "🚀 Starting Stock Hub with PostgreSQL..."
	DB_TYPE=postgres go run cmd/server/main.go

migrate-up: ## Застосувати міграції
	@echo "📊 Running migrations..."
	DB_TYPE=postgres go run cmd/migrate/main.go -command=up

migrate-down: ## Відкотити останню міграцію
	@echo "↩️  Rolling back last migration..."
	DB_TYPE=postgres go run cmd/migrate/main.go -command=down

migrate-down-all: ## Відкотити ВСІ міграції
	@echo "⚠️  Rolling back ALL migrations..."
	@DB_TYPE=postgres go run cmd/migrate/main.go -command=down 2>/dev/null || true
	@DB_TYPE=postgres go run cmd/migrate/main.go -command=down 2>/dev/null || true
	@DB_TYPE=postgres go run cmd/migrate/main.go -command=down 2>/dev/null || true
	@echo "✅ All migrations rolled back"

docker-up: ## Запустити PostgreSQL в Docker
	@echo "🐳 Starting PostgreSQL in Docker..."
	docker-compose up -d
	@echo "✅ PostgreSQL is running on localhost:5432"

docker-down: ## Зупинити PostgreSQL Docker
	@echo "🛑 Stopping PostgreSQL..."
	docker-compose down

docker-logs: ## Показати логи PostgreSQL
	docker-compose logs -f postgres

install-deps: ## Встановити залежності
	@echo "📦 Installing dependencies..."
	go get github.com/lib/pq
	go get github.com/gorilla/websocket
	go mod tidy
	@echo "✅ Dependencies installed"

setup-postgres: docker-up ## Повне налаштування PostgreSQL
	@echo "⏳ Waiting for PostgreSQL to start..."
	@sleep 3
	@echo "✅ PostgreSQL is ready!"
	@echo "📝 Create .env file:"
	@cp .env.example .env 2>/dev/null || echo ".env already exists"
	@echo "✅ Setup complete! Run: make run"

clean: ## Очистити тимчасові файли
	@echo "🧹 Cleaning..."
	rm -f stock_hub.db
	rm -f *.log
	@echo "✅ Clean complete"

test: ## Запустити тести
	@echo "🧪 Running tests..."
	go test -v ./...

build: ## Зібрати binary
	@echo "🔨 Building..."
	go build -o bin/stock_hub cmd/server/main.go
	@echo "✅ Binary created: bin/stock_hub"

psql: ## Підключитись до PostgreSQL
	@echo "🔌 Connecting to PostgreSQL..."
	psql -U postgres -d stock_hub

db-reset: ## Повністю очистити PostgreSQL БД та застосувати міграції
	@echo "🔄 Resetting PostgreSQL database..."
	@psql -U postgres -d stock_hub -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" > /dev/null 2>&1
	@echo "✅ Database reset"
	@echo "📊 Applying migrations..."
	@make migrate-up
	@echo "✅ Database ready!"

.DEFAULT_GOAL := help
