# Travel Expense Management API

A comprehensive Travel Expense Management API built with Go, Gin, GORM, PostgreSQL, and Redis. This API demonstrates backend architecture best practices including database design, third-party integrations, caching strategies, and comprehensive testing.

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Setup Instructions](#setup-instructions)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Deployment](#deployment)
- [Architecture Decisions](#architecture-decisions)

## ✨ Features

- **User Management**: Create and retrieve users
- **Expense Management**: Full CRUD operations with validation
- **Expense Reports**: Create, manage, and submit expense reports
- **Currency Conversion**: Automatic conversion to USD with Redis caching (6-hour TTL)
- **Caching Strategy**: Redis caching for users (1-hour), reports (30-min), and exchange rates (6-hour)
- **Database Migrations**: Goose migrations for version control
- **Comprehensive Validation**: Input validation with custom validators
- **Error Handling**: Structured error responses
- **Logging**: Structured logging with correlation IDs
- **Graceful Shutdown**: Context-aware shutdown handling

## 🛠 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Migrations**: Goose
- **Validation**: go-playground/validator
- **Logging**: Zap
- **Testing**: Testify with mocks

## 🏗 Architecture

### Design Patterns

- **Repository Pattern**: Abstraction layer for data access
- **Service Layer**: Business logic separation
- **DTO Pattern**: Request/Response transformation
- **Dependency Injection**: Loose coupling between components
- **Middleware Chain**: Reusable cross-cutting concerns

### Architecture Layers

```
┌─────────────────────────────────────────┐
│           HTTP Handlers                  │
│     (Request/Response, Validation)        │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│           Service Layer                  │
│     (Business Logic, Orchestration)       │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│        Repository Layer                   │
│     (Data Access, GORM Queries)          │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│        Database (PostgreSQL)              │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│     External Services                     │
│  (Currency API, Redis Cache)              │
└─────────────────────────────────────────┘
```

## 📁 Project Structure

```
flypro-assessment/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── handlers/                  # HTTP handlers
│   ├── services/                   # Business logic layer
│   │   └── mocks/                  # Service mocks for testing
│   ├── repository/                 # Data access layer
│   │   └── mocks/                  # Repository mocks
│   ├── models/                     # Database models
│   ├── dto/                        # Data Transfer Objects
│   ├── validators/                 # Custom validators
│   ├── middleware/                 # HTTP middleware
│   └── utils/                      # Utility functions
├── migrations/                     # Goose migration files
├── scripts/                        # Utility scripts (seed, etc.)
├── tests/                          # Test files
├── docker-compose.yml              # Docker services configuration
├── Dockerfile                      # Application container
├── Makefile                        # Build and deployment commands
├── go.mod                          # Go module dependencies
└── README.md                       # This file
```

## 🚀 Setup Instructions

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Redis 7 or higher
- Make (optional but recommended)
- Docker and Docker Compose (for containerized setup)

### Local Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/ayoashy/flypro-assessment-ayo.git
   cd flypro-assessment-ayo
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   # or
   go mod download
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start PostgreSQL and Redis**
   ```bash
   # Using Docker Compose
   docker-compose up -d postgres redis
   
   # Or use your local PostgreSQL and Redis instances
   ```

5. **Run database migrations**
   ```bash
   make migrate-up
   # or
   goose -dir migrations postgres "host=localhost port=5432 user=flypro_user password=flypro_password dbname=flypro_db sslmode=disable" up
   ```

6. **Seed database (optional)**
   ```bash
   make seed
   # or
   go run scripts/seed.go
   ```

7. **Build and run the application**
   ```bash
   make build
   ./flypro-assessment
   
   # or run directly
   make run
   ```

### Docker Setup

1. **Build and start all services**
   ```bash
   docker-compose up -d
   ```

2. **Check logs**
   ```bash
   docker-compose logs -f app
   ```

3. **Stop services**
   ```bash
   docker-compose down
   ```

### Environment Variables

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=flypro_user
DB_PASSWORD=flypro_password
DB_NAME=flypro_db
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Currency API Configuration
CURRENCY_API_KEY=your_api_key_here
CURRENCY_API_URL=https://api.exchangerate-api.com/v4/latest

# Application Environment
ENV=development
```

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api
```

### User Management

#### Create User
```http
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "name": "John Doe"
}
```

#### Get User
```http
GET /api/users/:id
```

### Expense Management

#### Create Expense
```http
POST /api/expenses
Content-Type: application/json
X-User-ID: 1

{
  "amount": 100.50,
  "currency": "USD",
  "category": "travel",
  "description": "Flight ticket",
  "receipt": "https://example.com/receipt.jpg"
}
```

#### List Expenses
```http
GET /api/expenses?page=1&per_page=10&category=travel&status=pending
X-User-ID: 1
```

#### Get Expense
```http
GET /api/expenses/:id
```

#### Update Expense
```http
PUT /api/expenses/:id
Content-Type: application/json
X-User-ID: 1

{
  "amount": 150.00,
  "status": "approved"
}
```

#### Delete Expense
```http
DELETE /api/expenses/:id
X-User-ID: 1
```

### Expense Reports

#### Create Report
```http
POST /api/reports
Content-Type: application/json
X-User-ID: 1

{
  "title": "Q1 2024 Travel Expenses"
}
```

#### List Reports
```http
GET /api/reports?page=1&per_page=10&status=draft
X-User-ID: 1
```

#### Get Report
```http
GET /api/reports/:id
```

#### Add Expenses to Report
```http
POST /api/reports/:id/expenses
Content-Type: application/json
X-User-ID: 1

{
  "expense_ids": [1, 2, 3]
}
```

#### Submit Report
```http
PUT /api/reports/:id/submit
X-User-ID: 1
```

### Response Format

#### Success Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "amount": 100.50,
    "currency": "USD"
  }
}
```

#### Error Response
```json
{
  "success": false,
  "error": {
    "type": "VALIDATION_ERROR",
    "message": "amount must be greater than 0",
    "field": "amount"
  }
}
```

#### Paginated Response
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

## 🗄 Database Schema

### Entity Relationship Diagram

```
┌─────────┐
│  Users  │
└────┬────┘
     │
     │ 1:N
     │
