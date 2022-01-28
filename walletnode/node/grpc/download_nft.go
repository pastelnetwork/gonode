package grpc

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/types"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type downloadNft struct {
	conn   *clientConn
	client pb.DownloadNftClient
}

func (service *downloadNft) Download(ctx context.Context, txid, timestamp, signature, ttxid string) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)

	in := &pb.DownloadRequest{
		Txid:      txid,
		Timestamp: timestamp,
		Signature: signature,
		Ttxid:     ttxid,
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

func (service *downloadNft) DownloadThumbnail(ctx context.Context, key []byte) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.DownloadThumbnailRequest{
		Key: key,
	}

	res, err := service.client.DownloadThumbnail(ctx, in)

	if err != nil {
		return nil, err
	}

	if res.Thumbnail == nil {
		return nil, errors.New("nil thumbnail")
	}

	return res.Thumbnail, nil
}

func (service *downloadNft) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

//NON IMPLEMETED--->
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

///<---

func newDownloadArtwork(conn *clientConn) node.DownloadNftInterface {
	return &downloadNft{
		conn:   conn,
		client: pb.NewDownloadNftClient(conn),
	}
}
