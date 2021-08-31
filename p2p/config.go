package p2p

import (
	"fmt"
	"path/filepath"
)

const (
	errValidationStr = "p2p config validation failed - missing val"
)

// Config contains settings of the p2p service
type Config struct {
	// the local IPv4 or IPv6 address
	ListenAddress string `mapstructure:"listen_address" json:"listen_address,omitempty"`

	// the local port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`

	// ip address to bootstrap
	BootstrapIP string `mapstructure:"bootstrap_ip" json:"bootstrap_ip,omitempty"`

	// port to bootstrap
	BootstrapPort int `mapstructure:"bootstrap_port" json:"bootstrap_port,omitempty"`

	// data directory for badger
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`
}

// SetWorkDir applies `workDir` to DataDir if it was not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.DataDir) {
		config.DataDir = filepath.Join(workDir, config.DataDir)
	}
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}

// Validate p2p configs
func (config *Config) Validate() error {
	if config.ListenAddress == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "listen_address")
	}
	if config.Port == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "port")
	}
	if config.DataDir == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "data_dir")
	}

	return nil
}
