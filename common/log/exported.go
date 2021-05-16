package log

import (
	"context"
)

// WithError adds an error to log entry.
func WithError(err error) *Entry {
	return NewDefaultEntry().WithError(err)
}

// WithErrorStack adds an `error` and `stack` to log Entry.
func WithErrorStack(err error) *Entry {
	return NewDefaultEntry().WithErrorStack(err)
}

// WithPrefix adds a prefix to log entry.
func WithPrefix(value string) *Entry {
	return NewDefaultEntry().WithField("prefix", value)
}

// WithContext adds an context to log entry, using the value defined inside of context.
func WithContext(ctx context.Context) *Entry {
	return NewDefaultEntry().WithContext(ctx)
}

// WithField adds a field to entry.
func WithField(key string, value interface{}) *Entry {
	return NewDefaultEntry().WithField(key, value)
}

// WithFields adds multiple fields to entry.
func WithFields(fields Fields) *Entry {
	return NewDefaultEntry().WithFields(fields)
}

// Debug logs a message at level Debug.
func Debug(args ...interface{}) {
	NewDefaultEntry().Debug(args...)
}

// Info logs a message at level Info.
func Info(args ...interface{}) {
	NewDefaultEntry().Info(args...)
}

// Print logs a message at level Info.
func Print(args ...interface{}) {
	NewDefaultEntry().Print(args...)
}

// Warn logs a message at level Warn.
func Warn(args ...interface{}) {
	NewDefaultEntry().Warn(args...)
}

// Error logs a message at level Error.
func Error(args ...interface{}) {
	NewDefaultEntry().Error(args...)
}

// Recover logs a message at level Panic with error stack.
func Recover(err error) {
	entry := NewDefaultEntry()
	if debugMode {
		entry.WithErrorStack(err)
	}
	entry.Panic(err)
}

// Panic logs a message at level Panic.
func Panic(args ...interface{}) {
	NewDefaultEntry().Panic(args...)
}

// Fatal logs a message at level Fatal.
func Fatal(args ...interface{}) {
	NewDefaultEntry().Fatal(args...)
}

// Debugln logs a message at level Debug.
func Debugln(args ...interface{}) {
	NewDefaultEntry().Debugln(args...)
}

// Infoln logs a message at level Info.
func Infoln(args ...interface{}) {
	NewDefaultEntry().Infoln(args...)
}

// Println logs a message at level Info.
func Println(args ...interface{}) {
	NewDefaultEntry().Println(args...)
}

// Warnln logs a message at level Warn.
func Warnln(args ...interface{}) {
	NewDefaultEntry().Warnln(args...)
}

// Errorln logs a message at level Error.
func Errorln(args ...interface{}) {
	NewDefaultEntry().Errorln(args...)
}

// Panicln logs a message at level Panic.
func Panicln(args ...interface{}) {
	NewDefaultEntry().Panicln(args...)
}

// Fatalln logs a message at level Fatal then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	NewDefaultEntry().Fatalln(args...)
}

// Debugf logs a message at level Debug.
func Debugf(format string, args ...interface{}) {
	NewDefaultEntry().Debugf(format, args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	NewDefaultEntry().Infof(format, args...)
}

// Printf logs a message at level Info.
func Printf(args ...interface{}) {
	NewDefaultEntry().Print(args...)
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...interface{}) {
	NewDefaultEntry().Warnf(format, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	NewDefaultEntry().Errorf(format, args...)
}

// Tracef logs a message at level Trace.
func Tracef(format string, args ...interface{}) {
	NewDefaultEntry().Tracef(format, args...)
}

// Panicf logs a message at level Panic.
func Panicf(format string, args ...interface{}) {
	NewDefaultEntry().Panicf(format, args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	NewDefaultEntry().Fatalf(format, args...)
}
