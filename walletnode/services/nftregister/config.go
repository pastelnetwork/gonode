package nftregister

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberRQIDSFiles uint32 = 1

	defaultNftRegTxMinConfirmations = 6
	defaultNftActTxMinConfirmations = 2
	defaultWaitTxnValidInterval     = 5

	defaultDDAndFingerprintsMax = 50
	defaultRQIDsMax             = 50

	defaultThumbnailSize = 224
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NFTRegTxMinConfirmations int    `mapstructure:"-" json:"reg_art_tx_min_confirmations,omitempty"`
	NFTActTxMinConfirmations int    `mapstructure:"-" json:"reg_act_tx_min_confirmations,omitempty"`
	WaitTxnValidInterval     uint32 `mapstructure:"-"`

	DDAndFingerprintsMax uint32 `mapstructure:"dd_and_fingerprints_max" json:"dd_and_fingerprints_max,omitempty"`
	RQIDsMax             uint32 `mapstructure:"rq_ids_max" json:"rq_ids_max,omitempty"`
	NumberRQIDSFiles     uint32 `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	thumbnailSize int

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                   *common.NewConfig(),
		NFTRegTxMinConfirmations: defaultNftRegTxMinConfirmations,
		NFTActTxMinConfirmations: defaultNftActTxMinConfirmations,
		WaitTxnValidInterval:     defaultWaitTxnValidInterval,
		DDAndFingerprintsMax:     defaultDDAndFingerprintsMax,
		RQIDsMax:                 defaultRQIDsMax,
		NumberRQIDSFiles:         defaultNumberRQIDSFiles,
		thumbnailSize:            defaultThumbnailSize,
	}
}
