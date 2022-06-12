package download

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/log"

	"github.com/pastelnetwork/gonode/bridge/node"
	"github.com/pastelnetwork/gonode/bridge/services/common"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
)

type Service struct {
	connRefreshInterval int
	connections         int
	th                  *thumbnailHandler
	dd                  *ddFpHandler
	blockTracker        *blocktracker.BlockCntTracker
}

func NewService(config *Config, pastelClient pastel.Client, nodeClient node.SNClientInterface) Service {
	meshHandlerOpts := common.MeshHandlerOpts{
		NodeMaker:     &DownloadingNodeMaker{},
		PastelHandler: mixins.NewPastelHandler(pastelClient),
		NodeClient:    nodeClient,
		Configs: &common.MeshHandlerConfig{
			ConnectToNextNodeDelay: config.ConnectToNextNodeDelay,
			ConnectToNodeTimeout:   config.ConnectToNodeTimeout,
			AcceptNodesTimeout:     config.AcceptNodesTimeout,
			PastelID:               config.PastelID,
			Passphrase:             config.Passphrase,
		},
	}

	return Service{
		connRefreshInterval: config.ConnectionRefreshInterval,
		connections:         config.Connections,
		th:                  newThumbnailHandler(common.NewMeshHandler(meshHandlerOpts)),
		dd:                  newDDFPHandler(common.NewMeshHandler(meshHandlerOpts)),
		blockTracker:        blocktracker.New(pastelClient),
	}
}

func (b Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return b.runThumbnailHelper(ctx)
	})

	group.Go(func() error {
		return b.runDDFpHandler(ctx)
	})

	return group.Wait()
}

func (s *Service) runThumbnailHelper(ctx context.Context) error {
	baseBlkCnt, err := s.blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	newCtx, cancel := context.WithCancel(ctx)
	if err := s.th.Connect(newCtx, s.connections, cancel); err != nil {
		log.WithContext(ctx).WithError(err).Error("connect and setup thumbnails fetchers failed")
	} else {
		log.WithContext(ctx).Info("thumbnail handler up")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			currBlkCnt, err := s.blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				return err
			}

			if currBlkCnt-baseBlkCnt < int32(s.connRefreshInterval) {
				break
			}

			if err := s.th.CloseAll(newCtx); err != nil {
				log.WithContext(ctx).WithError(err).Error("close thumbnails fetchers failed")
				break
			}
			log.WithContext(ctx).Info("closed current thumbnail handler connections.")

			newCtx, cancel := context.WithCancel(ctx)
			if err := s.th.Connect(newCtx, s.connections, cancel); err != nil {
				log.WithContext(ctx).WithError(err).Error("connect and setup thumbnails fetchers failed")
				break
			}
			baseBlkCnt = currBlkCnt
			log.WithContext(ctx).Info("thumbnail handler up")
		}
	}
}

func (s *Service) runDDFpHandler(ctx context.Context) error {
	baseBlkCnt, err := s.blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return err
	}

	newCtx, cancel := context.WithCancel(ctx)
	if err := s.dd.Connect(newCtx, s.connections, cancel); err != nil {
		log.WithContext(ctx).WithError(err).Error("connect and setup fp_dd fetchers failed")
	} else {
		log.WithContext(ctx).Info("dd & fp handler up")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(10 * time.Second):
			currBlkCnt, err := s.blockTracker.GetBlockCount()
			if err != nil {
				log.WithContext(ctx).WithError(err).Warn("failed to get block count")
				return err
			}

			if currBlkCnt-baseBlkCnt < int32(s.connRefreshInterval) {
				break
			}

			log.WithContext(ctx).Info("refresh dd & fp handler connections")

			if err := s.dd.CloseAll(newCtx); err != nil {
				log.WithContext(ctx).WithError(err).Error("close fp_dd fetchers failed")
				break
			}
			log.WithContext(ctx).Info("closed current dd & fp handler connections")

			newCtx, cancel := context.WithCancel(ctx)
			if err := s.dd.Connect(newCtx, s.connections, cancel); err != nil {
				log.WithContext(ctx).WithError(err).Error("connect and setup fp_dd fetchers failed")
				break
			}
			baseBlkCnt = currBlkCnt
			log.WithContext(ctx).Info("dd & fp handler up")
		}
	}
}

// FetchThumbnail gets the actual thumbnail data from the network as bytes to be wrapped by the calling function
func (s *Service) FetchThumbnail(ctx context.Context, txid string, numnails int) (data map[int][]byte, err error) {
	return s.th.Fetch(ctx, txid, numnails)
}

// FetchDDAndFpData gets the actual dd data from the network as bytes to be wrapped by the calling function
func (s *Service) FetchDupeDetectionData(ctx context.Context, txid string) (data []byte, err error) {
	return s.dd.Fetch(ctx, txid)
}
