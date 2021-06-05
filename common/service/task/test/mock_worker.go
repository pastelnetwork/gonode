package test

import (
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/mocks"
	"github.com/stretchr/testify/mock"
)

// Worker implementing task.Worker mock for testing purpose
type Worker struct {
	WorkerMock *mocks.Worker
}

// NewMockWorker new Worker instance
func NewMockWorker() *Worker {
	return &Worker{
		WorkerMock: &mocks.Worker{},
	}

}

// ListenOnTasks listening Tasks call and returns tasks from args
func (w *Worker) ListenOnTasks(tasks []task.Task) *Worker {
	w.WorkerMock.On("Tasks").Return(tasks)
	return w
}

// ListenOnTask listening Task call and return task from args
func (w *Worker) ListenOnTask(task task.Task) *Worker {
	w.WorkerMock.On("Task", mock.AnythingOfType("string")).Return(task)
	return w
}

// ListenOnAddTask listening AddTask call
func (w *Worker) ListenOnAddTask() *Worker {
	w.WorkerMock.On("AddTask", mock.Anything).Return(nil)
	return w
}

// ListenOnRemoveTask listening RemoveTask call
func (w *Worker) ListenOnRemoveTask() *Worker {
	w.WorkerMock.On("RemoveTask", mock.Anything).Return(nil)
	return w
}

// ListenOnRun listening Run call and returns error from args
func (w *Worker) ListenOnRun(returnErr error) *Worker {
	w.WorkerMock.On("Run", mock.Anything).Return(returnErr)
	return w
}
