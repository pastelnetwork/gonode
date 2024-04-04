package ddclient

import (
	"fmt"
	"path/filepath"
)

const (
	errValidationStr = "ddserver client validation failed - missing val"
)

// Config contains settings of the dd-server
type Config struct {
	// Host the queries IPv4 or IPv6 address
	Host string `mapstructure:"host" json:"host,omitempty"`

	// Port the queries port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`

	// DDFilesDir - the location of temporary folder to transfer image data to ddserver
	DDFilesDir string `mapstructure:"dd-temp-file-dir" json:"dd-temp-file-dir,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}

// SetWorkDir update working dir
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.DDFilesDir) {
		config.DDFilesDir = filepath.Join(workDir, config.DDFilesDir)
	}
}

// Validate raptorq configs
func (config *Config) Validate() error {
	if config.Host == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "host")
	}
	if config.Port == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "port")
	}

	if config.DDFilesDir == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "dd-temp-file-dir")
	}

	return nil
}
