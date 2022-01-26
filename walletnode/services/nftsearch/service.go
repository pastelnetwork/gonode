package nftsearch

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "nft-search"
)

// Service represents a service for the NFT search.
type NftSearchService struct {
	*task.Worker

	config        *Config
	nodeClient    node.ClientInterface
	pastelHandler *mixins.PastelHandler
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
func (service *NftSearchService) AddTask(p *nft.NftSearchPayload) string {

	request := FromNftSearchRequest(p)
	task := NewNftSearchTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// GetThumbnail gets thumbnail
func (service *NftSearchService) GetThumbnail(ctx context.Context, regTicket *pastel.RegTicket, pastelID string, passphrase string) (data []byte, err error) {
	nftGetSearchTask := NewNftGetSearchTask(service, pastelID, passphrase)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := nftGetSearchTask.thumbnail.Connect(ctx, 1, cancel); err != nil {
		return nil, errors.Errorf("connect and setup fetchers: %w", err)
	}
	data, err = nftGetSearchTask.thumbnail.FetchOne(ctx, regTicket.RegTicketData.NFTTicketData.AppTicketData.PreviewHash)
	if err != nil {
		return nil, errors.Errorf("fetch multiple thumbnails: %w", err)
	}

	return data, nftGetSearchTask.thumbnail.CloseAll(ctx)
}

// RegTicket pull NFT registration ticket from cNode & decodes base64 encoded fields
func (service *NftSearchService) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := service.pastelHandler.PastelClient.RegTicket(ctx, RegTXID)
	if err != nil {
		return nil, errors.Errorf("fetch: %w", err)
	}

	nftTicketData, err := pastel.DecodeNFTTicket(regTicket.RegTicketData.NFTTicket)
	if err != nil {
		return nil, errors.Errorf("convert NFT ticket: %w", err)
	}

	regTicket.RegTicketData.NFTTicketData = *nftTicketData

	return &regTicket, nil
}

// NewNftSearchService returns a new Service instance.
func NewNftSearchService(config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
) *NftSearchService {
	return &NftSearchService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
	}
}
