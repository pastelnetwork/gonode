package senseregister

import "github.com/pastelnetwork/gonode/common/service/artwork"

// Request represents a request for the registration sense
type Request struct {
	Image                 *artwork.File `json:"image"`
	AppPastelID           string        `json:"app_pastel_id"`
	AppPastelIDPassphrase string        `json:"app_pastel_id_passphrase"`
}
