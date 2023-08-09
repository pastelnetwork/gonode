package storagechallenge

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"golang.org/x/crypto/sha3"
	"strings"
	"time"
)

// ProcessStorageChallenge consists of:
//
//	Getting the file,
//	Hashing the indicated portion of the file ("responding"),
//	Identifying the proper validators,
//	Sending the response to all other supernodes
//	Saving challenge state
func (task *SCTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage types.Message) (*pb.StorageChallengeMessage, error) {
	// incoming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating storage challenge incoming data: ")
		return nil, err
	}
	log.WithContext(ctx).WithField("incoming_challenge", incomingChallengeMessage).Info("Incoming challenge validated")

	log.WithContext(ctx).Info("retrieving the file from hash")
	// Get the file to hash
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.Data.Challenge.FileHash, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not read file data in to memory")
		return nil, err
	}
	log.WithContext(ctx).Info("challenge file has been retrieved")

	// Get the hash of the chunk of the file we're supposed to hash
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeResponseHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.Data.Challenge.StartIndex, incomingChallengeMessage.Data.Challenge.EndIndex)
	log.WithContext(ctx).Info(fmt.Sprintf("hash for data generated against the indices:%s", challengeResponseHash))

	log.WithContext(ctx).Info("sending message to other SNs for verification")
	messageType := types.ResponseMessageType
	blockNumChallengeRespondedTo, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not get current block count")
		return nil, err
	}
	log.WithContext(ctx).Info(fmt.Sprintf("block num challenge responded to:%d", blockNumChallengeRespondedTo))

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeRespondedTo)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}

	//Create the response message to be validated
	outgoingResponseMessage := types.Message{
		MessageType: messageType,
		ChallengeID: incomingChallengeMessage.ChallengeID,
		Data: types.MessageData{
			ChallengerID: incomingChallengeMessage.Data.ChallengerID,
			Challenge: types.ChallengeData{
				Block:      incomingChallengeMessage.Data.Challenge.Block,
				Merkelroot: incomingChallengeMessage.Data.Challenge.Merkelroot,
				Timestamp:  incomingChallengeMessage.Data.Challenge.Timestamp,
				FileHash:   incomingChallengeMessage.Data.Challenge.FileHash,
				StartIndex: incomingChallengeMessage.Data.Challenge.StartIndex,
				EndIndex:   incomingChallengeMessage.Data.Challenge.EndIndex,
			},
			Observers:   append([]string(nil), incomingChallengeMessage.Data.Observers...),
			RecipientID: incomingChallengeMessage.Data.RecipientID,
			Response: types.ResponseData{
				Block:      blockNumChallengeRespondedTo,
				Merkelroot: blkVerbose1.MerkleRoot,
				Hash:       challengeResponseHash,
				Timestamp:  time.Now(),
			},
		},
	}

	// send to Supernodes to validate challenge response hash
	log.WithContext(ctx).WithField("challenge_id", outgoingResponseMessage.ChallengeID).Info("sending challenge response for verification")
	if err = task.sendVerifyStorageChallenge(ctx, outgoingResponseMessage); err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not send processed challenge message to node for verification")
		return nil, err
	}
	log.WithContext(ctx).Info("message sent to other SNs for verification")

	task.SaveChallengeMessageState(
		ctx,
		"respond",
		outgoingResponseMessage.ChallengeID,
		outgoingResponseMessage.Data.ChallengerID,
		outgoingResponseMessage.Data.Response.Block,
	)

	return nil, nil
}

func (task *SCTask) validateProcessingStorageChallengeIncomingData(incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ChallengeMessageType {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}

	return nil
}

func (task *SCTask) computeHashOfFileSlice(fileData []byte, challengeSliceStartIndex, challengeSliceEndIndex int) string {
	if len(fileData) < challengeSliceStartIndex || len(fileData) < challengeSliceEndIndex {
		return ""
	}

	challengeDataSlice := fileData[challengeSliceStartIndex:challengeSliceEndIndex]
	algorithm := sha3.New256()
	algorithm.Write(challengeDataSlice)
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Send our verification message to (default 10) other supernodes that might host this file.
func (task *SCTask) sendVerifyStorageChallenge(ctx context.Context, challengeMessage types.Message) error {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send process storage challenge")
		return err
	}

	data, err := json.Marshal(challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error marshaling the data")
	}

	msg := pb.StorageChallengeMessage{
		MessageType:     pb.StorageChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}

	for _, node := range nodesToConnect {
		if err := task.SendMessage(ctx, msg, node.ExtAddress); err != nil {
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
			ChallengeID:     challengeMessage.ChallengeID,
			FileHash:        challengeMessage.Data.Challenge.FileHash,
			ChallengingNode: challengeMessage.Data.ChallengerID,
			RespondingNode:  task.nodeID,
			Status:          types.ProcessedStorageChallengeStatus,
			StartingIndex:   int(challengeMessage.Data.Challenge.StartIndex),
			EndingIndex:     int(challengeMessage.Data.Challenge.EndIndex),
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
