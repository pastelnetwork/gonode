package client

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type storageChallengeGRPCClient struct {
	client pb.StorageChallengeClient
	conn   *clientConn
	sessID string
}

func (service *storageChallengeGRPCClient) SessID() string {
	return service.sessID
}

func (service *storageChallengeGRPCClient) Session(ctx context.Context, nodeID, sessID string) error {
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

func (service *storageChallengeGRPCClient) ProcessStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeMessage) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.ProcessStorageChallenge(ctx, &pb.ProcessStorageChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error sending request from storage challenge grpc client to server")
		return err
	}

	return nil
}

func (service *storageChallengeGRPCClient) VerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeMessage) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	_, err := service.client.VerifyStorageChallenge(ctx, &pb.VerifyStorageChallengeRequest{
		Data: challengeMessage,
	})
	if err != nil {
		return err
	}

	return nil
}

// VerifyEvaluationResult sends a gRPC request to observers
func (service *storageChallengeGRPCClient) VerifyEvaluationResult(ctx context.Context, challengeMessage *pb.StorageChallengeMessage) (types.Message, error) {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)
	res, err := service.client.VerifyEvaluationResult(ctx, &pb.VerifyEvaluationResultRequest{
		Data: challengeMessage,
	})
	if err != nil {
		return types.Message{}, err
	}

	msg := types.Message{
		ChallengeID:     res.Data.ChallengeId,
		MessageType:     types.MessageType(res.Data.MessageType),
		Sender:          res.Data.SenderId,
		SenderSignature: res.Data.SenderSignature,
	}

	if err := json.Unmarshal(res.Data.Data, &msg.Data); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error un-marshaling received challenge message")
		return msg, errors.Errorf("error un-marshaling the received challenge message")
	}

	return msg, nil
}

// BroadcastStorageChallengeResult broadcast the result to the entire network
func (service *storageChallengeGRPCClient) BroadcastStorageChallengeResult(ctx context.Context, req *pb.BroadcastStorageChallengeRequest) error {
	ctx = contextWithLogPrefix(ctx, service.conn.id)
	ctx = contextWithMDSessID(ctx, service.sessID)

	_, err := service.client.BroadcastStorageChallengeResult(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func newStorageChallengeGRPCClient(conn *clientConn) node.StorageChallengeInterface {
	return &storageChallengeGRPCClient{
		conn:   conn,
		client: pb.NewStorageChallengeClient(conn),
	}
}
