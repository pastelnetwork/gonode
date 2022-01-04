// Code generated by goa v3.5.3, DO NOT EDIT.
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

// UploadImageRequestBody is the type of the "cascade" service "uploadImage"
// endpoint HTTP request body.
type UploadImageRequestBody struct {
	// File to upload
	Bytes []byte `form:"file,omitempty" json:"file,omitempty" xml:"file,omitempty"`
	// For internal use
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// ActionDetailsRequestBody is the type of the "cascade" service
// "actionDetails" endpoint HTTP request body.
type ActionDetailsRequestBody struct {
	// 3rd party app's PastelID
	PastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
	// Hash (SHA3-256) of the Action Data
	ActionDataHash *string `form:"action_data_hash,omitempty" json:"action_data_hash,omitempty" xml:"action_data_hash,omitempty"`
	// The signature (base64) of the Action Data
	ActionDataSignature *string `form:"action_data_signature,omitempty" json:"action_data_signature,omitempty" xml:"action_data_signature,omitempty"`
}

// StartProcessingRequestBody is the type of the "cascade" service
// "startProcessing" endpoint HTTP request body.
type StartProcessingRequestBody struct {
	// Burn transaction ID
	BurnTxid *string `form:"burn_txid,omitempty" json:"burn_txid,omitempty" xml:"burn_txid,omitempty"`
	// App PastelID
	AppPastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
}

// UploadImageResponseBody is the type of the "cascade" service "uploadImage"
// endpoint HTTP response body.
type UploadImageResponseBody struct {
	// Uploaded image ID
	ImageID string `form:"image_id" json:"image_id" xml:"image_id"`
	// Uploaded image ID
	ExpiresIn string `form:"expires_in" json:"expires_in" xml:"expires_in"`
}

// ActionDetailsResponseBody is the type of the "cascade" service
// "actionDetails" endpoint HTTP response body.
type ActionDetailsResponseBody struct {
	// Estimated fee
	EstimatedFee float64 `form:"estimated_fee" json:"estimated_fee" xml:"estimated_fee"`
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

// UploadImageBadRequestResponseBody is the type of the "cascade" service
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

// UploadImageInternalServerErrorResponseBody is the type of the "cascade"
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

// ActionDetailsBadRequestResponseBody is the type of the "cascade" service
// "actionDetails" endpoint HTTP response body for the "BadRequest" error.
type ActionDetailsBadRequestResponseBody struct {
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

// ActionDetailsInternalServerErrorResponseBody is the type of the "cascade"
// service "actionDetails" endpoint HTTP response body for the
// "InternalServerError" error.
type ActionDetailsInternalServerErrorResponseBody struct {
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

// NewUploadImageResponseBody builds the HTTP response body from the result of
// the "uploadImage" endpoint of the "cascade" service.
func NewUploadImageResponseBody(res *cascadeviews.ImageView) *UploadImageResponseBody {
	body := &UploadImageResponseBody{
		ImageID:   *res.ImageID,
		ExpiresIn: *res.ExpiresIn,
	}
	return body
}

// NewActionDetailsResponseBody builds the HTTP response body from the result
// of the "actionDetails" endpoint of the "cascade" service.
func NewActionDetailsResponseBody(res *cascadeviews.ActionDetailResultView) *ActionDetailsResponseBody {
	body := &ActionDetailsResponseBody{
		EstimatedFee: *res.EstimatedFee,
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

// NewUploadImageBadRequestResponseBody builds the HTTP response body from the
// result of the "uploadImage" endpoint of the "cascade" service.
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
// from the result of the "uploadImage" endpoint of the "cascade" service.
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

// NewActionDetailsBadRequestResponseBody builds the HTTP response body from
// the result of the "actionDetails" endpoint of the "cascade" service.
func NewActionDetailsBadRequestResponseBody(res *goa.ServiceError) *ActionDetailsBadRequestResponseBody {
	body := &ActionDetailsBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewActionDetailsInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "actionDetails" endpoint of the "cascade"
// service.
func NewActionDetailsInternalServerErrorResponseBody(res *goa.ServiceError) *ActionDetailsInternalServerErrorResponseBody {
	body := &ActionDetailsInternalServerErrorResponseBody{
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

// NewUploadImagePayload builds a cascade service uploadImage endpoint payload.
func NewUploadImagePayload(body *UploadImageRequestBody) *cascade.UploadImagePayload {
	v := &cascade.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v
}

// NewActionDetailsPayload builds a cascade service actionDetails endpoint
// payload.
func NewActionDetailsPayload(body *ActionDetailsRequestBody, imageID string) *cascade.ActionDetailsPayload {
	v := &cascade.ActionDetailsPayload{
		PastelID:            *body.PastelID,
		ActionDataHash:      *body.ActionDataHash,
		ActionDataSignature: *body.ActionDataSignature,
	}
	v.ImageID = imageID

	return v
}

// NewStartProcessingPayload builds a cascade service startProcessing endpoint
// payload.
func NewStartProcessingPayload(body *StartProcessingRequestBody, imageID string, appPastelidPassphrase string) *cascade.StartProcessingPayload {
	v := &cascade.StartProcessingPayload{
		BurnTxid:    *body.BurnTxid,
		AppPastelID: *body.AppPastelID,
	}
	v.ImageID = imageID
	v.AppPastelidPassphrase = appPastelidPassphrase

	return v
}

// NewRegisterTaskStatePayload builds a cascade service registerTaskState
// endpoint payload.
func NewRegisterTaskStatePayload(taskID string) *cascade.RegisterTaskStatePayload {
	v := &cascade.RegisterTaskStatePayload{}
	v.TaskID = taskID

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

// ValidateActionDetailsRequestBody runs the validations defined on
// ActionDetailsRequestBody
func ValidateActionDetailsRequestBody(body *ActionDetailsRequestBody) (err error) {
	if body.PastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("app_pastelid", "body"))
	}
	if body.ActionDataHash == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("action_data_hash", "body"))
	}
	if body.ActionDataSignature == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("action_data_signature", "body"))
	}
	if body.PastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", *body.PastelID, "^[a-zA-Z0-9]"))
	}
	if body.PastelID != nil {
		if utf8.RuneCountInString(*body.PastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.PastelID, utf8.RuneCountInString(*body.PastelID), 86, true))
		}
	}
	if body.PastelID != nil {
		if utf8.RuneCountInString(*body.PastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.app_pastelid", *body.PastelID, utf8.RuneCountInString(*body.PastelID), 86, false))
		}
	}
	if body.ActionDataHash != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.action_data_hash", *body.ActionDataHash, "^[a-fA-F0-9]"))
	}
	if body.ActionDataHash != nil {
		if utf8.RuneCountInString(*body.ActionDataHash) < 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_hash", *body.ActionDataHash, utf8.RuneCountInString(*body.ActionDataHash), 64, true))
		}
	}
	if body.ActionDataHash != nil {
		if utf8.RuneCountInString(*body.ActionDataHash) > 64 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_hash", *body.ActionDataHash, utf8.RuneCountInString(*body.ActionDataHash), 64, false))
		}
	}
	if body.ActionDataSignature != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.action_data_signature", *body.ActionDataSignature, "^[a-zA-Z0-9\\/+]"))
	}
	if body.ActionDataSignature != nil {
		if utf8.RuneCountInString(*body.ActionDataSignature) < 152 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_signature", *body.ActionDataSignature, utf8.RuneCountInString(*body.ActionDataSignature), 152, true))
		}
	}
	if body.ActionDataSignature != nil {
		if utf8.RuneCountInString(*body.ActionDataSignature) > 152 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.action_data_signature", *body.ActionDataSignature, utf8.RuneCountInString(*body.ActionDataSignature), 152, false))
		}
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
