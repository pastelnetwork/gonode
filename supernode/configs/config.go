package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/pastel-client"
	"github.com/pastelnetwork/gonode/supernode/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
)

// Config contains configuration of all components of the SuperNode.
type Config struct {
	Main `mapstructure:",squash"`

	Pastel          *pastel.Config          `mapstructure:"pastel" json:"pastel,omitempty"`
	Server          *server.Config          `mapstructure:"server" json:"server,omitempty"`
	ArtworkRegister *artworkregister.Config `mapstructure:"artwork_register" json:"artwork_register,omitempty"`
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
		Main:            *NewMain(),
		Pastel:          pastel.NewConfig(),
		Server:          server.NewConfig(),
		ArtworkRegister: artworkregister.NewConfig(),
	}
}
