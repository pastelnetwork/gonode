package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"golang.org/x/crypto/sha3"
)

const (
	//ChallengeFailuresThreshold is a threshold which when exceeds will send the challenge for self-healing
	ChallengeFailuresThreshold = 5
)

// ProcessStorageChallenge consists of:
//
//	Getting the file,
//	Hashing the indicated portion of the file ("responding"),
//	Identifying the proper validators,
//	Sending the response to all other supernodes
//	Saving challenge state
func (task *SCTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) (*pb.StorageChallengeData, error) {
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Start processing storage challenge")

	// incoming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating storage challenge incoming data: ")
		return nil, err
	}
	log.WithContext(ctx).WithField("incoming_challenge", incomingChallengeMessage).Info("Incoming challenge validated")

	log.WithContext(ctx).Info("retrieving the file from hash")
	// Get the file to hash
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.ChallengeFile.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not read file data in to memory")
		return nil, err
	}
	log.WithContext(ctx).Info("challenge file has been retrieved")

	// Get the hash of the chunk of the file we're supposed to hash
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeResponseHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex, incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex)
	log.WithContext(ctx).Info(fmt.Sprintf("hash for data generated against the indices:%s", challengeResponseHash))

	log.WithContext(ctx).Info("sending message to other SNs for verification")
	challengeStatus := pb.StorageChallengeData_Status_RESPONDED
	messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeId + incomingChallengeMessage.RespondingMasternodeId + incomingChallengeMessage.ChallengeFile.FileHashToChallenge + challengeStatus.String() + messageType.String() + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)
	blockNumChallengeRespondedTo, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not get current block count")
		return nil, err
	}
	log.WithContext(ctx).Info(fmt.Sprintf("block num challenge responded to:%d", blockNumChallengeRespondedTo))

	log.WithContext(ctx).Info(fmt.Sprintf("NodeID:%s", task.nodeID))
	//Create the message to be validated
	outgoingChallengeMessage := &pb.StorageChallengeData{
		MessageId:                    messageID,
		MessageType:                  messageType,
		ChallengeStatus:              challengeStatus,
		BlockNumChallengeSent:        incomingChallengeMessage.BlockNumChallengeSent,
		BlockNumChallengeRespondedTo: blockNumChallengeRespondedTo,
		BlockNumChallengeVerified:    0,
		MerklerootWhenChallengeSent:  incomingChallengeMessage.MerklerootWhenChallengeSent,
		ChallengingMasternodeId:      incomingChallengeMessage.ChallengingMasternodeId,
		RespondingMasternodeId:       task.nodeID,
		ChallengeFile: &pb.StorageChallengeDataChallengeFile{
			FileHashToChallenge:      incomingChallengeMessage.ChallengeFile.FileHashToChallenge,
			ChallengeSliceStartIndex: incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex,
			ChallengeSliceEndIndex:   incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex,
		},
		ChallengeSliceCorrectHash: "",
		ChallengeResponseHash:     challengeResponseHash,
		ChallengeId:               incomingChallengeMessage.ChallengeId,
	}

	//Currently we write our own block when we responded to the storage challenge. Could be changed to calculated by verification node
	//	but that also has trade-offs.
	blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug(fmt.Sprintf("Supernode %s responded to storage challenge for file hash %s in %v blocks!", outgoingChallengeMessage.RespondingMasternodeId, outgoingChallengeMessage.ChallengeFile.FileHashToChallenge, blocksToRespondToStorageChallenge))

	// send to Supernodes to validate challenge response hash
	log.WithContext(ctx).WithField("challenge_id", outgoingChallengeMessage.ChallengeId).Info("sending challenge for verification")
	if err = task.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not send processed challenge message to node for verification")
		return nil, err
	}
	log.WithContext(ctx).Info("message sent to other SNs for verification")

	task.SaveChallengeMessageState(
		ctx,
		"respond",
		outgoingChallengeMessage.ChallengeId,
		outgoingChallengeMessage.ChallengingMasternodeId,
		outgoingChallengeMessage.BlockNumChallengeSent,
	)

	return outgoingChallengeMessage, err
}

func (task *SCTask) validateProcessingStorageChallengeIncomingData(incomingChallengeMessage *pb.StorageChallengeData) error {
	if incomingChallengeMessage.ChallengeStatus != pb.StorageChallengeData_Status_PENDING {
		return fmt.Errorf("incorrect status to processing storage challenge")
	}
	if incomingChallengeMessage.MessageType != pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_ISSUANCE_MESSAGE {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}
	return nil
}

