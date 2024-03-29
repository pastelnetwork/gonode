package cleaner

import (
	"github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/pastel"
)

type cleanupService struct {
	pastelClient pastel.Client
	p2p          node.HermesP2PInterface
	sync         *synchronizer.Synchronizer

	config Config

	currentNFTBlock    int
	currentActionBlock int
}

// NewCleanupService returns a new cleanup service
func NewCleanupService(pastelClient pastel.Client, hp2p node.HermesP2PInterface, config Config) (service.SvcInterface, error) {
	return &cleanupService{
		pastelClient:       pastelClient,
		p2p:                hp2p,
		currentNFTBlock:    1,
		currentActionBlock: 1,
		sync:               synchronizer.NewSynchronizer(pastelClient),
		config:             config,
	}, nil
}
