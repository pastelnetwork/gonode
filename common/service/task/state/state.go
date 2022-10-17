//go:generate mockery --name=State

package state

import (
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/storage/local"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
)

// State represents a state of the task.
type State interface {
	// Status returns the current status.
	Status() *Status

	// SetStatusNotifyFunc sets a function to be called after the state is updated.
	SetStatusNotifyFunc(fn func(status *Status))

	// RequiredStatus returns an error if the current status doen't match the given one.
	RequiredStatus(subStatus SubStatus) error

	// StatusHistory returns all history from the very beginning.
	StatusHistory() []*Status

	// UpdateStatus updates the status of the state by creating a new status with the given `status`.
	UpdateStatus(subStatus SubStatus)

	// SubscribeStatus returns a new subscription of the state.
	SubscribeStatus() func() <-chan *Status

	//SetStateLog set the wallet node task status log to the state status log
	SetStateLog(statusLog types.Fields)
}

type state struct {
	status  *Status
	history []*Status

	notifyFn func(status *Status)
	sync.RWMutex
	subsCh    []chan *Status
	taskID    string
	store     storage.LocalStoreInterface
	statusLog types.Fields
}

// Status implements State.Status()
func (state *state) Status() *Status {
	return state.status
}

// SetStatusNotifyFunc implements State.SetStatusNotifyFunc()
func (state *state) SetStatusNotifyFunc(fn func(status *Status)) {
	state.notifyFn = fn
}

// RequiredStatus implements State.RequiredStatus()
func (state *state) RequiredStatus(subStatus SubStatus) error {
	if state.status.Is(subStatus) {
		return nil
	}
	return errors.Errorf("required status %q, current %q", subStatus, state.status)
}

// StatusHistory implements State.StatusHistory()
func (state *state) StatusHistory() []*Status {
	state.RLock()
	defer state.RUnlock()

	return append(state.history, state.status)
}

// UpdateStatus implements State.UpdateStatus()
func (state *state) UpdateStatus(subStatus SubStatus) {
	state.Lock()
	defer state.Unlock()

	status := NewStatus(subStatus)
	state.history = append(state.history, state.status)
	state.status = status

	history := types.TaskHistory{CreatedAt: time.Now(), TaskID: state.taskID, Status: status.String()}
	if state.statusLog.IsValid() {
		history.Details = types.NewDetails(status.String(), state.statusLog)
	}

	if state.store != nil {
		if _, err := state.store.InsertTaskHistory(history); err != nil {
			log.WithError(err).Error("unable to store task status")
		}
	}

	if state.notifyFn != nil {
		state.notifyFn(status)
	}

	for _, subCh := range state.subsCh {
		subCh := subCh
		go func() {
			subCh <- status
		}()
	}
}

// SubscribeStatus implements State.SubscribeStatus()
func (state *state) SubscribeStatus() func() <-chan *Status {
	state.RLock()
	defer state.RUnlock()

	subCh := make(chan *Status)
	state.subsCh = append(state.subsCh, subCh)

	for _, status := range append(state.history, state.status) {
		status := status
		go func() {
			subCh <- status
		}()
	}

	sub := func() <-chan *Status {
		return subCh
	}
	return sub
}

func (state *state) SetStateLog(statusLog types.Fields) {
	state.statusLog = statusLog
}

// New returns a new state instance.
func New(subStatus SubStatus, taskID string) State {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithError(err).Error("error opening history db")
	}

	if store != nil {
		if _, err := store.InsertTaskHistory(types.TaskHistory{CreatedAt: time.Now(), TaskID: taskID,
			Status: subStatus.String()}); err != nil {
			log.WithError(err).Error("unable to store task status")
		}
	}

	return &state{
		status: NewStatus(subStatus),
		taskID: taskID,
		store:  store,
	}
}
