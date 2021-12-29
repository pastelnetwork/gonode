package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	dto "github.com/pastelnetwork/gonode/proto/supernode"
	"github.com/pastelnetwork/gonode/supernode/node"
	"google.golang.org/grpc/metadata"
)

type sendVerifyStorageChallengeMsg struct {
	*messaging.CommonProtoMsg
	Context                 context.Context
	VerifierMasternodesAddr []string
	*ChallengeMessage
}

func (s *sendVerifyStorageChallengeMsg) getAppContext() context.Context {
	if s.Context == nil {
		return context.Background()
	}

	if md, ok := metadata.FromIncomingContext(s.Context); ok {
		return context.FromContext(metadata.NewOutgoingContext(context.Background(), md))
	}
	return context.Background()
}

func newSendVerifyStorageChallengeMsg(ctx context.Context, verifierMasternodesAddr []string, challengeMsg *ChallengeMessage) *sendVerifyStorageChallengeMsg {
	return &sendVerifyStorageChallengeMsg{
		CommonProtoMsg:          &messaging.CommonProtoMsg{},
		Context:                 ctx,
		VerifierMasternodesAddr: verifierMasternodesAddr,
		ChallengeMessage:        challengeMsg,
	}
}

type sendProcessStorageChallengeMsg struct {
	*messaging.CommonProtoMsg
	Context                  context.Context
	ProcessingMasternodeAddr string
	*ChallengeMessage
}

func (s *sendProcessStorageChallengeMsg) getAppContext() context.Context {
	if s.Context == nil {
		return context.Background()
	}

	if md, ok := metadata.FromIncomingContext(s.Context); ok {
		return context.FromContext(metadata.NewOutgoingContext(context.Background(), md))
	}
	return context.Background()
}

func newSendProcessStorageChallengeMsg(ctx context.Context, processingMasternodeAddr string, challengeMsg *ChallengeMessage) *sendProcessStorageChallengeMsg {
	return &sendProcessStorageChallengeMsg{
		CommonProtoMsg:           &messaging.CommonProtoMsg{},
		Context:                  ctx,
		ProcessingMasternodeAddr: processingMasternodeAddr,
		ChallengeMessage:         challengeMsg,
	}
}

type domainActor struct {
	conn node.Client
}

func (d *domainActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *sendVerifyStorageChallengeMsg:
		d.OnSendVerifyStorageChallengeMessage(msg)
	case *sendProcessStorageChallengeMsg:
		d.OnSendProcessStorageChallengeMessage(msg)
	default:
		log.WithField("actor", "domain actor").Debug("Unhandled action", msg)
	}
}

func newDomainActor(conn node.Client) actor.Actor {
	return &domainActor{
		conn: conn,
	}
}

// OnSendVerifyStorageChallengeMessage handle event sending verity storage challenge message
func (d *domainActor) OnSendVerifyStorageChallengeMessage(msg *sendVerifyStorageChallengeMsg) {
	appCtx := msg.getAppContext()
	for _, verifyingMasternodeAddr := range msg.VerifierMasternodesAddr {
		log.Debug(verifyingMasternodeAddr)
		storageChallengeReq := &dto.VerifyStorageChallengeRequest{
			Data: &dto.StorageChallengeData{
				MessageId:                    msg.MessageID,
				MessageType:                  dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
				ChallengeStatus:              dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
				BlockNumChallengeSent:        msg.BlockNumChallengeSent,
				BlockNumChallengeRespondedTo: msg.BlockNumChallengeRespondedTo,
				BlockNumChallengeVerified:    0,
				MerklerootWhenChallengeSent:  msg.MerklerootWhenChallengeSent,
				ChallengingMasternodeId:      msg.ChallengingMasternodeID,
				RespondingMasternodeId:       msg.RespondingMasternodeID,
				ChallengeFile: &dto.StorageChallengeDataChallengeFile{
					FileHashToChallenge:      msg.FileHashToChallenge,
					ChallengeSliceStartIndex: int64(msg.ChallengeSliceStartIndex),
					ChallengeSliceEndIndex:   int64(msg.ChallengeSliceEndIndex),
				},
				ChallengeSliceCorrectHash: "",
				ChallengeResponseHash:     msg.ChallengeResponseHash,
				ChallengeId:               msg.ChallengeID,
			},
		}

		// connect to verifying masternode by gRPC
		conn, err := d.conn.Connect(msg.getAppContext(), verifyingMasternodeAddr)
		if err != nil {
			log.WithContext(appCtx).WithError(err).Errorf("could not connect to verifying masternode %s", verifyingMasternodeAddr)
			continue
		}

		// send gRPC response hash to verifying masternode
		if _, err = conn.StorageChallenge().VerifyStorageChallenge(msg.getAppContext(), storageChallengeReq); err != nil {
			log.WithContext(appCtx).WithError(err).Errorf("could send verify storage challenge request to verifying masternode %s", verifyingMasternodeAddr)
		}
	}
}

// OnSendProcessStorageChallengeMessage handle event sending processing stotage challenge message
func (d *domainActor) OnSendProcessStorageChallengeMessage(msg *sendProcessStorageChallengeMsg) {
	log.Debug(msg.ProcessingMasternodeAddr)
	storageChallengeReq := &dto.ProcessStorageChallengeRequest{
		Data: &dto.StorageChallengeData{
			MessageId:                    msg.MessageID,
			MessageType:                  dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
			ChallengeStatus:              dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
			BlockNumChallengeSent:        msg.BlockNumChallengeSent,
			BlockNumChallengeRespondedTo: 0,
			BlockNumChallengeVerified:    0,
			MerklerootWhenChallengeSent:  msg.MerklerootWhenChallengeSent,
			ChallengingMasternodeId:      msg.ChallengingMasternodeID,
			RespondingMasternodeId:       msg.RespondingMasternodeID,
			ChallengeFile: &dto.StorageChallengeDataChallengeFile{
				FileHashToChallenge:      msg.FileHashToChallenge,
				ChallengeSliceStartIndex: int64(msg.ChallengeSliceStartIndex),
				ChallengeSliceEndIndex:   int64(msg.ChallengeSliceEndIndex),
			},
			ChallengeSliceCorrectHash: "",
			ChallengeResponseHash:     "",
			ChallengeId:               msg.ChallengeID,
		},
	}
	appCtx := msg.getAppContext()
	// connect to processing masternode by gRPC
	conn, err := d.conn.Connect(appCtx, msg.ProcessingMasternodeAddr)
	if err != nil {
		log.WithContext(appCtx).WithError(err).Errorf("could not connect to challenging masternode %s", msg.ProcessingMasternodeAddr)
		return
	}
	// send gRPC process challenge request
	if _, err = conn.StorageChallenge().ProcessStorageChallenge(appCtx, storageChallengeReq); err != nil {
		log.WithContext(appCtx).WithError(err).Errorf("could send process storage challenge request to challenging masternode %s", msg.ProcessingMasternodeAddr)
	}
}
