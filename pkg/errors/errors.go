package errors

import "fmt"

// AppError represents an application error with HTTP status code
type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Err        error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors
func BadRequest(message string, err error) *AppError {
	return NewAppError("BAD_REQUEST", message, 400, err)
}

func Unauthorized(message string, err error) *AppError {
	return NewAppError("UNAUTHORIZED", message, 401, err)
}

func Forbidden(message string, err error) *AppError {
	return NewAppError("FORBIDDEN", message, 403, err)
}

func NotFound(message string, err error) *AppError {
	return NewAppError("NOT_FOUND", message, 404, err)
}

func Conflict(message string, err error) *AppError {
	return NewAppError("CONFLICT", message, 409, err)
}

func TooManyRequests(message string, err error) *AppError {
	return NewAppError("RATE_LIMIT_EXCEEDED", message, 429, err)
}

func InternalServerError(message string, err error) *AppError {
	return NewAppError("INTERNAL_SERVER_ERROR", message, 500, err)
}

func ServiceUnavailable(message string, err error) *AppError {
	return NewAppError("SERVICE_UNAVAILABLE", message, 503, err)
}
