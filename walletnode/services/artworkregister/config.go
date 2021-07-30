package artworkregister

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultNumberRQIDSFiles uint32 = 10

	connectToNextNodeDelay = time.Millisecond * 200
	acceptNodesTimeout     = connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accepted by primary nodes.

	thumbnailSize = 224

	regArtTxMinConfirmations = 10
	regArtTxTimeout          = 26 * time.Minute

	regActTxMinConfirmations = 5
	regActTxTimeout          = 13 * time.Minute
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`

	NumberRQIDSFiles uint32 `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"raptorq_service" json:"raptorq_service,omitempty"`
	RqFilesDir            string

	// BurnAddress
	BurnAddress string `mapstructure:"burn_address" json:"burn_address,omitempty"`

	RegArtTxMinConfirmations int           `mapstructure:"reg_art_tx_min_confirmations" json:"reg_art_tx_min_confirmations,omitempty"`
	RegArtTxTimeout          time.Duration `mapstructure:"reg_art_tx_timeout" json:"reg_art_tx_timeout,omitempty"`

	RegActTxMinConfirmations int           `mapstructure:"reg_act_tx_min_confirmations" json:"reg_act_tx_min_confirmations,omitempty"`
	RegActTxTimeout          time.Duration `mapstructure:"reg_act_tx_timeout" json:"reg_act_tx_timeout,omitempty"`

	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration

	thumbnailSize int
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,

		NumberRQIDSFiles: defaultNumberRQIDSFiles,

		RegArtTxMinConfirmations: regArtTxMinConfirmations,
		RegArtTxTimeout:          regArtTxTimeout,
		RegActTxMinConfirmations: regActTxMinConfirmations,
		RegActTxTimeout:          regActTxTimeout,
		connectToNextNodeDelay:   connectToNextNodeDelay,
		acceptNodesTimeout:       acceptNodesTimeout,

		thumbnailSize: thumbnailSize,
	}
}
