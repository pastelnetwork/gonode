// +build linux

package configurer

import (
	"os"
	"path/filepath"
)

var defaultConfigPaths = []string{
	"$HOME/.pastel",
	".",
}

// DefaultConfigPath returns the default config path for Linux OS.
func DefaultConfigPath(filename string) string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".pastel", filename)
}
