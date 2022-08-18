package cascaderegister

import (
	"context"
	"time"

	// Package image/jpeg is not used explicitly in the code below,
	// but is imported for its initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images. Same with png.
	_ "image/jpeg"
	_ "image/png"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix       = "cascade"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// CascadeRegistrationService represents a service for Cascade Open API
type CascadeRegistrationService struct {
	*task.Worker
	config *Config

	ImageHandler  *mixins.FilesHandler
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface

	rqClient rqnode.ClientInterface
}

// Run starts worker.
func (service *CascadeRegistrationService) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return service.ImageHandler.FileStorage.Run(ctx)
	})

	// Run worker service
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	return group.Wait()
}

// Tasks returns all tasks.
func (service *CascadeRegistrationService) Tasks() []*CascadeRegistrationTask {
	var tasks []*CascadeRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*CascadeRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the Cascade OpenAPI by the given id.
func (service *CascadeRegistrationService) GetTask(id string) *CascadeRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*CascadeRegistrationTask)
	}
	return nil
}

// AddTask create ticket request and start a new task with the given payload
func (service *CascadeRegistrationService) AddTask(p *cascade.StartProcessingPayload) (string, error) {
	request := FromStartProcessingPayload(p)

	// get image filename from storage based on image_id
	filename, err := service.ImageHandler.FileDb.Get(p.FileID)
	if err != nil {
		return "", errors.Errorf("get image filename from storage: %w", err)
	}

	// get image data from storage
	file, err := service.ImageHandler.FileStorage.File(string(filename))
	if err != nil {
		return "", errors.Errorf("get image data: %v", err)
	}
	request.Image = file
	request.FileName = string(filename)

	task := NewCascadeRegisterTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// StoreFile stores file into walletnode file storage
func (service *CascadeRegistrationService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.ImageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// CalculateFee stores file into walletnode file storage
func (service *CascadeRegistrationService) CalculateFee(ctx context.Context, fileID string) (float64, error) {
	fileData, err := service.ImageHandler.GetImgData(fileID)
	if err != nil {
		return 0.0, err
	}

	return service.pastelHandler.GetEstimatedCascadeFee(ctx, utils.GetFileSizeInMB(fileData))
}

// NewService returns a new Service instance
// 	NB: Because NewNftApiHandler calls AddTask, a CascadeRegisterTask will actually
//		be instantiated instead of a generic Task.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface, db storage.KeyValue, raptorqClient rqnode.ClientInterface,
) *CascadeRegistrationService {
	return &CascadeRegistrationService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		ImageHandler:  mixins.NewFilesHandler(fileStorage, db, defaultImageTTL),
		pastelHandler: mixins.NewPastelHandler(pastelClient),
		rqClient:      raptorqClient,
	}
}
