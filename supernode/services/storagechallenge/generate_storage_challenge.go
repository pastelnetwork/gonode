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
	// TODO: replacing this with real repository implementation
	// symbolFiles, err := s.repository.GetSymbolFiles(ctx)
	// if err != nil {
	// 	log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Error("could not get symbol files", actorLog.String("s.repository.GetSymbolFiles", err.Error()))
	// 	return err
	// }
	var symbolFiles = []*SymbolFile{}

	var mapSymbolFileByFileHash = make(map[string]*SymbolFile)
	for _, symbolFile := range symbolFiles {
		mapSymbolFileByFileHash[symbolFile.FileHash] = symbolFile
	}

	comparisonStringForFileHashSelection := currentBlockHash + challengingMasternodeID
	sliceOfFileHashesToChallenge := xordistance.GetNClosestXORDistanceStringToAGivenComparisonString(challengesPerMasternodePerBlock, comparisonStringForFileHashSelection, _symbolFiles(symbolFiles).GetListXORDistanceString())

	for idx, symbolFileHash := range sliceOfFileHashesToChallenge {
		challengeDataSize := mapSymbolFileByFileHash[symbolFileHash].FileLengthInBytes

		// TODO: replacing this with real repository implementation
		// selecting n closest node excepting challenger (current node)
		// xorDistances, err := s.repository.GetTopRankedXorDistanceMasternodeToFileHash(ctx, symbolFileHash, s.numberOfChallengeReplicas, s.nodeID)
		// if err != nil {
		// 	// ignore this file hash
		// 	log.With(actorLog.String("ACTOR", "GenerateStorageChallenges")).Warn(fmt.Sprintf("could not get top %v ranked xor of distance masternodes id to file hash %s", s.numberOfChallengeReplicas, symbolFileHash), actorLog.String("s.repository.GetTopRankedXorDistanceMasternodeToFileHash", err.Error()))
		// 	continue
		// }
		var xorDistances = []*XORDistance{}

		comparisonStringForMasternodeSelection := currentBlockHash + symbolFileHash + s.nodeID + helper.GetHashFromString(fmt.Sprint(idx))
		respondingMasternodesID := xordistance.GetNClosestXORDistanceStringToAGivenComparisonString(1, comparisonStringForMasternodeSelection, _xorDistances(xorDistances).GetListXORDistanceString())
		challengeStatus := Status_PENDING
		messageType := MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(uint64(challengeDataSize), symbolFileHash, currentBlockHash, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + respondingMasternodesID[0] + symbolFileHash + challengeStatus + messageType + currentBlockHash
		messageID := helper.GetHashFromString(messageIDInputData)
		timestampChallengeSent := time.Now().UnixMilli()
		challengeIDInputData := challengingMasternodeID + respondingMasternodesID[0] + symbolFileHash + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(timestampChallengeSent)
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
			RespondingMasternodeID:        respondingMasternodesID[0],
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

type _symbolFiles []*SymbolFile

func (s _symbolFiles) GetListXORDistanceString() []string {
	ret := make([]string, len(s))
	for idx, symbolFile := range s {
		ret[idx] = symbolFile.FileHash
	}

	return ret
}

type _xorDistances []*XORDistance

func (s _xorDistances) GetListXORDistanceString() []string {
	ret := make([]string, len(s))
	for idx, xorDistance := range s {
		ret[idx] = xorDistance.MasternodeID
	}

	return ret
}

func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingMasternodeId string) (int, int) {
	blockHashStringAsInt, _ := strconv.ParseInt(blockHashString, 16, 64)
	blockHashStringAsIntStr := fmt.Sprint(blockHashStringAsInt)
	stepSizeForIndicesStr := blockHashStringAsIntStr[len(blockHashStringAsIntStr)-1:] + blockHashStringAsIntStr[0:1]
	stepSizeForIndices, _ := strconv.ParseUint(stepSizeForIndicesStr, 10, 32)
	stepSizeForIndicesAsInt := int(stepSizeForIndices)
	comparisonString := blockHashString + fileHashString + challengingMasternodeId
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
	var max int = array[0]
	var min int = array[0]
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
