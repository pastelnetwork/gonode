package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	defaultSkipCallers = 5
)

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

// Entry represents a wrapped logrus.Entry
type Entry struct {
	*logrus.Entry
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{Entry: entry.Entry.WithField(key, value)}
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	return &Entry{Entry: entry.Entry.WithFields(logrus.Fields(fields))}
}

// WithError adds an field error to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return &Entry{Entry: entry.Entry.WithError(err)}
}

// WithContext adds a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	return &Entry{Entry: entry.Entry.WithContext(ctx)}
}

// WithPrefix adds a prefix as single field (using the key `prefix`) to the Entry.
func (entry *Entry) WithPrefix(value string) *Entry {
	return &Entry{Entry: entry.Entry.WithField("prefix", value)}
}

// Trace logs a message at level Trace.
func (entry *Entry) Trace(args ...interface{}) {
	entry.Entry.Trace(args...)
}

// Debug logs a message at level Debug.
func (entry *Entry) Debug(args ...interface{}) {
	entry.Entry.Debug(args...)
}

// Print logs a message at level Print.
func (entry *Entry) Print(args ...interface{}) {
	entry.Entry.Print(args...)
}

// Info logs a message at level Info.
func (entry *Entry) Info(args ...interface{}) {
	entry.Entry.Info(args...)
}

// Warn logs a message at level Warn.
func (entry *Entry) Warn(args ...interface{}) {
	entry.Entry.Warn(args...)
}

// Error logs a message at level Error.
func (entry *Entry) Error(args ...interface{}) {
	entry.Entry.Error(args...)
}

// Fatal logs a message at level Fatal.
func (entry *Entry) Fatal(args ...interface{}) {
	entry.Entry.Fatal(args...)
}

// Panic logs a message at level Panic.
func (entry *Entry) Panic(args ...interface{}) {
	entry.Entry.Panic(args...)
}

// Entry Printf family functions

// Tracef logs a message at level Trace.
func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.Entry.Tracef(format, args...)
}

// Debugf logs a message at level Debug.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Entry.Debugf(format, args...)
}

// Infof logs a message at level Info.
func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Entry.Infof(format, args...)
}

// Printf logs a message at level Print.
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Entry.Printf(format, args...)
}

// Warnf logs a message at level Warn.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Entry.Warnf(format, args...)
}

// Errorf logs a message at level Error.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Entry.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Entry.Fatalf(format, args...)
}

// Panicf logs a message at level Panic.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Entry.Panicf(format, args...)
}

// Entry Println family functions

// Traceln logs a message at level Trace.
func (entry *Entry) Traceln(args ...interface{}) {
	entry.Entry.Traceln(args...)
}

// Debugln logs a message at level Debug.
func (entry *Entry) Debugln(args ...interface{}) {
	entry.Entry.Debugln(args...)
}

// Infoln logs a message at level Info.
func (entry *Entry) Infoln(args ...interface{}) {
	entry.Entry.Infoln(args...)
}

// Println logs a message at level Print.
func (entry *Entry) Println(args ...interface{}) {
	entry.Entry.Println(args...)
}

// Warnln logs a message at level Warn.
func (entry *Entry) Warnln(args ...interface{}) {
	entry.Entry.Warnln(args...)
}

// Errorln logs a message at level Error.
func (entry *Entry) Errorln(args ...interface{}) {
	entry.Entry.Errorln(args...)
}

// Fatalln logs a message at level Fatal.
func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Entry.Fatalln(args...)
}

// Panicln logs a message at level Panic.
func (entry *Entry) Panicln(args ...interface{}) {
	entry.Entry.Panicln(args...)
}

// WithCaller adds caller field to the Entry.
func (entry *Entry) WithCaller(skip int) *Entry {
	if debugMode {
		entry.Data["file"] = fileInfo(skip)
	}
	return entry
}

// NewEntry returns a new Entry instance.
func NewEntry(logger *Logger) *Entry {
	return (&Entry{
		Entry: logrus.NewEntry(logger.Logger),
	})
}

// NewDefaultEntry returns a new Entry instance using default logger.
func NewDefaultEntry() *Entry {
	return NewEntry(DefaultLogger).WithCaller(defaultSkipCallers)
}
