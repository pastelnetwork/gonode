package artwork

import "time"

// JobSubscription represents status of the process that is consumed by subscribers.
type JobSubscription struct {
	Status JobStatus
	TXID   string
}

// Job represents the process of registering a work of art.
type Job struct {
	ID          int
	status      JobStatus
	TXID        string
	subscribers []chan *JobSubscription
}

func (job *Job) newSubscriber() chan *JobSubscription {
	sub := make(chan *JobSubscription, 1)
	sub <- &JobSubscription{Status: job.status}
	job.subscribers = append(job.subscribers, sub)
	return sub
}

// Subscribe returns a new subscription of the job status.
func (job *Job) Subscribe() chan *JobSubscription {
	sub := job.newSubscriber()

	// NOTE: for testing
	go func() {
		time.Sleep(time.Second)
		sub <- &JobSubscription{Status: JobStatusAccepted}
		time.Sleep(time.Second)
		sub <- &JobSubscription{Status: JobStatusActivation}
		time.Sleep(time.Second)
		sub <- &JobSubscription{Status: JobStatusActivated}
		close(sub)
	}()

	return sub
}

// NewJob returns a new Job instance.
func NewJob(id int) *Job {
	return &Job{
		ID: id,
	}
}
