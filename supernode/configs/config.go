package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/pastel"
)

// Config contains configuration of all components of the SuperNode.
type Config struct {
	Main `mapstructure:",squash"`

	Node   *Node          `mapstructure:"node" json:"node,omitempty"`
	Pastel *pastel.Config `mapstructure:"pastel-api" json:"pastel-api,omitempty"`
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
		Node:   NewNode(),
		Pastel: pastel.NewConfig(),
	}
}
