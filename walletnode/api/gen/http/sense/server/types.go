// Code generated by goa v3.7.6, DO NOT EDIT.
//
// sense HTTP server types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package server

import (
	"unicode/utf8"

	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	senseviews "github.com/pastelnetwork/gonode/walletnode/api/gen/sense/views"
	goa "goa.design/goa/v3/pkg"
)

// UploadImageRequestBody is the type of the "sense" service "uploadImage"
// endpoint HTTP request body.
type UploadImageRequestBody struct {
	// File to upload
	Bytes []byte `form:"file,omitempty" json:"file,omitempty" xml:"file,omitempty"`
	// For internal use
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// StartProcessingRequestBody is the type of the "sense" service
// "startProcessing" endpoint HTTP request body.
type StartProcessingRequestBody struct {
	// Burn transaction ID
	BurnTxid *string `form:"burn_txid,omitempty" json:"burn_txid,omitempty" xml:"burn_txid,omitempty"`
	// App PastelID
	AppPastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
}

// UploadImageResponseBody is the type of the "sense" service "uploadImage"
// endpoint HTTP response body.
type UploadImageResponseBody struct {
	// Uploaded image ID
	ImageID string `form:"image_id" json:"image_id" xml:"image_id"`
	// Image expiration
	ExpiresIn string `form:"expires_in" json:"expires_in" xml:"expires_in"`
	// Estimated fee
	EstimatedFee float64 `form:"estimated_fee" json:"estimated_fee" xml:"estimated_fee"`
}

// StartProcessingResponseBody is the type of the "sense" service
// "startProcessing" endpoint HTTP response body.
type StartProcessingResponseBody struct {
	// Task ID of processing task
	TaskID string `form:"task_id" json:"task_id" xml:"task_id"`
}

// RegisterTaskStateResponseBody is the type of the "sense" service
// "registerTaskState" endpoint HTTP response body.
type RegisterTaskStateResponseBody struct {
	// Date of the status creation
	Date string `form:"date" json:"date" xml:"date"`
	// Status of the registration process
	Status string `form:"status" json:"status" xml:"status"`
}

// GetTaskHistoryResponseBody is the type of the "sense" service
// "getTaskHistory" endpoint HTTP response body.
type GetTaskHistoryResponseBody []*TaskHistoryResponse

// DownloadResponseBody is the type of the "sense" service "download" endpoint
// HTTP response body.
type DownloadResponseBody struct {
	// File downloaded
	File []byte `form:"file" json:"file" xml:"file"`
}

// UploadImageBadRequestResponseBody is the type of the "sense" service
// "uploadImage" endpoint HTTP response body for the "BadRequest" error.
type UploadImageBadRequestResponseBody struct {
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

// UploadImageInternalServerErrorResponseBody is the type of the "sense"
// service "uploadImage" endpoint HTTP response body for the
// "InternalServerError" error.
type UploadImageInternalServerErrorResponseBody struct {
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

// StartProcessingBadRequestResponseBody is the type of the "sense" service
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

// StartProcessingInternalServerErrorResponseBody is the type of the "sense"
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

// RegisterTaskStateNotFoundResponseBody is the type of the "sense" service
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

// RegisterTaskStateInternalServerErrorResponseBody is the type of the "sense"
// service "registerTaskState" endpoint HTTP response body for the
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

// GetTaskHistoryNotFoundResponseBody is the type of the "sense" service
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

// GetTaskHistoryInternalServerErrorResponseBody is the type of the "sense"
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

// DownloadNotFoundResponseBody is the type of the "sense" service "download"
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

// DownloadInternalServerErrorResponseBody is the type of the "sense" service
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
}

// NewUploadImageResponseBody builds the HTTP response body from the result of
// the "uploadImage" endpoint of the "sense" service.
func NewUploadImageResponseBody(res *senseviews.ImageView) *UploadImageResponseBody {
	body := &UploadImageResponseBody{
		ImageID:      *res.ImageID,
		ExpiresIn:    *res.ExpiresIn,
		EstimatedFee: *res.EstimatedFee,
	}
	return body
}

// NewStartProcessingResponseBody builds the HTTP response body from the result
// of the "startProcessing" endpoint of the "sense" service.
func NewStartProcessingResponseBody(res *senseviews.StartProcessingResultView) *StartProcessingResponseBody {
	body := &StartProcessingResponseBody{
		TaskID: *res.TaskID,
	}
	return body
}

// NewRegisterTaskStateResponseBody builds the HTTP response body from the
// result of the "registerTaskState" endpoint of the "sense" service.
func NewRegisterTaskStateResponseBody(res *sense.TaskState) *RegisterTaskStateResponseBody {
	body := &RegisterTaskStateResponseBody{
		Date:   res.Date,
		Status: res.Status,
	}
	return body
}

// NewGetTaskHistoryResponseBody builds the HTTP response body from the result
// of the "getTaskHistory" endpoint of the "sense" service.
func NewGetTaskHistoryResponseBody(res []*sense.TaskHistory) GetTaskHistoryResponseBody {
	body := make([]*TaskHistoryResponse, len(res))
	for i, val := range res {
		body[i] = marshalSenseTaskHistoryToTaskHistoryResponse(val)
	}
	return body
}

// NewDownloadResponseBody builds the HTTP response body from the result of the
// "download" endpoint of the "sense" service.
func NewDownloadResponseBody(res *sense.DownloadResult) *DownloadResponseBody {
	body := &DownloadResponseBody{
		File: res.File,
	}
	return body
}

// NewUploadImageBadRequestResponseBody builds the HTTP response body from the
// result of the "uploadImage" endpoint of the "sense" service.
func NewUploadImageBadRequestResponseBody(res *goa.ServiceError) *UploadImageBadRequestResponseBody {
	body := &UploadImageBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUploadImageInternalServerErrorResponseBody builds the HTTP response body
// from the result of the "uploadImage" endpoint of the "sense" service.
func NewUploadImageInternalServerErrorResponseBody(res *goa.ServiceError) *UploadImageInternalServerErrorResponseBody {
	body := &UploadImageInternalServerErrorResponseBody{
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
// the result of the "startProcessing" endpoint of the "sense" service.
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
// body from the result of the "startProcessing" endpoint of the "sense"
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
// the result of the "registerTaskState" endpoint of the "sense" service.
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
// body from the result of the "registerTaskState" endpoint of the "sense"
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
// result of the "getTaskHistory" endpoint of the "sense" service.
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
// body from the result of the "getTaskHistory" endpoint of the "sense" service.
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
// result of the "download" endpoint of the "sense" service.
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
// from the result of the "download" endpoint of the "sense" service.
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

// NewUploadImagePayload builds a sense service uploadImage endpoint payload.
func NewUploadImagePayload(body *UploadImageRequestBody) *sense.UploadImagePayload {
	v := &sense.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v
}

// NewStartProcessingPayload builds a sense service startProcessing endpoint
// payload.
func NewStartProcessingPayload(body *StartProcessingRequestBody, imageID string, appPastelidPassphrase string) *sense.StartProcessingPayload {
	v := &sense.StartProcessingPayload{
		BurnTxid:    *body.BurnTxid,
		AppPastelID: *body.AppPastelID,
	}
	v.ImageID = imageID
	v.AppPastelidPassphrase = appPastelidPassphrase

	return v
}

// NewRegisterTaskStatePayload builds a sense service registerTaskState
// endpoint payload.
func NewRegisterTaskStatePayload(taskID string) *sense.RegisterTaskStatePayload {
	v := &sense.RegisterTaskStatePayload{}
	v.TaskID = taskID

	return v
}

// NewGetTaskHistoryPayload builds a sense service getTaskHistory endpoint
// payload.
func NewGetTaskHistoryPayload(taskID string) *sense.GetTaskHistoryPayload {
	v := &sense.GetTaskHistoryPayload{}
	v.TaskID = taskID

	return v
}

// NewDownloadPayload builds a sense service download endpoint payload.
func NewDownloadPayload(txid string, pid string, key string) *sense.DownloadPayload {
	v := &sense.DownloadPayload{}
	v.Txid = txid
	v.Pid = pid
	v.Key = key

	return v
}

// ValidateUploadImageRequestBody runs the validations defined on
// UploadImageRequestBody
func ValidateUploadImageRequestBody(body *UploadImageRequestBody) (err error) {
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
