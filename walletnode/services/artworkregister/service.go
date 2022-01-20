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
	"github.com/pastelnetwork/gonode/walletnode/services/mixins"
	"time"
)

const (
	logPrefix       = "nft-register"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// Service represents a service for the registration NFT.
type NftRegisterService struct {
	*task.Worker
	config *Config

	nodeClient    node.ClientInterface
	imageHandler  *mixins.ImageHandler
	pastelHandler *mixins.PastelHandler

	rqClient rqnode.Client
}

// Run starts worker. //TODO: make common with the same from SenseRegisterService
func (service *NftRegisterService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.imageHandler.FileStorage.Run(ctx)
	})
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *NftRegisterService) Tasks() []*NftRegistrationTask {
	var tasks []*NftRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the registration NFT by the given id.
func (service *NftRegisterService) GetTask(id string) *NftRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftRegistrationTask)
	}
	return nil
}

// AddTask runs a new task of the registration NFT and returns its taskID.
func (service *NftRegisterService) AddTask(p *artworks.RegisterPayload) (string, error) {
	request := FromNftRegisterPayload(p)

	// get image filename from storage based on image_id
	filename, err := service.imageHandler.FileDb.Get(p.ImageID)
	if err != nil {
		return "", errors.Errorf("get image filename from storage: %w", err)
	}

	// get image data from storage
	file, err := service.imageHandler.FileStorage.File(string(filename))
	if err != nil {
		return "", errors.Errorf("get image data: %v", err)
	}
	request.Image = file

	task := NewNFTRegistrationTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// StoreFile stores file into walletnode file storage. //TODO: make common with the same from SenseRegisterService
func (service *NftRegisterService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.imageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// NewService returns a new Service instance.
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
	raptorqClient rqnode.Client,
) *NftRegisterService {
	return &NftRegisterService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		imageHandler:  mixins.NewImageHandler(fileStorage, db, defaultImageTTL),
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		rqClient:      raptorqClient,
	}
}
