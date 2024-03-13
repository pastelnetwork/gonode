package supernode

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/healthcheckchallenge"
)

// HealthCheckChallengeGRPC represents grpc service for health check
type HealthCheckChallengeGRPC struct {
	pb.UnimplementedHealthCheckChallengeServer

	*common.HealthCheckChallenge
}

// Session represents a continuous health check challenge session stream.  This stream may be more fully implemented as health check challenge expands.
func (service *HealthCheckChallengeGRPC) Session(stream pb.HealthCheckChallenge_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *healthcheckchallenge.HCTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewHCTask()
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
func (service *HealthCheckChallengeGRPC) Desc() *grpc.ServiceDesc {
	return &pb.HealthCheckChallenge_ServiceDesc
}

// ProcessHealthCheckChallenge is the server side of healthcheck challenge processing GRPC comms
func (service *HealthCheckChallengeGRPC) ProcessHealthCheckChallenge(ctx context.Context, scRequest *pb.ProcessHealthCheckChallengeRequest) (*pb.ProcessHealthCheckChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debug("Process HealthCheck Challenge Request received from gRpc client")

	task := service.NewHCTask()

	msg := types.HealthCheckMessage{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.HealthCheckMessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	_, err := task.ProcessHealthCheckChallenge(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Processing Health Check Challenge from Server Side")
	}

	return &pb.ProcessHealthCheckChallengeReply{}, nil
}

// VerifyHealthCheckChallenge is the server side of health check challenge verification GRPC comms
func (service *HealthCheckChallengeGRPC) VerifyHealthCheckChallenge(ctx context.Context, scRequest *pb.VerifyHealthCheckChallengeRequest) (*pb.VerifyHealthCheckChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Verify Health Check Challenge Request received from gRpc client")
	task := service.NewHCTask()

	msg := types.HealthCheckMessage{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.HealthCheckMessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	_, err := task.VerifyHealthCheckChallenge(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying health check challenge")
	}

	return &pb.VerifyHealthCheckChallengeReply{}, nil
}

// VerifyHealthCheckEvaluationResult is the server side of verify evaluation result
func (service *HealthCheckChallengeGRPC) VerifyHealthCheckEvaluationResult(ctx context.Context, scRequest *pb.VerifyHealthCheckEvaluationResultRequest) (*pb.VerifyHealthCheckEvaluationResultReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Verify HealthCheck Evaluation Result request received from gRpc client")
	task := service.NewHCTask()

	msg := types.HealthCheckMessage{
		ChallengeID:     scRequest.Data.ChallengeId,
		MessageType:     types.HealthCheckMessageType(scRequest.Data.MessageType),
		Sender:          scRequest.Data.SenderId,
		SenderSignature: scRequest.Data.SenderSignature,
	}

	if err := json.Unmarshal(scRequest.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return nil, errors.Errorf("error un-marshaling the received challenge message")
	}

	resp, err := task.VerifyHealthCheckEvaluationResult(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying health check evaluation result")
		return nil, errors.Errorf("error verifying health check evaluation report")
	}

	d, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, errors.Errorf("error marshaling the health check evaluation result response")
	}

	return &pb.VerifyHealthCheckEvaluationResultReply{Data: &pb.HealthCheckChallengeMessage{
		ChallengeId:     resp.ChallengeID,
		MessageType:     pb.HealthCheckChallengeMessageMessageType(resp.MessageType),
		SenderId:        resp.Sender,
		SenderSignature: resp.SenderSignature,
		Data:            d,
	}}, nil
}

// BroadcastHealthCheckChallengeResult broadcast the message to the entire network
func (service *HealthCheckChallengeGRPC) BroadcastHealthCheckChallengeResult(ctx context.Context, scRequest *pb.BroadcastHealthCheckChallengeRequest) (*pb.BroadcastHealthCheckChallengeResponse, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("broadcast healthcheck challenge result request received from gRpc client")
	task := service.NewHCTask()

	msg := types.BroadcastHealthCheckMessage{
		ChallengeID: scRequest.ChallengeId,
		Challenger:  scRequest.Challenger,
		Recipient:   scRequest.Recipient,
		Observers:   scRequest.Observers,
	}

	_, err := task.BroadcastHealthCheckChallengeResult(ctx, msg)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error verifying broadcast result")
		return nil, errors.Errorf("error verifying broadcast report")
	}

	return nil, err
}

// BroadcastHealthCheckChallengeMetrics is the server side of broadcast healthcheck-challenge metrics
func (service *HealthCheckChallengeGRPC) BroadcastHealthCheckChallengeMetrics(ctx context.Context, m *pb.BroadcastHealthCheckChallengeMetricsRequest) (*pb.BroadcastHealthCheckChallengeMetricsReply, error) {
	log.WithContext(ctx).WithField("req", m).Debug("Broadcast HealthCheck-Challenge Metrics Request received from gRpc client")
	task := service.NewHCTask()

	req := types.ProcessBroadcastHealthCheckChallengeMetricsRequest{
		Data:     m.Data,
		SenderID: m.SenderId,
	}

	err := task.ProcessBroadcastHealthCheckChallengeMetrics(ctx, req)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error processing healthcheck-challenge metrics")
		return nil, err
	}

	return &pb.BroadcastHealthCheckChallengeMetricsReply{}, nil
}

// NewHealthCheckChallengeGRPC returns a new HealthCheckChallenge instance.
func NewHealthCheckChallengeGRPC(service *healthcheckchallenge.HCService) *HealthCheckChallengeGRPC {
	return &HealthCheckChallengeGRPC{
		HealthCheckChallenge: common.NewHealthCheckChallenge(service),
	}
}
