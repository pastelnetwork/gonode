package storagechallenge

import (
	"context"
	"fmt"
	"sync"
	"time"

	json "github.com/json-iterator/go"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// VerifyStorageChallenge : Verifying the storage challenge will occur only on the challenger node
func (task SCTask) VerifyStorageChallenge(ctx context.Context, incomingResponseMessage types.Message) (*pb.StorageChallengeMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyStorageChallenge").
		WithField("sc_challenge_id", incomingResponseMessage.ChallengeID)

	if err := task.validateVerifyingStorageChallengeIncomingData(ctx, incomingResponseMessage); err != nil {
		return nil, err
	}

	if err := task.SCService.P2PClient.EnableKey(ctx, incomingResponseMessage.Data.Challenge.FileHash); err != nil {
		logger.WithError(err).Error("error enabling the symbol file")
	}

	//if the message is received by one of the observer then save the challenge message
	if task.isObserver(incomingResponseMessage.Data.Observers) {
		if err := task.StoreChallengeMessage(ctx, incomingResponseMessage); err != nil {
			logger.WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing response message by the observer")
		} else {
			logger.WithField("node_id", task.nodeID).
				Debug("response message has been stored by the observer")
		}

		return nil, nil
	}

	//if not the challenger, should return, otherwise proceed
	if task.nodeID != incomingResponseMessage.Data.ChallengerID {
		logger.Debug("current node is not the challenger to verify the response message")

		return nil, nil
	}

	logger.Info("Storage challenge evaluation started") // Incoming challenge message validation
	//Challenger should also store the response message
	if err := task.StoreChallengeMessage(ctx, incomingResponseMessage); err != nil {
		logger.WithError(err).Error("error storing response message")
	} else {
		logger.Debug("response message has been stored by the challenger")
	}

	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingResponseMessage.Data.Challenge.FileHash, true)
	if err != nil {
		logger.WithError(err).Error("could not read queries file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())
		return nil, err
	}
	logger.WithContext(ctx).Debug("file has been retrieved for verification")

	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingResponseMessage.Data.Challenge.StartIndex, incomingResponseMessage.Data.Challenge.EndIndex)
	logger.Debug("hash of the data has been generated against the given indices")

	blockNumChallengeVerified, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).WithField("challengeID", incomingResponseMessage.ChallengeID).Error("could not get current block count")
		return nil, err
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeVerified)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}

	blockNumChallengeSent := incomingResponseMessage.Data.Challenge.Block
	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - blockNumChallengeSent

	// determine success or failure
	var isVerified bool
	if incomingResponseMessage.Data.Response.Hash == challengeCorrectHash {
		isVerified = true
		logger.Debug(fmt.Sprintf("Supernode %s correctly responded in %d blocks to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, blocksVerifyStorageChallengeInBlocks, incomingResponseMessage.Data.Challenge.FileHash))
	} else {
		log.Debug(fmt.Sprintf("Supernode %s failed by incorrectly responding to a storage challenge for file %s", incomingResponseMessage.Data.RecipientID, incomingResponseMessage.Data.Challenge.FileHash))
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
				Timestamp:  time.Now().UTC(),
			},
		},
		Sender:          task.nodeID,
		SenderSignature: nil,
	}

	if err := task.StoreStorageChallengeMetric(ctx, evaluationMessage); err != nil {
		logger.WithError(err).WithField("message_type", evaluationMessage.MessageType).Error(
			"error storing storage challenge metric")
	}

	// send to observers for affirmations
	log.WithContext(ctx).WithField("challenge_id", evaluationMessage.ChallengeID).Debug("sending evaluation message for affirmation to observers")
	affirmations, err := task.getAffirmationFromObservers(ctx, evaluationMessage)
	if err != nil {
		logger.WithError(err).Error("error sending evaluation message")
		return nil, err
	}

	if len(affirmations) < SuccessfulEvaluationThreshold {
		err := errors.Errorf("not enough affirmations have been received, failing the challenge")
		logger.WithError(err).Error("affirmations failed")

		return nil, err
	}
	logger.Debug("sufficient affirmations have been received")

	if err := task.SCService.P2PClient.EnableKey(ctx, evaluationMessage.Data.Challenge.FileHash); err != nil {
		logger.WithError(err).Error("error enabling the symbol file")
		return nil, err
	}
	logger.Debug("key has been enabled by challenger")

	//Broadcasting
	broadcastingMsg, err := task.prepareBroadcastingMessage(ctx, evaluationMessage, affirmations)
	if err != nil {
		err := errors.Errorf("error preparing broadcasting message from the received affirmations")
		logger.WithError(err).Error("creating broadcasting message failed")

		return nil, err
	}

	if broadcastingMsg == nil {
		logger.Debug("Unable to create broadcast message")
		return nil, nil
	}

	if err := task.sendBroadcastingMessage(ctx, broadcastingMsg); err != nil {
		logger.WithError(err).Error("broadcasting storage challenge result failed")
		return nil, err
	}

	logger.Info("storage challenge has been broadcast to the entire network")
	return nil, nil
}

