package grpc

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/types"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/hermes/service/node"
	pb "github.com/pastelnetwork/gonode/proto/hermes"
)

type hermesP2P struct {
	conn   *clientConn
	client pb.HermesP2PClient
}

// DownloadThumbnail downloads thumbnail
func (service *hermesP2P) Retrieve(ctx context.Context, txid string) (data []byte, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.RetrieveRequest{
		Key: txid,
	}

	log.WithContext(ctx).Println("Sending sn retrieve request")
	res, err := service.client.Retrieve(ctx, in)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// DownloadDDAndFingerprints downloads dd & fingerprint
func (service *hermesP2P) Delete(ctx context.Context, txid string) (err error) {
	ctx = service.contextWithLogPrefix(ctx)
	in := &pb.DeleteRequest{
		Key: txid,
	}

	_, err = service.client.Delete(ctx, in)
	if err != nil {
		return fmt.Errorf("grpc delete err: %w", err)
	}

	return nil
}

func (service *hermesP2P) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// MeshNodes NON IMPLEMETED--->
func (service *hermesP2P) MeshNodes(_ context.Context, _ []types.MeshedSuperNode) error {
	return nil
}

// SessID ...
func (service *hermesP2P) SessID() string {
	return ""
}

// Session ...
func (service *hermesP2P) Session(_ context.Context, _ bool) error {
	return nil
}

// AcceptedNodes ...
func (service *hermesP2P) AcceptedNodes(_ context.Context) (pastelIDs []string, err error) {
	return nil, nil
}

//  ConnectTo ...
func (service *hermesP2P) ConnectTo(_ context.Context, _ types.MeshedSuperNode) error {
	return nil
}

func newHermesP2P(conn *clientConn) node.HermesP2PInterface {
	return &hermesP2P{
		conn:   conn,
		client: pb.NewHermesP2PClient(conn),
	}
}
