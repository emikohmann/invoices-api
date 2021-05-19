package invoices

import (
	"encoding/json"
	"github.com/emikohmann/invoices-api/cmd/api/services/invoices"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestControllerImpl_Process(t *testing.T) {
	service := invoices.NewServiceMock()
	controller := NewInvoicesController(service)

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	var err error
	context.Request, err = http.NewRequest(http.MethodGet, "/invoices/process", strings.NewReader(""))
	if err != nil {
		t.Error(err)
		return
	}

	controller.Process(context)
	assert.EqualValues(t, http.StatusBadRequest, response.Code)

	apiErr := make(map[string]interface{})
	if err := json.Unmarshal(response.Body.Bytes(), &apiErr); err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, "bad_request", apiErr["code"])
	assert.EqualValues(t, http.StatusBadRequest, apiErr["status_code"])
	assert.EqualValues(t, "clientPhone is required", apiErr["message"])
}
