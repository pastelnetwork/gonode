package artworksearch

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch/thumbnail"
)

const (
	logPrefix = "nft-search"
)

// Service represents a service for the NFT search.
type NftSearchService struct {
	*task.Worker

	config        *Config
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
}

// Run starts worker.
func (service *NftSearchService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *NftSearchService) Tasks() []*NftSearchTask {
	var tasks []*NftSearchTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftSearchTask))
	}
	return tasks
}

// GetTask returns the task of the NFT search by the given id.
func (service *NftSearchService) GetTask(id string) *NftSearchTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftSearchTask)
	}
	return nil
}

// AddTask runs a new task of the NFT search and returns its taskID.
func (service *NftSearchService) AddTask(p *artworks.ArtSearchPayload) string {

	request := FromArtSearchRequest(p)
	task := NewNftSearchTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
) *NftSearchService {
	return &NftSearchService{
		Worker:        task.NewWorker(),
		config:        config,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		nodeClient:    nodeClient,
	}
}

// RegTicket pull NFT registration ticket from cNode & decodes base64 encoded fields
func (service *NftSearchService) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := service.pastelHandler.PastelClient.RegTicket(ctx, RegTXID)
	if err != nil {
		return nil, errors.Errorf("fetch: %w", err)
	}

	articketData, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		return nil, errors.Errorf("convert NFT ticket: %w", err)
	}

	regTicket.RegTicketData.NFTTicketData = *articketData

	return &regTicket, nil
}

// GetThumbnail gets thumbnail
func (service *NftSearchService) GetThumbnail(ctx context.Context, regTicket *pastel.RegTicket, secInfo *alts.SecInfo) (data []byte, err error) {
	thumbnailHelper := thumbnail.New(service.pastelHandler.PastelClient, service.nodeClient, service.config.ConnectToNodeTimeout)

	if err := thumbnailHelper.Connect(ctx, 1, secInfo); err != nil {
		return data, fmt.Errorf("connect Thumbnail helper: %w", err)
	}
	defer thumbnailHelper.Close()

	return thumbnailHelper.Fetch(ctx, regTicket.RegTicketData.NFTTicketData.AppTicketData.PreviewHash)
}
