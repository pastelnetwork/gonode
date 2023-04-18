package synchronizer

import (
	"github.com/pastelnetwork/gonode/pastel"
)

// Synchronizer is used to keep track of master-node sync status throughout hermes
type Synchronizer struct {
	//PastelClient to access cNode through client implementation
	PastelClient pastel.Client
}

// NewSynchronizer returns a new synchronizer
func NewSynchronizer(pastelClient pastel.Client) *Synchronizer {
	return &Synchronizer{
		PastelClient: pastelClient,
	}
}
