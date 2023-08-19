package storagechallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
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

// GenerateStorageChallenges is called from service run, generate storage challenges will determine if we should issue a storage challenge,
// and if so calculate and issue it.
func (task SCTask) GenerateStorageChallenges(ctx context.Context) error {
	log.WithContext(ctx).Println("Generate Storage Challenges invoked")
	// List all NFT tickets and API Cascade Tickets, get their raptorq ids
	// list all RQ symbol keys from nft and action tickets

	log.WithContext(ctx).Info("list symbol file keys from registered nft tickets")
	sliceOfFileHashes, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list symbol file keys")
		return err
	}
	log.WithContext(ctx).Info("symbol file keys from registered nft tickets have been retrieved")

	// Identify which raptorq files are currently hosted on this node
	sliceOfFileHashesStoredByLocalSupernode := task.getFilesStoredByLocalSN(ctx, sliceOfFileHashes)
	log.WithContext(ctx).Info("slice of file hashes stored by local SN has been retrieved")

	log.WithContext(ctx).Info("retrieving block no and verbose")
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
	log.WithContext(ctx).Info("block no and verbose retrieved")

	// Identify challengers
	sliceOfChallengingSupernodeIDsForBlock, _, challengesPerBlock, err := task.getChallengers(ctx, merkleroot)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error identifying challengers")
		return err
	}

	//If we are in the sliceOfChallengingSupernodeIDsForBlock, we need to generate a storage challenge, Otherwise we don't.
	var challengingSupernodeID string
	isMyNodeAChallenger := task.isMyNodeChallenger(ctx, sliceOfChallengingSupernodeIDsForBlock)
	if !isMyNodeAChallenger {
		return nil
	}
	challengingSupernodeID = task.nodeID

	// Identify which files should be challenged, their recipients and observers
	sliceOfFileHashesToChallenge := task.getChallengingFiles(ctx, merkleroot, challengingSupernodeID, challengesPerBlock, sliceOfFileHashesStoredByLocalSupernode)

	log.WithContext(ctx).Info("identifying challenge recipients and partial observers")
	sliceOfSupernodesToChallenge, challengeFileObservers := task.identifyChallengeRecipientsAndObserversAgainstChallengingFiles(ctx, sliceOfFileHashesToChallenge, merkleroot, challengingSupernodeID)

	for idx2, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		b, err := task.GetSymbolFileByKey(ctx, currentFileHashToChallenge, true)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("Supernode %s encountered an error generating storage challenges", challengingSupernodeID)
			continue
		}

		challengeDataSize := uint64(len(b))
		if challengeDataSize == 0 {
			log.WithContext(ctx).Errorf("Supernode %s encountered an invalid file while attempting to generate a storage challenge for file hash %s", challengingSupernodeID, currentFileHashToChallenge)
			continue
		}

		messageType := types.ChallengeMessageType
		respondingSupernodeID := sliceOfSupernodesToChallenge[idx2]

		// Identify challenge slice indices
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(challengeDataSize, currentFileHashToChallenge, merkleroot, challengingSupernodeID)

		//Create Challenge And Message IDs
		challengeIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(merkleroot) + fmt.Sprint(rand.Int63())
		challengeID := utils.GetHashFromString(challengeIDInputData)

		storageChallengeMessage := types.Message{
			MessageType: messageType,
			ChallengeID: challengeID,
			Data: types.MessageData{
				ChallengerID: challengingSupernodeID,
				Observers:    challengeFileObservers[idx2],
				RecipientID:  sliceOfSupernodesToChallenge[idx2],
				Challenge: types.ChallengeData{
					Block:      currentBlockCount,
					Merkelroot: merkleroot,
					Timestamp:  time.Now(),
					FileHash:   currentFileHashToChallenge,
					StartIndex: challengeSliceStartIndex,
					EndIndex:   challengeSliceEndIndex,
				},
			},
			Sender:          challengingSupernodeID,
			SenderSignature: nil,
		}

		log.WithContext(ctx).WithField("challenge_id", challengeID).Info("sending challenge for processing")
		// Send the storage challenge message for processing by the responder
		if err = task.sendProcessStorageChallenge(ctx, storageChallengeMessage); err != nil {
			log.WithContext(ctx).WithError(err).WithField("challenge_id", storageChallengeMessage.ChallengeID).
				Error("Error processing storage challenge: ")
		}
		// Save the challenge state
		task.SaveChallengeMessageState(ctx, "sent", challengeID, challengingSupernodeID, currentBlockCount)
	}

	log.WithContext(ctx).Info("files have been challenged")
	return nil
}

