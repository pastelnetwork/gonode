package senseregister

import (
	"context"
	"github.com/pastelnetwork/gonode/walletnode/services/common"
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
type Service struct {
	*task.Worker
	config *Config

	nodeClient    node.Client
	ImageHandler  *common.ImageHandler
	PastelHandler *common.PastelHandler
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
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
func (service *Service) Tasks() []*SenseRegisterTask {
	var tasks []*SenseRegisterTask
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*SenseRegisterTask))
	}
	return tasks
}

// SenseRegisterTask returns the task of the Sense OpenAPI by the given id.
func (service *Service) GetTask(id string) *SenseRegisterTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*SenseRegisterTask)
	}
	return nil
}

// AddTask create ticket request and start a new task with the given payload
func (service *Service) AddTask(p *sense.StartProcessingPayload) (string, error) {
	ticket := &SenseRegisterRequest{
		BurnTxID:              p.BurnTxid,
		AppPastelID:           p.AppPastelID,
		AppPastelIDPassphrase: p.AppPastelidPassphrase,
	}

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

	task := NewSenseRegisterTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID(), nil
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	pastelClient pastel.Client,
	nodeClient node.Client,
	fileStorage storage.FileStorageInterface,
	db storage.KeyValue,
) *Service {
	return &Service{
		Worker:        task.NewWorker(),
		config:        config,
		nodeClient:    nodeClient,
		ImageHandler:  common.NewImageHandler(fileStorage, db, defaultImageTTL),
		PastelHandler: common.NewPastelHandler(pastelClient),
	}
}
