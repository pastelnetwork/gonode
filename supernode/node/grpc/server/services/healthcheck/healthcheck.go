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

func (service *HealthCheck) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	// echos received message
	return &pb.PingReply{Reply: req.Msg}, nil
}

// Desc returns a description of the service.
func (service *HealthCheck) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheck_ServiceDesc
}

// NewDownloadArtwork returns a new DownloadArtwork instance.
func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}
