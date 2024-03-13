package client

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/errors"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type healthCheckChallengeGRPCClient struct {
	client pb.HealthCheckChallengeClient
	conn   *clientConn
	sessID string
}

func (service *healthCheckChallengeGRPCClient) SessID() string {
	return service.sessID
}

func (service *healthCheckChallengeGRPCClient) Session(ctx context.Context, nodeID, sessID string) error {
	service.sessID = sessID

	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	stream, err := service.client.Session(ctx)
	if err != nil {
		return errors.Errorf("open Health stream: %w", err)
	}

	req := &pb.SessionRequest{
		NodeID: nodeID,
	}

	if err := stream.Send(req); err != nil {
		return errors.Errorf("send Session request: %w", err)
	}

	_, err = stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		switch status.Code(err) {
		case codes.Canceled, codes.Unavailable:
			return nil
		}
		return errors.Errorf("receive Session response: %w", err)
	}

	go func() {
		defer service.conn.Close()
		for {
			if _, err := stream.Recv(); err != nil {
				return
			}
		}
	}()

	return nil
}

func (service *healthCheckChallengeGRPCClient) ProcessHealthCheckChallenge(ctx context.Context, challengeMessage *pb.HealthCheckChallengeMessage) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.ProcessHealthCheckChallenge(ctx, &pb.ProcessHealthCheckChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error sending request from healthcheck challenge grpc client to server")
		return err
	}

	return nil
}

func (service *healthCheckChallengeGRPCClient) VerifyHealthCheckChallenge(ctx context.Context, challengeMessage *pb.HealthCheckChallengeMessage) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.VerifyHealthCheckChallenge(ctx, &pb.VerifyHealthCheckChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		return err
	}

	return nil
}

// VerifyHealthCheckEvaluationResult sends a gRPC request to observers
func (service *healthCheckChallengeGRPCClient) VerifyHealthCheckEvaluationResult(ctx context.Context, challengeMessage *pb.HealthCheckChallengeMessage) (types.HealthCheckMessage, error) {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	res, err := service.client.VerifyHealthCheckEvaluationResult(ctx, &pb.VerifyHealthCheckEvaluationResultRequest{
		Data: challengeMessage,
	})
	if err != nil {
		return types.HealthCheckMessage{}, err
	}

	msg := types.HealthCheckMessage{
		ChallengeID:     res.Data.ChallengeId,
		MessageType:     types.HealthCheckMessageType(res.Data.MessageType),
		Sender:          res.Data.SenderId,
		SenderSignature: res.Data.SenderSignature,
	}

	if err := json.Unmarshal(res.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return msg, errors.Errorf("error un-marshaling the received challenge message")
	}

	return msg, nil
}

// BroadcastHealthCheckChallengeResult broadcast the result to the entire network
func (service *healthCheckChallengeGRPCClient) BroadcastHealthCheckChallengeResult(ctx context.Context, req *pb.BroadcastHealthCheckChallengeRequest) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	_, err := service.client.BroadcastHealthCheckChallengeResult(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// BroadcastHealthCheckChallengeMetrics broadcast the result to the entire network
func (service *healthCheckChallengeGRPCClient) BroadcastHealthCheckChallengeMetrics(ctx context.Context, req types.ProcessBroadcastHealthCheckChallengeMetricsRequest) error {
	_, err := service.client.BroadcastHealthCheckChallengeMetrics(ctx, &pb.BroadcastHealthCheckChallengeMetricsRequest{
		Data:     req.Data,
		SenderId: req.SenderID,
	})
	if err != nil {
		return err
	}

	return nil
}

func newHealthCheckChallengeGRPCClient(conn *clientConn) node.HealthCheckChallengeInterface {
	return &healthCheckChallengeGRPCClient{
		conn:   conn,
		client: pb.NewHealthCheckChallengeClient(conn),
	}
}
