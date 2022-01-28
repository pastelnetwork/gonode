package nftdownload

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

// Service represents a service for the registration NFT.
type NftDownloadService struct {
	*task.Worker

	config        *Config
	nodeClient    node.ClientInterface
	pastelHandler *mixins.PastelHandler
}

// Run starts worker.
func (service *NftDownloadService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *NftDownloadService) Tasks() []*NftDownloadTask {
	var tasks []*NftDownloadTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftDownloadTask))
	}
	return tasks
}

// GetTask returns the task of the NFT downloading by the given id.
func (service *NftDownloadService) GetTask(id string) *NftDownloadTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftDownloadTask)
	}
	return nil
}

// AddTask adds a new task of the NFT downloading and returns its taskID.
func (service *NftDownloadService) AddTask(p *nft.NftDownloadPayload) string {
	request := FromDownloadPayload(p)

	task := NewNftDownloadTask(service, request)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewNftDownloadService returns a new Service instance.
func NewNftDownloadService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface) *NftDownloadService {
	return &NftDownloadService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		pastelHandler: mixins.NewPastelHandler(pastelClient),
	}
}
