package storagechallenge

import (
	"context"
	"encoding/json"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"time"
)

const (
	timestampTolerance = 5 * time.Second
)

// VerifyEvaluationResult process the evaluation report from the challenger by the observers
func (task *SCTask) VerifyEvaluationResult(ctx context.Context, incomingEvaluationResult types.Message) (types.Message, error) {
	logger := log.WithContext(ctx).WithField("method", "VerifyEvaluationResult").WithField("challengeID", incomingEvaluationResult.ChallengeID)
	logger.Debug("Start verifying challenger's evaluation report") // Incoming evaluation report validation

	receivedAt := time.Now()

	//retrieve messages from the DB against challengeID
	challengeMessage, err := task.RetrieveChallengeMessage(ctx, incomingEvaluationResult.ChallengeID, int(types.ChallengeMessageType))
	if err != nil {
		err := errors.Errorf("error retrieving challenge message for evaluation")
		log.WithContext(ctx).WithError(err).Error(err.Error())
		return types.Message{}, err
	}

	responseMessage, err := task.RetrieveChallengeMessage(ctx, incomingEvaluationResult.ChallengeID, int(types.ResponseMessageType))
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving response message for evaluation")
		return types.Message{}, errors.Errorf("error retrieving response message for evaluation")
	}

	if challengeMessage == nil || responseMessage == nil {
		log.WithContext(ctx).Error("unable to retrieve challenge or response message")
		return types.Message{}, errors.Errorf("unable to retrieve storage or challenge message")
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
		log.WithContext(ctx).WithError(err).Error("error verifying signatures for evaluation")
		return types.Message{}, errors.Errorf("error verifying signatures for evalaution")
	}

	//Verify message timestamps
	isChallengeTSOk, isResponseTSOk, isEvaluationTSOk := task.verifyTimestampsForEvaluation(*challengeMessage,
		*responseMessage, incomingEvaluationResult, receivedAt)

	//Verify Evaluation Result

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

	var isEvaluationResultOk bool
	if incomingEvaluationResult.Data.ChallengerEvaluation.Hash == challengeCorrectHash {
		isEvaluationResultOk = true
	}

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
				IsChallengerSignatureOK: isChallengerSignatureOk,
				IsRecipientSignatureOK:  isRecipientSignatureOk,
				IsChallengeTimestampOK:  isChallengeTSOk,
				IsProcessTimestampOK:    isResponseTSOk,
				IsEvaluationTimestampOK: isEvaluationTSOk,
				IsEvaluationResultOK:    isEvaluationResultOk,
				TrueHash:                challengeCorrectHash,
				Timestamp:               time.Now(),
			},
		},
		Sender: task.nodeID,
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

func (task *SCTask) verifyTimestampsForEvaluation(challengeMessage, responseMessage, incomingEvaluationResult types.Message, evaluationMsgRecvTime time.Time) (challengeMsgTS, resMsgTS, evalMsgTS bool) {
	isChallengeTSOk := task.verifyMessageTimestamps(incomingEvaluationResult.Data.Challenge.Timestamp, challengeMessage.Data.Challenge.Timestamp, challengeMessage.CreatedAt)

	isResponseTSOk := task.verifyMessageTimestamps(incomingEvaluationResult.Data.Response.Timestamp, challengeMessage.Data.Response.Timestamp, responseMessage.CreatedAt)

	difference := incomingEvaluationResult.Data.ChallengerEvaluation.Timestamp.Sub(evaluationMsgRecvTime)
	if difference < 0 {
		difference = -difference
	}

	isEvaluationTSOk := difference <= timestampTolerance

	return isChallengeTSOk, isResponseTSOk, isEvaluationTSOk
}

func (task *SCTask) verifyMessageTimestamps(timeWhenMsgSent, sentTimeWhenMsgReceived, timeWhenMsgStored time.Time) bool {
	difference := timeWhenMsgStored.Sub(timeWhenMsgSent)
	if difference < 0 {
		difference = -difference
	}

	if timeWhenMsgSent.Equal(sentTimeWhenMsgReceived) {
		return difference <= timestampTolerance
	}

	return false
}