┌────▼──────────┐      N:N      ┌───────────────┐
│   Expenses   │◄─────────────►│ ExpenseReports│
└──────────────┘                └───────┬───────┘
                                       │
                                       │ 1:N
                                       │
                                  ┌────▼────┐
                                  │  Users  │
                                  └─────────┘
```

### Tables

#### Users
- `id` (PK)
- `email` (UNIQUE, INDEXED)
- `name`
- `created_at`
- `updated_at`
- `deleted_at` (soft delete)

#### Expenses
- `id` (PK)
- `user_id` (FK, INDEXED)
- `amount`
- `currency`
- `category` (INDEXED)
- `description`
- `receipt`
- `status` (INDEXED)
- `created_at`
- `updated_at`
- `deleted_at` (soft delete)

#### Expense Reports
- `id` (PK)
- `user_id` (FK, INDEXED)
- `title`
- `status` (INDEXED)
- `total`
- `created_at`
- `updated_at`
- `deleted_at` (soft delete)

#### Report Expenses (Join Table)
- `expense_report_id` (FK)
- `expense_id` (FK)
- Composite Primary Key

## 🧪 Testing

### Run Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

### Test Structure

Tests are organized by layer:
- **Service Layer Tests**: Business logic validation with mocks
- **Repository Tests**: Database operations (optional, requires test DB)
- **Handler Tests**: HTTP request/response validation

### Test Coverage

Current test coverage focuses on:
- Service layer business logic
- Request validation
- Error handling
- Mocking patterns

To achieve 10-20% coverage, tests are written for:
1. Expense service (create, get, list, update, delete)
2. Currency service conversion logic
3. Handler request validation

## 🔧 Makefile Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make migrate-up        # Run database migrations
make migrate-down      # Rollback migrations
make migrate-status    # Show migration status
make migrate-create    # Create new migration (NAME=name)
make seed              # Seed database with sample data
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make clean             # Clean build artifacts
```

## 🚢 Deployment

### Docker Deployment

1. **Build the image**
   ```bash
   docker build -t flypro-assessment .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 --env-file .env flypro-assessment
   ```

### Production Considerations

