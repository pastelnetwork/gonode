package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/proto"
	"github.com/pastelnetwork/gonode/supernode/services/nftregister"
	"google.golang.org/grpc/metadata"
)

// RegisterNft represents common grpc service for registration NFTs.
type RegisterNft struct {
	*nftregister.NftRegistrationService
}

// SessID retrieves SessID from the metadata.
func (service *RegisterNft) SessID(ctx context.Context) (string, bool) {
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

// TaskFromMD returns task by SessID from the metadata.
func (service *RegisterNft) TaskFromMD(ctx context.Context) (*nftregister.NftRegistrationTask, error) {
	sessID, ok := service.SessID(ctx)
	if !ok {
		return nil, errors.New("not found sessID in metadata")
	}

	task := service.Task(sessID)
	if task == nil {
		return nil, errors.Errorf("not found %q task", sessID)
	}
	return task, nil
}

// NewRegisterNft returns a new RegisterNft instance.
func NewRegisterNft(service *nftregister.NftRegistrationService) *RegisterNft {
	return &RegisterNft{
		NftRegistrationService: service,
	}
}
