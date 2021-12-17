package senseregister

import "github.com/pastelnetwork/gonode/common/service/artwork"

type Request struct {
	Image       *artwork.File `json:"image"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
}

// GetEstimatedFeeRequest
type GetEstimatedFeeRequest struct {
	Image *artwork.File `json:"image"`
}
