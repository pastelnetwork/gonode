package nftregister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/dupedetection/ddclient"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/supernode/node"
	"github.com/pastelnetwork/gonode/supernode/services/common"
)

const (
	logPrefix = "nft"
)

// NftRegistrationService represent nft service.
type NftRegistrationService struct {
	*common.SuperNodeService
	config *Config

	nodeClient node.ClientInterface
	ddClient   ddclient.DDServerClient
}

// Run starts task
func (service *NftRegistrationService) Run(ctx context.Context) error {
	return service.RunHelper(ctx, service.config.PastelID, logPrefix)
}

// NewNftRegistrationTask runs a new task of the registration Nft and returns its taskID.
func (service *NftRegistrationService) NewTask() *NftRegistrationTask {
	task := NewNftRegistrationTask(service)
	service.Worker.AddTask(task)
	return task
}

// Task returns the task of the registration Nft by the given id.
func (service *NftRegistrationService) Task(id string) *NftRegistrationTask {
	return service.Worker.Task(id).(*NftRegistrationTask)
}

// NewService returns a new Service instance.
func NewService(config *Config, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.ClientInterface, p2pClient p2p.Client, rqClient rqnode.ClientInterface, ddClient ddclient.DDServerClient) *NftRegistrationService {
	return &NftRegistrationService{
		SuperNodeService: common.NewSuperNodeService(fileStorage, pastelClient, p2pClient, rqClient),
		config:           config,
		nodeClient:       nodeClient,
		ddClient:         ddClient,
	}
}
