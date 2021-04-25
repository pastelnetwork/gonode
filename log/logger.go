package log

import (
	"github.com/pastelnetwork/go-commons/log/formatters"
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel = logrus.InfoLevel
)

// Logger represents a wrapped logrus.Logger
type Logger struct {
	*logrus.Logger
}

// Entry represents a wrapped logrus.Entry
type Entry struct {
	*logrus.Entry
}

// NewLogger returns a new Logger instance with default values
func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(formatters.Terminal)
	logger.SetLevel(defaultLevel)

	return &Logger{
		Logger: logger,
	}
}
