package store

import (
	"context"
	"io"
	golog "log"

	"github.com/hashicorp/go-hclog"
	"github.com/pastelnetwork/gonode/common/log"
)

// Logger wraps go-common logger to implement interface `github.com/hashicorp/hclog.Logger`.
type Logger struct {
	fields log.Fields
	ctx    context.Context
}

// Log emits the message and args at the provided level
func (logger *Logger) Log(level hclog.Level, msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Logf(level.String(), msg, args...)
}

// Debug emits the message and args at DEBUG level
func (logger *Logger) Debug(msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Debugf(msg, args...)
}

// Trace emits the message and args at TRACE level
func (logger *Logger) Trace(msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Tracef(msg, args...)
}

// Info emits the message and args at INFO level
func (logger *Logger) Info(msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Infof(msg, args...)
}

// Warn emits the message and args at WARN level
func (logger *Logger) Warn(msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Warnf(msg, args...)
}

// Error emits the message and args at ERROR level
func (logger *Logger) Error(msg string, args ...interface{}) {
	log.NewDefaultEntry().WithContext(logger.ctx).WithFields(logger.fields).Errorf(msg, args...)
}

// IsTrace indicates that the logger would emit TRACE level logs
func (logger *Logger) IsTrace() bool {
	return true
}

// IsDebug indicates that the logger would emit DEBUG level logs
func (logger *Logger) IsDebug() bool {
	return true
}

// IsInfo indicates that the logger would emit INFO level logs
func (logger *Logger) IsInfo() bool {
	return true
}

// IsWarn indicates that the logger would emit WARN level logs
func (logger *Logger) IsWarn() bool {
	return true
}

// IsError indicates that the logger would emit ERROR level logs
func (logger *Logger) IsError() bool {
	return true
}

// With returns a sub-Logger for which every emitted log message will contain
// the given key/value pairs. This is used to create a context specific
// Logger.
func (logger *Logger) With(args ...interface{}) hclog.Logger {
	if len(args)%2 != 0 {
		return logger
	}

	fields := make(log.Fields, len(logger.fields)+len(args)/2)

	for i := range logger.fields {
		fields[i] = logger.fields[i]
	}

	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		fields[key] = args[i+1]
	}

	return &Logger{
		fields: fields,
	}
}

// Named creates a new sub-Logger that a name decending from the current prefix.
func (logger *Logger) Named(prefix string) hclog.Logger {
	return logger
}

// ResetNamed creates a new sub-Logger that a name decending from the current prefix.
func (logger *Logger) ResetNamed(name string) hclog.Logger {
	return logger.Named(name)
}

// Accept implements the SinkAdapter interface
func (logger *Logger) Accept(name string, level hclog.Level, msg string, args ...interface{}) {}

// ImpliedArgs returns the loggers implied args
func (logger *Logger) ImpliedArgs() []interface{} {
	return []interface{}{}
}

// Name returns the loggers name
func (logger *Logger) Name() string {
	return ""
}

// SetLevel updates the logging level on-the-fly. This will affect all subloggers as well.
func (logger *Logger) SetLevel(level hclog.Level) {}

// StandardLogger creates a *log.Logger that will send it's data through this Logger.
func (logger *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *golog.Logger {
	return golog.New(logger.StandardWriter(nil), "", 0)
}

// StandardWriter returns a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (logger *Logger) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return log.DefaultLogger.Out
}

// NewLogger returns a new Logger instance.
func NewLogger(ctx context.Context) hclog.Logger {
	return &Logger{
		ctx:    ctx,
		fields: make(log.Fields),
	}
}
