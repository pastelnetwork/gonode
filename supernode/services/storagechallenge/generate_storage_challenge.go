package storagechallenge

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	json "github.com/json-iterator/go"

	"github.com/google/uuid"

	"github.com/mkmik/argsort"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// GenerateStorageChallenges is called from service run, generate storage challenges will determine if we should issue a storage challenge,
// and if so calculate and issue it.
func (task SCTask) GenerateStorageChallenges(ctx context.Context) error {
	logger := log.WithContext(ctx).WithField("Method", "GenerateStorageChallenge")
	logger.Info("invoked")

	logger.Debug("retrieving block no and verbose")
	currentBlockCount, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get current block count")
		return err
	}
	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, currentBlockCount)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose 1")
		return err
	}

	//current block hash
	merkleroot := blkVerbose1.MerkleRoot
	logger.Debug("block no and verbose retrieved")

	// Identify challengers
	sliceOfChallengingSupernodeIDsForBlock, _, _, err := task.getChallengers(ctx, merkleroot)
	if err != nil {
		logger.WithError(err).Error("error identifying challengers")
		return err
	}
	logger.WithField("node_id", task.nodeID).Info("challengers have been selected")

	//If we are in the sliceOfChallengingSupernodeIDsForBlock, we need to generate a storage challenge, Otherwise we don't.
	var challengingSupernodeID string
	isMyNodeAChallenger := task.isMyNodeChallenger(ctx, sliceOfChallengingSupernodeIDsForBlock)
	if !isMyNodeAChallenger {
		return nil
	}
	challengingSupernodeID = task.nodeID

	sliceOfFileHashes, err := task.ListSymbolFileKeysFromNFTAndActionTickets(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get list symbol file keys")
		return err
	}
	logger.Debug("symbol file keys from registered nft tickets have been retrieved")

	// Identify which raptorq files are currently hosted on this node
	sliceOfFileHashesStoredByLocalSupernode := task.getFilesStoredByLocalSN(ctx, sliceOfFileHashes)
	logger.Debug("slice of file hashes stored by local SN has been retrieved")

	if len(sliceOfFileHashesStoredByLocalSupernode) == 0 {
		logger.Info("no files are hosted on this node - exiting")
		return nil
	}

	// Identify which files should be challenged, their recipients and observers
	sliceOfFileHashesToChallenge := task.getChallengingFiles(ctx, merkleroot, challengingSupernodeID, defaultChallengeReplicas, sliceOfFileHashesStoredByLocalSupernode)
	logger.Debug("file hashes to challenge stored by local SN have been retrieved")

	sliceOfSupernodesToChallenge, challengeFileObservers := task.identifyChallengeRecipientsAndObserversAgainstChallengingFiles(ctx, sliceOfFileHashesToChallenge, merkleroot, challengingSupernodeID)
	logger.Debug("challenge recipients and partial observers have been identified")

	for idx2, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		b, err := task.GetSymbolFileByKey(ctx, currentFileHashToChallenge, true)
		if err != nil {
			logger.WithError(err).Errorf("Supernode %s encountered an error generating storage challenges", challengingSupernodeID)
			continue
		}

		challengeDataSize := uint64(len(b))
		if challengeDataSize == 0 {
			logger.Errorf("Supernode %s encountered an invalid file while attempting to generate a storage challenge for file hash %s", challengingSupernodeID, currentFileHashToChallenge)
			continue
		}

		messageType := types.ChallengeMessageType
		respondingSupernodeID := sliceOfSupernodesToChallenge[idx2]

		// Identify challenge slice indices
		challengeSliceStartIndex, challengeSliceEndIndex := getStorageChallengeSliceIndices(challengeDataSize, currentFileHashToChallenge, merkleroot, challengingSupernodeID)

		//Create Challenge And Message IDs
		challengeIDInputData := challengingSupernodeID + respondingSupernodeID + currentFileHashToChallenge + fmt.Sprint(challengeSliceStartIndex) + fmt.Sprint(challengeSliceEndIndex) + fmt.Sprint(merkleroot) + uuid.NewString()
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
					Timestamp:  time.Now().UTC(),
					FileHash:   currentFileHashToChallenge,
					StartIndex: challengeSliceStartIndex,
					EndIndex:   challengeSliceEndIndex,
				},
			},
			Sender:          challengingSupernodeID,
			SenderSignature: nil,
		}

		if err := task.StoreStorageChallengeMetric(ctx, storageChallengeMessage); err != nil {
			logger.WithField("challenge_id", challengeID).
				WithField("message_type", storageChallengeMessage.MessageType).Error(
				"error storing storage challenge metric")
		}

		// Send the storage challenge message for processing by the responder
		if err = task.sendProcessStorageChallenge(ctx, storageChallengeMessage); err != nil {
			logger.WithError(err).WithField("challenge_id", storageChallengeMessage.ChallengeID).
				Error("Error processing storage challenge: ")
		}

		logger.WithField("challenge_id", challengeID).WithField("file_hash", currentFileHashToChallenge).
			Info("challenge sent for processing to recipient & observers")
	}

	logger.Info("files have been challenged")
	return nil
}

