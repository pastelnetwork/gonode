package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log/hooks"
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

// WithDuration adds an field `duration` with value of the time elapsed since t.
func WithDuration(t time.Time) *Entry {
	return NewDefaultEntry().WithField(hooks.DurationFieldName, t)
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

// Logf logs a message at the provided level
func Logf(levelName, format string, args ...interface{}) {
	NewDefaultEntry().Logf(levelName, format, args...)
}

// Fatal logs an error at level Fatal with error stack.
func Fatal(err error) {
	entry := NewDefaultEntry()
	if debugMode {
		entry = entry.WithErrorStack(err)
	}
	entry.Fatal(err)
}

// FatalAndExit checks if there is an error, display it in the console and exit with a non-zero exit code. Otherwise, exit 0.
// Note that if the debugMode is true, this will print out the stack trace.
func FatalAndExit(err error) {
	if err == nil || errors.IsContextCanceled(err) {
		return
	}
	defer os.Exit(errors.ExitCode(err))

	Fatal(err)

	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
}

// WithSub return a log entry - with log level of given sub domain - using SetSubLevelName() to set sub level
// but even with specific default level, but sub domain logs are also filtered by DefaultLogger's log level
func WithSub(subName string) *Entry {
	return NewEntry(DefaultLogger, getSubLogLevel(subName))
}

// P2P return a log entry - of p2p subsystem
func P2P() *Entry {
	return NewEntry(DefaultLogger, getSubLogLevel(p2pSubName))
}

// MetaDB return a log entry - of metadb subsystem
func MetaDB() *Entry {
	return NewEntry(DefaultLogger, getSubLogLevel(metadbSubName))
}

// DD return a log entry of dupe detection subsystem
func DD() *Entry {
	return NewEntry(DefaultLogger, getSubLogLevel(ddSubName))
}
