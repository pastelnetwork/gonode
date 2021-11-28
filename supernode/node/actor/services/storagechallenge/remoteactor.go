package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	dto "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

type storageChallengeActor struct {
	service storagechallenge.StorageChallenge
}

func (s *storageChallengeActor) Receive(actorCtx actor.Context) {
	// Begin transaction, inject to context to go through main process
	ctx := context.Context{}
	switch msg := actorCtx.Message().(type) {
	case *dto.GenerateStorageChallengeRequest:
		s.GenerateStorageChallenges(ctx, msg)
	case *dto.StorageChallengeRequest:
		s.ProcessStorageChallenge(ctx, msg)
	case *dto.VerifyStorageChallengeRequest:
		s.VerifyStorageChallenge(ctx, msg)
	default:
		log.Debugf("Action not handled %#v", msg)
		// TODO: response with unhandled notice
	}
}

func (s *storageChallengeActor) GenerateStorageChallenges(ctx context.Context, req *dto.GenerateStorageChallengeRequest) (resp *dto.GenerateStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("actor", "GenerateStorageChallenges").Debug("handled generate storage challenge action")
	// validate request body
	es := validateGenerateStorageChallengeData(req)
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("actor", "GenerateStorageChallenges").Errorf("[GenerateStorageChallenge][Validation Error] %s", es)
		return &dto.GenerateStorageChallengeReply{}, err
	}

	// calling domain service to process bussiness logics
	err = s.service.GenerateStorageChallenges(ctx, req.GetCurrentBlockHash(), req.GetChallengingMasternodeId(), int(req.GetChallengesPerMasternodePerBlock()))
	return &dto.GenerateStorageChallengeReply{}, err
}

func (s *storageChallengeActor) ProcessStorageChallenge(ctx context.Context, req *dto.StorageChallengeRequest) (resp *dto.StorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("actor", "ProcessStorageChallenge").Debug("handled process storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("actor", "ProcessStorageChallenge").Errorf("[ProcessStorageChallenge][Validation Error] %s", es)
		return &dto.StorageChallengeReply{Data: req.GetData()}, err
	}

	// calling domain service to process bussiness logics
	err = s.service.ProcessStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	return &dto.StorageChallengeReply{Data: req.GetData()}, err
}

func (s *storageChallengeActor) VerifyStorageChallenge(ctx context.Context, req *dto.VerifyStorageChallengeRequest) (resp *dto.VerifyStorageChallengeReply, err error) {
	log.WithContext(ctx).WithField("actor", "VerifyStorageChallenge").Debug("handled verify storage challenge action")
	// validate request body
	es := validateStorageChallengeData(req.GetData(), "Data")
	if err := validationErrorStackWrap(es); err != nil {
		log.WithContext(ctx).WithField("actor", "VerifyStorageChallenge").Errorf("[VerifyStorageChallenge][Validation Error] %s", es)
		return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
	}
	// calling domain service to process bussiness logics
	err = s.service.VerifyStorageChallenge(ctx, mapChallengeMessage(req.GetData()))
	return &dto.VerifyStorageChallengeReply{Data: req.GetData()}, err
}

func newStorageChallengeActor(domainService storagechallenge.StorageChallenge) actor.Actor {
	return &storageChallengeActor{service: domainService}
}

func RegisterStorageChallengeApplicationActor(remoter messaging.Remoter, domainService storagechallenge.StorageChallenge) (*actor.PID, error) {
	return remoter.RegisterActor(newStorageChallengeActor(domainService), "STORAGE_CHALLENGE")
}
