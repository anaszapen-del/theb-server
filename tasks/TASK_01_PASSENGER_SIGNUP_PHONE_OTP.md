# PASSENGER-SIGNUP – Passenger Phone-Number-Based Sign Up

## Overview
Implement a phone-number-based signup flow for new passengers. Users enter their name and phone number, receive an OTP for verification, and are automatically logged in upon successful verification. This is the entry point for new passengers to create accounts in THEB.

## Source
- Type: Story
- Priority: High (MVP feature)
- Labels: Auth, Passenger, MVP
- Reporter: Product Team
- Linked Dependencies: None

## Scope

### In Scope
- Passenger signup endpoint (`POST /api/v1/auth/passenger/signup`) accepting name and phone number
- OTP generation and delivery (via existing SMS provider)
- OTP verification endpoint (`POST /api/v1/auth/passenger/verify`)
- Auto-creation of passenger profile in PostgreSQL upon successful verification
- JWT token generation and session management after signup completion
- Phone number validation (E.164 format)
- Duplicate phone number detection and rejection
- Secure storage of OTP codes (hashed, time-limited in Redis)

### Out of Scope
- SMS provider integration (assumed already available)
- Mobile app UI (handled by React Native frontend)
- Admin dashboard for passenger approval
- Wallet/payment system initialization
- Email verification (phone OTP only)

## Requirements

### Functional
- Accept name and phone number in signup request
- Validate phone number format (E.164)
- Check for duplicate phone numbers in database
- Generate and send OTP to phone number
- Verify OTP against stored code
- Create passenger user record with fields: name, phone_number, created_at, updated_at
- Generate access and refresh JWT tokens upon successful verification
- Automatically log in user after verification (return tokens)
- Ensure OTP codes are invalidated after use

### Non-Functional
- OTP verification must complete within <2 seconds
- OTP validity window: 10 minutes
- Rate limiting on signup endpoint (max 5 requests per phone per hour)
- Rate limiting on verify endpoint (max 10 attempts per OTP per hour)
- Secure storage: OTP codes hashed in Redis, never logged
- Support concurrent signup requests for different phone numbers
- Graceful handling of SMS delivery failures

## Data & Contracts

### New API Endpoints

#### POST /api/v1/auth/passenger/signup
**Request:**
```json
{
  "name": "Ahmed Ali",
  "phone_number": "+962791234567"
}
```

**Response (201 Created):**
```json
{
  "message": "OTP sent to phone number",
  "phone_number": "+962791234567",
  "expires_in": 600
}
```

**Errors:**
- 400: Invalid phone number format
- 400: Invalid name (empty or too short)
- 409: Phone number already exists
- 429: Rate limit exceeded
- 503: SMS delivery failed

#### POST /api/v1/auth/passenger/verify
**Request:**
```json
{
  "phone_number": "+962791234567",
  "otp_code": "123456"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "user_id": "uuid-v4",
  "phone_number": "+962791234567",
  "name": "Ahmed Ali",
  "role": "passenger",
  "expires_in": 3600
}
```

**Errors:**
- 400: Invalid OTP format
- 400: Phone number not found
- 401: OTP is invalid or expired
- 429: Too many verify attempts
- 500: Token generation failed

### Database Schema
**users table** (new fields for signup):
- `id` (UUID, primary key)
- `name` (string, not null, 2-100 chars)
- `phone_number` (string, unique, not null, E.164 format)
- `role` (enum: passenger | captain | both, default: passenger)
- `created_at` (timestamp)
- `updated_at` (timestamp)

**Existing or New** (assume auth service has):
- OTP storage in Redis: `otp:{phone_number}` with hashed code + expiry

## UI / Components
- **Signup Screen**: Input fields for name and phone number, "Send OTP" button
- **OTP Verification Screen**: OTP input field, countdown timer (10 min), "Verify" button, "Resend OTP" option
- Note: Mobile app handles UI; backend provides endpoints only

## State & Logic

