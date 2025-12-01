# [FE-TASK_01] Passenger Phone-Number-Based Sign Up

**Backend Task:** [TASK_01_PASSENGER_SIGNUP_PHONE_OTP](../TASK_01_PASSENGER_SIGNUP_PHONE_OTP.md)  
**Created:** December 1, 2025  
**Status:** Todo  
**Priority:** High (MVP Feature)

---

## Overview

Implement the passenger signup flow in the mobile app (React Native/Expo). Users enter their name and phone number, receive an OTP via SMS, verify the OTP, and are automatically logged in upon successful verification.

### User Stories

- As a new passenger, I want to sign up using my phone number so that I can request rides
- As a new passenger, I want to verify my phone number with an OTP so that my account is secure
- As a new passenger, I want to be automatically logged in after signup so that I can start using the app immediately

---

## API Endpoints

### 1. Passenger Signup (Send OTP)

**Endpoint:** `POST /api/v1/auth/passenger/signup`

**Authentication:** Not Required

**Request Headers:**
```json
{
  "Content-Type": "application/json"
}
```

**Request Body:**
```typescript
interface SignupRequest {
  name: string;          // User's full name (2-100 chars)
  phone_number: string;  // E.164 format (e.g., +962791234567)
}
```

**Success Response (201 Created):**
```typescript
interface SignupResponse {
  message: string;       // "OTP sent to phone number"
  phone_number: string;  // Normalized phone number
  expires_in: number;    // OTP expiry in seconds (300 = 5 minutes)
}
```

**Error Responses:**
- `400 Bad Request` - Invalid phone number format or name
  ```json
  {
    "error": "BAD_REQUEST",
    "message": "Invalid phone number format. Expected E.164 format (e.g., +962791234567)",
    "details": null
  }
  ```
- `409 Conflict` - Phone number already registered
  ```json
  {
    "error": "CONFLICT",
    "message": "Phone number already registered. Please login instead.",
    "details": null
  }
  ```
- `429 Too Many Requests` - Rate limit exceeded (max 5 requests per phone per hour)
  ```json
  {
    "error": "RATE_LIMIT_EXCEEDED",
    "message": "Too many signup attempts. Please try again later.",
    "details": null
  }
  ```
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - SMS delivery failed

**Example Request:**
```javascript
const signupPassenger = async (name, phoneNumber) => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/auth/passenger/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: name,
        phone_number: phoneNumber
      })
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Signup failed:', error);
    throw error;
  }
};
```

**Example Response:**
```json
{
  "message": "OTP sent to phone number",
  "phone_number": "+962791234567",
  "expires_in": 300
}
```

---

### 2. Verify OTP and Complete Signup

**Endpoint:** `POST /api/v1/auth/passenger/verify`

**Authentication:** Not Required

**Request Headers:**
```json
{
  "Content-Type": "application/json"
}
```

**Request Body:**
```typescript
interface VerifyRequest {
  phone_number: string;  // E.164 format
  otp_code: string;      // 6-digit OTP
  name: string;          // User's name (same as signup)
}
```

**Success Response (200 OK):**
```typescript
interface VerifyResponse {
  access_token: string;   // JWT access token (15 min expiry)
  refresh_token: string;  // JWT refresh token (7 days expiry)
  user_id: string;        // UUID v4
  phone_number: string;   // User's phone number
  name: string;           // User's name
  role: string;           // "passenger"
  expires_in: number;     // Access token expiry in seconds (900 = 15 min)
}
```

**Error Responses:**
- `400 Bad Request` - Invalid OTP format or missing fields
  ```json
  {
    "error": "BAD_REQUEST",
    "message": "Invalid request body",
    "details": {
      "validation_error": "Key: 'VerifyRequest.OTPCode' Error:Field validation for 'OTPCode' failed on the 'len' tag"
    }
  }
  ```
- `401 Unauthorized` - Invalid or expired OTP
  ```json
  {
    "error": "UNAUTHORIZED",
    "message": "Invalid OTP code",
    "details": null
  }
  ```