- Set `ENV=production` for production mode
- Use secure database credentials
- Configure Redis with password authentication
- Set up proper logging aggregation
- Implement rate limiting
- Use HTTPS/TLS
- Set up monitoring and alerting

## 🏛 Architecture Decisions

### 1. How would you handle concurrent expense approvals?

**Strategy**: Implement optimistic locking using version fields and database transactions.

```go
// Add version field to Expense model
type Expense struct {
    // ... existing fields
    Version int `gorm:"default:0"`
}

// In service layer
func (s *ExpenseService) ApproveExpense(ctx context.Context, id uint, version int) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        var expense Expense
        if err := tx.Where("id = ? AND version = ?", id, version).
            First(&expense).Error; err != nil {
            return errors.New("expense was modified")
        }
        
        expense.Status = "approved"
        expense.Version++
        return tx.Save(&expense).Error
    })
}
```

**Alternative**: Use database-level locking:
```sql
SELECT * FROM expenses WHERE id = ? FOR UPDATE;
```

### 2. What strategies would you use to scale this system?

**Horizontal Scaling**:
- Stateless API servers (multiple instances behind load balancer)
- Read replicas for PostgreSQL
- Redis cluster for distributed caching
- Sharding by user_id for very large datasets

**Vertical Scaling**:
- Connection pooling (already implemented)
- Database query optimization with indexes
- Redis caching to reduce database load

**Microservices Approach**:
- Separate services for: User Management, Expense Management, Reporting, Currency Conversion
- Event-driven architecture with message queue (Kafka/RabbitMQ)
- API Gateway for routing

**Caching Strategy**:
- Multi-layer caching (L1: in-memory, L2: Redis, L3: Database)
- Cache warming for frequently accessed data
- Cache invalidation strategies

**Background Jobs**:
- Async processing for currency conversion
- Background report generation
- Email notifications

### 3. How would you ensure data consistency across services?

**ACID Transactions**:
- Use database transactions for critical operations
- Implement compensating transactions for distributed scenarios

**Event Sourcing**:
- Store all changes as events
- Rebuild state from events
- Enable audit trail and replay

**Saga Pattern**:
- For distributed transactions across services
- Compensating actions for rollback

**Idempotency**:
- Idempotency keys for all mutations
- Idempotent API endpoints

**Eventual Consistency**:
- Accept eventual consistency for non-critical paths
- Use eventual consistency with conflict resolution for reporting

### 4. What monitoring and alerting would you implement?

**Metrics**:
- Prometheus metrics for:
  - Request rate (requests/sec)
  - Error rate (4xx, 5xx errors)
  - Response time (p50, p95, p99)
  - Database connection pool usage
  - Redis cache hit/miss ratio
  - Currency API call success rate

**Logging**:
- Structured logging with correlation IDs
- Log levels: DEBUG, INFO, WARN, ERROR
- Centralized logging (ELK Stack or similar)

**Health Checks**:
- `/health` endpoint (already implemented)
- Database connectivity check
- Redis connectivity check
- External API (currency) health check

**Alerting**:
- High error rate (> 1% for 5 minutes)
- Slow response times (p95 > 1s)
- Database connection pool exhaustion
- Redis unavailability
- Currency API failures

**Tracing**:
- Distributed tracing (OpenTelemetry/Jaeger)
- Track requests across services
- Performance bottleneck identification

## 📝 Additional Notes

### Currency API

The API uses ExchangeRate-API (free tier) by default. To use a different provider:
1. Update `CURRENCY_API_URL` in `.env`
2. Modify `fetchExchangeRate` in `internal/services/currency_service.go`

### Authentication

Currently, user ID is passed via `X-User-ID` header for testing. In production, implement:
- JWT-based authentication
- OAuth 2.0
- API key authentication

### Rate Limiting

Rate limiting middleware can be added using:
```go
// Using gin-rate-limit middleware
limiter := gin_rate.New(gin_rate.Store(redis.NewRateStore()))
router.Use(limiter.Limit(100, time.Minute))
```

## 📄 License

This project is part of a technical assessment.

## 👤 Author

Backend Engineer Assessment - FlyPro

---

**Note**: This is a demonstration project for assessment purposes. Focus on code quality, architecture decisions, and engineering best practices rather than feature completeness.