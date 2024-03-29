package supernode

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// StorageChallengeGRPC represents common grpc service for registration NFTs.
type StorageChallengeGRPC struct {
	pb.UnimplementedStorageChallengeServer

	*common.StorageChallenge
}

// Session represents a continuous storage challenge session stream.  This stream may be more fully implemented as storage challenge expands.
func (service *StorageChallengeGRPC) Session(stream pb.StorageChallenge_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *storagechallenge.SCTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewSCTask()
		isTaskNew = true
	}

	go func() {
		<-task.Done()
		cancel()
	}()

	if isTaskNew {
		defer task.Cancel()
	}

	peer, _ := peer.FromContext(ctx)
	log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream")
	defer log.WithContext(ctx).WithField("addr", peer.Addr).Debugf("Session stream closed")

	req, err := stream.Recv()
	if err != nil {
		return errors.Errorf("receive handshake request: %w", err)
	}
	log.WithContext(ctx).WithField("req", req).Debugf("Session request")

	if !isTaskNew {
		defer task.Cancel()
	}

	resp := &pb.SessionReply{
		SessID: task.ID(),
	}
	if err := stream.Send(resp); err != nil {
		return errors.Errorf("send handshake response: %w", err)
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

// Desc returns a description of the service.
func (service *StorageChallengeGRPC) Desc() *grpc.ServiceDesc {
	return &pb.StorageChallenge_ServiceDesc
}

// ProcessStorageChallenge is the server side of storage challenge processing GRPC comms
func (service *StorageChallengeGRPC) ProcessStorageChallenge(ctx context.Context, scRequest *pb.ProcessStorageChallengeRequest) (*pb.ProcessStorageChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Process Storage Challenge Request received from gRpc client")

	task := service.NewSCTask()

	msg := types.Message{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.MessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	_, err := task.ProcessStorageChallenge(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Processing Storage Challenge from Server Side")
	}

	return &pb.ProcessStorageChallengeReply{}, nil
}

// VerifyStorageChallenge is the server side of storage challenge verification GRPC comms
func (service *StorageChallengeGRPC) VerifyStorageChallenge(ctx context.Context, scRequest *pb.VerifyStorageChallengeRequest) (*pb.VerifyStorageChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Verify Storage Challenge Request received from gRpc client")
	task := service.NewSCTask()

	msg := types.Message{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.MessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	_, err := task.VerifyStorageChallenge(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying storage challenge")
	}

	return &pb.VerifyStorageChallengeReply{}, nil
}

// VerifyEvaluationResult is the server side of verify evaluation result
func (service *StorageChallengeGRPC) VerifyEvaluationResult(ctx context.Context, scRequest *pb.VerifyEvaluationResultRequest) (*pb.VerifyEvaluationResultReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Verify Evaluation Result request received from gRpc client")
	task := service.NewSCTask()

	msg := types.Message{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.MessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	resp, err := task.VerifyEvaluationResult(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying evaluation result")
		return nil, errors.Errorf("error verifying evaluation report")
	}

	d, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, errors.Errorf("error marshaling the evaluation result response")
	}

	return &pb.VerifyEvaluationResultReply{Data: &pb.StorageChallengeMessage{
		ChallengeId:     resp.ChallengeID,
		MessageType:     pb.StorageChallengeMessageMessageType(resp.MessageType),
		SenderId:        resp.Sender,
		SenderSignature: resp.SenderSignature,
		Data:            d,
	}}, nil
}

// BroadcastStorageChallengeResult broadcast the message to the entire network
func (service *StorageChallengeGRPC) BroadcastStorageChallengeResult(ctx context.Context, scRequest *pb.BroadcastStorageChallengeRequest) (*pb.BroadcastStorageChallengeResponse, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("broadcast storage challenge result request received from gRpc client")
	task := service.NewSCTask()

	msg := types.BroadcastMessage{
		ChallengeID: scRequest.ChallengeId,
		Challenger:  scRequest.Challenger,
		Recipient:   scRequest.Recipient,
		Observers:   scRequest.Observers,
	}

	_, err := task.BroadcastStorageChallengeResult(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying evaluation result")
		return nil, errors.Errorf("error verifying evaluation report")
	}

	return nil, err
}

// BroadcastStorageChallengeMetrics is the server side of broadcast storage-challenge metrics
func (service *StorageChallengeGRPC) BroadcastStorageChallengeMetrics(ctx context.Context, m *pb.BroadcastStorageChallengeMetricsRequest) (*pb.BroadcastStorageChallengeMetricsReply, error) {
	log.WithContext(ctx).WithField("req", m).Debug("Broadcast Storage-Challenge Metrics Request received from gRpc client")
	task := service.NewSCTask()

	req := types.ProcessBroadcastChallengeMetricsRequest{
		Data:     m.Data,
		SenderID: m.SenderId,
	}

	err := task.ProcessBroadcastStorageChallengeMetrics(ctx, req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error processing storage-challenge metrics")
		return nil, err
	}

	return &pb.BroadcastStorageChallengeMetricsReply{}, nil
}

// NewStorageChallengeGRPC returns a new StorageChallenge instance.
func NewStorageChallengeGRPC(service *storagechallenge.SCService) *StorageChallengeGRPC {
	return &StorageChallengeGRPC{
		StorageChallenge: common.NewStorageChallenge(service),
	}
}
