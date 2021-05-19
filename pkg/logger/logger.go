package logger

import (
	"bytes"
	"fmt"
	"github.com/emikohmann/invoices-api/cmd/api/config"
)

const (
	tagLevel       = "[level:%s] "
	tagApplication = "[application:%s] "
	levelInfo      = "INFO"
	levelWarn      = "WARN"
	levelError     = "ERROR"
	levelPanic     = "PANIC"
)

func Info(message string, args ...interface{}) {
	write(levelInfo, message, args...)
}

func Warn(message string, args ...interface{}) {
	write(levelWarn, message, args...)
}

func Error(message string, args ...interface{}) {
	write(levelError, message, args...)
}

func Panic(message string, args ...interface{}) {
	write(levelPanic, message, args...)
	panic(fmt.Sprintf(message, args...))
}

func write(level string, message string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(tagLevel, level))
	buffer.WriteString(fmt.Sprintf(tagApplication, config.ApplicationName))
	buffer.WriteString(message)
	result := fmt.Sprintf(buffer.String(), args...)
	fmt.Println(result)
}
