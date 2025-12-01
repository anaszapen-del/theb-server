# Auth Service Documentation

## Overview

The Auth service handles passenger and captain authentication using phone-based OTP verification. This service provides secure signup, login, and token management functionality.

## Features

- Phone number-based signup with OTP verification
- JWT token generation (access + refresh tokens)
- Secure OTP storage with bcrypt hashing in Redis
- Rate limiting to prevent abuse
- E.164 phone number validation
- Multi-language name support (Arabic, English)

## API Endpoints

### Passenger Signup

**Endpoint:** `POST /api/v1/auth/passenger/signup`

**Description:** Initiates passenger signup by sending an OTP to the provided phone number.

**Request Body:**
```json
{
  "name": "Ahmed Ali",
  "phone_number": "+962791234567"
}
```

**Success Response (201):**
```json
{
  "message": "OTP sent to phone number",
  "phone_number": "+962791234567",
  "expires_in": 300
}
```

**Error Responses:**
- `400` - Invalid phone number format or name
- `409` - Phone number already registered
- `429` - Rate limit exceeded (max 5 per hour)
- `503` - SMS delivery failed

### OTP Verification

**Endpoint:** `POST /api/v1/auth/passenger/verify`

**Description:** Verifies OTP and completes signup, returning authentication tokens.

**Request Body:**
```json
{
  "phone_number": "+962791234567",
  "otp_code": "123456",
  "name": "Ahmed Ali"
}
```

**Success Response (200):**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "phone_number": "+962791234567",
  "name": "Ahmed Ali",
  "role": "passenger",
  "expires_in": 900
}
```

**Error Responses:**
- `400` - Invalid OTP format
- `401` - Invalid or expired OTP
- `409` - User already exists
- `429` - Too many verification attempts (max 10 per hour)

## Configuration

### OTP Settings

```yaml
otp:
  expiry: 5m      # OTP validity period
  length: 6       # OTP code length (digits)
```

### JWT Settings

```yaml
jwt:
  secret: "your-secret-key"
  access_token_expiry: 15m    # 15 minutes
  refresh_token_expiry: 168h  # 7 days
```

### Rate Limiting

```yaml
rate_limit:
  per_minute: 100      # General rate limit
  otp_per_hour: 5      # OTP requests per phone per hour
```

## Security Features

### OTP Security

- OTP codes are hashed using bcrypt before storage
- Stored in Redis with automatic expiry (5 minutes)
- Rate limiting prevents brute force attacks
- OTP is deleted immediately after successful verification

### JWT Tokens

- Access tokens: Short-lived (15 minutes)
- Refresh tokens: Long-lived (7 days)
- Tokens include user ID, phone number, and role
- Signed with HMAC-SHA256

### Phone Number Validation

- E.164 format required: `+[country code][number]`
- Regex validation: `^\+[1-9]\d{1,14}$`
- Auto-normalization (removes spaces, dashes, parentheses)

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  phone_number VARCHAR(20) UNIQUE NOT NULL,
  role VARCHAR(20) NOT NULL DEFAULT 'passenger',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_users_phone_number ON users(phone_number);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

## Redis Keys

### OTP Storage

```
Key: otp:{phone_number}
Value: bcrypt_hashed_otp
Expiry: 5 minutes
```

### Rate Limiting

```
Key: signup_attempts:{phone_number}
Value: attempt_count
Expiry: 1 hour

Key: otp_attempts:{phone_number}
Value: verification_attempt_count
Expiry: 1 hour
```

## Testing

### Unit Tests

Run validation tests:
```bash
go test -v ./internal/service/auth/tests/... -run TestValidate
```

### Integration Tests

Run all auth tests:
```bash
go test -v ./internal/service/auth/tests/...
```

## Development Mode

### OTP Logging

In development mode, OTP codes are logged to the console instead of being sent via SMS:

```json
{
  "level": "INFO",
  "message": "SMS OTP (DEVELOPMENT MODE)",
  "phone": "+962791234567",
  "otp": "123456"
}
```

### Testing without Redis

The service can run without Redis in development mode, but OTP functionality will not work. Start Redis for full functionality:

```bash
redis-server
```

## Production Deployment

### SMS Provider Integration

Implement SMS sending in `auth_service.go`:

```go
func (s *AuthService) sendOTP(ctx context.Context, phoneNumber, otpCode string) error {
    // TODO: Integrate with SMS provider
    // Examples: Twilio, AWS SNS, or local Jordanian provider
    
    return smsProvider.Send(phoneNumber, fmt.Sprintf(
        "Your THEB verification code is: %s", otpCode,
    ))
}
```

### Environment Variables

```bash
export APP_ENV=production
export JWT_SECRET=your-production-secret
export REDIS_HOST=your-redis-host
export REDIS_PASSWORD=your-redis-password
export SMS_PROVIDER_API_KEY=your-sms-api-key
```

### Security Checklist

- [ ] Change JWT secret from default
- [ ] Enable HTTPS/TLS in production
- [ ] Configure proper CORS origins
- [ ] Set up SMS provider integration
- [ ] Enable Redis authentication
- [ ] Configure proper rate limits
- [ ] Set up monitoring and alerting
- [ ] Enable request logging

## Error Handling

All errors return a consistent format:

```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": {
    "additional": "context"
  }
}
```

### Error Codes

- `BAD_REQUEST` - Invalid input data
- `UNAUTHORIZED` - Invalid or expired OTP
- `CONFLICT` - Phone number already exists
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `INTERNAL_SERVER_ERROR` - Server error
- `SERVICE_UNAVAILABLE` - SMS service unavailable

## Architecture

```
handlers/
  - auth_handler.go       # HTTP request handlers

services/
  - auth_service.go       # Business logic

repositories/
  - user_repository.go    # Database operations

models/
  - user.go               # User data model

dtos/
  - signup.go             # Request/Response DTOs

utils/
  - validation.go         # Phone/name validation
  - jwt.go                # JWT token utilities
  - otp.go                # OTP management

tests/
  - validation_test.go    # Validation tests
  - integration_test.go   # API integration tests
```

## Future Enhancements

- [ ] Email verification (optional)
- [ ] Social login (Google, Facebook, Apple)
- [ ] Two-factor authentication (2FA)
- [ ] Biometric authentication
- [ ] Password-based login (optional)
- [ ] Session management
- [ ] Device tracking
- [ ] Suspicious activity detection
- [ ] Account recovery flow
