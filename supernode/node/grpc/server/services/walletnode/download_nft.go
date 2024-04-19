package walletnode

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/download"
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
	data, err := task.Download(ctx, m.GetTxid(), m.GetTimestamp(), m.GetSignature(), m.GetTtxid(), m.GetTtype(), m.GetSendHash())
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
func (service *DownloadNft) DownloadDDAndFingerprints(ctx context.Context, req *pb.DownloadDDAndFingerprintsRequest) (*pb.DownloadDDAndFingerprintsReply, error) {
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

// GetTopMNs implements walletnode.downloadNFT.GetTopMNs()
func (service *DownloadNft) GetTopMNs(ctx context.Context, req *pb.GetTopMNsRequest) (*pb.GetTopMNsReply, error) {
	log.WithContext(ctx).Debug("request for mn-top list has been received")

	var mnTopList pastel.MasterNodes
	var err error
	if req.Block == 0 {
		mnTopList, err = service.PastelClient.MasterNodesTop(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving mn-top list on the SN")
		}
	} else {
		mnTopList, err = service.PastelClient.MasterNodesTopN(ctx, int(req.Block))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving mn-top list on the SN")
		}
	}

	var mnList []string
	for _, mn := range mnTopList {
		mnList = append(mnList, mn.ExtAddress)
	}

	blnc, err := service.PastelClient.ZGetTotalBalance(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving balance on the SN")
		blnc = &pastel.GetTotalBalanceResponse{}
	}

	resp := &pb.GetTopMNsReply{
		MnTopList:   mnList,
		CurrBalance: int64(blnc.Total),
	}

	log.WithContext(ctx).WithField("mn-list", mnList).Debug("top mn-list has been returned")
	return resp, nil
}

// Desc returns a description of the service.
func (service *DownloadNft) Desc() *grpc.ServiceDesc {
	return &pb.DownloadNft_ServiceDesc
}

// NewDownloadNft returns a new DownloadNft instance.
func NewDownloadNft(service *download.NftDownloaderService) *DownloadNft {
	return &DownloadNft{
		DownloadNft: common.NewDownloadNft(service),
	}
}
