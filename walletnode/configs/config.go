package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/pastel-client"
	"github.com/pastelnetwork/gonode/walletnode/api"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel *pastel.Config `mapstructure:"pastel" json:"pastel,omitempty"`
	API    *api.Config    `mapstructure:"api" json:"api,omitempty"`
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
		Main: *NewMain(),

		Pastel: pastel.NewConfig(),
		API:    api.NewConfig(),
	}
}
