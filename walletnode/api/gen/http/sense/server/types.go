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

// ActionDetailsRequestBody is the type of the "sense" service "actionDetails"
// endpoint HTTP request body.
type ActionDetailsRequestBody struct {
	// 3rd party app's PastelID
	PastelID *string `form:"app_pastelid,omitempty" json:"app_pastelid,omitempty" xml:"app_pastelid,omitempty"`
	// Hash (SHA3-256) of the Action Data
	ActionDataHash *string `form:"action_data_hash,omitempty" json:"action_data_hash,omitempty" xml:"action_data_hash,omitempty"`
	// The signature (base64) of the Action Data
	ActionDataSignature *string `form:"action_data_signature,omitempty" json:"action_data_signature,omitempty" xml:"action_data_signature,omitempty"`
}

// UploadImageResponseBody is the type of the "sense" service "uploadImage"
// endpoint HTTP response body.
type UploadImageResponseBody struct {
	// Uploaded image ID
	ImageID string `form:"image_id" json:"image_id" xml:"image_id"`
	// Uploaded image ID
	ExpiresIn string `form:"expires_in" json:"expires_in" xml:"expires_in"`
}

// ActionDetailsResponseBody is the type of the "sense" service "actionDetails"
// endpoint HTTP response body.
type ActionDetailsResponseBody struct {
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

// ActionDetailsBadRequestResponseBody is the type of the "sense" service
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

// ActionDetailsInternalServerErrorResponseBody is the type of the "sense"
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

// NewUploadImageResponseBody builds the HTTP response body from the result of
// the "uploadImage" endpoint of the "sense" service.
func NewUploadImageResponseBody(res *senseviews.ImageView) *UploadImageResponseBody {
	body := &UploadImageResponseBody{
		ImageID:   *res.ImageID,
		ExpiresIn: *res.ExpiresIn,
	}
	return body
}

// NewActionDetailsResponseBody builds the HTTP response body from the result
// of the "actionDetails" endpoint of the "sense" service.
func NewActionDetailsResponseBody(res *senseviews.ActionDetailResultView) *ActionDetailsResponseBody {
	body := &ActionDetailsResponseBody{
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

// NewActionDetailsBadRequestResponseBody builds the HTTP response body from
// the result of the "actionDetails" endpoint of the "sense" service.
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
// body from the result of the "actionDetails" endpoint of the "sense" service.
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

// NewUploadImagePayload builds a sense service uploadImage endpoint payload.
func NewUploadImagePayload(body *UploadImageRequestBody) *sense.UploadImagePayload {
	v := &sense.UploadImagePayload{
		Bytes:    body.Bytes,
		Filename: body.Filename,
	}

	return v
}

// NewActionDetailsPayload builds a sense service actionDetails endpoint
// payload.
func NewActionDetailsPayload(body *ActionDetailsRequestBody, imageID string) *sense.ActionDetailsPayload {
	v := &sense.ActionDetailsPayload{
		PastelID:            *body.PastelID,
		ActionDataHash:      *body.ActionDataHash,
		ActionDataSignature: *body.ActionDataSignature,
	}
	v.ImageID = imageID

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
		err = goa.MergeErrors(err, goa.ValidatePattern("body.action_data_signature", *body.ActionDataSignature, "^[a-zA-Z0-9]+$"))
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
