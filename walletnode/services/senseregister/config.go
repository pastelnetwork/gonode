package senseregister

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay = 200 * time.Millisecond
	defaultAcceptNodesTimeout     = 30 * time.Second // = 3 * (2* ConnectToNodeTimeout)

	defaultThumbnailSize = 224

	defaultSenseRegTxMinConfirmations = 12
	defaultSenseActTxMinConfirmations = 5

	defaultDDAndFingerprintsMax = 50
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`

	//// BurnAddress
	//BurnAddress string `mapstructure:"-" json:"burn_address,omitempty"`

	SenseRegTxMinConfirmations int `mapstructure:"-" json:"sense_reg_tx_min_confirmations,omitempty"`
	SenseActTxMinConfirmations int `mapstructure:"-" json:"sense_act_tx_min_confirmations,omitempty"`

	DDAndFingerprintsMax uint32 `mapstructure:"dd_and_fingerprints_max" json:"dd_and_fingerprints_max,omitempty"`
	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration

	thumbnailSize int
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                     *common.NewConfig(),
		NumberSuperNodes:           defaultNumberSuperNodes,
		SenseRegTxMinConfirmations: defaultSenseRegTxMinConfirmations,
		SenseActTxMinConfirmations: defaultSenseActTxMinConfirmations,
		connectToNextNodeDelay:     defaultConnectToNextNodeDelay,
		acceptNodesTimeout:         defaultAcceptNodesTimeout,
		thumbnailSize:              defaultThumbnailSize,
		DDAndFingerprintsMax:       defaultDDAndFingerprintsMax,
	}
}
