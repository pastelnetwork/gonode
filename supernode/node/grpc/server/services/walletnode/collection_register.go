package walletnode

import (
	"context"
	"io"
	"runtime/debug"

	"google.golang.org/grpc"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/collectionregister"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// this implements SN's GRPC methods that are called by WNs during Collection Registration
// meaning - these methods implements server side of WN to SN GRPC communication

// RegisterCollection represents grpc service for registration Collection.
type RegisterCollection struct {
	pb.UnimplementedRegisterCollectionServer

	*common.RegisterCollection
}

// Session implements wallet-node.RegisterCollectionServer.Session()
func (service *RegisterCollection) Session(stream pb.RegisterCollection_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *collectionregister.CollectionRegistrationTask

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewCollectionRegistrationTask()
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

// AcceptedNodes implements wallet-node.RegisterCollectionServer.AcceptedNodes()
func (service *RegisterCollection) AcceptedNodes(ctx context.Context, req *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
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

// ConnectTo implements wallet-node.RegisterCollectionServer.ConnectTo()
func (service *RegisterCollection) ConnectTo(ctx context.Context, req *pb.ConnectToRequest) (*pb.ConnectToReply, error) {
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

// MeshNodes implements wallet-node.RegisterCollectionServer.MeshNodes
func (service *RegisterCollection) MeshNodes(ctx context.Context, req *pb.MeshNodesRequest) (*pb.MeshNodesReply, error) {
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

// Desc returns a description of the service.
func (service *RegisterCollection) Desc() *grpc.ServiceDesc {
	return &pb.RegisterCollection_ServiceDesc
}

// SendCollectionTicketForSignature implements wallet-node.RegisterCollection.SendCollectionTicketForSignature
func (service *RegisterCollection) SendCollectionTicketForSignature(ctx context.Context, req *pb.SendCollectionTicketForSignatureRequest) (retRes *pb.SendCollectionTicketForSignatureResponse, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendCollectionTicketForSignature")
		retErr = recErr
	})

	log.WithContext(ctx).WithField("req", req).Debug("SignCollectionTicket request")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	if err := task.ValidatePreBurnTransaction(ctx, req.BurnTxid); err != nil {
		return nil, errors.Errorf("error validating pre-burn-txid")
	}

	collectionRegTxID, err := task.ValidateAndRegister(ctx, req.CollectionTicket, req.CreatorSignature)
	if err != nil {
		return nil, errors.Errorf("validate & register: %w", err)
	}

	rsp := pb.SendCollectionTicketForSignatureResponse{
		CollectionRegTxid: collectionRegTxID,
	}

	return &rsp, nil
}

// GetTopMNs implements walletnode.RegisterCollectionServer.GetTopMNs()
func (service *RegisterCollection) GetTopMNs(ctx context.Context, req *pb.GetTopMNsRequest) (*pb.GetTopMNsReply, error) {
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

	blnc, err := service.PastelClient.ZGetTotalBalance(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving balance on the SN")
		blnc = &pastel.GetTotalBalanceResponse{}
	}

	resp := &pb.GetTopMNsReply{
		MnTopList:   mnList,
		CurrBalance: int64(blnc.Total),
	}

	log.WithContext(ctx).WithField("mn-list", mnList).Debug("top mn-list has been returned")
	return resp, nil
}

// NewRegisterCollection returns a new RegisterCollection instance.
func NewRegisterCollection(service *collectionregister.CollectionRegistrationService) *RegisterCollection {
	return &RegisterCollection{
		RegisterCollection: common.NewRegisterCollection(service),
	}
}
