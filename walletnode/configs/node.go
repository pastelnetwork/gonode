package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkSearch   artworksearch.Config   `mapstructure:",squash" json:"artwork_search,omitempty"`
	ArtworkDownload artworkdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	API             *api.Config            `mapstructure:"api" json:"api,omitempty"`

	SenseRegister senseregister.Config `mapstructure:"sense_register" json:"sense_register,omitempty"`

	UserdataProcess userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`
	BurnAddress     string                 `mapstructure:"burn_address" json:"burn_address,omitempty"`

	RegTxMinConfirmations int `mapstructure:"reg_tx_min_confirmations" json:"reg_tx_min_confirmation,omitempty"`
	ActTxMinConfirmations int `mapstructure:"act_tx_min_confirmations" json:"act_tx_min_confirmations,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkSearch:   *artworksearch.NewConfig(),
		ArtworkRegister: *artworkregister.NewConfig(),
		ArtworkDownload: *artworkdownload.NewConfig(),
		API:             api.NewConfig(),
		UserdataProcess: *userdataprocess.NewConfig(),
	}
}
