package register

import (
	"sync/atomic"

	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

var taskID uint32

// Task is the task of registering new artwork.
type Task struct {
	ticket *Ticket
	State  *state.State
	id     int
}

// ID returns state id.
func (task *Task) ID() int {
	return task.id
}

// NewTask returns a new Task instance.
func NewTask(ticket *Ticket) *Task {
	return &Task{
		id:     int(atomic.AddUint32(&taskID, 1)),
		ticket: ticket,
		State:  state.New(state.NewMessage(state.StatusStarted)),
	}
}
