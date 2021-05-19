package invoices

import (
	"fmt"
	"github.com/emikohmann/invoices-api/cmd/api/config"
	"github.com/emikohmann/invoices-api/cmd/api/domain/invoices"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"github.com/emikohmann/invoices-api/pkg/rest"
	"strings"
	"time"
)

const (
	friendListURL   = "https://interview-brubank-api.herokuapp.com/users/%s"
	periodSeparator = "/"
	periodElements  = 2
	periodFormat    = "2006-01-02"
)

type ServiceImpl struct {
	Client rest.Client
}

func NewInvoicesService(client rest.Client) ServiceImpl {
	return ServiceImpl{
		Client: client,
	}
}

func (service ServiceImpl) GetExtraData(phone string) (invoices.ExtraData, apierrors.ApiError) {
	var friendList invoices.ExtraData
	if apiErr := service.Client.Get(fmt.Sprintf(friendListURL, phone), &friendList); apiErr != nil {
		return friendList, apiErr
	}
	return friendList, nil
}

func (service ServiceImpl) ParseRequest(phone string, period string, input []string, extraData invoices.ExtraData) (invoices.Request, apierrors.ApiError) {
	var request invoices.Request

	// Parse periods
	from, to, apiErr := service.GetPeriods(period)
	if apiErr != nil {
		return request, apiErr
	}
	computed := invoices.Period{
		From: from,
		To:   to,
	}

	// Parse each call
	clientPhone := invoices.Phone(phone)
	calls, apiErr := service.GetCalls(input, extraData, clientPhone, computed)
	if apiErr != nil {
		return request, apiErr
	}

	return invoices.Request{
		Phone:  clientPhone,
		Period: computed,
		Calls:  calls,
	}, nil
}

func (service ServiceImpl) GetPeriods(input string) (time.Time, time.Time, apierrors.ApiError) {
	var from, to time.Time
	var err error

	components := strings.Split(input, periodSeparator)
	if len(components) < periodElements {
		return from, to, apierrors.NewBadRequest(fmt.Sprintf("Invalid period %s", input))
	}

	from, err = time.Parse(periodFormat, components[0])
	if err != nil {
		return from, to, apierrors.NewBadRequest(fmt.Sprintf("Invalid period %s: %s", input, err.Error()))
	}

	to, err = time.Parse(periodFormat, components[1])
	if err != nil {
		return from, to, apierrors.NewBadRequest(fmt.Sprintf("Invalid period %s: %s", input, err.Error()))
	}

	return from, to, nil
}

func (service ServiceImpl) GetCalls(input []string, extraData invoices.ExtraData, phone invoices.Phone, period invoices.Period) ([]invoices.Call, apierrors.ApiError) {
	calls := make([]invoices.Call, 0)
	for i, line := range input {
		if i == 0 || strings.TrimSpace(line) == "" {
			// first or empty line
			continue
		}

		call, apiErr := service.GetCall(line, extraData)
		if apiErr != nil {
			return nil, apiErr
		}

		if call.Date.Before(period.From) || call.Date.After(period.To) {
			// discard since call is not in required period
			continue
		}

		target := call.Source
		if call.ReverseCharge {
			target = call.Destination
		}
		if !strings.EqualFold(string(phone), string(target)) {
			// discard since client is not implicated
			continue
		}

		calls = append(calls, call)
	}
	return calls, nil
}

func (service ServiceImpl) GetCall(line string, extraData invoices.ExtraData) (invoices.Call, apierrors.ApiError) {
	var call invoices.Call

	components := strings.Split(line, ",")
	if len(components) < 5 {
		return call, apierrors.NewBadRequest(fmt.Sprintf("Invalid call %s", line))
	}

	if apiErr := call.FillSource(components[0]); apiErr != nil {
		return call, apiErr
	}

	if apiErr := call.FillDestination(components[1]); apiErr != nil {
		return call, apiErr
	}

	if apiErr := call.FillReverseCharge(components[2]); apiErr != nil {
		return call, apiErr
	}

	if apiErr := call.FillDuration(components[3]); apiErr != nil {
		return call, apiErr
	}

	if apiErr := call.FillDate(components[4]); apiErr != nil {
		return call, apiErr
	}

	if apiErr := call.FillCallType(extraData); apiErr != nil {
		return call, apiErr
	}

	return call, nil
}

func (service ServiceImpl) ComputePrices(request *invoices.Request) {
	request.ComputePrices()
}

func (service ServiceImpl) GenerateInvoice(request invoices.Request, extraData invoices.ExtraData) (invoices.Response, apierrors.ApiError) {
	counters := map[invoices.CallScope]time.Duration{
		invoices.CallTypeNational:      0,
		invoices.CallTypeInternational: 0,
		invoices.CallTypeFriends:       0,
	}
	var total float64
	for _, call := range request.Calls {
		counters[call.Scope] += call.ComputableDuration
		total += call.Price
	}

	return invoices.Response{
		User: invoices.User{
			Address: extraData.Address,
			Name:    extraData.Name,
		},
		Movements:                 request.Calls,
		TotalAmount:               total,
		TotalMinutesNational:      counters[invoices.CallTypeNational].Minutes(),
		TotalMinutesInternational: counters[invoices.CallTypeInternational].Minutes(),
		TotalMinutesFriends:       config.FriendsMaxDurationMinutes - request.RemainingFriendMinutes,
	}, nil
}
