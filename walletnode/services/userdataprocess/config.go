package userdataprocess

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	minimalNodeConfirmSuccess = 2
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	MinimalNodeConfirmSuccess int `mapstructure:"minimal_node_confirm_success" json:"minimal_node_confirm_success,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                    *common.NewConfig(),
		MinimalNodeConfirmSuccess: minimalNodeConfirmSuccess,
	}
}
