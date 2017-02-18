// Package log is a wrapper around github.com/Sirupsen/logrus to enable formatted
// options to be the same globally.
package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

var (
	logger = logrus.New()
)

type (
	// Fields are a map of fields to print out in a log message.
	Fields map[string]interface{}
)

func init() {
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	logger.Out = colorable.NewColorableStderr()
	logger.Level = logrus.DebugLevel
}

// Debug will send a log message to the debug level.
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info will send a log message to the info level.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Warn will send a log message to the warn level.
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Error will send a log message to the error level.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal will send a log message to the fatal level.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Panic will send a log message to the panic level.
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Printf will print a standard log message
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// WithFields adds fields to a log message.
func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}
