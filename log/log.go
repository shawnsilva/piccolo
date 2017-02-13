package log

import (
	"github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

var (
	logger = logrus.New()
)

type (
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

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func WithFields(fields Fields) *logrus.Entry {
	return logger.WithFields(logrus.Fields(fields))
}
