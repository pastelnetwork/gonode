package metadb

import (
	"fmt"
	"path/filepath"
)

const (
	errValidationStr = "metadb configs validation failed - missing val"
)

// Config contains settings of the rqlite server
type Config struct {
	// Let this instance be the leader
	IsLeader bool `mapstructure:"is_leader" json:"is_leader,omitempty"`

	// None voter node
	NoneVoter bool `mapstructure:"none_voter" json:"none_voter,omitempty"`

	// IPv4 or IPv6 address for listening by http server and raft
	ListenAddress string `mapstructure:"listen_address" json:"listen_address,omitempty"`

	// http server bind port. For HTTPS, set X.509 cert and key
	HTTPPort int `mapstructure:"http_port" json:"http_port,omitempty"`

	// raft port communication bind address
	RaftPort int `mapstructure:"raft_port" json:"raft_port,omitempty"`

	// data directory for rqlite
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`

	// BlockInterval is the block interval after which leadership is transferred
	BlockInterval int32 `mapstructure:"block_interval" json:"block_interval,omitempty"`
}

// SetWorkDir applies `workDir` to DataDir if it was not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.DataDir) {
		config.DataDir = filepath.Join(workDir, config.DataDir)
	}
}

// Validate metadb configs
func (config *Config) Validate() error {
	if config.ListenAddress == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "listen_address")
	}
	if config.HTTPPort == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "http_port")
	}
	if config.RaftPort == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "raft_port")
	}
	if config.DataDir == "" {
		return fmt.Errorf("%s: %s", errValidationStr, "data_dir")
	}
	if config.BlockInterval == 0 {
		return fmt.Errorf("%s: %s", errValidationStr, "block_interval")
	}

	return nil
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		BlockInterval: 48,
	}
}
