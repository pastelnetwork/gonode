package artworkdownload

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	defaultNumberSuperNodes = 3
)

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}
