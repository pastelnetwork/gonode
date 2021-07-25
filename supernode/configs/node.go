package configs

import (
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister            artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	Server                     *server.Config         `mapstructure:"server" json:"server,omitempty"`
	PastelID                   string                 `mapstructure:"pastel_id" json:"pastel_id,omitempty"`
	PassPhrase                 string                 `mapstructure:"pass_phrase" json:"pass_phrase,omitempty"`
	PreburntTxMinConfirmations int                    `mapstructure:"preburnt_tx_min_confirmation" json:"preburnt_tx_min_confirmation,omitempty"`
	// timeout in minute
	PreburntTxConfirmationTimeout int `mapstructure:"preburnt_tx_confirmation_timeout" json:"preburnt_tx_confirmation_timeout,omitempty"`
}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		Server:          server.NewConfig(),
	}
}
