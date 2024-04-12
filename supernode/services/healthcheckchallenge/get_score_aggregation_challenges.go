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
	logger := log.WithContext(ctx).WithField("method", "GetScoreAggregationHCChallenges")

	logger.Info("invoked")

	tracker, err := task.historyDB.GetScoreLastAggregatedAt(queries.HealthCheckChallengeScoreAggregationType)
	if err != nil {
		logger.WithError(err).Error("error retrieving score-aggregate tracker for healthcheck-challenges")
		return errors.Errorf("error retrieving score-aggregate tracker for healthcheck-challenges")
	}

	totalChallengesToBeAggregated, err := task.getChallengeIDsCountForScoreAggregation(ctx, tracker)
	if err != nil {
		logger.WithError(err).Error("error retrieving hc affirmations counts")
		return errors.Errorf("error retrieving hc affirmation messages count")
	}

	if totalChallengesToBeAggregated == 0 {
		logger.Info("no healthcheck-challenges found for score aggregation")
		return nil
	}

	totalBatches := task.calculateTotalBatches(totalChallengesToBeAggregated)

	for i := 0; i < totalBatches; i++ {
		batchOfChallengeIDs, err := task.retrieveChallengeIDsInBatches(ctx, tracker, i)
		if err != nil {
			logger.WithError(err).Error("error retrieving the healthcheck-challenge ids for score aggregation")
			return err
		}

		if len(batchOfChallengeIDs) == 0 {
			continue
		}

		err = task.BatchInsertChallengeIDs(ctx, batchOfChallengeIDs, false)
		if err != nil {
			logger.Error("error storing healthcheck challenge_ids for score aggregation")
			return err
		}
	}

	if err := task.historyDB.UpsertScoreLastAggregatedAt(queries.HealthCheckChallengeScoreAggregationType); err != nil {
		logger.WithError(err).Error("error storing affirmation type batch in score aggregation healthcheck-challenges")
		return errors.Errorf("error updating aggregated til for hc")
	}

	return nil
}

func (task *HCTask) retrieveChallengeIDsInBatches(ctx context.Context, tracker queries.ScoreAggregationTracker, batchNumber int) ([]string, error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return nil, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	var (
		zeroTime     time.Time
		challengeIDs []string
		before       time.Time
	)

	before = time.Now().UTC().Add(-6 * time.Hour)
	if !tracker.AggregatedTil.Valid {
		challengeIDs, err = store.GetDistinctHCChallengeIDs(zeroTime, before, batchNumber)
		if err != nil {
			return nil, err
		}
	} else {
		challengeIDs, err = store.GetDistinctHCChallengeIDs(tracker.AggregatedTil.Time, before, batchNumber)
		if err != nil {
			return nil, err
		}
	}

	return challengeIDs, nil
}

func (task *HCTask) getChallengeIDsCountForScoreAggregation(ctx context.Context, tracker queries.ScoreAggregationTracker) (int, error) {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return 0, err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	var (
		zeroTime          time.Time
		challenegIDsCount int
		before            time.Time
	)

	before = time.Now().UTC().Add(-6 * time.Hour)
	if !tracker.AggregatedTil.Valid {
		challenegIDsCount, err = store.GetDistinctHCChallengeIDsCountForScoreAggregation(zeroTime, before)
		if err != nil {
			return 0, err
		}
	} else {
		challenegIDsCount, err = store.GetDistinctHCChallengeIDsCountForScoreAggregation(tracker.AggregatedTil.Time, before)
		if err != nil {
			return 0, err
		}
	}

	return challenegIDsCount, nil
}

func (task *HCTask) calculateTotalBatches(count int) int {
	return int(math.Ceil(float64(count) / float64(batchSize)))
}

// BatchInsertChallengeIDs stores the challenge id to db for further verification
func (task *HCTask) BatchInsertChallengeIDs(ctx context.Context, batchOfChallengeIDs []string, isAggregated bool) error {
	store, err := queries.OpenHistoryDB()
	if err != nil {
		return err
	}
	if store != nil {
		defer store.CloseHistoryDB(ctx)
	}

	if store != nil {
		err = store.BatchInsertHCScoreAggregationChallenges(batchOfChallengeIDs, isAggregated)
		if err != nil {
			return err
		}
	}

	return nil
}
