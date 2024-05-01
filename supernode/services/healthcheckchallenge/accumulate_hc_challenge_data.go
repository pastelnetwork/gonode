package healthcheckchallenge

import (
	"context"
	"database/sql"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
)

// AccumulateHChallengeScoreData accumulates the healthcheck-challenge score data
func (task *HCTask) AccumulateHChallengeScoreData(ctx context.Context, challengeID string) error {
	logger := log.WithContext(ctx).WithField("method", "AccumulateHChallengeScoreData")

	msgs, err := task.getAffirmationMessages(ctx, challengeID)
	if err != nil {
		logger.WithError(err).Error("error retrieving affirmation messages")
		return errors.Errorf("error retrieving affirmation messages")
	}

	pingInfos, err := task.getPingInfos(ctx)
	if err != nil {
		logger.WithError(err).Error("error retrieving ping-infos for hc score aggregation")
		return errors.Errorf("error retrieving ping-infos for hc score aggregation")
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
		logger.WithError(err).Error("error accumulating challenger summary for hc score aggregation")
		return err
	}

	err = task.processRecipientEvaluation(recipientEvaluations, recipientID, successThreshold, pingInfos)
	if err != nil {
		logger.WithError(err).Error("error accumulating recipient summary for hc score aggregation")
		return err
	}

	for _, msg := range msgs {
		err = task.processObserverEvaluation(commonEval, msg.Data.ObserverEvaluation.IsEvaluationTimestampOK, msg.Sender, pingInfos)
		if err != nil {
			logger.WithField("node_id", msg.Sender).WithError(err).Error("error accumulating observer data for hc score aggregation")
		}
	}

	return nil
}
func (task *HCTask) getPingInfos(ctx context.Context) (types.PingInfos, error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	pingInfos, err := store.GetAllPingInfos()
	if err != nil {
		return pingInfos, err
	}

	return pingInfos, nil
}

func (task *HCTask) getAffirmationMessages(ctx context.Context, challengeID string) (types.HealthCheckChallengeMessages, error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	affMsgs, err := store.GetHCMetricsByChallengeIDAndMessageType(challengeID, types.HealthCheckAffirmationMessageType)
	if err != nil {
		return nil, err
	}

	var affirmationMsgs types.HealthCheckChallengeMessages
	for _, msg := range affMsgs {
		var hcMsg types.HealthCheckMessage

		err := json.Unmarshal(msg.Data, &hcMsg.Data)
		if err != nil {
			return nil, errors.Errorf("error unmarshaling affirmation msg")
		}

		affirmationMsgs = append(affirmationMsgs, types.HealthCheckMessage{
			ChallengeID:     challengeID,
			MessageType:     types.HealthCheckMessageType(msg.MessageType),
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
			Data:            hcMsg.Data,
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
	aggregatedScoreData, err := task.ScoreStore.GetAccumulativeHCData(challengerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, challengerID)

			aggregatedScoreData = types.AccumulativeChallengeData{
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

	if err := task.ScoreStore.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return err
	}

	return nil
}

func (task *HCTask) processRecipientEvaluation(recipientEvaluations int, recipientID string, successThreshold int, infos types.PingInfos) error {
	aggregatedScoreData, err := task.ScoreStore.GetAccumulativeHCData(recipientID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, recipientID)

			aggregatedScoreData = types.AccumulativeChallengeData{
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

	if err := task.ScoreStore.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return err
	}

	return nil
}

func (task *HCTask) processObserverEvaluation(commonEval, observerEval bool, observerID string, infos types.PingInfos) error {
	aggregatedScoreData, err := task.ScoreStore.GetAccumulativeHCData(observerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			nodeID, nodeIP := getNodeInfo(infos, observerID)

			aggregatedScoreData = types.AccumulativeChallengeData{
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

	if err := task.ScoreStore.UpsertAccumulativeHCData(aggregatedScoreData); err != nil {
		return err
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
