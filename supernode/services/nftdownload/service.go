package nftdownload

import (
	"context"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "Nft-download"
)

// Service represent Nft service.
type NftDownloadService struct {
	*common.SuperNodeService
	config *Config
}

// Run starts task
func (service *NftDownloadService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewNftDownloadingTask runs a new task of the downloading Nft and returns its taskID.
func (service *NftDownloadService) NewTask() *NftDownloadingTask {
	task := NewNftDownloadingTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the registration Nft by the given id.
func (service *NftDownloadService) Task(id string) *NftDownloadingTask {
	return service.Worker.Task(id).(*NftDownloadingTask)
}

// NewService returns a new Service instance.
func NewService(config *Config, pastelClient pastel.Client, p2pClient p2p.Client, rqClient rqnode.ClientInterface) *NftDownloadService {
	return &NftDownloadService{
		SuperNodeService: common.NewSuperNodeService(nil, pastelClient, p2pClient, rqClient),
		config:           config,
	}
}
