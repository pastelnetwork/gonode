package common

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/p2p"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
)

// SuperNodeServiceInterface common interface for Services
type SuperNodeServiceInterface interface {
	RunHelper(ctx context.Context) error
	NewTask() task.Task
	Task(id string) task.Task
}

// SuperNodeService common "class" for Services
type SuperNodeService struct {
	*task.Worker
	*files.Storage

	PastelClient pastel.Client
	P2PClient    p2p.Client
	RQClient     rqnode.ClientInterface
}

// run starts task
func (service *SuperNodeService) run(ctx context.Context, pastelID string, prefix string) error {
	ctx = log.ContextWithPrefix(ctx, prefix)

	if pastelID == "" {
		return errors.New("PastelID is not specified in the config file")
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return service.Worker.Run(ctx)
	})
	if service.Storage != nil {
		group.Go(func() error {
			return service.Storage.Run(ctx)
		})
	}
	return group.Wait()
}

// RunHelper common code for Service runner
func (service *SuperNodeService) RunHelper(ctx context.Context, pastelID string, prefix string) error {
	for {
		select {
		case <-ctx.Done():
			log.WithContext(ctx).Error("context done - closing sn service")
			return nil
		case <-time.After(5 * time.Second):
			if err := service.run(ctx, pastelID, prefix); err != nil {
				service.Worker = task.NewWorker()
				log.WithContext(ctx).WithError(err).Error("Service run failed, retrying")
			} else {
				log.WithContext(ctx).Error("run service failed - closing sn service")
				return nil
			}
		}
	}
}

// NewSuperNodeService creates SuperNodeService
func NewSuperNodeService(
	fileStorage storage.FileStorageInterface,
	pastelClient pastel.Client,
	p2pClient p2p.Client,
	rqClient rqnode.ClientInterface,
) *SuperNodeService {
	return &SuperNodeService{
		Worker:       task.NewWorker(),
		Storage:      files.NewStorage(fileStorage),
		PastelClient: pastelClient,
		P2PClient:    p2pClient,
		RQClient:     rqClient,
	}
}