- `409 Conflict` - User already exists (shouldn't happen in normal flow)
- `429 Too Many Requests` - Too many verification attempts (max 10 per hour)
  ```json
  {
    "error": "RATE_LIMIT_EXCEEDED",
    "message": "Too many verification attempts. Please request a new OTP.",
    "details": null
  }
  ```
- `500 Internal Server Error` - Server error

**Example Request:**
```javascript
const verifyOTP = async (phoneNumber, otpCode, name) => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/auth/passenger/verify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        phone_number: phoneNumber,
        otp_code: otpCode,
        name: name
      })
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message);
    }

    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Verification failed:', error);
    throw error;
  }
};
```

**Example Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "phone_number": "+962791234567",
  "name": "Ahmed Ali",
  "role": "passenger",
  "expires_in": 900
}
```

---

## Authentication & Authorization

### JWT Token Handling

After successful OTP verification, the backend returns:
- **Access Token**: Short-lived (15 minutes), used for all authenticated API calls
- **Refresh Token**: Long-lived (7 days), used to get new access tokens

**Storing Tokens:**
- Use Expo SecureStore for secure token storage
- Store both access and refresh tokens
- Store user profile data (user_id, name, phone_number, role)

**Using Tokens:**
All authenticated API requests must include the access token in the Authorization header:
```javascript
const makeAuthenticatedRequest = async (url, options = {}) => {
  const accessToken = await SecureStore.getItemAsync('access_token');
  
  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    }
  });

  return response;
};
```

**Token Refresh Flow:**
When access token expires (401 Unauthorized), refresh it using the refresh token:
```javascript
const refreshAccessToken = async () => {
  const refreshToken = await SecureStore.getItemAsync('refresh_token');
  
  // Implement refresh endpoint call here (future task)
  // For now, redirect to login on token expiry
};
```

---

## Data Models

```typescript
// Request DTOs
interface SignupRequest {
  name: string;
  phone_number: string;
}

interface VerifyRequest {
  phone_number: string;
  otp_code: string;
  name: string;
}

// Response DTOs
interface SignupResponse {
  message: string;
  phone_number: string;
  expires_in: number;
}

interface VerifyResponse {
  access_token: string;
  refresh_token: string;
  user_id: string;
  phone_number: string;
  name: string;
  role: string;
  expires_in: number;
}

interface ErrorResponse {
  error: string;
  message: string;
  details?: Record<string, any>;
}

// User Profile (stored locally after signup)
interface UserProfile {
  userId: string;
  name: string;
  phoneNumber: string;
  role: string;
  accessToken: string;
  refreshToken: string;
}
```

### Validation

**Name Validation:**
- Required
- Min length: 2 characters
- Max length: 100 characters
- Supports Arabic, English, numbers, spaces

**Phone Number Validation:**
- Required
- Must be in E.164 format: `+[country code][number]`
- Example: `+962791234567` (Jordan)
- Regex: `^\+[1-9]\d{1,14}$`

**OTP Code Validation:**
- Required
- Exactly 6 digits
- Only numbers allowed

---

## UI Components

### 1. SignupScreen

**Location:** `src/screens/auth/SignupScreen.tsx`

**Purpose:** Collect user's name and phone number, send OTP

**Components:**
- Text input for name (with validation)
- Phone input with country code picker
- "Send OTP" button (primary CTA)
- "Already have an account? Login" link
- Loading indicator
- Error message display

**Props:** None (initial screen)

**State:**
```typescript
const [name, setName] = useState('');
const [phoneNumber, setPhoneNumber] = useState('');
const [loading, setLoading] = useState(false);
const [error, setError] = useState<string | null>(null);
```

**User Interactions:**
1. User enters name
2. User enters phone number with country code (+962 for Jordan)
3. User taps "Send OTP" button
4. Validation runs (name length, phone format)
5. API call to signup endpoint
6. Navigate to OTPVerificationScreen on success
7. Show error message on failure

---

### 2. OTPVerificationScreen

**Location:** `src/screens/auth/OTPVerificationScreen.tsx`

**Purpose:** Verify OTP code and complete signup

**Components:**
- OTP input (6 digits, separate boxes or single input)
- Countdown timer showing OTP expiry (5 minutes)
- "Verify" button (primary CTA)
- "Resend OTP" button (enabled after countdown expires)
- Loading indicator
- Error message display
- Phone number display (masked, e.g., +962***4567)

**Props:**
```typescript
interface OTPVerificationScreenProps {
  phoneNumber: string;
  name: string;
  expiresIn: number;
}
```

**State:**
```typescript
const [otpCode, setOtpCode] = useState('');
const [loading, setLoading] = useState(false);
const [error, setError] = useState<string | null>(null);
const [timeRemaining, setTimeRemaining] = useState(expiresIn);
const [canResend, setCanResend] = useState(false);
```

**User Interactions:**
1. User enters 6-digit OTP
2. Auto-submit when 6 digits entered (optional)
3. User taps "Verify" button
4. API call to verify endpoint
5. Store tokens and user data in SecureStore
6. Navigate to main passenger screen on success
7. Show error message on failure
8. "Resend OTP" returns to SignupScreen or calls signup API again

---

## State Management

### Local Storage (SecureStore)

**Tokens:**
```javascript
import * as SecureStore from 'expo-secure-store';

