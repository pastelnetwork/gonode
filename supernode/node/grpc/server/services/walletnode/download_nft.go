package walletnode

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/nftdownload"
	"google.golang.org/grpc"
)

const (
	downloadImageBufferSize = 32 * 1024
)

// DownloadNft represents grpc service for downloading NFT.
type DownloadNft struct {
	pb.UnimplementedDownloadNftServer

	*common.DownloadNft
}

// Download downloads Nft by given txid, timestamp and signature.
func (service *DownloadNft) Download(m *pb.DownloadRequest, stream pb.DownloadNft_DownloadServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	// Create new task
	task := service.NewNftDownloadingTask()
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	// Call task download
	data, err := task.Download(ctx, m.GetTxid(), m.GetTimestamp(), m.GetSignature(), m.GetTtxid())
	if err != nil {
		return err
	}

	// Return restored image back to WalletNode
	remaining := len(data) % downloadImageBufferSize
	n := len(data) / downloadImageBufferSize

	for i := 0; i < n; i++ {
		res := &pb.DownloadReply{
			File: data[i*downloadImageBufferSize : (i+1)*downloadImageBufferSize],
		}
		if err := stream.Send(res); err != nil {
			return errors.Errorf("send image data: %w", err)
		}
	}
	if remaining > 0 {
		lastPackageIndex := n * downloadImageBufferSize
		res := &pb.DownloadReply{
			File: data[lastPackageIndex : lastPackageIndex+remaining],
		}
		if err := stream.Send(res); err != nil {
			return errors.Errorf("send image data: %w", err)
		}
	}

	return nil
}

// DownloadThumbnail returns thumbnail of given hash
func (service *DownloadNft) DownloadThumbnail(ctx context.Context, req *pb.DownloadThumbnailRequest) (*pb.DownloadThumbnailReply, error) {
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
	data, err := task.DownloadThumbnail(ctx, req.Key)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadThumbnailReply{
		Thumbnail: data,
	}, nil
}

// Desc returns a description of the service.
func (service *DownloadNft) Desc() *grpc.ServiceDesc {
	return &pb.DownloadNft_ServiceDesc
}

// NewDownloadNft returns a new DownloadNft instance.
func NewDownloadNft(service *nftdownload.NftDownloaderService) *DownloadNft {
	return &DownloadNft{
		DownloadNft: common.NewDownloadNft(service),
	}
}
