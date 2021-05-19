package rest

import (
	"encoding/json"
	"fmt"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"io/ioutil"
	"net/http"
)

type ClientImpl struct{}

func NewRestClient() ClientImpl {
	return ClientImpl{}
}

func (client ClientImpl) Get(url string, fill interface{}) apierrors.ApiError {
	response, err := http.Get(url)
	if err != nil {
		return apierrors.NewInternalServer(fmt.Sprintf("Error doing API call: %s", err.Error()))
	}
	if response == nil {
		return apierrors.NewInternalServer("Error doing API call: nil response")
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return apierrors.NewInternalServer(fmt.Sprintf("Error reading response data: %s", err.Error()))
	}
	if response.StatusCode != http.StatusOK {
		return apierrors.NewApiError("invalid_status", response.StatusCode, string(data))
	}
	if err := json.Unmarshal(data, &fill); err != nil {
		return apierrors.NewInternalServer(fmt.Sprintf("Error unmarshaling response: %s", err.Error()))
	}
	return nil
}
