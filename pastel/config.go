package pastel

const (
	defaultHostname = "localhost"
)

// Config contains settings of the Pastel client.
type Config struct {
	DataDir  string `mapstructure:"data-dir" json:"data-dir"`
	Hostname string `mapstructure:"hostname" json:"hostname"`
	Port     int    `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"-"`
}

// NewConfig returns a new Config instance
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
	}
}
