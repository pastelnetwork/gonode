package artworkdownload

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "nft-download"
)

// Service represents a service for the registration NFT.
type Service struct {
	*task.Worker

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.ClientInterface
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
func (service *Service) Tasks() []*NftDownloadTask {
	var tasks []*NftDownloadTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftDownloadTask))
	}
	return tasks
}

// GetTask returns the task of the NFT downloading by the given id.
func (service *Service) GetTask(id string) *NftDownloadTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftDownloadTask)
	}
	return nil
}

// AddTask adds a new task of the NFT downloading and returns its taskID.
func (service *Service) AddTask(ticket *Ticket) string {
	task := NewNftDownloadTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface) *Service {
	return &Service{
		config:       config,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
	}
}
