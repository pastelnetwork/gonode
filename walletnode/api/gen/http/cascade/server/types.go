// Code generated by goa v3.7.6, DO NOT EDIT.
//
// cascade HTTP server types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"unicode/utf8"

	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	cascadeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade/views"
	goa "goa.design/goa/v3/pkg"
)

// UploadAssetRequestBody is the type of the "cascade" service "uploadAsset"
// endpoint HTTP request body.
type UploadAssetRequestBody struct {
	// File to upload
	Bytes []byte `form:"file,omitempty" json:"file,omitempty" xml:"file,omitempty"`
	// For internal use
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// StartProcessingRequestBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP request body.
type StartProcessingRequestBody struct {
	// Burn transaction ID
	BurnTxid *string `form:"burn_txid,omitempty" json:"burn_txid,omitempty" xml:"burn_txid,omitempty"`
	// App PastelID
	AppPastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
	// To make it publicly accessible
	MakePubliclyAccessible *bool `form:"make_publicly_accessible,omitempty" json:"make_publicly_accessible,omitempty" xml:"make_publicly_accessible,omitempty"`
}

// UploadAssetResponseBody is the type of the "cascade" service "uploadAsset"
// endpoint HTTP response body.
type UploadAssetResponseBody struct {
	// Uploaded file ID
	FileID string `form:"file_id" json:"file_id" xml:"file_id"`
	// File expiration
	ExpiresIn string `form:"expires_in" json:"expires_in" xml:"expires_in"`
	// Estimated fee
	TotalEstimatedFee float64 `form:"total_estimated_fee" json:"total_estimated_fee" xml:"total_estimated_fee"`
	// The amount that's required to be preburned
	RequiredPreburnAmount float64 `form:"required_preburn_amount" json:"required_preburn_amount" xml:"required_preburn_amount"`
}

// StartProcessingResponseBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP response body.
type StartProcessingResponseBody struct {
	// Task ID of processing task
	TaskID string `form:"task_id" json:"task_id" xml:"task_id"`
}

// RegisterTaskStateResponseBody is the type of the "cascade" service
// "registerTaskState" endpoint HTTP response body.
type RegisterTaskStateResponseBody struct {
	// Date of the status creation
	Date string `form:"date" json:"date" xml:"date"`
	// Status of the registration process
	Status string `form:"status" json:"status" xml:"status"`
}

// GetTaskHistoryResponseBody is the type of the "cascade" service
// "getTaskHistory" endpoint HTTP response body.
type GetTaskHistoryResponseBody []*TaskHistoryResponse

// UploadAssetBadRequestResponseBody is the type of the "cascade" service
// "uploadAsset" endpoint HTTP response body for the "BadRequest" error.
type UploadAssetBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// UploadAssetInternalServerErrorResponseBody is the type of the "cascade"
// service "uploadAsset" endpoint HTTP response body for the
// "InternalServerError" error.
type UploadAssetInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// StartProcessingBadRequestResponseBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP response body for the "BadRequest" error.
type StartProcessingBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// StartProcessingInternalServerErrorResponseBody is the type of the "cascade"
// service "startProcessing" endpoint HTTP response body for the
// "InternalServerError" error.
type StartProcessingInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterTaskStateNotFoundResponseBody is the type of the "cascade" service
// "registerTaskState" endpoint HTTP response body for the "NotFound" error.
type RegisterTaskStateNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// RegisterTaskStateInternalServerErrorResponseBody is the type of the
// "cascade" service "registerTaskState" endpoint HTTP response body for the
// "InternalServerError" error.
type RegisterTaskStateInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// GetTaskHistoryNotFoundResponseBody is the type of the "cascade" service
// "getTaskHistory" endpoint HTTP response body for the "NotFound" error.
type GetTaskHistoryNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// GetTaskHistoryInternalServerErrorResponseBody is the type of the "cascade"
// service "getTaskHistory" endpoint HTTP response body for the
// "InternalServerError" error.
type GetTaskHistoryInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DownloadNotFoundResponseBody is the type of the "cascade" service "download"
// endpoint HTTP response body for the "NotFound" error.
type DownloadNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// DownloadInternalServerErrorResponseBody is the type of the "cascade" service
// "download" endpoint HTTP response body for the "InternalServerError" error.
type DownloadInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name string `form:"name" json:"name" xml:"name"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID string `form:"id" json:"id" xml:"id"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message string `form:"message" json:"message" xml:"message"`
	// Is the error temporary?
	Temporary bool `form:"temporary" json:"temporary" xml:"temporary"`
	// Is the error a timeout?
	Timeout bool `form:"timeout" json:"timeout" xml:"timeout"`
	// Is the error a server-side fault?
	Fault bool `form:"fault" json:"fault" xml:"fault"`
}

