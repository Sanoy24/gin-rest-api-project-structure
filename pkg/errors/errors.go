package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewAppError(code int, message, errorType string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}

// common error types

var (
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid credentials", "INVALID_CREDENTIALS")
	ErrUnAuthorized       = NewAppError(http.StatusUnauthorized, "Unauthorized access", "UNAUTHORIZED")
	ErrUserNotFound       = NewAppError(http.StatusNotFound, "User not found", "USER_NOT_FOUND")
	ErrUserExists         = NewAppError(http.StatusConflict, "User already exists", "USER_EXISTS")
	ErrInvalidInput       = NewAppError(http.StatusBadRequest, "Invalid input", "INVALID_INPUT")
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "Internal server error", "INTERNAL")
)
