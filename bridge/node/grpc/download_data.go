package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/types"

	"github.com/pastelnetwork/gonode/bridge/node"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/bridge"
)

type downloadData struct {
	conn   *clientConn
	client pb.DownloadDataClient
}

// DownloadThumbnail downloads thumbnail
func (service *downloadData) DownloadThumbnail(ctx context.Context, txid string, numNails int) (files map[int][]byte, err error) {
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

// DownloadDDAndFingerprints downloads dd & fingerprints
func (service *downloadData) DownloadDDAndFingerprints(ctx context.Context, txid string) (file []byte, err error) {
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

// Health checks bridge health
func (service *downloadData) Health(ctx context.Context) (err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.PingRequest{}

	_, err = service.client.Ping(ctx, in)
	return err
}

func (service *downloadData) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// MeshNodes NON IMPLEMETED--->
func (service *downloadData) MeshNodes(_ context.Context, _ []types.MeshedSuperNode) error {
	return nil
}

// SessID ...
func (service *downloadData) SessID() string {
	return ""
}

// Session ...
func (service *downloadData) Session(_ context.Context, _ bool) error {
	return nil
}

// AcceptedNodes ...
func (service *downloadData) AcceptedNodes(_ context.Context) (pastelIDs []string, err error) {
	return nil, nil
}

//  ConnectTo ...
func (service *downloadData) ConnectTo(_ context.Context, _ types.MeshedSuperNode) error {
	return nil
}

func newDownloadData(conn *clientConn) node.DownloadDataInterface {
	return &downloadData{
		conn:   conn,
		client: pb.NewDownloadDataClient(conn),
	}
}
