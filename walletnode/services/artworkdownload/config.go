package artworkdownload

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay = time.Millisecond * 200
	defaultAcceptNodesTimeout     = defaultConnectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	defaultConnectToNodeTimeout   = time.Second * 5
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`

	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration
	connectToNodeTimeout   time.Duration
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: defaultNumberSuperNodes,

		connectToNextNodeDelay: defaultConnectToNextNodeDelay,
		acceptNodesTimeout:     defaultAcceptNodesTimeout,
		connectToNodeTimeout:   defaultConnectToNodeTimeout,
	}
}