// sendProcessStorageChallenge will send the storage challenge message to the responding node
func (task *SCTask) sendProcessStorageChallenge(ctx context.Context, challengeMessage types.Message) error {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send process storage challenge")
		return err
	}

	if nodesToConnect == nil {
		return errors.Errorf("no nodes found to connect for send process storage challenge")
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

	var wg sync.WaitGroup
	var recipientNode pastel.MasterNode
	for _, node := range nodesToConnect {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			if challengeMessage.Data.RecipientID == node.ExtKey {
				recipientNode = node
				return
			}
			logger := log.WithContext(ctx).WithField("node_address", node.ExtAddress)

			if err := task.SendMessage(ctx, &msg, node.ExtAddress); err != nil {
				logger.WithError(err).Error("error sending generate storage challenge message for processing")
				return
			}

			logger.Debug("challenge message has been sent")
		}()
	}
	wg.Wait()
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).Debug("challenge message has been sent to observers")

	if err := task.SendMessage(ctx, &msg, recipientNode.ExtAddress); err != nil {
		log.WithContext(ctx).WithField("node_address", recipientNode.ExtAddress).WithError(err).Error("error sending storage challenge message to recipient for processing")
		return err
	}
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).
		WithField("node_address", recipientNode.ExtAddress).Debug("challenge message has been sent to recipient")

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
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeID).WithField("method", "sendProcessStorageChallenge").
				WithField("node_id", mn.ExtKey).Debug("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	logger := log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).
		WithField("message_type", challengeMessage.MessageType).WithField("node_id", task.nodeID)

	switch challengeMessage.MessageType {
	case types.ChallengeMessageType:
		//send message to partial observers;
		for _, po := range challengeMessage.Data.Observers {
			if value, ok := mapSupernodes[po]; ok {
				nodesToConnect = append(nodesToConnect, value)
			} else {
				logger.WithField("observer_id", po).Info("observer not found in masternode list")
			}
		}

		// and challenge recipient
		nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

		logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
		return nodesToConnect, nil
	case types.ResponseMessageType:
		//should process and send by recipient only
		if challengeMessage.Data.RecipientID == task.nodeID {
			//send message to partial observers for verification
			for _, po := range challengeMessage.Data.Observers {
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				} else {
					logger.WithField("observer_id", po).Info("observer not found in masternode list")
				}
			}

			//and challenger
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.ChallengerID])

			logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
			return nodesToConnect, nil
		}

		return nil, errors.Errorf("no recipients found to send this message")
	case types.EvaluationMessageType:
		//should process and send by challenger only
		if challengeMessage.Data.ChallengerID == task.nodeID {
			//send message to challenge recipient
			nodesToConnect = append(nodesToConnect, mapSupernodes[challengeMessage.Data.RecipientID])

			//and observers
			for _, po := range challengeMessage.Data.Observers {
				if value, ok := mapSupernodes[po]; ok {
					nodesToConnect = append(nodesToConnect, value)
				} else {
					logger.WithField("observer_id", po).Info("observer not found in masternode list")
				}
			}

			logger.WithField("nodes_to_connect", nodesToConnect).Debug("nodes to send msg have been selected")
			return nodesToConnect, nil
		}

		return nil, errors.Errorf("no recipients found to send this message")
	case types.AffirmationMessageType:
	// not required

	case types.BroadcastMessageType:
		return supernodes, nil

	default:
		return nil, errors.Errorf("no nodes found to send message")
	}

	return nil, err
}

// SendMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *SCTask) SendMessage(ctx context.Context, challengeMessage *pb.StorageChallengeMessage, processingSupernodeAddr string) error {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Debug("Sending storage challenge to processing supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendProcessStorageChallenge").Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	switch challengeMessage.MessageType {
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_CHALLENGE_MESSAGE:
		return storageChallengeIF.ProcessStorageChallenge(ctx, challengeMessage)
	case pb.StorageChallengeMessage_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE:
		return storageChallengeIF.VerifyStorageChallenge(ctx, challengeMessage)

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
	log.WithContext(ctx).WithField("len_slice_of_file_hashes", len(sliceOfFileHashes)).Debug("identifying which files are currently hosted on this node")

	totalLen := len(sliceOfFileHashes)
	sliceOfFileHashesStoredByLocalSupernode = make([]string, 0, totalLen)
	for i := 0; i < len(sliceOfFileHashes); i++ {
		if value, ok := task.SCService.localKeys.Load(sliceOfFileHashes[i]); ok {
			stored, isOk := value.(bool)
			if !isOk {
				log.WithContext(ctx).WithField("file_hash", sliceOfFileHashes[i]).Error("could not convert stored value to bool")
				continue
			}

			if !stored {
				log.WithContext(ctx).WithField("file_hash", sliceOfFileHashes[i]).Debug("file hash is in map but not stored")
				continue
			}

			sliceOfFileHashesStoredByLocalSupernode = append(sliceOfFileHashesStoredByLocalSupernode, sliceOfFileHashes[i])
		}
	}

	log.WithContext(ctx).WithField("total-possible-keys", totalLen).WithField("keys-found-len", len(sliceOfFileHashesStoredByLocalSupernode)).
		Info("files hosted on this node have been identified")

	return sliceOfFileHashesStoredByLocalSupernode
}

func (task SCTask) getChallengers(ctx context.Context, merkleRoot string) (sliceOfChallengers []string, challengersPerBlock, challengesPerBlock int, err error) {
	log.WithContext(ctx).Debug("retrieving list of all super nodes")
	listOfSupernodes, err := task.GetListOfSupernode(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get list of Supernode using pastel client")
		return nil, 0, 0, err
	}

	numberOfSupernodesDividedByThree := int(math.Ceil(float64(len(listOfSupernodes)) / 3))
	// Calculate the number of supernodes to issue challenge per block
	challengersPerBlock = numberOfSupernodesDividedByThree
	// Calculate the challenges per challenger
	challengesPerBlock = numberOfSupernodesDividedByThree

	// Identify challengers for this block
	sliceOfChallengers = task.GetNClosestSupernodeIDsToComparisonString(ctx, challengersPerBlock, merkleRoot, listOfSupernodes)

	return sliceOfChallengers, challengersPerBlock, challengesPerBlock, nil
}

func (task SCTask) getChallengingFiles(ctx context.Context, merkleRoot, challengerNodeID string, challengesPerBlock int, sliceOfFileHashesStoredByLocalSupernode []string) (sliceOfFileHashesToChallenge []string) {
	comparisonStringForFileHashSelection := merkleRoot + challengerNodeID
	comparisonStringHashForChallengeFileSelection := utils.GetHashFromString(comparisonStringForFileHashSelection)

	log.WithContext(ctx).Debug(fmt.Sprintf("comparison string hash for challenge file selection:%s", comparisonStringHashForChallengeFileSelection))
	sliceOfFileHashesToChallenge = task.GetNClosestFileHashesToAGivenComparisonString(ctx, challengesPerBlock, comparisonStringHashForChallengeFileSelection, sliceOfFileHashesStoredByLocalSupernode)
	log.WithContext(ctx).Debug(fmt.Sprintf("total files to challenge are:%d", len(sliceOfFileHashesToChallenge)))

	return sliceOfFileHashesToChallenge

}

