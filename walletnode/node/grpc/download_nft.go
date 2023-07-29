package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/types"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type downloadNft struct {
	conn   *clientConn
	client pb.DownloadNftClient
}

func (service *downloadNft) Download(ctx context.Context, txid, timestamp, signature, ttxid, ttype string) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)

	in := &pb.DownloadRequest{
		Txid:      txid,
		Timestamp: timestamp,
		Signature: signature,
		Ttxid:     ttxid,
		Ttype:     ttype,
	}

	var stream pb.DownloadNft_DownloadClient
	stream, err = service.client.Download(ctx, in)
	if err != nil {
		err = errors.Errorf("open stream: %w", err)
		return
	}
	defer stream.CloseSend()

	// Receive file
	for {
		var resp *pb.DownloadReply
		resp, err = stream.Recv()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		file = append(file, resp.File...)
	}

	return
}

func (service *downloadNft) DownloadThumbnail(ctx context.Context, txid string, numNails int) (files map[int][]byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.DownloadThumbnailRequest{
		Txid:     txid,
		Numnails: int32(numNails),
	}
	log.WithContext(ctx).Println("Sending sn download thumbnail request")
	res, err := service.client.DownloadThumbnail(ctx, in)

	if err != nil {
		return nil, err
	}

	if res.Thumbnailone == nil {
		return nil, errors.New("nil thumbnail")
	}

	if res.Thumbnailtwo == nil && numNails > 1 {
		return nil, errors.New("nil thumbnail2")
	}
	rMap := make(map[int][]byte)
	rMap[0] = res.Thumbnailone
	rMap[1] = res.Thumbnailtwo
	return rMap, nil
}

func (service *downloadNft) DownloadDDAndFingerprints(ctx context.Context, txid string) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.DownloadDDAndFingerprintsRequest{
		Txid: txid,
	}

	res, err := service.client.DownloadDDAndFingerprints(ctx, in)
	if err != nil {
		return nil, err
	}

	if res.File == nil {
		return nil, errors.New("nil dd and fingerprints file")
	}

	return res.File, nil
}

func (service *downloadNft) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// NON IMPLEMETED--->
func (service *downloadNft) MeshNodes(_ context.Context, _ []types.MeshedSuperNode) error {
	return nil
}
func (service *downloadNft) SessID() string {
	return ""
}
func (service *downloadNft) Session(_ context.Context, _ bool) error {
	return nil
}
func (service *downloadNft) AcceptedNodes(_ context.Context) (pastelIDs []string, err error) {
	return nil, nil
}
func (service *downloadNft) ConnectTo(_ context.Context, _ types.MeshedSuperNode) error {
	return nil
}

func (service *downloadNft) Close() error {
	return service.conn.Close()
}

// GetDupeDetectionDBHash implements node.RegisterNft.GetDupeDetectionDBHash
func (service *downloadNft) GetDupeDetectionDBHash(_ context.Context) (hash string, err error) {
	return "not implemented", nil
}

// GetTopMNs implements node.downloadNFT.GetTopMNs
func (service *downloadNft) GetTopMNs(ctx context.Context) (mnList *pb.GetTopMNsReply, err error) {
	ctx = service.contextWithLogPrefix(ctx)

	req := &pb.GetTopMNsRequest{}
	resp, err := service.client.GetTopMNs(ctx, req)
	if err != nil {
		return nil, errors.Errorf("WN request to SN for mn-top list: %w", err)
	}

	return resp, nil
}

// GetDDServerStats implements node.RegisterNft.GetDDServerStats
func (service *downloadNft) GetDDServerStats(_ context.Context) (stats *pb.DDServerStatsReply, err error) {
	return stats, nil
}

///<---

func newDownloadNft(conn *clientConn) node.DownloadNftInterface {
	return &downloadNft{
		conn:   conn,
		client: pb.NewDownloadNftClient(conn),
	}
}
