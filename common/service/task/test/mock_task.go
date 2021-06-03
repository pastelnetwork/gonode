package test

import (
	"github.com/pastelnetwork/gonode/common/service/task/mocks"
	"github.com/stretchr/testify/mock"
)

//Task implementing task.Task mock for testing purpose
type Task struct {
	TaskMock *mocks.Task
}

// NewMockTask new Task instance
func NewMockTask() *Task {
	return &Task{
		TaskMock: &mocks.Task{},
	}
}

// ListenOnID listening ID call and returns task id from args
func (t *Task) ListenOnID(id string) *Task {
	t.TaskMock.On("ID").Return(id)
	return t
}

// ListenOnUpdateStatus listening UpdateStatus call
func (t *Task) ListenOnUpdateStatus() *Task {
	t.TaskMock.On("UpdateStatus", mock.Anything).Return(nil)
	return t
}

// ListenOnSetStatusNotifyFunc listening SetStatusNotifyFunc call
func (t *Task) ListenOnSetStatusNotifyFunc() *Task {
	t.TaskMock.On("SetStatusNotifyFunc", mock.Anything).Return(nil)
	return t
}
