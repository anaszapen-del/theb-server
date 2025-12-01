# THEB Backend Setup Script for Windows PowerShell

Write-Host "🐺 THEB Backend Setup Script" -ForegroundColor Cyan
Write-Host "==============================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "✅ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed. Please install Go 1.21 or higher." -ForegroundColor Red
    exit 1
}

Write-Host ""

# Check if Docker is installed (optional)
try {
    $dockerVersion = docker --version
    Write-Host "✅ Docker is installed: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Docker is not installed (optional for local database setup)" -ForegroundColor Yellow
}

Write-Host ""

# Create .env file if it doesn't exist
if (-not (Test-Path .env)) {
    Write-Host "📝 Creating .env file from .env.example..." -ForegroundColor Yellow
    Copy-Item .env.example .env
    Write-Host "✅ .env file created. Please update it with your configuration." -ForegroundColor Green
} else {
    Write-Host "✅ .env file already exists" -ForegroundColor Green
}

Write-Host ""

# Install Go dependencies
Write-Host "📦 Installing Go dependencies..." -ForegroundColor Yellow
go mod download
go mod tidy
Write-Host "✅ Dependencies installed" -ForegroundColor Green

Write-Host ""

# Check if PostgreSQL is accessible
Write-Host "🔍 Checking PostgreSQL connection..." -ForegroundColor Yellow
try {
    $env:PGPASSWORD = "00962"
    $null = psql -h localhost -U postgres -d theb_db -c "\q" 2>$null
    Write-Host "✅ PostgreSQL is running and theb_db exists" -ForegroundColor Green
} catch {
    Write-Host "⚠️  PostgreSQL is not accessible or theb_db doesn't exist" -ForegroundColor Yellow
    Write-Host "   Run: docker-compose up -d postgres" -ForegroundColor Gray
    Write-Host "   Or create the database manually: CREATE DATABASE theb_db;" -ForegroundColor Gray
}

Write-Host ""

# Check if Redis is running
Write-Host "🔍 Checking Redis connection..." -ForegroundColor Yellow
try {
    $redisTest = redis-cli ping 2>$null
    if ($redisTest -eq "PONG") {
        Write-Host "✅ Redis is running" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Redis is not running" -ForegroundColor Yellow
        Write-Host "   Run: docker-compose up -d redis" -ForegroundColor Gray
    }
} catch {
    Write-Host "⚠️  redis-cli not found in PATH" -ForegroundColor Yellow
    Write-Host "   Run: docker-compose up -d redis" -ForegroundColor Gray
}

Write-Host ""
Write-Host "==============================" -ForegroundColor Cyan
Write-Host "🎉 Setup Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Update .env file with your configuration"
Write-Host "2. Start PostgreSQL and Redis:"
Write-Host "   docker-compose up -d postgres redis"
Write-Host "3. Run the application:"
Write-Host "   go run cmd/app/main.go"
Write-Host ""
Write-Host "For API documentation:" -ForegroundColor Cyan
Write-Host "   Visit http://localhost:8080/swagger/index.html"
Write-Host ""
