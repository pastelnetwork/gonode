package storage

import "github.com/pastelnetwork/gonode/common/types"

// LocalStoreInterface is interface for local sqlite store
type LocalStoreInterface interface {
	InsertTaskHistory(history types.TaskHistory) (int, error)
	QueryTaskHistory(taskID string) (history []types.TaskHistory, err error)
}
