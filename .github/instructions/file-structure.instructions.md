---
applyTo: "**/*"
---

# File Structure Guidelines

## Root Directory Structure

```
theb-backend/
├── cmd/
│   └── app/
│       └── main.go              # Application entry point
├── config/                      # Configuration files
│   ├── development.yaml
│   ├── staging.yaml
│   └── production.yaml
├── docs/                        # Documentation
│   ├── database_structure.md
│   ├── file-structure.md
│   └── theb_prd.md
├── internal/                    # Private application code
├── pkg/                         # Public shared packages
│   ├── errors/                  # Error types
│   ├── response/                # Response utilities
│   └── pagination/              # Pagination helpers
├── scripts/                     # Deployment scripts
├── tests/                       # E2E and integration tests
└── .github/                     # GitHub-specific configuration
    └── instructions/            # Copilot instructions
```

## Internal Directory Structure

```
internal/
├── app/                         # Application initialization
├── cache/                       # Redis cache utilities
├── config/                      # Internal config handling
├── container/                   # Dependency injection container
├── db/                          # Database connections
│   ├── postgres.go              # PostgreSQL connection (theb_db, password: 00962)
│   └── redis.go                 # Redis connection
├── logger/                      # Logging utilities
├── middleware/                  # HTTP middleware
│   ├── auth.go                  # JWT authentication
│   ├── cors.go                  # CORS configuration
│   ├── logger.go                # Request logging
│   └── rate_limit.go            # Rate limiting
├── realtime/                    # WebSocket handling
│   ├── handlers/                # WebSocket connection handlers
│   ├── hub.go                   # Connection pool manager
│   └── client.go                # WebSocket client handler
├── router/                      # Route definitions
└── service/                     # Business services
    ├── auth/                    # Auth service
    │   ├── handlers/            # HTTP handlers (OTP login, verify, refresh)
    │   ├── services/            # Business logic (token generation, validation)
    │   ├── repositories/        # Data access (user lookup, token storage)
    │   ├── models/              # GORM models (User, Token)
    │   ├── utils/               # Auth utilities (JWT helpers)
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Request/Response DTOs
    │   ├── docs/                # Module documentation
    │   └── entry.go             # Service registration
    ├── user/                    # User management
    │   ├── handlers/            # User CRUD handlers
    │   ├── services/            # User business logic
    │   ├── repositories/        # User data access
    │   ├── models/              # User, Passenger, Captain, Vehicle models
    │   ├── utils/               # User utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    ├── location/                # Location tracking
    │   ├── handlers/            # Location streaming, history
    │   ├── services/            # Geospatial queries, captain matching
    │   ├── repositories/        # Redis + PostgreSQL storage
    │   ├── models/              # CaptainLocation model
    │   ├── utils/               # Location utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    ├── order/                   # Order/Ride management
    │   ├── handlers/            # Ride request, accept, status update
    │   ├── services/            # Matching logic, lifecycle management
    │   ├── repositories/        # Order data access
    │   ├── models/              # Order, OrderEvent, Cancellation models
    │   ├── utils/               # Order utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    ├── payment/                 # Payment processing
    │   ├── handlers/            # Payment processing, fare calculation
    │   ├── services/            # Gateway integration, refunds
    │   ├── repositories/        # Payment data access
    │   ├── models/              # Payment model
    │   ├── utils/               # Payment utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    ├── rating/                  # Rating service
    │   ├── handlers/            # Rating submission, retrieval
    │   ├── services/            # Rating business logic
    │   ├── repositories/        # Rating data access
    │   ├── models/              # Rating model
    │   ├── utils/               # Rating utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    ├── notification/            # Notifications
    │   ├── handlers/            # Push notification handlers
    │   ├── services/            # Expo Push integration
    │   ├── repositories/        # Notification log
    │   ├── models/              # Notification model
    │   ├── utils/               # Notification utilities
    │   ├── tests/               # Unit tests
    │   ├── dtos/                # Data transfer objects
    │   ├── docs/                # Documentation
    │   └── entry.go             # Service registration
    └── captain/                 # Captain-specific features
        ├── handlers/            # Captain management handlers
        ├── services/            # Captain business logic
        ├── repositories/        # Captain data access
        ├── models/              # Captain, Vehicle models
        ├── utils/               # Captain utilities
        ├── tests/               # Unit tests
        ├── dtos/                # Data transfer objects
        ├── docs/                # Documentation
        └── entry.go             # Service registration
```

## Service Module Template

Every service module should follow this structure:

```
internal/service/[module-name]/
├── handlers/          # HTTP request handlers
├── services/          # Business logic layer
├── repositories/      # Data access layer
├── models/            # GORM models
├── utils/             # Module-specific utilities
├── tests/             # Unit and integration tests
├── dtos/              # Request/Response DTOs
├── docs/              # Module documentation
└── entry.go           # Entry point for service registration
```

### Service Module Guidelines

- **handlers/**: Thin HTTP handlers that parse requests and call services
- **services/**: Business logic implementation, orchestrates repositories
- **repositories/**: Data access layer using GORM for PostgreSQL operations
- **models/**: GORM models with tags, relationships, and database constraints
- **utils/**: Helper functions specific to this module
- **tests/**: Unit tests for services and repositories
- **dtos/**: Request and response data transfer objects with validation tags
- **docs/**: API documentation, usage examples, and module-specific notes
- **entry.go**: Registers handlers, services, and repositories in the DI container

## Naming Conventions

- **Directories**: lowercase, singular (e.g., `user`, not `users`)
- **Files**: snake_case for utilities, PascalCase for types
- **Go files**: descriptive names with clear purpose
- **Test files**: `*_test.go` pattern

## Import Organization

```go
import (
    // Standard library
    "context"
    "fmt"

    // Third-party packages
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // Internal packages
    "theb-backend/internal/config"
    "theb-backend/internal/service/user/models"
    "theb-backend/pkg/response"
)
```

## Database Configuration

- **Database**: PostgreSQL 13+
- **Database Name**: `theb_db`
- **Password**: `00962`
- **ORM**: GORM with auto-migrations or manual migration files
- **Cache**: Redis 6+ for real-time data and session management

## Configuration Files

- **Backend**: `config/` directory for environment-specific configs (YAML/JSON)
- Environment variables for secrets (API keys, database credentials)
- Validation of required configuration on startup
