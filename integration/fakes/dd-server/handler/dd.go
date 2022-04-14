package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/integration/fakes/dd-server/dupedetection"

	"github.com/pastelnetwork/gonode/integration/fakes/common/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// DDService represents grpc service for dd server.
type DDService struct {
	pb.UnimplementedDupeDetectionServerServer
	store storage.Store
}

// SessID retrieves SessID from the metadata.
func (service *DDService) SessID(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	mdVals := md.Get("sessID")
	if len(mdVals) == 0 {
		return "", false
	}
	return mdVals[0], true
}

// ImageRarenessScore request
func (service *DDService) ImageRarenessScore(ctx context.Context, req *pb.RarenessScoreRequest) (*pb.ImageRarenessScoreReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("ImageRarenessScore request")

	key := "rareness" + "*" + req.PastelIdOfSubmitter
	data, err := service.store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %w", err)
	}

	rep := &pb.ImageRarenessScoreReply{}
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, fmt.Errorf("unable to decode data: %w", err)
	}

	return rep, nil
}

// Desc returns a description of the service.
func (service *DDService) Desc() *grpc.ServiceDesc {
	return &pb.DupeDetectionServer_ServiceDesc
}

// NewDDService returns a new DDService instance.
func NewDDService(store storage.Store) *DDService {
	return &DDService{
		store: store,
	}
}
