package api

const (
	defaultHostname = ""
	defaultPort     = 29932
)

type Config struct {
	Hostname string `mapstructure:"hostname" json:"hostname,omitempty"`
	Port     int    `mapstructure:"port" json:"port,omitempty"`
	Username string `mapstructure:"username" json:"username,omitempty"`
	Password string `mapstructure:"password" json:"password,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
		Port:     defaultPort,
	}
}
