package apierrors

import (
	"fmt"
	"net/http"
)

type ApiError interface {
	Code() string
	Status() int
	Message() string
	Error() string
}

type apiError struct {
	ErrorCode    string `json:"code"`
	ErrorStatus  int    `json:"status_code"`
	ErrorMessage string `json:"message"`
}

func (error apiError) Code() string {
	return error.ErrorCode
}

func (error apiError) Status() int {
	return error.ErrorStatus
}

func (error apiError) Message() string {
	return error.ErrorMessage
}

func (error apiError) Error() string {
	return error.ErrorMessage
}

func NewBadRequest(message string) ApiError {
	return NewApiError("bad_request", http.StatusBadRequest, message)
}

func NewUnauthorized(caller string, resource string) ApiError {
	return NewApiError("unauthorized", http.StatusUnauthorized, fmt.Sprintf("Caller %s is not authorized to use the resource %s", caller, resource))
}

func NewNotFound(resource string) ApiError {
	return NewApiError("not_found", http.StatusNotFound, fmt.Sprintf("Resource %s not found", resource))
}

func NewConflict(message string) ApiError {
	return NewApiError("conflict", http.StatusConflict, message)
}

func NewInternalServer(message string) ApiError {
	return NewApiError("internal_server_error", http.StatusInternalServerError, message)
}

func NewApiError(code string, status int, message string) ApiError {
	return apiError{
		ErrorCode:    code,
		ErrorStatus:  status,
		ErrorMessage: message,
	}
}
