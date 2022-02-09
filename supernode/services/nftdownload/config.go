package nftdownload

import "github.com/pastelnetwork/gonode/supernode/services/common"

// Config contains settings of the registering Nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config: *common.NewConfig(),
	}
}
