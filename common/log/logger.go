package log

import (
	"github.com/pastelnetwork/gonode/common/log/formatters"
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
	logger.SetFormatter(formatters.Terminal)
	logger.SetLevel(defaultLevel)

	return &Logger{
		Logger: logger,
	}
}

// NewLoggerWithErrorLevel returns a new Logger instance with default values
func NewLoggerWithErrorLevel() *Logger {
	logger := logrus.New()
	logger.SetFormatter(formatters.Terminal)
	logger.SetLevel(logrus.ErrorLevel)

	return &Logger{
		Logger: logger,
	}
}

// V sets verbose level
func (l *Logger) V(level int) bool {
	switch level {
	case 0:
		return true // always log errors
	case 1:
		return true // always log warnings
	case 2:
		return false // always log info
	case 3:
		return false // never log debug
	default:
		return false
	}

}
