package invoices

import (
	"github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
)

type Service interface {
	GetExtraData(phone string) (invoices.ExtraData, apierrors.ApiError)
	ParseRequest(phone string, period string, input []string, extraData invoices.ExtraData) (invoices.Request, apierrors.ApiError)
	ComputePrices(request *invoices.Request)
	GenerateInvoice(request invoices.Request, extraData invoices.ExtraData) (invoices.Response, apierrors.ApiError)
}

