# AUTO-LOGIN-OTP – Auto Login After OTP Verification

## Overview
Enhance the OTP verification flow to automatically create an authenticated session and return JWT tokens upon successful verification. This eliminates an extra login step, providing a seamless signup experience where passengers are immediately logged in after OTP validation and can access the app without manual re-authentication.

## Source
- Type: Story
- Priority: High (MVP feature - directly related to passenger signup)
- Labels: Auth, Passenger, UX, MVP
- Reporter: Product Team
- Linked Dependencies: TASK_PASSENGER_SIGNUP_PHONE_OTP (must complete signup first)

## Scope

### In Scope
- Modify OTP verification endpoint to generate and return JWT tokens (access + refresh)
- Implement session creation logic during OTP verification
- Ensure tokens are issued immediately upon successful OTP validation
- Return all necessary user data for authenticated session (user_id, phone, name, role)
- Token response includes expiry times and token type metadata
- Secure token generation using HMAC-SHA256
- Proper token validation on subsequent requests
- Clear error responses for invalid/expired OTP (no tokens issued)

### Out of Scope
- Mobile app token storage (React Native side - use secure key storage libraries)
- HomeScreen navigation logic (mobile app responsibility)
- Token refresh endpoint implementation (handle in separate refresh token story)
- Multi-device session management
- Logout flow (handle in separate logout story)

## Requirements

### Functional
- Generate access JWT token (15-minute validity) upon successful OTP verification
- Generate refresh JWT token (7-day validity) upon successful OTP verification
- Return both tokens in verification response
- Include token metadata (token_type: "Bearer", expires_in: seconds)
- Include user context (user_id, phone_number, name, role)
- Ensure tokens are signed with application JWT secret
- Tokens contain necessary claims (sub: user_id, phone, role, iat, exp)
- If OTP invalid → return 401 with error message, no tokens issued
- If OTP expired → return 401 with specific "OTP_EXPIRED" error code
- Invalidate OTP immediately after successful verification (delete from Redis)
- Ensure no session can be created from invalid/expired OTP

### Non-Functional
- Token generation must complete within <100ms
- Token response must include all necessary fields for mobile app session
- Tokens must be cryptographically secure
- Token claims must be standardized (JWT RFC 7519 compliant)
- Support concurrent token requests from same user
- Tokens must be stateless (no server-side session store required)
- Token format must be parseable by standard JWT libraries

## Data & Contracts

### Modified API Endpoint

#### POST /api/v1/auth/passenger/verify (Enhanced)

**Request:**
```json
{
  "phone_number": "+962791234567",
  "otp_code": "123456"
}
```

**Response (200 OK) - Success:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "phone_number": "+962791234567",
    "name": "Ahmed Ali",
    "role": "passenger",
    "created_at": "2025-12-01T10:30:00Z"
  }
}
```

**Response (401 Unauthorized) - Invalid OTP:**
```json
{
  "error": "INVALID_OTP",
  "message": "OTP code is incorrect",
  "attempts_remaining": 9
}
```

**Response (401 Unauthorized) - Expired OTP:**
```json
{
  "error": "OTP_EXPIRED",
  "message": "OTP code has expired. Please request a new OTP.",
  "expires_at": "2025-12-01T10:40:00Z"
}
```

**Response (429 Too Many Requests) - Rate Limited:**
```json
{
  "error": "RATE_LIMIT_EXCEEDED",
  "message": "Too many verification attempts. Please try again later.",
  "retry_after": 3600
}
```

### JWT Token Structure

**Access Token Claims:**
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "+962791234567",
  "role": "passenger",
  "iat": 1701424200,
  "exp": 1701427800,
  "type": "access"
}
```

**Refresh Token Claims:**
```json
{
  "sub": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "+962791234567",
  "iat": 1701424200,
  "exp": 1702029000,
  "type": "refresh"
}
```

### Session Lifecycle

1. **Signup**: User enters name + phone → OTP sent
2. **Verify**: User enters OTP → Tokens generated → Session created
3. **Use App**: Mobile stores tokens in secure storage
4. **API Calls**: Mobile includes `Authorization: Bearer {access_token}` header
5. **Token Expiry**: After 15 min, mobile uses refresh_token to get new access_token
6. **Logout**: Mobile clears tokens from storage (separate flow)

