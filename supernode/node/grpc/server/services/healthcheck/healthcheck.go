package healthcheck

import (
	"context"

	pb "github.com/pastelnetwork/gonode/proto/healthcheck"

	"google.golang.org/grpc"
)

// HealthCheck represents grpc service for supernode healthcheck
type HealthCheck struct {
	pb.UnimplementedHealthCheckServer
}

// Ping will send a message to and get back a reply from supernode
func (service *HealthCheck) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	// echos received message
	return &pb.PingReply{Reply: req.Msg}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewHealthCheck returns a new HealthCheck instance.
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}
