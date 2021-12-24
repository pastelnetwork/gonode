package grpc

import (
	"context"

	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/artwork"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	senseUploadImageBufferSize = 32 * 1024
)

type registerSense struct {
	conn   *clientConn
	client pb.RegisterSenseClient

	sessID string
}

func (service *registerSense) SessID() string {
	return service.sessID
}

// Session implements node.RegisterArtwork.Session()
func (service *registerSense) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Session stream: %w", err)
	}

	req := &pb.SenseSessionRequest{
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

// AcceptedNodes implements node.RegisterArtwork.AcceptedNodes()
func (service *registerSense) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.SenseAcceptedNodesRequest{}
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

// ConnectTo implements node.RegisterArtwork.ConnectTo()
func (service *registerSense) ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.SenseConnectToRequest{
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
func (service *registerSense) MeshNodes(ctx context.Context, meshedNodes []types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.SenseMeshNodesRequest{
		Nodes: []*pb.SenseMeshNodesRequest_Node{},
	}

	for _, node := range meshedNodes {
		request.Nodes = append(request.Nodes, &pb.SenseMeshNodesRequest_Node{
			SessID: node.SessID,
			NodeID: node.NodeID,
		})
	}

	_, err := service.client.MeshNodes(ctx, request)

	return err
}

// SendRegMetadata send metadata of registration to SNs for next steps
func (service *registerSense) SendRegMetadata(ctx context.Context, regMetadata *types.NftRegMetadata) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.SenseSendRegMetadataRequest{
		BlockHash:       regMetadata.BlockHash,
		CreatorPastelID: regMetadata.CreatorPastelID,
	}

	_, err := service.client.SendRegMetadata(ctx, request)
	return err
}

// ProbeImage implements node.RegisterArtwork.ProbeImage()
func (service *registerSense) ProbeImage(ctx context.Context, image *artwork.File) ([]byte, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.ProbeImage(ctx)
	if err != nil {
		return nil, errors.Errorf("open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := image.Open()
	if err != nil {
		return nil, errors.Errorf("open file %q: %w", file.Name(), err)
	}
	defer file.Close()

	buffer := make([]byte, senseUploadImageBufferSize)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		req := &pb.SenseProbeImageRequest{
			Payload: buffer[:n],
		}
		if err := stream.Send(req); err != nil {
			return nil, errors.Errorf("send image data: %w", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("receive image response: %w", err)
	}

	return resp.CompressedSignedDDAndFingerprints, nil
}

func (service *registerSense) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerSense) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// SendSignedTicket
func (service *registerSense) SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, ddFp []byte) (int64, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := pb.SenseSendSignedNFTTicketRequest{
		NftTicket:        ticket,
		CreatorSignature: signature,
		DdFpFiles:        ddFp,
	}

	rsp, err := service.client.SendSignedNFTTicket(ctx, &req)
	if err != nil {
		return -1, err
	}

	return rsp.RegistrationFee, nil
}

func (service *registerSense) SendPreBurntFeeTxid(ctx context.Context, txid string) (string, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	log.WithContext(ctx).Debug("send burned txid to super node")
	req := pb.SenseSendPreBurntFeeTxidRequest{
		Txid: txid,
	}

	rsp, err := service.client.SendPreBurntFeeTxid(ctx, &req)
	if err != nil {
		return "", errors.Errorf("send burned txid to super node: %w", err)
	}

	// TODO: response from sending preburned TxId should be the TxId of RegActTicket
	return rsp.NFTRegTxid, nil
}

func newRegisterSense(conn *clientConn) node.RegisterSense {
	return &registerSense{
		conn:   conn,
		client: pb.NewRegisterSenseClient(conn),
	}
}
