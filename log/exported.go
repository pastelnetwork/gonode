package log

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultSkipCallers = 4
)

func logEntry() *logrus.Entry {
	return logEntryWithCallers(defaultLogger, defaultSkipCallers)
}

// WithError adds an error to log entry.
func WithError(err error) *logrus.Entry {
	return logEntry().WithError(err)
}

// WithContext adds an context to log entry, using the value defined inside of context.
func WithContext(ctx context.Context) *logrus.Entry {
	return logEntry().WithContext(ctx)
}

// WithField adds a field to entry.
func WithField(key string, value interface{}) *logrus.Entry {
	return logEntry().WithField(key, value)
}

// WithFields adds multiple fields to entry.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logEntry().WithFields(fields)
}

// WithTime overrides the time of logs generated with it.
func WithTime(t time.Time) *logrus.Entry {
	return logEntry().WithTime(t)
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	logEntry().Debug(args...)
}

// Info logs a message at level Info.
func Info(args ...interface{}) {
	logEntry().Info(args...)
}

// Print logs a message at level Info.
func Print(args ...interface{}) {
	logEntry().Print(args...)
}

// Warn logs a message at level Warn.
func Warn(args ...interface{}) {
	logEntry().Warn(args...)
}

// Error logs a message at level Error.
func Error(args ...interface{}) {
	logEntry().Error(args...)
}

// Panic logs a message at level Panic.
func Panic(args ...interface{}) {
	logEntry().Panic(args...)
}

// Fatal logs a message at level Fatal.
func Fatal(args ...interface{}) {
	logEntry().Fatal(args...)
}

// Debugln logs a message at level Debug.
func Debugln(args ...interface{}) {
	logEntry().Debugln(args...)
}

// Infoln logs a message at level Info.
func Infoln(args ...interface{}) {
	logEntry().Infoln(args...)
}

// Println logs a message at level Info.
func Println(args ...interface{}) {
	logEntry().Println(args...)
}

// Warnln logs a message at level Warn.
func Warnln(args ...interface{}) {
	logEntry().Warnln(args...)
}

// Errorln logs a message at level Error.
func Errorln(args ...interface{}) {
	logEntry().Errorln(args...)
}

// Panicln logs a message at level Panic.
func Panicln(args ...interface{}) {
	logEntry().Panicln(args...)
}

// Fatalln logs a message at level Fatal then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	logEntry().Fatalln(args...)
}

// Debugf logs a message at level Debug.
func Debugf(format string, args ...interface{}) {
	logEntry().Debugf(format, args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	logEntry().Infof(format, args...)
}

// Printf logs a message at level Info.
func Printf(args ...interface{}) {
	logEntry().Print(args...)
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...interface{}) {
	logEntry().Warnf(format, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	logEntry().Errorf(format, args...)
}

// Panicf logs a message at level Panic.
func Panicf(format string, args ...interface{}) {
	logEntry().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logEntry().Fatalf(format, args...)
}
