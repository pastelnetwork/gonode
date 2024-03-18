package healthcheckchallenge

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

const (
	verificationThreshold = 2 * time.Minute
)

// VerifyHealthCheckChallenge : Verifying the healthcheck challenge will occur only on the challenger node
func (task *HCTask) VerifyHealthCheckChallenge(ctx context.Context, incomingResponseMessage types.HealthCheckMessage) (*pb.StorageChallengeMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge").
		WithField("hc_challenge_id", incomingResponseMessage.ChallengeID)
	logger.Info("verify health-check challenge has been invoked") // Incoming challenge message validation

	if err := task.validateVerifyingStorageChallengeIncomingData(ctx, incomingResponseMessage); err != nil {
		return nil, err
	}

	//if the message is received by one of the observer then save the challenge message
	if task.isObserver(incomingResponseMessage.Data.Observers) {
		if err := task.StoreChallengeMessage(ctx, incomingResponseMessage); err != nil {
			logger.
				WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing response health-check message by the observer")
		} else {
			logger.WithField("node_id", task.nodeID).WithField("challenge_id",
				incomingResponseMessage.ChallengeID).
				Info("response health-check message has been stored by the observer")
		}

		return nil, nil
	}

	//if not the challenger, should return, otherwise proceed
	if task.nodeID != incomingResponseMessage.Data.ChallengerID {
		logger.WithField("node_id", task.nodeID).
			Debug("current node is not the health-check challenger to verify the response message")

		return nil, nil
	}

	//Challenger should also store the response message
	if err := task.StoreChallengeMessage(ctx, incomingResponseMessage); err != nil {
		logger.
			WithField("node_id", task.nodeID).
			WithError(err).
			Error("error storing response message")
	} else {
		logger.Debug("response health-check message has been stored by the challenger")
	}

	// Verifies the challenger's time against the recipient's time.
	isVerified := task.isVerifiedChallengerAndRecipientTimestamp(ctx, incomingResponseMessage.Data.Challenge.Timestamp, incomingResponseMessage.Data.Response.Timestamp)

	//Identify current block count to see if validation was performed within the mandated length of time (default one block)
	blockNumChallengeVerified, err := task.SuperNodeService.PastelClient.GetBlockCount(ctx)
	if err != nil {
		logger.WithError(err).Error("could not get current block count")
		return nil, err
	}

	blkVerbose1, err := task.SuperNodeService.PastelClient.GetBlockVerbose1(ctx, blockNumChallengeVerified)
	if err != nil {
		logger.WithError(err).Error("could not get current block verbose 1")
		return nil, err
	}

	blockNumChallengeSent := incomingResponseMessage.Data.Challenge.Block
	blocksVerifyStorageChallengeInBlocks := blockNumChallengeVerified - blockNumChallengeSent
	logger.WithField("no_of_blocks", blocksVerifyStorageChallengeInBlocks).Debug("No of blocks from the time HC challenge being sent and verified")

	evaluationMessage := types.HealthCheckMessage{
		MessageType: types.HealthCheckEvaluationMessageType,
		ChallengeID: incomingResponseMessage.ChallengeID,
		Data: types.HealthCheckMessageData{
			ChallengerID: incomingResponseMessage.Data.ChallengerID,
			RecipientID:  incomingResponseMessage.Data.RecipientID,
			Observers:    append([]string(nil), incomingResponseMessage.Data.Observers...),
			Challenge: types.HealthCheckChallengeData{
				Block:      incomingResponseMessage.Data.Challenge.Block,
				Merkelroot: incomingResponseMessage.Data.Challenge.Merkelroot,
				Timestamp:  incomingResponseMessage.Data.Challenge.Timestamp,
			},
			Response: types.HealthCheckResponseData{
				Block:      incomingResponseMessage.Data.Response.Block,
				Merkelroot: incomingResponseMessage.Data.Response.Merkelroot,
				Timestamp:  incomingResponseMessage.Data.Response.Timestamp,
			},
			ChallengerEvaluation: types.HealthCheckEvaluationData{
				Block:      blockNumChallengeVerified,
				Merkelroot: blkVerbose1.MerkleRoot,
				IsVerified: isVerified,
				Timestamp:  time.Now().UTC(),
			},
		},
		Sender:          task.nodeID,
		SenderSignature: nil,
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, evaluationMessage); err != nil {
		logger.WithField("message_type", evaluationMessage.MessageType).Error(
			"error storing health check challenge metric")
	}

	// send to observers for affirmations
	affirmations, err := task.getAffirmationFromObservers(ctx, evaluationMessage)
	if err != nil {
		logger.WithError(err).WithField("challengeID", evaluationMessage.ChallengeID).Error("error sending evaluation message")
		return nil, err
	}

	if len(affirmations) >= SuccessfulEvaluationThreshold {
		err := errors.Errorf("not enough affirmations have been received, failing the health-check challenge")
		logger.WithError(err).Error("affirmations failed")

		return nil, err
	}
	logger.Info("sufficient affirmations have been received")

	//Broadcasting
	broadcastingMsg, err := task.prepareBroadcastingMessage(ctx, evaluationMessage, affirmations)
	if err != nil {
		err := errors.Errorf("error preparing broadcasting message from the received affirmations")
		logger.WithError(err).Error("creating broadcasting message failed")

		return nil, err
	}

	if broadcastingMsg == nil {
		logger.Info("Unable to create broadcast health-check message")
		return nil, nil
	}

	if err := task.sendBroadcastingMessage(ctx, broadcastingMsg); err != nil {
		logger.WithError(err).Error("broadcasting storage challenge result failed")
		return nil, err
	}

	logger.Info("health-check challenge has been broadcast to the entire network")

	//blocksToRespondToStorageChallenge := outgoingChallengeMessage.BlockNumChallengeRespondedTo - incomingChallengeMessage.BlockNumChallengeSent
	//logger.WithField("method", "VerifyStorageChallenge").WithField("challengeID", incomingChallengeMessage.ChallengeId).Debug("Supernode " + outgoingChallengeMessage.RespondingMasternodeId + " responded to storage challenge for file hash " + outgoingChallengeMessage.ChallengeFile.FileHashToChallenge + " in " + fmt.Sprint(blocksToRespondToStorageChallenge) + " blocks!")

	return nil, nil
}

