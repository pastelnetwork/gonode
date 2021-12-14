// Code generated by goa v3.5.3, DO NOT EDIT.
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

// StartTaskRequestBody is the type of the "sense" service "startTask" endpoint
// HTTP request body.
type StartTaskRequestBody struct {
	// 3rd party app's PastelID
	PastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
	// Hash (SHA3-256) of the Action Data
	ActionDataHash []byte `form:"action_data_hash,omitempty" json:"action_data_hash,omitempty" xml:"action_data_hash,omitempty"`
	// The signature of the Action Data
	ActionDataSignature []byte `form:"action_data_signature,omitempty" json:"action_data_signature,omitempty" xml:"action_data_signature,omitempty"`
}

// UploadImageResponseBody is the type of the "sense" service "uploadImage"
// endpoint HTTP response body.
type UploadImageResponseBody struct {
	// Task ID
	TaskID string `form:"task_id" json:"task_id" xml:"task_id"`
}

// StartTaskResponseBody is the type of the "sense" service "startTask"
// endpoint HTTP response body.
type StartTaskResponseBody struct {
	// Estimated fee
	EstimatedFee float64 `form:"estimated_fee" json:"estimated_fee" xml:"estimated_fee"`
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

// StartTaskBadRequestResponseBody is the type of the "sense" service
// "startTask" endpoint HTTP response body for the "BadRequest" error.
type StartTaskBadRequestResponseBody struct {
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

// StartTaskInternalServerErrorResponseBody is the type of the "sense" service
// "startTask" endpoint HTTP response body for the "InternalServerError" error.
type StartTaskInternalServerErrorResponseBody struct {
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
// the "uploadImage" endpoint of the "sense" service.
func NewUploadImageResponseBody(res *senseviews.ImageUploadResultView) *UploadImageResponseBody {
	body := &UploadImageResponseBody{
		TaskID: *res.TaskID,
	}
	return body
}

// NewStartTaskResponseBody builds the HTTP response body from the result of
// the "startTask" endpoint of the "sense" service.
func NewStartTaskResponseBody(res *senseviews.StartActionDataResultView) *StartTaskResponseBody {
	body := &StartTaskResponseBody{
		EstimatedFee: *res.EstimatedFee,
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

// NewStartTaskBadRequestResponseBody builds the HTTP response body from the
// result of the "startTask" endpoint of the "sense" service.
func NewStartTaskBadRequestResponseBody(res *goa.ServiceError) *StartTaskBadRequestResponseBody {
	body := &StartTaskBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewStartTaskInternalServerErrorResponseBody builds the HTTP response body
// from the result of the "startTask" endpoint of the "sense" service.
func NewStartTaskInternalServerErrorResponseBody(res *goa.ServiceError) *StartTaskInternalServerErrorResponseBody {
	body := &StartTaskInternalServerErrorResponseBody{
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

// NewStartTaskPayload builds a sense service startTask endpoint payload.
func NewStartTaskPayload(body *StartTaskRequestBody, taskID string) *sense.StartTaskPayload {
	v := &sense.StartTaskPayload{
		PastelID:            *body.PastelID,
		ActionDataHash:      body.ActionDataHash,
		ActionDataSignature: body.ActionDataSignature,
	}
	v.TaskID = &taskID

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

// ValidateStartTaskRequestBody runs the validations defined on
// StartTaskRequestBody
func ValidateStartTaskRequestBody(body *StartTaskRequestBody) (err error) {
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
		err = goa.MergeErrors(err, goa.ValidatePattern("body.app_pastelid", *body.PastelID, "^[a-zA-Z0-9]+$"))
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
	return
}
