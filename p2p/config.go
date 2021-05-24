package p2p

var (
	defaultListenAddress = "0.0.0.0"
	defaultPort          = "0"
	defaultBootstrapIP   = ""
	defaultBootstrapPort = ""
	defualtUseStun       = true
)

// Config contains settings of the supernode server.
type Config struct {
	// The local IPv4 or IPv6 address
	ListenAddress string `mapstructure:"listen_address" json:"listen_address,omitempty"`
	// The local port to listen for connections on
	Port string `mapstructure:"port" json:"port,omitempty"`
	// IP Address to bootstrap against
	BootstrapIP string `mapstructure:"bootstrap_ip" json:"bootstrap_ip,omitempty"`
	// Port to bootstrap against
	BootstrapPort string `mapstructure:"bootstrap_port" json:"bootstrap_port,omitempty"`
	// Use STUN protocol for public addr discovery
	UseStun bool `mapstructure:"use_stun" json:"use_stun,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ListenAddress: defaultListenAddress,
		Port:          defaultPort,
		BootstrapIP:   defaultBootstrapIP,
		BootstrapPort: defaultBootstrapPort,
		UseStun:       defualtUseStun,
	}
}
