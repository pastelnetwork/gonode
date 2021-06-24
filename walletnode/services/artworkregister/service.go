package artworkregister

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/fs"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

const (
	logPrefix = "artwork"
)

// Service represents a service for the registration artwork.
type Service struct {
	*task.Worker
	*artwork.Storage

	config       *Config
	db           storage.KeyValue
	pastelClient pastel.Client
	nodeClient   node.Client
}

// Run starts worker.
func (service *Service) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// NOTE: Before releasing, should be reomved (for testing). Used to bypass REST API.
	if test := sys.GetStringEnv("TICKET", ""); test != "" {
		ticket := struct {
			Ticket
			ImagePath *string `json:"image_path"`
		}{}
		if err := json.Unmarshal([]byte(test), &ticket); err != nil {
			return errors.Errorf("failed to marshal ticket %q : %w", test, err)
		}

		if imagePath := ticket.ImagePath; imagePath != nil {
			baseDir, filename := filepath.Split(*imagePath)
			imageStorage := artwork.NewStorage(fs.NewFileStorage(baseDir))
			ticket.Image = artwork.NewFile(imageStorage, filename)
			if err := ticket.Image.SetFormatFromExtension(filepath.Ext(filename)); err != nil {
				return err
			}
			group.Go(func() error {
				return imageStorage.Run(ctx)
			})
		}

		group.Go(func() error {
			service.AddTask(&ticket.Ticket)
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
func (service *Service) Tasks() []*Task {
	var tasks []*Task
	for _, task := range service.Worker.Tasks() {
		tasks = append(tasks, task.(*Task))
	}
	return tasks
}

// Task returns the task of the registration artwork by the given id.
func (service *Service) Task(id string) *Task {
	if t := service.Worker.Task(id); t != nil {
		return t.(*Task)
	}
	return nil
}

// AddTask runs a new task of the registration artwork and returns its taskID.
func (service *Service) AddTask(ticket *Ticket) string {
	task := NewTask(service, ticket)
	service.Worker.AddTask(task)

	return task.ID()
}

// NewService returns a new Service instance.
func NewService(config *Config, db storage.KeyValue, fileStorage storage.FileStorage, pastelClient pastel.Client, nodeClient node.Client) *Service {
	return &Service{
		config:       config,
		db:           db,
		pastelClient: pastelClient,
		nodeClient:   nodeClient,
		Worker:       task.NewWorker(),
		Storage:      artwork.NewStorage(fileStorage),
	}
}
