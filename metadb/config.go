package metadb

const (
	defaultHTTPAddress  = "0.0.0.0:4001"
	defaultRaftAddress  = "0.0.0.0:4002"
	defaultRaftLogLevel = "error"
	defaultDataDir      = ".data"
)

// Config contains settings of the rqlite server
type Config struct {
	// http server bind address. For HTTPS, set X.509 cert and key
	HTTPAddress string `mapstructure:"http_address" json:"http_address,omitempty"`
	// unique name for node. If not set, set to hostname
	NodeID string `mapstructure:"node_id" json:"node_id,omitempty"`
	// raft communication bind address
	RaftAddress string `mapstructure:"raft_address" json:"raft_address,omitempty"`
	// comma-delimited list of nodes, through which a cluster can be joined (proto://host:port)"
	JoinAddress string `mapstructure:"join_address" json:"join_address,omitempty"`
	// set Discovery Service URL
	DiscoveryURL string `mapstructure:"discovery_url" json:"discovery_url,omitempty"`
	// set Discovery ID. If not set, Discovery Service not used
	DiscoveryID string `mapstructure:"discovery_id" json:"discovery_id,omitempty"`
	// minimum log level for Raft module
	RaftLogLevel string `mapstructure:"raft_log_level" json:"raft_log_level,omitempty"`
	// data directory for rqlite
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		HTTPAddress:  defaultHTTPAddress,
		RaftAddress:  defaultRaftAddress,
		RaftLogLevel: defaultRaftLogLevel,
		DataDir:      defaultDataDir,
	}
}
