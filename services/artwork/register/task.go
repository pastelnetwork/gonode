package register

import (
	"sync/atomic"

	"github.com/pastelnetwork/walletnode/services/artwork/register/state"
)

var taskID uint32

// Task is the task of registering new artwork.
type Task struct {
	*Service

	ID     int
	State  *state.State
	Ticket *Ticket
}

// NewTask returns a new Task instance.
func NewTask(service *Service, Ticket *Ticket) *Task {
	return &Task{
		ID:     int(atomic.AddUint32(&taskID, 1)),
		Ticket: Ticket,
		State:  state.New(state.NewMessage(state.StatusTaskStarted)),
	}
}
