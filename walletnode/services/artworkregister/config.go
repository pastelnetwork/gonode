package artworkregister

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultNumberRQIDSFiles uint32 = 1

	defaultConnectToNextNodeDelay = 200 * time.Millisecond
	defaultAcceptNodesTimeout     = 30 * time.Second // = 3 * (2* ConnectToNodeTimeout)

	defaultThumbnailSize = 224

	defaultRegArtTxMinConfirmations = 12
	defaultRegActTxMinConfirmations = 5

	defaultDDAndFingerprintsMax = 50
	defaultRQIDsMax             = 50
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`

	NumberRQIDSFiles uint32 `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string

	// BurnAddress
	BurnAddress string `mapstructure:"-" json:"burn_address,omitempty"`

	RegArtTxMinConfirmations int `mapstructure:"-" json:"reg_art_tx_min_confirmations,omitempty"`
	RegActTxMinConfirmations int `mapstructure:"-" json:"reg_act_tx_min_confirmations,omitempty"`

	DDAndFingerprintsMax uint32 `mapstructure:"dd_and_fingerprints_max" json:"dd_and_fingerprints_max,omitempty"`
	RQIDsMax             uint32 `mapstructure:"rq_ids_max" json:"rq_ids_max,omitempty"`
	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration

	thumbnailSize int
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                   *common.NewConfig(),
		NumberSuperNodes:         defaultNumberSuperNodes,
		NumberRQIDSFiles:         defaultNumberRQIDSFiles,
		RegArtTxMinConfirmations: defaultRegArtTxMinConfirmations,
		RegActTxMinConfirmations: defaultRegActTxMinConfirmations,
		connectToNextNodeDelay:   defaultConnectToNextNodeDelay,
		acceptNodesTimeout:       defaultAcceptNodesTimeout,
		thumbnailSize:            defaultThumbnailSize,
		DDAndFingerprintsMax:     defaultDDAndFingerprintsMax,
		RQIDsMax:                 defaultRQIDsMax,
	}
}