func (task *SCTask) computeHashOfFileSlice(fileData []byte, challengeSliceStartIndex, challengeSliceEndIndex int64) string {
	challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
	algorithm := sha3.New256()
	algorithm.Write(challengeDataSlice)
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Send our verification message to (default 10) other supernodes that might host this file.
func (task *SCTask) sendVerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error {
	//Get the full list of supernodes
	listOfSupernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendProcessStorageChallenge").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return err
	}
	log.WithContext(ctx).Info(fmt.Sprintf("list of supernodes have been retrieved for process storage challenge:%s", listOfSupernodes))

	//Filtering out the current node from the full list of supernodes
	//turn this into a map, so we don't have to do n iterations through supernodes in case there are lots of supernodes
	mapSupernodesWithoutCurrentNode := make(map[string]pastel.MasterNode)
	//list of supernode ext keys without current node to find the 10 closest nodes
	var sliceOfSupernodeKeysExceptCurrentNode []string
	for _, mn := range listOfSupernodes {
		if mn.ExtKey != task.nodeID && mn.ExtKey != challengeMessage.ChallengingMasternodeId {
			mapSupernodesWithoutCurrentNode[mn.ExtKey] = mn
			sliceOfSupernodeKeysExceptCurrentNode = append(sliceOfSupernodeKeysExceptCurrentNode, mn.ExtKey)
		}
	}
	log.WithContext(ctx).Info(fmt.Sprintf("current node has been filtered out from the supernodes list:%s", sliceOfSupernodeKeysExceptCurrentNode))

	//Finding 10 closest nodes to file hash
	sliceOfSupernodesClosestToFileHashExcludingCurrentNode := task.GetNClosestSupernodeIDsToComparisonString(ctx, 10, challengeMessage.ChallengeFile.FileHashToChallenge, sliceOfSupernodeKeysExceptCurrentNode)
	log.WithContext(ctx).Info(fmt.Sprintf("sliceOfSupernodesClosestToFileHashExcludingCurrentNode:%s", sliceOfSupernodesClosestToFileHashExcludingCurrentNode))

	var (
		countOfFailures  int
		wg               sync.WaitGroup
		responseMessages []pb.StorageChallengeData
	)
	err = nil
	// iterate through supernodes, connecting and sending the message
	for _, nodeToConnectTo := range sliceOfSupernodesClosestToFileHashExcludingCurrentNode {
		log.WithContext(ctx).WithField("outgoing_message", challengeMessage).Info("outgoing message from ProcessStorageChallenge")

		var (
			mn pastel.MasterNode
			ok bool
		)
		if mn, ok = mapSupernodesWithoutCurrentNode[nodeToConnectTo]; !ok {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendVerifyStorageChallenge").Warn(fmt.Sprintf("cannot get Supernode info of Supernode id %s", mapSupernodesWithoutCurrentNode))
			continue
		}
		//We use the ExtAddress of the supernode to connect
		processingSupernodeAddr := mn.ExtAddress
		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("Sending storage challenge for verification to processing supernode address: " + processingSupernodeAddr)

		log.WithContext(ctx).Info(fmt.Sprintf("establishing connection with node: %s", processingSupernodeAddr))
		nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Connection failed to establish with node: %s", processingSupernodeAddr))
			continue
		}
		storageChallengeIF := nodeClientConn.StorageChallenge()
		log.WithContext(ctx).Info(fmt.Sprintf("connection established with node:%s", nodeToConnectTo))

		log.WithContext(ctx).Info(fmt.Sprintf("sending challenge message for verification to node:%s", nodeToConnectTo))
		//Sends the verify storage challenge message to the connected verifying supernode

		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := storageChallengeIF.VerifyStorageChallenge(ctx, challengeMessage)
			if err != nil {
				log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("verifierSuperNodeAddress", nodeToConnectTo).Warn("Storage challenge verification failed or didn't process: ", err.Error())
				return
			}

			responseMessages = append(responseMessages, *res)
			log.WithContext(ctx).WithField("challenge_id", res.ChallengeId).
				Info("response has been received from verifying node")
		}()
		wg.Wait()
	}
	if err == nil {
		log.WithContext(ctx).Println("After calling storage process on " + challengeMessage.ChallengeId + " no nodes returned an error code in verification")
	}

	//Counting the number of challenges being failed by verifying nodes.
	var responseMessage pb.StorageChallengeData
	for _, responseMessage = range responseMessages {
		if responseMessage.ChallengeStatus == pb.StorageChallengeData_Status_FAILED_INCORRECT_RESPONSE {
			countOfFailures++
		}
	}

	if countOfFailures >= ChallengeFailuresThreshold {
		closestSupernodeToMerkelRootForSelfHealingChallenge := task.GetNClosestSupernodeIDsToComparisonString(ctx, 1, challengeMessage.ChallengeFile.FileHashToChallenge, sliceOfSupernodeKeysExceptCurrentNode)
		log.WithContext(ctx).Info(fmt.Sprintf("closestSupernodeToMerkelRootForSelfHealingChallenge:%s", closestSupernodeToMerkelRootForSelfHealingChallenge))

		var (
			sn pastel.MasterNode
			ok bool
		)
		if len(closestSupernodeToMerkelRootForSelfHealingChallenge) > 0 {
			if sn, ok = mapSupernodesWithoutCurrentNode[closestSupernodeToMerkelRootForSelfHealingChallenge[0]]; !ok {
				log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendProcessSelfHealingChallenge").Warn(fmt.Sprintf("cannot get Supernode info of Supernode id %s", mapSupernodesWithoutCurrentNode))
			}
		}

		//We use the ExtAddress of the supernode to connect
		processingSupernodeAddr := sn.ExtAddress
		log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("Sending self-healing challenge for processing to supernode address: " + processingSupernodeAddr)

		log.WithContext(ctx).Info(fmt.Sprintf("establishing connection with node: %s", processingSupernodeAddr))
		nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error(fmt.Sprintf("Connection failed to establish with node: %s", processingSupernodeAddr))
		}
		selfHealingChallengeIF := nodeClientConn.SelfHealingChallenge()
		log.WithContext(ctx).Info(fmt.Sprintf("connection established with node:%s", processingSupernodeAddr))

		log.WithContext(ctx).Info(fmt.Sprintf("sending self-healing challenge message for processing to node:%s", processingSupernodeAddr))
		//Sends the verify storage challenge message to the connected verifying supernode

		////Storing to DB for inspection by Self healing
		//log.WithContext(ctx).Info(fmt.Sprintf("Storage challenge total no of failures exceeds than:%d", ChallengeFailuresThreshold))
		//
		//store, err := local.OpenHistoryDB()
		//if err != nil {
		//	log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		//}
		//defer store.CloseHistoryDB(ctx)
		//
		//log.WithContext(ctx).Println("Storing failed challenge to DB for self healing inspection")
		//failedChallenge := types.FailedStorageChallenge{
		//	ChallengeID:    responseMessage.ChallengeId,
		//	Status:         types.CreatedSelfHealingStatus.String(),
		//	FileHash:       challengeMessage.ChallengeFile.FileHashToChallenge,
		//	RespondingNode: task.nodeID,
		//}
		//
		//_, err = store.InsertFailedStorageChallenge(failedChallenge)
		//if err != nil {
		//	log.WithContext(ctx).WithError(err).Error("Error storing failed challenge to DB")
		//}

		data := &pb.SelfHealingData{
			MessageId: challengeMessage.MerklerootWhenChallengeSent + closestSupernodeToMerkelRootForSelfHealingChallenge[0] +
				challengeMessage.ChallengeFile.FileHashToChallenge,
			MessageType:                 pb.SelfHealingData_MessageType_SELF_HEALING_ISSUANCE_MESSAGE,
			ChallengeStatus:             pb.SelfHealingData_Status_PENDING,
			MerklerootWhenChallengeSent: challengeMessage.MerklerootWhenChallengeSent,
			ChallengingMasternodeId:     task.nodeID,
			RespondingMasternodeId:      sn.ExtKey,
			ChallengeFile: &pb.SelfHealingDataChallengeFile{
				FileHashToChallenge: challengeMessage.ChallengeFile.FileHashToChallenge,
			},
			ChallengeId: challengeMessage.ChallengeId,
		}

		if err := selfHealingChallengeIF.ProcessSelfHealingChallenge(ctx, data); err != nil {
			log.WithContext(ctx).WithError(err).Error("Error sending self-healing challenge for processing")
		}
	}

	return err
	//return s.actor.Send(ctx, s.domainActorID, newSendVerifyStorageChallengeMsg(ctx, verifierSupernodesAddr, challengeMessage))
}
