package userdataprocess

import (
	"time"

	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)

const (
	connectToNextNodeDelay 		= time.Millisecond * 200
	acceptNodesTimeout     		= connectToNextNodeDelay * 10 // waiting 2 seconds (10 supernodes) for secondary nodes to be accpeted by primary nodes.
	connectTimeout        		= time.Second * 2
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes          int `mapstructure:"number_super_nodes" json:"number_super_nodes,omitempty"`
	MinimalNodeConfirmSuccess int `mapstructure:"minimal_node_confirm_success" json:"minimal_node_confirm_success,omitempty"`
	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration
	connectTimeout         time.Duration
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: 			userdata.DefaultNumberSuperNodes,
		MinimalNodeConfirmSuccess:	userdata.MinimalNodeConfirmSuccess,
		connectToNextNodeDelay: 	connectToNextNodeDelay,
		acceptNodesTimeout:     	acceptNodesTimeout,
		connectTimeout:         	connectTimeout,
	}
}
