package config

// MinerNode contains configuration details to call miner RPC
type MinerNode struct {
	PastelAPI PastelAPI `json:"pastel_api" mapstructure:"pastel_api"`
}
