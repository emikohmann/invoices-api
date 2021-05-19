package apierrors

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestApiError(t *testing.T) {
	apiErr := apiError{
		ErrorCode:    "test_code",
		ErrorStatus:  10,
		ErrorMessage: "Test message",
	}
	assert.EqualValues(t, "test_code", apiErr.Code())
	assert.EqualValues(t, 10, apiErr.Status())
	assert.EqualValues(t, "Test message", apiErr.Message())
	assert.EqualValues(t, "Test message", apiErr.Error())
}

func TestNewBadRequest(t *testing.T) {
	apiErr := NewBadRequest("Test message")
	assert.EqualValues(t, "bad_request", apiErr.Code())
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
	assert.EqualValues(t, "Test message", apiErr.Message())
}

func TestNewUnauthorized(t *testing.T) {
	apiErr := NewUnauthorized("user 1", "invoice 100")
	assert.EqualValues(t, "unauthorized", apiErr.Code())
	assert.EqualValues(t, http.StatusUnauthorized, apiErr.Status())
	assert.EqualValues(t, "Caller user 1 is not authorized to use the resource invoice 100", apiErr.Message())
}

func TestNewNotFound(t *testing.T) {
	apiErr := NewNotFound("invoice 100")
	assert.EqualValues(t, "not_found", apiErr.Code())
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status())
	assert.EqualValues(t, "Resource invoice 100 not found", apiErr.Message())
}

func TestNewConflict(t *testing.T) {
	apiErr := NewConflict("Test message")
	assert.EqualValues(t, "conflict", apiErr.Code())
	assert.EqualValues(t, http.StatusConflict, apiErr.Status())
	assert.EqualValues(t, "Test message", apiErr.Message())
}

func TestNewInternalServer(t *testing.T) {
	apiErr := NewInternalServer("Test message")
	assert.EqualValues(t, "internal_server_error", apiErr.Code())
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "Test message", apiErr.Message())
}