// sendProcessStorageChallenge will send the storage challenge message to the responding node
func (task *SCTask) sendProcessStorageChallenge(ctx context.Context, challengeMessage types.Message) error {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send process storage challenge")
		return err
	}

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing challenge message")
		return err
	}
	challengeMessage.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, challengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing storage challenge message")
	}

	msg := pb.StorageChallengeMessage{
		MessageType:     pb.StorageChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}

	for _, node := range nodesToConnect {
		logger := log.WithContext(ctx).WithField("node_address", node.ExtAddress)

		if err := task.SendMessage(ctx, msg, node.ExtAddress); err != nil {
			logger.WithError(err).Error("error sending storage challenge message for processing")
			continue
		}

		logger.Info("challenge message has been sent")
	}

	return nil
}

// This is how we programmatically determine which pieces of the file get read and hashed for challenging to determine if the supernodes are properly hosting.
func getStorageChallengeSliceIndices(totalDataLengthInBytes uint64, fileHashString string, blockHashString string, challengingSupernodeID string) (int, int) {
	if totalDataLengthInBytes < 200 {
		return 0, int(totalDataLengthInBytes) - 1
	}

	blockHashStringAsBigInt := new(big.Int)
	blockHashStringAsBigInt.SetString(blockHashString, 16)
	blockHashStringAsIntStr := blockHashStringAsBigInt.String()
	stepSizeForIndicesStr := blockHashStringAsIntStr[len(blockHashStringAsIntStr)-1:] + blockHashStringAsIntStr[0:1]
	stepSizeForIndices, _ := strconv.Atoi(stepSizeForIndicesStr)

	comparisonString := blockHashString + fileHashString + challengingSupernodeID
	sliceOfXorDistancesOfIndicesToBlockHash := make([]*big.Int, 0)
	sliceOfIndicesWithStepSize := make([]int, 0)
	totalDataLengthInBytesAsInt := int(totalDataLengthInBytes)

	for j := 0; j <= totalDataLengthInBytesAsInt; j += stepSizeForIndices {
		jAsString := fmt.Sprintf("%d", j)
		currentXorDistance := utils.ComputeXorDistanceBetweenTwoStrings(jAsString, comparisonString)
		sliceOfXorDistancesOfIndicesToBlockHash = append(sliceOfXorDistancesOfIndicesToBlockHash, currentXorDistance)
		sliceOfIndicesWithStepSize = append(sliceOfIndicesWithStepSize, j)
	}

	sliceOfSortedIndices := argsort.SortSlice(sliceOfXorDistancesOfIndicesToBlockHash, func(i, j int) bool {
		return sliceOfXorDistancesOfIndicesToBlockHash[i].Cmp(sliceOfXorDistancesOfIndicesToBlockHash[j]) < 0
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
	if array[0] > array[1] {
		return array[1], array[0]
	}

	return array[0], array[1]
}

// GetNodesAddressesToConnect basically retrieves the masternode address against the pastel-id from the list and return that
func (task *SCTask) GetNodesAddressesToConnect(ctx context.Context, challengeMessage types.Message) ([]pastel.MasterNode, error) {
	var nodesToConnect []pastel.MasterNode
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendProcessStorageChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		mapSupernodes[mn.ExtKey] = mn
	}

	switch challengeMessage.MessageType {
	case types.ChallengeMessageType:
		//send message to challenge recipient
		nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

		//and partial observers; that's not required we only need partial observers for verification though
		for _, po := range challengeMessage.Data.Observers {
			if value, ok := mapSupernodes[po]; ok {
				nodesToConnect = append(nodesToConnect, value)
			}
		}

		return nodesToConnect, nil
	case types.ResponseMessageType:
		if challengeMessage.Data.RecipientID == task.nodeID { //if challenge recipient is the responder
			//send message to challenger
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.ChallengerID])

			//and partial observers for verification
			for _, po := range challengeMessage.Data.Observers { //partial observers
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				}
			}

		}

		return nodesToConnect, nil
	case types.EvaluationMessageType:
		for _, po := range challengeMessage.Data.Observers { //partial observers
			if value, ok := mapSupernodes[po]; ok {
				nodesToConnect = append(nodesToConnect, value)
			}
		}

		return nodesToConnect, nil
	case types.AffirmationMessageType:
	// not required

	case types.BroadcastMessageType:
		return supernodes, nil

	default:
		return nil, errors.Errorf("no nodes found to be connect")
	}

	return nil, err
}

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *SCTask) SendMessage(ctx context.Context, challengeMessage pb.StorageChallengeMessage, processingSupernodeAddr string) error {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("Sending storage challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use nodeclient to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendProcessStorageChallenge").Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	switch challengeMessage.MessageType {
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_CHALLENGE_MESSAGE:
		return storageChallengeIF.ProcessStorageChallenge(ctx, &challengeMessage)
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE:
		return storageChallengeIF.VerifyStorageChallenge(ctx, &challengeMessage)

	default:
		log.WithContext(ctx).Info("message type not supported by any Process & Verify worker")
	}

	return nil
}

