package download

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "walletnode-nft-download"
)

// NftDownloadingService represents a service for the registration NFT.
type NftDownloadingService struct {
	*task.Worker

	config        *Config
	nodeClient    node.ClientInterface
	pastelHandler *mixins.PastelHandler
}

// Run starts worker.
func (service *NftDownloadingService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
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
func NewNftDownloadService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface) *NftDownloadingService {
	return &NftDownloadingService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
	}
}
