package client

import (
	"context"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/types"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type selfHealingGRPCClient struct {
	client pb.SelfHealingClient
	conn   *clientConn
	sessID string
}

func (service *selfHealingGRPCClient) SessID() string {
	return service.sessID
}

func (service *selfHealingGRPCClient) Session(ctx context.Context, nodeID, sessID string) error {
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
	log.WithContext(ctx).WithField("req", req).Debug("Session request")

	if err := stream.Send(req); err != nil {
		return errors.Errorf("send Session request: %w", err)
	}

	resp, err := stream.Recv()
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
	log.WithContext(ctx).WithField("resp", resp).Debug("Session response")

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

func (service *selfHealingGRPCClient) Ping(ctx context.Context, pingRequest *pb.PingRequest) (*pb.PingResponse, error) {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	res, err := service.client.Ping(ctx, pingRequest)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error sending request from ping and fetch node info worker")
		return nil, err
	}

	return res, nil
}

func (service *selfHealingGRPCClient) ProcessSelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingMessage) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.ProcessSelfHealingChallenge(ctx, &pb.ProcessSelfHealingChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error sending request from process self-healing challenge grpc client to server")
		return err
	}

	return nil
}
func (service *selfHealingGRPCClient) VerifySelfHealingChallenge(ctx context.Context, challengeMessage *pb.SelfHealingMessage) (types.SelfHealingMessage, error) {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	res, err := service.client.VerifySelfHealingChallenge(ctx, &pb.VerifySelfHealingChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		return types.SelfHealingMessage{}, err
	}

	msg := types.SelfHealingMessage{
		TriggerID:       res.Data.TriggerId,
		MessageType:     types.SelfHealingMessageType(res.Data.MessageType),
		SenderID:        res.Data.SenderId,
		SenderSignature: res.Data.SenderSignature,
	}

	if err := json.Unmarshal(res.Data.Data, &msg.SelfHealingMessageData); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return msg, errors.Errorf("error un-marshaling the received challenge message")
	}

	return msg, nil
}

func (service *selfHealingGRPCClient) BroadcastSelfHealingMetrics(ctx context.Context, req types.ProcessBroadcastMetricsRequest) error {
	_, err := service.client.BroadcastSelfHealingMetrics(ctx, &pb.BroadcastSelfHealingMetricsRequest{
		Type:            int64(req.Type),
		Data:            req.Data,
		SenderSignature: req.SenderSignature,
		SenderId:        req.SenderID,
	})
	if err != nil {
		return err
	}

	return nil
}

func newSelfHealingGRPCClient(conn *clientConn) node.SelfHealingChallengeInterface {
	return &selfHealingGRPCClient{
		conn:   conn,
		client: pb.NewSelfHealingClient(conn),
	}
}
