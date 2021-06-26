package metadb

import (
	"path/filepath"
)

const (
	defaultListenAddress = "0.0.0.0"
	defaultHTTPPort      = 4446
	defaultRaftPort      = 4447
	defaultDataDir       = "metadb"
)

// Config contains settings of the rqlite server
type Config struct {
	// IPv4 or IPv6 address for listening by http server and raft
	ListenAddress string `mapstructure:"listen_address" json:"listen_address,omitempty"`

	// http server bind port. For HTTPS, set X.509 cert and key
	HTTPPort int `mapstructure:"http_port" json:"http_port,omitempty"`

	// raft port communication bind address
	RaftPort int `mapstructure:"raft_port" json:"raft_port,omitempty"`

	// data directory for rqlite
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`
}

// SetWorkDir applies `workDir` to DataDir if it was not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.DataDir) {
		config.DataDir = filepath.Join(workDir, config.DataDir)
	}
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ListenAddress: defaultListenAddress,
		HTTPPort:      defaultHTTPPort,
		RaftPort:      defaultRaftPort,
		DataDir:       defaultDataDir,
	}
}
