package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkDownload artworkdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	API             *api.Config            `mapstructure:"api" json:"api,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		ArtworkDownload: *artworkdownload.NewConfig(),
		API:             api.NewConfig(),
	}
}
