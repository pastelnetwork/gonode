package configer

import (
	"path/filepath"

	"github.com/pastelnetwork/go-commons/errors"
	"github.com/spf13/viper"
)

var (
	defaultConfigPaths = []string{".", "$HOME/.pastel"}
)

// SetDefaultConfigPaths sets default paths for Viper to search for the config file in.
func SetDefaultConfigPaths(paths ...string) {
	defaultConfigPaths = paths
}

// ParseFile parses the config file from the given path `filename`, and assign it to the struct `config`.
func ParseFile(filename string, config interface{}) error {
	viper := viper.New()

	for _, configPath := range defaultConfigPaths {
		viper.AddConfigPath(configPath)
	}

	if dir, _ := filepath.Split(filename); dir != "" {
		viper.SetConfigFile(filename)
	} else {
		viper.SetConfigName(filename)
	}

	if err := viper.ReadInConfig(); err != nil {
		return errors.Errorf("could not read config file: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return errors.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}
