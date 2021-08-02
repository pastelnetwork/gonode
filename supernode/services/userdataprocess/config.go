package userdataprocess

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
	"github.com/pastelnetwork/gonode/common/service/userdata"
)


// Config contains settings of the process userdata.
type Config struct {
	common.Config        `mapstructure:",squash" json:"-"`
	NumberSuperNodes int `mapstructure:"number_super_nodes" json:"number_super_nodes,omitempty"`
	MinimalNodeConfirmSuccess int  `mapstructure:"minimal_node_confirm_success" json:"minimal_node_confirm_success,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		NumberSuperNodes: userdata.DefaultNumberSuperNodes,
		MinimalNodeConfirmSuccess: userdata.MinimalNodeConfirmSuccess,
	}
}