// isVerifiedChallengerAndRecipientTimeStamp Verifies the challenger's time against the recipient's time.
func (task *HCTask) isVerifiedChallengerAndRecipientTimestamp(ctx context.Context, challengeTime time.Time, responseTime time.Time) bool {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	if challengeTime.Before(responseTime) {
		if responseTime.Sub(challengeTime) < verificationThreshold {
			logger.Info("health check challenge has been verified")
			return true
		}
	}

	logger.Info("health check challenge has not been verified")
	return false
}

func (task *HCTask) validateVerifyingStorageChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	if incomingChallengeMessage.MessageType != types.HealthCheckResponseMessageType {
		return fmt.Errorf("incorrect message type to verify storage challenge")
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, incomingChallengeMessage); err != nil {
		logger.WithField("hc_challenge_id", incomingChallengeMessage.ChallengeID).
			WithField("message_type", incomingChallengeMessage.MessageType).Error(
			"error storing health check challenge metric")
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
func (task *HCTask) getAffirmationFromObservers(ctx context.Context, challengeMessage types.HealthCheckMessage) (map[string]types.HealthCheckMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		logger.WithError(err).Error("unable to find nodes to connect for getting affirmations")
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
func (task *HCTask) prepareEvaluationMessage(ctx context.Context, challengeMessage types.HealthCheckMessage) (pb.HealthCheckChallengeMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		logger.WithError(err).Error("error signing the challenge message")
		return pb.HealthCheckChallengeMessage{}, err
	}
	challengeMessage.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, challengeMessage); err != nil {
		logger.WithError(err).Error("error storing evaluation report message")
		return pb.HealthCheckChallengeMessage{}, err
	}

	return pb.HealthCheckChallengeMessage{
		MessageType:     pb.HealthCheckChallengeMessageMessageType(challengeMessage.MessageType),
		ChallengeId:     challengeMessage.ChallengeID,
		Data:            data,
		SenderId:        challengeMessage.Sender,
		SenderSignature: challengeMessage.SenderSignature,
	}, nil
}

// processEvaluationResults simply sends an evaluation msg, receive an evaluation result and process the results
func (task *HCTask) processEvaluationResults(ctx context.Context, nodesToConnect []pastel.MasterNode, evaluationMsg *pb.HealthCheckChallengeMessage) (affirmations map[string]types.HealthCheckMessage) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	affirmations = make(map[string]types.HealthCheckMessage)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, node := range nodesToConnect {
		node := node
		logger.WithField("node_id", node.ExtAddress).Info("processing sn address")
		wg.Add(1)

		go func() {
			defer wg.Done()

			affirmationResponse, err := task.sendEvaluationMessage(ctx, evaluationMsg, node.ExtAddress)
			if err != nil {
				logger.WithError(err).Error("error sending evaluation message for processing")
				return
			}

			if affirmationResponse != nil && affirmationResponse.ChallengeID != "" &&
				affirmationResponse.MessageType == types.HealthCheckAffirmationMessageType {

				if err := task.StoreHealthCheckChallengeMetric(ctx, *affirmationResponse); err != nil {
					logger.WithField("hc_challenge_id", affirmationResponse.ChallengeID).
						WithField("message_type", affirmationResponse.MessageType).Error(
						"error storing health check challenge metric")
				}

				isSuccess, err := task.isAffirmationSuccessful(ctx, affirmationResponse)
				if err != nil {
					logger.WithField("node_id", node.ExtKey).WithError(err).
						Error("unsuccessful affirmation by observer")
				}
				logger.WithField("is_success", isSuccess).Info("affirmation message has been verified")

				mu.Lock()
				affirmations[affirmationResponse.Sender] = *affirmationResponse
				mu.Unlock()

				if err := task.StoreChallengeMessage(ctx, *affirmationResponse); err != nil {
					logger.WithField("node_id", node.ExtKey).WithError(err).
						Error("error storing affirmation response from observer")
				}
			}
		}()
	}

	wg.Wait()

	return affirmations
}

