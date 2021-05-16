package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	servicescommon "github.com/pastelnetwork/gonode/walletnode/services/common"
)

// NodeServices contains configuration of the servcies.
type NodeServices struct {
	servicescommon.Config `mapstructure:",squash"`

	ArtworkRegister *artworkregister.Config `mapstructure:"artwork_register" json:"artwork_register,omitempty"`
}

// Node contains the SuperNode configuration itself.
type Node struct {
	NodeServices `mapstructure:",squash"`

	API *api.Config `mapstructure:"api" json:"api,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() *Node {
	node := &Node{
		NodeServices: NodeServices{
			Config:          *servicescommon.NewConfig(),
			ArtworkRegister: artworkregister.NewConfig(),
		},
		API: api.NewConfig(),
	}
	node.ArtworkRegister.Config = &node.Config

	return node
}
