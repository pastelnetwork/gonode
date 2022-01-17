package artworkregister

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"time"
)

const (
	logPrefix       = "nft-register"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// Service represents a service for the registration NFT.
type Service struct {
	*task.Worker
	config *Config

	nodeClient    node.ClientInterface
	ImageHandler  *common.ImageHandler
	PastelHandler *common.PastelHandler

	rqClient rqnode.Client
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.ImageHandler.FileStorage.Run(ctx)
	})
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *Service) Tasks() []*NftRegistrationTask {
	var tasks []*NftRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the registration NFT by the given id.
func (service *Service) GetTask(id string) *NftRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftRegistrationTask)
	}
	return nil
}

// AddTask runs a new task of the registration NFT and returns its taskID.
func (service *Service) AddTask(p *artworks.RegisterPayload) (string, error) {
	ticket := FromNftRegisterPayload(p)

	// get image filename from storage based on image_id
	filename, err := service.ImageHandler.FileDb.Get(p.ImageID)
	if err != nil {
		return "", errors.Errorf("get image filename from storage: %w", err)
	}

	// get image data from storage
	file, err := service.ImageHandler.FileStorage.File(string(filename))
	if err != nil {
		return "", errors.Errorf("get image data: %v", err)
	}
	ticket.Image = file

	task := NewNFTRegistrationTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// NewService returns a new Service instance.
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
	raptorqClient rqnode.Client,
) *Service {
	return &Service{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		ImageHandler:  common.NewImageHandler(fileStorage, db, defaultImageTTL),
		PastelHandler: common.NewPastelHandler(pastelClient),
		rqClient:      raptorqClient,
	}
}
