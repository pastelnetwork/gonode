package walletnode

import (
	"context"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"google.golang.org/grpc"
	"io"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
)

// this implements SN's GRPC methods that are called by WNs during Sense Registration
// meaning - these methods implements server side of WN to SN GRPC communication

// RegisterSense represents grpc service for registration Sense.
type RegisterCascade struct {
	pb.UnimplementedRegisterCascadeServer

	*common.RegisterCascade
}

// Session implements walletnode.RegisterSenseServer.Session()
func (service *RegisterCascade) Session(stream pb.RegisterCascade_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *cascaderegister.CascadeRegistrationTask

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
	log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receieve handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := task.NetworkHandler.Session(ctx, req.IsPrimary); err != nil {
		return err
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("send handshake response: %w", err)
	}
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")

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

// AcceptedNodes implements walletnode.RegisterSenseServer.AcceptedNodes()
func (service *RegisterCascade) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.NetworkHandler.AcceptedNodes(ctx)
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
	log.WithContext(ctx).WithField("resp", resp).Debug("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterSenseServer.ConnectTo()
func (service *RegisterCascade) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.NetworkHandler.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")
	return resp, nil
}

// MeshNodes implements walletnode.RegisterSenseServer.MeshNodes
func (service *RegisterCascade) MeshNodes(ctx context.Context, req *pb.MeshNodesRequest) (*pb.MeshNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("MeshNodes request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	meshedNodes := []types.MeshedSuperNode{}
	for _, node := range req.GetNodes() {
		meshedNodes = append(meshedNodes, types.MeshedSuperNode{
			NodeID: node.NodeID,
			SessID: node.SessID,
		})
	}

	err = task.NetworkHandler.MeshNodes(ctx, meshedNodes)
	return &pb.MeshNodesReply{}, err
}

// SendRegMetadata informs to SNs metadata required for registration request like current block hash, creator,..
func (service *RegisterCascade) SendRegMetadata(ctx context.Context, req *pb.SendRegMetadataRequest) (*pb.SendRegMetadataReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("SendRegMetadata request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	reqMetadata := &types.ActionRegMetadata{
		BlockHash:       req.BlockHash,
		CreatorPastelID: req.CreatorPastelID,
		BurnTxID:        req.BurnTxid,
	}

	err = task.SendRegMetadata(ctx, reqMetadata)

	return &pb.SendRegMetadataReply{}, err
}

// SendSignedActionTicket implements walletnode.RegisterSense.SendSignedActionTicket
func (service *RegisterCascade) SendSignedActionTicket(ctx context.Context, req *pb.SendSignedCascadeTicketRequest) (retRes *pb.SendSignedActionTicketReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedActionTicket")
		retErr = recErr
	})

	log.WithContext(ctx).WithField("req", req).Debug("SignTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	// Validate burn_txid
	err = task.ValidateBurnTxID(ctx)
	if err != nil {
		return nil, errors.Errorf("pre-burn txid is bad %w", err)
	}

	actionRegTxid, err := task.ValidateAndRegister(ctx, req.ActionTicket, req.CreatorSignature, req.RqFiles, req.EncodeParameters.Oti)
	if err != nil {
		return nil, errors.Errorf("get total storage fee %w", err)
	}

	rsp := pb.SendSignedActionTicketReply{
		ActionRegTxid: actionRegTxid,
	}

	return &rsp, nil
}

// SendActionAct informs to SN that walletnode activated action_reg
func (service *RegisterCascade) SendActionAct(ctx context.Context, req *pb.SendActionActRequest) (retRes *pb.SendActionActReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedActionTicket")
		retErr = recErr
	})

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return &pb.SendActionActReply{}, err
	}

	err = task.ValidateActionActAndStore(ctx, req.ActionRegTxid)

	return &pb.SendActionActReply{}, err
}

// Desc returns a description of the service.
func (service *RegisterCascade) Desc() *grpc.ServiceDesc {
	return &pb.RegisterCascade_ServiceDesc
}

// NewRegisterSense returns a new RegisterSense instance.
func NewRegisterCascade(service *cascaderegister.CascadeRegistrationService) *RegisterCascade {
	return &RegisterCascade{
		RegisterCascade: common.NewRegisterCascade(service),
	}
}
