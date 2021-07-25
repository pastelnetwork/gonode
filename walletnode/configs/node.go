package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister          artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	API                      *api.Config            `mapstructure:"api" json:"api,omitempty"`
	BurnAddress              string                 `mapstructure:"burn_address" json:"burn_address,omitempty"`
	RegArtTxMinConfirmations int                    `mapstructure:"reg_art_tx_min_confirmation" json:"reg_art_tx_min_confirmation,omitempty"`
	// Timeout in minutes
	RegArtTxTimeout int `mapstructure:"reg_art_tx_timeout" json:"reg_art_tx_timeout,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		API:             api.NewConfig(),
	}
}
