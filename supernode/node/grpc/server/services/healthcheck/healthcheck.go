package healthcheck

import (
	"context"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/p2p"
	pb "github.com/pastelnetwork/gonode/proto/healthcheck"
	"google.golang.org/grpc"
)

// StatsMngr is interface of StatsManger, return stats of system
type StatsMngr interface {
	// Stats returns stats of system
	Stats(ctx context.Context) (map[string]interface{}, error)
}

// HealthCheck represents grpc service for supernode healthcheck
type HealthCheck struct {
	StatsMngr
	p2pService p2p.P2P
	pb.UnimplementedHealthCheckServer
}

// Status will send a message to and get back a reply from supernode
func (service *HealthCheck) Status(ctx context.Context, _ *pb.StatusRequest) (*pb.StatusReply, error) {
	stats, err := service.Stats(ctx)
	if err != nil {
		return &pb.StatusReply{StatusInJson: ""}, errors.Errorf("failed to Stats(): %w", err)
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		return &pb.StatusReply{StatusInJson: ""}, errors.Errorf("failed to Marshal(): %w", err)
	}

	// echos received message
	return &pb.StatusReply{StatusInJson: string(jsonData)}, nil
}

// P2PStore store a data, and return given key
func (service *HealthCheck) P2PStore(ctx context.Context, in *pb.P2PStoreRequest) (*pb.P2PStoreReply, error) {
	value := in.GetValue()
	key, err := service.p2pService.Store(ctx, value)

	if err != nil {
		return &pb.P2PStoreReply{Key: ""}, err
	}
	return &pb.P2PStoreReply{Key: key}, nil
}

// P2PRetrieve get data of given key
func (service *HealthCheck) P2PRetrieve(ctx context.Context, in *pb.P2PRetrieveRequest) (*pb.P2PRetrieveReply, error) {
	key := in.GetKey()
	value, err := service.p2pService.Retrieve(ctx, key)

	if err != nil {
		return &pb.P2PRetrieveReply{Value: nil}, err
	}
	return &pb.P2PRetrieveReply{Value: value}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck(mngr StatsMngr, p2pService p2p.P2P) *HealthCheck {
	return &HealthCheck{
		StatsMngr:  mngr,
		p2pService: p2pService,
	}
}
