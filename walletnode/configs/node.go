package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/cascaderegister"
	"github.com/pastelnetwork/gonode/walletnode/services/collectionregister"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
	"github.com/pastelnetwork/gonode/walletnode/services/nftregister"
	"github.com/pastelnetwork/gonode/walletnode/services/nftsearch"
	"github.com/pastelnetwork/gonode/walletnode/services/senseregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	API *api.Config `mapstructure:"api" json:"api,omitempty"`

	NftRegister        nftregister.Config        `mapstructure:",squash" json:"nft_register,omitempty"`
	NftSearch          nftsearch.Config          `mapstructure:",squash" json:"nft_search,omitempty"`
	NftDownload        download.Config           `mapstructure:",squash" json:"nft_download,omitempty"`
	SenseRegister      senseregister.Config      `mapstructure:"sense_register" json:"sense_register,omitempty"`
	CascadeRegister    cascaderegister.Config    `mapstructure:"cascade_register" json:"cascade_register,omitempty"`
	CollectionRegister collectionregister.Config `mapstructure:"collection_register" json:"collection_register,omitempty"`

	// UserdataProcess userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`

	RegTxMinConfirmations int `mapstructure:"reg_tx_min_confirmations" json:"reg_tx_min_confirmation,omitempty"`
	ActTxMinConfirmations int `mapstructure:"act_tx_min_confirmations" json:"act_tx_min_confirmations,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		NftSearch:   *nftsearch.NewConfig(),
		NftRegister: *nftregister.NewConfig(),
		NftDownload: *download.NewConfig(),
		API:         api.NewConfig(),

		SenseRegister:      *senseregister.NewConfig(),
		CascadeRegister:    *cascaderegister.NewConfig(),
		CollectionRegister: *collectionregister.NewConfig(),
	}
}
