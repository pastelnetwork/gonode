package storagechallenge

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	json "github.com/json-iterator/go"

	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
)

const (
	timestampTolerance          = 30 * time.Second
	challengeTimestampTolerance = 60 * time.Second
)

// VerifyEvaluationResult process the evaluation report from the challenger by the observers
func (task SCTask) VerifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message) (types.Message, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyEvaluationResult").WithField("sc_challenge_id", incomingEvaluationResult.ChallengeID)

	if err := task.validateVerifyingStorageChallengeEvaluationReport(ctx, incomingEvaluationResult); err != nil {
		return types.Message{}, err
	}
	receivedAt := time.Now().UTC()

	//if observers then save the evaluation message
	if task.isObserver(incomingEvaluationResult.Data.Observers) {
		if err := task.StoreChallengeMessage(ctx, incomingEvaluationResult); err != nil {
			logger.WithError(err).Error("error storing challenge message")
		}

		logger.Debug("evaluation report by challenger has been stored by the observer")
	} else {
		logger.Debug("not the observer to process evaluation report")
		return types.Message{}, nil
	}

	logger.Info("Storage Challenge evaluation report verification started")

	if err := task.SCService.P2PClient.EnableKey(ctx, incomingEvaluationResult.Data.Challenge.FileHash); err != nil {
		logger.WithError(err).Error("error enabling the symbol file")
	}

	//retrieve messages from the DB against challengeID
	challengeMessage, responseMessage, err := task.retrieveChallengeAndResponseMessageForEvaluationReportVerification(ctx, incomingEvaluationResult.ChallengeID)
	if err != nil {
		logger.WithError(err).Error("error retrieving the challenge and response message for verification")
		return types.Message{}, err
	}

	if challengeMessage == nil || responseMessage == nil {
		logger.Error("unable to retrieve challenge or response message")
		return types.Message{}, errors.Errorf("unable to retrieve challenge or response message")
	}

	//Need to verify the following
	//	IsChallengeTimestampOK  bool      `json:"is_challenge_timestamp_ok"`
	//	IsProcessTimestampOK    bool      `json:"is_process_timestamp_ok"`
	//	IsEvaluationTimestampOK bool      `json:"is_evaluation_timestamp_ok"`
	//	IsRecipientSignatureOK  bool      `json:"is_recipient_signature_ok"`
	//	IsChallengerSignatureOK bool      `json:"is_challenger_signature_ok"`
	//	IsEvaluationResultOK    bool      `json:"is_evaluation_result_ok"`

	//Verify message signatures
	isChallengerSignatureOk, isRecipientSignatureOk, err := task.verifyMessageSignaturesForEvaluation(ctx, *challengeMessage, *responseMessage)
	if err != nil {
		logger.WithError(err).WithField("challenge_id", incomingEvaluationResult.ChallengeID).
			Error("error verifying signatures for evaluation")
		return types.Message{}, errors.Errorf("error verifying signatures for evalaution")
	}

	//Verify message timestamps
	isChallengeTSOk, isResponseTSOk, isEvaluationTSOk := task.verifyTimestampsForEvaluation(ctx, *challengeMessage,
		*responseMessage, incomingEvaluationResult, receivedAt)

	//Verify Merkelroot and BlockNumber
	isChallengeBlockAndMROk, isResponseBlockAndMROk := task.verifyMerkelrootAndBlockNum(ctx, *challengeMessage, *responseMessage, incomingEvaluationResult)

	//Get the file assuming we host it locally (if not, return)
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingEvaluationResult.Data.Challenge.FileHash, true)
	if err != nil {
		logger.WithError(err).Error("could not read queries file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())

		observerEvaluation := types.Message{
			MessageType: types.AffirmationMessageType,
			ChallengeID: incomingEvaluationResult.ChallengeID,
			Data: types.MessageData{
				ChallengerID: incomingEvaluationResult.Data.ChallengerID,
				RecipientID:  incomingEvaluationResult.Data.RecipientID,
				Observers:    append([]string(nil), incomingEvaluationResult.Data.Observers...),
				Challenge: types.ChallengeData{
					Block:      incomingEvaluationResult.Data.Challenge.Block,
					Merkelroot: incomingEvaluationResult.Data.Challenge.Merkelroot,
					FileHash:   incomingEvaluationResult.Data.Challenge.FileHash,
					StartIndex: incomingEvaluationResult.Data.Challenge.StartIndex,
					EndIndex:   incomingEvaluationResult.Data.Challenge.EndIndex,
					Timestamp:  incomingEvaluationResult.Data.Challenge.Timestamp,
				},
				Response: types.ResponseData{
					Block:      incomingEvaluationResult.Data.Response.Block,
					Merkelroot: incomingEvaluationResult.Data.Response.Merkelroot,
					Hash:       incomingEvaluationResult.Data.Response.Hash,
					Timestamp:  incomingEvaluationResult.Data.Response.Timestamp,
				},
				ChallengerEvaluation: types.EvaluationData{
					Block:      incomingEvaluationResult.Data.ChallengerEvaluation.Block,
					Merkelroot: incomingEvaluationResult.Data.ChallengerEvaluation.Merkelroot,
					Hash:       incomingEvaluationResult.Data.ChallengerEvaluation.Hash,
					IsVerified: incomingEvaluationResult.Data.ChallengerEvaluation.IsVerified,
					Timestamp:  time.Now().UTC(),
				},
				ObserverEvaluation: types.ObserverEvaluationData{
					IsChallengerSignatureOK: isChallengerSignatureOk && isChallengeBlockAndMROk,
					IsRecipientSignatureOK:  isRecipientSignatureOk && isResponseBlockAndMROk,
					IsChallengeTimestampOK:  isChallengeTSOk,
					IsProcessTimestampOK:    isResponseTSOk,
					IsEvaluationTimestampOK: isEvaluationTSOk,
					IsEvaluationResultOK:    challengeMessage.Data.ChallengerEvaluation.Hash == challengeMessage.Data.Response.Hash,
					TrueHash:                challengeMessage.Data.ChallengerEvaluation.Hash,
					Timestamp:               time.Now().UTC(),
				},
			},
			Sender: task.nodeID,
		}

		if err := task.StoreStorageChallengeMetric(ctx, observerEvaluation); err != nil {
			logger.WithError(err).Error("error storing storage challenge metric")
		}

		signature, _, err := task.SignMessage(ctx, observerEvaluation.Data)
		if err != nil {
			logger.WithError(err).Error("error signing evaluation response")
		}
		observerEvaluation.SenderSignature = signature

		return observerEvaluation, nil
	}

	//Compute the hash of the data at the indicated byte range
	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingEvaluationResult.Data.Challenge.StartIndex, incomingEvaluationResult.Data.Challenge.EndIndex)
	isEvaluationResultOk, reason := task.verifyEvaluationResult(ctx, incomingEvaluationResult, challengeCorrectHash)

	evaluationResultResponse := types.Message{
		MessageType: types.AffirmationMessageType,
		ChallengeID: incomingEvaluationResult.ChallengeID,
		Data: types.MessageData{
			ChallengerID: incomingEvaluationResult.Data.ChallengerID,
			RecipientID:  incomingEvaluationResult.Data.RecipientID,
			Observers:    append([]string(nil), incomingEvaluationResult.Data.Observers...),
			Challenge: types.ChallengeData{
				Block:      incomingEvaluationResult.Data.Challenge.Block,
				Merkelroot: incomingEvaluationResult.Data.Challenge.Merkelroot,
				FileHash:   incomingEvaluationResult.Data.Challenge.FileHash,
				StartIndex: incomingEvaluationResult.Data.Challenge.StartIndex,
				EndIndex:   incomingEvaluationResult.Data.Challenge.EndIndex,
				Timestamp:  incomingEvaluationResult.Data.Challenge.Timestamp,
			},
			Response: types.ResponseData{
				Block:      incomingEvaluationResult.Data.Response.Block,
				Merkelroot: incomingEvaluationResult.Data.Response.Merkelroot,
				Hash:       incomingEvaluationResult.Data.Response.Hash,
				Timestamp:  incomingEvaluationResult.Data.Response.Timestamp,
			},
			ChallengerEvaluation: types.EvaluationData{
				Block:      incomingEvaluationResult.Data.ChallengerEvaluation.Block,
				Merkelroot: incomingEvaluationResult.Data.ChallengerEvaluation.Merkelroot,
				Hash:       incomingEvaluationResult.Data.ChallengerEvaluation.Hash,
				IsVerified: incomingEvaluationResult.Data.ChallengerEvaluation.IsVerified,
				Timestamp:  incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp,
			},
			ObserverEvaluation: types.ObserverEvaluationData{
				IsChallengerSignatureOK: isChallengerSignatureOk && isChallengeBlockAndMROk,
				IsRecipientSignatureOK:  isRecipientSignatureOk && isResponseBlockAndMROk,
				IsChallengeTimestampOK:  isChallengeTSOk,
				IsProcessTimestampOK:    isResponseTSOk,
				IsEvaluationTimestampOK: isEvaluationTSOk,
				IsEvaluationResultOK:    isEvaluationResultOk,
				TrueHash:                challengeCorrectHash,
				Reason:                  reason,
				Timestamp:               time.Now().UTC(),
			},
		},
		Sender: task.nodeID,
	}

	if err := task.StoreStorageChallengeMetric(ctx, evaluationResultResponse); err != nil {
		logger.WithError(err).Error("error storing storage challenge metric")
	}

	signature, _, err := task.SignMessage(ctx, evaluationResultResponse.Data)
	if err != nil {
		logger.WithError(err).Error("error signing evaluation response")
	}
	evaluationResultResponse.SenderSignature = signature

	if err := task.StoreChallengeMessage(ctx, evaluationResultResponse); err != nil {
		logger.WithError(err).Error("error storing evaluation response ")
	}
	logger.WithField("node_id", task.nodeID).Info("evaluation report has been evaluated by the observer")

	return evaluationResultResponse, nil
}

