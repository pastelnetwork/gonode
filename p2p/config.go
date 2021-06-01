package p2p

var (
	defaultListenAddress = "0.0.0.0"
	defaultPort          = 4445
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
	// a driver-specific data source name, usually consisting of
	// at least a database name and connection information.
	DataSourceName string `mapstructure:"data_source_name" json:"data_source_name,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ListenAddresses: defaultListenAddress,
		Port:            defaultPort,
	}
}
