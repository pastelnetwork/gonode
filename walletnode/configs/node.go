package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkSearch   artworksearch.Config   `mapstructure:",squash" json:"artwork_search,omitempty"`
	ArtworkDownload artworkdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	API             *api.Config            `mapstructure:"api" json:"api,omitempty"`
	BurnAddress     string                 `mapstructure:"burn_address" json:"burn_address,omitempty"`

	RegArtTxMinConfirmations int `mapstructure:"reg_art_tx_min_confirmations" json:"reg_art_tx_min_confirmation,omitempty"`
	// Timeout in minutes
	RegArtTxTimeout int `mapstructure:"reg_art_tx_timeout" json:"reg_art_tx_timeout,omitempty"`

	RegActTxMinConfirmations int `mapstructure:"reg_act_tx_min_confirmations" json:"reg_act_tx_min_confirmations,omitempty"`
	// Timeout in minutes
	RegActTxTimeout int `mapstructure:"reg_act_tx_timeout" json:"reg_act_tx_timeout,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkSearch:   *artworksearch.NewConfig(),
		ArtworkRegister: *artworkregister.NewConfig(),
		ArtworkDownload: *artworkdownload.NewConfig(),
		API:             api.NewConfig(),
	}
}
