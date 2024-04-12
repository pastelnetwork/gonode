package healthcheckchallenge

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
)

const (
	timestampTolerance          = 30 * time.Second
	challengeTimestampTolerance = 120 * time.Second
)

// VerifyHealthCheckEvaluationResult process the evaluation report from the challenger by the observers
func (task *HCTask) VerifyHealthCheckEvaluationResult(ctx context.Context, incomingEvaluationResult types.HealthCheckMessage) (types.HealthCheckMessage, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckEvaluationResult").
		WithField("hc_challenge_id", incomingEvaluationResult.ChallengeID)

	logger.Info("verify health-check challenge evaluation result has been invoked")

	if err := task.validateVerifyingHealthCheckChallengeEvaluationReport(ctx, incomingEvaluationResult); err != nil {
		return types.HealthCheckMessage{}, err
	}
	logger.Debug("Incoming evaluation signature validated")

	receivedAt := time.Now().UTC()

	//if observers then save the evaluation message
	if task.isObserver(incomingEvaluationResult.Data.Observers) {
		logger.WithField("node_id", task.nodeID)

		if err := task.StoreChallengeMessage(ctx, incomingEvaluationResult); err != nil {
			logger.
				WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing challenge message")
		}

		logger.Debug("health-check challenge evaluation report by challenger has been stored by the observer")
	} else {
		logger.WithField("node_id", task.nodeID).Debug("not the observer to process evaluation report")
		return types.HealthCheckMessage{}, nil
	}

	//retrieve messages from the DB against challengeID
	challengeMessage, responseMessage, err := task.retrieveChallengeAndResponseMessageForEvaluationReportVerification(ctx, incomingEvaluationResult.ChallengeID)
	if err != nil {
		logger.WithError(err).Error("error retrieving the health-check challenge and response message for verification")
		return types.HealthCheckMessage{}, err
	}

	if challengeMessage == nil || responseMessage == nil {
		logger.Error("unable to retrieve health-check challenge or response message")
		return types.HealthCheckMessage{}, errors.Errorf("unable to retrieve challenge or response message")
	}

	//Verify message signatures
	isChallengerSignatureOk, isRecipientSignatureOk, err := task.verifyMessageSignaturesForEvaluation(ctx, *challengeMessage, *responseMessage)
	if err != nil {
		logger.WithError(err).
			Error("error verifying signatures for evaluation")
		return types.HealthCheckMessage{}, errors.Errorf("error verifying signatures for evalaution")
	}

	//Verify message timestamps
	isChallengeTSOk, isResponseTSOk, isEvaluationTSOk := task.verifyTimestampsForEvaluation(ctx, *challengeMessage,
		*responseMessage, incomingEvaluationResult, receivedAt)

	//Verify Merkelroot and BlockNumber
	isChallengeBlockAndMROk, isResponseBlockAndMROk := task.verifyMerkelrootAndBlockNum(ctx, *challengeMessage, *responseMessage, incomingEvaluationResult)

	evaluationResultResponse := types.HealthCheckMessage{
		MessageType: types.HealthCheckAffirmationMessageType,
		ChallengeID: incomingEvaluationResult.ChallengeID,
		Data: types.HealthCheckMessageData{
			ChallengerID: incomingEvaluationResult.Data.ChallengerID,
			RecipientID:  incomingEvaluationResult.Data.RecipientID,
			Observers:    append([]string(nil), incomingEvaluationResult.Data.Observers...),
			Challenge: types.HealthCheckChallengeData{
				Block:      incomingEvaluationResult.Data.Challenge.Block,
				Merkelroot: incomingEvaluationResult.Data.Challenge.Merkelroot,
				Timestamp:  incomingEvaluationResult.Data.Challenge.Timestamp,
			},
			Response: types.HealthCheckResponseData{
				Block:      incomingEvaluationResult.Data.Response.Block,
				Merkelroot: incomingEvaluationResult.Data.Response.Merkelroot,
				Timestamp:  incomingEvaluationResult.Data.Response.Timestamp,
			},
			ChallengerEvaluation: types.HealthCheckEvaluationData{
				Block:      incomingEvaluationResult.Data.ChallengerEvaluation.Block,
				Merkelroot: incomingEvaluationResult.Data.ChallengerEvaluation.Merkelroot,
				IsVerified: incomingEvaluationResult.Data.ChallengerEvaluation.IsVerified,
				Timestamp:  incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp,
			},
			ObserverEvaluation: types.HealthCheckObserverEvaluationData{
				IsChallengerSignatureOK: isChallengerSignatureOk && isChallengeBlockAndMROk,
				IsRecipientSignatureOK:  isRecipientSignatureOk && isResponseBlockAndMROk,
				IsChallengeTimestampOK:  isChallengeTSOk,
				IsProcessTimestampOK:    isResponseTSOk,
				IsEvaluationTimestampOK: isEvaluationTSOk,
				IsEvaluationResultOK:    isEvaluationTSOk,
				Timestamp:               time.Now().UTC(),
			},
		},
		Sender: task.nodeID,
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, evaluationResultResponse); err != nil {
		logger.
			WithField("message_type", evaluationResultResponse.MessageType).Error(
			"error storing health check challenge metric")
	}

	signature, _, err := task.SignMessage(ctx, evaluationResultResponse.Data)
	if err != nil {
		logger.WithError(err).Error("error signing evaluation response")
	}
	evaluationResultResponse.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, evaluationResultResponse); err != nil {
		logger.WithError(err).Error("error storing evaluation response ")
	}
	logger.WithField("node_id", task.nodeID).Debug("evaluation report has been evaluated by observer")

	return evaluationResultResponse, nil
}

