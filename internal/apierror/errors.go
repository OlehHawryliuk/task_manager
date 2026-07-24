package apierror

import (
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Details    string
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	ErrInvalidRequest = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_REQUEST",
		Message:    "Invalid request",
	}

	ErrUnauthorized = &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
	}

	ErrForbidden = &APIError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    "Access denied",
	}

	ErrNotFound = &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
	}

	ErrConflict = &APIError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    "Resource already exists",
	}

	ErrInternalServer = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "Internal server error",
	}

	ErrDatabaseError = &APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "DATABASE_ERROR",
		Message:    "Database operation failed",
	}

	ErrValidationError = &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
	}
)

func NewInvalidRequest(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_REQUEST",
		Message:    "Invalid request",
		Details:    details,
	}
}

func NewUnauthorized(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
		Details:    details,
	}
}

func NewForbidden(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusForbidden,
		Code:       "FORBIDDEN",
		Message:    "Access denied",
		Details:    details,
	}
}

func NewNotFound(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		Details:    details,
	}
}

func NewConflict(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    "Resource already exists",
		Details:    details,
	}
}

func NewValidationError(details string) *APIError {
	return &APIError{
		StatusCode: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
		Details:    details,
	}
}
