package walletnode

import (
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
)

// DownloadArtwork represents grpc service for downloading artwork.
type DownloadArtwork struct {
	pb.UnimplementedDownloadArtworkServer

	*common.DownloadArtwork
}

// Download downloads artwork by given txid, timestamp and signature.
func (service *DownloadArtwork) Download(m *DownloadRequest, stream DownloadArtwork_DownloadServer) error {
	// Create new task

	// Call task download

	// Return restored image back to WalletNode

	return nil
}

// Desc returns a description of the service.
func (service *DownloadArtwork) Desc() *grpc.ServiceDesc {
	return &pb.DownloadArtwork_ServiceDesc
}

// NewDownloadArtwork returns a new DownloadArtwork instance.
func NewDownloadArtwork(service *artworkregister.Service) *DownloadArtwork {
	return &DownloadArtwork{
		DownloadArtwork: common.NewDownloadArtwork(service),
	}
}