// RetrieveChallengeMessage retrieves the challenge messages from db for further verification
func (task SCTask) RetrieveChallengeMessage(ctx context.Context, challengeID string, msgType int) (challengeMessage *types.Message, err error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		msg, err := store.QueryStorageChallengeMessage(challengeID, msgType)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error retrieving challenge message from DB")
			return nil, err
		}
		log.WithContext(ctx).Debug("storage challenge msg retrieved")

		var data types.MessageData
		err = json.Unmarshal(msg.Data, &data)
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("Error un-marshaling message data")
			return nil, err
		}

		return &types.Message{
			ChallengeID:     msg.ChallengeID,
			MessageType:     types.MessageType(msg.MessageType),
			Data:            data,
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
			CreatedAt:       msg.CreatedAt,
			UpdatedAt:       msg.UpdatedAt,
		}, nil
	}

	return nil, nil
}

func (task SCTask) verifyMessageSignaturesForEvaluation(ctx context.Context, challengeMessage, responseMessage types.Message) (challengeMsgSig, ResMsgSig bool, err error) {
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

func (task SCTask) verifyTimestampsForEvaluation(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.Message, evaluationMsgRecvTime time.Time) (challengeMsgTS, resMsgTS, evalMsgTS bool) {
	isChallengeTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Challenge.Timestamp.UTC(),
		challengeMessage.Data.Challenge.Timestamp.UTC(), challengeMessage.CreatedAt.UTC(),
		challengeMessage.ChallengeID, challengeMessage.MessageType)
	log.WithContext(ctx).Debug("challenge message timestamps have been evaluated successfully")

	isResponseTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Response.Timestamp.UTC(),
		responseMessage.Data.Response.Timestamp.UTC(), responseMessage.CreatedAt.UTC(),
		responseMessage.ChallengeID, responseMessage.MessageType)
	log.WithContext(ctx).Debug("response message timestamps have been evaluated successfully")

	differenceBetweenEvaluationMsgSendAndRecvTime := incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp.Sub(evaluationMsgRecvTime)
	if differenceBetweenEvaluationMsgSendAndRecvTime < 0 {
		differenceBetweenEvaluationMsgSendAndRecvTime = -differenceBetweenEvaluationMsgSendAndRecvTime
	}

	var evaluationReportTSOk bool
	if differenceBetweenEvaluationMsgSendAndRecvTime <= timestampTolerance {
		log.WithContext(ctx).Debug("evaluation message timestamp is evaluated successfully")
		evaluationReportTSOk = true
	} else {
		log.WithContext(ctx).Debug("the time diff between the evaluation report send and receive time, exceeds the tolerance limit")
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
	log.WithContext(ctx).Debug("evaluation report & overall time taken by challenge has been evaluated successfully")

	return isChallengeTSOk, isResponseTSOk, isEvaluationTSOk
}

