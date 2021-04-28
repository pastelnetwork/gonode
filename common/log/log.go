package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

// DefaultLogger is logger with default settings
var DefaultLogger = NewLogger()

// SetLevelName parses and sets the defaultLogger level.
func SetLevelName(name string) error {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return err
	}
	DefaultLogger.SetLevel(level)

	return nil
}

// SetOutput sets the defaultLogger output.
func SetOutput(output io.Writer) {
	DefaultLogger.SetOutput(output)
}

// AddHook adds hook to an instance of defaultLogger.
func AddHook(hook logrus.Hook) {
	DefaultLogger.Hooks.Add(hook)
}
