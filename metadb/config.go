package metadb

const (
	defaultHTTPAddress = "0.0.0.0:4001"
	defaultRaftAddress = "0.0.0.0:4002"
	defaultDataDir     = ".data"
)

// Config contains settings of the rqlite server
type Config struct {
         // IPv4 or IPv6 address for listening by http server and raft
	ListenAddress string `mapstructure:"listen_address" json:"listen_address,omitempty"`
	
	// http server bind port. For HTTPS, set X.509 cert and key
	HTTPPort string `mapstructure:"http_port" json:"http_port,omitempty"`
	
	// raft port communication bind address
	RaftPort string `mapstructure:"raft_port" json:"raft_port,omitempty"`

	// data directory for rqlite
	DataDir string `mapstructure:"data_dir" json:"data_dir,omitempty"`		
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		HTTPAddress: defaultHTTPAddress,
		RaftAddress: defaultRaftAddress,
		DataDir:     defaultDataDir,
	}
}
