package debug

import (
	json "github.com/json-iterator/go"
)

const (
	defaultPort = 9091
)

// Config contains configuration of debug service
type Config struct {
	// HTTPPort the queries port to listen for connections on
	HTTPPort int `mapstructure:"http-port" json:"http-port,omitempty"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		HTTPPort: defaultPort,
	}
}
