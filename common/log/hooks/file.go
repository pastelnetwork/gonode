package hooks

import (
	"fmt"
	"math"
	"time"

	"github.com/pastelnetwork/gonode/common/log/formatters"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultFormatter = formatters.LogFile
	defaultMaxSize   = 104857600 // 100Mb
	defaultMaxAge    = time.Second * 86400
	defaultCompress  = false
)

// FileHook is a hook to handle writing to local log files.
type FileHook struct {
	fileLogger *lumberjack.Logger
	formatter  logrus.Formatter
}

// SetFormatter sets the format that will be used by hook.
func (hook *FileHook) SetFormatter(formatter logrus.Formatter) {
	hook.formatter = formatter
}

// SetMaxAge sets the maximum duration to retain old log files based on the timestamp encoded in their filename.
func (hook *FileHook) SetMaxAge(maxAge time.Duration) {
	days, _ := math.Modf(maxAge.Hours() / 24)
	hook.fileLogger.MaxAge = int(days)
}

// SetMaxSize sets the maximum size in megabytes of the log file before it gets rotated.
func (hook *FileHook) SetMaxSize(maxSize int) {
	size := maxSize / 1048576 // to Megabytes
	hook.fileLogger.MaxSize = size
}

// Fire writes the log file to defined filename path.
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to generate string for entry, %s", err)
	}

	_, err = hook.fileLogger.Write(msg)
	return err
}

// Levels returns configured log levels.
func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// NewFileHook returns new FileHook instance with default values.
func NewFileHook(filename string) *FileHook {
	hook := &FileHook{
		fileLogger: &lumberjack.Logger{
			Filename: filename,
			Compress: defaultCompress,
		},
	}

	hook.SetFormatter(defaultFormatter)
	hook.SetMaxAge(defaultMaxAge)
	hook.SetMaxSize(defaultMaxSize)

	return hook
}
