.PHONY: help build run test clean migrate-up migrate-down migrate-status migrate-create seed docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building application..."
	go build -o flypro-assessment ./cmd/server

run: ## Run the application
	@echo "Running application..."
	go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test -v -coverprofile=coverage.out ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f flypro-assessment
	rm -f coverage.out coverage.html

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@export PATH=$$PATH:$$HOME/go/bin && goose -dir migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@export PATH=$$PATH:$$HOME/go/bin && goose -dir migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" down

migrate-status: ## Show migration status
	@echo "Migration status..."
	@export PATH=$$PATH:$$HOME/go/bin && goose -dir migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" status

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)..."
	@export PATH=$$PATH:$$HOME/go/bin && goose -dir migrations create $(NAME) sql

seed: ## Seed database with sample data
	@echo "Seeding database..."
	go run scripts/seed.go

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

install-deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go install github.com/pressly/goose/v3/cmd/goose@latest

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vendor: ## Vendor dependencies
	@echo "Vendoring dependencies..."
	go mod vendor
