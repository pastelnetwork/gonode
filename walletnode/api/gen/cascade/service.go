// Code generated by goa v3.15.0, DO NOT EDIT.
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
	// Get the file registration details
	RegistrationDetails(context.Context, *RegistrationDetailsPayload) (res *Registration, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// APIKeyAuth implements the authorization logic for the APIKey security scheme.
	APIKeyAuth(ctx context.Context, key string, schema *security.APIKeyScheme) (context.Context, error)
}

// APIName is the name of the API as defined in the design.
const APIName = "walletnode"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "1.0"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "cascade"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [6]string{"uploadAsset", "startProcessing", "registerTaskState", "getTaskHistory", "download", "registrationDetails"}

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

type ActivationAttempt struct {
	// ID
	ID int
	// File ID
	FileID string
	// Activation Attempt At in datetime format
	ActivationAttemptAt string
	// Indicates if the activation was successful
	IsSuccessful *bool
	// Error Message
	ErrorMessage *string
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

type File struct {
	// File ID
	FileID string
	// Upload Timestamp in datetime format
	UploadTimestamp string
	// Path to the file
	Path *string
	// Index of the file
	FileIndex *string
	// Base File ID
	BaseFileID string
	// Task ID
	TaskID string
	// Registration Transaction ID
	RegTxid *string
	// Activation Transaction ID
	ActivationTxid *string
	// Required Burn Transaction Amount
	ReqBurnTxnAmount float64
	// Burn Transaction ID
	BurnTxnID *string
	// Required Amount
	ReqAmount float64
	// Indicates if the process is concluded
	IsConcluded *bool
	// Cascade Metadata Ticket ID
	CascadeMetadataTicketID string
	// UUID Key
	UUIDKey *string
	// Hash of the Original Big File
	HashOfOriginalBigFile string
	// Name of the Original Big File with Extension
	NameOfOriginalBigFileWithExt string
	// Size of the Original Big File
	SizeOfOriginalBigFile float64
	// Data Type of the Original Big File
	DataTypeOfOriginalBigFile string
	// Start Block
	StartBlock *int32
	// Done Block
	DoneBlock *int
	// List of registration attempts
	RegistrationAttempts []*RegistrationAttempt
	// List of activation attempts
	ActivationAttempts []*ActivationAttempt
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

// Registration is the result type of the cascade service registrationDetails
// method.
type Registration struct {
	// List of files
	Files []*File
}

type RegistrationAttempt struct {
	// ID
	ID int
	// File ID
	FileID string
	// Registration Started At in datetime format
	RegStartedAt string
	// Processor SNS
	ProcessorSns *string
	// Finished At in datetime format
	FinishedAt string
	// Indicates if the registration was successful
	IsSuccessful *bool
	// Error Message
	ErrorMessage *string
}

// RegistrationDetailsPayload is the payload type of the cascade service
// registrationDetails method.
type RegistrationDetailsPayload struct {
	// file ID
	FileID string
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
	// Address to use for registration fee
	SpendableAddress *string
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

// NewRegistration initializes result type Registration from viewed result type
// Registration.
func NewRegistration(vres *cascadeviews.Registration) *Registration {
	return newRegistration(vres.Projected)
}

// NewViewedRegistration initializes viewed result type Registration from
// result type Registration using the given view.
func NewViewedRegistration(res *Registration, view string) *cascadeviews.Registration {
	p := newRegistrationView(res)
	return &cascadeviews.Registration{Projected: p, View: "default"}
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

// newRegistration converts projected type Registration to service type
// Registration.
func newRegistration(vres *cascadeviews.RegistrationView) *Registration {
	res := &Registration{}
	if vres.Files != nil {
		res.Files = make([]*File, len(vres.Files))
		for i, val := range vres.Files {
			res.Files[i] = transformCascadeviewsFileViewToFile(val)
		}
	}
	return res
}

// newRegistrationView projects result type Registration to projected type
// RegistrationView using the "default" view.
func newRegistrationView(res *Registration) *cascadeviews.RegistrationView {
	vres := &cascadeviews.RegistrationView{}
	if res.Files != nil {
		vres.Files = make([]*cascadeviews.FileView, len(res.Files))
		for i, val := range res.Files {
			vres.Files[i] = transformFileToCascadeviewsFileView(val)
		}
	} else {
		vres.Files = []*cascadeviews.FileView{}
	}
	return vres
}

// transformCascadeviewsFileViewToFile builds a value of type *File from a
// value of type *cascadeviews.FileView.
func transformCascadeviewsFileViewToFile(v *cascadeviews.FileView) *File {
	if v == nil {
		return nil
	}
	res := &File{
		FileID:                       *v.FileID,
		UploadTimestamp:              *v.UploadTimestamp,
		Path:                         v.Path,
		FileIndex:                    v.FileIndex,
		BaseFileID:                   *v.BaseFileID,
		TaskID:                       *v.TaskID,
		RegTxid:                      v.RegTxid,
		ActivationTxid:               v.ActivationTxid,
		ReqBurnTxnAmount:             *v.ReqBurnTxnAmount,
		BurnTxnID:                    v.BurnTxnID,
		ReqAmount:                    *v.ReqAmount,
		IsConcluded:                  v.IsConcluded,
		CascadeMetadataTicketID:      *v.CascadeMetadataTicketID,
		UUIDKey:                      v.UUIDKey,
		HashOfOriginalBigFile:        *v.HashOfOriginalBigFile,
		NameOfOriginalBigFileWithExt: *v.NameOfOriginalBigFileWithExt,
		SizeOfOriginalBigFile:        *v.SizeOfOriginalBigFile,
		DataTypeOfOriginalBigFile:    *v.DataTypeOfOriginalBigFile,
		StartBlock:                   v.StartBlock,
		DoneBlock:                    v.DoneBlock,
	}
	if v.RegistrationAttempts != nil {
		res.RegistrationAttempts = make([]*RegistrationAttempt, len(v.RegistrationAttempts))
		for i, val := range v.RegistrationAttempts {
			res.RegistrationAttempts[i] = transformCascadeviewsRegistrationAttemptViewToRegistrationAttempt(val)
		}
	} else {
		res.RegistrationAttempts = []*RegistrationAttempt{}
	}
	if v.ActivationAttempts != nil {
		res.ActivationAttempts = make([]*ActivationAttempt, len(v.ActivationAttempts))
		for i, val := range v.ActivationAttempts {
			res.ActivationAttempts[i] = transformCascadeviewsActivationAttemptViewToActivationAttempt(val)
		}
	} else {
		res.ActivationAttempts = []*ActivationAttempt{}
	}

	return res
}

// transformCascadeviewsRegistrationAttemptViewToRegistrationAttempt builds a
// value of type *RegistrationAttempt from a value of type
// *cascadeviews.RegistrationAttemptView.
func transformCascadeviewsRegistrationAttemptViewToRegistrationAttempt(v *cascadeviews.RegistrationAttemptView) *RegistrationAttempt {
	res := &RegistrationAttempt{
		ID:           *v.ID,
		FileID:       *v.FileID,
		RegStartedAt: *v.RegStartedAt,
		ProcessorSns: v.ProcessorSns,
		FinishedAt:   *v.FinishedAt,
		IsSuccessful: v.IsSuccessful,
		ErrorMessage: v.ErrorMessage,
	}

	return res
}

// transformCascadeviewsActivationAttemptViewToActivationAttempt builds a value
// of type *ActivationAttempt from a value of type
// *cascadeviews.ActivationAttemptView.
func transformCascadeviewsActivationAttemptViewToActivationAttempt(v *cascadeviews.ActivationAttemptView) *ActivationAttempt {
	res := &ActivationAttempt{
		ID:                  *v.ID,
		FileID:              *v.FileID,
		ActivationAttemptAt: *v.ActivationAttemptAt,
		IsSuccessful:        v.IsSuccessful,
		ErrorMessage:        v.ErrorMessage,
	}

	return res
}

// transformFileToCascadeviewsFileView builds a value of type
// *cascadeviews.FileView from a value of type *File.
func transformFileToCascadeviewsFileView(v *File) *cascadeviews.FileView {
	res := &cascadeviews.FileView{
		FileID:                       &v.FileID,
		UploadTimestamp:              &v.UploadTimestamp,
		Path:                         v.Path,
		FileIndex:                    v.FileIndex,
		BaseFileID:                   &v.BaseFileID,
		TaskID:                       &v.TaskID,
		RegTxid:                      v.RegTxid,
		ActivationTxid:               v.ActivationTxid,
		ReqBurnTxnAmount:             &v.ReqBurnTxnAmount,
		BurnTxnID:                    v.BurnTxnID,
		ReqAmount:                    &v.ReqAmount,
		IsConcluded:                  v.IsConcluded,
		CascadeMetadataTicketID:      &v.CascadeMetadataTicketID,
		UUIDKey:                      v.UUIDKey,
		HashOfOriginalBigFile:        &v.HashOfOriginalBigFile,
		NameOfOriginalBigFileWithExt: &v.NameOfOriginalBigFileWithExt,
		SizeOfOriginalBigFile:        &v.SizeOfOriginalBigFile,
		DataTypeOfOriginalBigFile:    &v.DataTypeOfOriginalBigFile,
		StartBlock:                   v.StartBlock,
		DoneBlock:                    v.DoneBlock,
	}
	if v.RegistrationAttempts != nil {
		res.RegistrationAttempts = make([]*cascadeviews.RegistrationAttemptView, len(v.RegistrationAttempts))
		for i, val := range v.RegistrationAttempts {
			res.RegistrationAttempts[i] = transformRegistrationAttemptToCascadeviewsRegistrationAttemptView(val)
		}
	} else {
		res.RegistrationAttempts = []*cascadeviews.RegistrationAttemptView{}
	}
	if v.ActivationAttempts != nil {
		res.ActivationAttempts = make([]*cascadeviews.ActivationAttemptView, len(v.ActivationAttempts))
		for i, val := range v.ActivationAttempts {
			res.ActivationAttempts[i] = transformActivationAttemptToCascadeviewsActivationAttemptView(val)
		}
	} else {
		res.ActivationAttempts = []*cascadeviews.ActivationAttemptView{}
	}

	return res
}

// transformRegistrationAttemptToCascadeviewsRegistrationAttemptView builds a
// value of type *cascadeviews.RegistrationAttemptView from a value of type
// *RegistrationAttempt.
func transformRegistrationAttemptToCascadeviewsRegistrationAttemptView(v *RegistrationAttempt) *cascadeviews.RegistrationAttemptView {
	res := &cascadeviews.RegistrationAttemptView{
		ID:           &v.ID,
		FileID:       &v.FileID,
		RegStartedAt: &v.RegStartedAt,
		ProcessorSns: v.ProcessorSns,
		FinishedAt:   &v.FinishedAt,
		IsSuccessful: v.IsSuccessful,
		ErrorMessage: v.ErrorMessage,
	}

	return res
}

// transformActivationAttemptToCascadeviewsActivationAttemptView builds a value
// of type *cascadeviews.ActivationAttemptView from a value of type
// *ActivationAttempt.
func transformActivationAttemptToCascadeviewsActivationAttemptView(v *ActivationAttempt) *cascadeviews.ActivationAttemptView {
	res := &cascadeviews.ActivationAttemptView{
		ID:                  &v.ID,
		FileID:              &v.FileID,
		ActivationAttemptAt: &v.ActivationAttemptAt,
		IsSuccessful:        v.IsSuccessful,
		ErrorMessage:        v.ErrorMessage,
	}

	return res
}
