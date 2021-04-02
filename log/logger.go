package log

import (
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel = logrus.InfoLevel
)

// Logger represents a wrapped logrus.Logger
type Logger struct {
	*logrus.Logger
}

// NewLogger returns a new Logger instance with default values
func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(TextFormatter())
	logger.SetLevel(defaultLevel)

	return &Logger{
		Logger: logrus.New(),
	}
}