func (task SCTask) validateVerifyingStorageChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.Message) error {
	if incomingChallengeMessage.MessageType != types.ResponseMessageType {
		return fmt.Errorf("incorrect message type to verify storage challenge")
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

// Send our evaluation message to other observers that might host this file.
func (task SCTask) getAffirmationFromObservers(ctx context.Context, challengeMessage types.Message) (map[string]types.Message, error) {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for getting affirmations")
		return nil, err
	}

	if nodesToConnect == nil {
		return nil, errors.Errorf("no nodes found to connect for getting affirmations")
	}

	evaluationMsg, err := task.prepareEvaluationMessage(ctx, challengeMessage)
	if err != nil {
		return nil, err
	}

	return task.processEvaluationResults(ctx, nodesToConnect, &evaluationMsg), nil
}

// prepareEvaluationMessage prepares the evaluation message by gathering the data required for it
func (task SCTask) prepareEvaluationMessage(ctx context.Context, challengeMessage types.Message) (pb.StorageChallengeMessage, error) {
	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing the challenge message")
		return pb.StorageChallengeMessage{}, err
	}
	challengeMessage.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, challengeMessage); err != nil {
		log.WithContext(ctx).WithError(err).Error("error storing evaluation report message")
		return pb.StorageChallengeMessage{}, err
	}

	return pb.StorageChallengeMessage{
		MessageType:     pb.StorageChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}, nil
}

// processEvaluationResults simply sends an evaluation msg, receive an evaluation result and process the results
func (task SCTask) processEvaluationResults(ctx context.Context, nodesToConnect []pastel.MasterNode, evaluationMsg *pb.StorageChallengeMessage) (affirmations map[string]types.Message) {
	affirmations = make(map[string]types.Message)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, node := range nodesToConnect {
		node := node
		log.WithContext(ctx).WithField("node_id", node.ExtAddress).Debug("processing sn address")
		wg.Add(1)

		go func() {
			defer wg.Done()

			affirmationResponse, err := task.sendEvaluationMessage(ctx, evaluationMsg, node.ExtAddress)
			if err != nil {
				log.WithError(err).Debug("error sending evaluation message for processing")
				return
			}

			if affirmationResponse != nil && affirmationResponse.ChallengeID != "" &&
				affirmationResponse.MessageType == types.AffirmationMessageType {

				if err := task.StoreStorageChallengeMetric(ctx, *affirmationResponse); err != nil {
					log.WithContext(ctx).WithField("challenge_id", affirmationResponse.ChallengeID).
						WithField("message_type", affirmationResponse.MessageType).Error(
						"error storing storage challenge metric")
				}

				isSuccess, err := task.isAffirmationSuccessful(ctx, affirmationResponse)
				if err != nil {
					log.WithContext(ctx).WithField("node_id", node.ExtKey).WithError(err).
						Error("unsuccessful affirmation by observer")
				}
				log.WithContext(ctx).WithField("is_success", isSuccess).Debug("affirmation message has been verified")

				mu.Lock()
				affirmations[affirmationResponse.Sender] = *affirmationResponse
				mu.Unlock()

				if err := task.StoreChallengeMessage(ctx, *affirmationResponse); err != nil {
					log.WithContext(ctx).WithField("node_id", node.ExtKey).WithError(err).
						Error("error storing affirmation response from observer")
				}
			}
		}()
	}

	wg.Wait()

	return affirmations
}

// sendEvaluationMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task SCTask) sendEvaluationMessage(ctx context.Context, challengeMessage *pb.StorageChallengeMessage, processingSupernodeAddr string) (*types.Message, error) {
	log.WithContext(ctx).WithField("challenge_id", challengeMessage.ChallengeId).Debug("Sending evaluation message to supernode address: " + processingSupernodeAddr)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.ConnectSN(ctx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "SendScEvaluationMessage", err)
		return nil, fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	affirmationResult, err := storageChallengeIF.VerifyEvaluationResult(ctx, challengeMessage)
	if err != nil {
		log.WithContext(ctx).
			WithField("challengeID", challengeMessage.ChallengeId).
			WithField("method", "sendEvaluationMessage").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}

	log.WithContext(ctx).WithField("affirmation:", affirmationResult).Debug("affirmation result received")
	return &affirmationResult, nil
}

// isAffirmationSuccessful checks the affirmation result
func (task SCTask) isAffirmationSuccessful(ctx context.Context, msg *types.Message) (bool, error) {
	if msg.ChallengeID == "" {
		return false, nil
	}

	isVerified, err := task.VerifyMessageSignature(ctx, *msg)
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

// prepareBroadcastingMessage(ctx
func (task SCTask) prepareBroadcastingMessage(ctx context.Context, evaluationMsg types.Message, affirmations map[string]types.Message) (*pb.BroadcastStorageChallengeRequest, error) {
	challenger := make(map[string][]byte)
	recipient := make(map[string][]byte)
	obs := make(map[string][]byte, len(affirmations))

	//challenge msg
	challengeMsg, err := task.RetrieveChallengeMessage(ctx, evaluationMsg.ChallengeID, int(types.ChallengeMessageType))
	if err != nil {
		return nil, errors.Errorf("error retrieving challenge message")
	}

	challengeMsgBytes, err := json.Marshal(challengeMsg)
	if err != nil {
		return nil, errors.Errorf("error converting challenge message to bytes for broadcast")
	}
	challenger[challengeMsg.Sender] = challengeMsgBytes

	//response msg
	responseMsg, err := task.RetrieveChallengeMessage(ctx, evaluationMsg.ChallengeID, int(types.ResponseMessageType))
	if err != nil {
		return nil, errors.Errorf("error retrieving the response message")
	}

	responseMsgBytes, err := json.Marshal(responseMsg)
	if err != nil {
		return nil, errors.Errorf("error converting challenge message to bytes for broadcast")
	}
	recipient[responseMsg.Sender] = responseMsgBytes

	//affirmation msgs
	for pastelID, msg := range affirmations {
		affirmationMsg, err := json.Marshal(msg)
		if err != nil {
			log.WithContext(ctx).WithField("pastel_id", pastelID).WithError(err).Error("error converting affirmation msg to bytes")
			continue
		}

		obs[pastelID] = affirmationMsg
	}

	return &pb.BroadcastStorageChallengeRequest{
		ChallengeId: challengeMsg.ChallengeID,
		Challenger:  challenger, //Challenger's message
		Recipient:   recipient,  // Recipient's message
		Observers:   obs,        // Affirmations from observers
	}, nil
}

func (task SCTask) sendBroadcastingMessage(ctx context.Context, msg *pb.BroadcastStorageChallengeRequest) error {
	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, types.Message{MessageType: types.BroadcastMessageType})
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to find nodes to connect for send broadcast storage challenge")
		return err
	}

	if len(nodesToConnect) == 0 {
		return errors.Errorf("no nodes found to connect for send broadcast storage challenge")
	}

	// Create a buffered channel to limit the number of Goroutines
	sem := make(chan struct{}, 10) // Limit set to 10
	var wg sync.WaitGroup

	for _, node := range nodesToConnect {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Acquire a token from the semaphore channel
			sem <- struct{}{}
			defer func() { <-sem }() // Release the token back into the channel

			if msg == nil || node.ExtAddress == "" {
				return //not the valid state
			}

			if err := task.send(ctx, msg, node.ExtAddress); err != nil {
				log.WithError(err).Debug("error sending broadcast message for processing")
				return
			}
		}()
	}
	wg.Wait()

	return nil
}

func (task SCTask) send(ctx context.Context, req *pb.BroadcastStorageChallengeRequest, processingSupernodeAddr string) error {
	//Connect over grpc
	nodeClientConn, err := task.nodeClient.ConnectSN(ctx, processingSupernodeAddr)
	if err != nil {
		logError(ctx, "BroadcastStorageChallengeResult", err)
		return fmt.Errorf("Could not use nodeclient to connect to: " + processingSupernodeAddr)
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.StorageChallenge()

	return storageChallengeIF.BroadcastStorageChallengeResult(ctx, req)
}
