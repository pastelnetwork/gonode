package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	Main `mapstructure:",squash"`

	Server *api.Config `mapstructure:"api" json:"api,omitempty"`
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
		Server: api.NewConfig(),
	}
}
