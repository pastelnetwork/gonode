// Code generated by goa v3.13.1, DO NOT EDIT.
//
// sense service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package sense

import (
	"context"

	senseviews "github.com/pastelnetwork/gonode/walletnode/api/gen/sense/views"
	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// OpenAPI Sense service
type Service interface {
	// Upload the image
	UploadImage(context.Context, *UploadImagePayload) (res *Image, err error)
	// Start processing the image
	StartProcessing(context.Context, *StartProcessingPayload) (res *StartProcessingResult, err error)
	// Streams the state of the registration process.
	RegisterTaskState(context.Context, *RegisterTaskStatePayload, RegisterTaskStateServerStream) (err error)
	// Gets the history of the task's states.
	GetTaskHistory(context.Context, *GetTaskHistoryPayload) (res []*TaskHistory, err error)
	// Download sense result; duplication detection results file.
	Download(context.Context, *DownloadPayload) (res *DownloadResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "sense"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [5]string{"uploadImage", "startProcessing", "registerTaskState", "getTaskHistory", "download"}

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

type Details struct {
	// details regarding the status
	Message *string
	// important fields regarding status history
	Fields map[string]any
}

// DownloadPayload is the payload type of the sense service download method.
type DownloadPayload struct {
	// Nft Registration Request transaction ID
	Txid string
	// Owner's PastelID
	Pid string
	// Passphrase of the owner's PastelID
	Key string
}

// DownloadResult is the result type of the sense service download method.
type DownloadResult struct {
	// File downloaded
	File []byte
}

// GetTaskHistoryPayload is the payload type of the sense service
// getTaskHistory method.
type GetTaskHistoryPayload struct {
	// Task ID of the registration process
	TaskID string
}

// Image is the result type of the sense service uploadImage method.
type Image struct {
	// Uploaded image ID
	ImageID string
	// Image expiration
	ExpiresIn string
	// Estimated fee
	TotalEstimatedFee float64
	// The amount that's required to be preburned
	RequiredPreburnAmount float64
}

// RegisterTaskStatePayload is the payload type of the sense service
// registerTaskState method.
type RegisterTaskStatePayload struct {
	// Task ID of the registration process
	TaskID string
}

// StartProcessingPayload is the payload type of the sense service
// startProcessing method.
type StartProcessingPayload struct {
	// Uploaded image ID
	ImageID string
	// Burn transaction ID
	BurnTxid string
	// Act Collection TxID to add given ticket in collection
	CollectionActTxid *string
	// OpenAPI GroupID string
	OpenAPIGroupID string
	// App PastelID
	AppPastelID string
	// Passphrase of the owner's PastelID
	Key string
}

// StartProcessingResult is the result type of the sense service
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

// TaskState is the result type of the sense service registerTaskState method.
type TaskState struct {
	// Date of the status creation
	Date string
	// Status of the registration process
	Status string
}

// UploadImagePayload is the payload type of the sense service uploadImage
// method.
type UploadImagePayload struct {
	// File to upload
	Bytes []byte
	// For internal use
	Filename *string
}

// MakeUnAuthorized builds a goa.ServiceError from an error.
func MakeUnAuthorized(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "UnAuthorized", false, false, false)
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

// NewImage initializes result type Image from viewed result type Image.
func NewImage(vres *senseviews.Image) *Image {
	return newImage(vres.Projected)
}

// NewViewedImage initializes viewed result type Image from result type Image
// using the given view.
func NewViewedImage(res *Image, view string) *senseviews.Image {
	p := newImageView(res)
	return &senseviews.Image{Projected: p, View: "default"}
}

// NewStartProcessingResult initializes result type StartProcessingResult from
// viewed result type StartProcessingResult.
func NewStartProcessingResult(vres *senseviews.StartProcessingResult) *StartProcessingResult {
	return newStartProcessingResult(vres.Projected)
}

// NewViewedStartProcessingResult initializes viewed result type
// StartProcessingResult from result type StartProcessingResult using the given
// view.
func NewViewedStartProcessingResult(res *StartProcessingResult, view string) *senseviews.StartProcessingResult {
	p := newStartProcessingResultView(res)
	return &senseviews.StartProcessingResult{Projected: p, View: "default"}
}

// newImage converts projected type Image to service type Image.
func newImage(vres *senseviews.ImageView) *Image {
	res := &Image{}
	if vres.ImageID != nil {
		res.ImageID = *vres.ImageID
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

// newImageView projects result type Image to projected type ImageView using
// the "default" view.
func newImageView(res *Image) *senseviews.ImageView {
	vres := &senseviews.ImageView{
		ImageID:               &res.ImageID,
		ExpiresIn:             &res.ExpiresIn,
		TotalEstimatedFee:     &res.TotalEstimatedFee,
		RequiredPreburnAmount: &res.RequiredPreburnAmount,
	}
	return vres
}

// newStartProcessingResult converts projected type StartProcessingResult to
// service type StartProcessingResult.
func newStartProcessingResult(vres *senseviews.StartProcessingResultView) *StartProcessingResult {
	res := &StartProcessingResult{}
	if vres.TaskID != nil {
		res.TaskID = *vres.TaskID
	}
	return res
}

// newStartProcessingResultView projects result type StartProcessingResult to
// projected type StartProcessingResultView using the "default" view.
func newStartProcessingResultView(res *StartProcessingResult) *senseviews.StartProcessingResultView {
	vres := &senseviews.StartProcessingResultView{
		TaskID: &res.TaskID,
	}
	return vres
}
