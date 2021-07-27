package walletnode

import (
	"bufio"
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
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

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID, common.NodeTypePrimary); err != nil {
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
	request := UserdataProcessRequestSigned{}

	request.Userdata.Realname = req.Realname
	request.Userdata.FacebookLink = req.Facebook_link
	request.Userdata.TwitterLink=req.Twitter_link
	request.Userdata.NativeCurrency= req.Native_currency
	request.Userdata.Location= req.Location
	request.Userdata.PrimaryLanguage= req.Primary_language
	request.Userdata.Categories=req.Categories
	request.Userdata.Biography= req.Biography

	if req.AvatarImage.Content != nil && len(req.AvatarImage.Content) > 0 {
		req.AvatarImage.Content = make ([]byte, len(req.AvatarImage.Content))
		copy(request.Userdata.AvatarImage.Content,req.AvatarImage.Content)
	}
	request.Userdata.AvatarImage.Filename = req.AvatarImage.Filename

	if req.CoverPhoto.Content != nil && len(req.CoverPhoto.Content) > 0 {
		req.CoverPhoto.Content = make ([]byte, len(req.CoverPhoto.Content))
		copy(request.Userdata.CoverPhoto.Content,req.CoverPhoto.Content)
	}
	request.Userdata.CoverPhoto.Filename := req.CoverPhoto.Filename

	request.Userdata.ArtistPastelID  = req.ArtistPastelID
	request.Userdata.Timestamp   = req.Timestamp
	request.Userdata.PreviousBlockHash=req.PreviousBlockHash
	request.UserdataHash = req.UserdataHash 
	request.Signature = req.Signature


	processResult := task.supernodeProcessUserdata(ctx, request)
	if processResult.ResponseCode == userdata.ErrorOnContent {
		return &pb.UserdataReply {
			Response_code 		: processResult.ResponseCode
			Detail				: processResult.Detail
			Realname 			: processResult.Realname,
			Facebook_link 		: processResult.FacebookLink,
			Twitter_link 		: processResult.TwitterLink,
			Native_currency 	: processResult.NativeCurrency,
			Location 			: processResult.Location,
			Primary_language 	: processResult.PrimaryLanguage,
			Categories 			: processResult.Categories,
			Biography 			: processResult.Biography,
			Avatar_image		: processResult.AvatarImage,
			Cover_photo			: processResult.CoverPhoto,
		}
	}
	// Process actual write to rqlite db happen here
	<-task.NewAction(func(ctx context.Context) error {
		// Send data to SN contain the leader rqlite
		if err := task.ConnectTo(ctx, req.NodeID, req.SessID, common.NodeTypeLeader); err != nil {
			return err
		} else {
			if task.connectedToLeader != nil {
				if _, err := task.connectedToLeader.ProcessUserdata.SendUserdataToLeader(ctx, request); err != nil {
					return errors.Errorf("failed to send userdata to leader rqlite node %s at address %s %w", task.connectedToLeader.ID, task.connectedToLeader.Address, err)
				}
			} else {
				return return errors.Errorf("leader rqlite node object is empty")
			}
			
		}

		return nil
	})

	return &pb.UserdataReply {
		Response_code		: processResult.ResponseCode
		Detail				: processResult.Detail
	}
}

// SendUserdata implements walletnode.ProcessUserdataServer.SendUserdata()
func (service *ProcessUserdata) ReceiveUserdata(ctx context.Context, userpastelid string) (*pb.UserdataRequest, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SendUserdata request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	resutl, err:= task.retrieveUserdata(userpastelid)
	if err != nil {
		return nil, err
	}
	
	// Generate protobuf response respProto
	respProto := &pb.UserdataRequest{}

	respProto.Realname = result.Userdata.Realname
	respProto.FacebookLink = result.Userdata.FacebookLink
	respProto.TwitterLink = result.Userdata.TwitterLink
	respProto.NativeCurrency = result.Userdata.NativeCurrency
	respProto.Location = result.Userdata.Location
	respProto.PrimaryLanguage = result.Userdata.PrimaryLanguage
	respProto.Categories = result.Userdata.Categories
	respProto.Biography = result.Userdata.Biography
	respProto.AvatarImage = &pb.UserdataRequest_UserImageUpload {}

	if result.Userdata.AvatarImage.Content != nil && len(result.Userdata.AvatarImage.Content) > 0 {
		respProto.AvatarImage.Content = make ([]byte, len(result.Userdata.AvatarImage.Content))
		copy(respProto.AvatarImage.Content,result.Userdata.AvatarImage.Content)
	}
	respProto.AvatarImage.Filename = result.Userdata.AvatarImage.Filename

	respProto.CoverPhoto = &pb.UserdataRequest_UserImageUpload {}
	if result.Userdata.CoverPhoto.Content != nil && len(result.Userdata.CoverPhoto.Content) > 0 {
		respProto.CoverPhoto.Content = make ([]byte, len(result.Userdata.CoverPhoto.Content))
		copy(respProto.CoverPhoto.Content,result.Userdata.CoverPhoto.Content)
	}
	respProto.CoverPhoto.Filename = result.Userdata.CoverPhoto.Filename
	
	respProto.ArtistPastelID = result.Userdata.ArtistPastelID 
	respProto.Timestamp = result.Userdata.Timestamp
	respProto.PreviousBlockHash = result.Userdata.PreviousBlockHash

	return respProto
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
