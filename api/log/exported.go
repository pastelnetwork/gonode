package log

import "github.com/pastelnetwork/go-commons/log"

// Errorf logs a error statement.
func Errorf(format string, v ...interface{}) {
	log.Errorf(prefix+format, v...)
}

// Warningf logs a warn statement.
func Warningf(format string, v ...interface{}) {
	log.Warnf(prefix+format, v...)
}

// Infof logs a notice statement.
func Infof(format string, v ...interface{}) {
	log.Infof(prefix+format, v...)
}

// Debugf logs a debug statement.
func Debugf(format string, v ...interface{}) {
	log.Debugf(prefix+format, v...)
}
