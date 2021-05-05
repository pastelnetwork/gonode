package log

import (
	"github.com/pastelnetwork/gonode/common/log/formatters"
	"github.com/sirupsen/logrus"
)

const (
	defaultLevel = logrus.DebugLevel
)

// Logger represents a wrapped logrus.Logger
type Logger struct {
	*logrus.Logger
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
