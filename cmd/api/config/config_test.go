package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {
	assert.EqualValues(t, "invoices-api", ApplicationName)
}

func TestBusinessValues(t *testing.T) {
	assert.EqualValues(t, 150, FriendsMaxDurationMinutes)
	assert.EqualValues(t, 2.5, PriceNationalPerCall)
	assert.EqualValues(t, 20, PriceInternationalPerMinute)
}