// sendEvaluationMessage establish a connection with the processingSupernodeAddr and sends the given message to it.
func (task *HCTask) sendEvaluationMessage(ctx context.Context, challengeMessage *pb.HealthCheckChallengeMessage, processingSupernodeAddr string) (*types.HealthCheckMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge").WithField(
		"hc_challenge_id", challengeMessage.ChallengeId)

	//Connect over grpc
	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use node client to connect to: " + processingSupernodeAddr)
		logger.
			WithField("method", "sendEvaluationMessage").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}
	defer nodeClientConn.Close()

	storageChallengeIF := nodeClientConn.HealthCheckChallenge()

	affirmationResult, err := storageChallengeIF.VerifyHealthCheckEvaluationResult(ctx, challengeMessage)
	if err != nil {
		logger.
			WithField("method", "sendEvaluationMessage").
			WithField("node_address", processingSupernodeAddr).
			Warn(err.Error())
		return nil, err
	}

	logger.WithField("affirmation:", affirmationResult).Debug("affirmation result received")
	return &affirmationResult, nil
}

// isAffirmationSuccessful checks the affirmation result
func (task *HCTask) isAffirmationSuccessful(ctx context.Context, msg *types.HealthCheckMessage) (bool, error) {
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

// prepareBroadcastingMessage
func (task *HCTask) prepareBroadcastingMessage(ctx context.Context, evaluationMsg types.HealthCheckMessage, affirmations map[string]types.HealthCheckMessage) (*pb.BroadcastHealthCheckChallengeRequest, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	challenger := make(map[string][]byte)
	recipient := make(map[string][]byte)
	obs := make(map[string][]byte, len(affirmations))

	//challenge msg
	challengeMsg, err := task.RetrieveChallengeMessage(ctx, evaluationMsg.ChallengeID, int(types.HealthCheckChallengeMessageType))
	if err != nil {
		return nil, errors.Errorf("error retrieving challenge message")
	}

	challengeMsgBytes, err := json.Marshal(challengeMsg)
	if err != nil {
		return nil, errors.Errorf("error converting challenge message to bytes for broadcast")
	}
	challenger[challengeMsg.Sender] = challengeMsgBytes

	//response msg
	responseMsg, err := task.RetrieveChallengeMessage(ctx, evaluationMsg.ChallengeID, int(types.HealthCheckResponseMessageType))
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
			logger.WithField("pastel_id", pastelID).WithError(err).Error("error converting affirmation msg to bytes")
			continue
		}

		obs[pastelID] = affirmationMsg
	}

	return &pb.BroadcastHealthCheckChallengeRequest{
		ChallengeId: challengeMsg.ChallengeID,
		Challenger:  challenger, //Challenger's message
		Recipient:   recipient,  // Recipient's message
		Observers:   obs,        // Affirmations from observers
	}, nil
}

func (task *HCTask) sendBroadcastingMessage(ctx context.Context, msg *pb.BroadcastHealthCheckChallengeRequest) error {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, types.HealthCheckMessage{MessageType: types.HealthCheckBroadcastMessageType})
	if err != nil {
		logger.WithError(err).Error("unable to find nodes to connect for send broadcast storage challenge")
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

			logger := logger.WithField("node_address", node.ExtAddress)

			if msg == nil || node.ExtAddress == "" {
				return //not the valid state
			}

			if err := task.send(ctx, msg, node.ExtAddress); err != nil {
				logger.WithError(err).Error("error sending broadcast message for processing")
				return
			}
		}()
	}
	wg.Wait()

	return nil
}

func (task *HCTask) send(ctx context.Context, req *pb.BroadcastHealthCheckChallengeRequest, processingSupernodeAddr string) error {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	nodeClientConn, err := task.nodeClient.Connect(ctx, processingSupernodeAddr)
	if err != nil {
		err = fmt.Errorf("Could not use nodeclient to connect to: " + processingSupernodeAddr)
		logger.WithField("challengeID", req.ChallengeId).WithField("method", "sendBroadcastingMessage").Warn(err.Error())
		return err
	}
	defer nodeClientConn.Close()

	healthcheckChallengeIF := nodeClientConn.HealthCheckChallenge()

	return healthcheckChallengeIF.BroadcastHealthCheckChallengeResult(ctx, req)
}
