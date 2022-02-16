package storagechallenge

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

func (task StorageChallengeTask) GenerateStorageChallenges(ctx context.Context, challengesPerNodePerBlock int) error {
	log.WithContext(ctx).Println("Generating Storage Challenges called.")
	/* ------------------------------------------------- */
	/* ------- Original implementation goes here ------- */
	// list all RQ symbol keys from nft ticket
	sliceOfFileHashes, err := task.repository.ListSymbolFileKeysFromNFTTicket(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}

	sliceToCheckIfFileContainedByLocalSupernode := make([]bool, 0)
	for _, currentFileHash := range sliceOfFileHashes {
		_, err = task.repository.GetSymbolFileByKey(ctx, currentFileHash, true)
		if err == nil {
			sliceToCheckIfFileContainedByLocalSupernode = append(sliceToCheckIfFileContainedByLocalSupernode, true)
		} else {
			sliceToCheckIfFileContainedByLocalSupernode = append(sliceToCheckIfFileContainedByLocalSupernode, false)
		}
	}

	sliceOfFileHashesStoredByChallenger := make([]string, 0)
	for idx, currentFileContainedByLocalMN := range sliceToCheckIfFileContainedByLocalSupernode {
		if currentFileContainedByLocalMN {
			sliceOfFileHashesStoredByChallenger = append(sliceOfFileHashesStoredByChallenger, sliceOfFileHashes[idx])
		}
	}
	/* --------------------------------------------------------------- */
	/* ------ Prepare merkleroot and challenging Supernode id ------- */
	// collect current block number get merkle root from block verbose
	currentBlockCount, err := task.pclient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := task.pclient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return err
	}
	//current block hash
	merkleroot := blkVerbose1.MerkleRoot

	// get all Supernode by pastel client, choose randomly challenging Supernode id from this list
	listOfSupernode, err := task.repository.GetListOfSupernode(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list of Supernode using pastel client")
		return err
	}

	// get random Supernode in list tobe chalenging Supernode ignore this node
	challengingSupernode := task.repository.GetNClosestSupernodeIDsToComparisonString(ctx, 1, utils.GetHashFromString(merkleroot+fmt.Sprint(rand.Int63())), listOfSupernode, task.nodeID)
	if len(challengingSupernode) == 0 {
		err = fmt.Errorf("not enough Supernode to assign as challenging")
		log.WithContext(ctx).WithError(err).Error("could not get random challenging Supernode from list")
		return err
	}

	challengingSupernodeID := challengingSupernode[0]

	comparisonStringForFileHashSelection := merkleroot + challengingSupernodeID
	sliceOfFileHashesToChallenge := task.repository.GetNClosestFileHashesToAGivenComparisonString(ctx, challengesPerNodePerBlock, comparisonStringForFileHashSelection, sliceOfFileHashesStoredByChallenger)
	sliceOfSupernodesToChallenge := make([]string, len(sliceOfFileHashesToChallenge))

	log.WithContext(ctx).Infof("Challenging Supernode %s is now selecting file hashes to challenge this block, and then for each one, selecting which Supernode to challenge...", challengingSupernodeID)

	for idx1, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		sliceOfSupernodesStoringFileHash := task.repository.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, task.numberOfChallengeReplicas, currentFileHashToChallenge, task.nodeID)
		sliceOfSupernodesStoringFileHashExcludingChallenger := make([]string, 0)
		ignoreIDX := -1
		for idx, currentSupernodeID := range sliceOfSupernodesStoringFileHash {
			if currentSupernodeID == challengingSupernodeID {
				ignoreIDX = idx
			}
		}
		if ignoreIDX >= 0 {
			sliceOfSupernodesStoringFileHashExcludingChallenger = append(sliceOfSupernodesStoringFileHash[:ignoreIDX], sliceOfSupernodesStoringFileHash[ignoreIDX+1:]...)
		} else {
			sliceOfSupernodesStoringFileHashExcludingChallenger = append(sliceOfSupernodesStoringFileHashExcludingChallenger, sliceOfSupernodesStoringFileHash...)
		}
		comparisonStringForSupernodeSelection := merkleroot + currentFileHashToChallenge + challengingSupernodeID + utils.GetHashFromString(fmt.Sprint(idx1))

		respondingSupernodeIDs := task.repository.GetNClosestSupernodeIDsToComparisonString(ctx, 1, comparisonStringForSupernodeSelection, sliceOfSupernodesStoringFileHashExcludingChallenger)

		sliceOfSupernodesToChallenge[idx1] = respondingSupernodeIDs[0]
	}

	for idx2, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		b, err := task.repository.GetSymbolFileByKey(ctx, currentFileHashToChallenge, false)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Supernode %s encountered an error generating storage challenges", challengingSupernodeID)
			continue
		}
		challengeDataSize := uint64(len(b))
		if challengeDataSize == 0 {
			log.WithContext(ctx).Errorf("Supernode %s encountered an invalid file while attempting to generate a storage challenge for file hash %s", challengingSupernodeID, currentFileHashToChallenge)
			continue
		}

		respondingSupernodeID := sliceOfSupernodesToChallenge[idx2]
		challengeStatus := pb.StorageChallengeData_Status_PENDING
		messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(challengeDataSize, currentFileHashToChallenge, merkleroot, challengingSupernodeID)
		messageIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + challengeStatus.String() + messageType.String() + merkleroot
		messageID := utils.GetHashFromString(messageIDInputData)
		challengeIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(merkleroot) + fmt.Sprint(rand.Int63())
		challengeID := utils.GetHashFromString(challengeIDInputData)
		outgoingChallengeMessage := &pb.StorageChallengeData{
			MessageId:                    messageID,
			MessageType:                  messageType,
			ChallengeStatus:              challengeStatus,
			BlockNumChallengeSent:        currentBlockCount,
			BlockNumChallengeRespondedTo: 0,
			BlockNumChallengeVerified:    0,
			MerklerootWhenChallengeSent:  merkleroot,
			ChallengingMasternodeId:      challengingSupernodeID,
			RespondingMasternodeId:       respondingSupernodeID,
			ChallengeFile: &pb.StorageChallengeDataChallengeFile{
				FileHashToChallenge:      currentFileHashToChallenge,
				ChallengeSliceStartIndex: int64(challengeSliceStartIndex),
				ChallengeSliceEndIndex:   int64(challengeSliceEndIndex),
			},
			ChallengeSliceCorrectHash: "",
			ChallengeResponseHash:     "",
			ChallengeId:               challengeID,
		}

		if err = task.sendprocessStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
			log.WithContext(ctx).WithError(err).Error("could not send process storage challenge")
		}

		task.repository.SaveChallengMessageState(ctx, "sent", challengeID, challengingSupernodeID, currentBlockCount)
	}

	return nil
}

