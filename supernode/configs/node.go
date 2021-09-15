package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkdownload"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	Server          *server.Config         `mapstructure:"server" json:"server,omitempty"`
	PastelID        string                 `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
	PassPhrase      string                 `mapstructure:"pass_phrase" json:"pass_phrase,omitempty"`

	NumberConnectedNodes       int `mapstructure:"number_connected_nodes" json:"number_connected_nodes,omitempty"`
	PreburntTxMinConfirmations int `mapstructure:"preburnt_tx_min_confirmations" json:"preburnt_tx_min_confirmations,omitempty"`
	// timeout in minute
	PreburntTxConfirmationTimeout int                    `mapstructure:"preburnt_tx_confirmation_timeout" json:"preburnt_tx_confirmation_timeout,omitempty"`
	ArtworkDownload               artworkdownload.Config `mapstructure:",squash" json:"artwork_download,omitempty"`
	UserdataProcess               userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		ArtworkDownload: *artworkdownload.NewConfig(),
		UserdataProcess: *userdataprocess.NewConfig(),
		Server:          server.NewConfig(),
	}
}
