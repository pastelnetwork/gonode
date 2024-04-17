package storagechallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/types"
)

// AggregateChallengesScore calculates the score based on specified weightages
func (task *SCTask) AggregateChallengesScore(ctx context.Context) error {
	logger := log.WithContext(ctx).WithField("method", "AggregateChallengesScore")
	logger.Info("invoked")

	accumulativeSCNodesData, err := task.scoreStore.GetAccumulativeSCDataForAllNodes()
	if err != nil {
		logger.WithError(err).Error("error retrieving accumulative data for nodes")
		return errors.Errorf("error retrieving accumulative data for nodes")
	}

	accumulativeHCNodesData, err := task.scoreStore.GetAccumulativeHCDataForAllNodes()
	if err != nil {
		logger.WithError(err).Error("error retrieving accumulative data for nodes")
		return errors.Errorf("error retrieving accumulative data for nodes")
	}

	mapOfAggregatedScores := make(map[string]types.AggregatedScore)
	for _, nodeData := range accumulativeSCNodesData {
		aggregatedSCScore := AggregateChallengeScore(nodeData)
		mapOfAggregatedScores[nodeData.NodeID] = types.AggregatedScore{
			NodeID:                nodeData.NodeID,
			IPAddress:             nodeData.IPAddress,
			StorageChallengeScore: aggregatedSCScore,
		}
	}

	for _, nodeData := range accumulativeHCNodesData {
		aggregatedHCScore := AggregateChallengeScore(nodeData)
		if existingScore, exists := mapOfAggregatedScores[nodeData.NodeID]; exists {
			existingScore.HealthCheckChallengeScore = aggregatedHCScore
			mapOfAggregatedScores[nodeData.NodeID] = existingScore
		} else {
			mapOfAggregatedScores[nodeData.NodeID] = types.AggregatedScore{
				NodeID:                    nodeData.NodeID,
				IPAddress:                 nodeData.IPAddress,
				HealthCheckChallengeScore: aggregatedHCScore,
			}
		}
	}

	for nodeID, nodeData := range mapOfAggregatedScores {

		if err := task.scoreStore.UpsertAggregatedScore(nodeData); err != nil {
			logger.WithField("node_id", nodeID).WithError(err).Error("error updating aggregated score")
			continue
		}
	}

	return nil
}

func AggregateChallengeScore(data types.AccumulativeChallengeData) float64 {
	var challengerScore, recipientScore, observerScore float64

	if data.TotalChallengesAsChallengers > 0 {
		challengerScore = (float64(data.CorrectChallengerEvaluations) / float64(data.TotalChallengesAsChallengers)) * 40
	}

	if data.TotalChallengesAsRecipients > 0 {
		recipientScore = (float64(data.CorrectRecipientEvaluations) / float64(data.TotalChallengesAsRecipients)) * 40
	}

	if data.TotalChallengesAsObservers > 0 {
		observerScore = (float64(data.CorrectObserverEvaluations) / float64(data.TotalChallengesAsObservers)) * 20
	}

	return challengerScore + recipientScore + observerScore
}
