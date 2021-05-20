package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	API             api.Config             `mapstructure:"api" json:"api,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() *Node {
	return &Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		API:             *api.NewConfig(),
	}
}
