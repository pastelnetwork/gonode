package userdataprocess

import (
	"time"

	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultConnectToNextNodeDelay = 200 * time.Millisecond
	defaultAcceptNodesTimeout     = 30 * time.Second // = 3 * (2* ConnectToNodeTimeout)
	defaultConnectToNodeTimeout   = time.Second * 15
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	NumberSuperNodes          int `mapstructure:"number_super_nodes" json:"number_super_nodes,omitempty"`
	MinimalNodeConfirmSuccess int `mapstructure:"minimal_node_confirm_success" json:"minimal_node_confirm_success,omitempty"`
	// internal settings
	connectToNextNodeDelay time.Duration
	acceptNodesTimeout     time.Duration
	connectToNodeTimeout   time.Duration
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes:          userdata.DefaultNumberSuperNodes,
		MinimalNodeConfirmSuccess: userdata.MinimalNodeConfirmSuccess,
		connectToNextNodeDelay:    defaultConnectToNextNodeDelay,
		acceptNodesTimeout:        defaultAcceptNodesTimeout,
		connectToNodeTimeout:      defaultConnectToNodeTimeout,
	}
}
