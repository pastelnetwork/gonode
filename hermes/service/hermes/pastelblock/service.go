package pastelblock

import (
	"github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
)

type pastelBlockService struct {
	pastelClient pastel.Client
	store        store.PastelBlockStore
	sync         *synchronizer.Synchronizer
}

// NewPastelBlocksService returns a new pastel-block service
func NewPastelBlocksService(pbStore store.PastelBlockStore, pastelClient pastel.Client, s *synchronizer.Synchronizer) (service.SvcInterface, error) {
	return &pastelBlockService{
		pastelClient: pastelClient,
		store:        pbStore,
		sync:         s,
	}, nil
}
