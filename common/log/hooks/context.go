package hooks

import (
	"github.com/sirupsen/logrus"
)

// ContextHookFields is log fields
type ContextHookFields map[string]interface{}

// ContextHookFn is the callback function for formating message
type ContextHookFn func(ctxValue interface{}, msg string, fields ContextHookFields) (string, ContextHookFields)

// ContextHook represents a hook for logrus logger.
type ContextHook struct {
	contextKey interface{}
	fn         ContextHookFn
}

// Fire implements logrus.ContextHook.Fire()
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	ctxValue := entry.Context.Value(hook.contextKey)
	if ctxValue != nil {
		msg, fields := hook.fn(ctxValue, entry.Message, ContextHookFields(entry.Data))

		entry.Message = msg
		entry.Data = logrus.Fields(fields)
	}
	return nil
}

// Levels implements logrus.ContextHook.Levels()
func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewContextHook creates a new ContextHook instance
func NewContextHook(contextKey interface{}, fn ContextHookFn) *ContextHook {
	return &ContextHook{
		contextKey: contextKey,
		fn:         fn,
	}
}
