# Project Copilot Instructions

## Overview

- This repo contains the **THEB** backend services built with Golang. Mobile application (React Native) is separate; this is the backend server.
- Goal: Deliver scalable microservices supporting both **Passenger** and **Captain** roles with real-time location tracking, ride matching, and WebSocket communication.
- Database: PostgreSQL (database name: `theb_db`, password: `00962`)
- Keep instructions evidence-based: only document what exists; place open items in Unresolved Questions.

## Architecture

- **Microservices architecture** using Golang with clean architecture principles.
- Entry point: Each service has its own `cmd/<service-name>/main.go` file.
- Service modules: Organized under `internal/` with handlers, services, repositories, models, and DTOs.
- API routing: Gin framework with middleware for auth, CORS, logging, and rate limiting.
- Database: PostgreSQL (`theb_db`) for persistent data with GORM as ORM.
- Cache: Redis for real-time captain locations and session management.
- Real-time: WebSocket server for location streaming and ride status updates.
- Configuration: YAML files for different environments (development, staging, production) with environment variable overrides.

## Tech Stack

- **Golang 1.21+** with module support and clean architecture patterns.
- **Gin Framework** for HTTP routing and middleware.
- **GORM** as ORM for PostgreSQL database operations.
- **PostgreSQL 13+** (database: `theb_db`, password: `00962`) for persistent data.
- **Redis 6+** for caching, real-time captain locations, and session management.
- **WebSocket (Gorilla)** for real-time bidirectional communication.
- **JWT** for authentication and authorization.
- **Swag** for OpenAPI/Swagger documentation generation.
- **Air** for hot-reload during development.
- **Docker** for containerization and deployment.

## Key Workflows

- **Dev server**: `air` for hot-reload or `go run cmd/<service-name>/main.go` for specific service.
- **Database setup**: PostgreSQL connection to `theb_db` with password `00962`.
- **Migrations**: GORM auto-migrate or manual migration files in `migrations/`.
- **API documentation**: `swag init` generates Swagger docs from code annotations.
- **Testing**: `go test ./...` runs all tests with coverage reporting.
- **Build**: `go build -o bin/<service-name> cmd/<service-name>/main.go` creates executable.
- **Docker**: `docker-compose up` starts all services (PostgreSQL, Redis, backend services).

## Conventions

- **Services**: Each microservice has `handlers/`, `services/`, `repositories/`, `models/`, and `dtos/` subdirectories.
- **Entry points**: Each service has `entry.go` file that registers all dependencies.
- **Database models**: GORM models with proper tags, relationships, and indexes.
- **DTOs**: Separate request and response DTOs with validation tags.
- **Error handling**: Use custom error types with proper HTTP status codes.
- **Middleware**: Auth, CORS, logging, and rate limiting applied globally or per route.
- **API versioning**: All endpoints prefixed with `/api/v1/`.
- **Logging**: Structured logging with correlation IDs for request tracing.
- **Testing**: Unit tests for services and repositories, integration tests for handlers.

## Domain Model (Current State)

- **User**: Phone-based authentication, role selection (passenger | captain | both).
- **Captain**: User with vehicle details, license verification, online status, location tracking.
- **Vehicle**: Belongs to captain; has type, model, year, plate number.
- **Order/Ride**: Links Passenger, Captain, pickup/dropoff locations, status (requested → accepted → arriving → in-progress → completed → cancelled).
- **Location**: Real-time captain location (Redis + WebSocket streaming), historical location log (PostgreSQL).
- **Payment**: Ride fare calculation and payment processing with transaction recording.
- **Rating**: Post-ride passenger rating of captain with review text.
- **Cancellation**: Records who cancelled (passenger/captain/system), reason, timestamp.
- **Order Event**: Status change history for audit trail.

## Integration Points

- **Mobile App**: React Native (Expo) client consumes REST APIs and WebSocket endpoints.
- **Database**: PostgreSQL (`theb_db`, password: `00962`) for persistent data storage.
- **Cache**: Redis for real-time captain locations, session management, and geospatial queries.
- **WebSocket**: Real-time ride updates, location streaming, captain matching notifications.
- **Google Maps**: Directions API for routing calculations (backend integration).
- **Phone Auth**: OTP-based login via SMS provider integration (backend).
- **Payment Gateway**: Future integration (Stripe/local Jordanian gateway) for ride payments.
- **Push Notifications**: Expo Push API for ride alerts and captain notifications (backend sends).

## Cross-Cutting Concerns

- **Auth**: JWT-based authentication with access and refresh tokens. Phone OTP verification flow.
- **Security**: HTTPS only in production, input validation, SQL injection prevention, rate limiting.
- **Logging**: Structured logging with correlation IDs for request tracing and debugging.
- **Error Handling**: Custom error types with proper HTTP status codes and user-friendly messages.
- **Configuration**: Environment-specific configs (development, staging, production) via YAML/env vars.
- **Monitoring**: Health check endpoints, metrics collection, performance monitoring.
- **Database**: Connection pooling, transaction management, migration strategy.
- **Caching**: Redis for session management, real-time data, and performance optimization.

## Testing Strategy

- **Unit Tests**: Test business logic in services and utilities using Go's testing package.
- **Repository Tests**: Test data access layer with test database or mocks.
- **Integration Tests**: Test API endpoints with HTTP test server.
- **WebSocket Tests**: Test real-time communication with WebSocket client mocks.
- **Load Tests**: Test system performance under high load (Apache Bench, k6).
- **Coverage**: Target >70% code coverage for critical paths.
- **Mocking**: Use interfaces for dependencies to enable mocking in tests.

