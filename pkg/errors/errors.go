package errors

import "errors"

// AppError represents an application error with HTTP status and user-facing message.
type AppError struct {
	Code       string // Machine-readable code, e.g. "NOT_FOUND"
	Message    string // User-facing message
	StatusCode int    // HTTP status code
}

func (e *AppError) Error() string {
	return e.Message
}

// Ensure *AppError works with errors.As (wrap with fmt.Errorf so As can unwrap).
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewNotFound returns a 404 Not Found error with an optional custom message.
func NewNotFound(message string) *AppError {
	if message == "" {
		message = "The requested resource was not found."
	}
	return &AppError{Code: "NOT_FOUND", Message: message, StatusCode: 404}
}

// NewBadRequest returns a 400 Bad Request error with an optional custom message.
func NewBadRequest(message string) *AppError {
	if message == "" {
		message = "The request was invalid or malformed."
	}
	return &AppError{Code: "BAD_REQUEST", Message: message, StatusCode: 400}
}

// NewConflict returns a 409 Conflict error with an optional custom message.
func NewConflict(message string) *AppError {
	if message == "" {
		message = "A resource with this value already exists."
	}
	return &AppError{Code: "CONFLICT", Message: message, StatusCode: 409}
}

// NewValidation returns a 400 error for validation failures (e.g. field errors).
func NewValidation(message string) *AppError {
	if message == "" {
		message = "One or more fields failed validation."
	}
	return &AppError{Code: "VALIDATION_ERROR", Message: message, StatusCode: 400}
}

// NewUnauthorized returns a 401 Unauthorized error (e.g. invalid or missing token).
func NewUnauthorized(message string) *AppError {
	if message == "" {
		message = "Authentication required."
	}
	return &AppError{Code: "UNAUTHORIZED", Message: message, StatusCode: 401}
}

// NewForbidden returns a 403 Forbidden error (e.g. not allowed to access resource).
func NewForbidden(message string) *AppError {
	if message == "" {
		message = "You do not have permission to access this resource."
	}
	return &AppError{Code: "FORBIDDEN", Message: message, StatusCode: 403}
}

// AsAppError returns *AppError if err is or wraps an AppError; otherwise nil.
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
