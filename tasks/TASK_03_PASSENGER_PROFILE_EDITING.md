# PASSENGER-PROFILE-EDITING – Passenger Profile Editing

## Overview
Enable passengers to view and edit their personal profile details in the app. Editable fields include avatar URL, email, mobile number, date of birth, gender, and trusted phone numbers. All changes should sync instantly with the backend, supporting partial updates and validation.

## Source
- Type: Story
- Priority: High (core user feature)
- Labels: Passenger, Profile, Edit, MVP
- Reporter: Product Team
- Linked Dependencies: None

## Scope

### In Scope
- API endpoint to update passenger profile (`PATCH /api/v1/passenger/profile`)
- Support for partial updates (only changed fields sent)
- Editable fields: avatar_url, email, phone_number, date_of_birth, gender, trusted_phones (array)
- Trusted Contacts management (add/delete trusted phone numbers)
- Input validation (email format, phone format, DOB min age, gender values)
- Immediate backend sync and response
- Structured storage for trusted phone numbers

### Out of Scope
- Profile page layout/UI (handled by mobile app)
- Avatar image upload (only URL supported)
- Multi-language support for profile fields
- Admin editing of passenger profiles
- Profile deletion (handled separately)

## Requirements

### Functional
- Allow passengers to update any editable profile field
- Support PATCH requests for partial updates
- Validate email format, phone number format (E.164), DOB (min age 16), gender (enum: male, female, other)
- Allow adding/removing trusted phone numbers (array of objects)
- Save changes to backend immediately
- Return updated profile in response
- Show placeholders for empty fields in UI
- Reject invalid updates with clear error messages

### Non-Functional
- Profile update must complete within <500ms
- All updates must be atomic (transactional)
- Support concurrent updates (no data loss)
- Secure input validation (prevent injection)
- Consistent error response format
- Audit log for profile changes (for future compliance)

## Data & Contracts

### API Endpoint

#### PATCH /api/v1/passenger/profile
**Request (partial):**
```json
{
  "avatar_url": "https://example.com/avatar.jpg",
  "email": "ahmed@example.com",
  "phone_number": "+962791234567",
  "date_of_birth": "2000-01-01",
  "gender": "male",
  "trusted_phones": [
    { "name": "Mother", "phone": "+962790000001" },
    { "name": "Brother", "phone": "+962790000002" }
  ]
}
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "avatar_url": "https://example.com/avatar.jpg",
  "email": "ahmed@example.com",
  "phone_number": "+962791234567",
  "date_of_birth": "2000-01-01",
  "gender": "male",
  "trusted_phones": [
    { "name": "Mother", "phone": "+962790000001" },
    { "name": "Brother", "phone": "+962790000002" }
  ],
  "updated_at": "2025-12-01T11:00:00Z"
}
```

**Errors:**
- 400: Invalid email/phone/DOB/gender
- 409: Phone/email already in use
- 401: Unauthorized
- 500: Internal server error

### Database Schema
- `avatar_url` (string, nullable)
- `email` (string, unique, nullable)
- `phone_number` (string, unique, E.164)
- `date_of_birth` (date, nullable)
- `gender` (enum: male, female, other)
- `trusted_phones` (JSON array of objects: {name, phone})
- `updated_at` (timestamp)

## UI / Components
- **Edit Profile Screen**: Inputs for all editable fields
- **Trusted Contacts**: Add/delete trusted phone numbers
- **Validation**: Show errors for invalid input
- **Save Button**: Sends PATCH request
- **Placeholders**: For empty fields
- **Immediate update**: UI reflects changes after save

## State & Logic
- **Handler**: Parses PATCH request, validates input, calls service
- **Service**: Validates fields, applies partial updates, manages trusted contacts
- **Repository**: Updates only changed fields in DB, ensures atomicity
- **Trusted Phones**: Stored as JSON array, validated for format and uniqueness
- **Validation**: Email regex, phone E.164, DOB >= 16 years, gender in allowed values
- **Response**: Returns updated profile

