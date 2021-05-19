package invoices

import (
	"github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/emikohmann/invoices-api/pkg/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var (
	testService = NewInvoicesService(
		rest.NewRestClientMock(),
	)
)

func TestServiceImpl_GetExtraDataSuccess(t *testing.T) {
	extraData, apiErr := testService.GetExtraData("+5491167950940")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.EqualValues(t, 9, len(extraData.Friends))
	assert.EqualValues(t, "118 Mariam Locks", extraData.Address)
	assert.EqualValues(t, "Danielle Skiles", extraData.Name)
	assert.EqualValues(t, "+5491167950940", extraData.PhoneNumber)
}

func TestServiceImpl_GetExtraDataError(t *testing.T) {
	_, apiErr := testService.GetExtraData("+0000000000000")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusInternalServerError, apiErr.Status())
	assert.EqualValues(t, "Mocked error", apiErr.Error())
}

func TestServiceImpl_ParseRequest(t *testing.T) {
	request, apiErr := testService.ParseRequest(
		"+5491167950940",
		"2020-01-01/2020-03-30",
		[]string{
			"numero origen,numero destino,cobro revertido,duracion,fecha",
			"+191167970944,+191167970944,S,99,2020-11-04T08:30:53Z",
			"+191167980953,+5491167930920,N,152,2020-08-27T05:55:43Z",
			"+5491167940999,+5491167950940,S,33,2021-03-24T18:14:04Z",
			"+5491167930920,+5491167980954,N,102,2020-04-11T17:15:07Z",
			"+5491167940999,+5491167980951,N,62,2020-10-31T09:44:55Z",
			"+191167980953,+191167980952,N,33,2020-10-12T02:35:13Z",
			"+5491167950940,+5491167930920,N,164,2021-02-01T02:32:17Z",
			"+191167970944,+5491167930920,N,193,2020-12-04T13:16:47Z",
			"+191167970944,+5491167980951,S,98,2021-01-11T04:03:50Z",
			"+191167970944,+5491167930920,S,64,2020-06-14T00:31:48Z",
			"+5491167950940,+191167970944,N,386,2021-03-09T13:42:18Z",
			"+191167970944,+5491167910920,S,104,2020-12-28T07:45:51Z",
			"+5491167980951,+5491167920930,N,53,2021-04-06T14:02:00Z",
			"+191167980953,+5491167910920,N,133,2020-08-17T06:36:44Z",
			"+5491167980951,+5491167950940,S,18,2020-02-02T01:03:22Z",
		}, invoices.ExtraData{
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
			PhoneNumber: "+5491167950940",
		})
	if apiErr != nil {
		t.Error(apiErr)
		return
	}

	assert.EqualValues(t, 2020, request.Period.From.Year())
	assert.EqualValues(t, time.January, request.Period.From.Month())
	assert.EqualValues(t, 1, request.Period.From.Day())

	assert.EqualValues(t, 2020, request.Period.To.Year())
	assert.EqualValues(t, time.March, request.Period.To.Month())
	assert.EqualValues(t, 30, request.Period.To.Day())

	assert.EqualValues(t, "+5491167950940", request.Phone)
	assert.EqualValues(t, 1, len(request.Calls))
}

func TestServiceImpl_GenerateInvoice(t *testing.T) {
	request := invoices.Request{
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
	}

	testService.ComputePrices(&request)
	response, apiErr := testService.GenerateInvoice(request, invoices.ExtraData{
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
		PhoneNumber: "+5491167950940",
	})
	if apiErr != nil {
		t.Error(apiErr)
		return
	}

	assert.EqualValues(t, "118 Mariam Locks", response.User.Address)
	assert.EqualValues(t, "Danielle Skiles", response.User.Name)
	assert.EqualValues(t, 1, len(response.Movements))
	assert.False(t, response.Movements[0].Friends)
	assert.False(t, response.Movements[0].ReverseCharge)
	assert.EqualValues(t, 10*time.Second, response.Movements[0].Duration)
	assert.EqualValues(t, invoices.Phone("+5491167950940"), response.Movements[0].Source)
	assert.EqualValues(t, invoices.Phone("+5491167950941"), response.Movements[0].Destination)
	assert.EqualValues(t, time.Now().Add(24*10*time.Hour).Year(), response.Movements[0].Date.Year())
	assert.EqualValues(t, invoices.CallTypeNational, response.Movements[0].Scope)
	assert.EqualValues(t, 10*time.Second, response.Movements[0].ComputableDuration)
	assert.EqualValues(t, 2.5, response.Movements[0].Price)
}
