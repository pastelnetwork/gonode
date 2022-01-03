package storagechallenge

import (
	"fmt"

	appcontext "github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
)

func (s *service) VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *ChallengeMessage) error {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug(incomingChallengeMessage.MessageType)
	// incomming challenge message validation
	if err := s.validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	/* ----------------------------------------------- */
	/* ----- Main logic implementation goes here ----- */
	challengeFileData, err := s.repository.GetSymbolFileByKey(ctx, incomingChallengeMessage.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not read file data in to memory", "file.ReadFileIntoMemory", err.Error())
		return err
	}

	challengeCorrectHash := s.computeHashOfFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	messageType := storageChallengeVerificationMessage
	blockNumChallengeVerified, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not get current block count")
		return err
	}

	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - incomingChallengeMessage.BlockNumChallengeSent
	var challengeStatus string
	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (blocksVerifyStorageChallengeInBlocks <= s.storageChallengeExpiredBlocks) {
		challengeStatus = statusSucceeded
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug(fmt.Sprintf("masternode %s correctly responded in %d blocks to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, blocksVerifyStorageChallengeInBlocks, incomingChallengeMessage.FileHashToChallenge))
		s.repository.SaveChallengMessageState(ctx, "succeeded", incomingChallengeMessage.ChallengeID, incomingChallengeMessage.ChallengingMasternodeID, incomingChallengeMessage.BlockNumChallengeSent)
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {
		challengeStatus = statusFailedTimeout
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug(fmt.Sprintf("masternode %s  correctly responded in %d blocks to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingChallengeMessage.RespondingMasternodeID, blocksVerifyStorageChallengeInBlocks, incomingChallengeMessage.FileHashToChallenge))
		s.repository.SaveChallengMessageState(ctx, "timeout", incomingChallengeMessage.ChallengeID, incomingChallengeMessage.ChallengingMasternodeID, incomingChallengeMessage.BlockNumChallengeSent)
	} else {
		challengeStatus = statusFailedIncorrectResponse
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug(fmt.Sprintf("masternode %s failed by incorrectly responding to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, incomingChallengeMessage.FileHashToChallenge))
		s.repository.SaveChallengMessageState(ctx, "failed", incomingChallengeMessage.ChallengeID, incomingChallengeMessage.ChallengingMasternodeID, incomingChallengeMessage.BlockNumChallengeSent)
	}

	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)

	var outgoingChallengeMessage = &ChallengeMessage{
		MessageID:                    messageID,
		MessageType:                  messageType,
		ChallengeStatus:              challengeStatus,
		BlockNumChallengeSent:        incomingChallengeMessage.BlockNumChallengeSent,
		BlockNumChallengeRespondedTo: incomingChallengeMessage.BlockNumChallengeRespondedTo,
		BlockNumChallengeVerified:    blockNumChallengeVerified,
		MerklerootWhenChallengeSent:  incomingChallengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeID:      incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:       incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:          incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:     incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:       incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:    challengeCorrectHash,
		ChallengeResponseHash:        incomingChallengeMessage.ChallengeResponseHash,
		ChallengeID:                  incomingChallengeMessage.ChallengeID,
	}

	blocksToRespondToStorageChallengeInNanoseconds := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeID).Debug("masternode " + outgoingChallengeMessage.RespondingMasternodeID + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallengeInNanoseconds) + " nano seconds!")

	return nil
}

func (s *service) validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage *ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != statusResponded {
		return fmt.Errorf("incorrect status to verify storage challenge")
	}
	if incomingChallengeMessage.MessageType != storageChallengeResponseMessage {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}
	return nil
}