// RetrieveChallengeMessage retrieves the challenge messages from db for further verification
func (task *HCTask) RetrieveChallengeMessage(ctx context.Context, challengeID string, msgType int) (challengeMessage *types.HealthCheckMessage, err error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	store, err := queries.OpenHistoryDB()
	if err != nil {
		logger.WithError(err).Error("Error Opening DB")
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		msg, err := store.QueryHCChallengeMessage(challengeID, msgType)
		if err != nil {
			logger.WithError(err).Error("Error retrieving challenge message from DB")
			return nil, err
		}
		logger.Debug("health check challenge msg retrieved")

		var data types.HealthCheckMessageData
		err = json.Unmarshal(msg.Data, &data)
		if err != nil {
			logger.WithError(err).Error("Error un-marshaling message data")
			return nil, err
		}

		return &types.HealthCheckMessage{
			ChallengeID:     msg.ChallengeID,
			MessageType:     types.HealthCheckMessageType(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
			CreatedAt:       msg.CreatedAt,
			UpdatedAt:       msg.UpdatedAt,
		}, nil
	}

	return nil, nil
}

func (task *HCTask) verifyMessageSignaturesForEvaluation(ctx context.Context, challengeMessage, responseMessage types.HealthCheckMessage) (challengeMsgSig, ResMsgSig bool, err error) {
	isChallengerSignatureOk, err := task.VerifyMessageSignature(ctx, challengeMessage)
	if err != nil {
		return false, false, errors.Errorf("error verifying challenge message signature")
	}

	isRecipientSignatureOk, err := task.VerifyMessageSignature(ctx, responseMessage)
	if err != nil {
		return false, false, errors.Errorf("error verifying challenge message signature")
	}

	return isChallengerSignatureOk, isRecipientSignatureOk, nil
}

func (task *HCTask) verifyTimestampsForEvaluation(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.HealthCheckMessage, evaluationMsgRecvTime time.Time) (challengeMsgTS, resMsgTS, evalMsgTS bool) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	isChallengeTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Challenge.Timestamp.UTC(),
		challengeMessage.Data.Challenge.Timestamp.UTC(), challengeMessage.CreatedAt.UTC(),
		challengeMessage.ChallengeID, challengeMessage.MessageType)
	logger.Debug("health-check challenge message timestamps have been verified successfully")

	isResponseTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Response.Timestamp.UTC(),
		responseMessage.Data.Response.Timestamp.UTC(), responseMessage.CreatedAt.UTC(),
		responseMessage.ChallengeID, responseMessage.MessageType)
	logger.Debug("health-check response message timestamps have been verified successfully")

	differenceBetweenEvaluationMsgSendAndRecvTime := incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp.Sub(evaluationMsgRecvTime)
	if differenceBetweenEvaluationMsgSendAndRecvTime < 0 {
		differenceBetweenEvaluationMsgSendAndRecvTime = -differenceBetweenEvaluationMsgSendAndRecvTime
	}

	var evaluationReportTSOk bool
	if differenceBetweenEvaluationMsgSendAndRecvTime <= timestampTolerance {
		logger.Debug("health-check challenge evaluation message timestamp is verified successfully")
		evaluationReportTSOk = true
	} else {
		logger.Debug("the time diff between the health-check challenge evaluation report send and receive time, exceeds the tolerance limit")
		evaluationReportTSOk = false
	}

	overallTimeTakenByChallenge := evaluationMsgRecvTime.Sub(challengeMessage.Data.Challenge.Timestamp)
	if overallTimeTakenByChallenge < 0 {
		overallTimeTakenByChallenge = -overallTimeTakenByChallenge
	}

	overallChallengeTSOk := overallTimeTakenByChallenge <= challengeTimestampTolerance

	var isEvaluationTSOk bool
	if evaluationReportTSOk && overallChallengeTSOk {
		isEvaluationTSOk = true
	}
	logger.Debug("health-check challenge evaluation report & overall time taken by challenge has been verified successfully")

	return isChallengeTSOk, isResponseTSOk, isEvaluationTSOk
}

