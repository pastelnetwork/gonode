package test

import (
	"context"
	"testing"

	taskCommon "github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/mocks"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/stretchr/testify/mock"
)

const (
	// IDMethod represent ID name method
	IDMethod = "ID"

	// RunMethod represent Run name method
	RunMethod = "Run"

	// CancelMethod represent Cancel name method
	CancelMethod = "Cancel"

	// DoneMethod represent Donce name method
	DoneMethod = "Done"

	// RunActionMethod represent RunAction name method
	RunActionMethod = "RunAction"

	// NewActionMethod represent NewAction name method
	NewActionMethod = "NewAction"

	// StatusMethod represent Status state name method
	StatusMethod = "Status"

	// SetStatusNotifyFuncMethod represent SetStatusNotifyFunc state name method
	SetStatusNotifyFuncMethod = "SetStatusNotifyFunc"

	// RequiredStatusMethod represent RequiredStatus state name method
	RequiredStatusMethod = "RequiredStatus"

	// StatusHistoryMethod represent StatusHistory state name method
	StatusHistoryMethod = "StatusHistory"

	// UpdateStatusMethod represent UpdateStatus state name method
	UpdateStatusMethod = "UpdateStatus"

	// SubscribeStatusMethod represent SubscribeStatus state name method
	SubscribeStatusMethod = "SubscribeStatus"

	// ActionFn represent task.ActionFn function type
	ActionFn = "task.ActionFn"
)

//Task implementing task.Task mock for testing purpose
type Task struct {
	t *testing.T
	*mocks.Task
}

// NewMockTask new Task instance
func NewMockTask(t *testing.T) *Task {
	return &Task{
		t:    t,
		Task: &mocks.Task{},
	}
}

// ListenOnID listening ID call and returns task id from args
func (task *Task) ListenOnID(id string) *Task {
	task.On(IDMethod).Return(id)
	return task
}

// AssertIDCall ID call method assertion
func (task *Task) AssertIDCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, IDMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, IDMethod, expectedCalls)
	return task
}

// ListenOnRun listening Run call and returns error from args
func (task *Task) ListenOnRun(returnErr error) *Task {
	task.On(RunMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRunCall Run call assertion
func (task *Task) AssertRunCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, RunMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, RunMethod, expectedCalls)
	return task
}

// ListenOnCancel listening Cancel call
func (task *Task) ListenOnCancel() *Task {
	task.On(CancelMethod).Return(nil)
	return task
}

// AssertCancelCall Cancel call assertion
func (task *Task) AssertCancelCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, CancelMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, CancelMethod, expectedCalls)
	return task
}

// ListenOnDone listening Done call and returns chan from args
func (task *Task) ListenOnDone(c <-chan struct{}) *Task {
	task.On(DoneMethod).Return(c)
	return task
}

// AssertDoneCall assertion Done call
func (task *Task) AssertDoneCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, DoneMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, DoneMethod, expectedCalls)
	return task
}

// ListenOnNewAction listening NewAction call and return chan from args
func (task *Task) ListenOnNewAction(ctx context.Context, c chan struct{}) *Task {
	handleFunc := func(fn taskCommon.ActionFn) <-chan struct{} {
		defer close(c)
		fn(ctx)
		return c
	}
	task.On(NewActionMethod, mock.Anything).Return(handleFunc)
	return task
}

// AssertNewActionCall assertion NewAction call
func (task *Task) AssertNewActionCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, NewActionMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, NewActionMethod, expectedCalls)
	return task
}

// ListenOnRunAction listening RunAction call and returns error from args
func (task *Task) ListenOnRunAction(returnErr error) *Task {
	task.On(RunActionMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRunActionCall RunAction call assertion
func (task *Task) AssertRunActionCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, RunActionMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, RunActionMethod, expectedCalls)
	return task
}

// ListenOnSetStatusNotifyFunc listening SetStatusNotifyFunc call and return values from args
func (task *Task) ListenOnSetStatusNotifyFunc(arguments ...interface{}) *Task {
	// handleFunc := func(fn func(status *state.Status)) {
	// 	fn(status)
	// }
	task.On(SetStatusNotifyFuncMethod, mock.Anything).Return(arguments...)
	return task
}

// AssertSetStatusNotifyFuncCall SetSTatusNotifyFunc call assertion
func (task *Task) AssertSetStatusNotifyFuncCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, SetStatusNotifyFuncMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, SetStatusNotifyFuncMethod, expectedCalls)
	return task
}

// ListenOnRequiredStatus listening RequiredStatus call and returns error from args
func (task *Task) ListenOnRequiredStatus(returnErr error) *Task {
	task.On(RequiredStatusMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRequiredStatusCall assesrtion RequiredStatus call
func (task *Task) AssertRequiredStatusCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, RequiredStatusMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, RequiredStatusMethod, expectedCalls)
	return task
}

// ListenOnStatusHistory listening on StatusHistory and returns statuses from args
func (task *Task) ListenOnStatusHistory(statuses []*state.Status) *Task {
	task.On(StatusHistoryMethod).Return(statuses)
	return task
}

// AssertStatusHistory assertion status history method
func (task *Task) AssertStatusHistory(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, StatusHistoryMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, StatusHistoryMethod, expectedCalls)
	return task
}

// ListenOnUpdateStatus listening UpdateStatus call
func (task *Task) ListenOnUpdateStatus() *Task {
	task.On(UpdateStatusMethod, mock.Anything).Return(nil)
	return task
}

// AssertUpdateStatusCall UpdateStatus call assertion
func (task *Task) AssertUpdateStatusCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, UpdateStatusMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, UpdateStatusMethod, expectedCalls)
	return task
}

// ListenOnSubscribeStatus listening SubscribeStatus and returns status from args
func (task *Task) ListenOnSubscribeStatus(status *state.Status) *Task {
	handleFunc := func() func() <-chan *state.Status {
		ch := make(chan *state.Status)
		go func() {
			ch <- status
		}()

		sub := func() <-chan *state.Status {
			return ch
		}
		return sub
	}
	task.On(SubscribeStatusMethod).Return(handleFunc)
	return task
}

// AssertSubscribeStatusCall assertion SubscribeStatus call
func (task *Task) AssertSubscribeStatusCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, SubscribeStatusMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, SubscribeStatusMethod, expectedCalls)
	return task
}