### Signup Logic (Handler → Service → Repository)
1. **Handler** validates request structure
2. **Service**:
   - Validates phone format (E.164 regex)
   - Validates name (non-empty, 2-100 chars)
   - Checks if phone already exists (DB query)
   - Generates OTP code (6 digits)
   - Hashes and stores OTP in Redis with 10-min expiry
   - Sends OTP via SMS provider
   - Returns response with phone and OTP expiry
3. **Repository**: Handles phone existence check in `users` table

### Verification Logic (Handler → Service → Repository)
1. **Handler** validates request structure
2. **Service**:
   - Retrieves hashed OTP from Redis
   - Compares provided OTP with stored (bcrypt.Compare)
   - If invalid/expired, increment attempt counter (Redis)
   - If valid:
     - Delete OTP from Redis
     - Create user record in DB with name, phone, role='passenger'
     - Generate JWT access token (15 min expiry)
     - Generate JWT refresh token (7 days expiry)
     - Return tokens + user data
3. **Repository**: Handles user creation in `users` table

### State Storage
- **OTP codes**: Redis (temporary, hashed)
- **User data**: PostgreSQL
- **Tokens**: Issued to client (stored in JWT claims)
- **Rate limiting**: Redis counters (signup_{phone}, verify_{phone})

## Acceptance Criteria

- [x] Backend accepts POST `/api/v1/auth/passenger/signup` with name and phone number
- [x] System validates phone number format (E.164)
- [x] Duplicate phone numbers are rejected with 409 Conflict error
- [x] OTP code is generated and sent to phone number
- [x] OTP code stored securely (hashed in Redis)
- [x] Backend accepts POST `/api/v1/auth/passenger/verify` with phone and OTP
- [x] Successful OTP verification creates passenger user record in DB
- [x] User record includes: name, phone_number, created_at, updated_at, role='passenger'
- [x] JWT tokens (access + refresh) are generated after verification
- [x] User is automatically logged in (tokens returned)
- [x] OTP is invalidated after use
- [x] Rate limiting enforced on both endpoints
- [x] Error handling for SMS delivery failures
- [x] All endpoints return consistent error response format

## Edge Cases

- **Empty/null fields**: Reject with 400 Bad Request
- **Invalid phone format**: Reject with descriptive error message
- **Duplicate phone**: Show conflict error, suggest login instead
- **OTP expired**: User must request new OTP
- **OTP invalid**: Show error, allow retry (up to rate limit)
- **SMS fails to send**: Return 503, user can retry
- **Database error on user creation**: Rollback OTP deletion, return 500
- **Concurrent signup same phone**: First request wins, second gets duplicate error
- **Very long name (>100 chars)**: Truncate or reject, define behavior
- **Special characters in name**: Allow (support Arabic names like "أحمد")
- **Phone number with spaces/dashes**: Normalize to E.164 format or reject

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| SMS delivery delays | User can't verify quickly, poor UX | Set realistic OTP timeout (10 min), show timeout to user, allow resend |
| OTP code guessing attacks | Account takeover | Hash OTP in Redis, rate limit verify attempts (max 10/hour), log failed attempts |
| Database connection failures | User creation fails, but tokens issued | Wrap user creation in transaction, validate DB before issuing tokens |
| Phone number collision (non-unique) | Data corruption | Add unique constraint in DB, handle duplicate error gracefully |
| Large volume of OTP requests | DOS attack, SMS cost | Rate limit signup endpoint (5 req/phone/hour), validate phone format early |
| Token leakage | Session hijacking | Use secure token storage (HttpOnly cookies or secure local storage), set short expiry |
| Concurrent verify requests with same OTP | Invalid state | Lock OTP or mark used immediately after first verify |

## Open Questions / Assumptions

