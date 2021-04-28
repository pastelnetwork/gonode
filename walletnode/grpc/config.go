package grpc

const (
	defaultPort = 4444
)

// Config contains settings of the grpc client.
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
