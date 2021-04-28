package hooks

import (
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/sirupsen/logrus"
)

// ContextHookFormatFn is the callback function for formating message
type ContextHookFormatFn func(ctxValue interface{}, msg string) string

// ContextHook represents a hook for logrus logger.
type ContextHook struct {
	contextKey interface{}
	fn         func(entry *log.Entry, ctxValue interface{})
}

// Fire implements logrus.ContextHook.Fire()
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	ctxValue := entry.Context.Value(hook.contextKey)
	if ctxValue != nil {
		hook.fn(&log.Entry{Entry: entry}, ctxValue)
	}
	return nil
}

// Levels implements logrus.ContextHook.Levels()
func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewContextHook creates a new ContextHook instance
func NewContextHook(contextKey interface{}, fn func(entry *log.Entry, ctxValue interface{})) *ContextHook {
	return &ContextHook{
		contextKey: contextKey,
		fn:         fn,
	}
}
