package handler

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/pastelnetwork/gonode/integration/fakes/rq-service/raptorq"

	"github.com/pastelnetwork/gonode/integration/fakes/common/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RQService represents grpc service for rq server.
type RQService struct {
	pb.UnimplementedRaptorQServer
	store storage.Store
}

// SessID retrieves SessID from the metadata.
func (service *RQService) SessID(ctx context.Context) (string, bool) {
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

// Desc returns a description of the service.
func (service *RQService) Desc() *grpc.ServiceDesc {
	return &pb.RaptorQ_ServiceDesc
}

// Encode ...
func (service *RQService) Encode(ctx context.Context, req *pb.EncodeRequest) (*pb.EncodeReply, error) {
	key := "encode" + "*"
	data, err := service.store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %w", err)
	}

	rep := &pb.EncodeReply{}
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, fmt.Errorf("unable to decode data: %w", err)
	}

	return rep, nil
}

// EncodeMetaData ...
func (service *RQService) EncodeMetaData(ctx context.Context, req *pb.EncodeMetaDataRequest) (*pb.EncodeMetaDataReply, error) {
	key := "encodemetadata" + "*" + req.PastelId

	data, err := service.store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %w", err)
	}

	rep := &pb.EncodeMetaDataReply{}
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, fmt.Errorf("unable to decode data: %w", err)
	}

	return rep, nil
}

// Decode ...
func (service *RQService) Decode(ctx context.Context, req *pb.DecodeRequest) (*pb.DecodeReply, error) {
	key := "decode" + "*"
	data, err := service.store.Get(key)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch data: %w", err)
	}

	rep := &pb.DecodeReply{}
	if err := json.Unmarshal(data, &rep); err != nil {
		return nil, fmt.Errorf("unable to decode data: %w", err)
	}

	return rep, nil
}

// NewRQService returns a new DDService instance.
func NewRQService(store storage.Store) *RQService {
	return &RQService{
		store: store,
	}
}
