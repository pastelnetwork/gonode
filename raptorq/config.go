package raptorq

// Config contains settings of the p2p service
type Config struct {
	// the local IPv4 or IPv6 address
	Host string `mapstructure:"host" json:"host,omitempty"`

	// the local port to listen for connections on
	Port int `mapstructure:"port" json:"port,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{}
}