// Save tokens after successful verification
await SecureStore.setItemAsync('access_token', verifyResponse.access_token);
await SecureStore.setItemAsync('refresh_token', verifyResponse.refresh_token);

// Save user profile
await SecureStore.setItemAsync('user_profile', JSON.stringify({
  userId: verifyResponse.user_id,
  name: verifyResponse.name,
  phoneNumber: verifyResponse.phone_number,
  role: verifyResponse.role
}));
```

### Context/Redux

**Auth Context:**
```typescript
interface AuthContextType {
  user: UserProfile | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  signup: (name: string, phoneNumber: string) => Promise<SignupResponse>;
  verifyOTP: (phoneNumber: string, otpCode: string, name: string) => Promise<void>;
  logout: () => Promise<void>;
}
```

**Auth State:**
- `user`: Current user profile or null
- `isAuthenticated`: Boolean indicating if user is logged in
- `isLoading`: Boolean for async operations
- `signup()`: Function to call signup API
- `verifyOTP()`: Function to verify OTP and log in user
- `logout()`: Function to clear tokens and user data

---

## Error Handling

| Error Code | Scenario | User Message | Action |
|------------|----------|--------------|--------|
| 400 | Invalid phone format | "Please enter a valid phone number in international format (e.g., +962791234567)" | Show inline error, highlight phone input |
| 400 | Invalid name | "Name must be between 2 and 100 characters" | Show inline error, highlight name input |
| 409 | Phone already registered | "This phone number is already registered. Please log in instead." | Show error with "Go to Login" button |
| 429 | Signup rate limit | "Too many attempts. Please try again in 1 hour." | Show error with countdown |
| 429 | Verify rate limit | "Too many verification attempts. Please request a new OTP." | Show error with "Request New OTP" button |
| 401 | Invalid OTP | "Invalid verification code. Please try again." | Show inline error, clear OTP input |
| 401 | Expired OTP | "Verification code has expired. Please request a new one." | Show error with "Resend OTP" button |
| 503 | SMS delivery failed | "Failed to send verification code. Please try again." | Show error with "Retry" button |
| 500 | Server error | "Something went wrong. Please try again later." | Show error with "Retry" button |
| Network error | No internet | "No internet connection. Please check your network." | Show error, retry button |

**Error Display:**
- Use Toast notifications for non-critical errors
- Use inline error messages for form validation
- Use modal/alert for critical errors (e.g., phone already registered)
- All error messages should be in both English and Arabic

---

## Loading & Empty States

### Loading States

**Signup Screen:**
- Show spinner on "Send OTP" button when API call in progress
- Disable form inputs during loading
- Prevent multiple submissions

**OTP Verification Screen:**
- Show spinner on "Verify" button when API call in progress
- Disable OTP input during verification
- Show loading overlay if navigating after success

### Countdown Timer

**OTP Expiry Timer:**
- Display countdown: "Code expires in 4:32"
- Update every second
- Show warning when < 1 minute remaining (change color to red)
- When expired, show "Code expired" and enable "Resend OTP" button

### Success State

**After Successful Verification:**
- Show success message: "Account created successfully!"
- Brief success animation (checkmark, green background)
- Auto-navigate to main passenger screen after 1-2 seconds

---

## Testing Requirements

### Unit Tests

**API Service Functions:**
- [ ] `signupPassenger()` - success case
- [ ] `signupPassenger()` - network error
- [ ] `signupPassenger()` - validation error (400)
- [ ] `signupPassenger()` - conflict error (409)
- [ ] `verifyOTP()` - success case
- [ ] `verifyOTP()` - invalid OTP (401)
- [ ] `verifyOTP()` - expired OTP (401)
- [ ] `verifyOTP()` - rate limit (429)

**Validation Functions:**
- [ ] `validatePhoneNumber()` - valid numbers
- [ ] `validatePhoneNumber()` - invalid numbers
- [ ] `validateName()` - valid names (English, Arabic)
- [ ] `validateName()` - too short, too long
- [ ] `normalizePhoneNumber()` - various formats

**Token Management:**
- [ ] `saveTokens()` - stores tokens in SecureStore
- [ ] `getAccessToken()` - retrieves token
- [ ] `clearTokens()` - removes all tokens

### Component Tests

**SignupScreen:**
- [ ] Renders all input fields correctly
- [ ] Name validation works
- [ ] Phone validation works
- [ ] Shows error messages
- [ ] Disables button during loading
- [ ] Navigates to OTP screen on success

**OTPVerificationScreen:**
- [ ] Renders OTP input correctly
- [ ] Countdown timer updates every second
- [ ] Shows expiry warning
- [ ] Validates OTP format (6 digits)
- [ ] Shows error messages
- [ ] Navigates to main screen on success
- [ ] Resend button enabled after expiry

### Integration Tests

- [ ] Complete signup flow (signup → verify → logged in)
- [ ] Error handling flow (invalid phone → fix → success)
- [ ] Token storage and retrieval
- [ ] Navigation after successful verification
- [ ] Logout clears tokens and returns to auth screen

### Manual Testing Checklist

- [ ] Test with Jordan phone numbers (+962)
- [ ] Test with other country codes (+1, +44, etc.)
- [ ] Test with invalid phone formats
- [ ] Test name with Arabic characters (أحمد علي)
- [ ] Test name with special characters
- [ ] Test OTP verification with correct code
- [ ] Test OTP verification with wrong code
- [ ] Test OTP expiry (wait 5 minutes)
- [ ] Test rate limiting (multiple signup attempts)
- [ ] Test network offline/online scenarios
- [ ] Test on both iOS and Android
- [ ] Test on different screen sizes
- [ ] Test with screen reader (accessibility)

---

## Design Notes

### Colors (THEB Brand)

- **Wolf Black**: `#0D0D0D` – backgrounds, headers
- **Desert Gold**: `#D4A048` – primary buttons, CTAs, accents
- **Pure White**: `#FFFFFF` – text on dark backgrounds, input backgrounds
- **Sand Gray**: `#B6B0A2` – secondary text, disabled states
- **Success Green**: `#41A45A` – success messages, checkmarks
- **Alert Red**: `#D9534F` – error messages, warnings

