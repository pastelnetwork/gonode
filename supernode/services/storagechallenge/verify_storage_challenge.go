package storagechallenge

import (
	"fmt"

	appcontext "github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
)

func (s *service) VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *ChallengeMessage) error {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug(incomingChallengeMessage.MessageType)
	if err := s.validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	challengeFileData, err := s.repository.GetSymbolFileByKey(ctx, incomingChallengeMessage.FileHashToChallenge)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Error("could not read file data in to memory", "file.ReadFileIntoMemory", err.Error())
		return err
	}

	challengeCorrectHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	messageType := storageChallengeVerificationMessage
	blockNumChallengeVerified, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}

	blocksVerifyStorageChallengeInNanoseconds := blockNumChallengeVerified - incomingChallengeMessage.BlockNumChallengeSent
	var challengeStatus string
	// var analysisStatus = ALALYSIS_STATUS_TIMEOUT
	// defer func() {
	// 	s.saveChallengeAnalysis(ctx, incomingChallengeMessage.MerklerootWhenChallengeSent, incomingChallengeMessage.ChallengingMasternodeID, analysisStatus)
	// }()

	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (blocksVerifyStorageChallengeInNanoseconds <= s.storageChallengeExpiredBlocks) {
		challengeStatus = statusSucceeded
		// analysisStatus = ANALYSIS_STATUS_CORRECT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug(fmt.Sprintf("masternode %s correctly responded in %d milliseconds to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, blocksVerifyStorageChallengeInNanoseconds, incomingChallengeMessage.FileHashToChallenge))
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {
		challengeStatus = statusFailedTimeout
		// analysisStatus = ALALYSIS_STATUS_TIMEOUT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug(fmt.Sprintf("masternode %s  correctly responded in %d milliseconds to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingChallengeMessage.RespondingMasternodeID, blocksVerifyStorageChallengeInNanoseconds, incomingChallengeMessage.FileHashToChallenge))
	} else {
		challengeStatus = statusFailedIncorrectResponse
		// analysisStatus = ALALYSIS_STATUS_INCORRECT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug(fmt.Sprintf("masternode %s failed by incorrectly responding to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, incomingChallengeMessage.FileHashToChallenge))
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

	// TODO: replace with new repository implementation
	// if err := s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage); err != nil {
	// 	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Error("could not update new storage challenge message in to database", actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
	// 	return err
	// }

	blocksToRespondToStorageChallengeInNanoseconds := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug("masternode " + outgoingChallengeMessage.RespondingMasternodeID + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallengeInNanoseconds) + " nano seconds!")

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

// func (s *service) saveChallengeAnalysis(ctx appcontext.Context, blockHash, challengingMasternodeID string, challengeAnalysisStatus ChallengeAnalysisStatus) error {
// 	switch challengeAnalysisStatus {
// 	case ANALYSYS_STATUS_ISSUED:
// 		s.repository.IncreaseMasternodeTotalChallengesIssued(ctx, challengingMasternodeID)
// 		s.repository.IncreasePastelBlockTotalChallengesIssued(ctx, blockHash)
// 	case ANALYSIS_STATUS_RESPONDED_TO:
// 		s.repository.IncreaseMasternodeTotalChallengesRespondedTo(ctx, challengingMasternodeID)
// 		s.repository.IncreasePastelBlockTotalChallengesRespondedTo(ctx, blockHash)
// 	case ANALYSIS_STATUS_CORRECT:
// 		s.repository.IncreaseMasternodeTotalChallengesCorrect(ctx, challengingMasternodeID)
// 		s.repository.IncreasePastelBlockTotalChallengesCorrect(ctx, blockHash)
// 	case ALALYSIS_STATUS_INCORRECT:
// 		s.repository.IncreaseMasternodeTotalChallengesIncorrect(ctx, challengingMasternodeID)
// 		s.repository.IncreasePastelBlockTotalChallengesIncorrect(ctx, blockHash)
// 	case ALALYSIS_STATUS_TIMEOUT:
// 		s.repository.IncreaseMasternodeTotalChallengesTimeout(ctx, challengingMasternodeID)
// 		s.repository.IncreasePastelBlockTotalChallengesTimeout(ctx, blockHash)
// 	}

// 	return nil
// }
