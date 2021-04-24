package services

import (
	"context"

	pb "github.com/pastelnetwork/supernode-proto"
)

type ArtworkRegister struct {
	pb.UnimplementedArtworkRegisterServer
}

func (regsiger *ArtworkRegister) Connect(context.Context, *pb.ConnectRequest) (*pb.ConnectReply, error) {

	return nil, nil
}

func (regsiger *ArtworkRegister) AcceptConnection(context.Context, *pb.AcceptConnectionRequest) (*pb.AcceptConnectionReply, error) {

	return nil, nil
}

func NewArtworkRegister() pb.ArtworkRegisterServer {
	return &ArtworkRegister{}
}
