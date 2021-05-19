package invoices

import (
	"github.com/emikohmann/invoices-api/cmd/api/config"
	"github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"time"
)

type ServiceMock struct{}

func NewServiceMock() ServiceMock {
	return ServiceMock{}
}

func (service ServiceMock) GetExtraData(phone string) (invoices.ExtraData, apierrors.ApiError) {
	return invoices.ExtraData{
		Address: "118 Mariam Locks",
		Friends: []string{
			"+191167970944",
			"+5491167940999",
			"+191167980952",
			"+5491167930920",
			"+5491167920944",
			"+5491167980954",
			"+5491167980953",
			"+5491167980951",
			"+191167980953",
		},
		Name:        "Danielle Skiles",
		PhoneNumber: invoices.Phone(phone),
	}, nil
}

func (service ServiceMock) ParseRequest(phone string, period string, input []string, extraData invoices.ExtraData) (invoices.Request, apierrors.ApiError) {
	return invoices.Request{
		Phone: "+5491167950940",
		Period: invoices.Period{
			From: time.Now(),
			To:   time.Now().Add(24 * 30 * time.Hour),
		},
		Calls: []invoices.Call{
			{
				Source:          "+5491167950940",
				Destination:     "+5491167950941",
				ReverseCharge:   false,
				Duration:        10 * time.Second,
				DurationSeconds: 10,
				Date:            time.Now().Add(24 * 10 * time.Hour),
				Scope:           invoices.CallTypeNational,
				Friends:         false,
			},
		},
		RemainingFriendMinutes: 0,
	}, nil
}

func (service ServiceMock) ComputePrices(request *invoices.Request) {
	for i := range request.Calls {
		request.Calls[i].Price = config.PriceNationalPerCall
	}
}

func (service ServiceMock) GenerateInvoice(request invoices.Request, extraData invoices.ExtraData) (invoices.Response, apierrors.ApiError) {
	return invoices.Response{
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
	}, nil
}
