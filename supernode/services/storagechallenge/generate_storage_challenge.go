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

func (s *service) GenerateStorageChallenges(ctx context.Context, challengesPerNodePerBlock int) error {
	/* --------------------------------------------------------------- */
	/* ------ Prepare merkleroot and challenging masternode id ------- */
	// collect current block number get merkle root from block verbose
	currentBlockCount, err := s.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := s.pclient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return err
	}
	merkleroot := blkVerbose1.MerkleRoot

	// get all masternode by pastel client, choose randomly challenging masternode id from this list
	listOfMasternode, err := s.repository.GetListOfMasternode(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list of masternode using pastel client")
		return err
	}

	// get random masternode in list tobe chalenging masternode ignore this node
	challengingMasternode := s.repository.GetNClosestMasternodeIDsToComparisionString(ctx, 1, utils.GetHashFromString(merkleroot+fmt.Sprint(rand.Int63())), listOfMasternode, s.nodeID)
	if len(challengingMasternode) == 0 {
		err = fmt.Errorf("not enough masternode to assign as challenging")
		log.WithContext(ctx).WithError(err).Error("could not get random challenging masternode from list")
		return err
	}
	challengingMasternodeID := challengingMasternode[0]

	/* ------------------------------------------------- */
	/* ------- Original implementation goes here ------- */
	// list all RQ symbol keys from nft ticket
	sliceOfFileHashes, err := s.repository.ListSymbolFileKeysFromNFTTicket(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}

	sliceToCheckIfFileContainedByLocalMasternode := make([]bool, 0)
	for _, currentFileHash := range sliceOfFileHashes {
		_, err = s.repository.GetSymbolFileByKey(ctx, currentFileHash, true)
		if err == nil {
			sliceToCheckIfFileContainedByLocalMasternode = append(sliceToCheckIfFileContainedByLocalMasternode, true)
		} else {
			sliceToCheckIfFileContainedByLocalMasternode = append(sliceToCheckIfFileContainedByLocalMasternode, false)
		}
	}

	sliceOfFileHashesStoredByChallenger := make([]string, 0)
	for idx, currentFileContainedByLocalMN := range sliceToCheckIfFileContainedByLocalMasternode {
		if currentFileContainedByLocalMN {
			sliceOfFileHashesStoredByChallenger = append(sliceOfFileHashesStoredByChallenger, sliceOfFileHashes[idx])
		}
	}

	comparisonStringForFileHashSelection := merkleroot + challengingMasternodeID
	sliceOfFileHashesToChallenge := s.repository.GetNClosestFileHashesToAGivenComparisonString(ctx, challengesPerNodePerBlock, comparisonStringForFileHashSelection, sliceOfFileHashesStoredByChallenger)
	sliceOfMasternodesToChallenge := make([]string, len(sliceOfFileHashesToChallenge))

	log.WithContext(ctx).Infof("Challenging Masternode %s is now selecting file hashes to challenge this block, and then for each one, selecting which Masternode to challenge...", challengingMasternodeID)

	for idx1, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		sliceOfMasternodesStoringFileHash := s.repository.GetNClosestMasternodesToAGivenFileUsingKademlia(ctx, s.numberOfChallengeReplicas, currentFileHashToChallenge)
		sliceOfMasternodesStoringFileHashExcludingChallenger := make([]string, 0)
		for idx, currentMasternodeID := range sliceOfMasternodesStoringFileHash {
			if currentMasternodeID == challengingMasternodeID {
				sliceOfMasternodesStoringFileHashExcludingChallenger = append(sliceOfMasternodesStoringFileHash[:idx], sliceOfMasternodesStoringFileHash[idx+1:]...)
			}
		}
		comparisonStringForMasternodeSelection := merkleroot + currentFileHashToChallenge + challengingMasternodeID + utils.GetHashFromString(fmt.Sprint(idx1))
		respondingMasternodeIDs := s.repository.GetNClosestMasternodeIDsToComparisionString(ctx, 1, comparisonStringForMasternodeSelection, sliceOfMasternodesStoringFileHashExcludingChallenger)
		sliceOfMasternodesToChallenge[idx1] = respondingMasternodeIDs[0]
	}

	for idx2, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		b, err := s.repository.GetSymbolFileByKey(ctx, currentFileHashToChallenge, false)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("masternode %s encountered an error generating storage challenges", challengingMasternodeID)
			continue
		}
		challengeDataSize := uint64(len(b))
		if challengeDataSize == 0 {
			log.WithContext(ctx).Errorf("masternode %s encountered an invalid file while attempting to generate a storage challenge for file hash %s", challengingMasternodeID, currentFileHashToChallenge)
			continue
		}

		respondingMasternodeID := sliceOfMasternodesToChallenge[idx2]
		challengeStatus := statusPending
		messageType := storageChallengeIssuanceMessage
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(challengeDataSize, currentFileHashToChallenge, merkleroot, challengingMasternodeID)
		messageIDInputData := challengingMasternodeID + respondingMasternodeID + currentFileHashToChallenge + challengeStatus + messageType + merkleroot
		messageID := utils.GetHashFromString(messageIDInputData)
		challengeIDInputData := challengingMasternodeID + respondingMasternodeID + currentFileHashToChallenge + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(merkleroot) + fmt.Sprint(rand.Int63())
		challengeID := utils.GetHashFromString(challengeIDInputData)
		outgoingChallengeMessage := &ChallengeMessage{
			MessageID:                    messageID,
			MessageType:                  messageType,
			ChallengeStatus:              challengeStatus,
			BlockNumChallengeSent:        currentBlockCount,
			BlockNumChallengeRespondedTo: 0,
			BlockNumChallengeVerified:    0,
			MerklerootWhenChallengeSent:  merkleroot,
			ChallengingMasternodeID:      challengingMasternodeID,
			RespondingMasternodeID:       respondingMasternodeID,
			FileHashToChallenge:          currentFileHashToChallenge,
			ChallengeSliceStartIndex:     uint64(challengeSliceStartIndex),
			ChallengeSliceEndIndex:       uint64(challengeSliceEndIndex),
			ChallengeSliceCorrectHash:    "",
			ChallengeResponseHash:        "",
			ChallengeID:                  challengeID,
		}

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
