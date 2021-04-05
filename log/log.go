package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

var defaultLogger = NewLogger()

// SetLevelName parses and sets the defaultLogger level.
func SetLevelName(name string) error {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return err
	}
	defaultLogger.SetLevel(level)

	return nil
}

// SetOutput sets the defaultLogger output.
func SetOutput(output io.Writer) {
	defaultLogger.SetOutput(output)
}

// AddHook adds hook to an instance of defaultLogger.
func AddHook(hook logrus.Hook) {
	defaultLogger.Hooks.Add(hook)
}
