package storagechallenge

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
	"golang.org/x/crypto/sha3"
)

func (s *service) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *ChallengeMessage) error {
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").Debug("Start processing storage challenge")
	if err := s.validateProcessingStorageChallengeIncommingData(incomingChallengeMessage); err != nil {
		return err
	}

	// analysisStatus := ALALYSIS_STATUS_TIMEOUT

	// defer func() {
	// 	s.saveChallengeAnalysis(ctx, incomingChallengeMessage.BlockHashWhenChallengeSent, incomingChallengeMessage.ChallengingMasternodeID, analysisStatus)
	// }()

	// TODO: replace with new repository implementation
	// filePath, err := s.repository.GetFilePathFromFileHash(ctx, incomingChallengeMessage.FileHashToChallenge)
	// if err != nil {
	// 	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not get symbol file path from file hash", actorLog.String("s.repository.GetFilePathFromFileHash", err.Error()))
	// 	return err
	// }

	// TODO: replace with kademlia storage
	// challengeFileData, err := file.ReadFileIntoMemory(filePath)
	// if err != nil {
	// 	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not read file data in to memory", actorLog.String("file.ReadFileIntoMemory", err.Error()))
	// 	return err
	// }
	var challengeFileData = []byte("")
	challengeResponseHash := s.computeHashofFileSlice(challengeFileData, int(incomingChallengeMessage.ChallengeSliceStartIndex), int(incomingChallengeMessage.ChallengeSliceEndIndex))
	challengeStatus := Status_RESPONDED
	messageType := MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeID + incomingChallengeMessage.RespondingMasternodeID + incomingChallengeMessage.FileHashToChallenge + challengeStatus + messageType + incomingChallengeMessage.BlockHashWhenChallengeSent
	messageID := helper.GetHashFromString(messageIDInputData)
	timestampChallengeRespondedTo := time.Now().UnixMilli()

	var outgoingChallengeMessage = &ChallengeMessage{
		MessageID:                     messageID,
		MessageType:                   MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE,
		ChallengeStatus:               Status_RESPONDED,
		TimestampChallengeSent:        incomingChallengeMessage.TimestampChallengeSent,
		TimestampChallengeRespondedTo: timestampChallengeRespondedTo,
		TimestampChallengeVerified:    0,
		BlockHashWhenChallengeSent:    incomingChallengeMessage.BlockHashWhenChallengeSent,
		ChallengingMasternodeID:       incomingChallengeMessage.ChallengingMasternodeID,
		RespondingMasternodeID:        incomingChallengeMessage.RespondingMasternodeID,
		FileHashToChallenge:           incomingChallengeMessage.FileHashToChallenge,
		ChallengeSliceStartIndex:      incomingChallengeMessage.ChallengeSliceStartIndex,
		ChallengeSliceEndIndex:        incomingChallengeMessage.ChallengeSliceEndIndex,
		ChallengeSliceCorrectHash:     "",
		ChallengeResponseHash:         challengeResponseHash,
		ChallengeID:                   incomingChallengeMessage.ChallengeID,
	}

	// TODO: replace with new repository implementation
	// if err := s.repository.UpsertStorageChallengeMessage(ctx, outgoingChallengeMessage); err != nil {
	// 	log.With(actorLog.String("ACTOR", "ProcessStorageChallenge")).Error("could not update new storage challenge message in to database", actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
	// 	return err
	// }
	// analysisStatus = ANALYSIS_STATUS_RESPONDED_TO
	timeToRespondToStorageChallengeInSeconds := helper.ComputeElapsedTimeInSecondsBetweenTwoDatetimes(incomingChallengeMessage.TimestampChallengeSent, outgoingChallengeMessage.TimestampChallengeRespondedTo)
	log.WithField("method", "ProcessStorageChallenge").Debug("masternode %s responded to storage challenge for file hash %s in %d milli second!", outgoingChallengeMessage.RespondingMasternodeID, outgoingChallengeMessage.FileHashToChallenge, fmt.Sprint(timeToRespondToStorageChallengeInSeconds)+" seconds!")

	return s.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage)
}

func (s *service) validateProcessingStorageChallengeIncommingData(incomingChallengeMessage *ChallengeMessage) error {
	if incomingChallengeMessage.ChallengeStatus != Status_PENDING {
		return fmt.Errorf("incorrect status to processing storage challenge")
	}
	if incomingChallengeMessage.MessageType != MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}
	return nil
}

func (s *service) computeHashofFileSlice(file_data []byte, challenge_slice_start_index int, challenge_slice_end_index int) string {
	challenge_data_slice := file_data[challenge_slice_start_index:challenge_slice_end_index]
	algorithm := sha3.New256()
	algorithm.Write(challenge_data_slice)
	hash_of_data_slice := hex.EncodeToString(algorithm.Sum(nil))
	return hash_of_data_slice
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

	verifierMasterNodesClientPIDs := []*actor.PID{}
	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapMasternodes[challengeMessage.RespondingMasternodeID]; !ok {
		return fmt.Errorf("cannot get masternode info of masternode id %v", challengeMessage.RespondingMasternodeID)
	}
	verifierMasterNodesClientPIDs = append(verifierMasterNodesClientPIDs, actor.NewPID(mn.ExtAddress, "storage-challenge"))

	return s.remoter.Send(ctx, s.domainActorID, &verifyStorageChallengeMsg{VerifierMasterNodesClientPIDs: verifierMasterNodesClientPIDs, ChallengeMessage: challengeMessage})
}
