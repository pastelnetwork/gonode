package userdataprocess

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)


// Config contains settings of the registering artwork.
type Config struct {
	common.Config        `mapstructure:",squash" json:"-"`
	NumberConnectedNodes int `mapstructure:"number_connected_nodes" json:"number_connected_nodes,omitempty"`
	MinimalNodeConfirmSuccess int  `mapstructure:"minimal_node_confirm_success" json:"minimal_node_confirm_success,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberConnectedNodes: userdata.DefaultNumberSuperNodes - 1, // NumberConnectedNodes won't count the primary node
		MinimalNodeConfirmSuccess: userdata.MinimalNodeConfirmSuccess,
	}
}
