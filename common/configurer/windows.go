// +build windows

package configurer

import "path/filepath"

var defaultConfigPaths = []string{
	".",
	filepath.FromSlash("C:/Documents and Settings/Username/Application Data/Pastel"),
	filepath.FromSlash("C:/Users/Username/AppData/Roaming/Pastel"),
}
