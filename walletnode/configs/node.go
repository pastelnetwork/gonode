package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkSearch   artworksearch.Config   `mapstructure:",squash" json:"artwork_search,omitempty"`
	API             *api.Config            `mapstructure:"api" json:"api,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkSearch:   *artworksearch.NewConfig(),
		ArtworkRegister: *artworkregister.NewConfig(),
		API:             api.NewConfig(),
	}
}
