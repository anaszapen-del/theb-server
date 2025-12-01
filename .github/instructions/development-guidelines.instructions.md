---
applyTo: "**"
---

# Development Guidelines

## Backend Development (Golang)

### Environment Setup
- Go 1.21+ with module support
- PostgreSQL 13+ for development database
- Redis 6+ for caching and real-time data
- Docker for isolated development environments
- Air for hot-reload during development

### Environment Setup
- Go 1.21+ with module support
- PostgreSQL 13+ for development database
- Redis 6+ for caching and real-time data
- Docker for isolated development environments
- Air for hot-reload during development

### Configuration Management
- Environment variables for all configuration
- Separate configs for development/staging/production
- Sensitive data never committed to version control
- Configuration validation on startup
- Use `.env` files with `godotenv` or Viper
- Database: PostgreSQL (database: `theb_db`, password: `00962`)
- Cache: Redis for real-time data and session management

### Microservices Architecture
- **Auth Service**: Phone OTP authentication, JWT token generation
- **User Service**: User profiles, passenger/captain management
- **Location Service**: Real-time captain location tracking (Redis + PostgreSQL)
- **Order Service**: Ride requests, matching, status lifecycle
- **Payment Service**: Fare calculation, payment processing, refunds
- **Rating Service**: Post-ride ratings and reviews
- **Notification Service**: Push notifications via Expo Push

### Real-Time Features
- WebSocket server for location streaming and ride updates
- Redis Pub/Sub for event broadcasting across services
- Geospatial queries with Redis GEOADD/GEORADIUS for captain matching
- Connection pooling and load balancing for WebSocket servers
- Rate limiting to prevent abuse

### Database Design
- PostgreSQL for persistent data (users, orders, payments, ratings)
- GORM as ORM with migrations
- Indexing on frequently queried fields (phone, captain location, order status)
- Foreign key constraints for referential integrity
- Soft deletes for audit trail

### API Design
- RESTful endpoints for CRUD operations
- WebSocket endpoints for real-time features
- API versioning (e.g., `/api/v1/`)
- JSON request/response format
- Consistent error response structure
- OpenAPI/Swagger documentation with Swag

### Security
- JWT tokens for authentication (short-lived access + refresh tokens)
- Rate limiting on auth endpoints to prevent brute force
- Input validation and sanitization
- HTTPS only in production
- CORS configuration for mobile app origin
- Secure storage of OTP codes (hashed, time-limited)

### Testing
- Unit tests for business logic (services, repositories)
- Integration tests with test database
- API endpoint tests with httptest
- Mock external dependencies (payment gateways, SMS providers)
- Test coverage > 70% target

## Documentation

- API documentation with OpenAPI/Swagger (use Swag for Go)
- Code documentation with godoc for backend
- Inline comments for complex logic
- README files for all modules and services
- Architecture diagrams for system overview
- After type struct always add the @name tag if used in API documentation
- Add @ID tag to handler functions if used in API documentation

## Version Control

- Git flow with feature branches
- Descriptive commit messages (conventional commits recommended)
- Pull request reviews required for main/develop branches
- Keep commits atomic and focused
- Tag releases with semantic versioning (v1.0.0)

## CI/CD

- Automated linting and tests on PR
- Docker images for backend services
- Staged deployments (dev → staging → production)
- Environment-specific configuration injection
- Rollback strategy for failed deployments
