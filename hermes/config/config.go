package config

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	defaultLogMaxAgeInDays   = 3
	defaultLogMaxSizeInMB    = 100
	defaultLogCompress       = true
	defaultLogMaxBackups     = 10
	defaultCommonLogLevel    = "info"
	defaultSubSystemLogLevel = "error" // disable almost logs

)

// Config contains settings of the hermes server
type Config struct {
	LogConfig  *LogConfig `mapstructure:"log-config" json:"log-config,omitempty"`
	Quiet      bool       `mapstructure:"quiet" json:"quiet"`
	DdWorkDir  string     `mapstructure:"dd-service-dir" json:"dd-service-dir"`
	PastelID   string     `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
	PassPhrase string     `mapstructure:"pass_phrase" json:"pass_phrase,omitempty"`

	// ActivationDeadline is number of blocks allowed to pass for a ticket to be activated
	ActivationDeadline uint32         `mapstructure:"activation_deadline" json:"activation_deadline,omitempty"`
	WorkDir            string         `mapstructure:"work-dir" json:"work-dir"`
	Pastel             *pastel.Config `mapstructure:"-" json:"-"`
	P2P                *p2p.Config    `mapstructure:"p2p" json:"p2p,omitempty"`
}

// LogConfig represents log config
type LogConfig struct {
	File         string          `mapstructure:"log-file" json:"log-file,omitempty"`
	Compress     bool            `mapstructure:"log-compress" json:"log-compress,omitempty"`
	MaxSizeInMB  int             `mapstructure:"log-max-size-mb" json:"log-max-size-mb,omitempty"`
	MaxAgeInDays int             `mapstructure:"log-max-age-days" json:"log-max-age-days,omitempty"`
	MaxBackups   int             `mapstructure:"log-max-backups" json:"log-max-backups,omitempty"`
	Levels       *LogLevelConfig `mapstructure:"log-levels" json:"log-levels,omitempty"`
}

// LogLevelConfig contains log level configs for each subsystem
type LogLevelConfig struct {
	CommonLogLevel        string `mapstructure:"common" json:"common,omitempty"`
	P2PLogLevel           string `mapstructure:"p2p" json:"p2p,omitempty"`
	MetaDBLogLevel        string `mapstructure:"metadb" json:"metadb,omitempty"`
	DupeDetectionLogLevel string `mapstructure:"dd" json:"dd,omitempty"`
}

// String converts to string
func (config *Config) String() string {
	// The main purpose of using a custom converting is to avoid unveiling credentials.
	// All credentials fields must be tagged `json:"-"`.
	data, _ := json.Marshal(config)
	return string(data)
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		LogConfig: &LogConfig{
			Compress:     defaultLogCompress,
			MaxSizeInMB:  defaultLogMaxSizeInMB,
			MaxAgeInDays: defaultLogMaxAgeInDays,
			MaxBackups:   defaultLogMaxBackups,
			Levels: &LogLevelConfig{
				CommonLogLevel:        defaultCommonLogLevel,
				MetaDBLogLevel:        defaultSubSystemLogLevel,
				P2PLogLevel:           defaultSubSystemLogLevel,
				DupeDetectionLogLevel: defaultSubSystemLogLevel,
			},
		},
		Pastel: pastel.NewConfig(),
		P2P:    p2p.NewConfig(),
	}
}
