package artworkdownload

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3

	defaultConnectToNextNodeDelay = 200 * time.Millisecond
	defaultAcceptNodesTimeout     = 30 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout   = time.Second * 15
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
