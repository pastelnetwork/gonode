package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	dto "github.com/pastelnetwork/gonode/proto/supernode/storagechallenge"
	"github.com/pastelnetwork/gonode/supernode/node"
)

type sendVerifyStorageChallengeMsg struct {
	*messaging.CommonProtoMsg
	Context                       context.Context
	VerifierMasternodesClientPIDs []*actor.PID
	*ChallengeMessage
}

func newSendVerifyStorageChallengeMsg(ctx context.Context, verifierMasternodesClientPIDs []*actor.PID, challengeMsg *ChallengeMessage) *sendVerifyStorageChallengeMsg {
	return &sendVerifyStorageChallengeMsg{
		CommonProtoMsg:                &messaging.CommonProtoMsg{},
		Context:                       ctx,
		VerifierMasternodesClientPIDs: verifierMasternodesClientPIDs,
		ChallengeMessage:              challengeMsg,
	}
}

type sendProcessStorageChallengeMsg struct {
	*messaging.CommonProtoMsg
	Context                       context.Context
	ProcessingMasternodeClientPID *actor.PID
	*ChallengeMessage
}

func newSendProcessStorageChallengeMsg(ctx context.Context, processingMasternodeClientPID *actor.PID, challengeMsg *ChallengeMessage) *sendProcessStorageChallengeMsg {
	return &sendProcessStorageChallengeMsg{
		CommonProtoMsg:                &messaging.CommonProtoMsg{},
		Context:                       ctx,
		ProcessingMasternodeClientPID: processingMasternodeClientPID,
		ChallengeMessage:              challengeMsg,
	}
}

type domainActor struct {
	conn node.Client
}

func (d *domainActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *sendVerifyStorageChallengeMsg:
		d.OnSendVerifyStorageChallengeMessage(context, msg)
	case *sendProcessStorageChallengeMsg:
		d.OnSendProcessStorageChallengeMessage(context, msg)
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
func (d *domainActor) OnSendVerifyStorageChallengeMessage(ctx actor.Context, msg *sendVerifyStorageChallengeMsg) {
	for _, verifyingMasternodePID := range msg.VerifierMasternodesClientPIDs {
		log.Debug(verifyingMasternodePID.String())
		storageChallengeReq := &dto.VerifyStorageChallengeRequest{
			Data: &dto.StorageChallengeData{
				MessageId:                     msg.MessageID,
				MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
				ChallengeStatus:               dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
				TimestampChallengeSent:        msg.TimestampChallengeSent,
				TimestampChallengeRespondedTo: msg.TimestampChallengeRespondedTo,
				TimestampChallengeVerified:    0,
				BlockHashWhenChallengeSent:    msg.BlockHashWhenChallengeSent,
				ChallengingMasternodeId:       msg.ChallengingMasternodeID,
				RespondingMasternodeId:        msg.RespondingMasternodeID,
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

		// TODO: replace lines below with actor context
		bgContext := context.Background()
		conn, err := d.conn.Connect(bgContext, verifyingMasternodePID.GetAddress())
		if err != nil {
			return
		}
		conn.StorageChallenge().VerifyStorageChallenge(bgContext, storageChallengeReq)
		// ctx.Send(verifyingMasternodePID, storageChallengeReq)
	}
}

// OnSendProcessStorageChallengeMessage handle event sending processing stotage challenge message
func (d *domainActor) OnSendProcessStorageChallengeMessage(ctx actor.Context, msg *sendProcessStorageChallengeMsg) {
	log.Debug(msg.ProcessingMasternodeClientPID.String())
	storageChallengeReq := &dto.ProcessStorageChallengeRequest{
		Data: &dto.StorageChallengeData{
			MessageId:                     msg.MessageID,
			MessageType:                   dto.StorageChallengeDataMessageType(dto.StorageChallengeDataMessageType_value[msg.MessageType]),
			ChallengeStatus:               dto.StorageChallengeDataStatus(dto.StorageChallengeDataStatus_value[msg.ChallengeStatus]),
			TimestampChallengeSent:        msg.TimestampChallengeSent,
			TimestampChallengeRespondedTo: 0,
			TimestampChallengeVerified:    0,
			BlockHashWhenChallengeSent:    msg.BlockHashWhenChallengeSent,
			ChallengingMasternodeId:       msg.ChallengingMasternodeID,
			RespondingMasternodeId:        msg.RespondingMasternodeID,
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
	// TODO: replace lines below with actor context
	bgContext := context.Background()
	conn, err := d.conn.Connect(bgContext, msg.ProcessingMasternodeClientPID.GetAddress())
	if err != nil {
		return
	}
	conn.StorageChallenge().ProcessStorageChallenge(bgContext, storageChallengeReq)
	// ctx.Send(msg.ProcessingMasterNodesClientPID, storageChallengeReq)
}
