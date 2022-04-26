// Code generated by goa v3.6.2, DO NOT EDIT.
//
// cascade HTTP client types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package client

import (
	cascade "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade"
	cascadeviews "github.com/pastelnetwork/gonode/walletnode/api/gen/cascade/views"
	goa "goa.design/goa/v3/pkg"
)

// UploadAssetRequestBody is the type of the "cascade" service "uploadAsset"
// endpoint HTTP request body.
type UploadAssetRequestBody struct {
	// File to upload
	Bytes []byte `form:"file" json:"file" xml:"file"`
	// For internal use
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// StartProcessingRequestBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP request body.
type StartProcessingRequestBody struct {
	// Burn transaction ID
	BurnTxid string `form:"burn_txid" json:"burn_txid" xml:"burn_txid"`
	// App PastelID
	AppPastelID string `form:"app_pastelid" json:"app_pastelid" xml:"app_pastelid"`
}

// UploadAssetResponseBody is the type of the "cascade" service "uploadAsset"
// endpoint HTTP response body.
type UploadAssetResponseBody struct {
	// Uploaded file ID
	FileID *string `form:"file_id,omitempty" json:"file_id,omitempty" xml:"file_id,omitempty"`
	// File expiration
	ExpiresIn *string `form:"expires_in,omitempty" json:"expires_in,omitempty" xml:"expires_in,omitempty"`
	// Estimated fee
	EstimatedFee *float64 `form:"estimated_fee,omitempty" json:"estimated_fee,omitempty" xml:"estimated_fee,omitempty"`
}

// StartProcessingResponseBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP response body.
type StartProcessingResponseBody struct {
	// Task ID of processing task
	TaskID *string `form:"task_id,omitempty" json:"task_id,omitempty" xml:"task_id,omitempty"`
}

// RegisterTaskStateResponseBody is the type of the "cascade" service
// "registerTaskState" endpoint HTTP response body.
type RegisterTaskStateResponseBody struct {
	// Date of the status creation
	Date *string `form:"date,omitempty" json:"date,omitempty" xml:"date,omitempty"`
	// Status of the registration process
	Status *string `form:"status,omitempty" json:"status,omitempty" xml:"status,omitempty"`
}

// GetTaskHistoryResponseBody is the type of the "cascade" service
// "getTaskHistory" endpoint HTTP response body.
type GetTaskHistoryResponseBody struct {
	// List of past status strings
	List *string `form:"list,omitempty" json:"list,omitempty" xml:"list,omitempty"`
}

// DownloadResponseBody is the type of the "cascade" service "download"
// endpoint HTTP response body.
type DownloadResponseBody struct {
	// File downloaded
	File []byte `form:"file,omitempty" json:"file,omitempty" xml:"file,omitempty"`
}

// UploadAssetBadRequestResponseBody is the type of the "cascade" service
// "uploadAsset" endpoint HTTP response body for the "BadRequest" error.
type UploadAssetBadRequestResponseBody struct {
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

// UploadAssetInternalServerErrorResponseBody is the type of the "cascade"
// service "uploadAsset" endpoint HTTP response body for the
// "InternalServerError" error.
type UploadAssetInternalServerErrorResponseBody struct {
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

// StartProcessingBadRequestResponseBody is the type of the "cascade" service
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

// StartProcessingInternalServerErrorResponseBody is the type of the "cascade"
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

// RegisterTaskStateNotFoundResponseBody is the type of the "cascade" service
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

// RegisterTaskStateInternalServerErrorResponseBody is the type of the
// "cascade" service "registerTaskState" endpoint HTTP response body for the
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

// GetTaskHistoryNotFoundResponseBody is the type of the "cascade" service
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

// GetTaskHistoryInternalServerErrorResponseBody is the type of the "cascade"
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

// DownloadNotFoundResponseBody is the type of the "cascade" service "download"
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

// DownloadInternalServerErrorResponseBody is the type of the "cascade" service
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

// NewUploadAssetRequestBody builds the HTTP request body from the payload of
// the "uploadAsset" endpoint of the "cascade" service.
func NewUploadAssetRequestBody(p *cascade.UploadAssetPayload) *UploadAssetRequestBody {
	body := &UploadAssetRequestBody{
		Bytes:    p.Bytes,
		Filename: p.Filename,
	}
	return body
}

// NewStartProcessingRequestBody builds the HTTP request body from the payload
// of the "startProcessing" endpoint of the "cascade" service.
func NewStartProcessingRequestBody(p *cascade.StartProcessingPayload) *StartProcessingRequestBody {
	body := &StartProcessingRequestBody{
		BurnTxid:    p.BurnTxid,
		AppPastelID: p.AppPastelID,
	}
	return body
}

// NewUploadAssetAssetCreated builds a "cascade" service "uploadAsset" endpoint
// result from a HTTP "Created" response.
func NewUploadAssetAssetCreated(body *UploadAssetResponseBody) *cascadeviews.AssetView {
	v := &cascadeviews.AssetView{
		FileID:       body.FileID,
		ExpiresIn:    body.ExpiresIn,
		EstimatedFee: body.EstimatedFee,
	}

	return v
}

// NewUploadAssetBadRequest builds a cascade service uploadAsset endpoint
// BadRequest error.
func NewUploadAssetBadRequest(body *UploadAssetBadRequestResponseBody) *goa.ServiceError {
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

// NewUploadAssetInternalServerError builds a cascade service uploadAsset
// endpoint InternalServerError error.
func NewUploadAssetInternalServerError(body *UploadAssetInternalServerErrorResponseBody) *goa.ServiceError {
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

// NewStartProcessingResultViewCreated builds a "cascade" service
// "startProcessing" endpoint result from a HTTP "Created" response.
func NewStartProcessingResultViewCreated(body *StartProcessingResponseBody) *cascadeviews.StartProcessingResultView {
	v := &cascadeviews.StartProcessingResultView{
		TaskID: body.TaskID,
	}

	return v
}

// NewStartProcessingBadRequest builds a cascade service startProcessing
// endpoint BadRequest error.
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

// NewStartProcessingInternalServerError builds a cascade service
// startProcessing endpoint InternalServerError error.
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

// NewRegisterTaskStateTaskStateOK builds a "cascade" service
// "registerTaskState" endpoint result from a HTTP "OK" response.
func NewRegisterTaskStateTaskStateOK(body *RegisterTaskStateResponseBody) *cascade.TaskState {
	v := &cascade.TaskState{
		Date:   *body.Date,
		Status: *body.Status,
	}

	return v
}

// NewRegisterTaskStateNotFound builds a cascade service registerTaskState
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

// NewRegisterTaskStateInternalServerError builds a cascade service
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

// NewGetTaskHistoryTaskHistoryOK builds a "cascade" service "getTaskHistory"
// endpoint result from a HTTP "OK" response.
func NewGetTaskHistoryTaskHistoryOK(body *GetTaskHistoryResponseBody) *cascade.TaskHistory {
	v := &cascade.TaskHistory{
		List: *body.List,
	}

	return v
}

// NewGetTaskHistoryNotFound builds a cascade service getTaskHistory endpoint
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

// NewGetTaskHistoryInternalServerError builds a cascade service getTaskHistory
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

// NewDownloadResultOK builds a "cascade" service "download" endpoint result
// from a HTTP "OK" response.
func NewDownloadResultOK(body *DownloadResponseBody) *cascade.DownloadResult {
	v := &cascade.DownloadResult{
		File: body.File,
	}

	return v
}

// NewDownloadNotFound builds a cascade service download endpoint NotFound
// error.
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

// NewDownloadInternalServerError builds a cascade service download endpoint
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
		if !(*body.Status == "Task Started" || *body.Status == "Connected" || *body.Status == "Image Probed" || *body.Status == "Image And Thumbnail Uploaded" || *body.Status == "Status Gen ReptorQ Symbols" || *body.Status == "Preburn Registration Fee" || *body.Status == "Downloaded" || *body.Status == "Request Accepted" || *body.Status == "Request Registered" || *body.Status == "Request Activated" || *body.Status == "Error Insufficient Fee" || *body.Status == "Error Signatures Dont Match" || *body.Status == "Error Fingerprints Dont Match" || *body.Status == "Error ThumbnailHashes Dont Match" || *body.Status == "Error GenRaptorQ Symbols Failed" || *body.Status == "Error File Don't Match" || *body.Status == "Error Not Enough SuperNode" || *body.Status == "Error Find Responding SNs" || *body.Status == "Error Not Enough Downloaded Filed" || *body.Status == "Error Download Failed" || *body.Status == "Error Invalid Burn TxID" || *body.Status == "Task Failed" || *body.Status == "Task Rejected" || *body.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.status", *body.Status, []interface{}{"Task Started", "Connected", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Downloaded", "Request Accepted", "Request Registered", "Request Activated", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Error File Don't Match", "Error Not Enough SuperNode", "Error Find Responding SNs", "Error Not Enough Downloaded Filed", "Error Download Failed", "Error Invalid Burn TxID", "Task Failed", "Task Rejected", "Task Completed"}))
		}
	}
	return
}

// ValidateGetTaskHistoryResponseBody runs the validations defined on
// GetTaskHistoryResponseBody
func ValidateGetTaskHistoryResponseBody(body *GetTaskHistoryResponseBody) (err error) {
	if body.List == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("list", "body"))
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

// ValidateUploadAssetBadRequestResponseBody runs the validations defined on
// uploadAsset_BadRequest_response_body
func ValidateUploadAssetBadRequestResponseBody(body *UploadAssetBadRequestResponseBody) (err error) {
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

// ValidateUploadAssetInternalServerErrorResponseBody runs the validations
// defined on uploadAsset_InternalServerError_response_body
func ValidateUploadAssetInternalServerErrorResponseBody(body *UploadAssetInternalServerErrorResponseBody) (err error) {
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