## Code Generation & Tooling

- **Swag**: Generate OpenAPI/Swagger documentation from code annotations.
- **GORM**: Generate migrations from model definitions or use manual migration files.
- **Air**: Hot-reload for development (watches Go files and restarts server).
- **golangci-lint**: Comprehensive linting for Go code quality.
- **Docker Compose**: Orchestrate PostgreSQL, Redis, and backend services locally.
- Avoid speculative codegen until requirements are clear.

## Design System (THEB Brand)

### Colors
- **Wolf Black**: `#0D0D0D` – primary backgrounds, confidence.
- **Desert Gold**: `#D4A048` – accents, CTAs, Bedouin identity.
- **Pure White**: `#FFFFFF` – text, contrast.
- **Sand Gray**: `#B6B0A2` – soft UI elements.
- **Mafraq Blue**: `#3A6EA5` – subtle tech accent.
- **Success Green**: `#41A45A` – success states.
- **Alert Red**: `#D9534F` – errors, alerts.

### Typography
- **Primary Font**: Cairo (Arabic) / Inter (English).
- Style: Bold, clean, geometric.
- Used in API documentation and admin dashboard.

## API Architecture

### Microservices
- **Auth Service**: Phone OTP authentication, JWT token generation and validation.
- **User Service**: User profiles, passenger/captain management, vehicle information.
- **Location Service**: Real-time captain location tracking, geospatial queries, captain matching.
- **Order Service**: Ride requests, captain matching, status lifecycle management.
- **Payment Service**: Fare calculation, payment processing, transaction recording.
- **Rating Service**: Post-ride ratings and reviews.
- **Notification Service**: Push notifications via Expo Push API.
- **WebSocket Server**: Real-time bidirectional communication for location and ride updates.

### API Endpoints Structure
- `/api/v1/auth/*` - Authentication endpoints
- `/api/v1/users/*` - User management endpoints
- `/api/v1/captains/*` - Captain-specific endpoints
- `/api/v1/orders/*` - Ride/order management endpoints
- `/api/v1/payments/*` - Payment processing endpoints
- `/api/v1/ratings/*` - Rating and review endpoints
- `/api/v1/locations/*` - Location tracking endpoints
- `/ws/*` - WebSocket endpoints for real-time features

## Real-Time Features

- **Location Tracking**: Captain app streams location to backend via WebSocket when online.
- **Ride Matching**: Backend sends ride request to nearest captain; passenger gets captain details.
- **Live Updates**: Ride status changes (accepted, arriving, in-progress, completed) pushed via WebSocket.
- **Geospatial Queries**: Redis GEOADD/GEORADIUS for finding nearest captains within radius.
- **Connection Management**: WebSocket hub manages active connections, handles reconnections.
- **Event Broadcasting**: Redis Pub/Sub for coordinating events across multiple server instances.

## Examples

### Database Connection
```go
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// DSN: host=localhost user=postgres password=00962 dbname=theb_db port=5432
```

### JWT Authentication Middleware
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // Validate JWT token, set user context
    }
}
```

### WebSocket Connection Handler
```go
func HandleWebSocket(c *gin.Context) {
    conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
    client := &Client{conn: conn, hub: hub}
    hub.register <- client
}
```

### GORM Model Example
```go
type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Phone     string    `gorm:"unique;not null"`
    Name      string
    Role      string    `gorm:"type:varchar(20)"` // passenger, captain, both
    CreatedAt time.Time
}
```

## Do / Avoid

**Do:**
- Use clean architecture with clear separation of concerns (handlers, services, repositories).
- Implement proper error handling with custom error types and HTTP status codes.
- Use GORM for database operations with proper migrations.
- Apply middleware for authentication, logging, and rate limiting.
- Document APIs with Swag annotations for OpenAPI/Swagger generation.
- Write unit tests for services and repositories.
- Use Redis for real-time data and caching.
- Implement graceful shutdown for all services.

**Avoid:**
- Hardcoding configuration values; use environment variables or config files.
- Direct database access from handlers; always go through repositories.
- Ignoring errors; always handle and log errors properly.
- Blocking operations in WebSocket handlers; use goroutines for concurrent processing.
- Exposing sensitive data in API responses (passwords, tokens).
- Using SELECT * in queries; specify required fields.

## Maintenance Tips

- Update this file when adding: new microservices, external API integrations, environment variables, or database schema changes.
- Document environment variables (DB credentials, API keys, JWT secret) in config files with clear fallback behavior.
- Keep database migrations versioned and documented.
- Document API endpoints with Swag annotations for automatic documentation generation.

## Unresolved Questions

- API rate limiting strategy (per user, per IP, per endpoint)?
- Distributed tracing implementation (Jaeger, OpenTelemetry)?
- Service mesh for microservices communication (Istio, Linkerd)?
- Message queue for asynchronous tasks (RabbitMQ, Kafka)?
- Monitoring and alerting tools (Prometheus, Grafana)?
- CI/CD pipeline configuration (GitHub Actions, GitLab CI)?
- Load balancing strategy for WebSocket connections?
- Database backup and disaster recovery plan?
- Multi-region deployment strategy?
- API gateway for microservices (Kong, Traefik)?

---

Update responsibly: remove Unresolved items as decisions are implemented in code.
