package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/common/utils"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"golang.org/x/crypto/sha3"
	"strings"
)

// ProcessStorageChallenge consists of:
//
//	Getting the file,
//	Hashing the indicated portion of the file ("responding"),
//	Identifying the proper validators,
//	Sending the response to all other supernodes
//	Saving challenge state
func (task *SCTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage *pb.StorageChallengeData) (*pb.StorageChallengeData, error) {

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
	_ = outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent

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
	if len(fileData) < int(challengeSliceStartIndex) || len(fileData) < int(challengeSliceEndIndex) {
		return ""
	}

	challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
	algorithm := sha3.New256()
	algorithm.Write(challengeDataSlice)
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Send our verification message to (default 10) other supernodes that might host this file.
func (task *SCTask) sendVerifyStorageChallenge(ctx context.Context, challengeMessage *pb.StorageChallengeData) error {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, *challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send process storage challenge")
		return err
	}

	for _, node := range nodesToConnect {
		if err := task.SendMessage(ctx, *challengeMessage, node.ExtAddress); err != nil {
			log.WithContext(ctx).WithError(err).Error("error sending storage challenge message for processing")
			continue
		}
	}

	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
	}
	if store != nil {
		store.CloseHistoryDB(ctx)
	}

	if store != nil {
		log.WithContext(ctx).Println("Storing challenge logs to DB")
		storageChallengeLog := types.StorageChallenge{
			ChallengeID:     challengeMessage.ChallengeId,
			FileHash:        challengeMessage.ChallengeFile.FileHashToChallenge,
			ChallengingNode: challengeMessage.ChallengingMasternodeId,
			RespondingNode:  task.nodeID,
			Status:          types.ProcessedStorageChallengeStatus,
			StartingIndex:   int(challengeMessage.ChallengeFile.ChallengeSliceStartIndex),
			EndingIndex:     int(challengeMessage.ChallengeFile.ChallengeSliceEndIndex),
		}

		var verifyingNodes []string
		for _, vn := range nodesToConnect {
			verifyingNodes = append(verifyingNodes, vn.ExtKey)
		}
		storageChallengeLog.VerifyingNodes = strings.Join(verifyingNodes, ",")

		_, err = store.InsertStorageChallenge(storageChallengeLog)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error storing challenge log to DB")
			err = nil
		}
	}

	return nil
}
