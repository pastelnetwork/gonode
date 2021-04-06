package hooks

import (
	"github.com/sirupsen/logrus"
)

// ContextHookFormatFn is the callback function for formating message
type ContextHookFormatFn func(ctxValue interface{}, msg string) string

// ContextHook represents a hook for logrus logger.
type ContextHook struct {
	formatFn   ContextHookFormatFn
	contextKey interface{}
}

// Fire implements logrus.ContextHook.Fire()
func (hook *ContextHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	value := entry.Context.Value(hook.contextKey)
	if value != nil {
		entry.Message = hook.formatFn(value, entry.Message)
	}
	return nil
}

// Levels implements logrus.ContextHook.Levels()
func (hook *ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewContextHook creates a new ContextHook instance
func NewContextHook(contextKey interface{}, formatFn ContextHookFormatFn) *ContextHook {
	return &ContextHook{
		contextKey: contextKey,
		formatFn:   formatFn,
	}
}
