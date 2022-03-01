package storagechallenge

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

//Generate storage challenge tasks when a new block is detected
// Order of operations:
// List all NFT tickets, get their raptorq ids
// Identify which raptorq files are currently hosted on this node
// Identify current block number
// Identify merkle root from current block
// Get a list of supernodes
// Calculate the number of supernodes to issue challenge per block
// Calculate the challenges per challenger
// Identify which supernodes should issue challenges for this block
// determine if this node is a challenger, continue if so
// Identify which files should be challenged
// Identify which supernodes have our file
// Identify challenge slice indices
// Set relevant outgoing challenge message details
// Send the storage challenge message for processing by the responder
// Save the challenge state

//Called from service run, generate storage challenges will determine if we should issue a storage challenge,
// and if so calculate and issue it.
func (task StorageChallengeTask) GenerateStorageChallenges(ctx context.Context) error {
	log.WithContext(ctx).Println("Generating Storage Challenges called.")
	// List all NFT tickets, get their raptorq ids
	// list all RQ symbol keys from nft ticket
	sliceOfFileHashes, err := task.ListSymbolFileKeysFromNFTTicket(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}

	// Identify which raptorq files are currently hosted on this node
	sliceToCheckIfFileContainedByLocalSupernode := make([]bool, 0)
	for _, currentFileHash := range sliceOfFileHashes {
		_, err = task.GetSymbolFileByKey(ctx, currentFileHash, true)
		if err == nil {
			sliceToCheckIfFileContainedByLocalSupernode = append(sliceToCheckIfFileContainedByLocalSupernode, true)
		} else {
			sliceToCheckIfFileContainedByLocalSupernode = append(sliceToCheckIfFileContainedByLocalSupernode, false)
		}
	}

	sliceOfFileHashesStoredByLocalSupernode := make([]string, 0)
	for idx, currentFileContainedByLocalMN := range sliceToCheckIfFileContainedByLocalSupernode {
		if currentFileContainedByLocalMN {
			sliceOfFileHashesStoredByLocalSupernode = append(sliceOfFileHashesStoredByLocalSupernode, sliceOfFileHashes[idx])
		}
	}

	// Identify current block number
	// collect current block number get merkle root from block verbose
	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return err
	}
	//current block hash
	merkleroot := blkVerbose1.MerkleRoot

	// Get a list of supernodes
	// get all Supernode by pastel client, choose challenging Supernode id from this list
	listOfSupernodes, err := task.GetListOfSupernode(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list of Supernode using pastel client")
		return err
	}

	numberOfSupernodesDividedByThree := int(math.Ceil(float64(len(listOfSupernodes)) / 3))
	// Calculate the number of supernodes to issue challenge per block
	numberOfSupernodesToIssueChallengePerBlock := numberOfSupernodesDividedByThree
	// Calculate the challenges per challenger
	challengesPerSupernodePerBlock := numberOfSupernodesDividedByThree

	// Identify which supernodes should issue challenges for this block
	sliceOfChallengingSupernodeIDsForBlock := task.GetNClosestSupernodeIDsToComparisonString(ctx, numberOfSupernodesToIssueChallengePerBlock, merkleroot, listOfSupernodes)

	//If we are in the sliceOfChallengingSupernodeIDsForBlock, we need to generate a storage challenge.  Otherwise we don't.
	isMyNodeAChallenger := false
	for _, supernodeID := range sliceOfChallengingSupernodeIDsForBlock {
		if supernodeID == task.nodeID {
			isMyNodeAChallenger = true
		}
	}
	if !isMyNodeAChallenger {
		return nil
	}
	challengingSupernodeID := task.nodeID

	// Identify which files should be challenged
	comparisonStringForFileHashSelection := merkleroot + challengingSupernodeID
	sliceOfFileHashesToChallenge := task.GetNClosestFileHashesToAGivenComparisonString(ctx, challengesPerSupernodePerBlock, comparisonStringForFileHashSelection, sliceOfFileHashesStoredByLocalSupernode)
	sliceOfSupernodesToChallenge := make([]string, len(sliceOfFileHashesToChallenge))

	//log.WithContext(ctx).Infof("Challenging Supernode %s is now selecting file hashes to challenge this block, and then for each one, selecting which Supernode to challenge...", challengingSupernodeID)
	// Identify which supernodes have our file
	for idx1, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		sliceOfSupernodesStoringFileHashExcludingChallenger := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, task.numberOfChallengeReplicas, currentFileHashToChallenge, task.nodeID)

		comparisonStringForSupernodeSelection := merkleroot + currentFileHashToChallenge + challengingSupernodeID + utils.GetHashFromString(fmt.Sprint(idx1))

		respondingSupernodeIDs := task.GetNClosestSupernodeIDsToComparisonString(ctx, 1, comparisonStringForSupernodeSelection, sliceOfSupernodesStoringFileHashExcludingChallenger)

		sliceOfSupernodesToChallenge[idx1] = respondingSupernodeIDs[0]
	}
	// challenge those supernodes (they become responder) for the hash selected

	for idx2, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		b, err := task.GetSymbolFileByKey(ctx, currentFileHashToChallenge, false)
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
		// Identify challenge slice indices
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(challengeDataSize, currentFileHashToChallenge, merkleroot, challengingSupernodeID)
		messageIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + challengeStatus.String() + messageType.String() + merkleroot
		messageID := utils.GetHashFromString(messageIDInputData)
		challengeIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(merkleroot) + fmt.Sprint(rand.Int63())
		challengeID := utils.GetHashFromString(challengeIDInputData)
		// Set relevant outgoing challenge message details
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
		// Send the storage challenge message for processing by the responder
		if err = task.SendProcessStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error processing storage challenge: ")
		}
		// Save the challenge state
		task.SaveChallengeMessageState(ctx, "sent", challengeID, challengingSupernodeID, currentBlockCount)
	}

	return nil
}

// Send the storage challenge message to the responding node
func (task *StorageChallengeTask) SendProcessStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error {
	//Get masternodes (supernodes)
	Supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
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
	if mn, ok = mapSupernodes[challengeMessage.RespondingMasternodeId]; !ok {
		err = fmt.Errorf("cannot get Supernode info of Supernode id %v", challengeMessage.RespondingMasternodeId)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	//We use the ExtAddress of the supernode to connect
	processingSupernodeAddr := mn.ExtAddress
	log.WithContext(ctx).Println("Sending storage challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use nodeclient to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
		return err
	}
	storageChallengeIF := nodeClientConn.StorageChallenge()
	//Calls the ProcessStorageChallenge method on the connected supernode over GRPC.
	return storageChallengeIF.ProcessStorageChallenge(ctx, challengeMessage)
	// return s.actor.Send(ctx, s.domainActorID, newSendProcessStorageChallengeMsg(ctx, processingSupernodeAddr, challengeMessage))
}

//This is how we programmatically determine which pieces of the file get read and hashed for challenging to determine if the supernodes are properly hosting.
func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingSupernodeID string) (int, int) {
	//raptorq files are kb in length, so this is unlikely to be run but this check might be useful at a later date
	if totalDataLengthInBytes < 200 {
		return 0, int(totalDataLengthInBytes) - 1
	}
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
