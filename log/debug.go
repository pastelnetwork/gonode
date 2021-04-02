package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	runDir string

	// if `debugMode` is true, log record will contain `file` field, with a value that indicating where the log was called.
	debugMode bool
)

// SetDebugMode sets the state of debug mode
func SetDebugMode(isEnabled bool) {
	debugMode = isEnabled
}

func logEntryWithCallers(logger *Logger, skip int) *logrus.Entry {
	entry := logrus.NewEntry(logger.Logger)
	if debugMode {
		entry.Data["file"] = fileInfo(skip)
	}
	return entry
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		file, _ = filepath.Rel(runDir, file)
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func init() {
	runDir, _ = os.Getwd()
}
