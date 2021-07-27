package configs

import (
	"github.com/pastelnetwork/gonode/walletnode/api"
	"github.com/pastelnetwork/gonode/walletnode/services/artworkregister"
	"github.com/pastelnetwork/gonode/walletnode/services/userdataprocess"

)

// Node contains the SuperNode configuration itself.
type Node struct {
	// `squash` field cannot be pointer
	ArtworkRegister artworkregister.Config `mapstructure:",squash" json:"artwork_register,omitempty"`
	API             *api.Config            `mapstructure:"api" json:"api,omitempty"`

	UserdataProcess userdataprocess.Config `mapstructure:",squash" json:"userdata_process,omitempty"`

}

// NewNode returns a new Node instance
func NewNode() Node {
	return Node{
		ArtworkRegister: *artworkregister.NewConfig(),
		API:             api.NewConfig(),
		UserdataProcess: *userdataprocess.NewConfig(),
	}
}
