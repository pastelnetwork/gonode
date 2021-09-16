package healthcheck

import (
	"context"
	"encoding/json"

	"github.com/pastelnetwork/gonode/common/errors"
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
	pb.UnimplementedHealthCheckServer
}

// Ping will send a message to and get back a reply from supernode
func (service *HealthCheck) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingReply, error) {
	stats, err := service.Stats(ctx)
	if err != nil {
		return &pb.PingReply{Reply: ""}, errors.Errorf("failed to Stats(): %w", err)
	}

	jsonData, err := json.Marshal(stats)
	if err != nil {
		return &pb.PingReply{Reply: ""}, errors.Errorf("failed to Marshal(): %w", err)
	}

	// echos received message
	return &pb.PingReply{Reply: string(jsonData)}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck(mngr StatsMngr) *HealthCheck {
	return &HealthCheck{
		StatsMngr: mngr,
	}
}
