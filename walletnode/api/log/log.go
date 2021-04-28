package log

import (
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	prefix = "api"
)

func newEntry() *log.Entry {
	return log.NewDefaultEntry().WithPrefix(prefix)
}

// WithField adds a field to log entry.
func WithField(key string, value interface{}) *log.Entry {
	return newEntry().WithField(key, value)
}

// WithError adds an error to log entry.
func WithError(err error) *log.Entry {
	return newEntry().WithError(err)
}

// Infof logs a info statement.
func Infof(format string, v ...interface{}) {
	newEntry().Infof(format, v...)
}

// Warnf logs a warn statement.
func Warnf(format string, v ...interface{}) {
	newEntry().Warnf(format, v...)
}

// Fatalf logs a fatal statement.
func Fatalf(format string, v ...interface{}) {
	newEntry().Fatalf(format, v...)
}

// Errorf logs a error statement.
func Errorf(format string, v ...interface{}) {
	newEntry().Errorf(format, v...)
}

// Tracef logs a trace statement.
func Tracef(format string, v ...interface{}) {
	newEntry().Tracef(format, v...)
}

// Debugf logs a debug statement.
func Debugf(format string, v ...interface{}) {
	newEntry().Debugf(format, v...)
}
