package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/walletnode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type processUserdata struct {
	conn   *clientConn
	client pb.ProcessUserdataClient

	sessID string
}

func (service *processUserdata) SessID() string {
	return service.sessID
}

// Session implements node.ProcessUserdata.Session()
func (service *processUserdata) Session(ctx context.Context, isPrimary bool) error {
	ctx = service.contextWithLogPrefix(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Health stream: %w", err)
	}

	req := &pb.MDLSessionRequest{
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

// AcceptedNodes implements node.ProcessUserdata.AcceptedNodes()
func (service *processUserdata) AcceptedNodes(ctx context.Context) (pastelIDs []string, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.MDLAcceptedNodesRequest{}
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

// ConnectTo implements node.ProcessUserdata.ConnectTo()
func (service *processUserdata) ConnectTo(ctx context.Context, primaryNode types.MeshedSuperNode) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.MDLConnectToRequest{
		NodeID: primaryNode.NodeID,
		SessID: primaryNode.SessID,
	}
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return errors.Errorf("request to connect to primary node: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")

	return nil
}

// SendUserdata implements node.ProcessUserdata.SendUserdata()
func (service *processUserdata) SendUserdata(ctx context.Context, request *userdata.ProcessRequestSigned) (result *userdata.ProcessResult, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	if request == nil {
		return nil, errors.Errorf("input nil request")
	}

	if request.Userdata == nil {
		return nil, errors.Errorf("input nil userdata")
	}

	// Generate protobuf request reqProto
	reqProto := &pb.UserdataRequest{
		RealName:        request.Userdata.RealName,
		FacebookLink:    request.Userdata.FacebookLink,
		TwitterLink:     request.Userdata.TwitterLink,
		NativeCurrency:  request.Userdata.NativeCurrency,
		Location:        request.Userdata.Location,
		PrimaryLanguage: request.Userdata.PrimaryLanguage,
		Categories:      request.Userdata.Categories,
		Biography:       request.Userdata.Biography,
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  request.Userdata.AvatarImage.Content,
			Filename: request.Userdata.AvatarImage.Filename,
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  request.Userdata.CoverPhoto.Content,
			Filename: request.Userdata.CoverPhoto.Filename,
		},
		ArtistPastelID:    request.Userdata.ArtistPastelID,
		Timestamp:         request.Userdata.Timestamp,
		PreviousBlockHash: request.Userdata.PreviousBlockHash,
		UserdataHash:      request.UserdataHash,
		Signature:         request.Signature,
	}

	resp, err := service.client.SendUserdata(ctx, reqProto)
	if err != nil {
		return nil, errors.Errorf("send user data: %w", err)
	}

	// Convert protobuf response to UserdataProcessResult then return it
	result = &userdata.ProcessResult{
		ResponseCode:    resp.ResponseCode,
		Detail:          resp.Detail,
		RealName:        resp.RealName,
		FacebookLink:    resp.FacebookLink,
		TwitterLink:     resp.TwitterLink,
		NativeCurrency:  resp.NativeCurrency,
		Location:        resp.Location,
		PrimaryLanguage: resp.PrimaryLanguage,
		Categories:      resp.Categories,
		AvatarImage:     resp.AvatarImage,
		CoverPhoto:      resp.CoverPhoto,
	}

	return result, nil
}

// ReceiveUserdata implements node.ProcessUserdata.ReceiveUserdata()
func (service *processUserdata) ReceiveUserdata(ctx context.Context, userpastelid string) (result *userdata.ProcessRequest, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	reqProto := &pb.RetrieveRequest{
		Userpastelid: userpastelid,
	}
	resp, err := service.client.ReceiveUserdata(ctx, reqProto)
	if err != nil {
		return nil, errors.Errorf("receive data: %w", err)
	}

	var avatarImage userdata.UserImageUpload
	var coverPhoto userdata.UserImageUpload
	if resp.AvatarImage != nil {
		avatarImage = userdata.UserImageUpload{
			Content:  resp.AvatarImage.Content,
			Filename: resp.AvatarImage.Filename,
		}
	}
	if resp.CoverPhoto != nil {
		coverPhoto = userdata.UserImageUpload{
			Content:  resp.CoverPhoto.Content,
			Filename: resp.CoverPhoto.Filename,
		}
	}

	// Convert protobuf request to UserdataProcessRequest
	response := userdata.ProcessRequest{
		RealName:          resp.RealName,
		FacebookLink:      resp.FacebookLink,
		TwitterLink:       resp.TwitterLink,
		NativeCurrency:    resp.NativeCurrency,
		Location:          resp.Location,
		PrimaryLanguage:   resp.PrimaryLanguage,
		Categories:        resp.Categories,
		Biography:         resp.Biography,
		AvatarImage:       avatarImage,
		CoverPhoto:        coverPhoto,
		ArtistPastelID:    resp.ArtistPastelID,
		Timestamp:         resp.Timestamp,
		PreviousBlockHash: resp.PreviousBlockHash,
	}

	return &response, nil
}

func (service *processUserdata) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *processUserdata) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func newProcessUserdata(conn *clientConn) node.ProcessUserdata {
	return &processUserdata{
		conn:   conn,
		client: pb.NewProcessUserdataClient(conn),
	}
}
