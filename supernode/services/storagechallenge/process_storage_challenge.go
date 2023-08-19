package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"golang.org/x/crypto/sha3"
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
	log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").
		WithField("challengeID", incomingChallengeMessage.ChallengeID).
		Debug("Start processing storage challenge") // Incoming challenge message validation

	// incoming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(ctx, incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("Error validating storage challenge incoming data: ")
		return nil, err
	}
	log.WithContext(ctx).WithField("incoming_challenge", incomingChallengeMessage).Info("Incoming challenge validated")

	if incomingChallengeMessage.Data.RecipientID != task.nodeID { //if observers then save the challenge message & return
		if err := task.StoreChallengeMessage(ctx, incomingChallengeMessage); err != nil {
			log.WithContext(ctx).
				WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing challenge message")
		}

		return nil, nil
	}

	// Get the file to hash
	log.WithContext(ctx).Info("retrieving the file from hash")
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
		MessageType: types.ResponseMessageType,
		ChallengeID: incomingChallengeMessage.ChallengeID,
		Data: types.MessageData{
			ChallengerID: incomingChallengeMessage.Data.ChallengerID,
			RecipientID:  incomingChallengeMessage.Data.RecipientID,
			Observers:    append([]string(nil), incomingChallengeMessage.Data.Observers...),
			Challenge: types.ChallengeData{
				Block:      incomingChallengeMessage.Data.Challenge.Block,
				Merkelroot: incomingChallengeMessage.Data.Challenge.Merkelroot,
				Timestamp:  incomingChallengeMessage.Data.Challenge.Timestamp,
				FileHash:   incomingChallengeMessage.Data.Challenge.FileHash,
				StartIndex: incomingChallengeMessage.Data.Challenge.StartIndex,
				EndIndex:   incomingChallengeMessage.Data.Challenge.EndIndex,
			},
			Response: types.ResponseData{
				Block:      blockNumChallengeRespondedTo,
				Merkelroot: blkVerbose1.MerkleRoot,
				Hash:       challengeResponseHash,
				Timestamp:  time.Now(),
			},
		},
		Sender:          task.nodeID,
		SenderSignature: nil,
	}

	// send to Supernodes to validate challenge response hash
	log.WithContext(ctx).WithField("challenge_id", outgoingResponseMessage.ChallengeID).Info("sending challenge response for verification")
	if err = task.sendVerifyStorageChallenge(ctx, outgoingResponseMessage); err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not send processed challenge message to node for verification")
		return nil, err
	}
	log.WithContext(ctx).Info("message sent to other SNs for verification")

	task.SaveChallengeMessageState(ctx, "respond", outgoingResponseMessage.ChallengeID, outgoingResponseMessage.Data.ChallengerID, outgoingResponseMessage.Data.Response.Block)

	return nil, nil
}

func (task *SCTask) validateProcessingStorageChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ChallengeMessageType {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}

	isVerified, err := task.VerifyMessageSignature(ctx, incomingChallengeMessage)
	if err != nil {
		return errors.Errorf("error verifying sender's signature: %w", err)
	}

	if !isVerified {
		return errors.Errorf("not able to verify message signature")
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

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the challenge message")
		return err
	}
	challengeMessage.SenderSignature = signature

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

	return nil
}
