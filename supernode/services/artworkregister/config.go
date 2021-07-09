package artworkregister

import "github.com/pastelnetwork/gonode/supernode/services/common"

const (
	defaultNumberConnectedNodes = 2
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"raptorq_service" json:"raptorq_service,omitempty"`

	NumberConnectedNodes int `mapstructure:"number_connected_nodes" json:"number_connected_nodes,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberConnectedNodes: defaultNumberConnectedNodes,
	}
}
