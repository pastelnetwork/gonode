package test

import (
	"testing"

	"github.com/pastelnetwork/gonode/common/service/task/mocks"
	"github.com/stretchr/testify/mock"
)

//Task implementing task.Task mock for testing purpose
type Task struct {
	*mocks.Task
	IDMethod                  string
	runMethod                 string
	cancelMethod              string
	runActionMethod           string
	updateStatusMethod        string
	setStatusNotifyFuncMethod string
}

// NewMockTask new Task instance
func NewMockTask() *Task {
	return &Task{
		Task:                      &mocks.Task{},
		IDMethod:                  "ID",
		runMethod:                 "Run",
		cancelMethod:              "Cancel",
		runActionMethod:           "RunAction",
		updateStatusMethod:        "UpdateStatus",
		setStatusNotifyFuncMethod: "SetStatusNotifyFunc",
	}
}

// ListenOnID listening ID call and returns task id from args
func (task *Task) ListenOnID(id string) *Task {
	task.On(task.IDMethod).Return(id)
	return task
}

// AssertIDCall ID call method assertion
func (task *Task) AssertIDCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.IDMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.IDMethod, expectedCalls)
	return task
}

// ListenOnRun listening Run call and returns error from args
func (task *Task) ListenOnRun(returnErr error) *Task {
	task.On(task.runMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRunCall Run call assertion
func (task *Task) AssertRunCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.runMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.runMethod, expectedCalls)
	return task
}

// ListenOnCancel listening Cancel call
func (task *Task) ListenOnCancel() *Task {
	task.On(task.cancelMethod).Return(nil)
	return task
}

// AssertCancelCall Cancel call assertion
func (task *Task) AssertCancelCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.cancelMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.cancelMethod, expectedCalls)
	return task
}

// ListenOnRunAction listening RunAction call and returns error from args
func (task *Task) ListenOnRunAction(returnErr error) *Task {
	task.On(task.runActionMethod, mock.Anything).Return(returnErr)
	return task
}

// AssertRunActionCall RunAction call assertion
func (task *Task) AssertRunActionCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.runActionMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.runActionMethod, expectedCalls)
	return task
}

// ListenOnUpdateStatus listening UpdateStatus call
func (task *Task) ListenOnUpdateStatus() *Task {
	task.On(task.updateStatusMethod, mock.Anything).Return(nil)
	return task
}

// AssertUpdateStatusCall UpdateStatus call assertion
func (task *Task) AssertUpdateStatusCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.updateStatusMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.updateStatusMethod, expectedCalls)
	return task
}

// ListenOnSetStatusNotifyFunc listening SetStatusNotifyFunc call
func (task *Task) ListenOnSetStatusNotifyFunc() *Task {
	task.On(task.setStatusNotifyFuncMethod, mock.Anything).Return(nil)
	return task
}

// AssertSetStatusNotifyFuncCall SetSTatusNotifyFunc call assertion
func (task *Task) AssertSetStatusNotifyFuncCall(t *testing.T, expectedCalls int, arguments ...interface{}) *Task {
	if expectedCalls > 0 {
		task.AssertCalled(t, task.setStatusNotifyFuncMethod, arguments...)
	}
	task.AssertNumberOfCalls(t, task.setStatusNotifyFuncMethod, expectedCalls)
	return task
}
