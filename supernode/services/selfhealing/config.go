package selfhealing

import (
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

// Config self-healing config
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config: *common.NewConfig(),
	}
}
