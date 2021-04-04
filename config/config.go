package config

import (
	"encoding/json"

	"github.com/pastelnetwork/supernode/nats"
	"github.com/pastelnetwork/supernode/pastel"
)

// Config contains configuration of all components of the SuperNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel *pastel.Config `mapstructure:"pastel" json:"pastel"`
	Nats   *nats.Config   `mapstructure:"nats" json:"nats"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		Main:   *NewMain(),
		Pastel: pastel.NewConfig(),
		Nats:   nats.NewConfig(),
	}
}