func (task SCTask) verifyMessageTimestamps(ctx context.Context, timeWhenMsgSent, sentTimeWhenMsgReceived, timeWhenMsgStored time.Time, challengeID string, messageType types.MessageType) bool {
	logger := log.WithContext(ctx).
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

func (task SCTask) verifyMerkelrootAndBlockNum(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.Message) (isChallengeOk, isResOk bool) {
	logger := log.WithContext(ctx).
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

func (task SCTask) verifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message, correctHash string) (isAffirmationOk bool, reason string) {
	logger := log.WithContext(ctx).WithField("challenge_id", incomingEvaluationResult.ChallengeID)

	evaluationResultData := incomingEvaluationResult.Data

	//recipientHash = true, challengerHash = true, observerAffirmation = true
	if evaluationResultData.ChallengerEvaluation.Hash == correctHash &&
		evaluationResultData.ChallengerEvaluation.IsVerified &&
		evaluationResultData.Response.Hash == correctHash {
		isAffirmationOk = true

		logger.Debug("both hash by challenger and recipient were correct, and challenger evaluates correctly")
	}

	//recipientHash = false, challengerHash = false, observerAffirmation = true
	if evaluationResultData.Response.Hash != correctHash && !evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = true

		logger.Debug("recipient hash was not correct, and challenger evaluates correctly")
	}

	//recipientHash = true, challengerHash = false, observerAffirmation = false
	if evaluationResultData.Response.Hash == correctHash &&
		!evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = false

		reason = "recipient hash was correct, but challenger evaluated to false"

		logger.Debug("recipient hash was correct, but challenger evaluates it to false")
	}

	if evaluationResultData.Response.Hash != correctHash && evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = false

		reason = "recipient hash was not correct, but challenger evaluated it to true"

		logger.Debug("recipient hash was not correct, but challenger evaluates it to true")
	}

	return isAffirmationOk, reason
}

