package storagechallenge

import (
	"fmt"
	"time"

	appcontext "github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
)

func (s *service) VerifyStorageChallenge(ctx appcontext.Context, incomingChallengeMessage *ChallengeMessage) error {
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug(incomingChallengeMessage.MessageType)
	if err := s.validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	// TODO: replace with new repository implementation
	// filePath, err := s.repository.GetFilePathFromFileHash(ctx, incomingChallengeMessage.FileHashToChallenge)
	// if err != nil {
	// 	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Error("could not get symbol file path from file hash", actorLog.String("s.repository.GetFilePathFromFileHash", err.Error()))
	// 	return err
	// }

	// TODO: replace with read file data from kademlia
	// challengeFileData, err := file.ReadFileIntoMemory(filePath)
	// if err != nil {
	// 	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Error("could not read file data in to memory", "file.ReadFileIntoMemory", err.Error()))
	// 	return err
	// }
	var challengeFileData = []byte("")

	challengeCorrectHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	messageType := MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	TimestampChallengeVerified := time.Now().Unix()
	TimeVerifyStorageChallengeInMilliSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, TimestampChallengeVerified)
	var challengeStatus string
	// var analysisStatus = ALALYSIS_STATUS_TIMEOUT
	// defer func() {
	// 	s.saveChallengeAnalysis(ctx, incomingChallengeMessage.BlockHashWhenChallengeSent, incomingChallengeMessage.ChallengingMasternodeID, analysisStatus)
	// }()

	if (incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash) && (TimeVerifyStorageChallengeInMilliSeconds <= float64(s.storageChallengeExpiredAsMilliSeconds)) {
		challengeStatus = Status_SUCCEEDED
		// analysisStatus = ANALYSIS_STATUS_CORRECT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug("masternode %s correctly responded in %d milliseconds to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, TimeVerifyStorageChallengeInMilliSeconds, incomingChallengeMessage.FileHashToChallenge)
	} else if incomingChallengeMessage.ChallengeResponseHash == challengeCorrectHash {
		challengeStatus = Status_FAILED_TIMEOUT
		// analysisStatus = ALALYSIS_STATUS_TIMEOUT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug("masternode %s  correctly responded in %d milliseconds to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingChallengeMessage.RespondingMasternodeID, TimeVerifyStorageChallengeInMilliSeconds, incomingChallengeMessage.FileHashToChallenge)
	} else {
		challengeStatus = Status_FAILED_INCORRECT_RESPONSE
		// analysisStatus = ALALYSIS_STATUS_INCORRECT
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug("masternode %s failed by incorrectly responding to a storage challenge for file %s", incomingChallengeMessage.RespondingMasternodeID, incomingChallengeMessage.FileHashToChallenge)
	}

	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)

	var outgoingChallengeMessage = &ChallengeMessage{
		MessageID:                     messageID,
		MessageType:                   MessageType_STORAGE_CHALLENGE_VERIFICATION_MESSAGE,
		ChallengeStatus:               challengeStatus,
		TimestampChallengeSent:        incomingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: incomingChallengeMessage.TimestampChallengeRespondedTo,
		TimestampChallengeVerified:    TimestampChallengeVerified,
		BlockHashWhenChallengeSent:    incomingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeID:       incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:        incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:           incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     challengeCorrectHash,
		ChallengeResponseHash:         incomingChallengeMessage.ChallengeResponseHash,
		ChallengeID:                   incomingChallengeMessage.ChallengeID,
	}

	// TODO: replace with new repository implementation
	// if err := s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage); err != nil {
	// 	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Error("could not update new storage challenge message in to database", actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
	// 	return err
	// }

	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, outgoingChallengeMessage.TimestampChallengeRespondedTo)
	log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").Debug("masternode " + outgoingChallengeMessage.RespondingMasternodeID + " responded to storage challenge for file hash " + outgoingChallengeMessage.FileHashToChallenge + " in " + fmt.Sprint(timeToRespondToStorageChallengeInSeconds) + " seconds!")

	return nil
}

func (s *service) validateVerifyingStorageChallengeIncommingData(incomingChallengeMessage *ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != Status_RESPONDED {
		return fmt.Errorf("incorrect status to verify storage challenge")
	}
	if incomingChallengeMessage.MessageType != MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE {
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
