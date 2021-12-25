package storagechallenge

import (
	"fmt"
	"strconv"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
)

func (s *service) GenerateStorageChallenges(ctx context.Context, challengesPerMasternodePerBlock int) error {
	symbolFileKeys, err := s.repository.ListKeys(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}

	blockNumChallengeSent, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}

	blkVerbose1, err := s.pclient.GetBlockVerbose1(ctx, blockNumChallengeSent)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return err
	}

	merkleroot := blkVerbose1.MerkleRoot

	if len(symbolFileKeys) > 3 {
		symbolFileKeys = s.repository.GetNClosestXORDistanceFileHashesToComparisonString(ctx, 3, merkleroot, symbolFileKeys)
	}

	for idx, symbolFileHash := range symbolFileKeys {
		challengeDataSize := 0
		challengingMasternode := s.repository.GetNClosestXORDistanceMasternodesToComparisionString(ctx, 1, symbolFileHash)
		if len(challengingMasternode) == 0 {
			continue
		}
		challengingMasternodeID := challengingMasternode[0].PastelID
		comparisonStringForMasternodeSelection := merkleroot + symbolFileHash + s.nodeID + utils.GetHashFromString(fmt.Sprint(idx))
		respondingMasternodes := s.repository.GetNClosestXORDistanceMasternodesToComparisionString(ctx, 1, comparisonStringForMasternodeSelection)
		challengeStatus := statusPending
		messageType := storageChallengeIssuanceMessage
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(uint64(challengeDataSize), symbolFileHash, merkleroot, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + string(respondingMasternodes[0].ID) + symbolFileHash + challengeStatus + messageType + merkleroot
		messageID := utils.GetHashFromString(messageIDInputData)
		challengeIDInputData := challengingMasternodeID + string(respondingMasternodes[0].ID) + symbolFileHash + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(blockNumChallengeSent)
		challengeID := utils.GetHashFromString(challengeIDInputData)
		outgoinChallengeMessage := &ChallengeMessage{
			MessageID:                    messageID,
			MessageType:                  messageType,
			ChallengeStatus:              challengeStatus,
			BlockNumChallengeSent:        blockNumChallengeSent,
			BlockNumChallengeRespondedTo: 0,
			BlockNumChallengeVerified:    0,
			MerklerootWhenChallengeSent:  merkleroot,
			ChallengingMasternodeID:      challengingMasternodeID,
			RespondingMasternodeID:       string(respondingMasternodes[0].ID),
			FileHashToChallenge:          symbolFileHash,
			ChallengeSliceStartIndex:     uint64(challengeSliceStartIndex),
			ChallengeSliceEndIndex:       uint64(challengeSliceEndIndex),
			ChallengeSliceCorrectHash:    "",
			ChallengeResponseHash:        "",
			ChallengeID:                  challengeID,
		}

		s.sendprocessStorageChallenge(ctx, outgoinChallengeMessage)
	}

	return nil
}

func (s *service) sendprocessStorageChallenge(ctx context.Context, challengeMessage *ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("method", "sendprocessStorageChallenge").Warn("could not get masternode info: ", err.Error())
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
	processingMasternodeClientPID := actor.NewPID(mn.ExtAddress, "storage-challenge")

	return s.actor.Send(ctx, s.domainActorID, newSendProcessStorageChallengeMsg(ctx, processingMasternodeClientPID, challengeMessage))
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
		currentXorDistance := utils.ComputeXorDistanceBetweenTwoStrings(jAsString, comparisonString)
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
