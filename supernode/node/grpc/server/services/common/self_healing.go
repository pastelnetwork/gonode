package common

import (
	"context"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/selfhealing"
	"google.golang.org/grpc/metadata"
)

// SelfHealingChallenge represents common grpc service for self-healing challenges.
type SelfHealingChallenge struct {
	*selfhealing.SHService
}

// SessID retrieves SessID from the metadata.
func (service *SelfHealingChallenge) SessID(ctx context.Context) (string, bool) {
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

// NewSelfHealingChallenge returns a new SelfHealing instance.
func NewSelfHealingChallenge(service *selfhealing.SHService) *SelfHealingChallenge {
	return &SelfHealingChallenge{
		SHService: service,
	}
}
