// Code generated by goa v3.4.3, DO NOT EDIT.
//
// userdatas HTTP server types
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design -o api/

package server

import (
	"unicode/utf8"

	userdatas "github.com/pastelnetwork/gonode/walletnode/api/gen/userdatas"
	goa "goa.design/goa/v3/pkg"
)

// CreateUserdataRequestBody is the type of the "userdatas" service
// "createUserdata" endpoint HTTP request body.
type CreateUserdataRequestBody struct {
	// Real name of the user
	RealName *string `form:"realname,omitempty" json:"realname,omitempty" xml:"realname,omitempty"`
	// Facebook link of the user
	FacebookLink *string `form:"facebook_link,omitempty" json:"facebook_link,omitempty" xml:"facebook_link,omitempty"`
	// Twitter link of the user
	TwitterLink *string `form:"twitter_link,omitempty" json:"twitter_link,omitempty" xml:"twitter_link,omitempty"`
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string `form:"native_currency,omitempty" json:"native_currency,omitempty" xml:"native_currency,omitempty"`
	// Location of the user
	Location *string `form:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	// Primary language of the user, follow ISO 639-2 standard
	PrimaryLanguage *string `form:"primary_language,omitempty" json:"primary_language,omitempty" xml:"primary_language,omitempty"`
	// The categories of user's work, separate by ,
	Categories *string `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Biography of the user
	Biography *string `form:"biography,omitempty" json:"biography,omitempty" xml:"biography,omitempty"`
	// Avatar image of the user
	AvatarImage *UserImageUploadPayloadRequestBody `form:"avatar_image,omitempty" json:"avatar_image,omitempty" xml:"avatar_image,omitempty"`
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayloadRequestBody `form:"cover_photo,omitempty" json:"cover_photo,omitempty" xml:"cover_photo,omitempty"`
	// User's PastelID
	UserPastelID *string `form:"user_pastelid,omitempty" json:"user_pastelid,omitempty" xml:"user_pastelid,omitempty"`
	// Passphrase of the user's PastelID
	UserPastelIDPassphrase *string `form:"user_pastelid_passphrase,omitempty" json:"user_pastelid_passphrase,omitempty" xml:"user_pastelid_passphrase,omitempty"`
}

// UpdateUserdataRequestBody is the type of the "userdatas" service
// "updateUserdata" endpoint HTTP request body.
type UpdateUserdataRequestBody struct {
	// Real name of the user
	RealName *string `form:"realname,omitempty" json:"realname,omitempty" xml:"realname,omitempty"`
	// Facebook link of the user
	FacebookLink *string `form:"facebook_link,omitempty" json:"facebook_link,omitempty" xml:"facebook_link,omitempty"`
	// Twitter link of the user
	TwitterLink *string `form:"twitter_link,omitempty" json:"twitter_link,omitempty" xml:"twitter_link,omitempty"`
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string `form:"native_currency,omitempty" json:"native_currency,omitempty" xml:"native_currency,omitempty"`
	// Location of the user
	Location *string `form:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	// Primary language of the user, follow ISO 639-2 standard
	PrimaryLanguage *string `form:"primary_language,omitempty" json:"primary_language,omitempty" xml:"primary_language,omitempty"`
	// The categories of user's work, separate by ,
	Categories *string `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Biography of the user
	Biography *string `form:"biography,omitempty" json:"biography,omitempty" xml:"biography,omitempty"`
	// Avatar image of the user
	AvatarImage *UserImageUploadPayloadRequestBody `form:"avatar_image,omitempty" json:"avatar_image,omitempty" xml:"avatar_image,omitempty"`
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayloadRequestBody `form:"cover_photo,omitempty" json:"cover_photo,omitempty" xml:"cover_photo,omitempty"`
	// User's PastelID
	UserPastelID *string `form:"user_pastelid,omitempty" json:"user_pastelid,omitempty" xml:"user_pastelid,omitempty"`
	// Passphrase of the user's PastelID
	UserPastelIDPassphrase *string `form:"user_pastelid_passphrase,omitempty" json:"user_pastelid_passphrase,omitempty" xml:"user_pastelid_passphrase,omitempty"`
}