- **Q**: Which SMS provider is used? **A**: Assumed already integrated; backend calls existing SMS service
- **Q**: Should signup require email too? **A**: No, phone-only per requirements
- **Q**: Can users change phone number later? **A**: Out of scope; handle in profile update story
- **Q**: What is the exact OTP length (digits)? **A**: Assumed 6 digits (standard)
- **Q**: Should we validate that name contains only letters? **A**: No; allow Arabic, numbers, spaces to support local names
- **Q**: Is there a captcha on signup? **A**: No, not mentioned; rate limiting sufficient for MVP
- **Q**: Can passenger sign up with both roles at once? **A**: No, signup creates passenger role; captain conversion is separate

## Story Points

**5 points** – Reason: Medium complexity feature involving new API endpoints, OTP generation/verification logic, user creation, JWT token handling, rate limiting, and database schema. Requires coordination between handlers, services, repositories, and Redis. Good test coverage needed.

## Implementation Plan

1. **Database Setup**
   - Update `users` model (add/verify fields: name, phone_number, role, created_at, updated_at)
   - Create migration file for schema changes
   - Add unique constraint on phone_number
   - Run migration

2. **Auth Service DTOs & Models**
   - Create `SignupRequest` DTO (name, phone_number)
   - Create `SignupResponse` DTO (message, phone_number, expires_in)
   - Create `VerifyRequest` DTO (phone_number, otp_code)
   - Create `VerifyResponse` DTO (access_token, refresh_token, user_id, phone_number, name, role, expires_in)
   - Update `User` model with JSON tags

3. **Auth Service – Signup Handler**
   - Create `POST /api/v1/auth/passenger/signup` endpoint
   - Parse and validate request (name, phone)
   - Call signup service
   - Handle and return responses with proper status codes

4. **Auth Service – Signup Service Logic**
   - Implement phone format validation (E.164 regex)
   - Implement name validation (non-empty, length checks)
   - Call repository to check phone exists
   - Generate 6-digit OTP code
   - Hash OTP using bcrypt
   - Store hashed OTP in Redis with 10-min TTL
   - Call SMS provider to send OTP
   - Return success response
   - Add rate limiting decorator (5 req/phone/hour)

5. **Auth Service – Verify Handler**
   - Create `POST /api/v1/auth/passenger/verify` endpoint
   - Parse and validate request (phone_number, otp_code)
   - Call verify service
   - Handle and return responses with proper status codes

6. **Auth Service – Verify Service Logic**
   - Retrieve hashed OTP from Redis
   - Compare provided OTP (bcrypt.Compare)
   - Handle invalid/expired OTP
   - Create user record via repository (name, phone, role='passenger')
   - Generate access JWT token (15 min)
   - Generate refresh JWT token (7 days)
   - Delete OTP from Redis
   - Return tokens + user data
   - Add rate limiting decorator (10 attempts/OTP/hour)

7. **Auth Service – Repository**
   - Implement `CheckPhoneExists(phone_number)` – query users table
   - Implement `CreatePassenger(name, phone_number)` – insert user, return user ID

8. **Middleware & Utilities**
   - Ensure rate limiting middleware is applied to both endpoints
   - Ensure error response handler returns consistent format
   - Add correlation ID tracking for OTP requests (for debugging)

9. **Testing**
   - Unit tests for signup service (happy path, duplicate phone, invalid format, SMS failure)
   - Unit tests for verify service (valid OTP, expired OTP, invalid OTP, DB error)
   - Unit tests for repository (phone exists, create user)
   - Integration tests for both endpoints
   - Rate limiting tests
   - Edge case tests (empty fields, special characters, concurrent requests)

10. **Documentation**
    - Update API docs with Swag annotations
    - Add endpoint examples in `/docs/api_endpoints.md`
    - Document OTP workflow in architecture docs
    - Add error codes to `/docs/error_codes.md`

## TODO Checklist