func (task *HCTask) verifyMessageTimestamps(ctx context.Context, timeWhenMsgSent, sentTimeWhenMsgReceived, timeWhenMsgStored time.Time, challengeID string, messageType types.HealthCheckMessageType) bool {
	logger := log.
		WithField("challenge_id", challengeID).
		WithField("message_type", messageType.String()).
		WithField("node_id", task.nodeID).
		WithField("sent_time", timeWhenMsgSent).
		WithField("sent_time_upon_receiving", sentTimeWhenMsgReceived).
		WithField("time_when_stored", timeWhenMsgStored)

	difference := timeWhenMsgStored.Sub(timeWhenMsgSent)
	if difference < 0 {
		difference = -difference
	}

	var isTSOk bool
	if timeWhenMsgSent.Equal(sentTimeWhenMsgReceived) {

		if difference <= timestampTolerance {
			isTSOk = true
		} else {
			logger.Debug("the time difference when message has been sent and when message has been stored, " +
				"exceeds the tolerance limit")

			isTSOk = false
		}

	} else {
		logger.Debug("time when message sent and sent time when message received are not same")
	}

	return isTSOk
}

func (task *HCTask) verifyMerkelrootAndBlockNum(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.HealthCheckMessage) (isChallengeOk, isResOk bool) {
	logger := log.
		WithField("challenge_id", incomingEvaluationResult.ChallengeID).
		WithField("node_id", task.nodeID)

	isChallengeBlockOk := challengeMessage.Data.Challenge.Block == incomingEvaluationResult.Data.Challenge.Block
	if !isChallengeBlockOk {
		logger.Debug("challenge msg block num is different and not same when challenge sent and received")
	}

	isChallengeMerkelrootOk := challengeMessage.Data.Challenge.Merkelroot == incomingEvaluationResult.Data.Challenge.Merkelroot
	if !isChallengeMerkelrootOk {
		logger.Debug("challenge msg merkelroot is different and not same when challenge sent and received")
	}

	isResponseBlockOk := responseMessage.Data.Response.Block == incomingEvaluationResult.Data.Response.Block
	if !isResponseBlockOk {
		logger.Debug("response msg block num is different and not same when challenge sent and received")
	}

	isResponseMerkelrootOk := responseMessage.Data.Response.Merkelroot == incomingEvaluationResult.Data.Response.Merkelroot
	if !isResponseMerkelrootOk {
		logger.Debug("response msg merkelroot is different and not same when challenge sent and received")
	}

	return isChallengeBlockOk && isChallengeMerkelrootOk, isResponseBlockOk && isResponseMerkelrootOk
}

func (task *HCTask) validateVerifyingHealthCheckChallengeEvaluationReport(ctx context.Context, evaluationReport types.HealthCheckMessage) error {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	if evaluationReport.MessageType != types.HealthCheckEvaluationMessageType {
		return fmt.Errorf("incorrect message type to verify evaluation report")
	}

	if err := task.StoreHealthCheckChallengeMetric(ctx, evaluationReport); err != nil {
		logger.WithField("challenge_id", evaluationReport.ChallengeID).
			WithField("message_type", evaluationReport.MessageType).Error(
			"error storing health check challenge metric")
	}

	isVerified, err := task.VerifyMessageSignature(ctx, evaluationReport)
	if err != nil {
		return errors.Errorf("error verifying sender's signature: %w", err)
	}

	if !isVerified {
		return errors.Errorf("not able to verify message signature")
	}

	return nil
}

func (task *HCTask) retrieveChallengeAndResponseMessageForEvaluationReportVerification(ctx context.Context, challengeID string) (cMsg, rMsg *types.HealthCheckMessage, er error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyHealthCheckChallenge")

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 15 * time.Second
	b.InitialInterval = 1 * time.Second

	var challengeMessage *types.HealthCheckMessage
	var responseMessage *types.HealthCheckMessage

	// Retry fetching challenge message
	err := backoff.Retry(func() error {
		var err error
		challengeMessage, err = task.RetrieveChallengeMessage(ctx, challengeID, int(types.HealthCheckChallengeMessageType))
		if err != nil {
			if err == sql.ErrNoRows {
				return err // Will retry
			}
			return backoff.Permanent(err) // Won't retry
		}
		return nil // Success
	}, backoff.WithMaxRetries(b, 5))

	if err != nil {
		logger.WithError(err).Error("failed to retrieve challenge message")
		return nil, nil, errors.Errorf("failed to retrieve challenge message for evaluation affirmation")
	}

	// Reset backoff before next operation
	b.Reset()

	// Retry fetching response message
	err = backoff.Retry(func() error {
		var err error
		responseMessage, err = task.RetrieveChallengeMessage(ctx, challengeID, int(types.HealthCheckResponseMessageType))
		if err != nil {
			if err == sql.ErrNoRows {
				return err // Will retry
			}
			return backoff.Permanent(err) // Won't retry
		}
		return nil // Success
	}, backoff.WithMaxRetries(b, 5))

	if err != nil {
		logger.WithError(err).Error("failed to retrieve challenge message")
		return nil, nil, errors.Errorf("failed to retrieve challenge message for evaluation affirmation")
	}

	return challengeMessage, responseMessage, nil
}