// TaskHistoryResponse is used to define fields on response body types.
type TaskHistoryResponse struct {
	// Timestamp of the status creation
	Timestamp *string `form:"timestamp,omitempty" json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	// past status string
	Status string `form:"status" json:"status" xml:"status"`
	// message string (if any)
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// details of the status
	Details *DetailsResponse `form:"details,omitempty" json:"details,omitempty" xml:"details,omitempty"`
}

// DetailsResponse is used to define fields on response body types.
type DetailsResponse struct {
	// details regarding the status
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// important fields regarding status history
	Fields map[string]interface{} `form:"fields,omitempty" json:"fields,omitempty" xml:"fields,omitempty"`
}

// NewUploadAssetResponseBody builds the HTTP response body from the result of
// the "uploadAsset" endpoint of the "cascade" service.
func NewUploadAssetResponseBody(res *cascadeviews.AssetView) *UploadAssetResponseBody {
	body := &UploadAssetResponseBody{
		FileID:            *res.FileID,
		ExpiresIn:         *res.ExpiresIn,
		TotalEstimatedFee: *res.TotalEstimatedFee,
	}
	if res.RequiredPreburnAmount != nil {
		body.RequiredPreburnAmount = *res.RequiredPreburnAmount
	}
	if res.RequiredPreburnAmount == nil {
		body.RequiredPreburnAmount = 1
	}
	return body
}

// NewStartProcessingResponseBody builds the HTTP response body from the result
// of the "startProcessing" endpoint of the "cascade" service.
func NewStartProcessingResponseBody(res *cascadeviews.StartProcessingResultView) *StartProcessingResponseBody {
	body := &StartProcessingResponseBody{
		TaskID: *res.TaskID,
	}
	return body
}

// NewRegisterTaskStateResponseBody builds the HTTP response body from the
// result of the "registerTaskState" endpoint of the "cascade" service.
func NewRegisterTaskStateResponseBody(res *cascade.TaskState) *RegisterTaskStateResponseBody {
	body := &RegisterTaskStateResponseBody{
		Date:   res.Date,
		Status: res.Status,
	}
	return body
}

// NewGetTaskHistoryResponseBody builds the HTTP response body from the result
// of the "getTaskHistory" endpoint of the "cascade" service.
func NewGetTaskHistoryResponseBody(res []*cascade.TaskHistory) GetTaskHistoryResponseBody {
	body := make([]*TaskHistoryResponse, len(res))
	for i, val := range res {
		body[i] = marshalCascadeTaskHistoryToTaskHistoryResponse(val)
	}
	return body
}

// NewUploadAssetBadRequestResponseBody builds the HTTP response body from the
// result of the "uploadAsset" endpoint of the "cascade" service.
func NewUploadAssetBadRequestResponseBody(res *goa.ServiceError) *UploadAssetBadRequestResponseBody {
	body := &UploadAssetBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUploadAssetInternalServerErrorResponseBody builds the HTTP response body
// from the result of the "uploadAsset" endpoint of the "cascade" service.
func NewUploadAssetInternalServerErrorResponseBody(res *goa.ServiceError) *UploadAssetInternalServerErrorResponseBody {
	body := &UploadAssetInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStartProcessingBadRequestResponseBody builds the HTTP response body from
// the result of the "startProcessing" endpoint of the "cascade" service.
func NewStartProcessingBadRequestResponseBody(res *goa.ServiceError) *StartProcessingBadRequestResponseBody {
	body := &StartProcessingBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStartProcessingInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "startProcessing" endpoint of the "cascade"
// service.
func NewStartProcessingInternalServerErrorResponseBody(res *goa.ServiceError) *StartProcessingInternalServerErrorResponseBody {
	body := &StartProcessingInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterTaskStateNotFoundResponseBody builds the HTTP response body from
// the result of the "registerTaskState" endpoint of the "cascade" service.
func NewRegisterTaskStateNotFoundResponseBody(res *goa.ServiceError) *RegisterTaskStateNotFoundResponseBody {
	body := &RegisterTaskStateNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewRegisterTaskStateInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "registerTaskState" endpoint of the "cascade"
// service.
func NewRegisterTaskStateInternalServerErrorResponseBody(res *goa.ServiceError) *RegisterTaskStateInternalServerErrorResponseBody {
	body := &RegisterTaskStateInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetTaskHistoryNotFoundResponseBody builds the HTTP response body from the
// result of the "getTaskHistory" endpoint of the "cascade" service.
func NewGetTaskHistoryNotFoundResponseBody(res *goa.ServiceError) *GetTaskHistoryNotFoundResponseBody {
	body := &GetTaskHistoryNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetTaskHistoryInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "getTaskHistory" endpoint of the "cascade"
// service.
func NewGetTaskHistoryInternalServerErrorResponseBody(res *goa.ServiceError) *GetTaskHistoryInternalServerErrorResponseBody {
	body := &GetTaskHistoryInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadNotFoundResponseBody builds the HTTP response body from the
// result of the "download" endpoint of the "cascade" service.
func NewDownloadNotFoundResponseBody(res *goa.ServiceError) *DownloadNotFoundResponseBody {
	body := &DownloadNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewDownloadInternalServerErrorResponseBody builds the HTTP response body
// from the result of the "download" endpoint of the "cascade" service.
func NewDownloadInternalServerErrorResponseBody(res *goa.ServiceError) *DownloadInternalServerErrorResponseBody {
	body := &DownloadInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUploadAssetPayload builds a cascade service uploadAsset endpoint payload.
func NewUploadAssetPayload(body *UploadAssetRequestBody) *cascade.UploadAssetPayload {
	v := &cascade.UploadAssetPayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v
}

// NewStartProcessingPayload builds a cascade service startProcessing endpoint
// payload.
func NewStartProcessingPayload(body *StartProcessingRequestBody, fileID string, key string) *cascade.StartProcessingPayload {
	v := &cascade.StartProcessingPayload{
		BurnTxid:    *body.BurnTxid,
		AppPastelID: *body.AppPastelID,
	}
	if body.MakePubliclyAccessible != nil {
		v.MakePubliclyAccessible = *body.MakePubliclyAccessible
	}
	if body.MakePubliclyAccessible == nil {
		v.MakePubliclyAccessible = false
	}
	v.FileID = fileID
	v.Key = key

	return v
}

// NewRegisterTaskStatePayload builds a cascade service registerTaskState
// endpoint payload.
func NewRegisterTaskStatePayload(taskID string) *cascade.RegisterTaskStatePayload {
	v := &cascade.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v
}

// NewGetTaskHistoryPayload builds a cascade service getTaskHistory endpoint
// payload.
func NewGetTaskHistoryPayload(taskID string) *cascade.GetTaskHistoryPayload {
	v := &cascade.GetTaskHistoryPayload{}
	v.TaskID = taskID

	return v
}

// NewDownloadPayload builds a cascade service download endpoint payload.
func NewDownloadPayload(txid string, pid string, key string) *cascade.DownloadPayload {
	v := &cascade.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v
}

// ValidateUploadAssetRequestBody runs the validations defined on
// UploadAssetRequestBody
func ValidateUploadAssetRequestBody(body *UploadAssetRequestBody) (err error) {
	if body.Bytes == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
	}
	return
}

// ValidateStartProcessingRequestBody runs the validations defined on
// StartProcessingRequestBody
func ValidateStartProcessingRequestBody(body *StartProcessingRequestBody) (err error) {
	if body.BurnTxid == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("burn_txid", "body"))
	}
	if body.AppPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("app_pastelid", "body"))
	}
	if body.BurnTxid != nil {
		if utf8.RuneCountInString(*body.BurnTxid) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, true))
		}
	}
	if body.BurnTxid != nil {
		if utf8.RuneCountInString(*body.BurnTxid) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.burn_txid", *body.BurnTxid, utf8.RuneCountInString(*body.BurnTxid), 64, false))
		}
	}
	if body.AppPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", *body.AppPastelID, "^[a-zA-Z0-9]+$"))
	}
	if body.AppPastelID != nil {
		if utf8.RuneCountInString(*body.AppPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.AppPastelID, utf8.RuneCountInString(*body.AppPastelID), 86, true))
		}
	}
	if body.AppPastelID != nil {
		if utf8.RuneCountInString(*body.AppPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.AppPastelID, utf8.RuneCountInString(*body.AppPastelID), 86, false))
		}
	}
	return
}
