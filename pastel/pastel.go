package pastel

import (
	"github.com/pastelnetwork/go-pastel"
)

// Client represents pastel client
var Client *pastel.Client

// Init inits pastel client with the give config
func Init(config *Config) error {
	Client = pastel.NewClient(config.Hostname, config.Username, config.Password)
	return nil
}
