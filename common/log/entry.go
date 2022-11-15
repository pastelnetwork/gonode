package log

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log/hooks"
	"github.com/sirupsen/logrus"
)

const (
	defaultSkipCallers = 4
)

// Fields type, used to pass to `WithFields`.
type Fields logrus.Fields

// Entry represents a wrapped logrus.Entry
type Entry struct {
	*logrus.Entry
	level logrus.Level
}

// WithField adds a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{Entry: entry.Entry.WithField(key, value), level: entry.level}
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	return &Entry{Entry: entry.Entry.WithFields(logrus.Fields(fields)), level: entry.level}
}

// WithErrorStack adds an field `error` and `stack` to the Entry.
func (entry *Entry) WithErrorStack(err error) *Entry {
	return entry.WithError(err).WithField("stack", errors.ErrorStack(err))
}

// WithError adds an field `error` to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	fields := errors.ExtractFields(err)
	return &Entry{Entry: entry.Entry.WithError(err).WithFields(map[string]interface{}(fields)), level: entry.level}
}

// WithContext adds a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	ent := &Entry{Entry: entry.Entry.WithContext(ctx), level: entry.level}
	server := ctx.Value("server")
	return ent.WithField("server_ip", server)
}

// WithDuration adds an field `duration` with value of the time elapsed since t.
func (entry *Entry) WithDuration(t time.Time) *Entry {
	return &Entry{Entry: entry.Entry.WithField(hooks.DurationFieldName, t), level: entry.level}
}

// WithPrefix adds a prefix as single field (using the key `prefix`) to the Entry.
func (entry *Entry) WithPrefix(value string) *Entry {
	return &Entry{Entry: entry.Entry.WithField("prefix", value), level: entry.level}
}

// Trace logs a message at level Trace.
func (entry *Entry) Trace(args ...interface{}) {
	if entry.level >= logrus.TraceLevel {
		entry.Entry.Trace(args...)
	}
}

// Debug logs a message at level Debug.
func (entry *Entry) Debug(args ...interface{}) {
	if entry.level >= logrus.DebugLevel {
		entry.Entry.Debug(args...)
	}
}

// Print logs a message at level Print.
func (entry *Entry) Print(args ...interface{}) {
	entry.Entry.Print(args...)
}

// Info logs a message at level Info.
func (entry *Entry) Info(args ...interface{}) {
	if entry.level >= logrus.InfoLevel {
		entry.Entry.Info(args...)
	}
}

// Warn logs a message at level Warn.
func (entry *Entry) Warn(args ...interface{}) {
	if entry.level >= logrus.WarnLevel {
		entry.Entry.Warn(args...)
	}
}

// Error logs a message at level Error.
func (entry *Entry) Error(args ...interface{}) {
	if entry.level >= logrus.ErrorLevel {
		entry.Entry.Error(args...)
	}
}

// Fatal logs a message at level Fatal.
func (entry *Entry) Fatal(args ...interface{}) {
	if entry.level >= logrus.FatalLevel {
		entry.Entry.Log(logrus.FatalLevel, args...)
	}
}

// Entry Printf family functions

// Tracef logs a message at level Trace.
func (entry *Entry) Tracef(format string, args ...interface{}) {
	if entry.level >= logrus.TraceLevel {
		entry.Entry.Tracef(format, args...)
	}
}

// Debugf logs a message at level Debug.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.level >= logrus.DebugLevel {
		entry.Entry.Debugf(format, args...)
	}
}

// Infof logs a message at level Info.
func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.level >= logrus.InfoLevel {
		entry.Entry.Infof(format, args...)
	}
}

// Printf logs a message at level Print.
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Entry.Printf(format, args...)
}

// Warnf logs a message at level Warn.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.level >= logrus.WarnLevel {
		entry.Entry.Warnf(format, args...)
	}
}

// Errorf logs a message at level Error.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.level >= logrus.ErrorLevel {
		entry.Entry.Errorf(format, args...)
	}
}

// Fatalf logs a message at level Fatal.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.level >= logrus.FatalLevel {
		entry.Entry.Fatalf(format, args...)
	}
}

// Logf logs a message at the provided level
func (entry *Entry) Logf(levelName, format string, args ...interface{}) {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.DebugLevel
	}
	if entry.level >= level {
		entry.Entry.Logf(level, format, args...)
	}
}

// Entry Println family functions

// Traceln logs a message at level Trace.
func (entry *Entry) Traceln(args ...interface{}) {
	if entry.level >= logrus.TraceLevel {
		entry.Entry.Traceln(args...)
	}
}

// Debugln logs a message at level Debug.
func (entry *Entry) Debugln(args ...interface{}) {
	if entry.level >= logrus.DebugLevel {
		entry.Entry.Debugln(args...)
	}
}

// Infoln logs a message at level Info.
func (entry *Entry) Infoln(args ...interface{}) {
	if entry.level >= logrus.InfoLevel {
		entry.Entry.Infoln(args...)
	}
}

// Println logs a message at level Print.
func (entry *Entry) Println(args ...interface{}) {
	entry.Entry.Println(args...)
}

// Warnln logs a message at level Warn.
func (entry *Entry) Warnln(args ...interface{}) {
	if entry.level >= logrus.WarnLevel {
		entry.Entry.Warnln(args...)
	}
}

// Errorln logs a message at level Error.
func (entry *Entry) Errorln(args ...interface{}) {
	if entry.level >= logrus.ErrorLevel {
		entry.Entry.Errorln(args...)
	}
}

// Fatalln logs a message at level Fatal.
func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.level >= logrus.FatalLevel {
		entry.Entry.Fatalln(args...)
	}
}

// WithCaller adds caller field to the Entry.
func (entry *Entry) WithCaller(skip int) *Entry {
	if debugMode {
		entry.Data["file"] = fileInfo(skip)
	}
	return entry
}

// NewEntry returns a new Entry instance.
func NewEntry(logger *Logger, level logrus.Level) *Entry {
	return (&Entry{
		Entry: logrus.NewEntry(logger.Logger),
		level: level,
	})
}

// NewDefaultEntry returns a new Entry instance using default logger.
func NewDefaultEntry() *Entry {
	return NewEntry(DefaultLogger, currentDefaultLevel).WithCaller(defaultSkipCallers)
}
