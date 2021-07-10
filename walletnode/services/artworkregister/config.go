package artworkregister

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultNumberRQIDSFiles = 3

	connectToNextNodeDelay = time.Millisecond * 200
	acceptNodesTimeout     = connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	connectTimeout         = time.Second * 2

	thumbnailSize = 224
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`

	NumberRQIDSFiles int `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"raptorq_service" json:"raptorq_service,omitempty"`

	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration
	connectTimeout         time.Duration

	thumbnailSize int
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,

		NumberRQIDSFiles: defaultNumberRQIDSFiles,

		connectToNextNodeDelay: connectToNextNodeDelay,
		acceptNodesTimeout:     acceptNodesTimeout,
		connectTimeout:         connectTimeout,

		thumbnailSize: thumbnailSize,
	}
}