## UI / Components

Mobile app responsibility (not backend):
- **OTP Verification Screen**: Display OTP input, "Verify" button, timer countdown
- **Success State**: Show confirmation briefly, then navigate to Home
- **Error State**: Display error message (invalid OTP, expired, rate limit)
- **Resend OTP**: Allow retry if OTP expired
- **Token Storage**: Use secure storage (Keychain for iOS, Keystore for Android)
- **Auto-Login**: On app launch, check if tokens exist in secure storage → skip login
- **Token Refresh**: Silently refresh access token when expired (before API call fails)

Backend provides:
- Tokens in response
- Clear error codes for mobile to handle UI states

## State & Logic

### OTP Verification Flow (Enhanced)

**Handler** (thin wrapper):
```
1. Parse request (phone_number, otp_code)
2. Validate input format
3. Call service.VerifyOTP()
4. Return response with tokens or error
```

**Service** (business logic):
```
1. Retrieve OTP from Redis → get {hash, attempts, created_at}
2. If OTP not found → return OTP_EXPIRED error (401)
3. If attempts >= 10 → return RATE_LIMIT_EXCEEDED (429)
4. Compare provided OTP with stored hash (bcrypt.Compare)
   - If invalid → increment attempts in Redis, return INVALID_OTP (401)
   - If valid → continue
5. Retrieve user from DB by phone_number
6. Generate access_token (JWT, 15 min)
   - Claims: {sub: user_id, phone, role, iat, exp, type: "access"}
7. Generate refresh_token (JWT, 7 days)
   - Claims: {sub: user_id, phone, iat, exp, type: "refresh"}
8. Delete OTP from Redis (invalidate)
9. Create response with tokens + user data
10. Return 200 with tokens
```

**Repository**:
```
- GetUserByPhone(phone_number) → User struct or nil
```

**Token Generation**:
```
- Sign token with HMAC-SHA256 using JWT_SECRET env var
- Include standard claims (iat, exp, sub)
- Include custom claims (phone, role, type)
- Return encoded token string
```

### Error Handling

| Error | Code | HTTP Status | Action |
|-------|------|-------------|--------|
| OTP not found | OTP_EXPIRED | 401 | Delete Redis key, suggest resend |
| OTP invalid | INVALID_OTP | 401 | Increment counter, show attempts left |
| Too many attempts | RATE_LIMIT_EXCEEDED | 429 | Block for 1 hour, log security event |
| User not found | USER_NOT_FOUND | 401 | Rollback (should not happen if signup succeeded) |
| Token generation error | INTERNAL_ERROR | 500 | Log error, suggest retry |
| Database error | INTERNAL_ERROR | 500 | Log error, do not expose details |

### State Storage

- **OTP**: Redis (temporary, hashed)
  - Key: `otp:{phone_number}`
  - Value: `{hash, attempts, created_at, expires_at}`
  - TTL: 10 minutes (auto-cleanup)

- **User Session**: None (stateless JWT)
  - Tokens issued to client
  - Client stores in secure local storage
  - Client sends in Authorization header

- **Rate Limiting**: Redis counters
  - Key: `otp:verify:{phone_number}`
  - Counter incremented on each failed attempt
  - TTL: 1 hour

## Acceptance Criteria

- [x] OTP verification endpoint returns access_token and refresh_token
- [x] Tokens are JWT-encoded and signed with application secret
- [x] Access token has 15-minute validity
- [x] Refresh token has 7-day validity
- [x] Tokens include required claims (sub, phone, role, iat, exp)
- [x] Response includes token_type: "Bearer"
- [x] Response includes expires_in (seconds)
- [x] Response includes full user context (id, phone, name, role, created_at)
- [x] Invalid OTP returns 401 with INVALID_OTP error code
- [x] Expired OTP returns 401 with OTP_EXPIRED error code
- [x] Failed verification attempts are rate-limited (max 10 per hour)
- [x] Rate limit exceeded returns 429 with retry_after header
- [x] OTP is deleted from Redis immediately after successful verification
- [x] No session can be created from invalid/expired OTP
- [x] Token generation completes within <100ms
- [x] Concurrent verification requests are handled safely
- [x] Error responses include attempts_remaining for UX feedback

