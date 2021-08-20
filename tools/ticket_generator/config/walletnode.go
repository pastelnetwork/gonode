package config

// RestAPI contains information of the walletnode endpoint
// Used when register nft to the blockchain
type RestAPI struct {
	Hostname string `mapstructure:"hostname" json:"hostname"`
	Port     int    `mapstructure:"port" json:"port"`
}

// Artist contains registered information on blockchain
// Used when register nft to the block chain
type Artist struct {
	PastelID      string `mapstructure:"pastel_id" json:"pastel_id"`
	Passphrase    string `mapstructure:"passphrase" json:"passphrase"`
	SpendableAddr string `mapstructure:"spendable_addr" json:"spendable_addr"`
}

// WalletNode contains information to interact with associated cNode
// of go walletnode
type WalletNode struct {
	PastelAPI PastelAPI `mapstructure:"pastel_api" json:"pastel_api"`
	RestAPI   RestAPI   `mapstructure:"rest_api" json:"rest_api"`
	Artist    Artist    `mapstructure:"artist" json:"artist"`
}
