package formatters

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Global formatters
var (
	LogFile  = NewFileFormatter()
	Terminal = NewTerminalFormatter()
)

// NewFileFormatter returns a new Formatter instance for log file.
func NewFileFormatter() logrus.Formatter {
	return &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05.000",
		DisableColors:   true,
		ForceFormatting: true,
	}
}

// NewTerminalFormatter returns a new Formatter instance for terminal.
func NewTerminalFormatter() logrus.Formatter {
	formatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05.000",
		ForceFormatting: true,
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue+h",
		PrefixStyle:     "cyan",
		TimestampStyle:  "black+h",
	})

	return formatter
}
