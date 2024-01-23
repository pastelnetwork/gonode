package pastel

import (
	json "github.com/json-iterator/go"
)

const (
	defaultHostname          = "localhost"
	defaultMainnetPort       = 9932
	defaultTestnetPort       = 19932
	defaultBurnAddressTest   = "tPpasteLBurnAddressXXXXXXXXXXX3wy7u"
	defaultBurnAddressMain   = "PtpasteLBurnAddressXXXXXXXXXXbJ5ndd"
	defaultBurnAddressDevnet = "44oUgmZSL997veFEQDq569wv5tsT6KXf9QY7"
	defaultDevnetPort        = 29932
)

// Config Represents the structure of the `pastel.conf` file.
type Config struct {
	Hostname string `mapstructure:"rpcconnect"`
	Port     int    `mapstructure:"rpcport"`
	Username string `mapstructure:"rpcuser"`
	Password string `mapstructure:"rpcpassword"`
	Testnet  int    `mapstructure:"testnet"`
	Devnet   int    `mapstructure:"devnet"`
}

// MarshalJSON returns the JSON encoding.
func (config *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Hostname string `json:"hostname,omitempty"`
		Port     int    `json:"port,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"-"`
	}{
		Hostname: config.Hostname,
		Port:     config.port(),
		Username: config.Username,
		Password: config.Password,
	})
}

// Port returns port from external config.
// if port is not provided & testnet=1 it uses default testnet port
// else it uses default mainnet port
func (config *Config) port() int {
	if config.Port != 0 {
		return config.Port
	}

	if config.Testnet == 1 {
		return defaultTestnetPort
	}

	if config.Devnet == 1 {
		return defaultDevnetPort
	}

	return defaultMainnetPort
}

// BurnAddress ...
func (config *Config) BurnAddress() string {
	if config.Testnet == 1 {
		return defaultBurnAddressTest
	}

	if config.Devnet == 1 {
		return defaultBurnAddressDevnet
	}

	return defaultBurnAddressMain
}

// NewConfig returns a new Config instance.
func NewConfig() *Config {
	return &Config{
		Hostname: defaultHostname,
	}
}