// CreateUserdataResponseBody is the type of the "userdatas" service
// "createUserdata" endpoint HTTP response body.
type CreateUserdataResponseBody struct {
	// Result of the request is success or not
	ResponseCode int `form:"response_code" json:"response_code" xml:"response_code"`
	// The detail of why result is success/fail, depend on response_code
	Detail string `form:"detail" json:"detail" xml:"detail"`
	// Error detail on realname
	RealName *string `form:"realname,omitempty" json:"realname,omitempty" xml:"realname,omitempty"`
	// Error detail on facebook_link
	FacebookLink *string `form:"facebook_link,omitempty" json:"facebook_link,omitempty" xml:"facebook_link,omitempty"`
	// Error detail on twitter_link
	TwitterLink *string `form:"twitter_link,omitempty" json:"twitter_link,omitempty" xml:"twitter_link,omitempty"`
	// Error detail on native_currency
	NativeCurrency *string `form:"native_currency,omitempty" json:"native_currency,omitempty" xml:"native_currency,omitempty"`
	// Error detail on location
	Location *string `form:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	// Error detail on primary_language
	PrimaryLanguage *string `form:"primary_language,omitempty" json:"primary_language,omitempty" xml:"primary_language,omitempty"`
	// Error detail on categories
	Categories *string `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Error detail on biography
	Biography *string `form:"biography,omitempty" json:"biography,omitempty" xml:"biography,omitempty"`
	// Error detail on avatar
	AvatarImage *string `form:"avatar_image,omitempty" json:"avatar_image,omitempty" xml:"avatar_image,omitempty"`
	// Error detail on cover photo
	CoverPhoto *string `form:"cover_photo,omitempty" json:"cover_photo,omitempty" xml:"cover_photo,omitempty"`
}

// UpdateUserdataResponseBody is the type of the "userdatas" service
// "updateUserdata" endpoint HTTP response body.
type UpdateUserdataResponseBody struct {
	// Result of the request is success or not
	ResponseCode int `form:"response_code" json:"response_code" xml:"response_code"`
	// The detail of why result is success/fail, depend on response_code
	Detail string `form:"detail" json:"detail" xml:"detail"`
	// Error detail on realname
	RealName *string `form:"realname,omitempty" json:"realname,omitempty" xml:"realname,omitempty"`
	// Error detail on facebook_link
	FacebookLink *string `form:"facebook_link,omitempty" json:"facebook_link,omitempty" xml:"facebook_link,omitempty"`
	// Error detail on twitter_link
	TwitterLink *string `form:"twitter_link,omitempty" json:"twitter_link,omitempty" xml:"twitter_link,omitempty"`
	// Error detail on native_currency
	NativeCurrency *string `form:"native_currency,omitempty" json:"native_currency,omitempty" xml:"native_currency,omitempty"`
	// Error detail on location
	Location *string `form:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	// Error detail on primary_language
	PrimaryLanguage *string `form:"primary_language,omitempty" json:"primary_language,omitempty" xml:"primary_language,omitempty"`
	// Error detail on categories
	Categories *string `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Error detail on biography
	Biography *string `form:"biography,omitempty" json:"biography,omitempty" xml:"biography,omitempty"`
	// Error detail on avatar
	AvatarImage *string `form:"avatar_image,omitempty" json:"avatar_image,omitempty" xml:"avatar_image,omitempty"`
	// Error detail on cover photo
	CoverPhoto *string `form:"cover_photo,omitempty" json:"cover_photo,omitempty" xml:"cover_photo,omitempty"`
}

