package walletnode

import (
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
)

// Service represents grpc service for interaction between walletnode and supernode.
type Service struct {
	pb.UnimplementedWalletNodeServer
	artworkRegister *artworkregister.Service
}

// RegisterArtowrk is a stream for registration artwork.
func (service *Service) RegisterArtowrk(stream pb.WalletNode_RegisterArtowrkServer) error {
	registerArtwork := newRegisterArtowrk(stream)
	return registerArtwork.start(stream.Context(), service.artworkRegister)
}

// Desc returns a description of the service.
func (service *Service) Desc() *grpc.ServiceDesc {
	return &pb.WalletNode_ServiceDesc
}

// NewService returns a new Service instance.
func NewService(artworkRegister *artworkregister.Service) *Service {
	return &Service{
		artworkRegister: artworkRegister,
	}
}
