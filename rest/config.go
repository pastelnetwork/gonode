package rest

const (
	defaultHostname = "localhost"
	defaultPort     = 8080
)

// Config contains settings of the Pastel client.
type Config struct {
	Hostname string `mapstructure:"hostname" json:"hostname"`
	Port     int    `mapstructure:"port" json:"port"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
