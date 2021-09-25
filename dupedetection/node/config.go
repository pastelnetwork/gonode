package node

import "fmt"

const (
	errValidationStr = "ddserver client validation failed - missing val"
)

// Config contains settings of the dd-server
type Config struct {
	// the local IPv4 or IPv6 address
	Host string `mapstructure:"host" json:"host,omitempty"`

	// the local port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}

// Validate raptorq configs
func (config *Config) Validate() error {
	if config.Host == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "host")
	}
	if config.Port == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "port")
	}

	return nil
}
