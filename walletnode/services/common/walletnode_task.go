package common

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
)

// TaskCleanerFunc pointer to cleaner function
type TaskCleanerFunc func()

// TaskRunnerFunc pointer to runner function
type TaskRunnerFunc func(ctx context.Context) error

// WalletNodeTask represents base WalletNode task
type WalletNodeTask struct {
	task.Task

	LogPrefix string
	StatusLog types.Fields
	err       error
}

// RunHelper starts the task
func (task *WalletNodeTask) RunHelper(ctx context.Context, run TaskRunnerFunc, clean TaskCleanerFunc) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", task.LogPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	defer clean()
	if err := run(ctx); err != nil {
		task.err = err
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithErrorStack(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

// RemoveFile removes file from FS (TODO: move to gonode.common)
func (task *WalletNodeTask) RemoveFile(file *files.File) {
	if file != nil {
		log.Debugf("remove file: %s", file.Name())
		if err := file.Remove(); err != nil {
			log.Debugf("remove file failed: %s", err.Error())
		}
	}
}

// Error returns the task err
func (task *WalletNodeTask) Error() error {
	return task.err
}

// NewWalletNodeTask returns a new WalletNodeTask instance.
func NewWalletNodeTask(logPrefix string) *WalletNodeTask {
	wnt := &WalletNodeTask{
		Task:      task.New(StatusTaskStarted),
		LogPrefix: logPrefix,
		StatusLog: make(map[string]interface{}),
	}
	wnt.SetStateLog(wnt.StatusLog)

	return wnt
}
