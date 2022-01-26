package artworkregister

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberRQIDSFiles uint32 = 1

	defaultThumbnailSize = 224

	defaultNftRegTxMinConfirmations = 12
	defaultNftActTxMinConfirmations = 5

	defaultDDAndFingerprintsMax = 50
	defaultRQIDsMax             = 50
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberRQIDSFiles uint32 `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string

	NFTRegTxMinConfirmations int `mapstructure:"-" json:"reg_art_tx_min_confirmations,omitempty"`
	NFTActTxMinConfirmations int `mapstructure:"-" json:"reg_act_tx_min_confirmations,omitempty"`

	DDAndFingerprintsMax uint32 `mapstructure:"dd_and_fingerprints_max" json:"dd_and_fingerprints_max,omitempty"`
	RQIDsMax             uint32 `mapstructure:"rq_ids_max" json:"rq_ids_max,omitempty"`

	thumbnailSize int
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                   *common.NewConfig(),
		NumberRQIDSFiles:         defaultNumberRQIDSFiles,
		NFTRegTxMinConfirmations: defaultNftRegTxMinConfirmations,
		NFTActTxMinConfirmations: defaultNftActTxMinConfirmations,
		thumbnailSize:            defaultThumbnailSize,
		DDAndFingerprintsMax:     defaultDDAndFingerprintsMax,
		RQIDsMax:                 defaultRQIDsMax,
	}
}
