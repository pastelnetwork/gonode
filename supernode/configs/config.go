package configs

import (
	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/raptorq"
	"github.com/pastelnetwork/gonode/supernode/debug"
	healthcheck_lib "github.com/pastelnetwork/gonode/supernode/healthcheck"
)

const (
	defaultLogMaxAgeInDays   = 3
	defaultLogMaxSizeInMB    = 100
	defaultLogCompress       = true
	defaultLogMaxBackups     = 10
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

	Node   `mapstructure:"node" json:"node,omitempty"`
	Pastel *pastel.Config `mapstructure:"-" json:"-"`
	P2P    *p2p.Config    `mapstructure:"p2p" json:"p2p,omitempty"`
	//MetaDB       *metadb.Config          `mapstructure:"metadb" json:"metadb,omitempty"`
	//UserDB       *database.Config        `mapstructure:"userdb" json:"userdb,omitempty"`
	DDServer            *ddclient.Config        `mapstructure:"dd-server" json:"dd-server,omitempty"`
	RaptorQ             *raptorq.Config         `mapstructure:"raptorq" json:"raptorq,omitempty"`
	HealthCheck         *healthcheck_lib.Config `mapstructure:"health-check" json:"health-check,omitempty"`
	DebugService        *debug.Config           `mapstructure:"debug-service" json:"debug-service,omitempty"`
	RcloneStorageConfig *RcloneStorageConfig    `mapstructure:"storage-rclone" json:"storage-rclone,omitempty"`
}

// RcloneStorageConfig contains settings for RcloneStorage
type RcloneStorageConfig struct {
	SpecName   string `mapstructure:"spec_name" json:"spec_name,omitempty"`
	BucketName string `mapstructure:"bucket_name" json:"bucket_name,omitempty"`
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
		Node:   NewNode(),
		Pastel: pastel.NewConfig(),
		P2P:    p2p.NewConfig(),
		//MetaDB:       metadb.NewConfig(),
		//UserDB:       database.NewConfig(),
		RaptorQ:      raptorq.NewConfig(),
		DDServer:     ddclient.NewConfig(),
		HealthCheck:  healthcheck_lib.NewConfig(),
		DebugService: debug.NewConfig(),
	}
}

// OverridePastelIDAndPass overrides pastel id and pass phrase for all components
func (config *Config) OverridePastelIDAndPass(id, pass string) {
	// Override pastel id
	config.PastelID = id
	config.CascadeRegister.PastelID = id
	config.CollectionRegister.PastelID = id
	config.NftRegister.PastelID = id
	config.SenseRegister.PastelID = id
	config.NftDownload.PastelID = id
	config.StorageChallenge.PastelID = id
	config.SelfHealingChallenge.PastelID = id
	config.HealthCheckChallenge.PastelID = id

	// Override pass phrase
	config.PassPhrase = pass
	config.CascadeRegister.PassPhrase = pass
	config.CollectionRegister.PassPhrase = pass
	config.NftRegister.PassPhrase = pass
	config.SenseRegister.PassPhrase = pass
	config.NftDownload.PassPhrase = pass
	config.StorageChallenge.PassPhrase = pass
	config.SelfHealingChallenge.PassPhrase = pass
	config.HealthCheckChallenge.PassPhrase = pass
}
