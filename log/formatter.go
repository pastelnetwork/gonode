package log

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func TextFormatter() logrus.Formatter {
	formatter := &prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "Jan 02 15:04:05.000",
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
