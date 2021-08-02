// Code generated by goa v3.3.1, DO NOT EDIT.
//
// userdatas service
//
// Command:
// $ goa gen github.com/pastelnetwork/gonode/walletnode/api/design

package userdatas

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Pastel Process User Specified Data
type Service interface {
	// Create new user data
	CreateUserdata(context.Context, *CreateUserdataPayload) (res *UserdataProcessResult, err error)
	// Update user data for an existing user
	UpdateUserdata(context.Context, *UpdateUserdataPayload) (res *UserdataProcessResult, err error)
	// Gets the Userdata detail
	UserdataGet(context.Context, *UserdataGetPayload) (res *UserSpecifiedData, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "userdatas"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [3]string{"createUserdata", "updateUserdata", "userdataGet"}

// CreateUserdataPayload is the payload type of the userdatas service
// createUserdata method.
type CreateUserdataPayload struct {
	// Real name of the user
	Realname *string
	// Facebook link of the user
	FacebookLink *string
	// Twitter link of the user
	TwitterLink *string
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string
	// Location of the user
	Location *string
	// Primary language of the user
	PrimaryLanguage *string
	// The categories of user's work
	Categories *string
	// Biography of the user
	Biography *string
	// Avatar image of the user
	AvatarImage *UserImageUploadPayload
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayload
	// Artist's PastelID
	ArtistPastelID string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
}

// UserdataProcessResult is the result type of the userdatas service
// createUserdata method.
type UserdataProcessResult struct {
	// Result of the request is success or not
	ResponseCode int
	// The detail of why result is success/fail, depend on response_code
	Detail string
	// Error detail on realname
	Realname *string
	// Error detail on facebook_link
	FacebookLink *string
	// Error detail on twitter_link
	TwitterLink *string
	// Error detail on native_currency
	NativeCurrency *string
	// Error detail on location
	Location *string
	// Error detail on primary_language
	PrimaryLanguage *string
	// Error detail on categories
	Categories *string
	// Error detail on biography
	Biography *string
	// Error detail on avatar
	AvatarImage *string
	// Error detail on cover photo
	CoverPhoto *string
}

// UpdateUserdataPayload is the payload type of the userdatas service
// updateUserdata method.
type UpdateUserdataPayload struct {
	// Real name of the user
	Realname *string
	// Facebook link of the user
	FacebookLink *string
	// Twitter link of the user
	TwitterLink *string
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string
	// Location of the user
	Location *string
	// Primary language of the user
	PrimaryLanguage *string
	// The categories of user's work
	Categories *string
	// Biography of the user
	Biography *string
	// Avatar image of the user
	AvatarImage *UserImageUploadPayload
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayload
	// Artist's PastelID
	ArtistPastelID string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
}

// UserdataGetPayload is the payload type of the userdatas service userdataGet
// method.
type UserdataGetPayload struct {
	// Artist's PastelID
	Pastelid string
}

// UserSpecifiedData is the result type of the userdatas service userdataGet
// method.
type UserSpecifiedData struct {
	// Real name of the user
	Realname *string
	// Facebook link of the user
	FacebookLink *string
	// Twitter link of the user
	TwitterLink *string
	// Native currency of user in ISO 4217 Alphabetic Code
	NativeCurrency *string
	// Location of the user
	Location *string
	// Primary language of the user
	PrimaryLanguage *string
	// The categories of user's work
	Categories *string
	// Biography of the user
	Biography *string
	// Avatar image of the user
	AvatarImage *UserImageUploadPayload
	// Cover photo of the user
	CoverPhoto *UserImageUploadPayload
	// Artist's PastelID
	ArtistPastelID string
	// Passphrase of the artist's PastelID
	ArtistPastelIDPassphrase string
}

// User image upload payload
type UserImageUploadPayload struct {
	// File to upload
	Content []byte
	// File name of the user image
	Filename *string
}

// MakeBadRequest builds a goa.ServiceError from an error.
func MakeBadRequest(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "BadRequest",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "NotFound",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeInternalServerError builds a goa.ServiceError from an error.
func MakeInternalServerError(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "InternalServerError",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
