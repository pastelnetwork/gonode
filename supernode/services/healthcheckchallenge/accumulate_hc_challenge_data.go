package healthcheckchallenge

import (
	"context"
	"database/sql"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// AccumulateHChallengeScoreData accumulates the healthcheck-challenge score data
func (task *HCTask) AccumulateHChallengeScoreData(ctx context.Context, challengeID string) error {
	logger := log.WithContext(ctx).WithField("method", "AccumulateHChallengeScoreData")

	msgs, err := task.getAffirmationMessages(challengeID)
	if err != nil {
		logger.WithError(err).Error("error retrieving affirmation messages")
		return errors.Errorf("error retrieving affirmation messages")
	}

	pingInfos, err := task.historyDB.GetAllPingInfos()
	if err != nil {
		logger.WithError(err).Error("error retrieving affirmation messages")
		return errors.Errorf("error retrieving affirmation messages")
	}

	var (
		challengerEvaluations int
		challengerID          string
		recipientEvaluations  int
		recipientID           string
		observerEvaluation    = make(map[bool]int)
	)

	for _, msg := range msgs {
		if IsChallengerCorrectlyEvaluated(msg.Data) {
			challengerID = msg.Data.ChallengerID
			challengerEvaluations++
		}

		if IsRecipientCorrectlyEvaluated(msg.Data) {
			recipientID = msg.Data.RecipientID
			recipientEvaluations++
		}

		observerEvaluation[msg.Data.ObserverEvaluation.IsEvaluationResultOK] = observerEvaluation[msg.Data.ObserverEvaluation.IsEvaluationTimestampOK] + 1
	}

	commonEval := getCommonEvaluationResult(observerEvaluation)
	successThreshold := len(msgs) - 1

	err = task.processChallengerEvaluation(challengerEvaluations, challengerID, successThreshold, pingInfos)
	if err != nil {
		logger.WithError(err).Error("error accumulating challenger summary for sc score aggregation")
		return err
	}

	err = task.processRecipientEvaluation(recipientEvaluations, recipientID, successThreshold, pingInfos)
	if err != nil {
		logger.WithError(err).Error("error accumulating recipient summary for sc score aggregation")
		return err
	}

	for _, msg := range msgs {
		err = task.processObserverEvaluation(commonEval, msg.Data.ObserverEvaluation.IsEvaluationTimestampOK, msg.Sender, pingInfos)
		if err != nil {
			logger.WithField("node_id", msg.Sender).WithError(err).Error("error accumulating observer data for sc score aggregation")
		}
	}

	return nil
}

func (task *HCTask) getAffirmationMessages(challengeID string) (types.HealthCheckChallengeMessages, error) {
	affMsgs, err := task.historyDB.GetHCMetricsByChallengeIDAndMessageType(challengeID, types.AffirmationMessageType)
	if err != nil {
		return nil, errors.Errorf("error retrieving affirmation msgs")
	}

	var affirmationMsgs types.HealthCheckChallengeMessages
	for _, msg := range affMsgs {
		var scMsg types.Message

		err := json.Unmarshal(msg.Data, &scMsg.Data)
		if err != nil {
			return nil, errors.Errorf("error unmarshaling affirmation msg")
		}

		affirmationMsgs = append(affirmationMsgs, types.HealthCheckMessage{
			ChallengeID:     challengeID,
			MessageType:     types.HealthCheckMessageType(msg.MessageType),
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
		})
	}

	return affirmationMsgs, nil
}

func IsChallengerCorrectlyEvaluated(data types.HealthCheckMessageData) bool {
	return data.ObserverEvaluation.IsChallengeTimestampOK &&
		data.ObserverEvaluation.IsChallengerSignatureOK &&
		data.ObserverEvaluation.IsEvaluationResultOK
}

func IsRecipientCorrectlyEvaluated(data types.HealthCheckMessageData) bool {
	return data.ObserverEvaluation.IsProcessTimestampOK &&
		data.ObserverEvaluation.IsRecipientSignatureOK &&
		data.ObserverEvaluation.IsEvaluationResultOK
}

func getCommonEvaluationResult(evaluationMap map[bool]int) bool {
	var (
		highgestTally        int
		mostCommonEvaluation bool
	)

	for evaluationResult, tally := range evaluationMap {
		if tally > highgestTally {
			highgestTally = tally
			mostCommonEvaluation = evaluationResult
		}
	}

	return mostCommonEvaluation
}

func (task *HCTask) processChallengerEvaluation(challengerEvaluations int, challengerID string, successThreshold int, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeHCData(challengerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, challengerID)

			aggregatedScoreData = types.AggregatedSCScore{
				NodeID:    nodeID,
				IPAddress: nodeIP,
			}
		} else {
			return err
		}
	}

	if challengerEvaluations >= successThreshold {
		aggregatedScoreData.CorrectChallengerEvaluations = aggregatedScoreData.CorrectChallengerEvaluations + 1
	}

	aggregatedScoreData.TotalChallengesAsChallengers = aggregatedScoreData.TotalChallengesAsChallengers + 1

	if err := task.historyDB.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return errors.Errorf("error retrieving aggregated score data")
	}

	return nil
}

func (task *HCTask) processRecipientEvaluation(recipientEvaluations int, recipientID string, successThreshold int, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeHCData(recipientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, recipientID)

			aggregatedScoreData = types.AggregatedSCScore{
				NodeID:    nodeID,
				IPAddress: nodeIP,
			}
		} else {
			return err
		}
	}

	if recipientEvaluations >= successThreshold {
		aggregatedScoreData.CorrectRecipientEvaluations = aggregatedScoreData.CorrectRecipientEvaluations + 1
	}

	aggregatedScoreData.TotalChallengesAsRecipients = aggregatedScoreData.TotalChallengesAsRecipients + 1

	if err := task.historyDB.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return errors.Errorf("error retrieving aggregated score data")
	}

	return nil
}

func (task *HCTask) processObserverEvaluation(commonEval, observerEval bool, observerID string, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeHCData(observerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, observerID)

			aggregatedScoreData = types.AggregatedSCScore{
				NodeID:    nodeID,
				IPAddress: nodeIP,
			}
		} else {
			return err
		}
	}

	aggregatedScoreData.TotalChallengesAsObservers++

	if commonEval == observerEval {
		aggregatedScoreData.CorrectObserverEvaluations++
	}

	if err := task.historyDB.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return errors.Errorf("error retrieving aggregated score data")
	}

	return nil
}

func getNodeInfo(infos types.PingInfos, nodeID string) (id, nodeIP string) {
	for _, info := range infos {
		if info.SupernodeID == nodeID {
			return nodeID, info.IPAddress
		}
	}

	return nodeID, ""
}
