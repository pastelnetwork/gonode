package storagechallenge

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"sync"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
	"golang.org/x/crypto/sha3"
)

// ProcessStorageChallenge consists of:
//
//	Getting the file,
//	Hashing the indicated portion of the file ("responding"),
//	Identifying the proper validators,
//	Sending the response to all other supernodes
//	Saving challenge state
func (task *SCTask) ProcessStorageChallenge(ctx context.Context, incomingChallengeMessage types.Message) (*pb.StorageChallengeMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "ProcessStorageChallenge").
		WithField("sc_challenge_id", incomingChallengeMessage.ChallengeID)

	// incoming challenge message validation
	if err := task.validateProcessingStorageChallengeIncomingData(ctx, incomingChallengeMessage); err != nil {
		logger.WithError(err).Error("Error validating storage challenge incoming data: ")
		return nil, err
	}
	logger.Debug("Incoming challenge validated")

	//if the message is received by one of the observer then save the challenge message, lock the file & return
	if task.isObserver(incomingChallengeMessage.Data.Observers) {
		if err := task.SCService.P2PClient.DisableKey(ctx, incomingChallengeMessage.Data.Challenge.FileHash); err != nil {
			logger.WithError(err).Error("error locking the file")
		}

		if err := task.StoreChallengeMessage(ctx, incomingChallengeMessage); err != nil {
			logger.WithError(err).Error("error storing challenge message")
		}
		logger.WithField("node_id", task.nodeID).Debug("challenge message has been stored by the observer")

		return nil, nil
	}

	//if not the challenge recipient, should return, otherwise proceed
	if task.nodeID != incomingChallengeMessage.Data.RecipientID {
		logger.Debug("current node is not the challenge recipient to process challenge message")

		return nil, nil
	}

	logger.Info("Storage challenge processing started") // Incoming challenge message validation
	if err := task.StoreChallengeMessage(ctx, incomingChallengeMessage); err != nil {
		logger.WithError(err).Error("error storing challenge message")
	}

	// Get the file to hash
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingChallengeMessage.Data.Challenge.FileHash, true)
	if err != nil {
		logger.WithError(err).Error("could not read file data in to memory")

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
					Block:      incomingChallengeMessage.Data.Challenge.Block,
					Merkelroot: incomingChallengeMessage.Data.Challenge.Merkelroot,
					Hash:       "",
					Timestamp:  time.Now().UTC(),
				},
			},
			Sender:          task.nodeID,
			SenderSignature: nil,
		}

		if err := task.StoreStorageChallengeMetric(ctx, outgoingResponseMessage); err != nil {
			logger.WithError(err).Error("error storing storage challenge metric")
		}

		if err = task.sendVerifyStorageChallenge(ctx, outgoingResponseMessage); err != nil {
			logger.WithError(err).Error("could not send processed challenge message to node for verification")
			return nil, err
		}

		return nil, nil
	}

	// Get the hash of the chunk of the file we're supposed to hash
	challengeResponseHash := task.computeHashOfFileSlice(challengeFileData, incomingChallengeMessage.Data.Challenge.StartIndex, incomingChallengeMessage.Data.Challenge.EndIndex)
	logger.Info("hash for data generated against the indices")

	blockNumChallengeRespondedTo, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get current block count")
		return nil, err
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeRespondedTo)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose 1")
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
				Timestamp:  time.Now().UTC(),
			},
		},
		Sender:          task.nodeID,
		SenderSignature: nil,
	}

	if err := task.StoreStorageChallengeMetric(ctx, outgoingResponseMessage); err != nil {
		logger.WithError(err).Error("error storing storage challenge metric")
	}

	if err = task.sendVerifyStorageChallenge(ctx, outgoingResponseMessage); err != nil {
		logger.WithError(err).Error("could not send processed challenge message to node for verification")
		return nil, err
	}
	logger.Info("message sent to other SNs for verification")

	return nil, nil
}

func (task *SCTask) validateProcessingStorageChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ChallengeMessageType {
		return fmt.Errorf("incorrect message type to processing storage challenge")
	}

	if err := task.StoreStorageChallengeMetric(ctx, incomingChallengeMessage); err != nil {
		log.WithContext(ctx).WithField("challenge_id", incomingChallengeMessage.ChallengeID).
			WithField("message_type", incomingChallengeMessage.MessageType).Error(
			"error storing storage challenge metric")
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
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send verify storage challenge")
		return err
	}

	if nodesToConnect == nil {
		return errors.Errorf("no nodes found to connect to send verify storage challenge")
	}

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the challenge message")
		return err
	}
	challengeMessage.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, challengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing the response message")
	}

	msg := pb.StorageChallengeMessage{
		MessageType:     pb.StorageChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}

	var challengerNode pastel.MasterNode
	var wg sync.WaitGroup
	for _, node := range nodesToConnect {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			if challengeMessage.Data.ChallengerID == node.ExtKey {
				challengerNode = node
				return
			}

			logger := log.WithContext(ctx).WithField("node_address", node.ExtAddress)

			if err := task.SendMessage(ctx, &msg, node.ExtAddress); err != nil {
				logger.WithError(err).Error("error sending process storage challenge message for processing")
				return
			}
		}()
	}
	wg.Wait()
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).Debug("response message has been sent to observers")

	if err := task.SendMessage(ctx, &msg, challengerNode.ExtAddress); err != nil {
		log.WithContext(ctx).WithField("node_address", challengerNode.ExtAddress).WithError(err).Error("error sending response message to challenger for verification")
		return err
	}
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeID).
		WithField("node_address", challengerNode.ExtAddress).Debug("response message has been sent to challenger")

	return nil
}

func (task *SCTask) isObserver(sliceOfObservers []string) bool {
	for _, ob := range sliceOfObservers {
		if ob == task.nodeID {
			return true
		}
	}

	return false
}
