package common

import (
	"github.com/pastelnetwork/gonode/supernode/services/nftdownload"
)

// DownloadNft represents common grpc service for downloading NFTs.
type DownloadNft struct {
	*nftdownload.NftDownloadService
}

// NewDownloadNft returns a new DownloadNft instance.
func NewDownloadNft(service *nftdownload.NftDownloadService) *DownloadNft {
	return &DownloadNft{
		NftDownloadService: service,
	}
}
