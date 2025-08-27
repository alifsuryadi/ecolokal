#!/bin/bash

# EcoLokal Quick Start Script

echo "🌱 EcoLokal Quick Start"
echo "======================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Install swag if not present
if ! command -v swag &> /dev/null; then
    echo "📚 Installing Swagger CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate swagger docs
echo "📖 Generating API documentation..."
swag init -g cmd/api/main.go

# Build application
echo "🔨 Building application..."
go build -o bin/ecolokal ./cmd/api

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "⚙️  Creating .env file from template..."
    cp .env.example .env
    echo "📝 Please edit .env file with your database configuration"
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Setup PostgreSQL database"
echo "2. Edit .env file with your database configuration"
echo "3. Run database migrations:"
echo "   psql -h localhost -p 5432 -U postgres -c \"CREATE DATABASE ecolokal;\""
echo "   psql -h localhost -p 5432 -U postgres -d ecolokal -f migrations/001_initial_schema.sql"
echo "4. Start the application:"
echo "   ./bin/ecolokal"
echo ""
echo "Or use Docker:"
echo "   docker-compose up -d db"
echo "   docker exec -i ecolokal-db psql -U postgres -d postgres < setup.sql"
echo "   docker exec -i ecolokal-db psql -U postgres -d ecolokal < migrations/001_initial_schema.sql"
echo "   ./bin/ecolokal"
echo ""
echo "API Documentation: http://localhost:8080/api/docs/index.html"