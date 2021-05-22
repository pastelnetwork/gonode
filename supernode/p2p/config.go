package p2p

var (
	ip    = "0.0.0.0"
	port  = "0"
	bIP   = ""
	bPort = ""
	useStun  = true
)

// Config contains settings of the supernode server.
type Config struct {
	// The local IPv4 or IPv6 address
	IP string

	// The local port to listen for connections on
	Port string

	// IP Address to bootstrap against
	BootstrapIP string

	// Port to bootstrap against
	BootstrapPort string

	// Use STUN protocol for public addr discovery
	UseStun bool
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		IP:            ip,
		Port:          port,
		BootstrapIP:   bIP,
		BootstrapPort: bPort,
		UseStun:       useStun,
	}
}
