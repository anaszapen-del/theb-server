# Developer Task Executor Prompt

## Purpose

This prompt is for an automated developer agent. The agent picks up a technical task file from `/tasks`, creates a Git branch, does the work (code, tests, docs), checks that all Acceptance Criteria are met, then pushes the branch and opens a Pull Request. After completion, the agent generates a frontend task file with API endpoints documentation for frontend implementation.

---

## High-Level Steps

1. **Find and Read Task:** Get the task file (from user or most recent). Read its structure.
2. **Get Task Info:** Pull out the task ID, title, acceptance criteria, requirements, plan, checklist, and assumptions.
3. **Create Git Branch:**

- Name: `feature/<TASK_ID>-<short-title>` (lowercase, use hyphens, max 60 chars after ID).
- Create from `main` (or base branch).
- Document branch name in task notes.

4. **Plan Work:**

- Make a checklist from acceptance criteria and TODOs.
- Expand each acceptance criterion into clear steps.
- Note any unclear points as assumptions.

5. **Create Test Cases Document:**

- Create or update a test cases document (one test case for each acceptance criterion: Preconditions, Steps, Expected Result).
- Save in task folder or `/tests/test-cases/`.

6. **Check Related Code:**

- Search for related code to reuse or update.

7. **Do the Work (Repeat until done):**

- Update or create types and data validation.
- Update or add service functions (API calls).
- Implement repository layer (database operations).
- Implement handlers (HTTP endpoints).
- Add or update middleware if needed.
- Run build, lint, and tests. Fix any issues.
- Check performance.

8. **API Documentation:**

- Document all created/updated API endpoints with Swagger annotations.
- Include request/response examples.
- Document authentication requirements.

9. **Manual Testing:**

- Test API endpoints manually (Postman, curl, or similar).
- Verify error handling and edge cases.
- Make sure all acceptance criteria are met.

10. **Docs & Changelog:**

- Update API documentation if needed.
- Add module-specific documentation.

11. **Push & PR:**

- Commit with clear messages following conventional commits.
- Push branch to remote.
- Open a Pull Request with clear title and description.

12. **Generate Frontend Task:**

- Create a frontend task file in `/tasks/frontend/` with:
  - Feature overview and requirements
  - All API endpoints with complete specifications
  - Request/response schemas
  - Authentication requirements
  - Error handling details
  - UI/UX requirements
  - State management needs
  - Testing requirements

13. **Completion Summary:**

- List actions taken, acceptance criteria coverage, files changed, quality checks, performance notes, and next steps.
- Confirm frontend task file created.

---

## Task File Structure

Expect task files to have sections like:

- Title (with Task ID)
- Overview
- Scope (In/Out)
- Requirements
- Data & Contracts
- API Endpoints (for backend tasks)
- Database Schema
- Business Logic
- Acceptance Criteria (checkboxes)
- Edge Cases
- Risks & Mitigations
- Open Questions / Assumptions
- Story Points
- Implementation Plan
- TODO Checklist
- Testing Strategy
- Performance Considerations
- Definition of Done

---

## Frontend Task Generation

After completing the backend task, generate a comprehensive frontend task file with:

### File Structure
Save as: `/tasks/frontend/FE-<TASK_ID>-<feature-name>.md`

### Required Sections

1. **Title & Overview**
   - Feature name and description
   - Link to backend task
   - User stories

2. **API Integration**
   - Complete list of API endpoints with:
     - HTTP method and path
     - Request headers (authentication, content-type)
     - Request body schema with types
     - Response body schema with types
     - Status codes and error responses
     - Example requests and responses

3. **Authentication & Authorization**
   - JWT token handling
   - Token refresh flow
   - Role-based access requirements

4. **Data Models**
   - TypeScript interfaces for all DTOs
   - Validation schemas
   - Default values

5. **UI Components**
   - Screen/component structure
   - Props and state requirements
   - User interactions

6. **State Management**
   - What data needs to be stored
   - Where (local state, context, Redux, etc.)
   - Data flow

7. **Error Handling**
   - API error scenarios
   - User-facing error messages
   - Retry logic

8. **Loading & Empty States**
   - Loading indicators
   - Empty state designs
   - Skeleton screens

9. **Testing Requirements**
   - Unit tests for services/hooks
   - Component tests
   - Integration tests for API calls

10. **Acceptance Criteria**
    - Functional requirements
    - UI/UX requirements
    - Performance requirements

11. **Design Assets**
    - Links to designs or screenshots
    - Color codes, spacing, typography

---

## Internal Checklists

### Start

- [ ] Task file loaded
- [ ] Task requirements understood
- [ ] Assumptions/questions noted
- [ ] Plan refined
- [ ] Git branch created

### Implementation

