package healthcheckchallenge

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// ProcessHealthCheckChallenge processes the health-check challenge by responding in time
func (task *HCTask) ProcessHealthCheckChallenge(ctx context.Context, incomingChallengeMessage types.HealthCheckMessage) (*pb.HealthCheckChallengeMessage, error) {
	logger := log.WithContext(ctx).WithField("Method", "ProcessHealthCheckChallenge").
		WithField("hc_challenge_id", incomingChallengeMessage.ChallengeID)

	logger.Info("process health-check challenge has been invoked")

	// incoming challenge message validation
	if err := task.validateProcessingHealthCheckChallengeIncomingData(ctx, incomingChallengeMessage); err != nil {
		logger.WithError(err).Error("Error validating healthcheck challenge incoming data: ")
		return nil, err
	}
	logger.Debug("incoming health-check challenge validated")

	//if the message is received by one of the observer then save the challenge message, lock the file & return
	if task.isObserver(incomingChallengeMessage.Data.Observers) {
		logger.WithField("node_id", task.nodeID).WithField("hc_challenge_id",
			incomingChallengeMessage.ChallengeID)

		if err := task.StoreChallengeMessage(ctx, incomingChallengeMessage); err != nil {
			logger.WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing health-check challenge message")
		}
		logger.Debug("health-check challenge message has been stored by the observer")

		return nil, nil
	}

	//if not the challenge recipient, should return, otherwise proceed
	if task.nodeID != incomingChallengeMessage.Data.RecipientID {
		logger.WithField("node_id", task.nodeID).
			Debug("current node is not the health-check challenge recipient to process challenge message")

		return nil, nil
	}

	if err := task.StoreChallengeMessage(ctx, incomingChallengeMessage); err != nil {
		logger.WithField("node_id", task.nodeID).
			WithError(err).
			Error("error storing health-check challenge message")
	}
	logger.Debug("health-check challenge message has been stored by the recipient")

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

	outgoingResponseMessage := types.HealthCheckMessage{
		MessageType: types.HealthCheckResponseMessageType,
		ChallengeID: incomingChallengeMessage.ChallengeID,
		Data: types.HealthCheckMessageData{
			ChallengerID: incomingChallengeMessage.Data.ChallengerID,
			RecipientID:  incomingChallengeMessage.Data.RecipientID,
			Observers:    incomingChallengeMessage.Data.Observers,
			Challenge: types.HealthCheckChallengeData{
				Block:      incomingChallengeMessage.Data.Challenge.Block,
				Merkelroot: incomingChallengeMessage.Data.Challenge.Merkelroot,
				Timestamp:  incomingChallengeMessage.Data.Challenge.Timestamp,
			},
			Response: types.HealthCheckResponseData{
				Block:      blockNumChallengeRespondedTo,
				Merkelroot: blkVerbose1.MerkleRoot,
				Timestamp:  time.Now().UTC(),
			},
		},
		Sender:          task.nodeID,
		SenderSignature: nil,
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, outgoingResponseMessage); err != nil {
		logger.WithField("message_type", outgoingResponseMessage.MessageType).Error(
			"error storing health check challenge metric")
	}

	// send to Supernodes to validate challenge response hash
	if err = task.sendVerifyHealthCheckChallenge(ctx, outgoingResponseMessage); err != nil {
		logger.WithError(err).WithField("challengeID", incomingChallengeMessage.ChallengeID).Error("could not send processed challenge message to node for verification")
		return nil, err
	}
	logger.Info("health-check message sent to other SNs for verification")

	return nil, nil
}

func (task *HCTask) validateProcessingHealthCheckChallengeIncomingData(ctx context.Context, incomingChallengeMessage types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("Method", "ProcessHealthCheckChallenge")

	if incomingChallengeMessage.MessageType != types.HealthCheckChallengeMessageType {
		return fmt.Errorf("incorrect message type to processing health check challenge")
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

// Send our verification message to (default 10) other supernodes that might host this file.
func (task *HCTask) sendVerifyHealthCheckChallenge(ctx context.Context, challengeMessage types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("Method", "ProcessHealthCheckChallenge")

	nodesToConnect, err := task.GetNodesAddressesToConnect(ctx, challengeMessage)
	if err != nil {
		logger.WithError(err).Error("unable to find nodes to connect for send verify healthcheck challenge")
		return err
	}

	if nodesToConnect == nil {
		return errors.Errorf("no nodes found to connect to send verify healthcheck challenge")
	}

	signature, data, err := task.SignMessage(ctx, challengeMessage.Data)
	if err != nil {
		logger.WithError(err).Error("error signing the challenge message")
		return err
	}
	challengeMessage.SenderSignature = signature

	msg := pb.HealthCheckChallengeMessage{
		MessageType:     pb.HealthCheckChallengeMessageMessageType(challengeMessage.MessageType),
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

			logger := logger.WithField("node_address", node.ExtAddress)

			if err := task.SendMessage(ctx, &msg, node.ExtAddress); err != nil {
				log.WithError(err).WithField("node_address", challengerNode.ExtAddress).
					Debug("error sending process healthcheck challenge message for processing")
				return
			}

			logger.Debug("response message has been sent")
		}()
	}
	wg.Wait()
	logger.WithField("challenge_id", challengeMessage.ChallengeID).Debug("response message has been sent to observers")

	if err := task.SendMessage(ctx, &msg, challengerNode.ExtAddress); err != nil {
		log.WithField("node_address", challengerNode.ExtAddress).
			WithError(err).Debug("error sending response message to challenger for verification")
		return err
	}
	logger.WithField("hc_challenge_id", challengeMessage.ChallengeID).
		WithField("node_address", challengerNode.ExtAddress).Debug("response message has been sent to challenger")

	return nil
}

func (task *HCTask) isObserver(sliceOfObservers []string) bool {
	for _, ob := range sliceOfObservers {
		if ob == task.nodeID {
			return true
		}
	}

	return false
}
