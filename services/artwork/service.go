package artwork

// const logPrefix = "[artwork]"

// Service represent artwork service.
type Service struct {
	jobs map[int]*Job
}

// RegisterJob returns the job of the registration artwork.
func (service *Service) RegisterJob(jobID int) *Job {
	if job, ok := service.jobs[jobID]; ok {
		return job
	}
	return nil
}

// Register runs a new job of the registration artwork and returns its jobID.
func (service *Service) Register(artwork *Artwork) (int, error) {
	// NOTE: for testing
	job := NewJob(5)
	job.status = JobStatusStarted
	service.jobs[job.ID] = job

	return job.ID, nil
}

// NewService returns a new Service instance.
func NewService() *Service {
	return &Service{
		jobs: make(map[int]*Job),
	}
}