func (task SCTask) isMyNodeChallenger(ctx context.Context, sliceOfChallengingSupernodeIDsForBlock []string) bool {
	for _, supernodeID := range sliceOfChallengingSupernodeIDsForBlock {
		if supernodeID == task.nodeID {
			return true
		}
	}

	log.WithContext(ctx).Info("exit because this node is not a challenger")
	return false
}

func (task SCTask) getFilesStoredByLocalSN(ctx context.Context, sliceOfFileHashes []string) (sliceOfFileHashesStoredByLocalSupernode []string) {
	log.WithContext(ctx).Info("identifying which files are currently hosted on this node")

	for _, currentFileHash := range sliceOfFileHashes {
		_, err := task.GetSymbolFileByKey(ctx, currentFileHash, true)
		if err == nil {
			sliceOfFileHashesStoredByLocalSupernode = append(sliceOfFileHashesStoredByLocalSupernode, currentFileHash)
		}
	}
	log.WithContext(ctx).Info("files hosted on this node have been identified")

	return sliceOfFileHashesStoredByLocalSupernode
}

func (task SCTask) getChallengers(ctx context.Context, merkleRoot string) (sliceOfChallengers []string, challengersPerBlock, challengesPerBlock int, err error) {
	log.WithContext(ctx).Info("retrieving list of all super nodes")
	listOfSupernodes, err := task.GetListOfSupernode(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list of Supernode using pastel client")
		return nil, 0, 0, err
	}
	log.WithContext(ctx).Info("list of supernodes have been retrieved")

	log.WithContext(ctx).Info("identifying no of challengers & challenges to issue for this block")
	numberOfSupernodesDividedByThree := int(math.Ceil(float64(len(listOfSupernodes)) / 3))
	// Calculate the number of supernodes to issue challenge per block
	challengersPerBlock = numberOfSupernodesDividedByThree
	// Calculate the challenges per challenger
	challengesPerBlock = numberOfSupernodesDividedByThree

	// Identify challengers for this block
	sliceOfChallengers = task.GetNClosestSupernodeIDsToComparisonString(ctx, challengersPerBlock, merkleRoot, listOfSupernodes)
	log.WithContext(ctx).Info("challengers have been selected")

	return sliceOfChallengers, challengersPerBlock, challengesPerBlock, nil
}

