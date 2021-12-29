package storagechallenge

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
)

func (s *service) GenerateStorageChallenges(ctx context.Context, _ int) error {
	// list all symbol keys from nft ticket
	symbolFileKeys, err := s.repository.ListSymbolFileKeysFromNFTTicket(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}

	// collect current block number
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

	// get merkle root from current block verbose
	merkleroot := blkVerbose1.MerkleRoot
	log.WithContext(ctx).Infof("merkleRoot generate storage challenges: %s", merkleroot)
	fileHashComparisonString := utils.GetHashFromString(merkleroot + fmt.Sprint(rand.Int63()))

	// the default of number of symbol keys that challenger will generate the challenge for is 3
	if len(symbolFileKeys) > 3 {
		symbolFileKeys = s.repository.GetNClosestXORDistanceToComparisonString(ctx, 3, fileHashComparisonString, symbolFileKeys)
	}
	log.WithContext(ctx).Infof("closest 3 symbol file keys: %v", symbolFileKeys)

	for _, symbolFileKey := range symbolFileKeys {
		b, err := s.repository.GetSymbolFileByKey(ctx, symbolFileKey, false)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("could not symbolfile by key")
			continue
		}
		challengeDataSize := len(b)
		listMasternodesSavingSymbolFileHash := s.repository.GetNClosestXORDistanceMasternodesToComparisionString(ctx, 6, symbolFileKey, s.nodeID)
		closestMasternodesToMerkleRoot := s.repository.GetNClosestXORDistanceToComparisonString(ctx, 2, utils.GetHashFromString(merkleroot), listMasternodesSavingSymbolFileHash)
		if len(closestMasternodesToMerkleRoot) < 2 {
			continue
		}
		challengingMasternodeID := closestMasternodesToMerkleRoot[0]
		respondingMasternodeID := closestMasternodesToMerkleRoot[1]
		challengeStatus := statusPending
		messageType := storageChallengeIssuanceMessage
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(uint64(challengeDataSize), symbolFileKey, merkleroot, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + respondingMasternodeID + symbolFileKey + challengeStatus + messageType + merkleroot
		messageID := utils.GetHashFromString(messageIDInputData)
		challengeIDInputData := challengingMasternodeID + respondingMasternodeID + symbolFileKey + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(blockNumChallengeSent)
		challengeID := utils.GetHashFromString(challengeIDInputData)
		outgoingChallengeMessage := &ChallengeMessage{
			MessageID:                    messageID,
			MessageType:                  messageType,
			ChallengeStatus:              challengeStatus,
			BlockNumChallengeSent:        blockNumChallengeSent,
			BlockNumChallengeRespondedTo: 0,
			BlockNumChallengeVerified:    0,
			MerklerootWhenChallengeSent:  merkleroot,
			ChallengingMasternodeID:      challengingMasternodeID,
			RespondingMasternodeID:       respondingMasternodeID,
			FileHashToChallenge:          symbolFileKey,
			ChallengeSliceStartIndex:     uint64(challengeSliceStartIndex),
			ChallengeSliceEndIndex:       uint64(challengeSliceEndIndex),
			ChallengeSliceCorrectHash:    "",
			ChallengeResponseHash:        "",
			ChallengeID:                  challengeID,
		}

		// s.sendStatictis(ctx, outgoingChallengeMessage, "sent")
		s.sendprocessStorageChallenge(ctx, outgoingChallengeMessage)
	}

	return nil
}

func (s *service) sendprocessStorageChallenge(ctx context.Context, challengeMessage *ChallengeMessage) error {
	masternodes, err := s.pclient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendprocessStorageChallenge").WithError(err).Warn("could not get masternode extra: ", err.Error())
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
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	processingMasternodeAddr := mn.ExtAddress

	return s.actor.Send(ctx, s.domainActorID, newSendProcessStorageChallengeMsg(ctx, processingMasternodeAddr, challengeMessage))
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
