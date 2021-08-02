package artworkdownload

import "github.com/pastelnetwork/gonode/supernode/services/common"

// Config contains settings of the registering artwork.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"raptorq_service" json:"raptorq_service,omitempty"`
	RqFilesDir            string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}
