---
applyTo: ["*"]
---

# Response, Error & Pagination Utilities

## Quick Response Functions

```go
// Success responses
response.Success(c, data, "message")
response.Created(c, data, "message")
response.NoContent(c)

// Error responses
response.BadRequest(c, "message", details)
response.Unauthorized(c, "message")
response.NotFound(c, "message")
response.InternalServerError(c, "message")
```

## Response Handler

```go
handler := response.NewHandler()
handler.HandleError(c, err)           // Auto-detects error type
handler.HandleSuccess(c, data, "msg")
handler.HandlePaginated(c, data, pagination, "msg")
```

## Error Types

```go
// AppError with HTTP status
err := errors.NewBadRequest("message")
err := errors.NewNotFound("message")
err := errors.NewConflict("message")

// Validation errors
valErr := errors.NewValidationErrors()
valErr.Add("field", "error message", "value")
appErr := valErr.ToAppError()

// Database/External service errors
dbErr := errors.NewDatabaseError("SELECT", "users", originalErr)
extErr := errors.NewExternalServiceError("service", "code", "msg", originalErr)
```

## Pagination

```go
// Offset-based pagination
params := pagination.FromQuery(c)  // ?page=2&per_page=20
paginationMeta := pagination.NewPagination(params, total)
response.Paginated(c, data, paginationMeta, "message")

// Cursor-based pagination
cursor := pagination.FromCursorQuery(c)  // ?after=cursor&limit=20
cursorResult := pagination.NewCursorResult(data, "next", "prev", hasNext, hasPrev)
```

## Fluent API

```go
// Success response
response.NewResponseBuilder(c).
    WithData(data).
    WithMessage("success").
    WithVersion("v1").
    SendOK()

// Error response
response.NewErrorBuilder(c).
    WithCode("ERROR_CODE").
    WithMessage("error message").
    SendBadRequest()
```

## Response Structure

```go
type Response struct {
    Success   bool        `json:"success"`
    Message   string      `json:"message,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Error     *APIError   `json:"error,omitempty"`
    Meta      *Meta       `json:"meta,omitempty"`
    Timestamp time.Time   `json:"timestamp"`
    RequestID string      `json:"request_id,omitempty"`
}
```
