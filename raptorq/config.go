package raptorq

import "fmt"

const (
	errValidationStr = "raptorq validation failed - missing val"
	defaultHost      = "localhost"
	defaultPort      = 50051
)

// Config contains settings of the p2p service
type Config struct {
	// the queries IPv4 or IPv6 address
	Host string `mapstructure:"host" json:"host,omitempty"`

	// the queries port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Host: defaultHost,
		Port: defaultPort,
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

	return nil
}
