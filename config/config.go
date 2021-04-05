package config

import (
	"encoding/json"

	"github.com/pastelnetwork/walletnode/nats"
	"github.com/pastelnetwork/walletnode/pastel"
	"github.com/pastelnetwork/walletnode/rest"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel *pastel.Config `mapstructure:"pastel" json:"pastel"`
	Nats   *nats.Config   `mapstructure:"nats" json:"nats"`
	Rest   *rest.Config   `mapstructure:"rest" json:"rest"`
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
		Rest:   rest.NewConfig(),
	}
}
