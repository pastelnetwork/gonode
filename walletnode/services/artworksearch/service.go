package artworksearch

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/artworksearch/thumbnail"
)

const (
	logPrefix = "artwork"
)

// Service represents a service for the artwork search.
type Service struct {
	*task.Worker
	pastelClient pastel.Client
	nodeClient   node.Client
	config       *Config
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})

	return group.Wait()
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*NftSearchTask {
	var tasks []*NftSearchTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftSearchTask))
	}
	return tasks
}

// NftSearchTask returns the task of the artwork search by the given id.
func (service *Service) GetTask(id string) *NftSearchTask {
	return service.Worker.Task(id).(*NftSearchTask)
}

// AddTask runs a new task of the artwork search and returns its taskID.
func (service *Service) AddTask(request *ArtSearchRequest) string {
	task := NewNftSearchTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.Client) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		Worker:       task.NewWorker(),
		nodeClient:   nodeClient,
	}
}

// RegTicket pull NFT registration ticket from cNode & decodes base64 encoded fields
func (service *Service) RegTicket(ctx context.Context, RegTXID string) (*pastel.RegTicket, error) {
	regTicket, err := service.pastelClient.RegTicket(ctx, RegTXID)
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
func (service *Service) GetThumbnail(ctx context.Context, regTicket *pastel.RegTicket, secInfo *alts.SecInfo) (data []byte, err error) {
	thumbnailHelper := thumbnail.New(service.pastelClient, service.nodeClient, service.config.ConnectToNodeTimeout)

	if err := thumbnailHelper.Connect(ctx, 1, secInfo); err != nil {
		return data, fmt.Errorf("connect Thumbnail helper: %w", err)
	}
	defer thumbnailHelper.Close()

	return thumbnailHelper.Fetch(ctx, regTicket.RegTicketData.NFTTicketData.AppTicketData.PreviewHash)
}
