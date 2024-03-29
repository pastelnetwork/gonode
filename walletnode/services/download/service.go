package download

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage/queries"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
)

const (
	logPrefix = "wN-download"
)

// NftDownloadingService represents a service for the registration NFT.
type NftDownloadingService struct {
	*task.Worker

	cleanup       *CleanupService
	config        *Config
	nodeClient    node.ClientInterface
	pastelHandler *mixins.PastelHandler
	historyDB     queries.LocalStoreInterface
}

// Run starts worker.
func (service *NftDownloadingService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})

	group.Go(func() error {
		return service.cleanup.Run(ctx)
	})

	return group.Wait()
}

// Tasks returns all tasks.
func (service *NftDownloadingService) Tasks() []*NftDownloadingTask {
	var tasks []*NftDownloadingTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftDownloadingTask))
	}
	return tasks
}

// ValidateUser validates the user by the given id and pass.
func (service *NftDownloadingService) ValidateUser(ctx context.Context, id string, pass string) bool {
	return common.ValidateUser(ctx, service.pastelHandler.PastelClient, id, pass) &&
		common.IsPastelIDTicketRegistered(ctx, service.pastelHandler.PastelClient, id)
}

// GetTask returns the task of the NFT downloading by the given id.
func (service *NftDownloadingService) GetTask(id string) *NftDownloadingTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftDownloadingTask)
	}
	return nil
}

// AddTask adds a new task of the NFT downloading and returns its taskID.
func (service *NftDownloadingService) AddTask(p *nft.DownloadPayload, ticketType string) string {
	request := FromDownloadPayload(p, ticketType)

	task := NewNftDownloadTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewNftDownloadService returns a new Service instance.
func NewNftDownloadService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface, historyDB queries.LocalStoreInterface) *NftDownloadingService {
	return &NftDownloadingService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		cleanup:       NewCleanupService(config.StaticDir),
		historyDB:     historyDB,
	}
}
