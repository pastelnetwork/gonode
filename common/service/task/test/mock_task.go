package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task/mocks"
	"github.com/stretchr/testify/mock"
)

const (
	// IDMethod represent ID name method
	IDMethod = "ID"

	// RunMethod represent Run name method
	RunMethod = "Run"

	// CancelMethod represent Cancel name method
	CancelMethod = "Cancel"

	// RunActionMethod represent RunAction name method
	RunActionMethod = "RunAction"

	// NewActionMethod represents NewAction name method
	NewActionMethod = "NewAction"

	// UpdateStatusMethod represent UpdateStatus name method
	UpdateStatusMethod = "UpdateStatus"

	// SetStatusNotifyFuncMethod represent SetStatusNotifyFunc name method
	SetStatusNotifyFuncMethod = "SetStatusNotifyFunc"

	// RequiredStatusMethod represents RequiredStatus name method
	RequiredStatusMethod = "RequiredStatus"
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

// ListenOnSetStatusNotifyFunc listening SetStatusNotifyFunc call
func (task *Task) ListenOnSetStatusNotifyFunc() *Task {
	task.On(SetStatusNotifyFuncMethod, mock.Anything).Return(nil)
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

// ListenOnRequiredStatus listening RequiredStatus call
func (task *Task) ListenOnRequiredStatus(returnErr error) *Task {
	task.On(RequiredStatusMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRequiredStatusCall RequiredStatus call assertion
func (task *Task) AssertRequiredStatusCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, RequiredStatusMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, RequiredStatusMethod, expectedCalls)
	return task
}

// ListenOnNewAction listening NewAction call
func (task *Task) ListenOnNewAction(returnChan <-chan struct{}) *Task {
	task.On(NewActionMethod, mock.Anything).Return(returnChan)
	return task
}

// AssertNewActionCall NewAction call assertion
func (task *Task) AssertNewActionCall(expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(task.t, NewActionMethod, arguments...)
	}
	task.AssertNumberOfCalls(task.t, NewActionMethod, expectedCalls)
	return task
}
