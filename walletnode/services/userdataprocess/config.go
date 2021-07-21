package userdataprocess

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes 	= 10
	minimalNodeConfirmSuccess 	= 8
	connectToNextNodeDelay 		= time.Millisecond * 200
	acceptNodesTimeout     		= connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	connectTimeout        		= time.Second * 2
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes int `mapstructure:"number_supernodes" json:"number_supernodes,omitempty"`
	MinimalNodeConfirmSuccess int `mapstructure:"number_confirm_success" json:"number_confirm_success,omitempty"`
	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration
	connectTimeout         time.Duration
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: 			defaultNumberSuperNodes,
		MinimalNodeConfirmSuccess:	minimalNodeConfirmSuccess,
		connectToNextNodeDelay: 	connectToNextNodeDelay,
		acceptNodesTimeout:     	acceptNodesTimeout,
		connectTimeout:         	connectTimeout,
	}
}
