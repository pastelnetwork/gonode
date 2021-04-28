package server

const (
	defaultListenAddresses = "0.0.0.0"
	defaultPort            = 4444
)

// Config contains settings of the supernode server.
type Config struct {
	ListenAddresses string `mapstructure:"listen_addresses" json:"listen_addresses,omitempty"`
	Port            int    `mapstructure:"port" json:"port,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ListenAddresses: defaultListenAddresses,
		Port:            defaultPort,
	}
}
