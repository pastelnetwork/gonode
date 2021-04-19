package pastel

const (
	defaultHostname = "localhost"
)

// Config contains settings of the Pastel client.
type Config struct {
	DataDir  string `mapstructure:"data-dir" json:"data-dir,omitempty"`
	Hostname string `mapstructure:"hostname" json:"hostname,omitempty"`
	Port     int    `mapstructure:"port" json:"port,omitempty"`
	Username string `mapstructure:"username" json:"username,omitempty"`
	Password string `mapstructure:"password" json:"-"`
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
	}
}
