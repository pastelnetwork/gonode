package walletnode

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"io"
	"runtime/debug"

	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/gonode/supernode/services/cascaderegister"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/walletnode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	commonSrvc "github.com/pastelnetwork/gonode/supernode/services/common"
)

// this implements SN's GRPC methods that are called by WNs during Sense Registration
// meaning - these methods implements server side of WN to SN GRPC communication

// RegisterCascade represents grpc service for registration Sense.
type RegisterCascade struct {
	pb.UnimplementedRegisterCascadeServer

	*common.RegisterCascade
}

// Session implements walletnode.RegisterCascadeServer.Session()
func (service *RegisterCascade) Session(stream pb.RegisterCascade_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *cascaderegister.CascadeRegistrationTask

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewCascadeRegistrationTask()
	}
	go func() {
		<-task.Done()
		cancel()
	}()
	defer task.Cancel()

	peer, _ := peer.FromContext(ctx)

	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debug("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receieve handshake request: %w", err)
	}

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
			case codes.Canceled:
				log.WithContext(ctx).WithError(err).Error("handshake stream canceled")
				return nil
			case codes.Unavailable:
				return nil
			}
			return errors.Errorf("handshake stream closed: %w", err)
		}
	}
}

// AcceptedNodes implements walletnode.RegisterSenseServer.AcceptedNodes()
func (service *RegisterCascade) AcceptedNodes(ctx context.Context, _ *pb.AcceptedNodesRequest) (*pb.AcceptedNodesReply, error) {
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, err
	}

	nodes, err := task.NetworkHandler.AcceptedNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("err from network handler.accepted-nodes: %w", err)
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
	log.WithContext(ctx).WithField("burn-txn", req.BurnTxid).Info("SendRegMetadata  rcvd")
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
	if err != nil {
		return nil, fmt.Errorf("send reg metadata: %w", err)
	}

	log.WithContext(ctx).WithField("burn-txn", req.BurnTxid).Info("SendRegMetadata responded")

	return &pb.SendRegMetadataReply{}, nil
}

// UploadAsset implements walletnode.RegisterNft.UploadAssetWithThumbnail
func (service *RegisterCascade) UploadAsset(stream pb.RegisterCascade_UploadAssetServer) (retErr error) {
	ctx := stream.Context()
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenUploadAsset")
		retErr = recErr
	})

	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return errors.Errorf("task not found %w", err)
	}

	asset := service.Storage.NewFile()
	assetFile, err := asset.Create()
	if err != nil {
		return errors.Errorf("open asset file %q: %w", assetFile.Name(), err)
	}
	log.WithContext(ctx).WithField("filename", assetFile.Name()).Info("UploadAsset request")

	assetWriter := bufio.NewWriter(assetFile)

	assetSize := int64(0)
	hash := make([]byte, 0)

	err = func() error {
		defer assetFile.Close()
		defer assetWriter.Flush()

		for {
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				if status.Code(err) == codes.Canceled {
					return errors.New("connection closed")
				}
				return errors.Errorf("receive UploadAsset: %w", err)
			}

			if assetPiece := req.GetAssetPiece(); assetPiece != nil {
				n, err := assetWriter.Write(assetPiece)
				assetSize += int64(n)
				if err != nil {
					return errors.Errorf("write to file %q: %w", assetFile.Name(), err)
				}
			} else {
				if metaData := req.GetMetaData(); metaData != nil {
					if metaData.Size != assetSize {
						return errors.Errorf("incomplete payload, send = %d receive=%d", metaData.Size, assetSize)
					}

					if metaData.Hash == nil {
						return errors.Errorf("empty hash")
					}
					hash = metaData.Hash

					if metaData.Format != "" {
						if err := asset.SetFormatFromExtension(metaData.Format); err != nil {
							return errors.Errorf("set format %s for file %s", metaData.Format, asset.Name())
						}
					}
				}
			}
		}

		return nil
	}()

	if err != nil {
		return err
	}

	err = func() error {
		hasher := sha3.New256()
		f, fileErr := asset.Open()
		if fileErr != nil {
			return errors.Errorf("open file %w", fileErr)
		}
		defer f.Close()

		if _, err := io.Copy(hasher, f); err != nil {
			return errors.Errorf("hash file failed %w", err)
		}
		hashFromPayload := hasher.Sum(nil)
		if !bytes.Equal(hashFromPayload, hash) {
			log.WithContext(ctx).WithField("Filename", assetFile.Name()).Debugf("calculated from payload %s", base64.URLEncoding.EncodeToString(hashFromPayload))
			log.WithContext(ctx).WithField("Filename", assetFile.Name()).Debugf("sent by client %s", base64.URLEncoding.EncodeToString(hash))
			return errors.Errorf("wrong hash")
		}

		return nil
	}()

	if err != nil {
		return fmt.Errorf("hash file: %w", err)
	}

	err = task.UploadAsset(ctx, asset)
	if err != nil {
		return fmt.Errorf("upload asset: %w", err)
	}

	resp := &pb.UploadAssetReply{}

	if err := stream.SendAndClose(resp); err != nil {
		return errors.Errorf("send UploadAsset response: %w", err)
	}

	log.WithContext(ctx).WithField("filename", assetFile.Name()).Info("UploadAsset request success")

	return nil
}

