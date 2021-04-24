package nats

const (
	defaultPort = 4222
)

// Config contains settings of the nats-server.
type Config struct {
	Hostname string `mapstructure:"hostname" json:"hostname,omitempty"`
	Port     int    `mapstructure:"port" json:"port,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Port: defaultPort,
	}
}