func (task *StorageChallengeTask) sendprocessStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error {
	Supernodes, err := task.pclient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range Supernodes {
		mapSupernodes[mn.ExtKey] = mn
	}

	var mn pastel.MasterNode
	var ok bool
	if mn, ok = mapSupernodes[challengeMessage.ChallengingMasternodeId]; !ok {
		err = fmt.Errorf("cannot get Supernode info of Supernode id %v", challengeMessage.ChallengingMasternodeId)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	processingSupernodeAddr := mn.ExtAddress
	log.WithContext(ctx).Println("Sending storage challenge to processing supernode address: " + processingSupernodeAddr)

	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use nodeclient to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	storageChallengeIF := nodeClientConn.StorageChallenge()
	//Sends the process storage challenge message to the connected processing supernode
	return storageChallengeIF.ProcessStorageChallenge(ctx, challengeMessage)
	// return s.actor.Send(ctx, s.domainActorID, newSendProcessStorageChallengeMsg(ctx, processingSupernodeAddr, challengeMessage))
}

func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingSupernodeID string) (int, int) {
	blockHashStringAsInt, _ := strconv.ParseInt(blockHashString, 16, 64)
	blockHashStringAsIntStr := fmt.Sprint(blockHashStringAsInt)
	stepSizeForIndicesStr := blockHashStringAsIntStr[len(blockHashStringAsIntStr)-1:] + blockHashStringAsIntStr[0:1]
	stepSizeForIndices, _ := strconv.ParseUint(stepSizeForIndicesStr, 10, 32)
	stepSizeForIndicesAsInt := int(stepSizeForIndices)
	comparisonString := blockHashString + fileHashString + challengingSupernodeID
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
