package client

import (
	"context"
	"fmt"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	"github.com/pastelnetwork/gonode/proto"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
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
		return errors.Errorf("open Health stream: %w", err)
	}

	req := &pb.MDLSessionRequest{
		NodeID: nodeID,
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
		UserdataHash:       dataSigned.UserdataHash,
		UserdataResultHash: dataSigned.UserdataResultHash,
		HashSignature:      dataSigned.HashSignature,
		NodeID:             dataSigned.NodeID,
	})

	if err != nil {
		return userdata.SuperNodeReply{}, err
	}

	return userdata.SuperNodeReply{
		ResponseCode: resp.ResponseCode,
		Detail:       resp.Detail,
	}, nil
}

func (service *processUserdata) SendUserdataToLeader(ctx context.Context, finalUserdata userdata.ProcessRequestSigned) (userdata.SuperNodeReply, error) {
	ctx = service.contextWithLogPrefix(ctx)
	ctx = service.contextWithMDSessID(ctx)

	if finalUserdata.Userdata == nil {
		return userdata.SuperNodeReply{}, errors.Errorf("input nil userdata")
	}

	// Generate protobuf request reqProto
	reqProto := &pb.UserdataRequest{
		RealName:        finalUserdata.Userdata.RealName,
		FacebookLink:    finalUserdata.Userdata.FacebookLink,
		TwitterLink:     finalUserdata.Userdata.TwitterLink,
		NativeCurrency:  finalUserdata.Userdata.NativeCurrency,
		Location:        finalUserdata.Userdata.Location,
		PrimaryLanguage: finalUserdata.Userdata.PrimaryLanguage,
		Categories:      finalUserdata.Userdata.Categories,
		Biography:       finalUserdata.Userdata.Biography,
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  finalUserdata.Userdata.AvatarImage.Content,
			Filename: finalUserdata.Userdata.AvatarImage.Filename,
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  finalUserdata.Userdata.CoverPhoto.Content,
			Filename: finalUserdata.Userdata.CoverPhoto.Filename,
		},
		ArtistPastelID:    finalUserdata.Userdata.ArtistPastelID,
		Timestamp:         finalUserdata.Userdata.Timestamp,
		PreviousBlockHash: finalUserdata.Userdata.PreviousBlockHash,
		UserdataHash:      finalUserdata.UserdataHash,
		Signature:         finalUserdata.Signature,
	}

	// Send the data
	resp, err := service.client.SendUserdataToLeader(ctx, reqProto)

	if err != nil {
		return userdata.SuperNodeReply{}, err
	}

	return userdata.SuperNodeReply{
		ResponseCode: resp.ResponseCode,
		Detail:       resp.Detail,
	}, nil
}

func newProcessUserdata(conn *clientConn) node.ProcessUserdata {
	return &processUserdata{
		conn:   conn,
		client: pb.NewProcessUserdataClient(conn),
	}
}
