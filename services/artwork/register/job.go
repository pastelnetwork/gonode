package register

import (
	"sync/atomic"

	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

var (
	jobID uint32
)

type Job struct {
	State *state.State
	id    int
}

// ID returns state id.
func (job *Job) ID() int {
	return job.id
}

func NewJob() *Job {
	return &Job{
		id:    int(atomic.AddUint32(&jobID, 1)),
		State: state.New(state.NewMessage(state.StatusStarted)),
	}
}
