package bridge

import (
	"context"

	"github.com/pastelnetwork/gonode/bridge/services/download"

	pb "github.com/pastelnetwork/gonode/proto/bridge"

	"google.golang.org/grpc"
)

// Service represents grpc service for rq server.
type Service struct {
	pb.UnimplementedDownloadDataServer
	download download.Service
}

// NewService returns a new Bridge Service instance.
func NewService(download download.Service) *Service {
	return &Service{
		download: download,
	}
}

// DownloadThumbnail downloads thumbnail
func (s *Service) DownloadThumbnail(ctx context.Context, req *pb.DownloadThumbnailRequest) (r *pb.DownloadThumbnailReply, err error) {
	resp, err := s.download.FetchThumbnail(ctx, req.Txid, int(req.Numnails))
	if err != nil {
		return nil, err
	}

	return &pb.DownloadThumbnailReply{Thumbnailone: resp[0], Thumbnailtwo: resp[1]}, nil
}

// DownloadDDAndFingerprints downloads dd & fp
func (s *Service) DownloadDDAndFingerprints(ctx context.Context, req *pb.DownloadDDAndFingerprintsRequest) (r *pb.DownloadDDAndFingerprintsReply, err error) {
	resp, err := s.download.FetchDupeDetectionData(ctx, req.Txid)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadDDAndFingerprintsReply{File: resp}, nil
}

// Ping returns a description of the service.
func (s *Service) Ping(ctx context.Context, req *pb.PingRequest) (r *pb.PingReply, err error) {
	return &pb.PingReply{}, nil
}

// Desc returns a description of the service.
func (s *Service) Desc() *grpc.ServiceDesc {
	return &pb.DownloadData_ServiceDesc
}
