package chainreorg

import (
	"github.com/pastelnetwork/gonode/hermes/common"
	svc "github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
)

type chainReorgService struct {
	pastelClient pastel.Client
	store        store.PastelBlockStore
	sync         *synchronizer.Synchronizer
	blockMsgChan chan<- common.BlockMessage
}

// NewChainReorgService returns a new chain-reorg service
func NewChainReorgService(store store.PastelBlockStore, pastelClient pastel.Client, blockMsgChan chan<- common.BlockMessage) (svc.SvcInterface, error) {
	return &chainReorgService{
		pastelClient: pastelClient,
		store:        store,
		sync:         synchronizer.NewSynchronizer(pastelClient),
		blockMsgChan: blockMsgChan,
	}, nil
}
