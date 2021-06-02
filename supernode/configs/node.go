package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/rqlite"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	Server          *server.Config          `mapstructure:"server" json:"server,omitempty"`
	RQLite          *rqlite.Config         `mapstructure:"rqlite" json:"rqlite,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		Server:          server.NewConfig(),
		RQLite:          rqlite.NewConfig(),
	}
}
