package healthcheckchallenge

import (
	"context"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"math"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	batchSize = 500
)

// GetScoreAggregationChallenges gets the challenges and store it for score aggregation
func (task *HCTask) GetScoreAggregationChallenges(ctx context.Context) error {
	logger := log.WithContext(ctx).WithField("method", "GetScoreAggregationChallenges")

	logger.Info("invoked")

	tracker, err := task.historyDB.GetScoreLastAggregatedAt(queries.HealthCheckChallengeScoreAggregationType)
	if err != nil {
		logger.WithError(err).Error("error retrieving healthcheck-challenge score aggregate tracker")
		return errors.Errorf("error score-aggregate tracker for healthcheck-challenges")
	}

	totalChallengesToBeAggregated, err := task.getChallengeIDsCountForScoreAggregation(tracker)
	if err != nil {
		logger.WithError(err).Error("error retrieving hc affirmations counts")
		return errors.Errorf("error retrieving hc affirmation messages count")
	}

	if totalChallengesToBeAggregated == 0 {
		logger.Info("no healthcheck-challenges found for score aggregation")
	}

	totalBatches := task.calculateTotalBatches(totalChallengesToBeAggregated)

	for i := 0; i < totalBatches; i++ {
		batchOfChallengeIDs, err := task.retrieveChallengeIDsInBatches(tracker, i)
		if err != nil {
			logger.Error("error retrieving the healthcheck-challenge ids for score aggregation")
		}

		err = task.historyDB.BatchInsertHCScoreAggregationChallenges(batchOfChallengeIDs, false)
		if err != nil {
			logger.Error("error storing healthcheck challenge_ids for score aggregation")
			return err
		}
	}

	if err := task.historyDB.UpsertScoreLastAggregatedAt(queries.HealthCheckChallengeScoreAggregationType); err != nil {
		logger.Error("error storing affirmation type batch in score aggregation healthcheck-challenges")
		return errors.Errorf("error updating aggregated til")
	}

	return nil
}

func (task *HCTask) retrieveChallengeIDsInBatches(tracker queries.ScoreAggregationTracker, batchNumber int) ([]string, error) {
	var (
		zeroTime     time.Time
		err          error
		challengeIDs []string
		before       time.Time
	)

	before = time.Now().UTC().Add(-6 * time.Hour)
	if !tracker.AggregatedTil.Valid {
		challengeIDs, err = task.historyDB.GetDistinctHCChallengeIDs(zeroTime, before, batchNumber)
		if err != nil {
			return nil, err
		}
	} else {
		challengeIDs, err = task.historyDB.GetDistinctHCChallengeIDs(tracker.AggregatedTil.Time, before, batchNumber)
		if err != nil {
			return nil, err
		}
	}

	return challengeIDs, nil
}

func (task *HCTask) getChallengeIDsCountForScoreAggregation(tracker queries.ScoreAggregationTracker) (int, error) {
	var (
		zeroTime          time.Time
		err               error
		challenegIDsCount int
		before            time.Time
	)

	before = time.Now().UTC().Add(-6 * time.Hour)
	if !tracker.AggregatedTil.Valid {
		challenegIDsCount, err = task.historyDB.GetDistinctHCChallengeIDsCountForScoreAggregation(zeroTime, before)
		if err != nil {
			return 0, err
		}
	} else {
		challenegIDsCount, err = task.historyDB.GetDistinctHCChallengeIDsCountForScoreAggregation(tracker.AggregatedTil.Time, before)
		if err != nil {
			return 0, err
		}
	}

	return challenegIDsCount, nil
}

func (task *HCTask) calculateTotalBatches(count int) int {
	return int(math.Ceil(float64(count) / float64(batchSize)))
}
