package fingerprint

import (
	"github.com/pastelnetwork/gonode/hermes/common"
	"github.com/pastelnetwork/gonode/hermes/service"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/synchronizer"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/hermes/store"
	"github.com/pastelnetwork/gonode/pastel"
)

type fingerprintService struct {
	pastelClient pastel.Client
	store        store.DDStore
	p2p          node.HermesP2PInterface
	sync         *synchronizer.Synchronizer

	config Config

	latestNFTBlockHeight   int
	latestSenseBlockHeight int

	blockMsgChan <-chan common.BlockMessage
}

// NewFingerprintService returns a new fingerprint service
func NewFingerprintService(fgStore store.DDStore, pastelClient pastel.Client, hp2p node.HermesP2PInterface, config Config, blockMsgChan <-chan common.BlockMessage) (service.SvcInterface, error) {
	return &fingerprintService{
		pastelClient: pastelClient,
		store:        fgStore,
		p2p:          hp2p,
		config:       config,
		sync:         synchronizer.NewSynchronizer(pastelClient),
		blockMsgChan: blockMsgChan,
	}, nil
}
