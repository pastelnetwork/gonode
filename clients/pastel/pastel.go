package pastel

import (
	"github.com/pastelnetwork/go-pastel"
)

func New(config *Config) *pastel.Client {
	return pastel.NewClient(config.Hostname, config.Username, config.Password)
}
