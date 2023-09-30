// Code generated by goa v3.13.1, DO NOT EDIT.
//
// HTTP request path constructors for the cascade service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"fmt"
)

// UploadAssetCascadePath returns the URL path to the cascade service uploadAsset HTTP endpoint.
func UploadAssetCascadePath() string {
	return "/openapi/cascade/upload"
}

// StartProcessingCascadePath returns the URL path to the cascade service startProcessing HTTP endpoint.
func StartProcessingCascadePath(fileID string) string {
	return fmt.Sprintf("/openapi/cascade/start/%v", fileID)
}

// RegisterTaskStateCascadePath returns the URL path to the cascade service registerTaskState HTTP endpoint.
func RegisterTaskStateCascadePath(taskID string) string {
	return fmt.Sprintf("/openapi/cascade/start/%v/state", taskID)
}

// GetTaskHistoryCascadePath returns the URL path to the cascade service getTaskHistory HTTP endpoint.
func GetTaskHistoryCascadePath(taskID string) string {
	return fmt.Sprintf("/openapi/cascade/%v/history", taskID)
}

// DownloadCascadePath returns the URL path to the cascade service download HTTP endpoint.
func DownloadCascadePath() string {
	return "/openapi/cascade/download"
}
