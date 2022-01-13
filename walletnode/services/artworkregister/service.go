package artworkregister

import (
	"context"
	"encoding/json"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/pastel"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "artwork"
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker
	*files.Storage

	config       *Config
	db           storage.KeyValue
	pastelClient pastel.Client
	nodeClient   node.Client
	rqClient     rqnode.Client
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// NOTE: Before releasing, should be reomved (for testing). Used to bypass REST API.
	if test := sys.GetStringEnv("TICKET", ""); test != "" {
		ticket := struct {
			Request
			ImagePath *string `json:"image_path"`
		}{}
		if err := json.Unmarshal([]byte(test), &ticket); err != nil {
			return errors.Errorf("marshal ticket %q : %w", test, err)
		}

		if imagePath := ticket.ImagePath; imagePath != nil {
			baseDir, filename := filepath.Split(*imagePath)
			imageStorage := files.NewStorage(fs.NewFileStorage(baseDir))
			ticket.Image = files.NewFile(imageStorage, filename)
			if err := ticket.Image.SetFormatFromExtension(filepath.Ext(filename)); err != nil {
				return err
			}
			group.Go(func() error {
				return imageStorage.Run(ctx)
			})
		}

		group.Go(func() error {
			service.AddTask(&ticket.Request)
			return nil
		})
	}

	group.Go(func() error {
		return service.Storage.Run(ctx)
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

// Task returns the task of the registration artwork by the given id.
func (service *Service) GetTask(id string) *NftRegistrationTask {
	if t := service.Worker.Task(id); t != nil {
		return t.(*NftRegistrationTask)
	}
	return nil
}

// AddTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) AddTask(ticket *Request) string {
	task := NewNFTRegistrationTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, db storage.KeyValue, fileStorage storage.FileStorageInterface, pastelClient pastel.Client, nodeClient node.Client, raptorqClient rqnode.Client) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		rqClient:     raptorqClient,
		Worker:       task.NewWorker(),
		Storage:      files.NewStorage(fileStorage),
	}
}
