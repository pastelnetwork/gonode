package storagechallenge

import (
	"encoding/hex"
	"fmt"

	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"golang.org/x/crypto/sha3"
)

func (s *service) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error {
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug("Start processing storage challenge")

	// incomming challenge message validation
	if err := s.validateProcessingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	/* ----------------------------------------------- */
	/* ----- Main logic implementation goes here ----- */
	challengeFileData, err := s.repository.GetSymbolFileByKey(ctx, incomingChallengeMessage.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not read file data in to memory")
		return err
	}
	challengeResponseHash := s.computeHashOfFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	challengeStatus := statusResponded
	messageType := storageChallengeResponseMessage
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)
	blockNumChallengeRespondedTo, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not get current block count")
		return err
	}

	var outgoingChallengeMessage = &ChallengeMessage{
		MessageID:                    messageID,
		MessageType:                  messageType,
		ChallengeStatus:              challengeStatus,
		BlockNumChallengeSent:        incomingChallengeMessage.BlockNumChallengeSent,
		BlockNumChallengeRespondedTo: blockNumChallengeRespondedTo,
		BlockNumChallengeVerified:    0,
		MerklerootWhenChallengeSent:  incomingChallengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeID:      incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:       incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:          incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:     incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:       incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:    "",
		ChallengeResponseHash:        challengeResponseHash,
		ChallengeID:                  incomingChallengeMessage.ChallengeID,
	}

	blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug(fmt.Sprintf("masternode %s responded to storage challenge for file hash %s in %v nano second!", outgoingChallengeMessage.RespondingMasternodeID, outgoingChallengeMessage.FileHashToChallenge, blocksToRespondToStorageChallenge))

	// send to verifying masternode to validate challenge response hash
	if err = s.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
		return err
	}
	s.repository.SaveChallengMessageState(ctx, "respond", outgoingChallengeMessage.ChallengeID, outgoingChallengeMessage.ChallengingMasternodeID, outgoingChallengeMessage.BlockNumChallengeSent)
	return nil
}

func (s *service) validateProcessingStorageChallengeIncommingData(incomingChallengeMessage *ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != statusPending {
		return fmt.Errorf("incorrect status to processing storage challenge")
	}
	if incomingChallengeMessage.MessageType != storageChallengeIssuanceMessage {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}
	return nil
}

func (s *service) computeHashOfFileSlice(fileData []byte, challengeSliceStartIndex int, challengeSliceEndIndex int) string {
	challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
	algorithm := sha3.New256()
	algorithm.Write(challengeDataSlice)
	return hex.EncodeToString(algorithm.Sum(nil))
}

func (s *service) sendVerifyStorageChallenge(ctx context.Context, challengeMessage *ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	mapMasternodes := make(map[string]pastel.MasterNode)
	for _, mn := range masternodes {
		mapMasternodes[mn.ExtKey] = mn
	}

	verifierMasternodesAddr := []string{}
	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapMasternodes[challengeMessage.RespondingMasternodeID]; !ok {
		return fmt.Errorf("cannot get masternode info of masternode id %v", challengeMessage.RespondingMasternodeID)
	}
	verifierMasternodesAddr = append(verifierMasternodesAddr, mn.ExtAddress)

	return s.actor.Send(ctx, s.domainActorID, newSendVerifyStorageChallengeMsg(ctx, verifierMasternodesAddr, challengeMessage))
}