### Typography

- **Primary Font**: Cairo (Arabic) / Inter (English)
- **Headings**: Bold, 24-32px
- **Body Text**: Regular, 16px
- **Input Labels**: Medium, 14px
- **Button Text**: Bold, 16px

### Spacing

- **Container Padding**: 20px
- **Input Vertical Spacing**: 16px
- **Button Height**: 52px
- **Border Radius**: 8px (buttons, inputs)

### OTP Input Styling

- 6 separate boxes or single input with mask
- Each box: 48x48px, centered digit
- Active box: Desert Gold border
- Filled box: Desert Gold background, White text
- Error state: Alert Red border

### Design Assets

- Phone number input with country flag picker
- OTP input with auto-focus and auto-advance
- Loading spinner matching brand colors
- Success checkmark animation

---

## Implementation Notes

### Phone Number Input

Use a library like `react-native-phone-number-input` for:
- Country code picker with flags
- Auto-formatting
- Validation
- Jordan (+962) as default country

### OTP Input

Use a library like `react-native-otp-textinput` or `react-native-confirmation-code-field` for:
- Auto-focus on first digit
- Auto-advance to next box
- Paste support (full OTP from SMS)
- Auto-submit on completion (optional)

### SMS Auto-Read (Android Only)

Use `react-native-otp-verify` to auto-read OTP from SMS on Android:
- Request SMS permission
- Listen for incoming SMS
- Parse OTP code automatically
- Fill OTP input

