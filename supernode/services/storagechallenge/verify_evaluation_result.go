package storagechallenge

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
)

const (
	timestampTolerance          = 5 * time.Second
	challengeTimestampTolerance = 15 * time.Second
)

// VerifyEvaluationResult process the evaluation report from the challenger by the observers
func (task *SCTask) VerifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message) (types.Message, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyEvaluationResult").WithField("challengeID", incomingEvaluationResult.ChallengeID)
	logger.Debug("Start verifying challenger's evaluation report") // Incoming evaluation report validation

	receivedAt := time.Now()

	//if observers then save the evaluation message
	if task.isObserver(incomingEvaluationResult.Data.Observers) {
		logger := log.WithContext(ctx).WithField("node_id", task.nodeID)

		if err := task.StoreChallengeMessage(ctx, incomingEvaluationResult); err != nil {
			log.WithContext(ctx).
				WithField("node_id", task.nodeID).
				WithError(err).
				Error("error storing challenge message")
		}

		logger.Info("evaluation report by challenger has been stored by the observer")
	} else {
		log.WithContext(ctx).WithField("node_id", task.nodeID).Info("not the observer to process evaluation report")
		return types.Message{}, nil
	}

	//retrieve messages from the DB against challengeID
	challengeMessage, err := task.RetrieveChallengeMessage(ctx, incomingEvaluationResult.ChallengeID, int(types.ChallengeMessageType))
	if err != nil {
		err := errors.Errorf("error retrieving challenge message for evaluation")
		log.WithContext(ctx).WithField("challenge_id", incomingEvaluationResult.ChallengeID).WithError(err).Error(err.Error())
		return types.Message{}, err
	}

	responseMessage, err := task.RetrieveChallengeMessage(ctx, incomingEvaluationResult.ChallengeID, int(types.ResponseMessageType))
	if err != nil {
		log.WithContext(ctx).WithField("challenge_id", incomingEvaluationResult.ChallengeID).WithError(err).
			Error("error retrieving response message for evaluation")
		return types.Message{}, errors.Errorf("error retrieving response message for evaluation")
	}

	if challengeMessage == nil || responseMessage == nil {
		log.WithContext(ctx).Error("unable to retrieve challenge or response message")
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
		log.WithContext(ctx).WithError(err).WithField("challenge_id", incomingEvaluationResult.ChallengeID).
			Error("error verifying signatures for evaluation")
		return types.Message{}, errors.Errorf("error verifying signatures for evalaution")
	}

	//Verify message timestamps
	isChallengeTSOk, isResponseTSOk, isEvaluationTSOk := task.verifyTimestampsForEvaluation(ctx, *challengeMessage,
		*responseMessage, incomingEvaluationResult, receivedAt)

	//Verify Merkelroot and BlockNumber
	isChallengeBlockAndMROk, isResponseBlockAndMROk := task.verifyMerkelrootAndBlockNum(ctx, *challengeMessage, *responseMessage, incomingEvaluationResult)

	//Get the file assuming we host it locally (if not, return)
	log.WithContext(ctx).Info("getting the file from hash to verify evaluation report")
	challengeFileData, err := task.GetSymbolFileByKey(ctx, incomingEvaluationResult.Data.Challenge.FileHash, true)
	if err != nil {
		log.WithContext(ctx).WithField("method", "VerifyEvaluationResult").WithField("challengeID", incomingEvaluationResult.ChallengeID).Error("could not read local file data in to memory, so not continuing with verification.", "file.ReadFileIntoMemory", err.Error())
		return types.Message{}, err
	}
	log.WithContext(ctx).Info("file has been retrieved for verification of evaluation report")

	//Compute the hash of the data at the indicated byte range
	log.WithContext(ctx).Info("generating hash for the data against given indices")
	challengeCorrectHash := task.computeHashOfFileSlice(challengeFileData, incomingEvaluationResult.Data.Challenge.StartIndex, incomingEvaluationResult.Data.Challenge.EndIndex)
	log.WithContext(ctx).Info("hash of the data has been generated against the given indices")

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
				Timestamp:               time.Now(),
			},
		},
		Sender: task.nodeID,
	}

	signature, _, err := task.SignMessage(ctx, evaluationResultResponse.Data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error signing evaluation response")
	}
	evaluationResultResponse.SenderSignature = signature

	if err := task.SCService.P2PClient.EnableKey(ctx, evaluationResultResponse.Data.Challenge.FileHash); err != nil {
		log.WithContext(ctx).WithError(err).Error("error enabling the symbol file")
		return evaluationResultResponse, err
	}

	return evaluationResultResponse, nil
}

