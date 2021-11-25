package storagechallenge

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/utils/helper"
	"github.com/pastelnetwork/storage-challenges/utils/xordistance"
)

func (s *service) GenerateStorageChallenges(ctx context.Context, currentBlockHash string, challengingMasternodeID string, challengesPerMasternodePerBlock int) error {
	symbolFileKeys := s.repository.ListKeys(ctx)

	comparisonStringForFileHashSelection := currentBlockHash + challengingMasternodeID
	sliceOfFileHashesToChallenge := s.repository.GetNClosestXORDistanceFileHashesToComparisonString(ctx, challengesPerMasternodePerBlock, comparisonStringForFileHashSelection, symbolFileKeys)

	for idx, symbolFileHash := range sliceOfFileHashesToChallenge {
		challengeDataSize := 0

		comparisonStringForMasternodeSelection := currentBlockHash + symbolFileHash + s.nodeID + helper.GetHashFromString(fmt.Sprint(idx))
		respondingMasternodes := s.repository.GetNClosestXORDistanceMasternodesToComparisionString(ctx, 1, comparisonStringForMasternodeSelection)
		challengeStatus := statusPending
		messageType := storageChallengeIssuanceMessage
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(uint64(challengeDataSize), symbolFileHash, currentBlockHash, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + string(respondingMasternodes[0].ID) + symbolFileHash + challengeStatus + messageType + currentBlockHash
		messageID := helper.GetHashFromString(messageIDInputData)
		timestampChallengeSent := time.Now().UnixMilli()
		challengeIDInputData := challengingMasternodeID + string(respondingMasternodes[0].ID) + symbolFileHash + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(timestampChallengeSent)
		challengeID := helper.GetHashFromString(challengeIDInputData)
		outgoinChallengeMessage := &ChallengeMessage{
			MessageID:                     messageID,
			MessageType:                   messageType,
			ChallengeStatus:               challengeStatus,
			TimestampChallengeSent:        timestampChallengeSent,
			TimestampChallengeRespondedTo: 0,
			TimestampChallengeVerified:    0,
			BlockHashWhenChallengeSent:    currentBlockHash,
			ChallengingMasternodeID:       challengingMasternodeID,
			RespondingMasternodeID:        string(respondingMasternodes[0].ID),
			FileHashToChallenge:           symbolFileHash,
			ChallengeSliceStartIndex:      uint64(challengeSliceStartIndex),
			ChallengeSliceEndIndex:        uint64(challengeSliceEndIndex),
			ChallengeSliceCorrectHash:     "",
			ChallengeResponseHash:         "",
			ChallengeID:                   challengeID,
		}

		// TODO: replacing this with real repository implementation
		// err = s.repository.UpsertStorageChallengeMessage(ctx, outgoinChallengeMessage)
		// if err != nil {
		// 	log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Warn(fmt.Sprintf("could not update storage challenge into storage: %v", outgoinChallengeMessage), actorLog.String("s.repository.UpsertStorageChallengeMessage", err.Error()))
		// 	continue
		// }
		// s.saveChallengeAnalysis(ctx, outgoinChallengeMessage.BlockHashWhenChallengeSent, outgoinChallengeMessage.ChallengingMasternodeID, ANALYSYS_STATUS_ISSUED)

		s.sendprocessStorageChallenge(ctx, outgoinChallengeMessage)
	}

	return nil
}

func (s *service) sendprocessStorageChallenge(ctx context.Context, challengeMessage *ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("method", "sendprocessStorageChallenge").Warn("could not get masternode info: %v", err.Error())
		return err
	}

	mapMasternodes := make(map[string]pastel.MasterNode)
	for _, mn := range masternodes {
		mapMasternodes[mn.ExtKey] = mn
	}

	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapMasternodes[challengeMessage.ChallengingMasternodeID]; !ok {
		err = fmt.Errorf("cannot get masternode info of masternode id %v", challengeMessage.ChallengingMasternodeID)
		log.WithContext(ctx).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	processingMasterNodesClientPID := actor.NewPID(mn.ExtAddress, "storage-challenge")

	return s.remoter.Send(ctx, s.domainActorID, &processStorageChallengeMsg{ProcessingMasterNodesClientPID: processingMasterNodesClientPID, ChallengeMessage: challengeMessage})
}

func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingMasternodeID string) (int, int) {
	blockHashStringAsInt, _ := strconv.ParseInt(blockHashString, 16, 64)
	blockHashStringAsIntStr := fmt.Sprint(blockHashStringAsInt)
	stepSizeForIndicesStr := blockHashStringAsIntStr[len(blockHashStringAsIntStr)-1:] + blockHashStringAsIntStr[0:1]
	stepSizeForIndices, _ := strconv.ParseUint(stepSizeForIndicesStr, 10, 32)
	stepSizeForIndicesAsInt := int(stepSizeForIndices)
	comparisonString := blockHashString + fileHashString + challengingMasternodeID
	sliceOfXorDistancesOfIndicesToBlockHash := make([]uint64, 0)
	sliceOfIndicesWithStepSize := make([]int, 0)
	totalDataLengthInBytesAsInt := int(totalDataLengthInBytes)
	for j := 0; j <= totalDataLengthInBytesAsInt; j += stepSizeForIndicesAsInt {
		jAsString := fmt.Sprintf("%d", j)
		currentXorDistance := xordistance.ComputeXorDistanceBetweenTwoStrings(jAsString, comparisonString)
		sliceOfXorDistancesOfIndicesToBlockHash = append(sliceOfXorDistancesOfIndicesToBlockHash, currentXorDistance)
		sliceOfIndicesWithStepSize = append(sliceOfIndicesWithStepSize, j)
	}
	sliceOfSortedIndices := argsort.SortSlice(sliceOfXorDistancesOfIndicesToBlockHash, func(i, j int) bool {
		return sliceOfXorDistancesOfIndicesToBlockHash[i] < sliceOfXorDistancesOfIndicesToBlockHash[j]
	})
	sliceOfSortedIndicesWithStepSize := make([]int, 0)
	for _, currentSortedIndex := range sliceOfSortedIndices {
		sliceOfSortedIndicesWithStepSize = append(sliceOfSortedIndicesWithStepSize, sliceOfIndicesWithStepSize[currentSortedIndex])
	}
	firstTwoSortedIndices := sliceOfSortedIndicesWithStepSize[0:2]
	challengeSliceStartIndex, challengeSliceEndIndex := minMax(firstTwoSortedIndices)
	return challengeSliceStartIndex, challengeSliceEndIndex
}

func minMax(array []int) (int, int) {
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
