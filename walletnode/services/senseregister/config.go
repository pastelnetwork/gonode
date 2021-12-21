package senseregister

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes       = 3
	defaultConnectToNextNodeDelay = 200 * time.Millisecond
	defaultAcceptNodesTimeout     = 30 * time.Second // = 3 * (2* ConnectToNodeTimeout)

)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes       int           `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`
	connectToNextNodeDelay time.Duration `mapstructure:"connect_to_next_node_delay" json:"connect_to_next_node_delay,omitempty"`
	acceptNodesTimeout     time.Duration `mapstructure:"accept_nodes_timeout" json:"accept_nodes_timeout,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Config:                 *common.NewConfig(),
		NumberSuperNodes:       defaultNumberSuperNodes,
		connectToNextNodeDelay: defaultConnectToNextNodeDelay,
		acceptNodesTimeout:     defaultAcceptNodesTimeout,
	}
}
