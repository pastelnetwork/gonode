package bridge

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/bridge"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/download"
	"google.golang.org/grpc"
)

// DownloadData represents grpc service for downloading NFT.
type DownloadData struct {
	pb.UnimplementedDownloadDataServer

	*common.DownloadNft
}

// DownloadThumbnail returns thumbnail of given hash
func (service *DownloadData) DownloadThumbnail(ctx context.Context, req *pb.DownloadThumbnailRequest) (*pb.DownloadThumbnailReply, error) {
	log.WithContext(ctx).Println("Received download thumbnail request")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Create new task
	task := service.NewNftDownloadingTask()
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	// Call task download thumbnail
	log.WithContext(ctx).WithField("txid", req.Txid).WithField("numnails", req.Numnails).Println("Downloading thumbnail")
	data, err := task.DownloadThumbnail(ctx, req.Txid, req.Numnails)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadThumbnailReply{
		Thumbnailone: data[0],
		Thumbnailtwo: data[1],
	}, nil
}

// DownloadDDAndFingerprints returns ONE DDAndFingerprints file out of the 50 that exist as long as it decodes properly.
func (service *DownloadData) DownloadDDAndFingerprints(ctx context.Context, req *pb.DownloadDDAndFingerprintsRequest) (*pb.DownloadDDAndFingerprintsReply, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Create new task
	task := service.NewNftDownloadingTask()
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	// Call task download thumbnail
	data, err := task.DownloadDDAndFingerprints(ctx, req.Txid)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadDDAndFingerprintsReply{
		File: data,
	}, nil
}

// Desc returns a description of the service.
func (service *DownloadData) Desc() *grpc.ServiceDesc {
	return &pb.DownloadData_ServiceDesc
}

// NewDownloadData returns a new DownloadData instance.
func NewDownloadData(service *download.NftDownloaderService) *DownloadData {
	return &DownloadData{
		DownloadNft: common.NewDownloadNft(service),
	}
}
