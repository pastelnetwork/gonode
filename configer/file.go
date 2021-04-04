package configer

import (
	"path/filepath"
	"strings"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/spf13/viper"
)

var (
	defaultConfigPaths = []string{".", "$HOME/.pastel"}
	defaultConfigType  = "yml"
)

// SetConfigPaths sets paths for Viper to search for the config file in.
func SetConfigPaths(paths ...string) {
	defaultConfigPaths = paths
}

// ParseFile parses the config file from the given path `filename`, and assign it to the struct `val`.
func ParseFile(filename string, val interface{}) error {
	viper := viper.New()

	configPath, filename := filepath.Split(filename)

	if configPath != "" {
		viper.AddConfigPath(configPath)
	} else {
		for _, configPath := range defaultConfigPaths {
			viper.AddConfigPath(configPath)
		}
	}

	if ext := filepath.Ext(filename); ext != "" {
		filename = strings.TrimSuffix(filename, ext)
		viper.SetConfigType(ext[1:])
	} else {
		viper.SetConfigType(defaultConfigType)
	}

	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return errors.Errorf("error reading Cfg file: %s", err)
	}

	if err := viper.Unmarshal(val); err != nil {
		return errors.Errorf("unable to decode into struct, %v", err)
	}
	return nil
}
