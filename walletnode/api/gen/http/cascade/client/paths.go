// Code generated by goa v3.4.3, DO NOT EDIT.
//
// HTTP request path constructors for the cascade service.
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	"fmt"
)

// UploadImageCascadePath returns the URL path to the cascade service uploadImage HTTP endpoint.
func UploadImageCascadePath() string {
	return "/openapi/cascade/upload"
}

// ActionDetailsCascadePath returns the URL path to the cascade service actionDetails HTTP endpoint.
func ActionDetailsCascadePath(imageID string) string {
	return fmt.Sprintf("/openapi/cascade/details/%v", imageID)
}

// StartProcessingCascadePath returns the URL path to the cascade service startProcessing HTTP endpoint.
func StartProcessingCascadePath(imageID string) string {
	return fmt.Sprintf("/openapi/cascade/start/%v", imageID)
}

// RegisterTaskStateCascadePath returns the URL path to the cascade service registerTaskState HTTP endpoint.
func RegisterTaskStateCascadePath(taskID string) string {
	return fmt.Sprintf("/openapi/cascade/start/%v/state", taskID)
}
