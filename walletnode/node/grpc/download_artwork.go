package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
)

type downloadArtwork struct {
	conn   *clientConn
	client pb.DownloadArtworkClient
}

func (service *downloadArtwork) Download(ctx context.Context, txid, timestamp, signature, ttxid string) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)

	in := &pb.DownloadRequest{
		Txid:      txid,
		Timestamp: timestamp,
		Signature: signature,
		Ttxid:     ttxid,
	}

	var stream pb.DownloadArtwork_DownloadClient
	stream, err = service.client.Download(ctx, in)
	if err != nil {
		err = errors.Errorf("failed to open stream: %w", err)
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

func (service *downloadArtwork) DownloadThumbnail(ctx context.Context, key []byte) (file []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.DownloadThumbnailRequest{
		Key: key,
	}

	res, err := service.client.DownloadThumbnail(ctx, in)

	return res.Thumbnail, err
}

func (service *downloadArtwork) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newDownloadArtwork(conn *clientConn) node.DownloadArtwork {
	return &downloadArtwork{
		conn:   conn,
		client: pb.NewDownloadArtworkClient(conn),
	}
}
