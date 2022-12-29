package supernode

import (
	"context"
	"io"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node/grpc/server/services/common"
	"github.com/pastelnetwork/gonode/supernode/services/selfhealing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// SelfHealingChallengeGRPC represents common grpc service for self-healing.
type SelfHealingChallengeGRPC struct {
	pb.UnimplementedSelfHealingServer

	*common.SelfHealingChallenge
}

// Session represents a continuous storage challenge session stream.  This stream may be more fully implemented as storage challenge expands.
func (service *SelfHealingChallengeGRPC) Session(stream pb.SelfHealing_SessionServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	var task *selfhealing.SHTask
	isTaskNew := false

	if sessID, ok := service.SessID(ctx); ok {
		if task = service.Task(sessID); task == nil {
			return errors.Errorf("not found %q task", sessID)
		}
	} else {
		task = service.NewSHTask()
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
func (service *SelfHealingChallengeGRPC) Desc() *grpc.ServiceDesc {
	return &pb.StorageChallenge_ServiceDesc
}

// ProcessSelfHealingChallenge is the server side of self-healing challenge processing GRPC comms
func (service *SelfHealingChallengeGRPC) ProcessSelfHealingChallenge(ctx context.Context, scRequest *pb.ProcessSelfHealingChallengeRequest) (*pb.ProcessSelfHealingChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Info("Process self-healing challenge request received from gRpc client")

	task := service.NewSHTask()
	err := task.ProcessSelfHealingChallenge(ctx, scRequest.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Processing Self-Healing Challenge from Server Side")
	}

	return &pb.ProcessSelfHealingChallengeReply{}, nil
}

// VerifySelfHealingChallenge is the server side of self-healing challenge verification GRPC comms
func (service *SelfHealingChallengeGRPC) VerifySelfHealingChallenge(ctx context.Context, scRequest *pb.VerifySelfHealingChallengeRequest) (*pb.VerifySelfHealingChallengeReply, error) {
	log.WithContext(ctx).WithField("req", scRequest).Debugf("Verify Self-Healing Request received from gRpc client")
	//task := service.NewSHTask()

	//data, err := task.VerifySelfHealingChallenge(ctx, scRequest.Data)
	//if err != nil {
	//	log.WithContext(ctx).WithError(err).Error("Error verifying Self-Healing")
	//}

	return &pb.VerifySelfHealingChallengeReply{}, nil
}

// NewSelfHealingChallengeGRPC returns a new SelfHealing instance.
func NewSelfHealingChallengeGRPC(service *selfhealing.SHService) *SelfHealingChallengeGRPC {
	return &SelfHealingChallengeGRPC{
		SelfHealingChallenge: common.NewSelfHealingChallenge(service),
	}
}