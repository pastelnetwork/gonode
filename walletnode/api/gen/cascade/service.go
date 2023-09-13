// Code generated by goa v3.13.0, DO NOT EDIT.
//
// cascade service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package cascade

import (
	"context"

	cascadeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// OpenAPI Cascade service
type Service interface {
	// Upload the asset file
	UploadAsset(context.Context, *UploadAssetPayload) (res *Asset, err error)
	// Start processing the image
	StartProcessing(context.Context, *StartProcessingPayload) (res *StartProcessingResult, err error)
	// Streams the state of the registration process.
	RegisterTaskState(context.Context, *RegisterTaskStatePayload, RegisterTaskStateServerStream) (err error)
	// Gets the history of the task's states.
	GetTaskHistory(context.Context, *GetTaskHistoryPayload) (res []*TaskHistory, err error)
	// Download cascade Artifact.
	Download(context.Context, *DownloadPayload) (res *FileDownloadResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "cascade"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [5]string{"uploadAsset", "startProcessing", "registerTaskState", "getTaskHistory", "download"}

// RegisterTaskStateServerStream is the interface a "registerTaskState"
// endpoint server stream must satisfy.
type RegisterTaskStateServerStream interface {
	// Send streams instances of "TaskState".
	Send(*TaskState) error
	// Close closes the stream.
	Close() error
}

// RegisterTaskStateClientStream is the interface a "registerTaskState"
// endpoint client stream must satisfy.
type RegisterTaskStateClientStream interface {
	// Recv reads instances of "TaskState" from the stream.
	Recv() (*TaskState, error)
}

// Asset is the result type of the cascade service uploadAsset method.
type Asset struct {
	// Uploaded file ID
	FileID string
	// File expiration
	ExpiresIn string
	// Estimated fee
	TotalEstimatedFee float64
	// The amount that's required to be preburned
	RequiredPreburnAmount float64
}

type Details struct {
	// details regarding the status
	Message *string
	// important fields regarding status history
	Fields map[string]any
}

// DownloadPayload is the payload type of the cascade service download method.
type DownloadPayload struct {
	// Nft Registration Request transaction ID
	Txid string
	// Owner's PastelID
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// FileDownloadResult is the result type of the cascade service download method.
type FileDownloadResult struct {
	// File path
	FileID string
}

// GetTaskHistoryPayload is the payload type of the cascade service
// getTaskHistory method.
type GetTaskHistoryPayload struct {
	// Task ID of the registration process
	TaskID string
}

// RegisterTaskStatePayload is the payload type of the cascade service
// registerTaskState method.
type RegisterTaskStatePayload struct {
	// Task ID of the registration process
	TaskID string
}

// StartProcessingPayload is the payload type of the cascade service
// startProcessing method.
type StartProcessingPayload struct {
	// Uploaded asset file ID
	FileID string
	// Burn transaction ID
	BurnTxid string
	// App PastelID
	AppPastelID string
	// To make it publicly accessible
	MakePubliclyAccessible bool
	// Passphrase of the owner's PastelID
	Key string
}

// StartProcessingResult is the result type of the cascade service
// startProcessing method.
type StartProcessingResult struct {
	// Task ID of processing task
	TaskID string
}

type TaskHistory struct {
	// Timestamp of the status creation
	Timestamp *string
	// past status string
	Status string
	// message string (if any)
	Message *string
	// details of the status
	Details *Details
}

// TaskState is the result type of the cascade service registerTaskState method.
type TaskState struct {
	// Date of the status creation
	Date string
	// Status of the registration process
	Status string
}

// UploadAssetPayload is the payload type of the cascade service uploadAsset
// method.
type UploadAssetPayload struct {
	// File to upload
	Bytes []byte
	// For internal use
	Filename *string
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "BadRequest", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "NotFound", false, false, false)
}

// MakeInternalServerError builds a goa.ServiceError from an error.
func MakeInternalServerError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "InternalServerError", false, false, false)
}

// NewAsset initializes result type Asset from viewed result type Asset.
func NewAsset(vres *cascadeviews.Asset) *Asset {
	return newAsset(vres.Projected)
}

// NewViewedAsset initializes viewed result type Asset from result type Asset
// using the given view.
func NewViewedAsset(res *Asset, view string) *cascadeviews.Asset {
	p := newAssetView(res)
	return &cascadeviews.Asset{Projected: p, View: "default"}
}

// NewStartProcessingResult initializes result type StartProcessingResult from
// viewed result type StartProcessingResult.
func NewStartProcessingResult(vres *cascadeviews.StartProcessingResult) *StartProcessingResult {
	return newStartProcessingResult(vres.Projected)
}

// NewViewedStartProcessingResult initializes viewed result type
// StartProcessingResult from result type StartProcessingResult using the given
// view.
func NewViewedStartProcessingResult(res *StartProcessingResult, view string) *cascadeviews.StartProcessingResult {
	p := newStartProcessingResultView(res)
	return &cascadeviews.StartProcessingResult{Projected: p, View: "default"}
}

// newAsset converts projected type Asset to service type Asset.
func newAsset(vres *cascadeviews.AssetView) *Asset {
	res := &Asset{}
	if vres.FileID != nil {
		res.FileID = *vres.FileID
	}
	if vres.ExpiresIn != nil {
		res.ExpiresIn = *vres.ExpiresIn
	}
	if vres.TotalEstimatedFee != nil {
		res.TotalEstimatedFee = *vres.TotalEstimatedFee
	}
	if vres.RequiredPreburnAmount != nil {
		res.RequiredPreburnAmount = *vres.RequiredPreburnAmount
	}
	if vres.TotalEstimatedFee == nil {
		res.TotalEstimatedFee = 1
	}
	if vres.RequiredPreburnAmount == nil {
		res.RequiredPreburnAmount = 1
	}
	return res
}

// newAssetView projects result type Asset to projected type AssetView using
// the "default" view.
func newAssetView(res *Asset) *cascadeviews.AssetView {
	vres := &cascadeviews.AssetView{
		FileID:                &res.FileID,
		ExpiresIn:             &res.ExpiresIn,
		TotalEstimatedFee:     &res.TotalEstimatedFee,
		RequiredPreburnAmount: &res.RequiredPreburnAmount,
	}
	return vres
}

// newStartProcessingResult converts projected type StartProcessingResult to
// service type StartProcessingResult.
func newStartProcessingResult(vres *cascadeviews.StartProcessingResultView) *StartProcessingResult {
	res := &StartProcessingResult{}
	if vres.TaskID != nil {
		res.TaskID = *vres.TaskID
	}
	return res
}

// newStartProcessingResultView projects result type StartProcessingResult to
// projected type StartProcessingResultView using the "default" view.
func newStartProcessingResultView(res *StartProcessingResult) *cascadeviews.StartProcessingResultView {
	vres := &cascadeviews.StartProcessingResultView{
		TaskID: &res.TaskID,
	}
	return vres
}
