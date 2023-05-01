// Code generated by goa v3.7.6, DO NOT EDIT.
//
// sense HTTP client types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package client

import (
	sense "github.com/pastelnetwork/gonode/walletnode/api/gen/sense"
	senseviews "github.com/pastelnetwork/gonode/walletnode/api/gen/sense/views"
	goa "goa.design/goa/v3/pkg"
)

// UploadImageRequestBody is the type of the "sense" service "uploadImage"
// endpoint HTTP request body.
type UploadImageRequestBody struct {
	// File to upload
	Bytes []byte `form:"file" json:"file" xml:"file"`
	// For internal use
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// StartProcessingRequestBody is the type of the "sense" service
// "startProcessing" endpoint HTTP request body.
type StartProcessingRequestBody struct {
	// Burn transaction ID
	BurnTxid string `form:"burn_txid" json:"burn_txid" xml:"burn_txid"`
	// App PastelID
	AppPastelID string `form:"app_pastelid" json:"app_pastelid" xml:"app_pastelid"`
}

// UploadImageResponseBody is the type of the "sense" service "uploadImage"
// endpoint HTTP response body.
type UploadImageResponseBody struct {
	// Uploaded image ID
	ImageID *string `form:"image_id,omitempty" json:"image_id,omitempty" xml:"image_id,omitempty"`
	// Image expiration
	ExpiresIn *string `form:"expires_in,omitempty" json:"expires_in,omitempty" xml:"expires_in,omitempty"`
	// Estimated fee
	TotalEstimatedFee *float64 `form:"total_estimated_fee,omitempty" json:"total_estimated_fee,omitempty" xml:"total_estimated_fee,omitempty"`
	// The amount that's required to be preburned
	RequiredPreburnAmount *float64 `form:"required_preburn_amount,omitempty" json:"required_preburn_amount,omitempty" xml:"required_preburn_amount,omitempty"`
}

// StartProcessingResponseBody is the type of the "sense" service
// "startProcessing" endpoint HTTP response body.
type StartProcessingResponseBody struct {
	// Task ID of processing task
	TaskID *string `form:"task_id,omitempty" json:"task_id,omitempty" xml:"task_id,omitempty"`
}

// RegisterTaskStateResponseBody is the type of the "sense" service
// "registerTaskState" endpoint HTTP response body.
type RegisterTaskStateResponseBody struct {
	// Date of the status creation
	Date *string `form:"date,omitempty" json:"date,omitempty" xml:"date,omitempty"`
	// Status of the registration process
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
}

// GetTaskHistoryResponseBody is the type of the "sense" service
// "getTaskHistory" endpoint HTTP response body.
type GetTaskHistoryResponseBody []*TaskHistoryResponse

// DownloadResponseBody is the type of the "sense" service "download" endpoint
// HTTP response body.
type DownloadResponseBody struct {
	// File downloaded
	File []byte `form:"file,omitempty" json:"file,omitempty" xml:"file,omitempty"`
}

// UploadImageBadRequestResponseBody is the type of the "sense" service
// "uploadImage" endpoint HTTP response body for the "BadRequest" error.
type UploadImageBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// UploadImageInternalServerErrorResponseBody is the type of the "sense"
// service "uploadImage" endpoint HTTP response body for the
// "InternalServerError" error.
type UploadImageInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// StartProcessingBadRequestResponseBody is the type of the "sense" service
// "startProcessing" endpoint HTTP response body for the "BadRequest" error.
type StartProcessingBadRequestResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// StartProcessingInternalServerErrorResponseBody is the type of the "sense"
// service "startProcessing" endpoint HTTP response body for the
// "InternalServerError" error.
type StartProcessingInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// RegisterTaskStateNotFoundResponseBody is the type of the "sense" service
// "registerTaskState" endpoint HTTP response body for the "NotFound" error.
type RegisterTaskStateNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// RegisterTaskStateInternalServerErrorResponseBody is the type of the "sense"
// service "registerTaskState" endpoint HTTP response body for the
// "InternalServerError" error.
type RegisterTaskStateInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// GetTaskHistoryNotFoundResponseBody is the type of the "sense" service
// "getTaskHistory" endpoint HTTP response body for the "NotFound" error.
type GetTaskHistoryNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// GetTaskHistoryInternalServerErrorResponseBody is the type of the "sense"
// service "getTaskHistory" endpoint HTTP response body for the
// "InternalServerError" error.
type GetTaskHistoryInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// DownloadNotFoundResponseBody is the type of the "sense" service "download"
// endpoint HTTP response body for the "NotFound" error.
type DownloadNotFoundResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// DownloadInternalServerErrorResponseBody is the type of the "sense" service
// "download" endpoint HTTP response body for the "InternalServerError" error.
type DownloadInternalServerErrorResponseBody struct {
	// Name is the name of this class of errors.
	Name *string `form:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
	// ID is a unique identifier for this particular occurrence of the problem.
	ID *string `form:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
	// Message is a human-readable explanation specific to this occurrence of the
	// problem.
	Message *string `form:"message,omitempty" json:"message,omitempty" xml:"message,omitempty"`
	// Is the error temporary?
	Temporary *bool `form:"temporary,omitempty" json:"temporary,omitempty" xml:"temporary,omitempty"`
	// Is the error a timeout?
	Timeout *bool `form:"timeout,omitempty" json:"timeout,omitempty" xml:"timeout,omitempty"`
	// Is the error a server-side fault?
	Fault *bool `form:"fault,omitempty" json:"fault,omitempty" xml:"fault,omitempty"`
}

// TaskHistoryResponse is used to define fields on response body types.
type TaskHistoryResponse struct {
	// Timestamp of the status creation
	Timestamp *string `form:"timestamp,omitempty" json:"timestamp,omitempty" xml:"timestamp,omitempty"`
	// past status string
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
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

// NewUploadImageRequestBody builds the HTTP request body from the payload of
// the "uploadImage" endpoint of the "sense" service.
func NewUploadImageRequestBody(p *sense.UploadImagePayload) *UploadImageRequestBody {
	body := &UploadImageRequestBody{
		Bytes:    p.Bytes,
		Filename: p.Filename,
	}
	return body
}

// NewStartProcessingRequestBody builds the HTTP request body from the payload
// of the "startProcessing" endpoint of the "sense" service.
func NewStartProcessingRequestBody(p *sense.StartProcessingPayload) *StartProcessingRequestBody {
	body := &StartProcessingRequestBody{
		BurnTxid:    p.BurnTxid,
		AppPastelID: p.AppPastelID,
	}
	return body
}

// NewUploadImageImageCreated builds a "sense" service "uploadImage" endpoint
// result from a HTTP "Created" response.
func NewUploadImageImageCreated(body *UploadImageResponseBody) *senseviews.ImageView {
	v := &senseviews.ImageView{
		ImageID:               body.ImageID,
		ExpiresIn:             body.ExpiresIn,
		TotalEstimatedFee:     body.TotalEstimatedFee,
		RequiredPreburnAmount: body.RequiredPreburnAmount,
	}

	return v
}

// NewUploadImageBadRequest builds a sense service uploadImage endpoint
// BadRequest error.
func NewUploadImageBadRequest(body *UploadImageBadRequestResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewUploadImageInternalServerError builds a sense service uploadImage
// endpoint InternalServerError error.
func NewUploadImageInternalServerError(body *UploadImageInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewStartProcessingResultViewCreated builds a "sense" service
// "startProcessing" endpoint result from a HTTP "Created" response.
func NewStartProcessingResultViewCreated(body *StartProcessingResponseBody) *senseviews.StartProcessingResultView {
	v := &senseviews.StartProcessingResultView{
		TaskID: body.TaskID,
	}

	return v
}

// NewStartProcessingBadRequest builds a sense service startProcessing endpoint
// BadRequest error.
func NewStartProcessingBadRequest(body *StartProcessingBadRequestResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewStartProcessingInternalServerError builds a sense service startProcessing
// endpoint InternalServerError error.
func NewStartProcessingInternalServerError(body *StartProcessingInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewRegisterTaskStateTaskStateOK builds a "sense" service "registerTaskState"
// endpoint result from a HTTP "OK" response.
func NewRegisterTaskStateTaskStateOK(body *RegisterTaskStateResponseBody) *sense.TaskState {
	v := &sense.TaskState{
		Date:   *body.Date,
		Status: *body.Status,
	}

	return v
}

// NewRegisterTaskStateNotFound builds a sense service registerTaskState
// endpoint NotFound error.
func NewRegisterTaskStateNotFound(body *RegisterTaskStateNotFoundResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewRegisterTaskStateInternalServerError builds a sense service
// registerTaskState endpoint InternalServerError error.
func NewRegisterTaskStateInternalServerError(body *RegisterTaskStateInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewGetTaskHistoryTaskHistoryOK builds a "sense" service "getTaskHistory"
// endpoint result from a HTTP "OK" response.
func NewGetTaskHistoryTaskHistoryOK(body []*TaskHistoryResponse) []*sense.TaskHistory {
	v := make([]*sense.TaskHistory, len(body))
	for i, val := range body {
		v[i] = unmarshalTaskHistoryResponseToSenseTaskHistory(val)
	}

	return v
}

// NewGetTaskHistoryNotFound builds a sense service getTaskHistory endpoint
// NotFound error.
func NewGetTaskHistoryNotFound(body *GetTaskHistoryNotFoundResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewGetTaskHistoryInternalServerError builds a sense service getTaskHistory
// endpoint InternalServerError error.
func NewGetTaskHistoryInternalServerError(body *GetTaskHistoryInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewDownloadResultOK builds a "sense" service "download" endpoint result from
// a HTTP "OK" response.
func NewDownloadResultOK(body *DownloadResponseBody) *sense.DownloadResult {
	v := &sense.DownloadResult{
		File: body.File,
	}

	return v
}

// NewDownloadNotFound builds a sense service download endpoint NotFound error.
func NewDownloadNotFound(body *DownloadNotFoundResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// NewDownloadInternalServerError builds a sense service download endpoint
// InternalServerError error.
func NewDownloadInternalServerError(body *DownloadInternalServerErrorResponseBody) *goa.ServiceError {
	v := &goa.ServiceError{
		Name:      *body.Name,
		ID:        *body.ID,
		Message:   *body.Message,
		Temporary: *body.Temporary,
		Timeout:   *body.Timeout,
		Fault:     *body.Fault,
	}

	return v
}

// ValidateRegisterTaskStateResponseBody runs the validations defined on
// RegisterTaskStateResponseBody
func ValidateRegisterTaskStateResponseBody(body *RegisterTaskStateResponseBody) (err error) {
	if body.Date == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("date", "body"))
	}
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	if body.Status != nil {
		if !(*body.Status == "Task Started" || *body.Status == "Connected" || *body.Status == "Validating Burn Txn" || *body.Status == "Image Probed" || *body.Status == "Image And Thumbnail Uploaded" || *body.Status == "Status Gen ReptorQ Symbols" || *body.Status == "Preburn Registration Fee" || *body.Status == "Downloaded" || *body.Status == "Request Accepted" || *body.Status == "Request Registered" || *body.Status == "Request Activated" || *body.Status == "Error Setting up mesh of supernodes" || *body.Status == "Error Sending Reg Metadata" || *body.Status == "Error Uploading Image" || *body.Status == "Error Converting Image to Bytes" || *body.Status == "Error Encoding Image" || *body.Status == "Error Creating Ticket" || *body.Status == "Error Signing Ticket" || *body.Status == "Error Uploading Ticket" || *body.Status == "Error Activating Ticket" || *body.Status == "Error Probing Image" || *body.Status == "Error Generating DD and Fingerprint IDs" || *body.Status == "Error comparing suitable storage fee with task request maximum fee" || *body.Status == "Error balance not sufficient" || *body.Status == "Error getting hash of the image" || *body.Status == "Error sending signed ticket to SNs" || *body.Status == "Error checking balance" || *body.Status == "Error burning reg fee to get reg ticket id" || *body.Status == "Error validating reg ticket txn id" || *body.Status == "Error validating activate ticket txn id" || *body.Status == "Error Insufficient Fee" || *body.Status == "Error Signatures Dont Match" || *body.Status == "Error Fingerprints Dont Match" || *body.Status == "Error ThumbnailHashes Dont Match" || *body.Status == "Error GenRaptorQ Symbols Failed" || *body.Status == "Error File Don't Match" || *body.Status == "Error Not Enough SuperNode" || *body.Status == "Error Find Responding SNs" || *body.Status == "Error Not Enough Downloaded Filed" || *body.Status == "Error Download Failed" || *body.Status == "Error Invalid Burn TxID" || *body.Status == "Task Failed" || *body.Status == "Task Rejected" || *body.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.status", *body.Status, []interface{}{"Task Started", "Connected", "Validating Burn Txn", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Downloaded", "Request Accepted", "Request Registered", "Request Activated", "Error Setting up mesh of supernodes", "Error Sending Reg Metadata", "Error Uploading Image", "Error Converting Image to Bytes", "Error Encoding Image", "Error Creating Ticket", "Error Signing Ticket", "Error Uploading Ticket", "Error Activating Ticket", "Error Probing Image", "Error Generating DD and Fingerprint IDs", "Error comparing suitable storage fee with task request maximum fee", "Error balance not sufficient", "Error getting hash of the image", "Error sending signed ticket to SNs", "Error checking balance", "Error burning reg fee to get reg ticket id", "Error validating reg ticket txn id", "Error validating activate ticket txn id", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Error File Don't Match", "Error Not Enough SuperNode", "Error Find Responding SNs", "Error Not Enough Downloaded Filed", "Error Download Failed", "Error Invalid Burn TxID", "Task Failed", "Task Rejected", "Task Completed"}))
		}
	}
	return
}

// ValidateDownloadResponseBody runs the validations defined on
// DownloadResponseBody
func ValidateDownloadResponseBody(body *DownloadResponseBody) (err error) {
	if body.File == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("file", "body"))
	}
	return
}

// ValidateUploadImageBadRequestResponseBody runs the validations defined on
// uploadImage_BadRequest_response_body
func ValidateUploadImageBadRequestResponseBody(body *UploadImageBadRequestResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateUploadImageInternalServerErrorResponseBody runs the validations
// defined on uploadImage_InternalServerError_response_body
func ValidateUploadImageInternalServerErrorResponseBody(body *UploadImageInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateStartProcessingBadRequestResponseBody runs the validations defined
// on startProcessing_BadRequest_response_body
func ValidateStartProcessingBadRequestResponseBody(body *StartProcessingBadRequestResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateStartProcessingInternalServerErrorResponseBody runs the validations
// defined on startProcessing_InternalServerError_response_body
func ValidateStartProcessingInternalServerErrorResponseBody(body *StartProcessingInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateRegisterTaskStateNotFoundResponseBody runs the validations defined
// on registerTaskState_NotFound_response_body
func ValidateRegisterTaskStateNotFoundResponseBody(body *RegisterTaskStateNotFoundResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateRegisterTaskStateInternalServerErrorResponseBody runs the
// validations defined on registerTaskState_InternalServerError_response_body
func ValidateRegisterTaskStateInternalServerErrorResponseBody(body *RegisterTaskStateInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateGetTaskHistoryNotFoundResponseBody runs the validations defined on
// getTaskHistory_NotFound_response_body
func ValidateGetTaskHistoryNotFoundResponseBody(body *GetTaskHistoryNotFoundResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateGetTaskHistoryInternalServerErrorResponseBody runs the validations
// defined on getTaskHistory_InternalServerError_response_body
func ValidateGetTaskHistoryInternalServerErrorResponseBody(body *GetTaskHistoryInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateDownloadNotFoundResponseBody runs the validations defined on
// download_NotFound_response_body
func ValidateDownloadNotFoundResponseBody(body *DownloadNotFoundResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateDownloadInternalServerErrorResponseBody runs the validations defined
// on download_InternalServerError_response_body
func ValidateDownloadInternalServerErrorResponseBody(body *DownloadInternalServerErrorResponseBody) (err error) {
	if body.Name == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("name", "body"))
	}
	if body.ID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("id", "body"))
	}
	if body.Message == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("message", "body"))
	}
	if body.Temporary == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("temporary", "body"))
	}
	if body.Timeout == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("timeout", "body"))
	}
	if body.Fault == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("fault", "body"))
	}
	return
}

// ValidateTaskHistoryResponse runs the validations defined on
// TaskHistoryResponse
func ValidateTaskHistoryResponse(body *TaskHistoryResponse) (err error) {
	if body.Status == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("status", "body"))
	}
	return
}
