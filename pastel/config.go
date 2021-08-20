package pastel

import (
	"encoding/json"
)

const (
	defaultHostname    = "localhost"
	defaultMainnetPort = 9932
	defaultTestnetPort = 19932
)

// ExternalConfig represents the structure of the `pastel.conf` file.
type ExternalConfig struct {
	Hostname string `mapstructure:"rpcconnect"`
	Port     int    `mapstructure:"rpcport"`
	Username string `mapstructure:"rpcuser"`
	Password string `mapstructure:"rpcpassword"`
	Testnet  int    `mapstructure:"testnet"`
}

// Config contains settings of the Pastel client.
type Config struct {
	*ExternalConfig

	Hostname *string `mapstructure:"hostname"`
	Port     *int    `mapstructure:"port"`
	Username *string `mapstructure:"username"`
	Password *string `mapstructure:"password"`
}

// MarshalJSON returns the JSON encoding.
func (config *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Hostname string `json:"hostname,omitempty"`
		Port     int    `json:"port,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"-"`
	}{
		Hostname: config.hostname(),
		Port:     config.port(),
		Username: config.username(),
		Password: config.password(),
	})
}

// Hostname returns node hostname if it is specified, otherwise returns hostname from external config.
func (config *Config) hostname() string {
	if config.Hostname != nil {
		return *config.Hostname
	}
	return config.ExternalConfig.Hostname
}

// Port returns node port if it is specified, otherwise returns port from external config.
func (config *Config) port() int {
	if config.Port != nil {
		return *config.Port
	}

	if config.ExternalConfig.Port != 0 {
		return config.ExternalConfig.Port
	}

	if config.ExternalConfig.Testnet == 1 {
		return defaultTestnetPort
	}

	return defaultMainnetPort
}

// Username returns username port if it is specified, otherwise returns username from external config.
func (config *Config) username() string {
	if config.Username != nil {
		return *config.Username
	}
	return config.ExternalConfig.Username
}

// Password returns password port if it is specified, otherwise returns password from external config.
func (config *Config) password() string {
	if config.Password != nil {
		return *config.Password
	}
	return config.ExternalConfig.Password
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		ExternalConfig: &ExternalConfig{
			Hostname: defaultHostname,
		},
	}
}
