package pastel

import (
	"github.com/pastelnetwork/go-pastel"
)

// New returns a new instance of the pastel client.
func New(config *Config) *pastel.Client {
	return pastel.NewClient(config.Hostname, config.Username, config.Password)
}