- [ ] Create/update `User` GORM model with name, phone_number, role fields
- [ ] Create database migration for signup schema changes
- [ ] Run migration and verify schema
- [ ] Create SignupRequest and SignupResponse DTOs
- [ ] Create VerifyRequest and VerifyResponse DTOs
- [ ] Implement `CheckPhoneExists()` repository method
- [ ] Implement `CreatePassenger()` repository method
- [ ] Write unit tests for repository methods
- [ ] Implement phone number validation utility (E.164 format)
- [ ] Implement name validation utility
- [ ] Implement signup service logic (OTP generation, SMS sending)
- [ ] Implement verify service logic (OTP verification, user creation, token generation)
- [ ] Write unit tests for signup service
- [ ] Write unit tests for verify service
- [ ] Create signup handler (`POST /api/v1/auth/passenger/signup`)
- [ ] Create verify handler (`POST /api/v1/auth/passenger/verify`)
- [ ] Add rate limiting to signup endpoint
- [ ] Add rate limiting to verify endpoint
- [ ] Write integration tests for signup endpoint
- [ ] Write integration tests for verify endpoint
- [ ] Add Swag documentation to handlers
- [ ] Test error scenarios (duplicate, invalid OTP, SMS failure)
- [ ] Test edge cases (special characters, long names, concurrent requests)
- [ ] Manual testing with mobile app
- [ ] Code review
- [ ] Update API documentation in `/docs`
- [ ] Merge to main branch

## Testing Strategy

### Unit Tests
- **Repository**: `TestCheckPhoneExists()`, `TestCheckPhoneExists_NotFound()`, `TestCreatePassenger()`, `TestCreatePassenger_Error()`
- **Service (Signup)**: `TestSignup_Success()`, `TestSignup_DuplicatePhone()`, `TestSignup_InvalidPhoneFormat()`, `TestSignup_InvalidName()`, `TestSignup_SMSFailure()`, `TestSignup_OTPGeneration()`
- **Service (Verify)**: `TestVerify_ValidOTP()`, `TestVerify_InvalidOTP()`, `TestVerify_ExpiredOTP()`, `TestVerify_UserCreationError()`, `TestVerify_TokenGeneration()`
- **Utilities**: `TestPhoneValidation_E164()`, `TestPhoneValidation_Invalid()`, `TestNameValidation_Valid()`, `TestNameValidation_Empty()`

### Integration Tests
- `TestSignup_EndToEnd()` – full signup flow
- `TestVerify_EndToEnd()` – full verification flow
- `TestRateLimit_Signup()` – exceed rate limit
- `TestRateLimit_Verify()` – exceed rate limit
- `TestDuplicatePhoneFlow()` – second signup with same phone

### Manual QA
- Mobile app: Test signup flow end-to-end
- Verify OTP is sent correctly
- Verify user is logged in after verification
- Verify tokens are valid (decode JWT)
- Test error messages on mobile UI
- Test with various phone formats (international, local)
- Test with Arabic names

### Load Testing
- Simulate 100 concurrent signup requests
- Verify rate limiting works correctly
- Check database performance with large user dataset

## Performance Considerations

- **OTP Lookup**: Redis key lookup (O(1)) – very fast
- **Phone Existence Check**: Index on `users.phone_number` for fast lookup
- **Token Generation**: JWT encoding is fast; no DB call needed
- **SMS Delivery**: Async (fire-and-forget) to avoid blocking request
- **Database Connection**: Use connection pooling; keep signup transaction short
- **Redis Memory**: OTP codes stored only 10 minutes; auto-cleanup via TTL
- **Concurrent Verifies**: Use Redis transactions or locks to prevent race conditions

## Definition of Done

- [x] All acceptance criteria met
- [x] All TODO items completed
- [x] Unit tests pass (coverage >80% for auth service)
- [x] Integration tests pass
- [x] Rate limiting working correctly
- [x] Error handling tested and documented
- [x] API docs updated with Swag annotations
- [x] Code reviewed and approved
- [x] No critical bugs in manual testing
- [x] Performance acceptable (<500ms for signup, <2s for verify)
- [x] Merged to main branch
- [x] Ready for mobile app integration
