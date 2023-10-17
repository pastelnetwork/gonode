package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type registerCollection struct {
	conn   *clientConn
	client pb.RegisterCollectionClient
	sessID string
}

func (service *registerCollection) SessID() string {
	return service.sessID
}

// Session implements node.RegisterCollection.Session()
func (service *registerCollection) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Session stream: %w", err)
	}

	req := &pb.SessionRequest{
		IsPrimary: isPrimary,
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("send Session request: %w", err)
	}

	resp, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		switch status.Code(err) {
		case codes.Canceled, codes.Unavailable:
			return nil
		}
		return errors.Errorf("receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")
	service.sessID = resp.SessID

	go func() {
		defer service.conn.Close()
		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
		}
	}()

	return nil
}

// AcceptedNodes implements node.RegisterCollection.AcceptedNodes()
func (service *registerCollection) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeID)
	}
	return ids, nil
}

// ConnectTo implements node.RegisterCollection.ConnectTo()
func (service *registerCollection) ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.ConnectToRequest{
		NodeID: primaryNode.NodeID,
		SessID: primaryNode.SessID,
	}
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return err
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")

	return nil
}

// MeshNodes informs SNs which SNs are connected to do NFT request
func (service *registerCollection) MeshNodes(ctx context.Context, meshedNodes []types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.MeshNodesRequest{
		Nodes: []*pb.MeshNodesRequest_Node{},
	}

	for _, node := range meshedNodes {
		request.Nodes = append(request.Nodes, &pb.MeshNodesRequest_Node{
			SessID: node.SessID,
			NodeID: node.NodeID,
		})
	}

	_, err := service.client.MeshNodes(ctx, request)

	return err
}

// SendTicketForSignature sends the collection ticket to be signed by other SNs
func (service *registerCollection) SendTicketForSignature(ctx context.Context, ticket []byte, signature []byte, burnTxID string) (string, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := pb.SendCollectionTicketForSignatureRequest{
		CollectionTicket: ticket,
		CreatorSignature: signature,
		BurnTxid:         burnTxID,
	}

	rsp, err := service.client.SendCollectionTicketForSignature(ctx, &req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("send signed ticket failed")

		return "", err
	}

	return rsp.CollectionRegTxid, nil
}

// GetDupeDetectionDBHash implements node.RegisterCollection.GetDupeDetectionDBHash
func (service *registerCollection) GetDupeDetectionDBHash(_ context.Context) (hash string, err error) {
	return "not implemented", nil
}

// GetDDServerStats implements node.RegisterNft.GetDDServerStats
func (service *registerCollection) GetDDServerStats(_ context.Context) (stats *pb.DDServerStatsReply, err error) {
	//not implemented
	return stats, nil
}

// GetTopMNs implements node.RegisterCollection.GetTopMNs
func (service *registerCollection) GetTopMNs(ctx context.Context, block int) (mnList *pb.GetTopMNsReply, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.GetTopMNsRequest{
		Block: int64(block),
	}
	resp, err := service.client.GetTopMNs(ctx, req)
	if err != nil {
		return nil, errors.Errorf("WN request to SN for mn-top list: %w", err)
	}

	return resp, nil
}

func (service *registerCollection) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func (service *registerCollection) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func newRegisterCollection(conn *clientConn) node.RegisterCollectionInterface {
	return &registerCollection{
		client: pb.NewRegisterCollectionClient(conn),
		conn:   conn,
	}
}
