package invoices

import (
	"fmt"
	"github.com/emikohmann/invoices-api/cmd/api/services/invoices"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"github.com/emikohmann/invoices-api/pkg/logger"
	requestUtils "github.com/emikohmann/invoices-api/pkg/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	formPhone   = "clientPhone"
	formPeriod  = "invoicePeriod"
	formInput   = "input.csv"
	formMaxSize = 1 * 1024 * 1024
)

type ControllerImpl struct {
	Service invoices.Service
}

func NewInvoicesController(service invoices.Service) ControllerImpl {
	return ControllerImpl{
		Service: service,
	}
}

func (controller ControllerImpl) Process(c *gin.Context) {
	// Validate phone
	phone := strings.TrimSpace(c.Request.FormValue(formPhone))
	if phone == "" {
		apiErr := apierrors.NewBadRequest(fmt.Sprintf("%s is required", formPhone))
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Validate period
	period := strings.TrimSpace(c.Request.FormValue(formPeriod))
	if period == "" {
		apiErr := apierrors.NewBadRequest(fmt.Sprintf("%s is required", formPeriod))
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Validate file
	input, apiErr := requestUtils.GetFile(c.Request, formInput, formMaxSize)
	if apiErr != nil {
		logger.Error("Error getting input file: %s", apiErr.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}
	if len(input) == 0 {
		apiErr := apierrors.NewBadRequest(fmt.Sprintf("The input file %s is empty", formInput))
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Get friend list
	extraData, apiErr := controller.Service.GetExtraData(phone)
	if apiErr != nil {
		logger.Error("Error getting friend list: %s", apiErr.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Parse entire request and calls
	request, apiErr := controller.Service.ParseRequest(phone, period, input, extraData)
	if apiErr != nil {
		logger.Error("Error parsing request: %s", apiErr.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	// Compute prices
	controller.Service.ComputePrices(&request)

	// Generate final invoice
	response, apiErr := controller.Service.GenerateInvoice(request, extraData)
	if apiErr != nil {
		logger.Error("Error generating invoice: %s", apiErr.Error())
		c.JSON(apiErr.Status(), apiErr)
		return
	}

	c.JSON(http.StatusOK, response)
}
