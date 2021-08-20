package config

// MinderNode contains configuration details to call miner RPC
type MinderNode struct {
	PastelAPI PastelAPI `json:"pastel_api"`
}
