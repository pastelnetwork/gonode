package senseregister

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "sense"
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	pastelClient pastel.Client
	nodeClient   node.Client
	db           storage.KeyValue
}

// Run starts worker.
func (service Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// Run storage service
	group.Go(func() error {
		return service.Storage.Run(ctx)
	})

	// Run worker service
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})

	return group.Wait()
}

// AddTask adds a task to the worker.
func (service *Service) AddTask(ticket *Request) string {
	task := NewTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// VerifyImageSignature verifies the signature of the image
func (service *Service) VerifyImageSignature(ctx context.Context, imgData []byte, signature string, pastelID string) error {
	ok, err := service.pastelClient.Verify(ctx, imgData, signature, pastelID, pastel.SignAlgorithmED448)
	if err != nil {
		return err
	}

	if !ok {
		return errors.Errorf("signature verification failed")
	}

	return nil
}

// GetEstimatedFee returns the estimated sense fee for the given image
func (service *Service) GetEstimatedFee(ctx context.Context, ImgSizeInMb int64) (float64, error) {
	actionFees, err := service.pastelClient.GetActionFee(ctx, ImgSizeInMb)
	if err != nil {
		return 0, err
	}

	return actionFees.SenseFee, nil
}

// NewService returns a new Service instance
func NewService(
	config *Config,
	fileStorage storage.FileStorage,
	pastelClient pastel.Client,
	nodeClient node.Client,
	db storage.KeyValue,
) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
	}
}
