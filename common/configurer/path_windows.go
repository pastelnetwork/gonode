// +build windows

package configurer

import (
	"os"
	"path/filepath"
	"syscall"
)

const (
	beforeVistaAppDir = "Application Data"
	fromVistaAppDir   = "AppData/Roaming"
)

var defaultConfigPaths = []string{
	file.Path("$HOME", beforeVistaAppDir, "Pastel"),
	file.Path("$HOME", fromVistaAppDir, "Pastel"),
	".",
}

// DefaultConfigPath returns the default config path for Windows OS.
func DefaultConfigPath(filename string) string {
	homeDir := os.UserHomeDir()
	appDir := beforeVistaAppDir

	v, _ := syscall.GetVersion()
	if v&0xff > 5 {
		appDir = fromVistaAppDir
	}
	return filepath.Join(homeDir, filepath.FromSlash(appDir), "Pastel", filename)
}
