package senseregister

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultSenseRegTxMinConfirmations = 12
	defaultSenseActTxMinConfirmations = 5
	defaultWaitTxnValidInterval       = 5

	defaultDDAndFingerprintsMax = 50
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	SenseRegTxMinConfirmations int    `mapstructure:"-" json:"sense_reg_tx_min_confirmations,omitempty"`
	SenseActTxMinConfirmations int    `mapstructure:"-" json:"sense_act_tx_min_confirmations,omitempty"`
	WaitTxnValidInterval       uint32 `mapstructure:"-"`

	DDAndFingerprintsMax uint32 `mapstructure:"dd_and_fingerprints_max" json:"dd_and_fingerprints_max,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                     *common.NewConfig(),
		SenseRegTxMinConfirmations: defaultSenseRegTxMinConfirmations,
		SenseActTxMinConfirmations: defaultSenseActTxMinConfirmations,
		DDAndFingerprintsMax:       defaultDDAndFingerprintsMax,
		WaitTxnValidInterval:       defaultWaitTxnValidInterval,
	}
}
