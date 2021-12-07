package log

import (
	"testing"
)

// TestLoggers is sample app
func TestLoggers(_ *testing.T) {
	WithSub("hello").WithField("core", "value").Warn("Warn1")
	WithSub("hello").WithField("core", "value").Info("Info1")
	WithSub("hello").WithField("core", "value").Error("Error1")
	SetSubLevelName("hello", "error")
	SetLevelName("error")
	WithSub("hello").WithField("core", "value").Warn("Warn2")
	WithSub("hello").WithField("core", "value").Info("Info2")
	WithSub("hello").WithField("core", "value").Error("Error2")
	WithField("core", "value").Error("default logging")
}
