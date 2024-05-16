package storagechallenge

import (
	"context"
	"math"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage/queries"
	"github.com/pastelnetwork/gonode/common/storage/scorestore"

	"github.com/pastelnetwork/gonode/common/errors"
)

const (
	batchSize = 500
)

// GetScoreAggregationChallenges gets the challenges and store it for score aggregation
func (task *SCTask) GetScoreAggregationChallenges(ctx context.Context) error {
	logger := log.WithContext(ctx).WithField("method", "GetScoreAggregationChallenges")

	logger.Info("invoked")

	tracker, err := task.getStorageChallengeTracker(ctx)
	if err != nil {
		logger.WithError(err).Error("error retrieving storage-challenge score aggregate tracker")
		return errors.Errorf("error score-aggregate tracker for storage challenges")
	}

	totalChallengesToBeAggregated, err := task.getChallengeIDsCountForScoreAggregation(ctx, tracker)
	if err != nil {
		logger.WithError(err).Error("error retrieving sc affirmations counts")
		return errors.Errorf("error retrieving affirmation messages count")
	}

	if totalChallengesToBeAggregated == 0 {
		logger.Info("no storage-challenges found for score aggregation")
		return nil
	}

	totalBatches := task.calculateTotalBatches(totalChallengesToBeAggregated)

	store, err := scorestore.OpenScoringDb()
	if err != nil {
		return err
	}
	if store != nil {
		defer store.CloseDB(ctx)
	}

	for i := 0; i < totalBatches; i++ {
		batchOfChallengeIDs, err := task.retrieveChallengeIDsInBatches(ctx, tracker, i)
		if err != nil {
			logger.WithError(err).Error("error retrieving the challenge ids for score aggregation")
			return err
		}

		if len(batchOfChallengeIDs) == 0 {
			continue
		}

		if store != nil {
			err = store.BatchInsertScoreAggregationChallenges(batchOfChallengeIDs, false)
			if err != nil {
				logger.WithError(err).Error("error storing storage_challenge_ids for score aggregation")
				return err
			}
		}
	}

	if err := task.ScoreStore.UpsertScoreLastAggregatedAt(scorestore.StorageChallengeScoreAggregationType); err != nil {
		logger.WithError(err).Error("error updating aggregated til for storage-challenge score aggregation tracker")
		return errors.Errorf("error updating aggregated til")
	}

	return nil
}

func (task *SCTask) getStorageChallengeTracker(ctx context.Context) (scorestore.ScoreAggregationTracker, error) {
	tracker, err := task.ScoreStore.GetScoreLastAggregatedAt(scorestore.StorageChallengeScoreAggregationType)
	if err != nil {
		return scorestore.ScoreAggregationTracker{}, err
	}

	return tracker, nil
}

func (task *SCTask) retrieveChallengeIDsInBatches(ctx context.Context, tracker scorestore.ScoreAggregationTracker, batchNumber int) ([]string, error) {
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
		challengeIDs, err = store.GetDistinctChallengeIDs(zeroTime, before, batchNumber)
		if err != nil {
			return nil, err
		}
	} else {
		challengeIDs, err = store.GetDistinctChallengeIDs(tracker.AggregatedTil.Time, before, batchNumber)
		if err != nil {
			return nil, err
		}
	}

	return challengeIDs, nil
}

func (task *SCTask) getChallengeIDsCountForScoreAggregation(ctx context.Context, tracker scorestore.ScoreAggregationTracker) (int, error) {
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
		challenegIDsCount, err = store.GetDistinctChallengeIDsCountForScoreAggregation(zeroTime, before)
		if err != nil {
			return 0, err
		}
	} else {
		challenegIDsCount, err = store.GetDistinctChallengeIDsCountForScoreAggregation(tracker.AggregatedTil.Time, before)
		if err != nil {
			return 0, err
		}
	}

	return challenegIDsCount, nil
}

func (task *SCTask) calculateTotalBatches(count int) int {
	return int(math.Ceil(float64(count) / float64(batchSize)))
}
