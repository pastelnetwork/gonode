package senseregister

import (
	"context"
	"github.com/pastelnetwork/gonode/mixins"
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
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix       = "sense"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// Service represents a service for Sense Open API
type SenseRegistrationService struct {
	*task.Worker
	config *Config

	ImageHandler  *mixins.FilesHandler
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface
}

// Run starts worker.
func (service *SenseRegistrationService) Run(ctx context.Context) error {
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

// Run starts worker. //TODO: make common with the same from NftRegisterService
func (service *SenseRegistrationService) Tasks() []*SenseRegisterTask {
	var tasks []*SenseRegisterTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*SenseRegisterTask))
	}
	return tasks
}

// SenseRegisterTask returns the task of the Sense OpenAPI by the given id.
func (service *SenseRegistrationService) GetTask(id string) *SenseRegisterTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*SenseRegisterTask)
	}
	return nil
}

// AddTask create ticket request and start a new task with the given payload
func (service *SenseRegistrationService) AddTask(p *sense.StartProcessingPayload) (string, error) {
	request := FromSenseRegisterPayload(p)

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
	request.Image = file

	task := NewSenseRegisterTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// StoreFile stores file into walletnode file storage. //TODO: make common with the same from NftRegisterService
func (service *SenseRegistrationService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.ImageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// StoreFile stores file into walletnode file storage. //TODO: make common with the same from NftRegisterService
func (service *SenseRegistrationService) ValidateDetailsAndCalculateFee(ctx context.Context, fileID string, fileSignature string, pastelID string) (float64, error) {
	fileData, err := service.ImageHandler.GetImgData(fileID)
	if err != nil {
		return 0.0, err
	}

	fileDataInMb := int64(len(fileData)) / (1024 * 1024)

	// Validate image signature
	ok, err := service.pastelHandler.VerifySignature(ctx,
		fileData,
		fileSignature,
		pastelID,
		pastel.SignAlgorithmED448)
	if err != nil {
		return 0.0, err
	}
	if !ok {
		return 0.0, errors.Errorf("Signature doesn't match")
	}

	return service.pastelHandler.GetEstimatedActionFee(ctx, fileDataInMb)
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
) *SenseRegistrationService {
	return &SenseRegistrationService{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		ImageHandler:  mixins.NewFilesHandler(fileStorage, db, defaultImageTTL),
		pastelHandler: mixins.NewPastelHandler(pastelClient),
	}
}
