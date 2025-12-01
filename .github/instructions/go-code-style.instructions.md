---
applyTo: "**/*"
---

# Go Code Style Guidelines

## Formatting and Standards

- Use `gofmt` and `golint` standards
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use structured logging with correlation IDs
- Implement proper error handling with context

## Architecture Patterns

- **Clean Architecture**: Separate concerns with dependency injection
- **Repository Pattern**: All database operations through repositories
- **Service Layer**: Business logic in dedicated service structs
- **Handler Layer**: HTTP endpoints as thin handlers

## Code Organization

- Each module must be independently testable
- Use dependency injection for testability
- All database operations through GORM
- Follow RESTful API design principles
- Implement graceful shutdown patterns
- One responsibility per file

## Error Handling

- Use context for request-scoped operations
- Return structured errors with proper HTTP status codes
- Log errors with correlation IDs for debugging
- Handle database connection errors gracefully

## THEB-Specific Backend Services

### Auth Service
- Phone-based OTP authentication (no email/password)
- JWT token generation (access + refresh tokens)
- Token validation and refresh endpoints
- Rate limiting on OTP requests to prevent abuse
- Store OTP codes securely (hashed, time-limited in Redis)

### User Service
- User profile management (passenger + captain roles)
- Captain vehicle information and verification
- Role-based data access (passenger | captain | both)
- License and identity verification workflow
- Support for multi-role users

### Location Service
- Real-time captain location tracking via WebSocket
- Redis GEOADD/GEORADIUS for geospatial queries
- Historical location logging to PostgreSQL
- Captain matching based on proximity
- Location privacy (only visible to matched passenger)
- Optimize for low-latency updates (< 100ms)

### Order Service
- Ride request creation and validation
- Nearest captain matching algorithm
- Order status lifecycle management (requested → accepted → arriving → in-progress → completed)
- Cancellation handling (passenger/captain/system)
- Order event history tracking
- Real-time status updates via WebSocket

### Payment Service
- Fare calculation based on distance and time
- Payment gateway integration (Stripe or local Jordanian provider)
- Transaction processing and recording
- Refund handling for cancelled rides
- Payment status tracking
- Integration with order completion

### Rating Service
- Post-ride captain rating by passenger
- Rating aggregation and captain score calculation
- Rating history and trends
- Moderation for inappropriate feedback

### Notification Service
- Expo Push Notifications integration
- Ride alerts (request, acceptance, arrival, completion)
- Captain notifications (new ride requests, cancellations)
- Push token management
- Delivery status tracking

### WebSocket Server
- Real-time bidirectional communication
- Connection pooling and management
- Event broadcasting (ride updates, location streaming)
- Redis Pub/Sub for multi-instance coordination
- Heartbeat/ping-pong for connection health
- Graceful reconnection handling

## Documentation

- API documentation with OpenAPI/Swagger
- Code documentation with godoc
- README files for all modules
- Use Swag for API documentation
- After type struct always add the @name tag to the struct if it is used in the API documentation
- Add @ID tag to the handler function if it is used in the API documentation

## THEB API Examples

### Auth Endpoints
```go
// @Summary Login with phone number
// @ID auth-login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Phone number"
// @Success 200 {object} OTPResponse
// @Router /api/v1/auth/login [post]

// @Summary Verify OTP
// @ID auth-verify
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Phone and OTP"
// @Success 200 {object} TokenResponse
// @Router /api/v1/auth/verify [post]
```

### Order Endpoints
```go
// @Summary Request a ride
// @ID order-request
// @Tags Order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body RideRequest true "Pickup and dropoff locations"
// @Success 200 {object} OrderResponse
// @Router /api/v1/orders/request [post]

// @Summary Accept ride (captain)
// @ID order-accept
// @Tags Order
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} OrderResponse
// @Router /api/v1/orders/{id}/accept [post]
```

### Location Endpoints
```go
// @Summary Stream captain location (WebSocket)
// @ID location-stream
// @Tags Location
// @Security BearerAuth
// @Router /ws/location/stream [get]

// @Summary Get captain location history
// @ID location-history
// @Tags Location
// @Security BearerAuth
// @Param captain_id query string true "Captain ID"
// @Success 200 {array} LocationPoint
// @Router /api/v1/location/history [get]
```
