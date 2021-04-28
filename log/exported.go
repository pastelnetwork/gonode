package log

import (
	"context"
)

func newEntry() *Entry {
	return NewDefaultEntry()
}

// WithError adds an error to log entry.
func WithError(err error) *Entry {
	return newEntry().WithError(err)
}

// WithPrefix adds a prefix to log entry.
func WithPrefix(value string) *Entry {
	return newEntry().WithField("prefix", value)
}

// WithContext adds an context to log entry, using the value defined inside of context.
func WithContext(ctx context.Context) *Entry {
	return newEntry().WithContext(ctx)
}

// WithField adds a field to entry.
func WithField(key string, value interface{}) *Entry {
	return newEntry().WithField(key, value)
}

// WithFields adds multiple fields to entry.
func WithFields(fields Fields) *Entry {
	return newEntry().WithFields(fields)
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	newEntry().Debug(args...)
}

// Info logs a message at level Info.
func Info(args ...interface{}) {
	newEntry().Info(args...)
}

// Print logs a message at level Info.
func Print(args ...interface{}) {
	newEntry().Print(args...)
}

// Warn logs a message at level Warn.
func Warn(args ...interface{}) {
	newEntry().Warn(args...)
}

// Error logs a message at level Error.
func Error(args ...interface{}) {
	newEntry().Error(args...)
}

// Panic logs a message at level Panic.
func Panic(args ...interface{}) {
	newEntry().Panic(args...)
}

// Fatal logs a message at level Fatal.
func Fatal(args ...interface{}) {
	newEntry().Fatal(args...)
}

// Debugln logs a message at level Debug.
func Debugln(args ...interface{}) {
	newEntry().Debugln(args...)
}

// Infoln logs a message at level Info.
func Infoln(args ...interface{}) {
	newEntry().Infoln(args...)
}

// Println logs a message at level Info.
func Println(args ...interface{}) {
	newEntry().Println(args...)
}

// Warnln logs a message at level Warn.
func Warnln(args ...interface{}) {
	newEntry().Warnln(args...)
}

// Errorln logs a message at level Error.
func Errorln(args ...interface{}) {
	newEntry().Errorln(args...)
}

// Panicln logs a message at level Panic.
func Panicln(args ...interface{}) {
	newEntry().Panicln(args...)
}

// Fatalln logs a message at level Fatal then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	newEntry().Fatalln(args...)
}

// Debugf logs a message at level Debug.
func Debugf(format string, args ...interface{}) {
	newEntry().Debugf(format, args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	newEntry().Infof(format, args...)
}

// Printf logs a message at level Info.
func Printf(args ...interface{}) {
	newEntry().Print(args...)
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...interface{}) {
	newEntry().Warnf(format, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	newEntry().Errorf(format, args...)
}

// Tracef logs a message at level Trace.
func Tracef(format string, args ...interface{}) {
	newEntry().Tracef(format, args...)
}

// Panicf logs a message at level Panic.
func Panicf(format string, args ...interface{}) {
	newEntry().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	newEntry().Fatalf(format, args...)
}
