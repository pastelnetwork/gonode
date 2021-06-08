package configs

import (
	"encoding/json"
	"github.com/pastelnetwork/gonode/pastel"
	superCfg "github.com/pastelnetwork/gonode/supernode/configs"
	walletCfg "github.com/pastelnetwork/gonode/walletnode/configs"
)

const (
	defaultLogLevel = "info"

	ZksnarkParamsURL = "https://z.cash/downloads/"
)

var (
	ZksnarkParamsNames = []string{
		"sapling-spend.params",
		"sapling-output.params",
		"sprout-proving.key",
		"sprout-verifying.key",
		"sprout-groth16.params",
	}
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	DefaultDir string `json:"-"`

	LogLevel string `mapstructure:"log-level" json:"log-level,omitempty"`
	LogFile  string `mapstructure:"log-file" json:"log-file,omitempty"`
	Quite    bool   `mapstructure:"quite" json:"quite"`
	TempDir  string `mapstructure:"temp-dir" json:"temp-dir"`
	WorkDir  string `mapstructure:"work-dir" json:"work-dir"`

	Pastel     *pastel.Config    `mapstructure:"pastel-api" json:"pastel-api,omitempty"`
	Init       *Init             `mapstructure:"init" json:"init"`
	Start      *Start            `mapstructure:"start" json:"start"`
	WalletNode *walletCfg.Config `json:"-"`
	SuperNode  *superCfg.Config  `json:"-"`
}

// String : returns string from Config fields
func (config *Config) String() (string, error) {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data), nil
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		LogLevel:   defaultLogLevel,
		Init:       NewInit(),
		Start:      NewStart(),
		Pastel:     pastel.NewConfig(),
		WalletNode: walletCfg.New(),
		SuperNode:  superCfg.New(),
	}
}
