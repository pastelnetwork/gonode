// Code generated by goa v3.4.3, DO NOT EDIT.
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

// ActionDetailsRequestBody is the type of the "sense" service "actionDetails"
// endpoint HTTP request body.
type ActionDetailsRequestBody struct {
	// 3rd party app's PastelID
	PastelID string `form:"app_pastelid" json:"app_pastelid" xml:"app_pastelid"`
	// Hash (SHA3-256) of the Action Data
	ActionDataHash string `form:"action_data_hash" json:"action_data_hash" xml:"action_data_hash"`
	// The signature (base64) of the Action Data
	ActionDataSignature string `form:"action_data_signature" json:"action_data_signature" xml:"action_data_signature"`
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
}

// ActionDetailsResponseBody is the type of the "sense" service "actionDetails"
// endpoint HTTP response body.
type ActionDetailsResponseBody struct {
	// Estimated fee
	EstimatedFee *float64 `form:"estimated_fee,omitempty" json:"estimated_fee,omitempty" xml:"estimated_fee,omitempty"`
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

// ActionDetailsBadRequestResponseBody is the type of the "sense" service
// "actionDetails" endpoint HTTP response body for the "BadRequest" error.
type ActionDetailsBadRequestResponseBody struct {
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

// ActionDetailsInternalServerErrorResponseBody is the type of the "sense"
// service "actionDetails" endpoint HTTP response body for the
// "InternalServerError" error.
type ActionDetailsInternalServerErrorResponseBody struct {
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

// NewUploadImageRequestBody builds the HTTP request body from the payload of
// the "uploadImage" endpoint of the "sense" service.
func NewUploadImageRequestBody(p *sense.UploadImagePayload) *UploadImageRequestBody {
	body := &UploadImageRequestBody{
		Bytes:    p.Bytes,
		Filename: p.Filename,
	}
	return body
}

// NewActionDetailsRequestBody builds the HTTP request body from the payload of
// the "actionDetails" endpoint of the "sense" service.
func NewActionDetailsRequestBody(p *sense.ActionDetailsPayload) *ActionDetailsRequestBody {
	body := &ActionDetailsRequestBody{
		PastelID:            p.PastelID,
		ActionDataHash:      p.ActionDataHash,
		ActionDataSignature: p.ActionDataSignature,
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
		ImageID:   body.ImageID,
		ExpiresIn: body.ExpiresIn,
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

// NewActionDetailsActionDetailResultCreated builds a "sense" service
// "actionDetails" endpoint result from a HTTP "Created" response.
func NewActionDetailsActionDetailResultCreated(body *ActionDetailsResponseBody) *senseviews.ActionDetailResultView {
	v := &senseviews.ActionDetailResultView{
		EstimatedFee: body.EstimatedFee,
	}

	return v
}

// NewActionDetailsBadRequest builds a sense service actionDetails endpoint
// BadRequest error.
func NewActionDetailsBadRequest(body *ActionDetailsBadRequestResponseBody) *goa.ServiceError {
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

// NewActionDetailsInternalServerError builds a sense service actionDetails
// endpoint InternalServerError error.
func NewActionDetailsInternalServerError(body *ActionDetailsInternalServerErrorResponseBody) *goa.ServiceError {
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
		if !(*body.Status == "Task Started" || *body.Status == "Connected" || *body.Status == "Image Probed" || *body.Status == "Image And Thumbnail Uploaded" || *body.Status == "Status Gen ReptorQ Symbols" || *body.Status == "Preburn Registration Fee" || *body.Status == "Ticket Accepted" || *body.Status == "Ticket Registered" || *body.Status == "Ticket Activated" || *body.Status == "Error Insufficient Fee" || *body.Status == "Error Signatures Dont Match" || *body.Status == "Error Fingerprints Dont Match" || *body.Status == "Error ThumbnailHashes Dont Match" || *body.Status == "Error GenRaptorQ Symbols Failed" || *body.Status == "Task Rejected" || *body.Status == "Task Completed") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.status", *body.Status, []interface{}{"Task Started", "Connected", "Image Probed", "Image And Thumbnail Uploaded", "Status Gen ReptorQ Symbols", "Preburn Registration Fee", "Ticket Accepted", "Ticket Registered", "Ticket Activated", "Error Insufficient Fee", "Error Signatures Dont Match", "Error Fingerprints Dont Match", "Error ThumbnailHashes Dont Match", "Error GenRaptorQ Symbols Failed", "Task Rejected", "Task Completed"}))
		}
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

// ValidateActionDetailsBadRequestResponseBody runs the validations defined on
// actionDetails_BadRequest_response_body
func ValidateActionDetailsBadRequestResponseBody(body *ActionDetailsBadRequestResponseBody) (err error) {
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

// ValidateActionDetailsInternalServerErrorResponseBody runs the validations
// defined on actionDetails_InternalServerError_response_body
func ValidateActionDetailsInternalServerErrorResponseBody(body *ActionDetailsInternalServerErrorResponseBody) (err error) {
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
