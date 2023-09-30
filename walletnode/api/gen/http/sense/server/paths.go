// Code generated by goa v3.13.1, DO NOT EDIT.
//
// HTTP request path constructors for the sense service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"fmt"
)

// UploadImageSensePath returns the URL path to the sense service uploadImage HTTP endpoint.
func UploadImageSensePath() string {
	return "/openapi/sense/upload"
}

// StartProcessingSensePath returns the URL path to the sense service startProcessing HTTP endpoint.
func StartProcessingSensePath(imageID string) string {
	return fmt.Sprintf("/openapi/sense/start/%v", imageID)
}

// RegisterTaskStateSensePath returns the URL path to the sense service registerTaskState HTTP endpoint.
func RegisterTaskStateSensePath(taskID string) string {
	return fmt.Sprintf("/openapi/sense/start/%v/state", taskID)
}

// GetTaskHistorySensePath returns the URL path to the sense service getTaskHistory HTTP endpoint.
func GetTaskHistorySensePath(taskID string) string {
	return fmt.Sprintf("/openapi/sense/%v/history", taskID)
}

// DownloadSensePath returns the URL path to the sense service download HTTP endpoint.
func DownloadSensePath() string {
	return "/openapi/sense/download"
}