## Edge Cases

- **Concurrent verify requests**: First valid request succeeds; subsequent requests get invalid OTP error (already deleted)
- **User tries verify after signup timeout**: OTP expired, return OTP_EXPIRED with resend option
- **Invalid phone format in verify**: Reject with 400 Bad Request
- **OTP code with leading/trailing spaces**: Trim before validation
- **Multiple OTP requests sent**: Latest OTP is valid; previous ones replaced in Redis
- **Verify with wrong phone number**: Return 401 (phone not found or OTP invalid)
- **Token generation fails**: Return 500, do not create partial session
- **User deleted between signup and verify**: Handle gracefully, return USER_NOT_FOUND (rare edge case)
- **Database connection drops during verify**: Rollback, return 500
- **OTP attempt counter reaches limit**: Block further attempts, require full hour wait
- **Mobile sends expired token to API**: Auth middleware rejects with 401 Unauthorized
- **Mobile loses tokens (app crash)**: User must log in again (separate flow)

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Token leakage (network sniffing) | Session hijacking, account takeover | Use HTTPS only in production, mobile uses secure storage (Keychain/Keystore) |
| Weak JWT secret | Tokens forged/modified by attacker | Use strong random JWT_SECRET (32+ bytes), rotate periodically, never commit to repo |
| Race condition: OTP deleted before response sent | Valid user gets error after successful verify | Use atomic Redis operations (GETDEL command), log if race detected |
| Token claims not validated | Malformed tokens accepted | Validate all claims on each API request, check exp time, verify signature |
| Concurrent requests with same OTP | Multiple sessions created | Use Redis transactions or atomic operations, ensure only first succeeds |
| Token generation slow | API response timeout, poor UX | Cache JWT_SECRET in memory, benchmark token generation (<100ms target) |
| OTP brute force via rate limit bypass | Account takeover | Implement IP-based rate limiting in addition to phone-based, use captcha on signup |
| Refresh token not implemented yet | Old tokens never expire, security risk | Implement refresh endpoint before this story completes, short access token life (15 min) |
| Error messages expose too much | Information disclosure | Use generic errors in production (e.g., "Invalid credentials"), log details server-side |
| Mobile doesn't store tokens securely | Tokens visible in plaintext on device | Document secure storage requirement, use WebSecure Android/Keychain iOS |

## Open Questions / Assumptions

- **Q**: What is the JWT_SECRET? **A**: Assumed to be a 32+ byte random string stored in environment variables (not committed)
- **Q**: Should tokens be stored in database/Redis for revocation? **A**: No - assume stateless JWT for MVP, revocation via blacklist in future if needed
- **Q**: What if user changes phone after signup? **A**: Out of scope; assume phone immutable for now
- **Q**: Should we support multiple devices for same user? **A**: No - single device session for MVP; multi-device handled in future
- **Q**: How long should rate limit block last? **A**: Assumed 1 hour (can be configurable)
- **Q**: Should we log all failed OTP attempts? **A**: Yes - for security audits and fraud detection
- **Q**: Can tokens be manually refreshed immediately? **A**: Yes - refresh_token can be used to get new access_token anytime (implement in separate story)
- **Q**: Should OTP expire if user requests new one? **A**: Yes - new OTP replaces old one, invalidates old in Redis

## Story Points

**3 points** – Reason: Low to medium complexity. Builds directly on signup story. Main work is token generation (using JWT library), minimal new database queries, straightforward error handling. Requires careful testing of edge cases (concurrent requests, token validation) but limited scope change to existing OTP verification endpoint.

## Implementation Plan

1. **Update Configuration**
   - Add JWT_SECRET to config (load from env var)
   - Add token expiry times (access: 15 min, refresh: 7 days)
   - Validate JWT_SECRET on startup (must be non-empty)

