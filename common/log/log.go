package log

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	metadbSubName = "metadb"
	p2pSubName    = "p2p"
	ddSubName     = "dd"
)

// DefaultLogger is logger with default settings
var DefaultLogger = NewLogger()
var subLevels = map[string]logrus.Level{}
var currentDefaultLevel = logrus.InfoLevel
var mtx = sync.RWMutex{}

// SetLevelName parses and sets the defaultLogger level.
func SetLevelName(name string) error {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return err
	}

	currentDefaultLevel = level

	return nil
}

// SetOutput sets the defaultLogger output.
func SetOutput(output io.Writer) {
	DefaultLogger.SetOutput(output)
}

// SetSubLevelName parses and sets the level for a subsystem with given name.
func SetSubLevelName(subName string, levelName string) error {
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		return err
	}
	mtx.Lock()
	defer mtx.Unlock()
	subLevels[subName] = level
	return nil
}

// SetP2PLogLevelName sets log level of p2p subsystem
func SetP2PLogLevelName(levelName string) error {
	return SetSubLevelName(p2pSubName, levelName)
}

// SetMetaDBLogLevelName sets log level of metadb subsystem
func SetMetaDBLogLevelName(levelName string) error {
	return SetSubLevelName(metadbSubName, levelName)
}

// SetDDLogLevelName sets log level of dupe subsystem
func SetDDLogLevelName(levelName string) error {
	return SetSubLevelName(ddSubName, levelName)
}

func getSubLogLevel(subName string) logrus.Level {
	mtx.RLock()
	defer mtx.RUnlock()
	level, ok := subLevels[subName]
	if !ok {
		// default level
		return currentDefaultLevel
	}

	return level
}

// AddHook adds hook to an instance of defaultLogger.
func AddHook(hook logrus.Hook) {
	DefaultLogger.Hooks.Add(hook)
}
