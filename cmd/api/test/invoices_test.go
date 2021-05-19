package test

import (
	"encoding/json"
	"github.com/emikohmann/invoices-api/cmd/api/app"
	controller "github.com/emikohmann/invoices-api/cmd/api/controllers/invoices"
	domain "github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProcessInvoice(t *testing.T) {
	response := httptest.NewRecorder()
	_, router := gin.CreateTestContext(response)
	invoices := controller.NewControllerMock()
	app.UrlMapping(router, invoices)

	var reader io.Reader = nil
	request, err := http.NewRequest(http.MethodPost, "/invoices/process", reader)
	if err != nil {
		panic(err)
	}

	router.ServeHTTP(response, request)
	assert.EqualValues(t, http.StatusOK, response.Code)

	var data domain.Response
	if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
		t.Error(err)
		return
	}

	assert.EqualValues(t, "118 Mariam Locks", data.User.Address)
	assert.EqualValues(t, "Danielle Skiles", data.User.Name)
	assert.EqualValues(t, 1, len(data.Movements))
	assert.False(t, data.Movements[0].Friends)
	assert.False(t, data.Movements[0].ReverseCharge)
	assert.EqualValues(t, 10, data.Movements[0].DurationSeconds)
	assert.EqualValues(t, domain.Phone("+5491167950940"), data.Movements[0].Source)
	assert.EqualValues(t, domain.Phone("+5491167950941"), data.Movements[0].Destination)
	assert.EqualValues(t, time.Now().Add(24*10*time.Hour).Year(), data.Movements[0].Date.Year())
	assert.EqualValues(t, domain.CallTypeNational, data.Movements[0].Scope)
	assert.EqualValues(t, 2.5, data.Movements[0].Price)
}
