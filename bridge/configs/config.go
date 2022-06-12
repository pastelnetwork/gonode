package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/bridge/services/server"

	"github.com/pastelnetwork/gonode/bridge/services/download"

	"github.com/pastelnetwork/gonode/pastel"
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
	LogConfig *LogConfig `mapstructure:"log-config" json:"log-config,omitempty"`
	Quiet     bool       `mapstructure:"quiet" json:"quiet"`
	TempDir   string     `mapstructure:"temp-dir" json:"temp-dir"`
	WorkDir   string     `mapstructure:"work-dir" json:"work-dir"`

	Pastel   *pastel.Config   `mapstructure:"-" json:"-"`
	Download *download.Config `mapstructure:"download" json:"download,omitempty"`
	Server   *server.Config   `mapstructure:"server" json:"server,omitempty"`
}

// LogConfig contains log configs
type LogConfig struct {
	Level        string `mapstructure:"log-level" json:"log-level,omitempty"`
	File         string `mapstructure:"log-file" json:"log-file,omitempty"`
	Compress     bool   `mapstructure:"log-compress" json:"log-compress,omitempty"`
	MaxSizeInMB  int    `mapstructure:"log-max-size-mb" json:"log-max-size-mb,omitempty"`
	MaxAgeInDays int    `mapstructure:"log-max-age-days" json:"log-max-age-days,omitempty"`
	MaxBackups   int    `mapstructure:"log-max-backups" json:"log-max-backups,omitempty"`
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
		LogConfig: &LogConfig{
			Level:        defaultLogLevel,
			Compress:     defaultLogCompress,
			MaxAgeInDays: defaultLogMaxAgeInDays,
			MaxBackups:   defaultLogMaxBackups,
			MaxSizeInMB:  defaultLogMaxSizeInMB,
		},

		Pastel:   pastel.NewConfig(),
		Download: download.NewConfig(),
		Server:   server.NewConfig(),
	}
}
