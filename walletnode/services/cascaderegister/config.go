package cascaderegister

import (
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"path/filepath"
)

const (
	defaultNumberRQIDSFiles uint32 = 1

	defaultCascadeRegTxMinConfirmations = 6
	defaultCascadeActTxMinConfirmations = 2
	defaultWaitTxnValidInterval         = 5

	defaultRQIDsMax = 50
)

var (
	defaultCascadeFilesDir = filepath.Join(configurer.DefaultPath(), "files")
)

// Config contains settings of the registering nft.
type Config struct {
	common.Config `mapstructure:",squash" json:"-"`

	CascadeRegTxMinConfirmations int    `mapstructure:"-" json:"cascade_reg_tx_min_confirmations,omitempty"`
	CascadeActTxMinConfirmations int    `mapstructure:"-" json:"cascade_act_tx_min_confirmations,omitempty"`
	WaitTxnValidInterval         uint32 `mapstructure:"-"`

	RQIDsMax         uint32 `mapstructure:"rq_ids_max" json:"rq_ids_max,omitempty"`
	NumberRQIDSFiles uint32 `mapstructure:"number_rqids_files" json:"number_rqids_files,omitempty"`

	// raptorq service
	RaptorQServiceAddress string `mapstructure:"-" json:"-"`
	RqFilesDir            string
	CascadeFilesDir       string
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Config:                       *common.NewConfig(),
		CascadeRegTxMinConfirmations: defaultCascadeRegTxMinConfirmations,
		CascadeActTxMinConfirmations: defaultCascadeActTxMinConfirmations,
		WaitTxnValidInterval:         defaultWaitTxnValidInterval,
		RQIDsMax:                     defaultRQIDsMax,
		NumberRQIDSFiles:             defaultNumberRQIDSFiles,
		CascadeFilesDir:              defaultCascadeFilesDir,
	}
}
