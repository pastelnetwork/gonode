package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"golang.org/x/crypto/sha3"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

//ProcessStorageChallenge consists of:
//	 Getting the file,
//   Hashing the indicated portion of the file ("responding"),
//   Identifying the proper validators,
//   Sending the response to all other supernodes
//	 Saving challenge state
func (task *SCTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) (*pb.StorageChallengeData, error) {
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Start processing storage challenge")

	// incoming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		return nil, err
	}

	// Get the file to hash
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.ChallengeFile.FileHashToChallenge, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not read file data in to memory")
		return nil, err
	}
	// Get the hash of the chunk of the file we're supposed to hash
	challengeResponseHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex, incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex)
	challengeStatus := pb.StorageChallengeData_Status_RESPONDED
	messageType := pb.StorageChallengeData_MessageType_STORAGE_CHALLENGE_RESPONSE_MESSAGE
	messageIDInputData := incomingChallengeMessage.ChallengingMasternodeId + incomingChallengeMessage.RespondingMasternodeId + incomingChallengeMessage.ChallengeFile.FileHashToChallenge + challengeStatus.String() + messageType.String() + incomingChallengeMessage.MerklerootWhenChallengeSent
	messageID := utils.GetHashFromString(messageIDInputData)
	blockNumChallengeRespondedTo, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not get current block count")
		return nil, err
	}

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
			ChallengeSliceStartIndex: int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceStartIndex),
			ChallengeSliceEndIndex:   int64(incomingChallengeMessage.ChallengeFile.ChallengeSliceEndIndex),
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
	if err = task.sendVerifyStorageChallenge(ctx, outgoingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeId).Error("could not send processed challenge message to verifying node")
		return nil, err
	}

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
	//Get the <default 10> closest super node id's to the file hash.
	sliceOfSupernodesStoringFileHashExcludingChallenger := task.GetNClosestSupernodesToAGivenFileUsingKademlia(ctx, task.numberOfVerifyingNodes, challengeMessage.ChallengeFile.FileHashToChallenge, task.nodeID)
	//turn this into a map so we don't have to do 10n iterations through supernodes in case there are lots of supernodes
	mapOfSupernodesStoringFHEC := make(map[string]bool)
	// using bool instead of empty struct to make evaluation below easier to understand
	for _, s := range sliceOfSupernodesStoringFileHashExcludingChallenger {
		mapOfSupernodesStoringFHEC[s] = true
	}

	Supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	SupernodeAddressesIgnoringThisNode := []string{}
	for _, mn := range Supernodes {
		//only find supernodes that aren't this one and that are within the ten closest to the file hash
		if mn.ExtKey != task.nodeID && mapOfSupernodesStoringFHEC[mn.ExtKey] {
			SupernodeAddressesIgnoringThisNode = append(SupernodeAddressesIgnoringThisNode, mn.ExtAddress)
		}
	}

	err = nil
	// iterate through supernodes, connecting and sending the message
	for _, nodeToConnectTo := range SupernodeAddressesIgnoringThisNode {
		nodeClientConn, err := task.nodeClient.Connect(ctx, nodeToConnectTo)
		if err != nil {
			err = fmt.Errorf("Could not use nodeclient to connect to: " + nodeToConnectTo)
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("method", "sendprocessStorageChallenge").Warn(err.Error())
			return err
		}
		storageChallengeIF := nodeClientConn.StorageChallenge()
		//Sends the verify storage challenge message to the connected verifying supernode
		err = storageChallengeIF.VerifyStorageChallenge(ctx, challengeMessage)
		if err != nil {
			log.WithContext(ctx).WithField("challengeID", challengeMessage.ChallengeId).WithField("verifierSuperNodeAddress", nodeToConnectTo).Warn("Storage challenge verification failed or didn't process: ", err.Error())
		}
	}
	if err == nil {
		log.WithContext(ctx).Println("After calling storage process on " + challengeMessage.ChallengeId + " no nodes returned an error code in verification")
	}
	return err
	//return s.actor.Send(ctx, s.domainActorID, newSendVerifyStorageChallengeMsg(ctx, verifierSupernodesAddr, challengeMessage))
}
