# EcoLokal - Bank Sampah API

EcoLokal adalah REST API untuk mengelola ekosistem Bank Sampah di tingkat RT/RW atau kelurahan, dengan sistem penjemputan sampah terjadwal dan manajemen poin reward.

## 🌟 Fitur Utama

- **Manajemen Pengguna dengan 3 Role**: Warga, Petugas, Admin
- **Katalog Jenis Sampah**: Dengan sistem poin per kilogram
- **Sistem Penjadwalan Penjemputan**: Request dan assign ke petugas
- **Sistem Poin Reward**: Otomatis dihitung berdasarkan berat aktual
- **Riwayat Transaksi**: Tracking poin masuk dan keluar

## 🛠️ Tech Stack

- **Go 1.21+**
- **Gin Web Framework**
- **PostgreSQL 12+**
- **JWT Authentication**
- **Swagger Documentation**
- **Docker & Docker Compose**

## 📋 Prerequisites

- **Go 1.21 or higher** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional) - [Download Docker](https://www.docker.com/get-started)
- **Make** (optional, for using Makefile commands)

## 🚀 Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/alifsuryadi/ecolokal.git
cd ecolokal
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Install Development Tools
```bash
# Install Swagger CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Optional: Install linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 4. Setup Database

#### Option A: Using Docker (Recommended)
```bash
# Start PostgreSQL database with Docker Compose
docker-compose up -d db

# Wait for database to be ready, then run migrations
docker exec -i ecolokal-db psql -U postgres -d postgres < setup.sql
docker exec -i ecolokal-db psql -U postgres -d ecolokal < migrations/001_initial_schema.sql
```

#### Option B: Local PostgreSQL
```bash
# Create database
psql -U postgres -c "CREATE DATABASE ecolokal;"

# Run migrations
psql -h localhost -p 5432 -U postgres -d ecolokal -f migrations/001_initial_schema.sql
```

### 5. Configure Environment
Copy the environment configuration:
```bash
cp .env.example .env
```

Edit `.env` file with your configuration:
```env
# Application Configuration
APP_PORT=8080
APP_ENV=development

# Database Configuration  
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=ecolokal
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24
```

### 6. Generate API Documentation
```bash
# Generate swagger documentation
swag init -g cmd/api/main.go
```

### 7. Build and Run
```bash
# Build the application
go build -o bin/ecolokal.exe ./cmd/api

# Run the application
./bin/ecolokal.exe
```

Or using Makefile:
```bash
# Build and run
make dev

# Or step by step
make build
make run
```

## 🐳 Docker Development

### Full Stack with Docker Compose
```bash
# Start database and application
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app
```

### Database Only
```bash
# Start only database
make docker-up

# Stop database
make docker-down
```

## 📚 API Documentation

After starting the application, access the Swagger documentation at:
- **Swagger UI**: http://localhost:8080/api/docs/index.html
- **Swagger JSON**: http://localhost:8080/api/docs/doc.json

## 🔧 Development Commands

### Using Makefile
```bash
make help                # Show all available commands
make build               # Build the application
make run                 # Run the application
make dev                 # Build and run
make test                # Run tests
make swagger             # Generate swagger docs
make clean               # Clean build artifacts
make docker-up           # Start database with Docker
make docker-down         # Stop database
make docker-full         # Start full stack
make migrate             # Run database migrations
make deps                # Download dependencies
make fmt                 # Format code
make lint                # Run linter
```

### Manual Commands
```bash
# Build
go build -o bin/ecolokal.exe ./cmd/api

# Run tests
go test -v ./...

# Format code
go fmt ./...

# Generate docs
swag init -g cmd/api/main.go

# Clean modules
go clean -modcache
```

## 📁 Project Structure

```
ecolokal/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── domain/                  # Domain models
│   ├── repository/              # Data access layer
│   ├── usecase/                 # Business logic layer
│   └── delivery/
│       └── http/
│           ├── handler/         # HTTP handlers
│           ├── middleware/      # HTTP middlewares
│           └── router/          # Route configuration
├── pkg/
│   ├── database/                # Database connection
│   ├── utils/                   # Utility functions
│   └── validator/               # Custom validators
├── migrations/                  # Database migrations
├── docs/                        # Auto-generated swagger docs
├── bin/                         # Compiled binaries
├── docker-compose.yml           # Docker services configuration
├── Dockerfile                   # Docker image configuration
├── Makefile                     # Development commands
├── .env                         # Environment variables
├── .env.example                 # Environment template
├── setup.sql                    # Database setup script
└── README.md                    # This file
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### Default Users
After running migrations, you'll have:
- **Admin**: `admin@ecolokal.com` (password needs to be hashed and set in migration)

## 📊 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Users
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update user profile
- `GET /api/users/points` - Get user points

### Pickups (Warga)
- `POST /api/pickups` - Create pickup request
- `GET /api/pickups` - Get user's pickup requests

### Pickups (Petugas)
- `GET /api/pickups/assigned` - Get assigned pickups
- `PUT /api/pickups/:id/status` - Update pickup status
- `PUT /api/pickups/:id/items` - Update pickup items with actual weights

### Admin
- `GET /api/admin/pickups` - Get all pickup requests
- `PUT /api/admin/pickups/:id/assign` - Assign pickup to petugas
- `POST /api/admin/waste-types` - Create waste type
- `PUT /api/admin/waste-types/:id` - Update waste type
- `DELETE /api/admin/waste-types/:id` - Delete waste type

### Waste Types
- `GET /api/waste-types` - Get all waste types

### Transactions
- `GET /api/transactions` - Get user transactions
- `POST /api/transactions` - Create transaction (point redemption)

## 🧪 Testing

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/usecase
```

## 🚀 Deployment

### Environment Variables for Production
```env
APP_ENV=production
JWT_SECRET=your-very-secure-secret-key
DB_HOST=your-production-db-host
DB_PASSWORD=your-secure-db-password
```

### Docker Production
```bash
# Build production image
docker build -t ecolokal:latest .

# Run with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Database Migration in Production
```bash
# Run migrations
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

## 📝 Development Guidelines

### Code Structure
- **Domain Layer**: Business entities and rules
- **Repository Layer**: Data access abstraction
- **Usecase Layer**: Business logic implementation
- **Handler Layer**: HTTP request/response handling

### Adding New Features
1. Define domain models in `internal/domain/`
2. Create repository interface and implementation
3. Implement business logic in usecase
4. Add HTTP handlers
5. Update routes in `internal/delivery/http/router/`
6. Add swagger documentation
7. Write tests

### Database Changes
1. Create new migration file in `migrations/`
2. Update domain models if needed
3. Update repository methods
4. Test thoroughly

## 🐛 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL is running
   docker-compose ps
   
   # Check logs
   docker-compose logs db
   ```

2. **Build Errors**
   ```bash
   # Clean and rebuild
   go clean -modcache
   go mod download
   go mod tidy
   ```

3. **Port Already in Use**
   ```bash
   # Change APP_PORT in .env file or kill the process
   netstat -tulpn | grep :8080
   ```

4. **Swagger Documentation Not Generated**
   ```bash
   # Install swag and regenerate
   go install github.com/swaggo/swag/cmd/swag@latest
   swag init -g cmd/api/main.go
   ```

### Logs
- Application logs are written to stdout
- Database logs: `docker-compose logs db`
- Application logs: `docker-compose logs app`

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -am 'Add new feature'`)
4. Push to branch (`git push origin feature/new-feature`)
5. Create Pull Request

### Code Style
- Use `go fmt` for formatting
- Follow Go naming conventions
- Add comments for exported functions
- Write unit tests for business logic

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

- **Alif Suryadi** - Developer

## 📞 Support

For support and questions:
- Email: support@ecolokal.com
- Issues: [GitHub Issues](https://github.com/alifsuryadi/ecolokal/issues)

## 🔄 Changelog

### v1.0.0
- Initial release
- User authentication and authorization
- Pickup request management
- Waste type management
- Point system
- Transaction history
- Swagger API documentation
