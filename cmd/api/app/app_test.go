package app

import (
	"github.com/emikohmann/invoices-api/cmd/api/controllers/invoices"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestUrlMapping(t *testing.T) {
	controller := invoices.NewControllerMock()
	response := httptest.NewRecorder()
	_, router := gin.CreateTestContext(response)
	UrlMapping(router, controller)
}
