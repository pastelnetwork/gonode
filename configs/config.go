package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/go-pastel"
	"github.com/pastelnetwork/walletnode/api"
	"github.com/pastelnetwork/walletnode/services/artwork/register"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel          *pastel.Config   `mapstructure:"pastel" json:"pastel,omitempty"`
	Rest            *api.Config      `mapstructure:"api" json:"api,omitempty"`
	ArtworkRegister *register.Config `mapstructure:"artwork_register" json:"artwork_register,omitempty"`
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

		Pastel:          pastel.NewConfig(),
		ArtworkRegister: register.NewConfig(),
		Rest:            api.NewConfig(),
	}
}
