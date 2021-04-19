package artwork

import (
	"context"
	"time"

	"github.com/pastelnetwork/walletnode/dao"
	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

const ServiceName = "artwork"

// const logPrefix = "[artwork]"

// Service represent artwork service.
type Service struct {
	db     dao.KeyValue
	worker *register.Worker
	jobs   []*register.Job
}

func (service *Service) Name() string {
	return ServiceName
}

// Jobs returns all jobs.
func (service *Service) Jobs() []*register.Job {
	return service.jobs
}

func (service *Service) Run(ctx context.Context) error {
	return service.worker.Run(ctx)
}

// Job returns the job of the registration artwork.
func (service *Service) Job(jobID int) *register.Job {
	for _, job := range service.jobs {
		if job.ID() == jobID {
			return job
		}
	}
	return nil
}

// Register runs a new job of the registration artwork and returns its jobID.
func (service *Service) Register(ctx context.Context, artwork *Artwork) (int, error) {
	// NOTE: for testing
	job := register.NewJob()
	service.jobs = append(service.jobs, job)
	service.worker.AddJob(ctx, job)

	return job.ID(), nil
}

// New returns a new Service instance.
func New(db dao.KeyValue) *Service {
	ser := &Service{
		db:     db,
		worker: register.NewWorker(),
	}

	job := register.NewJob()
	ser.jobs = append(ser.jobs, job)
	go func() {
		// NOTE: for testing
		time.Sleep(time.Second)
		job.State.Update(state.NewMessage(state.StatusAccepted))
		time.Sleep(time.Second)
		job.State.Update(state.NewMessage(state.StatusActivation))
		time.Sleep(time.Second)
		msg := state.NewMessage(state.StatusActivated)
		msg.Latest = true
		job.State.Update(msg)
	}()

	return ser
}
