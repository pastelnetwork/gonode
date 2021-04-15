package arts

import (
	"net/http"

	arts "github.com/pastelnetwork/walletnode/api/gen/arts"
	"github.com/pastelnetwork/walletnode/api/services"
)

var errBadRequest = func(val ...interface{}) *arts.BadRequest {
	return &arts.BadRequest{
		InnerError: services.InnerError(http.StatusBadRequest, val...),
	}
}
