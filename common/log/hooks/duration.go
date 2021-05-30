package hooks

import (
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// DurationFieldName is field key.
	DurationFieldName = "duration"
)

// DurationHook represents a hook for logrus logger.
type DurationHook struct{}

// Fire implements logrus.DurationHook.Fire()
func (hook *DurationHook) Fire(entry *logrus.Entry) error {
	field, ok := entry.Data[DurationFieldName]
	if !ok {
		return nil
	}

	start, ok := field.(time.Time)
	if !ok {
		return nil
	}
	delete(entry.Data, DurationFieldName)
	entry.Data[DurationFieldName] = time.Since(start).String()

	return nil
}

// Levels implements logrus.DurationHook.Levels()
func (hook *DurationHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewDurationHook creates a new DurationHook instance
func NewDurationHook() *DurationHook {
	return &DurationHook{}
}
