package config

// PastelAPI contains configuration to call cNode RPC
type PastelAPI struct {
	Hostname   string `mapstructure:"hostname" json:"hostname"`
	Port       int    `mapstructure:"port" json:"port"`
	Username   string `mapstructure:"username" json:"username"`
	Passphrase string `mapstructure:"passphrase" json:"passphrase"`
}
