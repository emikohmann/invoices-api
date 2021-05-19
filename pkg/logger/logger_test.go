package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTags(t *testing.T) {
	assert.EqualValues(t, "[level:%s] ", tagLevel)
	assert.EqualValues(t, "[application:%s] ", tagApplication)
}

func TestLevels(t *testing.T) {
	assert.EqualValues(t, "INFO", levelInfo)
	assert.EqualValues(t, "WARN", levelWarn)
	assert.EqualValues(t, "ERROR", levelError)
	assert.EqualValues(t, "PANIC", levelPanic)
}

func TestInfo(t *testing.T) {
	Info("This is a test message")
	Info("This is a test message with args %s", "test_arg")
}

func TestWarn(t *testing.T) {
	Warn("This is a test message")
	Warn("This is a test message with args %s", "test_arg")
}

func TestError(t *testing.T) {
	Error("This is a test message")
	Error("This is a test message with args %s", "test_arg")
}

func TestPanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec == nil {
			t.Error("The test did not panic")
		}
	} ()
	Panic("This is a test message")
	Panic("This is a test message with args %s", "test_arg")
}

func TestWrite(t *testing.T) {
	write("TEST_LEVEL", "This is a test message")
	write("TEST_LEVEL", "This is a test message with args %s", "test_arg")
}
