package walletnode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/service/userdata"
	pb "github.com/pastelnetwork/gonode/metadb/network/proto/walletnode"
	"github.com/pastelnetwork/gonode/metadb/network/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/userdataprocess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ProcessUserdata represents grpc service for processing userdata.
type ProcessUserdata struct {
	pb.UnimplementedProcessUserdataServer

	*common.ProcessUserdata
}

// Session implements walletnode.ProcessUserdataServer.Session()
func (service *ProcessUserdata) Session(stream pb.ProcessUserdata_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *userdataprocess.Task

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("failed to receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if err := task.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("failed to send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("Session response")

	for {
		if _, err := stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			switch status.Code(err) {
			case codes.Canceled, codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}

// AcceptedNodes implements walletnode.ProcessUserdataServer.AcceptedNodes()
func (service *ProcessUserdata) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("AcceptedNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.AcceptedNodes(ctx)
	if err != nil {
		return nil, err
	}

	var peers []*pb.AcceptedNodesReply_Peer
	for _, node := range nodes {
		peers = append(peers, &pb.AcceptedNodesReply_Peer{
			NodeID: node.ID,
		})
	}

	resp := &pb.AcceptedNodesReply{
		Peers: peers,
	}
	log.WithContext(ctx).WithField("resp", resp).Debugf("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.ProcessUserdataServer.ConnectTo()
func (service *ProcessUserdata) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID, userdata.NodeTypePrimary); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")
	return resp, nil
}

// SendUserdata implements walletnode.ProcessUserdataServer.SendUserdata()
func (service *ProcessUserdata) SendUserdata(ctx context.Context, req *pb.UserdataRequest) (*pb.UserdataReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdata request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}
	// Convert protobuf request to UserdataProcessRequest
	request := userdata.UserdataProcessRequestSigned{
		Userdata: &userdata.UserdataProcessRequest{
			Realname:        req.Realname,
			FacebookLink:    req.FacebookLink,
			TwitterLink:     req.TwitterLink,
			NativeCurrency:  req.NativeCurrency,
			Location:        req.Location,
			PrimaryLanguage: req.PrimaryLanguage,
			Categories:      req.Categories,
			Biography:       req.Biography,
			AvatarImage: userdata.UserImageUpload{
				Content:  req.AvatarImage.Content,
				Filename: req.AvatarImage.Filename,
			},
			CoverPhoto: userdata.UserImageUpload{
				Content:  req.CoverPhoto.Content,
				Filename: req.CoverPhoto.Filename,
			},
			ArtistPastelID:    req.ArtistPastelID,
			Timestamp:         req.Timestamp,
			PreviousBlockHash: req.PreviousBlockHash,
		},
		UserdataHash: req.UserdataHash,
		Signature:    req.Signature,
	}

	processResult, err := task.SupernodeProcessUserdata(ctx, &request)
	if err != nil {
		return nil, errors.Errorf("SupernodeProcessUserdata can not process %w", err)
	}
	if processResult.ResponseCode == userdata.ErrorOnContent {
		return &pb.UserdataReply{
			ResponseCode:    processResult.ResponseCode,
			Detail:          processResult.Detail,
			Realname:        processResult.Realname,
			FacebookLink:    processResult.FacebookLink,
			TwitterLink:     processResult.TwitterLink,
			NativeCurrency:  processResult.NativeCurrency,
			Location:        processResult.Location,
			PrimaryLanguage: processResult.PrimaryLanguage,
			Categories:      processResult.Categories,
			Biography:       processResult.Biography,
			AvatarImage:     processResult.AvatarImage,
			CoverPhoto:      processResult.CoverPhoto,
		}, nil
	}
	// Process actual write to rqlite db happen here
	<-task.NewAction(func(ctx context.Context) error {
		// Send data to SN contain the leader rqlite
		/* NEED TO FIX!!!
		if err := task.ConnectTo(ctx, req.NodeID, req.SessID, userdata.NodeTypeLeader); err != nil {
			return err
		} else {
			if task.connectedToLeader != nil {
				if _, err := task.connectedToLeader.ProcessUserdata.SendUserdataToLeader(ctx, request); err != nil {
					return errors.Errorf("failed to send userdata to leader rqlite %w", err)
				}
			} else {
				return errors.Errorf("leader rqlite node object is empty")
			}

		}*/

		return nil
	})

	return &pb.UserdataReply{
		ResponseCode: processResult.ResponseCode,
		Detail:       processResult.Detail,
	}, nil
}

// ReceiveUserdata implements walletnode.ProcessUserdataServer.ReceiveUserdata()
func (service *ProcessUserdata) ReceiveUserdata(ctx context.Context, req *pb.RetrieveRequest) (*pb.UserdataRequest, error) {
	log.WithContext(ctx).WithField("userpastelid", req.Userpastelid).Debugf("ReceiveUserdata request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	userpastelid := req.Userpastelid
	result, err := task.ReceiveUserdata(ctx, userpastelid)
	if err != nil {
		return nil, err
	}

	// Generate protobuf response respProto
	respProto := &pb.UserdataRequest{
		Realname:        result.Realname,
		FacebookLink:    result.FacebookLink,
		TwitterLink:     result.TwitterLink,
		NativeCurrency:  result.NativeCurrency,
		Location:        result.Location,
		PrimaryLanguage: result.PrimaryLanguage,
		Categories:      result.Categories,
		Biography:       result.Biography,
		AvatarImage: &pb.UserdataRequest_UserImageUpload{
			Content:  result.AvatarImage.Content,
			Filename: result.AvatarImage.Filename,
		},
		CoverPhoto: &pb.UserdataRequest_UserImageUpload{
			Content:  result.CoverPhoto.Content,
			Filename: result.CoverPhoto.Filename,
		},
		ArtistPastelID:    result.ArtistPastelID,
		Timestamp:         result.Timestamp,
		PreviousBlockHash: result.PreviousBlockHash,
	}

	return respProto, nil
}

// Desc returns a description of the service.
func (service *ProcessUserdata) Desc() *grpc.ServiceDesc {
	return &pb.ProcessUserdata_ServiceDesc
}

// NewProcessUserdata returns a new ProcessUserdata instance.
func NewProcessUserdata(service *userdataprocess.Service) *ProcessUserdata {
	return &ProcessUserdata{
		ProcessUserdata: common.NewProcessUserdata(service),
	}
}
