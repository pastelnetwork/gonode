package metadb

import (
	"fmt"
	"path/filepath"
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
}

// SetWorkDir applies `workDir` to DataDir if it was not specified as an absolute path.
func (config *Config) SetWorkDir(workDir string) {
	if !filepath.IsAbs(config.DataDir) {
		config.DataDir = filepath.Join(workDir, config.DataDir)
	}
}

// GetExposedAddr returns IPv4 or IPv6 Addr along with port
func (config *Config) GetExposedAddr() string {
	return fmt.Sprintf("%s:%v", config.ListenAddress, config.HTTPPort)
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}
