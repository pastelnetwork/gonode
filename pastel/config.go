package pastel

const (
	defaultConfigFile = "pastel.conf"
	defaultHostname   = "localhost"
	defaultPort       = 9932
)

// ExternalConfig represents the structure of the `pastel.conf` file.
type ExternalConfig struct {
	Hostname string `mapstructure:"rpcconnect" json:"hostname,omitempty"`
	Port     int    `mapstructure:"rpcport" json:"port,omitempty"`
	Username string `mapstructure:"rpcuser" json:"username,omitempty"`
	Password string `mapstructure:"rpcpassword" json:"-"`
}

// Config contains settings of the Pastel client.
type Config struct {
	*ExternalConfig

	ConfigFile string `mapstructure:"config-file" json:"config-file,omitempty"`

	Node struct {
		Hostname *string `mapstructure:"hostname" json:"hostname,omitempty"`
		Port     *int    `mapstructure:"port" json:"port,omitempty"`
		Username *string `mapstructure:"username" json:"username,omitempty"`
		Password *string `mapstructure:"password" json:"-"`
	} `mapstructure:",squash"`
}

// Hostname returns node hostname if it is specified, otherwise returns hostname from external config.
func (config *Config) Hostname() string {
	if config.Node.Hostname != nil {
		return *config.Node.Hostname
	}
	return config.ExternalConfig.Hostname
}

// Port returns node port if it is specified, otherwise returns port from external config.
func (config *Config) Port() int {
	if config.Node.Port != nil {
		return *config.Node.Port
	}
	return config.ExternalConfig.Port
}

// Username returns username port if it is specified, otherwise returns username from external config.
func (config *Config) Username() string {
	if config.Node.Username != nil {
		return *config.Node.Username
	}
	return config.ExternalConfig.Username
}

// Password returns password port if it is specified, otherwise returns password from external config.
func (config *Config) Password() string {
	if config.Node.Password != nil {
		return *config.Node.Password
	}
	return config.ExternalConfig.Password
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ConfigFile: defaultConfigFile,
		ExternalConfig: &ExternalConfig{
			Hostname: defaultHostname,
			Port:     defaultPort,
		},
	}
}
