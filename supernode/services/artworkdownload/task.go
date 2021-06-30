package artworkdownload

import (
	"context"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
)

// Task is the task of registering new artwork.
type Task struct {
	task.Task
	*Service

	ResampledArtwork *artwork.File
	Artwork          *artwork.File

	acceptedMu sync.Mutex
	accpeted   Nodes

	connectedTo *Node
}

// Run starts the task
func (task *Task) Run(ctx context.Context) error {
	ctx = task.context(ctx)
	defer log.WithContext(ctx).Debug("Task canceled")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debugf("States updated")
	})

	return task.RunAction(ctx)
}

// Download downloads image and return the image.
func (task *Task) Download(_ context.Context, txid, timestamp, signature, ttxid string) ([]byte, error) {
	if err := task.RequiredStatus(StatusConnected); err != nil {
		return nil, err
	}

	var file []byte

	<-task.NewAction(func(ctx context.Context) error {
		// Validate timestamp is not older then 10 minutes
		// Get Art Registration ticket by txid

		if len(ttxid) == 0 {
			// validate timestamp signature with PastelID from Trade ticket
			// by calling command `pastelid verify timestamp-string sig PastelID passphrase`
		} else {
			// Get list of non sold Trade ticket owened by the owner of the PastelID from request
			// by calling command `tickets list trade available`

			// Validate that Trade ticket with ttxid is in the list

			// Validate timestamp signature with PastelID from Trade ticket
			// by calling command `pastelid verify timestamp-string sig PastelID passphrase`
		}

		// Get the list of “symbols/chunks” - rq_ids - from Art Registration ticket and request them from Kademlia

		// Validate that the hash of each “symbol/chunk” matches its id

		// When all “symbols/chunks” are received, pass them to the Kademlia server to decode (also passing encoder parameters: rq_coti" and “rq_ssoti”)

		// Validate hash of the restored image matches the image hash in the Art Reistration ticket (data_hash)
		return nil
	})

	return file, nil
}

func (task *Task) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, task.ID()))
}

// NewTask returns a new Task instance.
func NewTask(service *Service) *Task {
	return &Task{
		Task:    task.New(StatusTaskStarted),
		Service: service,
	}
}