## Acceptance Criteria
- [x] PATCH `/api/v1/passenger/profile` supports partial updates
- [x] Editable fields: avatar_url, email, phone_number, date_of_birth, gender, trusted_phones
- [x] Trusted contacts can be added/removed
- [x] Input validation for all fields
- [x] Invalid updates rejected with error
- [x] Changes saved instantly to backend
- [x] Response includes updated profile
- [x] UI reflects changes immediately
- [x] Placeholders shown for empty fields
- [x] No restart required after update
- [x] Consistent error response format

## Edge Cases
- **Empty fields**: Show placeholder, allow clearing
- **Duplicate phone/email**: Reject with 409 error
- **Invalid email/phone/DOB/gender**: Reject with 400 error
- **Concurrent updates**: Last write wins, atomic transaction
- **Trusted phone with invalid format**: Reject
- **Trusted phone duplicate**: Reject
- **Remove all trusted phones**: Allow (empty array)
- **Avatar URL invalid**: Reject
- **DOB under 16**: Reject
- **Gender not in allowed values**: Reject
- **Partial update (only one field)**: Accept

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Invalid input accepted | Data corruption | Strict validation, unit tests |
| Race condition on concurrent updates | Data loss | Use DB transactions, last write wins |
| Trusted phone duplicates | Confusing contacts | Validate uniqueness in array |
| Email/phone collision | Account takeover | Unique constraints, error handling |
| Slow DB writes | Poor UX | Optimize queries, index fields |
| Injection attacks | Security breach | Sanitize all inputs |
| Partial update fails | Inconsistent data | Use atomic PATCH, rollback on error |

## Open Questions / Assumptions
- **Q**: Can avatar be uploaded? **A**: No, only URL supported
- **Q**: Is email required? **A**: Optional, can be empty
- **Q**: Can phone number be changed? **A**: Yes, must be unique and valid
- **Q**: How many trusted phones allowed? **A**: Assume up to 5
- **Q**: What gender values allowed? **A**: male, female, other
- **Q**: Is DOB required? **A**: Optional, but if provided must be >= 16 years
- **Q**: Can trusted phones be empty? **A**: Yes

## Story Points
**5 points** – Reason: Medium complexity. Requires PATCH endpoint, partial update logic, validation, trusted contacts management, atomic DB updates, and comprehensive error handling.

## Implementation Plan
1. Update passenger model to include all editable fields
2. Add trusted_phones as JSON array in DB
3. Create migration for schema changes
4. Implement PATCH `/api/v1/passenger/profile` handler
5. Create DTOs for request/response
6. Implement service logic for partial updates and validation
7. Add trusted contacts add/delete logic
8. Ensure atomic DB updates (transaction)
9. Implement input validation utilities
10. Write unit tests for service and validation
11. Write integration tests for PATCH endpoint
12. Add error response standardization
13. Update API docs with Swag annotations
14. Manual QA with mobile app
15. Code review and merge

## TODO Checklist
- [ ] Update passenger model with editable fields
- [ ] Add trusted_phones JSON array to DB
- [ ] Create migration for schema changes
- [ ] Implement PATCH handler for profile editing
- [ ] Create request/response DTOs
- [ ] Implement service logic for partial updates
- [ ] Add trusted contacts management logic
- [ ] Ensure atomic DB updates
- [ ] Implement input validation utilities
- [ ] Write unit tests for service and validation
- [ ] Write integration tests for PATCH endpoint
- [ ] Standardize error responses
- [ ] Update API docs with Swag
- [ ] Manual QA with mobile app
- [ ] Code review
- [ ] Merge to main branch

## Testing Strategy
### Unit Tests
- Service: valid/invalid updates, partial updates, trusted contacts logic
- Validation: email, phone, DOB, gender, trusted phones
### Integration Tests
- PATCH endpoint: full/partial updates, error cases, concurrent updates
### Manual QA
- Edit profile in app, verify backend sync, error handling, placeholders

## Performance Considerations
- DB indexed on phone/email for fast lookup
- Atomic updates to prevent race conditions
- Response time target: <500ms

## Definition of Done
- [x] All acceptance criteria met
- [x] All TODOs completed
- [x] Unit/integration tests pass
- [x] Manual QA successful
- [x] API docs updated
- [x] Code reviewed and merged
- [x] Ready for mobile app integration
