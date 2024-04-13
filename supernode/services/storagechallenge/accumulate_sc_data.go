package storagechallenge

import (
	"context"
	"database/sql"
	json "github.com/json-iterator/go"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/types"
)

// AccumulateStorageChallengeScoreData accumulates the storage-challenge score data
func (task *SCTask) AccumulateStorageChallengeScoreData(ctx context.Context, challengeID string) error {
	logger := log.WithContext(ctx).WithField("method", "AccumulateStorageChallengeScoreData")

	msgs, err := task.getAffirmationMessages(ctx, challengeID)
	if err != nil {
		logger.WithError(err).Error("error retrieving affirmation messages")
		return errors.Errorf("error retrieving affirmation messages")
	}

	pingInfos, err := task.getPingInfos(ctx)
	if err != nil {
		logger.WithError(err).Error("error retrieving ping infos for sc score aggregation")
		return errors.Errorf("error retrieving ping infos for sc score aggregation")
	}

	var (
		challengerEvaluations int
		challengerID          string
		recipientEvaluations  int
		recipientID           string
		observerEvaluation    = make(map[string]int)
	)

	for _, msg := range msgs {
		logger.WithField("msg", msg.Data).Debug("msg processing")

		if IsChallengerCorrectlyEvaluated(msg.Data) {
			challengerID = msg.Data.ChallengerID
			challengerEvaluations++
		}

		if IsRecipientCorrectlyEvaluated(msg.Data) {
			recipientID = msg.Data.RecipientID
			recipientEvaluations++
		}

		observerEvaluation[msg.Data.ObserverEvaluation.TrueHash] = observerEvaluation[msg.Data.ObserverEvaluation.TrueHash] + 1
	}

	commonHash := getCorrectHash(observerEvaluation)
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
		err = task.processObserverEvaluation(commonHash, msg.Data.ObserverEvaluation.TrueHash, msg.Sender, pingInfos)
		if err != nil {
			logger.WithField("node_id", msg.Sender).WithError(err).Error("error accumulating observer data for sc score aggregation")
		}
	}

	return nil
}

func (task *SCTask) getAffirmationMessages(ctx context.Context, challengeID string) (types.StorageChallengeMessages, error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	affMsgs, err := store.GetMetricsByChallengeIDAndMessageType(challengeID, types.AffirmationMessageType)
	if err != nil {
		return nil, err
	}

	var affirmationMsgs types.StorageChallengeMessages
	for _, msg := range affMsgs {
		var scMsg types.Message

		err := json.Unmarshal(msg.Data, &scMsg.Data)
		if err != nil {
			return nil, errors.Errorf("error unmarshaling affirmation msg")
		}

		affirmationMsgs = append(affirmationMsgs, types.Message{
			ChallengeID:     challengeID,
			MessageType:     types.MessageType(msg.MessageType),
			Sender:          msg.Sender,
			SenderSignature: msg.SenderSignature,
			Data:            scMsg.Data,
		})
	}

	return affirmationMsgs, nil
}

func IsChallengerCorrectlyEvaluated(data types.MessageData) bool {
	return data.ObserverEvaluation.IsChallengeTimestampOK &&
		data.ObserverEvaluation.IsChallengerSignatureOK &&
		data.ObserverEvaluation.IsEvaluationResultOK
}

func IsRecipientCorrectlyEvaluated(data types.MessageData) bool {
	return data.ObserverEvaluation.IsProcessTimestampOK &&
		data.ObserverEvaluation.IsRecipientSignatureOK &&
		data.Response.Hash == data.ObserverEvaluation.TrueHash
}

func getCorrectHash(hashMap map[string]int) (correctHash string) {
	var (
		highgestTally  int
		mostCommonHash string
	)

	for hash, tally := range hashMap {
		if tally > highgestTally {
			highgestTally = tally
			mostCommonHash = hash
		}
	}

	return mostCommonHash
}

func (task *SCTask) processChallengerEvaluation(challengerEvaluations int, challengerID string, successThreshold int, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeSCData(challengerID)
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

	aggregatedScoreData.TotalChallengesAsChallengers = aggregatedScoreData.TotalChallengesAsChallengers + 1
	if challengerEvaluations >= successThreshold {
		aggregatedScoreData.CorrectChallengerEvaluations = aggregatedScoreData.CorrectChallengerEvaluations + 1
	}

	if err := task.historyDB.UpsertAccumulativeSCData(aggregatedScoreData); err != nil {
		return err
	}

	return nil
}

func (task *SCTask) processRecipientEvaluation(recipientEvaluations int, recipientID string, successThreshold int, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeSCData(recipientID)
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

	if err := task.historyDB.UpsertAccumulativeSCData(aggregatedScoreData); err != nil {
		return err
	}

	return nil
}

func (task *SCTask) processObserverEvaluation(commonHash string, observerTrueHash, observerID string, infos types.PingInfos) error {
	aggregatedScoreData, err := task.historyDB.GetAccumulativeSCData(observerID)
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
	if observerTrueHash == commonHash {
		aggregatedScoreData.CorrectObserverEvaluations++
	}

	if err := task.historyDB.UpsertAccumulativeSCData(aggregatedScoreData); err != nil {
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

func (task *SCTask) getPingInfos(ctx context.Context) (types.PingInfos, error) {
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
