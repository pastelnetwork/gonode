package artwork

import (
	"context"
	"time"

	"github.com/pastelnetwork/walletnode/services/artwork/register"
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

// const logPrefix = "[artwork]"

// Service represent artwork service.
type Service struct {
	jobs []*register.Job
}

// Jobs returns all jobs.
func (service *Service) Jobs() []*register.Job {
	return service.jobs
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
	job := register.NewJob(5)
	service.jobs = append(service.jobs, job)

	return job.ID(), nil
}

// NewService returns a new Service instance.
func NewService() *Service {
	// return &Service{
	// 	jobs: make(map[int]*register.Job),
	// }

	ser := &Service{}

	job := register.NewJob(5)
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
