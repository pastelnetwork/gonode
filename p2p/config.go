package p2p

import "path/filepath"

const (
	defaultListenAddress = "0.0.0.0"
	defaultPort          = 4445
	defaultDataDir       = "p2p"
)

// Config contains settings of the supernode server.
type Config struct {
	// The local IPv4 or IPv6 address
	ListenAddresses string `mapstructure:"listen_addresses" json:"listen_address,omitempty"`

	// The local port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`

	// IP Address to bootstrap against
	BootstrapIP string `mapstructure:"bootstrap_ip" json:"bootstrap_ip,omitempty"`

	// Port to bootstrap against
	BootstrapPort string `mapstructure:"bootstrap_port" json:"bootstrap_port,omitempty"`

	// Use STUN protocol for public addr discovery
	UseStun bool `mapstructure:"use_stun" json:"use_stun,omitempty"`

	// Default value is assigned automatically, by concatenated
	// CLI parameter `--work-dir` and `p2p` as a subfolder
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`

	// if not specified, persistent database store is used
	MemoryDB bool `mapstructure:"memory_db" json:"memory_db,omitempty"`
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
		ListenAddresses: defaultListenAddress,
		Port:            defaultPort,
		DataDir:         defaultDataDir,
	}
}
