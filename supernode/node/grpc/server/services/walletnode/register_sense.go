package walletnode

import (
	"bufio"
	"context"
	"io"
	"runtime/debug"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/walletnode/register_sense"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/senseregister"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// RegisterSense represents grpc service for registration artwork.
type RegisterSense struct {
	pb.UnimplementedRegisterSenseServer

	*common.RegisterSense
}

// Session implements walletnode.RegisterSenseServer.Session()
func (service *RegisterSense) Session(stream pb.RegisterSense_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *senseregister.Task

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

// AcceptedNodes implements walletnode.RegisterSenseServer.AcceptedNodes()
func (service *RegisterSense) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
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

// ConnectTo implements walletnode.RegisterSenseServer.ConnectTo()
func (service *RegisterSense) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
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

// MeshNodes implements walletnode.RegisterSenseServer.MeshNodes
func (service *RegisterSense) MeshNodes(ctx context.Context, req *pb.MeshNodesRequest) (*pb.MeshNodesReply, error) {
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
func (service *RegisterSense) SendRegMetadata(ctx context.Context, req *pb.SendRegMetadataRequest) (*pb.SendRegMetadataReply, error) {
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

// ProbeImage implements walletnode.RegisterSenseServer.ProbeImage()
func (service *RegisterSense) ProbeImage(stream pb.RegisterSense_ProbeImageServer) (retErr error) {
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
	file, err := image.Create()
	if err != nil {
		return errors.Errorf("open file %q: %w", file.Name(), err)
	}
	defer file.Close()
	log.WithContext(ctx).WithField("filename", file.Name()).Debug("ProbeImage request")

	wr := bufio.NewWriter(file)
	defer wr.Flush()

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.Canceled {
				return errors.New("connection closed")
			}
			return errors.Errorf("receive ProbeImage: %w", err)
		}

		if _, err := wr.Write(req.Image); err != nil {
			return errors.Errorf("write to file %q: %w", file.Name(), err)
		}
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

// SendSignedActionTicket implements walletnode.RegisterSense.SendSignedActionTicket
func (service *RegisterSense) SendSignedActionTicket(ctx context.Context, req *pb.SendSignedActionTicketRequest) (retRes *pb.SendSignedActionTicketReply, retErr error) {
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

// Desc returns a description of the service.
func (service *RegisterSense) Desc() *grpc.ServiceDesc {
	return &pb.RegisterSense_ServiceDesc
}

// NewRegisterSense returns a new RegisterSense instance.
func NewRegisterSense(service *senseregister.Service) *RegisterSense {
	return &RegisterSense{
		RegisterSense: common.NewRegisterSense(service),
	}
}
