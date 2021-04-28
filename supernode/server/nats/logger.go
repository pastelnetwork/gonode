package nats

import "github.com/pastelnetwork/gonode/common/log"

const (
	logPrefix = "[nats]"
)

// Logger wraps go-common logger to implement interface `github.com/nats-io/nats-server.Logger`.
type Logger struct {
	prefix string
}

// Noticef logs a notice statement.
func (logger *Logger) Noticef(format string, v ...interface{}) {
	log.Infof(logger.prefix+format, v...)
}

// Warnf logs a warn statement.
func (logger *Logger) Warnf(format string, v ...interface{}) {
	log.Warnf(logger.prefix+format, v...)
}

// Fatalf logs a fatal statement.
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(logger.prefix+format, v...)
}

// Errorf logs a error statement.
func (logger *Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(logger.prefix+format, v...)
}

// Tracef logs a trace statement.
func (logger *Logger) Tracef(format string, v ...interface{}) {
	log.Tracef(logger.prefix+format, v...)
}

// Debugf logs a debug statement.
func (logger *Logger) Debugf(format string, v ...interface{}) {
	log.Debugf(logger.prefix+format, v...)
}

// NewLogger returns a new Logger instance.
func NewLogger() *Logger {
	return &Logger{
		prefix: logPrefix,
	}
}
