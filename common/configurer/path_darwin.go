// +build darwin

package configurer

import (
	"os"
	"path/filepath"
)

var defaultConfigPaths = []string{
	"$HOME/Library/Application Support/Pastel",
	".",
}

// DefaultConfigPath returns the default config path for darwin OS.
func DefaultConfigPath(filename string) string {
	homeDir, _ := os.UserConfigDir()
	return filepath.Join(homeDir, "Pastel", filename)
}