// SendSignedActionTicket implements walletnode.RegisterSense.SendSignedActionTicket
func (service *RegisterCascade) SendSignedActionTicket(ctx context.Context, req *pb.SendSignedCascadeTicketRequest) (retRes *pb.SendSignedActionTicketReply, retErr error) {
	defer errors.Recover(func(recErr error) {
		log.WithContext(ctx).WithField("stack-strace", string(debug.Stack())).Error("PanicWhenSendSignedActionTicket")
		retErr = recErr
	})

	log.WithContext(ctx).Info("SignTicket request rcvd")
	task, err := service.TaskFromMD(ctx)
	if err != nil {
		return nil, errors.Errorf("get task from metada %w", err)
	}

	// Validate burn_txid
	err = task.PastelHandler.ValidateBurnTxID(ctx, task.ActionTicketRegMetadata.BurnTxID, float64(task.ActionTicketRegMetadata.EstimatedFee))
	if err != nil {
		task.UpdateStatus(commonSrvc.StatusErrorInvalidBurnTxID)
		log.WithContext(ctx).WithError(err).Error("validate burn request failure")
		return nil, errors.Errorf("pre-burn txid is bad %w", err)
	}
	log.WithContext(ctx).Info("Burn txn validated")

	actionRegTxid, err := task.ValidateAndRegister(ctx, req.ActionTicket, req.CreatorSignature, req.RqFiles)
	if err != nil {
		return nil, errors.Errorf("get total storage fee: %w", err)
	}

	rsp := pb.SendSignedActionTicketReply{
		ActionRegTxid: actionRegTxid,
	}
	log.WithContext(ctx).WithField("rsp.retxid", actionRegTxid).Info("SignTicket request success")

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

	err = task.ValidateActionActAndConfirm(ctx, req.ActionRegTxid)

	return &pb.SendActionActReply{}, err
}

// GetTopMNs implements walletnode.RegisterCascadeServer.GetTopMNs()
func (service *RegisterCascade) GetTopMNs(ctx context.Context, req *pb.GetTopMNsRequest) (*pb.GetTopMNsReply, error) {
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

// Desc returns a description of the service.
func (service *RegisterCascade) Desc() *grpc.ServiceDesc {
	return &pb.RegisterCascade_ServiceDesc
}

// NewRegisterCascade returns a new RegisterSense instance.
func NewRegisterCascade(service *cascaderegister.CascadeRegistrationService) *RegisterCascade {
	return &RegisterCascade{
		RegisterCascade: common.NewRegisterCascade(service),
	}
}
