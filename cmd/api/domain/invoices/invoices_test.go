package invoices

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestPhone_GetCountryCode(t *testing.T) {
	for _, testCase := range []struct {
		InputNumber         string
		ExpectedCountryCode string
	}{
		{
			InputNumber:         "+513512813948",
			ExpectedCountryCode: "+51",
		},
		{
			InputNumber:         "+523512813948",
			ExpectedCountryCode: "+52",
		},
		{
			InputNumber:         "+533512813948",
			ExpectedCountryCode: "+53",
		},
		{
			InputNumber:         "+543512813948",
			ExpectedCountryCode: "+54",
		},
		{
			InputNumber:         "+553512813948",
			ExpectedCountryCode: "+55",
		},
		{
			InputNumber:         "+563512813948",
			ExpectedCountryCode: "+56",
		},
		{
			InputNumber:         "+573512813948",
			ExpectedCountryCode: "+57",
		},
		{
			InputNumber:         "+583512813948",
			ExpectedCountryCode: "+58",
		},
		{
			InputNumber:         "+593512813948",
			ExpectedCountryCode: "+59",
		},
	} {
		countryCode, err := Phone(testCase.InputNumber).GetCountryCode()
		if err != nil {
			t.Error(err)
			return
		}
		assert.EqualValues(t, testCase.ExpectedCountryCode, countryCode)
	}
}

func TestCall_FillSource(t *testing.T) {
	call := Call{}

	apiErr := call.FillSource("+5491167950941")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.EqualValues(t, Phone("+5491167950941"), call.Source)
}

func TestCall_FillSourceError(t *testing.T) {
	call := Call{}

	apiErr := call.FillSource("")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestCall_FillDestination(t *testing.T) {
	call := Call{
		Source: "+5491167950940",
	}

	apiErr := call.FillDestination("+5491167950941")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.EqualValues(t, Phone("+5491167950941"), call.Destination)
}

func TestCall_FillDestinationError(t *testing.T) {
	call := Call{
		Source: "+5491167950940",
	}

	apiErr := call.FillDestination("")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestCall_FillReverseChargeTrue(t *testing.T) {
	call := Call{
		Source:      "+5491167950940",
		Destination: "+5491167950941",
	}

	apiErr := call.FillReverseCharge("S")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.True(t, call.ReverseCharge)
}

func TestCall_FillReverseChargeFalse(t *testing.T) {
	call := Call{
		Source:      "+5491167950940",
		Destination: "+5491167950941",
	}

	apiErr := call.FillReverseCharge("N")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.False(t, call.ReverseCharge)
}

func TestCall_FillReverseChargeError(t *testing.T) {
	call := Call{
		Source:      "+5491167950940",
		Destination: "+5491167950941",
	}

	apiErr := call.FillReverseCharge("INVALID")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestCall_FillDuration(t *testing.T) {
	call := Call{
		Source:        "+5491167950940",
		Destination:   "+5491167950941",
		ReverseCharge: false,
	}

	apiErr := call.FillDuration("30")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.EqualValues(t, 30*time.Second, call.Duration)
	assert.EqualValues(t, 30, call.DurationSeconds)
}

func TestCall_FillDurationError(t *testing.T) {
	call := Call{
		Source:        "+5491167950940",
		Destination:   "+5491167950941",
		ReverseCharge: false,
	}

	apiErr := call.FillDuration("INVALID")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestCall_FillDate(t *testing.T) {
	call := Call{
		Source:          "+5491167950940",
		Destination:     "+5491167950941",
		ReverseCharge:   false,
		Duration:        10 * time.Second,
		DurationSeconds: 10,
	}

	apiErr := call.FillDate("2020-01-05T12:05:43.000Z")
	if apiErr != nil {
		t.Error(apiErr)
		return
	}
	assert.EqualValues(t, 2020, call.Date.Year())
	assert.EqualValues(t, time.January, call.Date.Month())
	assert.EqualValues(t, 5, call.Date.Day())
}

func TestCall_FillDateError(t *testing.T) {
	call := Call{
		Source:          "+5491167950940",
		Destination:     "+5491167950941",
		ReverseCharge:   false,
		Duration:        10 * time.Second,
		DurationSeconds: 10,
	}

	apiErr := call.FillDate("INVALID")
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestCall_FillCallType(t *testing.T) {
	call := Call{
		Source:          "+5491167950940",
		Destination:     "+5491167950941",
		ReverseCharge:   false,
		Duration:        10 * time.Second,
		DurationSeconds: 10,
		Date:            time.Now().Add(24 * 10 * time.Hour),
	}

	apiErr := call.FillCallType(ExtraData{
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
	assert.EqualValues(t, CallTypeNational, call.Scope)
	assert.False(t, call.Friends)
}

func TestCall_FillCallTypeError(t *testing.T) {
	call := Call{
		Source:          "",
		Destination:     "",
		ReverseCharge:   false,
		Duration:        10 * time.Second,
		DurationSeconds: 10,
		Date:            time.Now().Add(24 * 10 * time.Hour),
	}

	apiErr := call.FillCallType(ExtraData{
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
	assert.NotNil(t, apiErr)
	assert.EqualValues(t, http.StatusBadRequest, apiErr.Status())
}

func TestRequest_ComputePrices(t *testing.T) {
	request := Request{
		Phone: "+5491167950940",
		Period: Period{
			From: time.Now(),
			To:   time.Now().Add(24 * 30 * time.Hour),
		},
		Calls: []Call{
			{
				Source:          "+5491167950940",
				Destination:     "+5491167950941",
				ReverseCharge:   false,
				Duration:        10 * time.Second,
				DurationSeconds: 10,
				Date:            time.Now().Add(24 * 10 * time.Hour),
				Scope:           CallTypeNational,
				Friends:         false,
			},
		},
		RemainingFriendMinutes: 0,
	}
	request.ComputePrices()
	assert.EqualValues(t, 10*time.Second, request.Calls[0].ComputableDuration)
	assert.EqualValues(t, 2.5, request.Calls[0].Price)
}
