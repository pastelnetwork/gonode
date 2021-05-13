package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	servicescommon "github.com/pastelnetwork/gonode/supernode/services/common"
)

// NodeServices contains configuration of the servcies.
type NodeServices struct {
	servicescommon.Config `mapstructure:",squash"`

	ArtworkRegister *artworkregister.Config `mapstructure:"artwork_register" json:"artwork_register,omitempty"`
}

// Node contains the SuperNode configuration itself.
type Node struct {
	NodeServices `mapstructure:",squash"`

	Server *server.Config `mapstructure:"server" json:"server,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() *Node {
	node := &Node{
		NodeServices: NodeServices{
			Config:          *servicescommon.NewConfig(),
			ArtworkRegister: artworkregister.NewConfig(),
		},
		Server: server.NewConfig(),
	}
	node.ArtworkRegister.Config = &node.Config

	return node
}
