package supernode

import (
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
)

// Service represents grpc service for interaction between two supernodes.
type Service struct {
	pb.UnimplementedSuperNodeServer

	artworkRegister *artworkregister.Service
}

// RegisterArtowrk is a stream for registration artwork.
func (service *Service) RegisterArtowrk(stream pb.SuperNode_RegisterArtowrkServer) error {
	registerArtwork := newRegisterArtowrk(stream)
	return registerArtwork.start(stream.Context(), service.artworkRegister)
}

// Desc returns a description of the service.
func (service *Service) Desc() *grpc.ServiceDesc {
	return &pb.SuperNode_ServiceDesc
}

// NewService returns a new Service instance.
func NewService(artworkRegister *artworkregister.Service) *Service {
	return &Service{
		artworkRegister: artworkRegister,
	}
}
