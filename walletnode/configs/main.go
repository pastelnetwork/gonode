package configs

import (
	"os"
	"path/filepath"
)

const (
	defaultLogLevel = "info"
)

// Main contains main config of the App
type Main struct {
	LogLevel string `mapstructure:"log-level" json:"log-level,omitempty"`
	LogFile  string `mapstructure:"log-file" json:"log-file,omitempty"`
	Quiet    bool   `mapstructure:"quiet" json:"quiet"`
	WorkDir  string `mapstructure:"work-dir" json:"work-dir"`
}

// NewMain returns a new Main instance.
func NewMain() *Main {
	return &Main{
		LogLevel: defaultLogLevel,
		WorkDir:  filepath.Join(os.TempDir(), "walletnode"),
	}
}
