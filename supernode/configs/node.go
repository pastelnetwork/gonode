package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	mdlserver "github.com/pastelnetwork/gonode/metadb/network/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	UserdataProcess userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`
	Server          *server.Config         `mapstructure:"server" json:"server,omitempty"`
	MDLServer		*mdlserver.Config      `mapstructure:"mdlserver" json:"mdlserver,omitempty"`
	PastelID        string                 `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		UserdataProcess: *userdataprocess.NewConfig(),
		Server:          server.NewConfig(),
		MDLServer:       mdlserver.NewConfig(),
	}
}
