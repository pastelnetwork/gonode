package models

// MasterNodeListConf represents multiple configurations, that can be retrieved using the command `masternode list-conf`.
type MasterNodeListConf map[string]MasterNodeConfig

// MasterNodeConfig represents pastel masternode configuration.
type MasterNodeConfig struct {
	Alias       string `json:"alias"`
	Address     string `json:"address"`
	PrivateKey  string `json:"privateKey"`
	TXHash      string `json:"txHash"`
	OutputIndex string `json:"outputIndex"`
	ExtAddress  string `json:"extAddress"`
	ExtKey      string `json:"extKey"`
	ExtCfg      string `json:"extCfg"`
	Status      string `json:"status"`
}
