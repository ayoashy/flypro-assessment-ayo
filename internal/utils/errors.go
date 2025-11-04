package utils

import (
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND"
	ErrorTypeConflict   ErrorType = "CONFLICT"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
	ErrorTypeBadRequest ErrorType = "BAD_REQUEST"
)

// AppError represents an application error
type AppError struct {
	Type     ErrorType `json:"type"`
	Message  string    `json:"message"`
	Field    string    `json:"field,omitempty"`
	Code     int       `json:"code"`
	Internal error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

// NewValidationError creates a validation error
func NewValidationError(field, message string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Field:   field,
		Code:    http.StatusBadRequest,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Code:    http.StatusNotFound,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    http.StatusConflict,
	}
}

// NewInternalError creates an internal error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:     ErrorTypeInternal,
		Message:  message,
		Code:     http.StatusInternalServerError,
		Internal: err,
	}
}

// NewBadRequestError creates a bad request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Code:    http.StatusBadRequest,
	}
}
