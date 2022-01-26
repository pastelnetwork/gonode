package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/nftdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/nftregister"
	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister nftregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkSearch   nftsearch.Config   `mapstructure:",squash" json:"artwork_search,omitempty"`
	ArtworkDownload nftdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	API             *api.Config        `mapstructure:"api" json:"api,omitempty"`

	SenseRegister senseregister.Config `mapstructure:"sense_register" json:"sense_register,omitempty"`

	UserdataProcess userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`

	RegTxMinConfirmations int `mapstructure:"reg_tx_min_confirmations" json:"reg_tx_min_confirmation,omitempty"`
	ActTxMinConfirmations int `mapstructure:"act_tx_min_confirmations" json:"act_tx_min_confirmations,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkSearch:   *nftsearch.NewConfig(),
		ArtworkRegister: *nftregister.NewConfig(),
		ArtworkDownload: *nftdownload.NewConfig(),
		API:             api.NewConfig(),
		UserdataProcess: *userdataprocess.NewConfig(),
		SenseRegister:   *senseregister.NewConfig(),
	}
}
