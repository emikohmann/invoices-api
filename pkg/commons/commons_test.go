package commons

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringContains(t *testing.T) {
	assert.True(t, StringContains([]string{"a", "b", "c"}, "a"))
	assert.True(t, StringContains([]string{"a", "b", "c"}, "b"))
	assert.True(t, StringContains([]string{"a", "b", "c"}, "c"))

	assert.True(t, StringContains([]string{"a", "b"}, "a"))
	assert.True(t, StringContains([]string{"a", "b"}, "b"))
	assert.False(t, StringContains([]string{"a", "b"}, "c"))

	assert.True(t, StringContains([]string{"a"}, "a"))
	assert.False(t, StringContains([]string{"a"}, "b"))
	assert.False(t, StringContains([]string{"a"}, "c"))

	assert.False(t, StringContains([]string{}, "a"))
	assert.False(t, StringContains([]string{}, "b"))
	assert.False(t, StringContains([]string{}, "c"))
}