func (task SCTask) validateVerifyingStorageChallengeEvaluationReport(ctx context.Context, evaluationReport types.Message) error {
	if evaluationReport.MessageType != types.EvaluationMessageType {
		return fmt.Errorf("incorrect message type to verify evaluation report")
	}

	if err := task.StoreStorageChallengeMetric(ctx, evaluationReport); err != nil {
		log.WithContext(ctx).WithField("challenge_id", evaluationReport.ChallengeID).
			WithField("message_type", evaluationReport.MessageType).Error(
			"error storing storage challenge metric")
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

func (task SCTask) retrieveChallengeAndResponseMessageForEvaluationReportVerification(ctx context.Context, challengeID string) (cMsg, rMsg *types.Message, er error) {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 15 * time.Second
	b.InitialInterval = 1 * time.Second

	var challengeMessage *types.Message
	var responseMessage *types.Message

	// Retry fetching challenge message
	err := backoff.Retry(func() error {
		var err error
		challengeMessage, err = task.RetrieveChallengeMessage(ctx, challengeID, int(types.ChallengeMessageType))
		if err != nil {
			if err == sql.ErrNoRows {
				return err // Will retry
			}
			return backoff.Permanent(err) // Won't retry
		}
		return nil // Success
	}, backoff.WithMaxRetries(b, 5))

	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to retrieve challenge message")
		return nil, nil, errors.Errorf("failed to retrieve challenge message for evaluation affirmation")
	}

	// Reset backoff before next operation
	b.Reset()

	// Retry fetching response message
	err = backoff.Retry(func() error {
		var err error
		responseMessage, err = task.RetrieveChallengeMessage(ctx, challengeID, int(types.ResponseMessageType))
		if err != nil {
			if err == sql.ErrNoRows {
				return err // Will retry
			}
			return backoff.Permanent(err) // Won't retry
		}
		return nil // Success
	}, backoff.WithMaxRetries(b, 5))

	if err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to retrieve challenge message")
		return nil, nil, errors.Errorf("failed to retrieve challenge message for evaluation affirmation")
	}

	return challengeMessage, responseMessage, nil
}
