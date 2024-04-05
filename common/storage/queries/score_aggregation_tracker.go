package queries

import (
	"database/sql"
	"github.com/pastelnetwork/gonode/common/errors"
	"time"
)

type ScoreAggregationTrackerChallengeType int

const (
	StorageChallengeScoreAggregationType ScoreAggregationTrackerChallengeType = iota + 1
	HealthCheckChallengeScoreAggregationType
)

type ScoreAggregationTracker struct {
	ChallengeType ScoreAggregationTrackerChallengeType `db:"challenge_type"`
	AggregatedTil sql.NullTime                         `db:"aggregated_til"`
}

type ScoreAggregationTrackerQueries interface {
	GetScoreLastAggregatedAt(challengeType ScoreAggregationTrackerChallengeType) (ScoreAggregationTracker, error)
	UpsertScoreLastAggregatedAt(challengeType ScoreAggregationTrackerChallengeType) error
}

// GetScoreLastAggregatedAt retrieves the last aggregated for the given challenge-type
func (s *SQLiteStore) GetScoreLastAggregatedAt(challengeType ScoreAggregationTrackerChallengeType) (ScoreAggregationTracker, error) {
	const query = `SELECT challenge_type, aggregated_til FROM score_aggregation_tracker WHERE challenge_type = ?`

	var tracker ScoreAggregationTracker

	row := s.db.QueryRow(query, challengeType)

	err := row.Scan(&tracker.ChallengeType, &tracker.AggregatedTil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ScoreAggregationTracker{
				ChallengeType: challengeType,
				AggregatedTil: sql.NullTime{Valid: false},
			}, nil
		}

		return tracker, err
	}

	return tracker, nil
}

// UpsertScoreLastAggregatedAt updates or insert score-aggregation-tracker
func (s *SQLiteStore) UpsertScoreLastAggregatedAt(challengeType ScoreAggregationTrackerChallengeType) error {
	const upsertQuery = `
INSERT INTO score_aggregation_tracker (challenge_type, aggregated_til)
VALUES (?, ?)
ON CONFLICT(challenge_type)
DO UPDATE SET aggregated_til = excluded.aggregated_til;`

	_, err := s.db.Exec(upsertQuery, challengeType, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}
