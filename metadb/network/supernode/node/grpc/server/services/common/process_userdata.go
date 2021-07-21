package common

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb/network/proto"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
	"google.golang.org/grpc/metadata"
)

// ProcessUserdata represents common grpc service for registration artwork.
type ProcessUserdata struct {
	*processuserdata.Service
}

// SessID retrieves SessID from the metadata.
func (service *ProcessUserdata) SessID(ctx context.Context) (string, bool) {
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
func (service *ProcessUserdata) TaskFromMD(ctx context.Context) (*processuserdata.Task, error) {
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

// NewProcessUserdata returns a new ProcessUserdata instance.
func NewProcessUserdata(service *processuserdata.Service) *ProcessUserdata {
	return &ProcessUserdata{
		Service: service,
	}
}
