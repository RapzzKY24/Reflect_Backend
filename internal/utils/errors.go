package utils

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

func Unauthorized(message string) *AppError {
	return NewAppError(http.StatusUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}

func NotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

func Conflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}

func InternalServerError(message string) *AppError {
	return NewAppError(http.StatusInternalServerError, message)
}