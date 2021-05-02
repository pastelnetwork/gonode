package services

import (
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/services/artworkregister"
	"google.golang.org/grpc"
)

// WalletNode represents grpc server
type WalletNode struct {
	pb.UnimplementedWalletNodeServer
	artworkRegister *artworkregister.Service
}

// RegisterArtowrk is responsible for communication between walletnote and supernodes.
func (server *WalletNode) RegisterArtowrk(stream pb.WalletNode_RegisterArtowrkServer) error {
	registerArtwork := newRegisterArtowrk(stream)
	return registerArtwork.start(stream.Context(), server.artworkRegister)
}

// Desc returns a description of the server.
func (server *WalletNode) Desc() *grpc.ServiceDesc {
	return &pb.WalletNode_ServiceDesc
}

// NewWalletNode returns a new WalletNode instance.
func NewWalletNode(artworkRegister *artworkregister.Service) *WalletNode {
	return &WalletNode{
		artworkRegister: artworkRegister,
	}
}
