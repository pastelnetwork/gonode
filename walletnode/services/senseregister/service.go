package senseregister

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/mixins"

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
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

const (
	logPrefix       = "sense"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// SenseRegistrationService represents a service for Sense Open API
type SenseRegistrationService struct {
	*task.Worker
	config *Config

	downloadHandler *download.NftDownloadingService
	ImageHandler    *mixins.FilesHandler
	pastelHandler   *mixins.PastelHandler
	nodeClient      node.ClientInterface
	historyDB       storage.LocalStoreInterface
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

// Tasks starts worker. //TODO: make common with the same from NftRegisterService
func (service *SenseRegistrationService) Tasks() []*SenseRegistrationTask {
	var tasks []*SenseRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*SenseRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the Sense OpenAPI by the given id.
func (service *SenseRegistrationService) GetTask(id string) *SenseRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*SenseRegistrationTask)
	}
	return nil
}

// ValidateUser validates the user by the given id and pass.
func (service *SenseRegistrationService) ValidateUser(ctx context.Context, id string, pass string) bool {
	return common.ValidateUser(ctx, service.pastelHandler.PastelClient, id, pass) &&
		common.IsPastelIDTicketRegistered(ctx, service.pastelHandler.PastelClient, id)
}

// AddTask create ticket request and start a new task with the given payload
func (service *SenseRegistrationService) AddTask(p *sense.StartProcessingPayload) (string, error) {
	request := FromStartProcessingPayload(p)

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

// StoreFile stores file into walletnode file storage
func (service *SenseRegistrationService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.ImageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// CalculateFee stores file into walletnode file storage
func (service *SenseRegistrationService) CalculateFee(ctx context.Context, fileID string) (float64, error) {
	fileData, err := service.ImageHandler.GetImgData(fileID)
	if err != nil {
		return 0.0, err
	}

	return service.pastelHandler.GetEstimatedSenseFee(ctx, utils.GetFileSizeInMB(fileData))
}

// NewService returns a new Service instance
//
//	NB: Because NewNftApiHandler calls AddTask, a SenseRegisterTask will actually
//		be instantiated instead of a generic Task.
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.ClientInterface,
	fileStorage storage.FileStorageInterface,
	downloadService *download.NftDownloadingService,
	db storage.KeyValue,
	historyDB storage.LocalStoreInterface,
) *SenseRegistrationService {
	return &SenseRegistrationService{
		Worker:          task.NewWorker(),
		config:          config,
		nodeClient:      nodeClient,
		downloadHandler: downloadService,
		ImageHandler:    mixins.NewFilesHandler(fileStorage, db, defaultImageTTL),
		pastelHandler:   mixins.NewPastelHandler(pastelClient),
		historyDB:       historyDB,
	}
}