### Keyboard Behavior

- Use `KeyboardAvoidingView` to prevent keyboard from covering inputs
- Dismiss keyboard on scroll or tap outside
- Numeric keyboard for OTP input
- Phone keyboard for phone number input

### Accessibility

- Use proper ARIA labels for inputs
- Support screen readers (VoiceOver, TalkBack)
- Keyboard navigation support
- High contrast mode support
- Font scaling support

### Performance

- Debounce API calls for signup (prevent double-submission)
- Optimize countdown timer (use `requestAnimationFrame` or intervals)
- Lazy load country picker data
- Cache country codes locally

---

## Acceptance Criteria

### Functional Requirements

- [ ] User can enter name (2-100 chars, supports Arabic)
- [ ] User can enter phone number with country code picker
- [ ] Phone number auto-formats as user types
- [ ] Name and phone validation works with inline errors
- [ ] "Send OTP" button disabled until valid input
- [ ] API call to signup endpoint works correctly
- [ ] OTP sent message displayed with phone number
- [ ] Navigate to OTP verification screen
- [ ] OTP input accepts 6 digits only
- [ ] Countdown timer shows OTP expiry time
- [ ] "Verify" button disabled until 6 digits entered
- [ ] API call to verify endpoint works correctly
- [ ] Tokens stored securely in SecureStore
- [ ] User profile stored locally
- [ ] Navigate to main passenger screen after verification
- [ ] "Resend OTP" button enabled after expiry
- [ ] Error messages displayed for all error scenarios
- [ ] Rate limiting handled gracefully

### UI/UX Requirements

- [ ] UI matches THEB brand design (colors, typography)
- [ ] Smooth animations and transitions
- [ ] Loading indicators shown during API calls
- [ ] Form inputs disabled during loading
- [ ] Error messages clear and actionable
- [ ] Success feedback after verification
- [ ] Responsive on all screen sizes (iPhone SE to iPad)
- [ ] Works in portrait and landscape
- [ ] Supports dark mode (if app has dark mode)
- [ ] RTL layout for Arabic language

### Performance Requirements

- [ ] Signup API call completes in < 2 seconds
- [ ] Verify API call completes in < 2 seconds
- [ ] No UI lag when typing
- [ ] Countdown timer updates smoothly
- [ ] App doesn't crash on network errors
- [ ] Handles slow network gracefully

### Security Requirements

- [ ] Tokens stored in SecureStore (encrypted)
- [ ] No tokens logged to console
- [ ] No sensitive data in error messages
- [ ] OTP not visible in UI (masked input optional)
- [ ] Phone number masked in verification screen

### Accessibility Requirements

- [ ] All inputs have proper labels
- [ ] Screen reader support (VoiceOver, TalkBack)
- [ ] Keyboard navigation works
- [ ] Color contrast meets WCAG AA standards
- [ ] Font scaling supported (up to 200%)

---

## Definition of Done

- [ ] Code implemented and follows React Native best practices
- [ ] All acceptance criteria met
- [ ] Unit tests written and passing (>80% coverage)
- [ ] Component tests written and passing
- [ ] Integration tests written and passing
- [ ] API integration tested with real backend
- [ ] Error handling verified for all scenarios
- [ ] UI tested on iOS and Android simulators
- [ ] UI tested on physical devices (iOS and Android)
- [ ] Performance acceptable (no lag, smooth animations)
- [ ] Accessibility verified with screen reader
- [ ] Code reviewed and approved
- [ ] Documentation updated (if needed)
- [ ] No console warnings or errors
- [ ] Ready for QA testing

---

## Notes

- OTP codes are logged to backend console in development mode for testing
- In production, integrate with SMS provider (Twilio, AWS SNS, or Jordanian provider)
- Consider adding auto-SMS read feature for Android (future enhancement)
- Consider adding biometric authentication after signup (future enhancement)
- For testing, use test phone numbers that don't trigger real SMS sends
- Backend OTP expiry is 5 minutes (300 seconds)
- Backend rate limits: 5 signups per phone per hour, 10 verifications per hour
