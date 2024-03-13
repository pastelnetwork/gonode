package common

import (
	"context"

	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/healthcheckchallenge"
	"google.golang.org/grpc/metadata"
)

// HealthCheckChallenge represents common grpc service for health checks
type HealthCheckChallenge struct {
	*healthcheckchallenge.HCService
}

// SessID retrieves SessID from the metadata.
func (service *HealthCheckChallenge) SessID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	mdVals := md.Get(proto.MetadataKeySessID)
	if len(mdVals) == 0 {
		return "", false
	}
	return mdVals[0], true
}

// NewHealthCheckChallenge returns a new HealthCheckChallenge instance.
func NewHealthCheckChallenge(service *healthcheckchallenge.HCService) *HealthCheckChallenge {
	return &HealthCheckChallenge{
		HCService: service,
	}
}