func (task SCTask) identifyChallengeRecipientsAndObserversAgainstChallengingFiles(ctx context.Context, sliceOfFileHashesToChallenge []string, merkleRoot, challengerNodeID string) ([]string, map[int][]string) {
	mapOfchallengeFilePartialObservers := make(map[int][]string)
	sliceOfSupernodesToChallenge := make([]string, len(sliceOfFileHashesToChallenge))

	// Identify which supernodes have our file to challenge
	for idx1, currentFileHashToChallenge := range sliceOfFileHashesToChallenge {
		log.WithContext(ctx).WithField("no_of_challenge_replicas", defaultNumberOfChallengeReplicas).Debug("no of closest nodes to find against the file hash")
		sliceOfSupernodesStoringFileHashExcludingChallenger := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, defaultNumberOfChallengeReplicas, currentFileHashToChallenge, challengerNodeID)
		log.WithContext(ctx).WithField("slice_of_supernodes_ex_challenger_len", len(sliceOfSupernodesStoringFileHashExcludingChallenger)).Debug("slice of closest sns excluding challenger have been retrieved")

		comparisonStringForSupernodeSelection := merkleRoot + currentFileHashToChallenge + challengerNodeID + utils.GetHashFromString(fmt.Sprint(idx1))

		respondingSupernodeIDs := task.GetNClosestSupernodeIDsToComparisonString(ctx, 1, comparisonStringForSupernodeSelection, sliceOfSupernodesStoringFileHashExcludingChallenger)
		if len(respondingSupernodeIDs) < 1 {
			log.WithContext(ctx).WithField("file_hash", currentFileHashToChallenge).Debug("no closest nodes have found against the file")
			continue
		}

		sliceOfSupernodesToChallenge[idx1] = respondingSupernodeIDs[0] //challenge recipient
		log.WithContext(ctx).WithField("challenge_recipient", respondingSupernodeIDs[0]).Debug("challenge recipient has been selected")

		//partial observers
		for _, sn := range sliceOfSupernodesStoringFileHashExcludingChallenger {
			if sn == respondingSupernodeIDs[0] {
				continue
			}

			mapOfchallengeFilePartialObservers[idx1] = append(mapOfchallengeFilePartialObservers[idx1], sn)
		}

		log.WithContext(ctx).WithField("partial_observers_map", mapOfchallengeFilePartialObservers).Debug("partial observers have been selected")
	}

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
		log.WithContext(ctx).Debug("store")
		storageChallengeLog := types.StorageChallengeLogMessage{
			ChallengeID:     msg.ChallengeID,
			MessageType:     int(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
		}

		err = store.InsertStorageChallengeMessage(storageChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Debug("message already exists, not updating")

				return nil
			}

			log.WithContext(ctx).WithError(err).Error("Error storing challenge message to DB")
			return err
		}
	}

	return nil
}

// StoreStorageChallengeMetric stores the challenge message to db for further verification
func (task *SCTask) StoreStorageChallengeMetric(ctx context.Context, msg types.Message) error {
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
		log.WithContext(ctx).Debug("store")
		storageChallengeLog := types.StorageChallengeMetric{
			ChallengeID: msg.ChallengeID,
			MessageType: int(msg.MessageType),
			Data:        data,
			SenderID:    msg.Sender,
		}

		err = store.InsertStorageChallengeMetric(storageChallengeLog)
		if err != nil {
			if strings.Contains(err.Error(), ErrUniqueConstraint.Error()) {
				log.WithContext(ctx).WithField("challenge_id", msg.ChallengeID).
					WithField("message_type", msg.MessageType).
					WithField("sender_id", msg.Sender).
					Debug("message already exists, not updating")

				return nil
			}

			return err
		}
	}

	return nil
}

// SignMessage signs the message using sender's pastelID and passphrase
func (task SCTask) SignMessage(ctx context.Context, data types.MessageData) (sig []byte, dat []byte, err error) {
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
