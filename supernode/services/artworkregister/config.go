package artworkregister

import (
	"time"

	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	defaultNumberConnectedNodes          = 2
	defaultPreburntTxMinConfirmations    = 3
	defaultPreburntTxConfirmationTimeout = 8 * time.Minute
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string

	NumberConnectedNodes int `mapstructure:"-" json:"number_connected_nodes,omitempty"`

	PreburntTxMinConfirmations    int           `mapstructure:"-" json:"preburnt_tx_min_confirmations,omitempty"`
	PreburntTxConfirmationTimeout time.Duration `mapstructure:"-" json:"preburnt_tx_confirmation_timeout,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberConnectedNodes:          defaultNumberConnectedNodes,
		PreburntTxMinConfirmations:    defaultPreburntTxMinConfirmations,
		PreburntTxConfirmationTimeout: defaultPreburntTxConfirmationTimeout,
	}
}
