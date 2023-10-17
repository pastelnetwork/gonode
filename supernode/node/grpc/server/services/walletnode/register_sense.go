package walletnode

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/senseregister"
)

// this implements SN's GRPC methods that are called by WNs during Sense Registration
// meaning - these methods implements server side of WN to SN GRPC communication

// RegisterSense represents grpc service for registration Sense.
type RegisterSense struct {
	pb.UnimplementedRegisterSenseServer

	*common.RegisterSense
}

// Session implements walletnode.RegisterSenseServer.Session()
func (service *RegisterSense) Session(stream pb.RegisterSense_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *senseregister.SenseRegistrationTask

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewSenseRegistrationTask()
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
func (service *RegisterSense) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
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

// GetDDDatabaseHash implements walletnode.RegisterSenseServer.GetDupeDetectionDBHash()
func (service *RegisterSense) GetDDDatabaseHash(ctx context.Context, _ *pb.GetDBHashRequest) (*pb.DBHashReply, error) {
	hash, err := service.GetDupeDetectionDatabaseHash(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("GetDDDatabaseHash")
		return nil, fmt.Errorf("GetDDDatabaseHash: %w", err)
	}

	resp := &pb.DBHashReply{
		Hash: hash,
	}

	log.WithContext(ctx).WithField("hash len", len(hash)).Debug("DupeDetection DB hash returned")

	return resp, nil
}

// GetDDServerStats implements walletnode.RegisterSenseServer.GetDDServerStats()
func (service *RegisterSense) GetDDServerStats(ctx context.Context, _ *pb.DDServerStatsRequest) (*pb.DDServerStatsReply, error) {
	log.WithContext(ctx).Info("request for dd-server stats has been received")

	resp, err := service.GetDupeDetectionServerStats(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving dd-server stats from SN gRPC call")
	}

	stats := &pb.DDServerStatsReply{
		MaxConcurrent:  resp.MaxConcurrent,
		WaitingInQueue: resp.WaitingInQueue,
		Executing:      resp.Executing,
	}

	log.WithContext(ctx).WithField("stats", stats).Info("DD server stats returned")
	return stats, nil
}

// GetTopMNs implements walletnode.RegisterSenseServer.GetTopMNs()
func (service *RegisterSense) GetTopMNs(ctx context.Context, req *pb.GetTopMNsRequest) (*pb.GetTopMNsReply, error) {
	log.WithContext(ctx).Debug("request for mn-top list has been received")

	var mnTopList pastel.MasterNodes
	var err error
	if req.Block == 0 {
		mnTopList, err = service.PastelClient.MasterNodesTop(ctx)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving mn-top list on the SN")
		}
	} else {
		mnTopList, err = service.PastelClient.MasterNodesTopN(ctx, int(req.Block))
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("error retrieving mn-top list on the SN")
		}
	}

	var mnList []string
	for _, mn := range mnTopList {
		mnList = append(mnList, mn.ExtAddress)
	}

	resp := &pb.GetTopMNsReply{
		MnTopList: mnList,
	}

	log.WithContext(ctx).WithField("mn-list", mnList).Debug("top mn-list has been returned")
	return resp, nil
}

// ConnectTo implements walletnode.RegisterSenseServer.ConnectTo()
func (service *RegisterSense) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
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

	err = task.NetworkHandler.MeshNodes(ctx, meshedNodes)
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
		BlockHeight:     req.BlockHeight,
		Timestamp:       req.Timestamp,
		GroupID:         req.GroupId,
		CollectionTxID:  req.CollectionTxid,
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
		log.WithContext(ctx).WithError(err).Error("Unable to retrieve TaskFromMD")
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
				if subErr == io.EOF {
					break
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}
				return errors.Errorf("receive ProbeImage: %w", err)
			}

			if _, subErr := wr.Write(req.Payload); subErr != nil {
				return errors.Errorf("write to file %q: %w", file.Name(), subErr)
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
		log.WithContext(ctx).WithError(err).Error("updateFormatError")
		return errors.Errorf("add image format: %w", err)
	}

	exists, err := task.HashExists(ctx, image)
	if exists || err != nil {
		resp := &pb.ProbeImageReply{
			IsExisting:                        exists,
			CompressedSignedDDAndFingerprints: []byte{},
		}

		if err != nil {
			log.WithContext(ctx).WithError(err).Error("hash exist check failed")
			resp.ErrString = err.Error()
		}

		return stream.SendAndClose(resp)
	}

	if err := task.CalculateFee(ctx, image); err != nil {
		log.WithContext(ctx).WithError(err).Error("calculate fee failed")
		return errors.Errorf("calculate fee: %w", err)
	}

	// Validate burn_txid
	isValidBurnTxID := false
	var compressedDDFingerAndScores []byte

	err = task.ValidateBurnTxID(ctx, 20)
	if err == nil {
		isValidBurnTxID = true

		compressedDDFingerAndScores, err = task.ProbeImage(ctx, image)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("probe image failed")
			return fmt.Errorf("task.ProbeImage: %w", err)
		}
	} else {
		log.WithContext(ctx).WithError(err).Error("validate burn txid failed")
	}

	resp := &pb.ProbeImageReply{
		CompressedSignedDDAndFingerprints: compressedDDFingerAndScores,
		IsValidBurnTxid:                   isValidBurnTxID,
	}

	if err != nil && !isValidBurnTxID {
		resp.ErrString = err.Error()
	}

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send ProbeImage response: %w", err)
	}

	return nil
}

// SendSignedActionTicket implements walletnode.RegisterSense.SendSignedActionTicket
func (service *RegisterSense) SendSignedActionTicket(ctx context.Context, req *pb.SendSignedSenseTicketRequest) (retRes *pb.SendSignedActionTicketReply, retErr error) {
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
		return nil, errors.Errorf("validate & register: %w", err)
	}

	rsp := pb.SendSignedActionTicketReply{
		ActionRegTxid: actionRegTxid,
	}

	return &rsp, nil
}

// SendActionAct informs to SN that walletnode activated action_reg
func (service *RegisterSense) SendActionAct(ctx context.Context, req *pb.SendActionActRequest) (retRes *pb.SendActionActReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedActionTicket")
		retErr = recErr
	})

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return &pb.SendActionActReply{}, err
	}

	err = task.ValidateActionActAndConfirm(ctx, req.ActionRegTxid)

	return &pb.SendActionActReply{}, err
}

// Desc returns a description of the service.
func (service *RegisterSense) Desc() *grpc.ServiceDesc {
	return &pb.RegisterSense_ServiceDesc
}

// NewRegisterSense returns a new RegisterSense instance.
func NewRegisterSense(service *senseregister.SenseRegistrationService) *RegisterSense {
	return &RegisterSense{
		RegisterSense: common.NewRegisterSense(service),
	}
}
