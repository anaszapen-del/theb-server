#!/bin/bash

echo "🐺 THEB Backend Setup Script"
echo "=============================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Check if Docker is installed (optional)
if command -v docker &> /dev/null; then
    echo "✅ Docker is installed: $(docker --version)"
else
    echo "⚠️  Docker is not installed (optional for local database setup)"
fi

echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
else
    echo "✅ .env file already exists"
fi

echo ""

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies installed"

echo ""

# Check if PostgreSQL is running
echo "🔍 Checking PostgreSQL connection..."
if command -v psql &> /dev/null; then
    if PGPASSWORD=00962 psql -h localhost -U postgres -d theb_db -c '\q' 2>/dev/null; then
        echo "✅ PostgreSQL is running and theb_db exists"
    else
        echo "⚠️  PostgreSQL is not accessible or theb_db doesn't exist"
        echo "   Run: docker-compose up -d postgres"
        echo "   Or create the database manually: CREATE DATABASE theb_db;"
    fi
else
    echo "⚠️  psql not found in PATH. Make sure PostgreSQL is installed."
fi

echo ""

# Check if Redis is running
echo "🔍 Checking Redis connection..."
if command -v redis-cli &> /dev/null; then
    if redis-cli ping &> /dev/null; then
        echo "✅ Redis is running"
    else
        echo "⚠️  Redis is not running"
        echo "   Run: docker-compose up -d redis"
    fi
else
    echo "⚠️  redis-cli not found in PATH"
fi

echo ""
echo "=============================="
echo "🎉 Setup Complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your configuration"
echo "2. Start PostgreSQL and Redis:"
echo "   docker-compose up -d postgres redis"
echo "3. Run the application:"
echo "   make run  (or)  go run cmd/app/main.go"
echo ""
echo "For hot-reload development:"
echo "   air"
echo ""
echo "For API documentation:"
echo "   Visit http://localhost:8080/swagger/index.html"
echo ""
