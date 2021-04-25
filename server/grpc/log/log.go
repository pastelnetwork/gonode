package log

import (
	"github.com/pastelnetwork/go-commons/log"
)

const (
	prefix = "[grpc]"
)

// WithField adds a field to log entry.
func WithField(key string, value interface{}) *log.Entry {
	return &log.Entry{log.WithPrefix(prefix).WithField(key, value)}
}

// WithError adds an error to log entry.
func WithError(err error) *log.Entry {
	return &log.Entry{log.WithPrefix(prefix).WithError(err)}
}

// Infof logs a info statement.
func Infof(format string, v ...interface{}) {
	log.WithPrefix(prefix).Infof(format, v...)
}

// Warnf logs a warn statement.
func Warnf(format string, v ...interface{}) {
	log.WithPrefix(prefix).Warnf(format, v...)
}

// Fatalf logs a fatal statement.
func Fatalf(format string, v ...interface{}) {
	log.WithPrefix(prefix).Fatalf(format, v...)
}

// Errorf logs a error statement.
func Errorf(format string, v ...interface{}) {
	log.WithPrefix(prefix).Errorf(format, v...)
}

// Tracef logs a trace statement.
func Tracef(format string, v ...interface{}) {
	log.WithPrefix(prefix).Tracef(format, v...)
}

// Debugf logs a debug statement.
func Debugf(format string, v ...interface{}) {
	log.WithPrefix(prefix).Debugf(format, v...)
}
