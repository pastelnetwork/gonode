package storage

import "github.com/pastelnetwork/gonode/common/types"

type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
}
