# THEB Backend

THEB (ذيب) - Ride-Hailing Backend API for Mafraq, Jordan

## Overview

This is the backend server for the THEB ride-hailing application, built with Golang using clean architecture principles. The system supports both Passenger and Captain roles with real-time location tracking, ride matching, and WebSocket communication.

## Features

- 🔐 Phone-based OTP authentication with JWT tokens
- 👤 User management (Passenger & Captain roles)
- 📍 Real-time location tracking with Redis & WebSocket
- 🚗 Ride request and captain matching
- 💳 Payment processing and fare calculation
- ⭐ Post-ride rating system
- 🔔 Push notifications via Expo
- 📊 Real-time ride status updates

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL (theb_db, password: 00962)
- **Cache**: Redis 6+
- **ORM**: GORM
- **WebSocket**: Gorilla WebSocket
- **Auth**: JWT
- **Documentation**: Swag (OpenAPI/Swagger)

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+
- Redis 6+
- Make (optional, for Makefile commands)
- Docker & Docker Compose (optional, for containerized setup)

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd server
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 4. Start PostgreSQL and Redis

Using Docker Compose:
```bash
docker-compose up -d postgres redis
```

Or install them locally and ensure they're running.

### 5. Run database migrations

```bash
# Migrations will be auto-applied on first run
# Or manually run migration files from migrations/
```

### 6. Run the application

Development mode with hot-reload:
```bash
air
```

Or standard go run:
```bash
go run cmd/app/main.go
```

The server will start at `http://localhost:8080`

## API Documentation

Once the server is running, access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

Generate/update Swagger docs:
```bash
make swagger
```

## Project Structure

```
theb-backend/
├── cmd/app/                 # Application entry point
├── config/                  # Configuration files
├── internal/                # Private application code
│   ├── app/                 # Application initialization
│   ├── service/             # Business services
│   │   ├── auth/            # Authentication service
│   │   ├── user/            # User management
│   │   ├── location/        # Location tracking
│   │   ├── order/           # Ride/Order management
│   │   ├── payment/         # Payment processing
│   │   ├── rating/          # Rating system
│   │   ├── notification/    # Notifications
│   │   └── captain/         # Captain features
│   ├── middleware/          # HTTP middleware
│   ├── router/              # Route definitions
│   ├── db/                  # Database connections
│   └── ...
├── pkg/                     # Public packages
└── docs/                    # Documentation
```

## Available Commands

Using Makefile:
```bash
make help              # Show all available commands
make install           # Install dependencies
make run               # Run the application
make build             # Build binary
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make format            # Format code
make swagger           # Generate API documentation
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
```

## Configuration

Configuration files are located in `config/` directory:
- `development.yaml` - Development environment
- `production.yaml` - Production environment

Environment can be set via `APP_ENV` environment variable.

## Database

- **Name**: `theb_db`
- **Password**: `00962`
- **Port**: `5432`

Key tables:
- users, captains, orders, payments, ratings, locations_history, notifications

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
make test-coverage
```

## Deployment

### Docker

Build and run with Docker Compose:
```bash
docker-compose up --build
```

### Manual Deployment

1. Build the binary:
```bash
make build
```

2. Set environment variables for production

3. Run the binary:
```bash
./bin/theb-backend
```

## Contributing

1. Follow the coding style guidelines in `.github/instructions/`
2. Write tests for new features
3. Update API documentation
4. Submit pull requests for review

## License

[Your License Here]

## Contact

- **Project**: THEB (ذيب)
- **Website**: [theb.app](http://theb.app)
- **Support**: support@theb.app
