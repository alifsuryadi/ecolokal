# EcoLokal Makefile

.PHONY: build run dev test clean docker-up docker-down migrate help

# Build the application
build:
	@echo "Building EcoLokal..."
	@go build -o bin/ecolokal.exe ./cmd/api

# Run the application
run:
	@echo "Running EcoLokal..."
	@./bin/ecolokal.exe

# Build and run
dev: build run

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Generate swagger docs
swagger:
	@echo "Generating swagger docs..."
	@swag init -g cmd/api/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf docs/

# Start database with Docker Compose
docker-up:
	@echo "Starting database..."
	@docker-compose up -d db

# Stop database
docker-down:
	@echo "Stopping database..."
	@docker-compose down

# Start full stack with Docker Compose
docker-full:
	@echo "Starting full application stack..."
	@docker-compose up -d

# Migrate database (requires PostgreSQL connection)
migrate:
	@echo "Running database migrations..."
	@psql -h localhost -p 5432 -U postgres -d ecolokal -f migrations/001_initial_schema.sql

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  dev          - Build and run the application"
	@echo "  test         - Run tests"
	@echo "  swagger      - Generate swagger documentation"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-up    - Start database with Docker"
	@echo "  docker-down  - Stop database"
	@echo "  docker-full  - Start full stack with Docker"
	@echo "  migrate      - Run database migrations"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  install-tools- Install development tools"
	@echo "  help         - Show this help"