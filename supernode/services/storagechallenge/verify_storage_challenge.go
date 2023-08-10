package storagechallenge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/pastel"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"

	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

const (
	// SuccessfulEvaluationThreshold is the affirmation threshold required for broadcasting
	SuccessfulEvaluationThreshold = 2
)

// VerifyStorageChallenge : Verifying the storage challenge will occur only on the challenger node
//
//	 On receipt of challenge message we:
//			Validate it
//	 	Get the file assuming we host it locally (if not, return)
//	 	Compute the hash of the data at the indicated byte range
//			If the hash is correct and within the given byte range, success is indicated otherwise failure is indicated via SaveChallengeMessageState
func (task *SCTask) VerifyStorageChallenge(ctx context.Context, incomingResponseMessage types.Message) (*pb.StorageChallengeData, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID)
	logger.Debug("Start verifying storage challenge") // Incoming challenge message validation

	if err := task.validateVerifyingStorageChallengeIncomingData(ctx, incomingResponseMessage); err != nil {
		return nil, err
	}

	if incomingResponseMessage.Data.ChallengerID != task.nodeID { //if observers then save the response message & return
		if err := task.StoreChallengeMessage(ctx, incomingResponseMessage); err != nil {
			log.WithContext(ctx).
				WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing challenge message")
		}

		return nil, nil
	}

	//Get the file assuming we host it locally (if not, return)
	log.WithContext(ctx).Info("getting the file from hash to verify challenge")
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingResponseMessage.Data.Challenge.FileHash, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Error("could not read local file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())
		return nil, err
	}
	log.WithContext(ctx).Info("file has been retrieved for verification")

	//Compute the hash of the data at the indicated byte range
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingResponseMessage.Data.Challenge.StartIndex, incomingResponseMessage.Data.Challenge.EndIndex)
	log.WithContext(ctx).Info("hash of the data has been generated against the given indices")

	//Identify current block count to see if validation was performed within the mandated length of time (default one block)
	blockNumChallengeVerified, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", incomingResponseMessage.ChallengeID).Error("could not get current block count")
		return nil, err
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeVerified)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}

	blockNumChallengeSent := incomingResponseMessage.Data.Challenge.Block
	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - blockNumChallengeSent
	log.WithContext(ctx).WithField("no_of_blocks", blocksVerifyStorageChallengeInBlocks).Info("No of blocks from the time challenge being sent and verified")

	//Preparing the evaluation report
	log.WithContext(ctx).Info("determining the challenge outcome based on calculated and received hash")

	// determine success or failure
	var saveStatus string
	var isVerified bool
	if (incomingResponseMessage.Data.Response.Hash == challengeCorrectHash) && (blocksVerifyStorageChallengeInBlocks <= task.storageChallengeExpiredBlocks) {
		saveStatus = "succeeded"
		isVerified = true
		logger.Info(fmt.Sprintf("Supernode %s correctly responded in %d blocks to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, blocksVerifyStorageChallengeInBlocks, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "succeeded", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	} else if incomingResponseMessage.Data.Response.Hash == challengeCorrectHash {
		saveStatus = "timeout"
		logger.Info(fmt.Sprintf("Supernode %s  correctly responded in %d blocks to a storage challenge for file %s, but was too slow so failed the challenge anyway!", incomingResponseMessage.Data.RecipientID, blocksVerifyStorageChallengeInBlocks, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "timeout", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	} else {
		saveStatus = "failed"
		log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingResponseMessage.ChallengeID).Debug(fmt.Sprintf("Supernode %s failed by incorrectly responding to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, incomingResponseMessage.Data.Challenge.FileHash))
		task.SaveChallengeMessageState(ctx, "failed", incomingResponseMessage.ChallengeID, incomingResponseMessage.Data.ChallengerID, incomingResponseMessage.Data.Challenge.Block)
	}

	evaluationMessage := types.Message{
		MessageType: types.EvaluationMessageType,
		ChallengeID: incomingResponseMessage.ChallengeID,
		Data: types.MessageData{
			ChallengerID: incomingResponseMessage.Data.ChallengerID,
			RecipientID:  incomingResponseMessage.Data.RecipientID,
			Observers:    append([]string(nil), incomingResponseMessage.Data.Observers...),
			Challenge: types.ChallengeData{
				Block:      incomingResponseMessage.Data.Challenge.Block,
				Merkelroot: incomingResponseMessage.Data.Challenge.Merkelroot,
				FileHash:   incomingResponseMessage.Data.Challenge.FileHash,
				StartIndex: incomingResponseMessage.Data.Challenge.StartIndex,
				EndIndex:   incomingResponseMessage.Data.Challenge.EndIndex,
				Timestamp:  incomingResponseMessage.Data.Challenge.Timestamp,
			},
			Response: types.ResponseData{
				Block:      incomingResponseMessage.Data.Response.Block,
				Merkelroot: incomingResponseMessage.Data.Response.Merkelroot,
				Hash:       incomingResponseMessage.Data.Response.Hash,
				Timestamp:  incomingResponseMessage.Data.Response.Timestamp,
			},
			ChallengerEvaluation: types.EvaluationData{
				Block:      blockNumChallengeVerified,
				Merkelroot: blkVerbose1.MerkleRoot,
				Hash:       challengeCorrectHash,
				IsVerified: isVerified,
				Timestamp:  time.Now(),
			},
		},
		Sender: task.nodeID,
	}

	//Store the challenge evaluation message on the challenger side
	if err := task.StoreChallengeMessage(ctx, evaluationMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing challenger evaluation challenge message")
	}

	// send to observers for affirmations
	log.WithContext(ctx).WithField("challenge_id", evaluationMessage.ChallengeID).Info("sending evaluation message for affirmation to observers")
	affirmations, err := task.getAffirmationFromObservers(ctx, evaluationMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).WithField("challengeID", evaluationMessage.ChallengeID).Error("error sending evaluation message")
		return nil, err
	}

	if len(affirmations) < SuccessfulEvaluationThreshold {
		err := errors.Errorf("not enough affirmations have been received, failing the challenge")
		log.WithContext(ctx).WithError(err).Error("affirmations failed")

		return nil, err
	}

	//Broadcasting

	//blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	//log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Supernode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.ChallengeFile.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallenge) + " blocks!")

	task.SaveChallengeMessageState(
		ctx,
		saveStatus,
		incomingResponseMessage.ChallengeID,
		incomingResponseMessage.Data.ChallengerID,
		incomingResponseMessage.Data.ChallengerEvaluation.Block,
	)

	return nil, nil
}

func (task *SCTask) validateVerifyingStorageChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ResponseMessageType {
		return fmt.Errorf("incorrect message type to verify storage challenge")
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

// Send our evaluation message to other observers that might host this file.
func (task *SCTask) getAffirmationFromObservers(ctx context.Context, challengeMessage types.Message) (map[string]bool, error) {
	affirmations := make(map[string]bool)

	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send process storage challenge")
		return affirmations, err
	}

	evaluationMsg, err := task.prepareEvaluationMessage(ctx, challengeMessage)
	if err != nil {
		return affirmations, err
	}

	task.processEvaluationResults(ctx, nodesToConnect, evaluationMsg, affirmations)

	return affirmations, nil
}

// prepareEvaluationMessage prepares the evaluation message by gathering the data required for it
func (task *SCTask) prepareEvaluationMessage(ctx context.Context, challengeMessage types.Message) (pb.StorageChallengeMessage, error) {
	data, err := json.Marshal(challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error marshaling the data")
		return pb.StorageChallengeMessage{}, err
	}

	signature, err := task.SignMessage(ctx, data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the challenge message")
		return pb.StorageChallengeMessage{}, err
	}

	challengeMessage.SenderSignature = signature
	return pb.StorageChallengeMessage{
		MessageType:     pb.StorageChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}, nil
}

// processEvaluationResults simply sends an evaluation msg, receive an evaluation result and process the results
func (task *SCTask) processEvaluationResults(ctx context.Context, nodesToConnect []pastel.MasterNode, evaluationMsg pb.StorageChallengeMessage, affirmations map[string]bool) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, node := range nodesToConnect {
		node := node

		wg.Add(1)
		go func(node pastel.MasterNode) {
			defer wg.Done()

			affirmationResponse, err := task.sendEvaluationMessage(ctx, evaluationMsg, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("error sending evaluation message for processing")
				return
			}

			isSuccess, err := task.isAffirmationSuccessful(ctx, affirmationResponse)
			if err != nil {
				log.WithContext(ctx).WithField("node_id", node.ExtKey).WithError(err).
					Error("unsuccessful affirmation by observer")
				return
			}

			if isSuccess {
				mu.Lock()
				affirmations[node.ExtKey] = isSuccess
				mu.Unlock()
			}
		}(node)
	}

	wg.Wait()
}

// sendEvaluationMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task SCTask) sendEvaluationMessage(ctx context.Context, challengeMessage pb.StorageChallengeMessage, processingSupernodeAddr string) (types.Message, error) {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Info("Sending evaluation message to supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
		log.WithContext(ctx).
			WithField("challengeID", challengeMessage.ChallengeId).
			WithField("method", "sendEvaluationMessage").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return types.Message{}, err
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	affirmationResult, err := storageChallengeIF.VerifyEvaluationResult(ctx, &challengeMessage)
	if err != nil {
		log.WithContext(ctx).
			WithField("challengeID", challengeMessage.ChallengeId).
			WithField("method", "sendEvaluationMessage").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return types.Message{}, err
	}

	return affirmationResult, nil
}

// isAffirmationSuccessful checks the affirmation result
func (task *SCTask) isAffirmationSuccessful(ctx context.Context, msg types.Message) (bool, error) {
	isVerified, err := task.VerifyMessageSignature(ctx, msg)
	if err != nil {
		return false, errors.Errorf("unable to verify observer signature")
	}

	if !isVerified {
		return false, errors.Errorf("observer signature is not correct")
	}

	observerEvaluation := msg.Data.ObserverEvaluation

	if !observerEvaluation.IsEvaluationResultOK {
		return false, errors.Errorf("evaluation result is not correct")
	}

	if !observerEvaluation.IsChallengerSignatureOK {
		return false, errors.Errorf("challenger signature is not correct")
	}

	if !observerEvaluation.IsRecipientSignatureOK {
		return false, errors.Errorf("recipient signature is not correct")
	}

	if !observerEvaluation.IsChallengeTimestampOK {
		return false, errors.Errorf("challenge timestamp is not correct")
	}

	if !observerEvaluation.IsProcessTimestampOK {
		return false, errors.Errorf("response timestamp is not correct")
	}

	if !observerEvaluation.IsEvaluationTimestampOK {
		return false, errors.Errorf("evaluation timestamp is not correct")
	}

	return true, nil
}
