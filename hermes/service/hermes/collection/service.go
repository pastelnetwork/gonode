package collection

import (
	"github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/store"

	"github.com/pastelnetwork/gonode/pastel"
)

type collectionService struct {
	pastelClient pastel.Client
	store        store.CollectionStore
	sync         *synchronizer.Synchronizer

	latestCollectionBlockHeight int
}

// NewCollectionService returns a new collection service
func NewCollectionService(cStore store.CollectionStore, pastelClient pastel.Client, s *synchronizer.Synchronizer) (service.SvcInterface, error) {
	return &collectionService{
		pastelClient: pastelClient,
		store:        cStore,
		sync:         s,
	}, nil
}
