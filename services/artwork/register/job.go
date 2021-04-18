package register

import (
	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

type Job struct {
	State *state.State
	id    int
}

// ID returns state id.
func (job *Job) ID() int {
	return job.id
}

func NewJob(id int) *Job {
	return &Job{
		State: state.New(state.NewMessage(state.StatusStarted)),
		id:    id,
	}
}
