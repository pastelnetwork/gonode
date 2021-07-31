package client

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/network/proto"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/supernode"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/metadb/network/supernode/node"
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
func (service *processUserdata) Session(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("failed to open Health stream: %w", err)
	}

	req := &pb.MDLSessionRequest{
		NodeID: nodeID,
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

func (service *processUserdata) contextWithMDSessID(ctx context.Context) context.Context {
	md := metadata.Pairs(proto.MetadataKeySessID, service.sessID)
	return metadata.NewOutgoingContext(ctx, md)
}

func (service *processUserdata) contextWithLogPrefix(ctx context.Context) context.Context {
	return log.ContextWithPrefix(ctx, fmt.Sprintf("%s-%s", logPrefix, service.conn.id))
}

func (service *processUserdata) SendUserdataToPrimary(ctx context.Context, dataSigned userdata.SuperNodeRequest) (userdata.SuperNodeReply, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)
	resp, err := service.client.SendUserdataToPrimary(ctx, &pb.SuperNodeRequest{
		UserdataHash			: dataSigned.UserdataHash,
		UserdataResultHash		: dataSigned.UserdataResultHash,
		HashSignature			: dataSigned.HashSignature,
		NodeID					: dataSigned.NodeID,
	})

	return userdata.SuperNodeReply {
		ResponseCode	: resp.ResponseCode,
		Detail			: resp.Detail,
	}, err
}

func (service *processUserdata) SendUserdataToLeader(ctx context.Context, finalUserdata userdata.UserdataProcessRequestSigned) (userdata.SuperNodeReply, error){
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	// Generate protobuf request reqProto
	reqProto := &pb.UserdataRequest{}

	reqProto.Realname = finalUserdata.Userdata.Realname
	reqProto.FacebookLink = finalUserdata.Userdata.FacebookLink
	reqProto.TwitterLink = finalUserdata.Userdata.TwitterLink
	reqProto.NativeCurrency = finalUserdata.Userdata.NativeCurrency
	reqProto.Location = finalUserdata.Userdata.Location
	reqProto.PrimaryLanguage = finalUserdata.Userdata.PrimaryLanguage
	reqProto.Categories = finalUserdata.Userdata.Categories
	reqProto.Biography = finalUserdata.Userdata.Biography
	reqProto.AvatarImage = &pb.UserdataRequest_UserImageUpload {}

	if finalUserdata.Userdata.AvatarImage.Content != nil && len(finalUserdata.Userdata.AvatarImage.Content) > 0 {
		reqProto.AvatarImage.Content = make ([]byte, len(finalUserdata.Userdata.AvatarImage.Content))
		copy(reqProto.AvatarImage.Content,finalUserdata.Userdata.AvatarImage.Content)
	}
	reqProto.AvatarImage.Filename = finalUserdata.Userdata.AvatarImage.Filename

	reqProto.CoverPhoto = &pb.UserdataRequest_UserImageUpload {}
	if finalUserdata.Userdata.CoverPhoto.Content != nil && len(finalUserdata.Userdata.CoverPhoto.Content) > 0 {
		reqProto.CoverPhoto.Content = make ([]byte, len(finalUserdata.Userdata.CoverPhoto.Content))
		copy(reqProto.CoverPhoto.Content,finalUserdata.Userdata.CoverPhoto.Content)
	}
	reqProto.CoverPhoto.Filename = finalUserdata.Userdata.CoverPhoto.Filename
	
	reqProto.ArtistPastelID = finalUserdata.Userdata.ArtistPastelID 
	reqProto.Timestamp = finalUserdata.Userdata.Timestamp  
	reqProto.PreviousBlockHash = finalUserdata.Userdata.PreviousBlockHash
	reqProto.UserdataHash = finalUserdata.UserdataHash
	reqProto.Signature = finalUserdata.Signature

	// Send the data
	resp, err := service.client.SendUserdataToLeader(ctx, reqProto)

	return userdata.SuperNodeReply {
		ResponseCode	: resp.ResponseCode,
		Detail			: resp.Detail,
	}, err
}

func newProcessUserdata(conn *clientConn) node.ProcessUserdata {
	return &processUserdata{
		conn:   conn,
		client: pb.NewProcessUserdataClient(conn),
	}
}