2. **JWT Token Service**
   - Create `TokenService` interface with methods:
     - `GenerateAccessToken(userID, phone, role) → token string`
     - `GenerateRefreshToken(userID, phone) → token string`
   - Implement token generation using standard Go JWT library (github.com/golang-jwt/jwt)
   - Implement token claims struct with required fields
   - Handle signing with HMAC-SHA256

3. **Update Verify Response DTO**
   - Create `VerifyResponse` struct with fields:
     - `AccessToken string`
     - `RefreshToken string`
     - `TokenType string` (default: "Bearer")
     - `ExpiresIn int64` (seconds)
     - `User UserResponse`
   - Create `UserResponse` struct with:
     - `ID uuid.UUID`
     - `PhoneNumber string`
     - `Name string`
     - `Role string`
     - `CreatedAt time.Time`

4. **Update Auth Service – Verify Logic**
   - Modify `VerifyOTP()` service method to:
     - Keep existing OTP validation logic
     - After valid OTP, retrieve user from DB
     - Call TokenService to generate access_token
     - Call TokenService to generate refresh_token
     - Return new `VerifyResponse` with tokens + user data
     - Delete OTP from Redis only after successful token generation

5. **Update Verify Handler**
   - Modify response to use new `VerifyResponse` struct
   - Ensure error responses return appropriate HTTP status codes
   - Include error code and message for mobile app error handling
   - Add attempts_remaining field for rate limit feedback

6. **Error Response Standardization**
   - Create error response DTO:
     - `Error string` (error code: INVALID_OTP, OTP_EXPIRED, RATE_LIMIT_EXCEEDED)
     - `Message string` (user-friendly description)
     - `AttemptsRemaining int` (optional)
     - `RetryAfter int` (optional, for rate limiting)
   - Update error handler middleware to use consistent format

7. **Rate Limiting Enhancement**
   - Add attempt counter to OTP verify (track failed attempts per phone)
   - Block after 10 failed attempts (return 429)
   - Attempt counter TTL: 1 hour
   - Return retry_after header on 429

8. **Testing**
   - Unit tests for TokenService (token generation, claim validation)
   - Unit tests for verify service (OTP validation + token generation)
   - Unit tests for error scenarios (invalid OTP, expired OTP, rate limit)
   - Integration tests for full verify flow (happy path + errors)
   - Concurrent request tests (multiple verify requests same OTP)
   - Token validation tests (decode JWT, verify claims, check expiry)
   - Rate limiting tests (exceed attempt limit)

9. **Documentation**
   - Update Swag API docs for verify endpoint
   - Document response structure (tokens, user data)
   - Document error codes and status codes
   - Add example JWT payload to API docs
   - Document mobile app requirements (secure token storage)
   - Add flowchart: Signup → Verify → Auto-Login → Home

10. **Security Review**
    - Ensure JWT_SECRET is strong
    - Verify HTTPS is enforced in production
    - Check that tokens don't expose sensitive data
    - Audit rate limiting effectiveness
    - Review token expiry times

## TODO Checklist

- [ ] Read JWT library documentation (golang-jwt/jwt)
- [ ] Create config entries for JWT_SECRET, token expiry times
- [ ] Add JWT_SECRET validation on startup
- [ ] Create TokenService interface and implementation
- [ ] Implement GenerateAccessToken() method
- [ ] Implement GenerateRefreshToken() method
- [ ] Create JWT claims struct with required fields
- [ ] Create VerifyResponse DTO
- [ ] Create UserResponse DTO
- [ ] Update error response DTO format
- [ ] Modify VerifyOTP() service to generate tokens
- [ ] Modify verify handler to use new response structure
- [ ] Add attempt counter tracking in Redis
- [ ] Implement rate limiting for verify endpoint
- [ ] Add HTTP 429 response for rate limit exceeded
- [ ] Write unit tests for TokenService
- [ ] Write unit tests for verify service with tokens
- [ ] Write unit tests for invalid OTP handling
- [ ] Write unit tests for expired OTP handling
- [ ] Write unit tests for rate limiting
- [ ] Write integration tests for full verify flow
- [ ] Write integration tests for error scenarios
- [ ] Test concurrent verify requests
- [ ] Test token decoding and claim validation
- [ ] Update Swag documentation for verify endpoint
- [ ] Add endpoint example to API docs
- [ ] Document error codes and HTTP status codes
- [ ] Create token structure example in docs
- [ ] Document mobile secure storage requirement
- [ ] Code review
- [ ] Security audit (JWT_SECRET, HTTPS, rate limit)
- [ ] Manual testing with mobile app
- [ ] Load testing (concurrent token generation)
- [ ] Merge to main branch

