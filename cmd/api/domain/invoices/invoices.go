package invoices

import (
	"errors"
	"fmt"
	"github.com/emikohmann/invoices-api/cmd/api/config"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"github.com/emikohmann/invoices-api/pkg/commons"
	"strconv"
	"strings"
	"time"
)

type CallScope string

const (
	CallTypeNational      CallScope = "national"
	CallTypeInternational CallScope = "international"
	CallTypeFriends       CallScope = "friends"
)

const (
	reverseChargeTrue  string = "S"
	reverseChargeFalse string = "N"
	callFormat         string = "2006-01-02T15:04:05Z"
)

var (
	validReverseCharge = []string{
		reverseChargeTrue,
		reverseChargeFalse,
	}
)

type Phone string

type ExtraData struct {
	Address     string   `json:"address"`
	Friends     []string `json:"friends"`
	Name        string   `json:"name"`
	PhoneNumber Phone    `json:"phone_number"`
}

type Request struct {
	Phone                  Phone   `json:"phone"`
	Period                 Period  `json:"period"`
	Calls                  []Call  `json:"calls"`
	RemainingFriendMinutes float64 `json:"remaining_friend_minutes"`
}

type Period struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Call struct {
	Source             Phone         `json:"source"`
	Destination        Phone         `json:"destination"`
	ReverseCharge      bool          `json:"reverse_charge"`
	Duration           time.Duration `json:"-"`
	ComputableDuration time.Duration `json:"-"`
	DurationSeconds    float64       `json:"duration_seconds"`
	Date               time.Time     `json:"date"`
	Scope              CallScope     `json:"scope"`
	Friends            bool          `json:"friends"`
	Price              float64       `json:"price"`
}

type Response struct {
	User                      User    `json:"user"`
	Movements                 []Call  `json:"movements"`
	TotalAmount               float64 `json:"total_amount"`
	TotalMinutesNational      float64 `json:"total_minutes_national"`
	TotalMinutesInternational float64 `json:"total_minutes_international"`
	TotalMinutesFriends       float64 `json:"total_minutes_friends"`
}

type User struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

func (phone Phone) GetCountryCode() (string, error) {
	if len(phone) < 3 {
		return "", errors.New("phone number is too short")
	}
	return string(phone)[0:3], nil
}

func (call *Call) FillSource(input string) apierrors.ApiError {
	sourceInput := strings.TrimSpace(input)
	if sourceInput == "" {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid source input: %s", input))
	}
	call.Source = Phone(sourceInput)
	return nil
}

func (call *Call) FillDestination(input string) apierrors.ApiError {
	destinationInput := strings.TrimSpace(input)
	if destinationInput == "" {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid destination input: %s", input))
	}
	call.Destination = Phone(destinationInput)
	return nil
}

func (call *Call) FillReverseCharge(input string) apierrors.ApiError {
	if !commons.StringContains(validReverseCharge, input) {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid reverse charge %s", input))
	}
	call.ReverseCharge = strings.EqualFold(input, reverseChargeTrue)
	return nil
}

func (call *Call) FillDuration(input string) apierrors.ApiError {
	seconds, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid duration %s: %s", input, err.Error()))
	}
	call.Duration = time.Duration(seconds) * time.Second
	call.DurationSeconds = call.Duration.Seconds()
	return nil
}

func (call *Call) FillDate(input string) apierrors.ApiError {
	var err error
	call.Date, err = time.Parse(callFormat, input)
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid date %s: %s", input, err.Error()))
	}
	return nil
}

func (call *Call) FillCallType(friends ExtraData) apierrors.ApiError {
	sourceCode, err := call.Source.GetCountryCode()
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid source country code %s: %s", sourceCode, err.Error()))
	}

	destinationCode, err := call.Destination.GetCountryCode()
	if err != nil {
		return apierrors.NewBadRequest(fmt.Sprintf("Invalid destination country code %s: %s", destinationCode, err.Error()))
	}

	if sourceCode == destinationCode {
		call.Scope = CallTypeNational
	} else {
		call.Scope = CallTypeInternational
	}

	target := call.Destination
	if call.ReverseCharge {
		target = call.Source
	}
	call.Friends = commons.StringContains(friends.Friends, string(target))
	return nil
}

func (request *Request) ComputePrices() {
	request.RemainingFriendMinutes = config.FriendsMaxDurationMinutes
	for i := range request.Calls {
		if request.Calls[i].Friends && request.RemainingFriendMinutes > 0 {
			request.RemainingFriendMinutes -= request.Calls[i].Duration.Minutes()
			if request.RemainingFriendMinutes >= 0 {
				request.Calls[i].ComputableDuration = 0
				request.Calls[i].Price = 0
				continue
			}
			request.Calls[i].ComputableDuration = time.Duration(-request.RemainingFriendMinutes) * time.Minute
			request.RemainingFriendMinutes = 0
		} else {
			request.Calls[i].ComputableDuration = request.Calls[i].Duration
		}
		if request.Calls[i].ComputableDuration > 0 {
			if request.Calls[i].Scope == CallTypeNational {
				request.Calls[i].Price = config.PriceNationalPerCall
			}
			if request.Calls[i].Scope == CallTypeInternational {
				request.Calls[i].Price = request.Calls[i].ComputableDuration.Minutes() * config.PriceInternationalPerMinute
			}
		}
	}
}