- [ ] Types and validation updated
- [ ] Repository layer implemented (database operations)
- [ ] Service layer implemented (business logic)
- [ ] Handlers implemented (HTTP endpoints)
- [ ] Middleware applied (auth, validation, etc.)
- [ ] DTOs created with proper validation tags
- [ ] Swagger annotations added
- [ ] Each code path for the task is implemented
- [ ] Error handling implemented
- [ ] Performance checked
- [ ] Docs updated (if needed)

### Testing

- [ ] Unit tests written for services
- [ ] Repository tests written
- [ ] Handler tests written
- [ ] All tests passing
- [ ] Manual API testing completed
- [ ] Edge cases covered

### Completion

- [ ] All acceptance criteria met
- [ ] No errors/warnings
- [ ] Lint & typecheck clean
- [ ] Branch pushed & PR created
- [ ] Frontend task file generated in `/tasks/frontend/`
- [ ] API endpoints fully documented
- [ ] Summary output generated

---

## Frontend Task Template

When generating the frontend task file, use this template:

```markdown
# [FE-<TASK_ID>] <Feature Name>

**Backend Task:** [Link to backend task or PR]
**Created:** <Date>
**Status:** Todo

---

## Overview

Brief description of what this feature does from the user's perspective.

### User Stories

- As a [role], I want to [action] so that [benefit]

---

## API Endpoints

### 1. [Endpoint Name]

**Endpoint:** `<METHOD> /api/v1/path`

**Authentication:** Required (Bearer Token) / Not Required

**Request Headers:**
```json
{
  "Authorization": "Bearer <token>",
  "Content-Type": "application/json"
}
```

**Request Body:**
```typescript
interface RequestBody {
  field1: string;
  field2: number;
  // Include all fields with types
}
```

**Success Response (200):**
```typescript
interface SuccessResponse {
  data: {
    // Response structure
  };
  message: string;
}
```

**Error Responses:**
- `400 Bad Request` - Validation error
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

**Example Request:**
```javascript
const response = await fetch('/api/v1/path', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    field1: 'value',
    field2: 123
  })
});
```

**Example Response:**
```json
{
  "data": { ... },
  "message": "Success"
}
```

---

## Data Models

```typescript
// TypeScript interfaces for all DTOs
interface UserProfile {
  id: string;
  name: string;
  // All fields with types
}

// Validation schemas (if using Zod, Yup, etc.)
const userProfileSchema = z.object({
  // Validation rules
});
```

---

## UI Components

### Screens/Components to Build

1. **ComponentName**
   - Purpose: ...
   - Location: `src/screens/` or `src/components/`
   - Props: ...
   - State: ...

---

## State Management

- **API Data:** Store user profile in Context/Redux
- **Form State:** Local state with React Hook Form
- **Loading/Error States:** Local state

---

## Error Handling

| Error Code | Scenario | User Message | Action |
|------------|----------|--------------|--------|
| 400 | Invalid input | "Please check your input" | Show field errors |
| 401 | Token expired | "Please log in again" | Redirect to login |
| 404 | Not found | "Resource not found" | Show empty state |

---

## Loading & Empty States

- **Loading:** Show skeleton screen while fetching data
- **Empty:** Display message when no data available
- **Error:** Show error message with retry button

---

## Acceptance Criteria

- [ ] User can [action]
- [ ] API integration working correctly
- [ ] Error handling implemented
- [ ] Loading states shown
- [ ] UI matches design
- [ ] Responsive on all screen sizes
- [ ] Accessible (proper labels, keyboard navigation)
- [ ] Tests written and passing

---

## Testing Requirements

- [ ] Unit tests for API service functions
- [ ] Unit tests for custom hooks
- [ ] Component tests for UI logic
- [ ] Integration tests for API calls
- [ ] Manual testing checklist completed

---

## Design Notes

- Colors: [List colors used]
- Typography: [Font sizes, weights]
- Spacing: [Padding, margins]
- Design link: [Figma/Design file link]

---

## Implementation Notes

- Assumptions made
- Edge cases to consider
- Performance considerations
- Accessibility requirements

---

## Definition of Done

- [ ] Code implemented and reviewed
- [ ] Tests written and passing (>80% coverage)
- [ ] API integration tested
- [ ] Error handling verified
- [ ] UI tested on iOS and Android
- [ ] Performance acceptable
- [ ] Accessibility verified
- [ ] Documentation updated
```

---

## Example Frontend Task Output

After completing a user authentication backend task, generate:

**File:** `/tasks/frontend/FE-AUTH-001-phone-login.md`

Include:
- POST /api/v1/auth/login endpoint details
- POST /api/v1/auth/verify endpoint details
- Request/response TypeScript interfaces
- JWT token handling instructions
- Login screen requirements
- OTP verification screen requirements
- Error scenarios and messages
- Testing checklist
