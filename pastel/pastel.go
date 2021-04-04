package pastel

import (
	"github.com/pastelnetwork/go-pastel"
)

var Client *pastel.Client

func Init(config *Config) error {
	Client = pastel.NewClient(config.Hostname, config.Username, config.Password)
	return nil
}
