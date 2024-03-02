package nftsearch

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// Config contains settings of the registering nft.
type Config struct {
	*common.Config `mapstructure:",squash" json:"-"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config: common.NewConfig(),
	}
}
