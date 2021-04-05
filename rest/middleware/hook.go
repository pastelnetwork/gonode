package middleware

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"goa.design/goa/v3/middleware"
)

// Hook represents a hook for logrus logger.
type Hook struct{}

// Fire implements logrus.Hook.Fire()
func (hook *Hook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	value := entry.Context.Value(middleware.RequestIDKey)
	if value != nil {
		entry.Message = fmt.Sprintf("%v [%v] %s", logPrefix, value, entry.Message)
	}
	return nil
}

// Levels implements logrus.Hook.Levels()
func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewHook creates a new Hook instance
func NewHook() *Hook {
	return &Hook{}
}
