package configs

import (
	"encoding/json"

	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/metadb"
	"github.com/pastelnetwork/gonode/metadb/database"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/supernode/debug"
	healthcheck_lib "github.com/pastelnetwork/gonode/supernode/healthcheck"
)

const (
	defaultLogMaxAgeInDays   = 2
	defaultLogMaxSizeInMB    = 10
	defaultLogCompress       = true
	defaultLogMaxBackups     = 20
	defaultCommonLogLevel    = "info"
	defaultSubSystemLogLevel = "error" // disable almost logs
)

// Config contains configuration of all components of the SuperNode.
type Config struct {
	DefaultDir string     `json:"-"`
	LogConfig  *LogConfig `mapstructure:"log-config" json:"log-config,omitempty"`
	Quiet      bool       `mapstructure:"quiet" json:"quiet"`
	TempDir    string     `mapstructure:"temp-dir" json:"temp-dir"`
	WorkDir    string     `mapstructure:"work-dir" json:"work-dir"`
	RqFilesDir string     `mapstructure:"rq-files-dir" json:"rq-files-dir"`
	DdWorkDir  string     `mapstructure:"dd-service-dir" json:"dd-service-dir"`

	Node         `mapstructure:"node" json:"node,omitempty"`
	Pastel       *pastel.Config          `mapstructure:"-" json:"-"`
	P2P          *p2p.Config             `mapstructure:"p2p" json:"p2p,omitempty"`
	MetaDB       *metadb.Config          `mapstructure:"metadb" json:"metadb,omitempty"`
	UserDB       *database.Config        `mapstructure:"userdb" json:"userdb,omitempty"`
	DDServer     *ddclient.Config        `mapstructure:"dd-server" json:"dd-server,omitempty"`
	RaptorQ      *raptorq.Config         `mapstructure:"raptorq" json:"raptorq,omitempty"`
	HealthCheck  *healthcheck_lib.Config `mapstructure:"health-check" json:"health-check,omitempty"`
	DebugService *debug.Config           `mapstructure:"debug-service" json:"debug-service,omitempty"`
}

// LogConfig contains log configs
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
		Node:         NewNode(),
		Pastel:       pastel.NewConfig(),
		P2P:          p2p.NewConfig(),
		MetaDB:       metadb.NewConfig(),
		UserDB:       database.NewConfig(),
		RaptorQ:      raptorq.NewConfig(),
		DDServer:     ddclient.NewConfig(),
		HealthCheck:  healthcheck_lib.NewConfig(),
		DebugService: debug.NewConfig(),
	}
}
