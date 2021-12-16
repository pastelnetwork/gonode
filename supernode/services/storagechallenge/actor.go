package storagechallenge

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/messaging"
	dto "github.com/pastelnetwork/gonode/proto/supernode"
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
	for _, verifyingMasternodePID := range msg.VerifierMasternodesClientPIDs {
		log.Debug(verifyingMasternodePID.String())
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

		conn, err := d.conn.Connect(msg.Context, verifyingMasternodePID.GetAddress())
		if err != nil {
			return
		}
		conn.StorageChallenge().VerifyStorageChallenge(msg.Context, storageChallengeReq)
		// ctx.Send(verifyingMasternodePID, storageChallengeReq)
	}
}

// OnSendProcessStorageChallengeMessage handle event sending processing stotage challenge message
func (d *domainActor) OnSendProcessStorageChallengeMessage(msg *sendProcessStorageChallengeMsg) {
	log.Debug(msg.ProcessingMasternodeClientPID.String())
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
	conn, err := d.conn.Connect(msg.Context, msg.ProcessingMasternodeClientPID.GetAddress())
	if err != nil {
		return
	}
	conn.StorageChallenge().ProcessStorageChallenge(msg.Context, storageChallengeReq)
	// ctx.Send(msg.ProcessingMasterNodesClientPID, storageChallengeReq)
}
