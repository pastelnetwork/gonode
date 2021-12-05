package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/raptorq"
)

const (
	defaultLogMaxAgeInDays = 3
	defaultLogMaxSizeInMB  = 100
	defaultLogCompress     = true
	defaultLogMaxBackups   = 10
	defaultLogLevel        = "info"
)

// Config contains configuration of all components of the WalletNode.
type Config struct {
	LogLevel        string `mapstructure:"log-level" json:"log-level,omitempty"`
	LogFile         string `mapstructure:"log-file" json:"log-file,omitempty"`
	LogCompress     bool   `mapstructure:"log-compress" json:"log-compress,omitempty"`
	LogMaxSizeInMB  int    `mapstructure:"log-max-size-mb" json:"log-max-size-mb,omitempty"`
	LogMaxAgeInDays int    `mapstructure:"log-max-age-days" json:"log-max-age-days,omitempty"`
	LogMaxBackups   int    `mapstructure:"log-max-backups" json:"log-max-backups,omitempty"`
	Quiet           bool   `mapstructure:"quiet" json:"quiet"`
	TempDir         string `mapstructure:"temp-dir" json:"temp-dir"`
	WorkDir         string `mapstructure:"work-dir" json:"work-dir"`
	RqFilesDir      string `mapstructure:"rq-files-dir" json:"rq-files-dir"`

	Node    `mapstructure:"node" json:"node,omitempty"`
	Pastel  *pastel.Config  `mapstructure:"-" json:"-"`
	RaptorQ *raptorq.Config `mapstructure:"raptorq" json:"raptorq,omitempty"`
}

func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// New returns a new Config instance
func New() *Config {
	return &Config{
		LogLevel:        defaultLogLevel,
		LogCompress:     defaultLogCompress,
		LogMaxAgeInDays: defaultLogMaxAgeInDays,
		LogMaxBackups:   defaultLogMaxBackups,
		LogMaxSizeInMB:  defaultLogMaxSizeInMB,

		Node:    NewNode(),
		Pastel:  pastel.NewConfig(),
		RaptorQ: raptorq.NewConfig(),
	}
}