func (task SCTask) getChallengingFiles(ctx context.Context, merkleRoot, challengerNodeID string, challengesPerBlock int, sliceOfFileHashesStoredByLocalSupernode []string) (sliceOfFileHashesToChallenge []string) {
	log.WithContext(ctx).Info("identifying the files to challenge")
	comparisonStringForFileHashSelection := merkleRoot + challengerNodeID
	comparisonStringHashForChallengeFileSelection := utils.GetHashFromString(comparisonStringForFileHashSelection)

	log.WithContext(ctx).Info(fmt.Sprintf("comparison string hash for challenge file selection:%s", comparisonStringHashForChallengeFileSelection))
	sliceOfFileHashesToChallenge = task.GetNClosestFileHashesToAGivenComparisonString(ctx, challengesPerBlock, comparisonStringHashForChallengeFileSelection, sliceOfFileHashesStoredByLocalSupernode)
	log.WithContext(ctx).Info(fmt.Sprintf("total files to challenge are:%d", len(sliceOfFileHashesToChallenge)))

	return sliceOfFileHashesToChallenge

}

func (task SCTask) identifyChallengeRecipientsAndObserversAgainstChallengingFiles(ctx context.Context, sliceOfFileHashesToChallenge []string, merkleRoot, challengerNodeID string) ([]string, map[int][]string) {
	log.WithContext(ctx).Info("identifying challenge recipients and partial observers")
	mapOfchallengeFilePartialObservers := make(map[int][]string)
	sliceOfSupernodesToChallenge := make([]string, len(sliceOfFileHashesToChallenge))

	// Identify which supernodes have our file to challenge
	for idx1, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		sliceOfSupernodesStoringFileHashExcludingChallenger := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, task.numberOfChallengeReplicas, currentFileHashToChallenge, task.nodeID)
		comparisonStringForSupernodeSelection := merkleRoot + currentFileHashToChallenge + challengerNodeID + utils.GetHashFromString(fmt.Sprint(idx1))

		respondingSupernodeIDs := task.GetNClosestSupernodeIDsToComparisonString(ctx, 1, comparisonStringForSupernodeSelection, sliceOfSupernodesStoringFileHashExcludingChallenger)
		if len(respondingSupernodeIDs) < 1 {
			log.WithContext(ctx).WithField("file_hash", currentFileHashToChallenge).Info("no closest nodes have found against the file")
			continue
		}

		sliceOfSupernodesToChallenge[idx1] = respondingSupernodeIDs[0] //challenge recipient
		log.WithContext(ctx).Info("challenge recipient has been selected")

		//partial observers
		for _, sn := range sliceOfSupernodesStoringFileHashExcludingChallenger {
			if sn == respondingSupernodeIDs[0] {
				continue
			}

			mapOfchallengeFilePartialObservers[idx1] = append(mapOfchallengeFilePartialObservers[idx1], sn)
		}

		log.WithContext(ctx).Info("partial observers have been selected")
	}

	log.WithContext(ctx).Info("challenge recipients and partial observers against the files have been selected")
	return sliceOfSupernodesToChallenge, mapOfchallengeFilePartialObservers
}

// StoreChallengeMessage stores the challenge message to db for further verification
func (task SCTask) StoreChallengeMessage(ctx context.Context, msg types.Message) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	data, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}

	if store != nil {
		log.WithContext(ctx).Info("store")
		storageChallengeLog := types.StorageChallengeLogMessage{
			ChallengeID:     msg.ChallengeID,
			MessageType:     int(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
		}

		err = store.InsertStorageChallengeMessage(storageChallengeLog)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing challenge message to DB")
			return err
		}
	}

	return nil
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task SCTask) SignMessage(ctx context.Context, data types.MessageData) ([]byte, []byte, error) {
	d, err := json.Marshal(data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error marshaling the data")
	}

	signature, err := task.PastelClient.Sign(ctx, d, task.config.PastelID, task.config.PassPhrase, pastel.SignAlgorithmED448)
	if err != nil {
		return nil, nil, errors.Errorf("error signing storage challenge message: %w", err)
	}

	return signature, d, nil
}

// VerifyMessageSignature verifies the sender's signature on message
func (task SCTask) VerifyMessageSignature(ctx context.Context, msg types.Message) (bool, error) {
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return false, errors.Errorf("unable to marshal message data")
	}

	isVerified, err := task.PastelClient.Verify(ctx, data, string(msg.SenderSignature), msg.Sender, pastel.SignAlgorithmED448)
	if err != nil {
		return false, errors.Errorf("error verifying storage challenge message: %w", err)
	}

	return isVerified, nil
}
