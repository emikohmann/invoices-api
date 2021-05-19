package app

import (
	controller "github.com/emikohmann/invoices-api/cmd/api/controllers/invoices"
	service "github.com/emikohmann/invoices-api/cmd/api/services/invoices"
	"github.com/emikohmann/invoices-api/pkg/logger"
	"github.com/emikohmann/invoices-api/pkg/rest"
	"github.com/gin-gonic/gin"
)

func StartApp() {
	// Instance router and dependencies
	router := gin.Default()
	restClient := rest.NewRestClient()
	invoicesService := service.NewInvoicesService(restClient)
	invoicesController := controller.NewInvoicesController(invoicesService)

	// URL mapping
	UrlMapping(router, invoicesController)

	// Run
	Run(router)
}

func UrlMapping(router *gin.Engine, invoices controller.Controller) {
	router.POST("/invoices/process", invoices.Process)
}

func Run(router *gin.Engine) {
	if err := router.Run(); err != nil {
		logger.Panic("Error running server, %s", err.Error())
	}
}
