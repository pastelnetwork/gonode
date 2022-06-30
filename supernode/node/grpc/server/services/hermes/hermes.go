package hermes

import (
	"context"
	"errors"
	"fmt"

	"github.com/pastelnetwork/gonode/p2p"
	pb "github.com/pastelnetwork/gonode/proto/hermes"
	"google.golang.org/grpc"
)

// Hermes is  used for providing hermes access to p2p network
type Hermes struct {
	p2pClient p2p.Client
	pb.UnimplementedHermesP2PServer
}

// Retrieve will retrieve data from p2p
func (service *Hermes) Retrieve(ctx context.Context, req *pb.RetrieveRequest) (*pb.RetrieveReply, error) {
	data, err := service.p2pClient.Retrieve(ctx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("retieve err: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("retrieve; data not found")
	}

	return &pb.RetrieveReply{Data: data}, nil
}

// Delete will delete data from p2p
func (service *Hermes) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteReply, error) {
	if err := service.p2pClient.Delete(ctx, req.Key); err != nil {
		return nil, fmt.Errorf("p2p delete failed: %w", err)
	}

	return &pb.DeleteReply{}, nil
}

// Desc returns a description of the service.
func (service *Hermes) Desc() *grpc.ServiceDesc {
	return &pb.HermesP2P_ServiceDesc
}

// NewHermes returns a new HealthCheck instance.
func NewHermes(p2pClient p2p.Client) *Hermes {
	return &Hermes{
		p2pClient: p2pClient,
	}
}
