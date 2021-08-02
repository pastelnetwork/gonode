package common

import (
	"github.com/pastelnetwork/gonode/supernode/services/artworkdownload"
)

// DownloadArtwork represents common grpc service for downloading artwork.
type DownloadArtwork struct {
	*artworkdownload.Service
}

// NewDownloadArtwork returns a new DownloadArtwork instance.
func NewDownloadArtwork(service *artworkdownload.Service) *DownloadArtwork {
	return &DownloadArtwork{
		Service: service,
	}
}
