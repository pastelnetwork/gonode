package artworkregister

import "github.com/pastelnetwork/gonode/pastel"

const (
	defaultNumberConnectedNodes = 2
)

// Config contains settings of the registering artwork.
type Config struct {
	NumberConnectedNodes int `mapstructure:"number_connected_nodes" json:"number_connected_nodes,omitempty"`

	node *pastel.MasterNodeConfig
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberConnectedNodes: defaultNumberConnectedNodes,
	}
}
