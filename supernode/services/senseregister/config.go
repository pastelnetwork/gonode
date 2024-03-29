package senseregister

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	defaultNumberConnectedNodes       = 2
	defaultPreburntTxMinConfirmations = 3
)

// Config contains settings of the registering Nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
	DDDatabase            string `mapstructure:"-" json:"-"`

	NumberConnectedNodes       int `mapstructure:"-" json:"number_connected_nodes,omitempty"`
	PreburntTxMinConfirmations int `mapstructure:"-" json:"preburnt_tx_min_confirmations,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                     *common.NewConfig(),
		NumberConnectedNodes:       defaultNumberConnectedNodes,
		PreburntTxMinConfirmations: defaultPreburntTxMinConfirmations,
	}
}
