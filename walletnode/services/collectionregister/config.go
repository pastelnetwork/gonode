package collectionregister

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultCollectionRegTxMinConfirmations = 6
	defaultCollectionActTxMinConfirmations = 2
	defaultWaitTxnValidInterval            = 5
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	CollectionRegTxMinConfirmations int    `mapstructure:"-" json:"collection_reg_tx_min_confirmations,omitempty"`
	CollectionActTxMinConfirmations int    `mapstructure:"-" json:"collection_act_tx_min_confirmations,omitempty"`
	WaitTxnValidInterval            uint32 `mapstructure:"-"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                          *common.NewConfig(),
		CollectionRegTxMinConfirmations: defaultCollectionRegTxMinConfirmations,
		CollectionActTxMinConfirmations: defaultCollectionActTxMinConfirmations,
		WaitTxnValidInterval:            defaultWaitTxnValidInterval,
	}
}