## Testing Strategy

### Unit Tests
- **TokenService**: `TestGenerateAccessToken()`, `TestGenerateAccessToken_ValidClaims()`, `TestGenerateRefreshToken()`, `TestGenerateAccessToken_InvalidSecret()`, `TestTokenExpiry()`
- **Verify Service**: `TestVerifyOTP_ValidOTP_ReturnsTokens()`, `TestVerifyOTP_InvalidOTP()`, `TestVerifyOTP_ExpiredOTP()`, `TestVerifyOTP_UserNotFound()`, `TestVerifyOTP_AttemptCounter()`, `TestVerifyOTP_RateLimit()`, `TestVerifyOTP_ConcurrentRequests()`
- **Response Mapping**: `TestVerifyResponse_IncludesAllFields()`, `TestVerifyResponse_UserData()`, `TestVerifyResponse_TokenMetadata()`

### Integration Tests
- `TestVerifyFlow_FullSignupToLogin()` – Signup → Verify → Tokens → Use in API call
- `TestVerifyFlow_InvalidOTPError()` – Wrong OTP returns 401
- `TestVerifyFlow_ExpiredOTPError()` – Expired OTP returns 401
- `TestVerifyFlow_RateLimitExceeded()` – 10+ failed attempts returns 429
- `TestVerifyFlow_ConcurrentVerifications()` – Multiple requests, only first succeeds
- `TestVerifyFlow_TokenValidation()` – Decode tokens, verify claims, check signature
- `TestVerifyFlow_OTPDeletedAfterSuccess()` – Verify OTP can't be reused

### Manual QA
- Mobile app: Full signup → verify → auto-login → home navigation flow
- Verify tokens are issued in response
- Verify mobile can decode and store tokens
- Verify tokens work in subsequent API calls
- Test invalid OTP shows appropriate error
- Test expired OTP shows resend option
- Test rate limiting after 10 failed attempts
- Test concurrent verify requests (simulate with mobile)
- Test token expiry after 15 minutes

### Load Testing
- Generate 1000 concurrent verify requests
- Measure token generation time (<100ms target)
- Verify no race conditions with concurrent requests
- Check Redis performance under load

## Performance Considerations

- **Token Generation**: JWT encoding is fast; ensure JWT_SECRET loaded in memory (not file I/O per request)
- **OTP Deletion**: Use Redis GETDEL for atomic read+delete (prevent race conditions)
- **User Lookup**: Index on `users.phone_number` for fast retrieval
- **Rate Limit Counter**: Redis counter increments are O(1), fast
- **Concurrent Requests**: Minimize database locks, use Redis atomic operations
- **Token Size**: JWT tokens are URL-safe strings; no payload size issues
- **Response Time Target**: <500ms for full verify flow (OTP validation + token generation + user lookup)

## Definition of Done

- [x] Tokens (access + refresh) returned in verify response
- [x] Tokens are valid JWT with correct claims
- [x] Access token has 15-minute expiry
- [x] Refresh token has 7-day expiry
- [x] Error responses include error codes for mobile handling
- [x] Rate limiting enforced (max 10 attempts per hour)
- [x] OTP deleted immediately after successful verification
- [x] Concurrent requests handled safely (no race conditions)
- [x] Unit tests pass (>80% coverage for auth service)
- [x] Integration tests pass
- [x] Manual testing with mobile app successful
- [x] API docs updated with Swag annotations
- [x] Error codes documented
- [x] Token structure documented for mobile developers
- [x] Security review completed
- [x] Code reviewed and approved
- [x] Merged to main branch
- [x] Ready for mobile app to implement token storage
