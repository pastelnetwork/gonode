package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	"github.com/pastelnetwork/gonode/supernode/services/storagechallenge"
)

type containsAppContext interface {
	getAppContext() context.Context
}

type generateStorageChallenge struct {
	*messaging.CommonProtoMsg
	Context                         context.Context
	ChallengingMasternodeID         string
	ChallengesPerMasternodePerBlock int32
	CurrentBlockHash                string
}

func newGenerateStorageChallengeMsg(ctx context.Context, currentBlockHash, challengingMasternodeID string, challengesPerMasternodePerBlock int32) *generateStorageChallenge {
	return &generateStorageChallenge{
		CommonProtoMsg:                  &messaging.CommonProtoMsg{},
		Context:                         ctx,
		ChallengingMasternodeID:         challengingMasternodeID,
		CurrentBlockHash:                currentBlockHash,
		ChallengesPerMasternodePerBlock: challengesPerMasternodePerBlock,
	}
}

func (msg *generateStorageChallenge) getAppContext() context.Context {
	if msg.Context == nil {
		return context.Background()
	}
	return msg.Context
}

type processStorageChallenge struct {
	*messaging.CommonProtoMsg
	Context context.Context
	*storagechallenge.ChallengeMessage
}

func (msg *processStorageChallenge) getAppContext() context.Context {
	if msg.Context == nil {
		return context.Background()
	}
	return msg.Context
}

func newProcessStorageChallengeMsg(ctx context.Context, challengeMsg *storagechallenge.ChallengeMessage) *processStorageChallenge {
	return &processStorageChallenge{
		CommonProtoMsg:   &messaging.CommonProtoMsg{},
		Context:          ctx,
		ChallengeMessage: challengeMsg,
	}
}

type verifyStorageChallenge struct {
	*messaging.CommonProtoMsg
	Context context.Context
	*storagechallenge.ChallengeMessage
}

func (msg *verifyStorageChallenge) getAppContext() context.Context {
	if msg.Context == nil {
		return context.Background()
	}
	return msg.Context
}

func newVerifyStorageChallengeMsg(ctx context.Context, challengeMsg *storagechallenge.ChallengeMessage) *verifyStorageChallenge {
	return &verifyStorageChallenge{
		CommonProtoMsg:   &messaging.CommonProtoMsg{},
		Context:          ctx,
		ChallengeMessage: challengeMsg,
	}
}

type applicationActor struct {
	domainService storagechallenge.StorageChallenge
	logger        *log.Logger
}

func (a *applicationActor) Receive(ctx actor.Context) {
	var err error
	var appCtx context.Context
	msgInterface := ctx.Message()
	withAppCtx, ok := msgInterface.(containsAppContext)
	if ok {
		appCtx = withAppCtx.getAppContext()
	} else {
		appCtx = context.Background()
	}
	appCtx = appCtx.WithActorContext(ctx)

	logger := a.logger.WithContext(appCtx)

	switch msg := ctx.Message().(type) {
	case *generateStorageChallenge:
		logger.Debugf("receive %#v action", msg)
		err = a.domainService.GenerateStorageChallenges(appCtx, msg.CurrentBlockHash, msg.ChallengingMasternodeID, int(msg.ChallengesPerMasternodePerBlock))
	case *processStorageChallenge:
		logger.Debugf("receive %#v action", msg)
		err = a.domainService.ProcessStorageChallenge(appCtx, msg.ChallengeMessage)
	case *verifyStorageChallenge:
		logger.Debugf("receive %#v action", msg)
		// calling domain service to process bussiness logics
		err = a.domainService.VerifyStorageChallenge(appCtx, msg.ChallengeMessage)
	default:
		logger.Debugf("Unhandled action %#v", msg)
	}

	if err != nil {
		logger.WithError(err).Errorf("action %#v failed", msgInterface)
	}
}
