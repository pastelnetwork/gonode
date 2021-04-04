package nats

import "github.com/pastelnetwork/go-commons/log"

const formatPrefix = "[nats]"

// Logger wraps go-common logger
type Logger struct{}

// Noticef logs a notice statement
func (Logger) Noticef(format string, v ...interface{}) {
	log.Infof(formatPrefix+format, v...)
}

// Warnf logs a warn statement
func (Logger) Warnf(format string, v ...interface{}) {
	log.Warnf(formatPrefix+format, v...)
}

// Fatalf logs a fatal statement
func (Logger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(formatPrefix+format, v...)
}

// Errorf logs a error statement
func (Logger) Errorf(format string, v ...interface{}) {
	log.Errorf(formatPrefix+format, v...)
}

// Tracef logs a trace statement
func (Logger) Tracef(format string, v ...interface{}) {
	log.Tracef(formatPrefix+format, v...)
}

// Debugf logs a debug statement
func (Logger) Debugf(format string, v ...interface{}) {
	log.Debugf(formatPrefix+format, v...)
}
