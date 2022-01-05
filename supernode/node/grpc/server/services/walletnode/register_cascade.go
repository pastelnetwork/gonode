package walletnode

import (
	"bufio"
	"context"
	"io"
	"runtime/debug"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/walletnode/register_cascade"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterCascade represents grpc service for registration artwork.
type RegisterCascade struct {
	pb.UnimplementedRegisterCascadeServer

	*common.RegisterCascade
}

// Session implements walletnode.RegisterCascadeServer.Session()
func (service *RegisterCascade) Session(stream pb.RegisterCascade_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *cascaderegister.Task

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

	if err := task.Session(ctx, req.IsPrimary); err != nil {
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

// AcceptedNodes implements walletnode.RegisterCascadeServer.AcceptedNodes()
func (service *RegisterCascade) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("AcceptedNodes request")
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
	log.WithContext(ctx).WithField("resp", resp).Debug("AcceptedNodes response")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterCascadeServer.ConnectTo()
func (service *RegisterCascade) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
	log.WithContext(ctx).WithField("req", req).Debug("ConnectTo request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ConnectTo(ctx, req.NodeID, req.SessID); err != nil {
		return nil, err
	}

	resp := &pb.ConnectToReply{}
	log.WithContext(ctx).WithField("resp", resp).Debug("ConnectTo response")
	return resp, nil
}

// MeshNodes implements walletnode.RegisterCascadeServer.MeshNodes
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

	err = task.MeshNodes(ctx, meshedNodes)
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

// ProbeImage implements walletnode.RegisterCascadeServer.ProbeImage()
func (service *RegisterCascade) ProbeImage(stream pb.RegisterCascade_ProbeImageServer) (retErr error) {
	ctx := stream.Context()

	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenProbeImage")
		retErr = recErr
	})

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return err
	}

	image := service.Storage.NewFile()

	// Do file receiving
	receiveFunc := func() error {
		file, subErr := image.Create()
		if subErr != nil {
			return errors.Errorf("open file %q: %w", file.Name(), err)
		}
		defer file.Close()
		log.WithContext(ctx).WithField("filename", file.Name()).Debug("ProbeImage request")

		wr := bufio.NewWriter(file)
		defer wr.Flush()

		for {
			req, subErr := stream.Recv()
			if subErr != nil {
				if err == io.EOF {
					break
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}
				return errors.Errorf("receive ProbeImage: %w", err)
			}

			if _, subErr := wr.Write(req.Image); subErr != nil {
				return errors.Errorf("write to file %q: %w", file.Name(), err)
			}
		}
		return nil
	}

	err = receiveFunc()
	if err != nil {
		return err
	}

	err = image.UpdateFormat()
	if err != nil {
		return errors.Errorf("add image format: %w", err)
	}

	// Validate burn_txid
	isValidBurnTxID := false
	var compressedDDFingerAndScores []byte

	err = task.ValidateBurnTxID(ctx)
	if err == nil {
		isValidBurnTxID = true

		compressedDDFingerAndScores, err = task.ProbeImage(ctx, image)
		if err != nil {
			return err
		}
	}

	resp := &pb.ProbeImageReply{
		CompressedSignedDDAndFingerprints: compressedDDFingerAndScores,
		IsValidBurnTxid:                   isValidBurnTxID,
	}

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send ProbeImage response: %w", err)
	}
	return nil
}

// SendSignedActionTicket implements walletnode.RegisterCascade.SendSignedActionTicket
func (service *RegisterCascade) SendSignedActionTicket(ctx context.Context, req *pb.SendSignedActionTicketRequest) (retRes *pb.SendSignedActionTicketReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedActionTicket")
		retErr = recErr
	})

	log.WithContext(ctx).WithField("req", req).Debug("SignTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	actionRegTxid, err := task.ValidateAndRegister(ctx, req.ActionTicket, req.CreatorSignature, req.DdFpFiles)
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

// NewRegisterCascade returns a new RegisterCascade instance.
func NewRegisterCascade(service *cascaderegister.Service) *RegisterCascade {
	return &RegisterCascade{
		RegisterCascade: common.NewRegisterCascade(service),
	}
}
