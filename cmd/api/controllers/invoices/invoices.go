package invoices

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Process(c *gin.Context)
}
