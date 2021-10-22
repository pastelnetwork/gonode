package debug

import (
	"encoding/json"
)

const (
	defaultPort = 8080
)

// Config contains configuration of debug service
type Config struct {
	// the local port to listen for connections on
	HttpPort int `mapstructure:"port" json:"http-port,omitempty"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func NewConfig() *Config {
	return &Config{
		HttpPort: defaultPort,
	}
}
