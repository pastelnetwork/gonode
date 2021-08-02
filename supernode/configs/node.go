package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	ArtworkDownload artworkdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	Server          *server.Config         `mapstructure:"server" json:"server,omitempty"`
	PastelID        string                 `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		Server:          server.NewConfig(),
	}
}
