package common

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/task"
	"github.com/pastelnetwork/gonode/common/service/task/state"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/files"
)

// TaskCleanerFunc pointer to func that removes artefacts
type TaskCleanerFunc func()

// SuperNodeTask base "class" for Task
type SuperNodeTask struct {
	task.Task

	LogPrefix string
}

// RunHelper common code for Task runner
func (task *SuperNodeTask) RunHelper(ctx context.Context, clean TaskCleanerFunc) error {
	ctx = task.context(ctx)
	log.WithContext(ctx).Debug("Start task")
	defer log.WithContext(ctx).Info("Task canceled")
	defer task.Cancel()

	task.SetStatusNotifyFunc(func(status *state.Status) {
		log.WithContext(ctx).WithField("status", status.String()).Debug("States updated")
	})

	defer clean()

	return task.RunAction(ctx)
}

func (task *SuperNodeTask) context(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", task.LogPrefix, task.ID()))
}

// RemoveFile removes file from FS (TODO: move to gonode.common)
func (task *SuperNodeTask) RemoveFile(file *files.File) {
	if file != nil {
		log.Debugf("remove file: %s", file.Name())
		if err := file.Remove(); err != nil {
			log.Debugf("remove file failed: %s", err.Error())
		}
	}
}

// NewSuperNodeTask returns a new Task instance.
func NewSuperNodeTask(logPrefix string, historyDB storage.LocalStoreInterface) *SuperNodeTask {
	snt := &SuperNodeTask{
		Task:      task.New(StatusTaskStarted),
		LogPrefix: logPrefix,
	}

	snt.InitialiseHistoryDB(historyDB)

	return snt
}
