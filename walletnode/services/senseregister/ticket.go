package senseregister

import "github.com/pastelnetwork/gonode/common/service/artwork"

// Request represents a request for the registration sense
type Request struct {
	Image       *artwork.File `json:"image"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
}
