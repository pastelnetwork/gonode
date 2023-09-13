package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/pastelnetwork/gonode/common/blocktracker"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/storage/files"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
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

	err := run(ctx)
	if err != nil {
		task.err = err
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithError(err).Error("Task is rejected")
		return nil
	}

	task.UpdateStatus(StatusTaskCompleted)
	return nil
}

// RunHelper with Retry
func (task *WalletNodeTask) RunHelperWithRetry(ctx context.Context, run TaskRunnerFunc, clean TaskCleanerFunc, maxRetries int, pt pastel.Client, resetFunc func(), setErrFunc func(err error)) error {
	ctx = log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", task.LogPrefix, task.ID()))

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status).Debug("States updated")
	})

	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Debug("End task")

	defer clean()

	retries := 0
	blockTracker := blocktracker.New(pt)
	baseBlkCnt, err := blockTracker.GetBlockCount()
	if err != nil {
		log.WithContext(ctx).WithError(err).Warn("failed to get block count")
		return nil
	}

	err = run(ctx)
	for doRetry(err) && retries < maxRetries {
		retries++
		log.WithContext(ctx).WithError(err).WithField("current-try", retries).WithField("base-block", baseBlkCnt).Error("Task is failed, retrying")
		resetFunc()

		if err := blockTracker.WaitTillNextBlock(ctx, baseBlkCnt); err != nil {
			log.WithContext(ctx).WithError(err).Warn("failed to wait till next block")
			return nil
		}

		err = run(ctx)
	}

	if err != nil {
		task.UpdateStatus(StatusTaskRejected)
		log.WithContext(ctx).WithError(err).Error("Task is rejected")
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

// SetError returns the task err
func (task *WalletNodeTask) SetError(err error) {
	task.err = err
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

func doRetry(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "retry")
}
