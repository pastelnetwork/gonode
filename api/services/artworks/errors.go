package artworks

import (
	"net/http"

	"github.com/pastelnetwork/walletnode/api/gen/artworks"
	"github.com/pastelnetwork/walletnode/api/services"
)

var errBadRequest = func(val ...interface{}) *artworks.BadRequest {
	return &artworks.BadRequest{
		InnerError: services.InnerError(http.StatusBadRequest, val...),
	}
}

var errInternalServerError = func(val ...interface{}) *artworks.InternalServerError {
	return &artworks.InternalServerError{
		InnerError: services.InnerError(http.StatusInternalServerError, val...),
	}
}
