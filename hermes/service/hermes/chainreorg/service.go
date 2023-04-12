package chainreorg

import (
	svc "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
)

type chainReorgService struct {
	pastelClient pastel.Client
	store        store.PastelBlockStore
	sync         *synchronizer.Synchronizer
}

// NewChainReorgService returns a new chain-reorg service
func NewChainReorgService(store store.PastelBlockStore, pastelClient pastel.Client, s *synchronizer.Synchronizer) (svc.SvcInterface, error) {
	return &chainReorgService{
		pastelClient: pastelClient,
		store:        store,
		sync:         s,
	}, nil
}
