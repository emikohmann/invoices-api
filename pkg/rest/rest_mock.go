package rest

import (
	"encoding/json"
	"fmt"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"strings"
)

type ClientMock struct{}

func NewRestClientMock() ClientMock {
	return ClientMock{}
}

const (
	mockedResponse = `{
	    "address": "118 Mariam Locks",
	    "friends": [
	        "+191167970944",
	        "+5491167940999",
	        "+191167980952",
	        "+5491167930920",
	        "+5491167920944",
	        "+5491167980954",
	        "+5491167980953",
	        "+5491167980951",
	        "+191167980953"
	    ],
	    "name": "Danielle Skiles",
	    "phone_number": "+5491167950940"
	}`
)

func (ClientMock) Get(url string, fill interface{}) apierrors.ApiError {
	if !strings.Contains(url, "+5491167950940") {
		return apierrors.NewInternalServer("Mocked error")
	}
	if err := json.Unmarshal([]byte(mockedResponse), &fill); err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("Error unmarshaling mocked response: %s", err.Error()))
	}
	return nil
}
