package configurer

import (
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/spf13/viper"
)

// SetDefaultConfigPaths sets default paths for Viper to search for the config file in.
func SetDefaultConfigPaths(paths ...string) {
	defaultConfigPaths = paths
}

// ParseFile parses the config file from the given path `filename`, and assign it to the struct `config`.
func ParseFile(filename string, config interface{}) error {
	var configType string

	switch filepath.Ext(filename) {
	case ".conf":
		configType = "env"
	}

	return parseFile(filename, configType, config)
}

// ParseJSONFile parses json config file from the given path `filename`, and assign it to the struct `config`.
func ParseJSONFile(filename string, config interface{}) error {
	return parseFile(filename, "json", config)
}

func parseFile(filename, configType string, config interface{}) error {
	conf := viper.New()

	for _, configPath := range defaultConfigPaths {
		conf.AddConfigPath(configPath)
	}

	if dir, _ := filepath.Split(filename); dir != "" {
		conf.SetConfigFile(filename)
	} else {
		conf.SetConfigName(filename)
	}

	if configType != "" {
		conf.SetConfigType(configType)
	}

	if err := conf.ReadInConfig(); err != nil {
		return errors.Errorf("could not read config file: %w", err)
	}

	if err := conf.Unmarshal(&config); err != nil {
		return errors.Errorf("unable to decode into struct, %w", err)
	}

	return nil
}
