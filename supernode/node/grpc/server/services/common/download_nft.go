package common

import (
	"github.com/pastelnetwork/gonode/supernode/services/download"
)

// DownloadNft represents common grpc service for downloading NFTs.
type DownloadNft struct {
	*download.NftDownloaderService
}

// NewDownloadNft returns a new DownloadNft instance.
func NewDownloadNft(service *download.NftDownloaderService) *DownloadNft {
	return &DownloadNft{
		NftDownloaderService: service,
	}
}
