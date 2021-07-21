package grpc

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
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
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		IsPrimary: isPrimary,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("failed to send Session request: %w", err)
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
		return errors.Errorf("failed to receive Session response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")
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

	req := &pb.AcceptedNodesRequest{}
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")

	resp, err := service.client.AcceptedNodes(ctx, req)
	if err != nil {
		return nil, errors.Errorf("failed to request to accepted secondary nodes: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")

	var ids []string
	for _, peer := range resp.Peers {
		ids = append(ids, peer.NodeID)
	}
	return ids, nil
}

// ConnectTo implements node.ProcessUserdata.ConnectTo()
func (service *processUserdata) ConnectTo(ctx context.Context, nodeID, sessID string) error {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	req := &pb.ConnectToRequest{
		NodeID: nodeID,
		SessID: sessID,
	}
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")

	resp, err := service.client.ConnectTo(ctx, req)
	if err != nil {
		return errors.Errorf("failed to request to connect to primary node: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")

	return nil
}

// SendUserdata implements node.ProcessUserdata.SendUserdata()
func (service *processUserdata) SendUserdata(ctx context.Context, request *userdata.UserdataProcessRequestSigned)
		(result *userdata.UserdataProcessResult, err error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.SendUserdata(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to open stream: %w", err)
	}
	defer stream.CloseSend()

	for {
		// Generate protobuf request reqProto
		reqProto := &pb.UserdataRequest{}

		reqProto.Realname = request.Userdata.Realname,
		reqProto.Facebook_link = request.Userdata.FacebookLink,
		reqProto.Twitter_link = request.Userdata.TwitterLink,
		reqProto.Native_currency = request.Userdata.NativeCurrency,
		reqProto.Location = request.Userdata.Location
		reqProto.Primary_language = request.Userdata.PrimaryLanguage
		reqProto.Categories = request.Userdata.Categories
		reqProto.Biography = request.Userdata.Biography
		reqProto.AvatarImage = &pb.UserdataRequest_UserImageUpload {}
	
		if request.Userdata.AvatarImage.Content != nil && len(request.Userdata.AvatarImage.Content) > 0 {
			image.AvatarImage.Content = make ([]byte, len(request.Userdata.AvatarImage.Content))
			copy(reqProto.AvatarImage.Content,request.Userdata.AvatarImage.Content)
		}
		reqProto.AvatarImage.Filename := request.Userdata.AvatarImage.Filename

		if request.Userdata.CoverPhoto.Content != nil && len(request.Userdata.CoverPhoto.Content) > 0 {
			image.CoverPhoto.Content = make ([]byte, len(request.Userdata.CoverPhoto.Content))
			copy(reqProto.CoverPhoto.Content,request.Userdata.CoverPhoto.Content)
		}
		reqProto.CoverPhoto.Filename := request.Userdata.CoverPhoto.Filename
		
		reqProto.ArtistPastelID = request.Userdata.ArtistPastelID 
		reqProto.Timestamp = request.Userdata.Timestamp  
		reqProto.PreviousBlockHash = request.Userdata.PreviousBlockHash
		reqProto.UserdataHash = request.Userdata.UserdataHash
		reqProto.Signature = request.Userdata.Signature

		// Send the request to the protobuf stream
		if err := stream.Send(reqProto); err != nil {
			return nil, errors.Errorf("failed to userdata to protobuf stream: %w", err).WithField("reqID", service.conn.id)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, errors.Errorf("failed to receive send image response: %w", err)
	}
	log.WithContext(ctx).WithField("fingerprintLenght", len(resp.Fingerprint)).Debugf("ProbeImage response")

	
	// Convert protobuf response to UserdataProcessResult then return it
	result := &userdata.UserdataProcessResult {
		ResponseCode: 		resp.UserdataReply.ResponseCode,
		Detail:				resp.UserdataReply.Detail,
		Realname:			resp.UserdataReply.Realname,
		FacebookLink:		resp.UserdataReply.FacebookLink,
		TwitterLink:		resp.UserdataReply.TwitterLink,
		NativeCurrency:		resp.UserdataReply.NativeCurrency,
		Location:			resp.UserdataReply.Location,
		PrimaryLanguage:	resp.UserdataReply.PrimaryLanguage,
		Categories:			resp.UserdataReply.Categories,
		AvatarImage:		resp.UserdataReply.AvatarImage,
		CoverPhoto:			resp.UserdataReply.CoverPhoto,
	}
	
	return result, nil
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
