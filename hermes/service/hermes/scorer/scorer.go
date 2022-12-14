package scorer

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/store"

	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer/hash"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer/network"
	"github.com/pastelnetwork/gonode/hermes/service/hermes/scorer/thumbnail"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

//  Scorer is scorer service
type Scorer struct {
	thumbnailChlngr *thumbnail.ThumbnailsChallenger
	imgHashChlngr   *hash.ImageHashChallenger
}

// Start initiates scoring challenges
func (s *Scorer) Start(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)
	/*group.Go(func() error {
		return s.imgHashChlngr.Run(gctx)
	})*/

	group.Go(func() error {
		return s.thumbnailChlngr.Run(gctx)
	})

	if err := group.Wait(); err != nil {
		log.WithContext(gctx).WithError(err).Errorf("scorer service run failure")
	}

	return nil
}

// New creates a new scorer
func New(config *Config, pastelClient pastel.Client, sn node.SNClientInterface, store store.ScoreStore) *Scorer {
	meshHandlerOpts := network.MeshHandlerOpts{
		NodeMaker:     &network.RegisterNftNodeMaker{},
		PastelHandler: mixins.NewPastelHandler(pastelClient),
		NodeClient:    sn,
		Configs: &network.MeshHandlerConfig{
			ConnectToNextNodeDelay: config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     config.AcceptNodesTimeout,
			MinSNs:                 config.NumberSuperNodes,
			PastelID:               config.CreatorPastelID,
			Passphrase:             config.CreatorPastelIDPassphrase,
		},
		Store: store,
	}

	imgChlngr := hash.NewImageHashChallenger(network.NewMeshHandler(meshHandlerOpts), pastelClient, store,
		config.CreatorPastelID, config.CreatorPastelIDPassphrase)
	thumbnailMeshOpts := meshHandlerOpts
	thumbnailMeshOpts.NodeMaker = &network.NftSearchingNodeMaker{}
	thumbnailMeshOpts.Configs.UseMaxNodes = true

	return &Scorer{
		imgHashChlngr:   imgChlngr,
		thumbnailChlngr: thumbnail.NewThumbnailsChallenger(network.NewMeshHandler(thumbnailMeshOpts), pastelClient),
	}
}
