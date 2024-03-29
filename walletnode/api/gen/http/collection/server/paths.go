// Code generated by goa v3.15.0, DO NOT EDIT.
//
// HTTP request path constructors for the collection service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"fmt"
)

// RegisterCollectionCollectionPath returns the URL path to the collection service registerCollection HTTP endpoint.
func RegisterCollectionCollectionPath() string {
	return "/collection/register"
}

// RegisterTaskStateCollectionPath returns the URL path to the collection service registerTaskState HTTP endpoint.
func RegisterTaskStateCollectionPath(taskID string) string {
	return fmt.Sprintf("/collection/%v/state", taskID)
}

// GetTaskHistoryCollectionPath returns the URL path to the collection service getTaskHistory HTTP endpoint.
func GetTaskHistoryCollectionPath(taskID string) string {
	return fmt.Sprintf("/collection/%v/history", taskID)
}
