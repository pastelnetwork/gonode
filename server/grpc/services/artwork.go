package services

import (
	"fmt"
	"io"
	"log"

	pb "github.com/pastelnetwork/supernode-proto"
	"github.com/pastelnetwork/supernode/services/artworkregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Artwork struct {
	pb.UnimplementedArtworkServer

	register *artworkregister.Service
}

func (regsiger *Artwork) Register(stream pb.Artwork_RegisterServer) error {
	ctx := stream.Context()

	if peer, ok := peer.FromContext(ctx); ok {
		fmt.Println(peer.Addr.String())
		//fmt.Println(peer.AuthInfo.AuthType())
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("stop-2")
			return nil
		}
		if err != nil {
			if status.Code(err) == codes.Canceled {
				log.Println("stream closed (context cancelled)")
				return nil
			}
			fmt.Println("stop-3", err)
			return err
		}

		select {
		case <-ctx.Done():
			fmt.Println("stop-1")
			return nil
		default:
			switch req := req.GetTestOneof().(type) {
			case *pb.RegisterRequest_PrimaryNode:
				fmt.Println("primary", req.PrimaryNode.ConnID)
			case *pb.RegisterRequest_SecondaryNode:
				fmt.Println("secondary", req.SecondaryNode.ConnID)
			default:
				fmt.Println("wrong call")
			}
		}
	}
	return nil
}

// Desc returns a description of the service.
func (regsiger *Artwork) Desc() *grpc.ServiceDesc {
	return &pb.Artwork_ServiceDesc
}

// NewArtwork returns a new Artwork instance.
func NewArtwork(register *artworkregister.Service) *Artwork {
	return &Artwork{
		register: register,
	}
}
