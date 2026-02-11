.PHONY: help run build test clean

help: ## Show help
	@echo "Email Service - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run email service
	@echo "📧 Starting Email Service..."
	go run cmd/server/main.go

build: ## Build binary
	@echo "🔨 Building email service..."
	go build -o bin/email-service cmd/server/main.go
	@echo "✅ Binary created: bin/email-service"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f *.log
	@echo "✅ Clean complete"

tidy: ## Tidy dependencies
	@echo "📦 Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies tidied"

.DEFAULT_GOAL := help
