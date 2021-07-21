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

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debugf("ConnectTo response")
	return resp, nil
}

// SendSignedArtTicket implements walletnode.RegisterArtwork.SendSignedArtTicket
func (service *RegisterArtwork) SendSignedArtTicket(ctx context.Context, req *pb.SendSignedArtTicketRequest) (*pb.SendSignedArtTicketReply, error) {
	log.WithContext(ctx).WithField("req", req).Debugf("SignTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("failed to get task from metada %w", err)
	}

	registrationFee, err := task.GetRegistrationFee(ctx, req.ArtTicket, req.ArtistSignature, req.Key1, req.Key2, req.EncodeFiles, req.EncodeParameters.Oti)
	if err != nil {
		return nil, errors.Errorf("failed to get total storage fee %w", err)
	}

	rsp := pb.SendSignedArtTicketReply{
		RegistrationFee: registrationFee,
	}

	return &rsp, nil
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
