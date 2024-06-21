package download

import (
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config   `mapstructure:",squash" json:"-"`
	StaticDir       string `mapstructure:"static_dir" json:"static_dir"`
	CascadeFilesDir string `mapstructure:"cascade_files_dir" json:"cascade_files_dir"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config: *common.NewConfig(),
	}
}