// GetUserdataResponseBody is the type of the "userdatas" service "getUserdata"
// endpoint HTTP response body.
type GetUserdataResponseBody struct {
	// Real name of the user
	RealName *string `form:"realname,omitempty" json:"realname,omitempty" xml:"realname,omitempty"`
	// Facebook link of the user
	FacebookLink *string `form:"facebook_link,omitempty" json:"facebook_link,omitempty" xml:"facebook_link,omitempty"`
	// Twitter link of the user
	TwitterLink *string `form:"twitter_link,omitempty" json:"twitter_link,omitempty" xml:"twitter_link,omitempty"`
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string `form:"native_currency,omitempty" json:"native_currency,omitempty" xml:"native_currency,omitempty"`
	// Location of the user
	Location *string `form:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	// Primary language of the user, follow ISO 639-2 standard
	PrimaryLanguage *string `form:"primary_language,omitempty" json:"primary_language,omitempty" xml:"primary_language,omitempty"`
	// The categories of user's work, separate by ,
	Categories *string `form:"categories,omitempty" json:"categories,omitempty" xml:"categories,omitempty"`
	// Biography of the user
	Biography *string `form:"biography,omitempty" json:"biography,omitempty" xml:"biography,omitempty"`
	// Avatar image of the user
	AvatarImage *UserImageUploadPayloadResponseBody `form:"avatar_image,omitempty" json:"avatar_image,omitempty" xml:"avatar_image,omitempty"`
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayloadResponseBody `form:"cover_photo,omitempty" json:"cover_photo,omitempty" xml:"cover_photo,omitempty"`
	// User's PastelID
	UserPastelID string `form:"user_pastelid" json:"user_pastelid" xml:"user_pastelid"`
	// Passphrase of the user's PastelID
	UserPastelIDPassphrase string `form:"user_pastelid_passphrase" json:"user_pastelid_passphrase" xml:"user_pastelid_passphrase"`
}

// CreateUserdataBadRequestResponseBody is the type of the "userdatas" service
// "createUserdata" endpoint HTTP response body for the "BadRequest" error.
type CreateUserdataBadRequestResponseBody struct {
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

// CreateUserdataInternalServerErrorResponseBody is the type of the "userdatas"
// service "createUserdata" endpoint HTTP response body for the
// "InternalServerError" error.
type CreateUserdataInternalServerErrorResponseBody struct {
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

// UpdateUserdataBadRequestResponseBody is the type of the "userdatas" service
// "updateUserdata" endpoint HTTP response body for the "BadRequest" error.
type UpdateUserdataBadRequestResponseBody struct {
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

// UpdateUserdataInternalServerErrorResponseBody is the type of the "userdatas"
// service "updateUserdata" endpoint HTTP response body for the
// "InternalServerError" error.
type UpdateUserdataInternalServerErrorResponseBody struct {
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

// GetUserdataBadRequestResponseBody is the type of the "userdatas" service
// "getUserdata" endpoint HTTP response body for the "BadRequest" error.
type GetUserdataBadRequestResponseBody struct {
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

// GetUserdataNotFoundResponseBody is the type of the "userdatas" service
// "getUserdata" endpoint HTTP response body for the "NotFound" error.
type GetUserdataNotFoundResponseBody struct {
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

// GetUserdataInternalServerErrorResponseBody is the type of the "userdatas"
// service "getUserdata" endpoint HTTP response body for the
// "InternalServerError" error.
type GetUserdataInternalServerErrorResponseBody struct {
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

// UserImageUploadPayloadResponseBody is used to define fields on response body
// types.
type UserImageUploadPayloadResponseBody struct {
	// File to upload (byte array of the file content)
	Content []byte `form:"content" json:"content" xml:"content"`
	// File name of the user image
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// UserImageUploadPayloadRequestBody is used to define fields on request body
// types.
type UserImageUploadPayloadRequestBody struct {
	// File to upload (byte array of the file content)
	Content []byte `form:"content,omitempty" json:"content,omitempty" xml:"content,omitempty"`
	// File name of the user image
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
}

// NewCreateUserdataResponseBody builds the HTTP response body from the result
// of the "createUserdata" endpoint of the "userdatas" service.
func NewCreateUserdataResponseBody(res *userdatas.UserdataProcessResult) *CreateUserdataResponseBody {
	body := &CreateUserdataResponseBody{
		ResponseCode:    res.ResponseCode,
		Detail:          res.Detail,
		RealName:        res.RealName,
		FacebookLink:    res.FacebookLink,
		TwitterLink:     res.TwitterLink,
		NativeCurrency:  res.NativeCurrency,
		Location:        res.Location,
		PrimaryLanguage: res.PrimaryLanguage,
		Categories:      res.Categories,
		Biography:       res.Biography,
		AvatarImage:     res.AvatarImage,
		CoverPhoto:      res.CoverPhoto,
	}
	return body
}

// NewUpdateUserdataResponseBody builds the HTTP response body from the result
// of the "updateUserdata" endpoint of the "userdatas" service.
func NewUpdateUserdataResponseBody(res *userdatas.UserdataProcessResult) *UpdateUserdataResponseBody {
	body := &UpdateUserdataResponseBody{
		ResponseCode:    res.ResponseCode,
		Detail:          res.Detail,
		RealName:        res.RealName,
		FacebookLink:    res.FacebookLink,
		TwitterLink:     res.TwitterLink,
		NativeCurrency:  res.NativeCurrency,
		Location:        res.Location,
		PrimaryLanguage: res.PrimaryLanguage,
		Categories:      res.Categories,
		Biography:       res.Biography,
		AvatarImage:     res.AvatarImage,
		CoverPhoto:      res.CoverPhoto,
	}
	return body
}

// NewGetUserdataResponseBody builds the HTTP response body from the result of
// the "getUserdata" endpoint of the "userdatas" service.
func NewGetUserdataResponseBody(res *userdatas.UserSpecifiedData) *GetUserdataResponseBody {
	body := &GetUserdataResponseBody{
		RealName:               res.RealName,
		FacebookLink:           res.FacebookLink,
		TwitterLink:            res.TwitterLink,
		NativeCurrency:         res.NativeCurrency,
		Location:               res.Location,
		PrimaryLanguage:        res.PrimaryLanguage,
		Categories:             res.Categories,
		Biography:              res.Biography,
		UserPastelID:           res.UserPastelID,
		UserPastelIDPassphrase: res.UserPastelIDPassphrase,
	}
	if res.AvatarImage != nil {
		body.AvatarImage = marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadResponseBody(res.AvatarImage)
	}
	if res.CoverPhoto != nil {
		body.CoverPhoto = marshalUserdatasUserImageUploadPayloadToUserImageUploadPayloadResponseBody(res.CoverPhoto)
	}
	return body
}

// NewCreateUserdataBadRequestResponseBody builds the HTTP response body from
// the result of the "createUserdata" endpoint of the "userdatas" service.
func NewCreateUserdataBadRequestResponseBody(res *goa.ServiceError) *CreateUserdataBadRequestResponseBody {
	body := &CreateUserdataBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewCreateUserdataInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "createUserdata" endpoint of the "userdatas"
// service.
func NewCreateUserdataInternalServerErrorResponseBody(res *goa.ServiceError) *CreateUserdataInternalServerErrorResponseBody {
	body := &CreateUserdataInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateUserdataBadRequestResponseBody builds the HTTP response body from
// the result of the "updateUserdata" endpoint of the "userdatas" service.
func NewUpdateUserdataBadRequestResponseBody(res *goa.ServiceError) *UpdateUserdataBadRequestResponseBody {
	body := &UpdateUserdataBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewUpdateUserdataInternalServerErrorResponseBody builds the HTTP response
// body from the result of the "updateUserdata" endpoint of the "userdatas"
// service.
func NewUpdateUserdataInternalServerErrorResponseBody(res *goa.ServiceError) *UpdateUserdataInternalServerErrorResponseBody {
	body := &UpdateUserdataInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetUserdataBadRequestResponseBody builds the HTTP response body from the
// result of the "getUserdata" endpoint of the "userdatas" service.
func NewGetUserdataBadRequestResponseBody(res *goa.ServiceError) *GetUserdataBadRequestResponseBody {
	body := &GetUserdataBadRequestResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetUserdataNotFoundResponseBody builds the HTTP response body from the
// result of the "getUserdata" endpoint of the "userdatas" service.
func NewGetUserdataNotFoundResponseBody(res *goa.ServiceError) *GetUserdataNotFoundResponseBody {
	body := &GetUserdataNotFoundResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewGetUserdataInternalServerErrorResponseBody builds the HTTP response body
// from the result of the "getUserdata" endpoint of the "userdatas" service.
func NewGetUserdataInternalServerErrorResponseBody(res *goa.ServiceError) *GetUserdataInternalServerErrorResponseBody {
	body := &GetUserdataInternalServerErrorResponseBody{
		Name:      res.Name,
		ID:        res.ID,
		Message:   res.Message,
		Temporary: res.Temporary,
		Timeout:   res.Timeout,
		Fault:     res.Fault,
	}
	return body
}

// NewCreateUserdataPayload builds a userdatas service createUserdata endpoint
// payload.
func NewCreateUserdataPayload(body *CreateUserdataRequestBody) *userdatas.CreateUserdataPayload {
	v := &userdatas.CreateUserdataPayload{
		RealName:               body.RealName,
		FacebookLink:           body.FacebookLink,
		TwitterLink:            body.TwitterLink,
		NativeCurrency:         body.NativeCurrency,
		Location:               body.Location,
		PrimaryLanguage:        body.PrimaryLanguage,
		Categories:             body.Categories,
		Biography:              body.Biography,
		UserPastelID:           *body.UserPastelID,
		UserPastelIDPassphrase: *body.UserPastelIDPassphrase,
	}
	if body.AvatarImage != nil {
		v.AvatarImage = unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.AvatarImage)
	}
	if body.CoverPhoto != nil {
		v.CoverPhoto = unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.CoverPhoto)
	}

	return v
}

// NewUpdateUserdataPayload builds a userdatas service updateUserdata endpoint
// payload.
func NewUpdateUserdataPayload(body *UpdateUserdataRequestBody) *userdatas.UpdateUserdataPayload {
	v := &userdatas.UpdateUserdataPayload{
		RealName:               body.RealName,
		FacebookLink:           body.FacebookLink,
		TwitterLink:            body.TwitterLink,
		NativeCurrency:         body.NativeCurrency,
		Location:               body.Location,
		PrimaryLanguage:        body.PrimaryLanguage,
		Categories:             body.Categories,
		Biography:              body.Biography,
		UserPastelID:           *body.UserPastelID,
		UserPastelIDPassphrase: *body.UserPastelIDPassphrase,
	}
	if body.AvatarImage != nil {
		v.AvatarImage = unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.AvatarImage)
	}
	if body.CoverPhoto != nil {
		v.CoverPhoto = unmarshalUserImageUploadPayloadRequestBodyToUserdatasUserImageUploadPayload(body.CoverPhoto)
	}

	return v
}

// NewGetUserdataPayload builds a userdatas service getUserdata endpoint
// payload.
func NewGetUserdataPayload(pastelid string) *userdatas.GetUserdataPayload {
	v := &userdatas.GetUserdataPayload{}
	v.Pastelid = pastelid

	return v
}

// ValidateCreateUserdataRequestBody runs the validations defined on
// CreateUserdataRequestBody
func ValidateCreateUserdataRequestBody(body *CreateUserdataRequestBody) (err error) {
	if body.UserPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("user_pastelid", "body"))
	}
	if body.UserPastelIDPassphrase == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("user_pastelid_passphrase", "body"))
	}
	if body.RealName != nil {
		if utf8.RuneCountInString(*body.RealName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.realname", *body.RealName, utf8.RuneCountInString(*body.RealName), 256, false))
		}
	}
	if body.FacebookLink != nil {
		if utf8.RuneCountInString(*body.FacebookLink) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.facebook_link", *body.FacebookLink, utf8.RuneCountInString(*body.FacebookLink), 128, false))
		}
	}
	if body.TwitterLink != nil {
		if utf8.RuneCountInString(*body.TwitterLink) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.twitter_link", *body.TwitterLink, utf8.RuneCountInString(*body.TwitterLink), 128, false))
		}
	}
	if body.NativeCurrency != nil {
		if utf8.RuneCountInString(*body.NativeCurrency) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, true))
		}
	}
	if body.NativeCurrency != nil {
		if utf8.RuneCountInString(*body.NativeCurrency) > 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, false))
		}
	}
	if body.Location != nil {
		if utf8.RuneCountInString(*body.Location) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.location", *body.Location, utf8.RuneCountInString(*body.Location), 256, false))
		}
	}
	if body.PrimaryLanguage != nil {
		if utf8.RuneCountInString(*body.PrimaryLanguage) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.primary_language", *body.PrimaryLanguage, utf8.RuneCountInString(*body.PrimaryLanguage), 30, false))
		}
	}
	if body.Biography != nil {
		if utf8.RuneCountInString(*body.Biography) > 1024 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.biography", *body.Biography, utf8.RuneCountInString(*body.Biography), 1024, false))
		}
	}
	if body.AvatarImage != nil {
		if err2 := ValidateUserImageUploadPayloadRequestBody(body.AvatarImage); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if body.CoverPhoto != nil {
		if err2 := ValidateUserImageUploadPayloadRequestBody(body.CoverPhoto); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if body.UserPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.user_pastelid", *body.UserPastelID, "^[a-zA-Z0-9]+$"))
	}
	if body.UserPastelID != nil {
		if utf8.RuneCountInString(*body.UserPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_pastelid", *body.UserPastelID, utf8.RuneCountInString(*body.UserPastelID), 86, true))
		}
	}
	if body.UserPastelID != nil {
		if utf8.RuneCountInString(*body.UserPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_pastelid", *body.UserPastelID, utf8.RuneCountInString(*body.UserPastelID), 86, false))
		}
	}
	return
}

// ValidateUpdateUserdataRequestBody runs the validations defined on
// UpdateUserdataRequestBody
func ValidateUpdateUserdataRequestBody(body *UpdateUserdataRequestBody) (err error) {
	if body.UserPastelID == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("user_pastelid", "body"))
	}
	if body.UserPastelIDPassphrase == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("user_pastelid_passphrase", "body"))
	}
	if body.RealName != nil {
		if utf8.RuneCountInString(*body.RealName) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.realname", *body.RealName, utf8.RuneCountInString(*body.RealName), 256, false))
		}
	}
	if body.FacebookLink != nil {
		if utf8.RuneCountInString(*body.FacebookLink) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.facebook_link", *body.FacebookLink, utf8.RuneCountInString(*body.FacebookLink), 128, false))
		}
	}
	if body.TwitterLink != nil {
		if utf8.RuneCountInString(*body.TwitterLink) > 128 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.twitter_link", *body.TwitterLink, utf8.RuneCountInString(*body.TwitterLink), 128, false))
		}
	}
	if body.NativeCurrency != nil {
		if utf8.RuneCountInString(*body.NativeCurrency) < 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, true))
		}
	}
	if body.NativeCurrency != nil {
		if utf8.RuneCountInString(*body.NativeCurrency) > 3 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.native_currency", *body.NativeCurrency, utf8.RuneCountInString(*body.NativeCurrency), 3, false))
		}
	}
	if body.Location != nil {
		if utf8.RuneCountInString(*body.Location) > 256 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.location", *body.Location, utf8.RuneCountInString(*body.Location), 256, false))
		}
	}
	if body.PrimaryLanguage != nil {
		if utf8.RuneCountInString(*body.PrimaryLanguage) > 30 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.primary_language", *body.PrimaryLanguage, utf8.RuneCountInString(*body.PrimaryLanguage), 30, false))
		}
	}
	if body.Biography != nil {
		if utf8.RuneCountInString(*body.Biography) > 1024 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.biography", *body.Biography, utf8.RuneCountInString(*body.Biography), 1024, false))
		}
	}
	if body.AvatarImage != nil {
		if err2 := ValidateUserImageUploadPayloadRequestBody(body.AvatarImage); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if body.CoverPhoto != nil {
		if err2 := ValidateUserImageUploadPayloadRequestBody(body.CoverPhoto); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if body.UserPastelID != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.user_pastelid", *body.UserPastelID, "^[a-zA-Z0-9]+$"))
	}
	if body.UserPastelID != nil {
		if utf8.RuneCountInString(*body.UserPastelID) < 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_pastelid", *body.UserPastelID, utf8.RuneCountInString(*body.UserPastelID), 86, true))
		}
	}
	if body.UserPastelID != nil {
		if utf8.RuneCountInString(*body.UserPastelID) > 86 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.user_pastelid", *body.UserPastelID, utf8.RuneCountInString(*body.UserPastelID), 86, false))
		}
	}
	return
}

// ValidateUserImageUploadPayloadRequestBody runs the validations defined on
// UserImageUploadPayloadRequestBody
func ValidateUserImageUploadPayloadRequestBody(body *UserImageUploadPayloadRequestBody) (err error) {
	if body.Content == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("content", "body"))
	}
	if body.Filename != nil {
		err = goa.MergeErrors(err, goa.ValidatePattern("body.filename", *body.Filename, "^.*\\.(png|PNG|jpeg|JPEG|jpg|JPG)$"))
	}
	return
}
