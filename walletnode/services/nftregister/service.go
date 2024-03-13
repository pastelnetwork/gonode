package nftregister

import (
	"context"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/mixins"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	rqgrpc "github.com/pastelnetwork/gonode/raptorq/node/grpc"
	"github.com/pastelnetwork/gonode/walletnode/api/gen/nft"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
	"github.com/pastelnetwork/gonode/walletnode/services/download"
)

const (
	logPrefix       = "nft-register"
	defaultImageTTL = time.Second * 3600 // 1 hour
)

// NftRegistrationService represents a service for the registration NFT.
type NftRegistrationService struct {
	*task.Worker
	config *Config

	ImageHandler  *mixins.FilesHandler
	pastelHandler *mixins.PastelHandler
	nodeClient    node.ClientInterface

	rqClient        rqnode.ClientInterface
	downloadHandler *download.NftDownloadingService
	historyDB       storage.LocalStoreInterface
}

// Run starts worker. //TODO: make common with the same from SenseRegisterService
func (service *NftRegistrationService) Run(ctx context.Context) error {
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
func (service *NftRegistrationService) Tasks() []*NftRegistrationTask {
	var tasks []*NftRegistrationTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*NftRegistrationTask))
	}
	return tasks
}

// GetTask returns the task of the registration NFT by the given id.
func (service *NftRegistrationService) GetTask(id string) *NftRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftRegistrationTask)
	}
	return nil
}

// AddTask runs a new task of the registration NFT and returns its taskID.
func (service *NftRegistrationService) AddTask(p *nft.RegisterPayload) (string, error) {
	request := FromNftRegisterPayload(p)

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
	request.FileName = string(filename)

	task := NewNFTRegistrationTask(service, request)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// StoreFile stores file into walletnode file storage. //TODO: make common with the same from SenseRegisterService
func (service *NftRegistrationService) StoreFile(ctx context.Context, fileName *string) (string, string, error) {
	return service.ImageHandler.StoreFileNameIntoStorage(ctx, fileName)
}

// CalculateFee stores file into walletnode file storage
func (service *NftRegistrationService) CalculateFee(ctx context.Context, fileID string) (float64, float64, error) {
	fileData, err := service.ImageHandler.GetImgData(fileID)
	if err != nil {
		return 0.0, 0.0, err
	}

	fileDataInMb := utils.GetFileSizeInMB(fileData)
	fee, err := service.pastelHandler.PastelClient.NFTStorageFee(ctx, int(fileDataInMb))
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("get network fee failure: %w", err)
	}

	return fileDataInMb, fee.EstimatedNftStorageFeeMax * 0.1, nil
}

// ValidateUser validates user credentials
func (service *NftRegistrationService) ValidateUser(ctx context.Context, id string, pass string) bool {
	return common.ValidateUser(ctx, service.pastelHandler.PastelClient, id, pass) &&
		common.IsPastelIDTicketRegistered(ctx, service.pastelHandler.PastelClient, id)
}

// NewService returns a new Service instance.
// NB: although it might appear that a generic task is instantiated here with NewWorker, because of the way the API calls
//
//	are  handled in NftApiHandler, an NftRegistrationTask will actually be created via AddTask.
func NewService(config *Config, pastelClient pastel.Client, nodeClient node.ClientInterface, fileStorage storage.FileStorageInterface,
	db storage.KeyValue, downloadService *download.NftDownloadingService, historyDB storage.LocalStoreInterface) *NftRegistrationService {
	return &NftRegistrationService{
		Worker:          task.NewWorker(),
		config:          config,
		nodeClient:      nodeClient,
		ImageHandler:    mixins.NewFilesHandler(fileStorage, db, defaultImageTTL),
		pastelHandler:   mixins.NewPastelHandler(pastelClient),
		downloadHandler: downloadService,
		rqClient:        rqgrpc.NewClient(),
		historyDB:       historyDB,
	}
}
