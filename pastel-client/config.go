package pastel

const (
	defaultHostname = "localhost"
	defaultPort     = 9932
)

// Config contains settings of the Pastel client.
type Config struct {
	Hostname string `mapstructure:"hostname" json:"hostname,omitempty"`
	Port     int    `mapstructure:"port" json:"port,omitempty"`
	Username string `mapstructure:"username" json:"username,omitempty"`
	Password string `mapstructure:"password" json:"-"`

	ExtKey string `mapstructure:"ext_key" json:"ext_key,omitempty"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
