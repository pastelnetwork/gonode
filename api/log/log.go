package log

import "github.com/pastelnetwork/go-commons/log"

const (
	Prefix = "[rest]"
)

// Errorf logs a error statement.
func Errorf(format string, v ...interface{}) {
	log.Errorf(Prefix+format, v...)
}

// Warningf logs a warn statement.
func Warningf(format string, v ...interface{}) {
	log.Warnf(Prefix+format, v...)
}

// Infof logs a notice statement.
func Infof(format string, v ...interface{}) {
	log.Infof(Prefix+format, v...)
}

// Debugf logs a debug statement.
func Debugf(format string, v ...interface{}) {
	log.Debugf(Prefix+format, v...)
}
