//go:build windows
// +build windows

package configurer

import (
	"os"
	"path"
	"path/filepath"
	"syscall"
)

const (
	beforeVistaAppDir = "Application Data"
	sinceVistaAppDir  = "AppData/Roaming"
)

var defaultConfigPaths = []string{
	path.Join("$HOME", beforeVistaAppDir, "Pastel"),
	path.Join("$HOME", sinceVistaAppDir, "Pastel"),
	".",
}

// DefaultPath returns the default config path for Windows OS.
func DefaultPath() string {
	homeDir, _ := os.UserHomeDir()
	appDir := beforeVistaAppDir

	v, _ := syscall.GetVersion()
	if v&0xff > 5 {
		appDir = sinceVistaAppDir
	}
	return filepath.Join(homeDir, filepath.FromSlash(appDir), "Pastel")
}
