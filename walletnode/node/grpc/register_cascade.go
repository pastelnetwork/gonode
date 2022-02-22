package grpc

import (
	"context"
	"encoding/base64"

	"github.com/pastelnetwork/gonode/common/storage/files"
	rqnode "github.com/pastelnetwork/gonode/raptorq/node"
	"golang.org/x/crypto/sha3"

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

const (
	cascadeUploadAssetBufferSize = 32 * 1024
)

type registerCascade struct {
	conn   *clientConn
	client pb.RegisterCascadeClient

	sessID string
}

func (service *registerCascade) SessID() string {
	return service.sessID
}

// Session implements node.RegisterNft.Session()
func (service *registerCascade) Session(ctx context.Context, isPrimary bool) error {
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

// AcceptedNodes implements node.RegisterNft.AcceptedNodes()
func (service *registerCascade) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
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

// ConnectTo implements node.RegisterNft.ConnectTo()
func (service *registerCascade) ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error {
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
func (service *registerCascade) MeshNodes(ctx context.Context, meshedNodes []types.MeshedSuperNode) error {
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

// SendRegMetadata send metadata of registration to SNs for next steps
func (service *registerCascade) SendRegMetadata(ctx context.Context, regMetadata *types.ActionRegMetadata) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.SendRegMetadataRequest{
		BlockHash:       regMetadata.BlockHash,
		CreatorPastelID: regMetadata.CreatorPastelID,
		BurnTxid:        regMetadata.BurnTxID,
	}

	_, err := service.client.SendRegMetadata(ctx, request)
	return err
}

// SendActionAct send action act to SNs for next steps
func (service *registerCascade) SendActionAct(ctx context.Context, actionRegTxid string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	request := &pb.SendActionActRequest{
		ActionRegTxid: actionRegTxid,
	}

	_, err := service.client.SendActionAct(ctx, request)
	return err
}

func (service *registerCascade) UploadAsset(ctx context.Context, asset *files.File) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	log.WithContext(ctx).Debug("Start upload asset to node")
	stream, err := service.client.UploadAsset(ctx)
	if err != nil {
		return errors.Errorf("open stream: %w", err)
	}
	defer stream.CloseSend()

	file, err := asset.Open()
	if err != nil {
		return errors.Errorf("open file %q: %w", file.Name(), err)
	}
	defer file.Close()

	buffer := make([]byte, cascadeUploadAssetBufferSize)
	lastPiece := false
	payloadSize := 0
	for {
		n, err := file.Read(buffer)
		payloadSize += n
		if err != nil && err == io.EOF {
			log.WithContext(ctx).WithField("Filename", file.Name()).Debug("EOF")
			lastPiece = true
			break
		} else if err != nil {
			return errors.Errorf("read file: %w", err)
		}

		req := &pb.UploadAssetRequest{
			Payload: &pb.UploadAssetRequest_AssetPiece{
				AssetPiece: buffer[:n],
			},
		}
		if err := stream.Send(req); err != nil {
			return errors.Errorf("send image data: %w", err)
		}
	}
	log.WithContext(ctx).Debugf("Encoded Asset Size :%d\n", payloadSize)
	if !lastPiece {
		return errors.Errorf("read all asset data failed")
	}

	file.Seek(0, io.SeekStart)
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, file); err != nil {
		return errors.Errorf("compute asset hash:%w", err)
	}
	hash := hasher.Sum(nil)
	log.WithContext(ctx).WithField("Filename", file.Name()).Debugf("hash: %s", base64.URLEncoding.EncodeToString(hash))

	metaData := &pb.UploadAssetRequest{
		Payload: &pb.UploadAssetRequest_MetaData_{
			MetaData: &pb.UploadAssetRequest_MetaData{
				Hash:   hash[:],
				Size:   int64(payloadSize),
				Format: asset.Format().String(),
			},
		},
	}

	if err := stream.Send(metaData); err != nil {
		return errors.Errorf("send asset metadata: %w", err)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return errors.Errorf("receive upload asset response: %w", err)
	}
	return nil
}

func (service *registerCascade) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *registerCascade) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

// SendSignedTicket
func (service *registerCascade) SendSignedTicket(ctx context.Context, ticket []byte, signature []byte, rqIdsFile []byte, encoderParams rqnode.EncoderParameters) (string, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := pb.SendSignedCascadeTicketRequest{
		ActionTicket:     ticket,
		CreatorSignature: signature,
		RqFiles:          rqIdsFile,
		EncodeParameters: &pb.EncoderParameters{Oti: encoderParams.Oti},
	}

	rsp, err := service.client.SendSignedActionTicket(ctx, &req)
	if err != nil {
		return "", err
	}

	return rsp.ActionRegTxid, nil
}

func newRegisterCascade(conn *clientConn) node.RegisterCascadeInterface {
	return &registerCascade{
		conn:   conn,
		client: pb.NewRegisterCascadeClient(conn),
	}
}
