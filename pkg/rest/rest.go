package rest

import "github.com/emikohmann/invoices-api/pkg/apierrors"

type Client interface {
	Get(url string, fill interface{}) apierrors.ApiError
}