// RetrieveChallengeMessage retrieves the challenge messages from db for further verification
func (task SCTask) RetrieveChallengeMessage(ctx context.Context, challengeID string, msgType int) (challengeMessage *types.Message, err error) {
	store, err := local.OpenHistoryDB()
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
		log.WithContext(ctx).Info("storage challenge msg retrieved")

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

func (task *SCTask) verifyMessageSignaturesForEvaluation(ctx context.Context, challengeMessage, responseMessage types.Message) (challengeMsgSig, ResMsgSig bool, err error) {
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

func (task *SCTask) verifyTimestampsForEvaluation(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.Message, evaluationMsgRecvTime time.Time) (challengeMsgTS, resMsgTS, evalMsgTS bool) {
	isChallengeTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Challenge.Timestamp,
		challengeMessage.Data.Challenge.Timestamp, challengeMessage.CreatedAt,
		challengeMessage.ChallengeID, challengeMessage.MessageType)

	isResponseTSOk := task.verifyMessageTimestamps(ctx, incomingEvaluationResult.Data.Response.Timestamp,
		challengeMessage.Data.Response.Timestamp, responseMessage.CreatedAt,
		responseMessage.ChallengeID, responseMessage.MessageType)

	differenceBetweenEvaluationMsgSendAndRecvTime := incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp.Sub(evaluationMsgRecvTime)
	if differenceBetweenEvaluationMsgSendAndRecvTime < 0 {
		differenceBetweenEvaluationMsgSendAndRecvTime = -differenceBetweenEvaluationMsgSendAndRecvTime
	}

	var evaluationReportTSOk bool
	if differenceBetweenEvaluationMsgSendAndRecvTime <= timestampTolerance {
		log.WithContext(ctx).Info("evaluation message timestamp is verified successfully")
		evaluationReportTSOk = true
	} else {
		log.WithContext(ctx).Info("the time diff between the evaluation report send and receive time, exceeds the tolerance limit")
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

	return isChallengeTSOk, isResponseTSOk, isEvaluationTSOk
}

func (task *SCTask) verifyMessageTimestamps(ctx context.Context, timeWhenMsgSent, sentTimeWhenMsgReceived, timeWhenMsgStored time.Time, challengeID string, messageType types.MessageType) bool {
	logger := log.WithContext(ctx).
		WithField("challenge_id", challengeID).
		WithField("message_type", messageType.String()).
		WithField("node_id", task.nodeID)

	difference := timeWhenMsgStored.Sub(timeWhenMsgSent)
	if difference < 0 {
		difference = -difference
	}

	var isTSOk bool
	if timeWhenMsgSent.Equal(sentTimeWhenMsgReceived) {

		if difference <= timestampTolerance {
			isTSOk = true
		} else {
			logger.Info("the time difference when message has been sent and when message has been stored, " +
				"exceeds the tolerance limit")

			isTSOk = false
		}

	} else {
		logger.Info("time when message sent and sent time when message received are not same")
	}
	logger.Info("timestamps have been verified successfully")

	return isTSOk
}

func (task *SCTask) verifyMerkelrootAndBlockNum(ctx context.Context, challengeMessage, responseMessage, incomingEvaluationResult types.Message) (isChallengeOk, isResOk bool) {
	logger := log.WithContext(ctx).
		WithField("challenge_id", incomingEvaluationResult.ChallengeID).
		WithField("node_id", task.nodeID)

	isChallengeBlockOk := challengeMessage.Data.Challenge.Block == incomingEvaluationResult.Data.Challenge.Block
	if !isChallengeBlockOk {
		logger.Info("challenge msg block num is different and not same when challenge sent and received")
	}

	isChallengeMerkelrootOk := challengeMessage.Data.Challenge.Merkelroot == incomingEvaluationResult.Data.Challenge.Merkelroot
	if !isChallengeMerkelrootOk {
		logger.Info("challenge msg merkelroot is different and not same when challenge sent and received")
	}

	isResponseBlockOk := responseMessage.Data.Challenge.Block == incomingEvaluationResult.Data.Challenge.Block
	if !isResponseBlockOk {
		logger.Info("response msg block num is different and not same when challenge sent and received")
	}

	isResponseMerkelrootOk := responseMessage.Data.Challenge.Merkelroot == incomingEvaluationResult.Data.Challenge.Merkelroot
	if !isResponseMerkelrootOk {
		logger.Info("response msg merkelroot is different and not same when challenge sent and received")
	}

	return isChallengeBlockOk && isChallengeMerkelrootOk, isResponseBlockOk && isResponseMerkelrootOk
}

func (task *SCTask) verifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message, correctHash string) (isAffirmationOk bool, reason string) {
	logger := log.WithContext(ctx).WithField("challenge_id", incomingEvaluationResult.ChallengeID)

	evaluationResultData := incomingEvaluationResult.Data

	//recipientHash = true, challengerHash = true, observerAffirmation = true
	if evaluationResultData.ChallengerEvaluation.Hash == correctHash &&
		evaluationResultData.ChallengerEvaluation.IsVerified &&
		evaluationResultData.Response.Hash == correctHash {
		isAffirmationOk = true

		logger.Info("both hash by challenger and recipient were correct, and challenger evaluates correctly")
	}

	//recipientHash = false, challengerHash = false, observerAffirmation = true
	if evaluationResultData.Response.Hash != correctHash && !evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = true

		logger.Info("recipient hash was not correct, and challenger evaluates correctly")
	}

	//recipientHash = true, challengerHash = false, observerAffirmation = false
	if evaluationResultData.Response.Hash == correctHash &&
		!evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = false

		reason = "recipient hash was correct, but challenger evaluated to false"

		logger.Info("recipient hash was correct, but challenger evaluates it to false")
	}

	if evaluationResultData.Response.Hash != correctHash && evaluationResultData.ChallengerEvaluation.IsVerified {
		isAffirmationOk = false

		reason = "recipient hash was not correct, but challenger evaluated it to true"

		logger.Info("recipient hash was not correct, but challenger evaluates it to true")
	}

	return isAffirmationOk, reason
}
