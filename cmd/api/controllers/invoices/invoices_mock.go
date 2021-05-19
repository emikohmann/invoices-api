package invoices

import (
	"github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type ControllerMock struct{}

func NewControllerMock() ControllerMock {
	return ControllerMock{}
}

func (controller ControllerMock) Process(c *gin.Context) {
	c.JSON(http.StatusOK, invoices.Response{
		User: invoices.User{
			Address: "118 Mariam Locks",
			Name:    "Danielle Skiles",
		},
		Movements: []invoices.Call{
			{
				Source:             "+5491167950940",
				Destination:        "+5491167950941",
				ReverseCharge:      false,
				Duration:           10 * time.Second,
				ComputableDuration: 10 * time.Second,
				DurationSeconds:    10,
				Date:               time.Now().Add(24 * 10 * time.Hour),
				Scope:              invoices.CallTypeNational,
				Friends:            false,
				Price:              2.5,
			},
		},
		TotalAmount:               2.5,
		TotalMinutesNational:      10,
		TotalMinutesInternational: 0,
		TotalMinutesFriends:       0,
	})
}
