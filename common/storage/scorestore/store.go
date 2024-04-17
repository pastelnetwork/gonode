package scorestore

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/log"
)

const minVerifications = 3

const createAccumulativeSCData string = `
CREATE TABLE IF NOT EXISTS accumulative_sc_data (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    node_id TEXT NOT NULL,
    ip_address TEXT,
    total_challenges_as_challengers INTEGER,
    total_challenges_as_recipients INTEGER,
    total_challenges_as_observers INTEGER,
    correct_challenger_evaluations INTEGER,
    correct_recipient_evaluations INTEGER,
    correct_observer_evaluation INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createAccumulativeHCData string = `
CREATE TABLE IF NOT EXISTS accumulative_hc_data (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    node_id TEXT NOT NULL,
    ip_address TEXT,
    total_challenges_as_challengers INTEGER,
    total_challenges_as_recipients INTEGER,
    total_challenges_as_observers INTEGER,
    correct_challenger_evaluations INTEGER,
    correct_recipient_evaluations INTEGER,
    correct_observer_evaluation INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createSCScoreAggregationQueue string = `
CREATE TABLE IF NOT EXISTS sc_score_aggregation_queue (
    challenge_id TEXT PRIMARY KEY NOT NULL,
    is_aggregated BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createScoreAggregationTracker string = `
CREATE TABLE IF NOT EXISTS score_aggregation_tracker (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    challenge_type INTEGER NOT NULL,
    aggregated_til DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createScoreAggregationTrackerUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS score_aggregation_tracker_unique ON score_aggregation_tracker(challenge_type);
`

const createAccumulativeSCDataUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS accumulative_sc_data_unique_index 
ON accumulative_sc_data(node_id);
`

const createAccumulativeHCDataUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS accumulative_hc_data_unique_index 
ON accumulative_hc_data(node_id);
`

const createAggregatedSCChallengesUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS sc_score_aggregation_queue_unique_index 
ON sc_score_aggregation_queue(challenge_id);
`

const createHCScoreAggregationQueue string = `
CREATE TABLE IF NOT EXISTS hc_score_aggregation_queue (
    challenge_id TEXT PRIMARY KEY NOT NULL,
    is_aggregated BOOLEAN NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createAggregatedHCChallengesUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS hc_score_aggregation_queue_unique_index 
ON hc_score_aggregation_queue(challenge_id);
`

const createAggregatedChallengeScores string = `
CREATE TABLE IF NOT EXISTS aggregated_challenge_scores (
    node_id TEXT PRIMARY KEY NOT NULL,
    ip_address TEXT,
    storage_challenge_score REAL,
    healthcheck_challenge_score REAL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
`

const createAggregatedChallengeScoresUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS aggregated_challenge_scores_unique_index 
ON aggregated_challenge_scores(node_id);
`

const (
	scoreDBName = "score.db"
	emptyString = ""
)

// ScoringStore handles sqlite ops for scoring
type ScoringStore struct {
	db *sqlx.DB
}

// CloseDB closes score database
func (s *ScoringStore) CloseDB(ctx context.Context) {
	if err := s.db.Close(); err != nil {
		log.WithContext(ctx).WithError(err).Error("error closing score db")
	}
}

// OpenScoringDb opens scoring DB
func OpenScoringDb() (ScoreStorageInterface, error) {
	dbFile := filepath.Join(configurer.DefaultPath(), scoreDBName)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}

	if _, err := db.Exec(createAccumulativeSCData); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSCScoreAggregationQueue); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createScoreAggregationTracker); err != nil {
		return nil, fmt.Errorf("cannot create table: %w", err)
	}

	if _, err := db.Exec(createAccumulativeSCDataUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAggregatedSCChallengesUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createScoreAggregationTrackerUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createHCScoreAggregationQueue); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAggregatedHCChallengesUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAccumulativeHCData); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAccumulativeHCDataUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAggregatedChallengeScores); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	if _, err := db.Exec(createAggregatedChallengeScoresUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create unique index: %w", err)
	}

	pragmas := []string{
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-262144;",
		"PRAGMA busy_timeout=120000;",
		"PRAGMA journal_mode=WAL;",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("cannot set sqlite database parameter: %w", err)
		}
	}

	return &ScoringStore{
		db: db,
	}, nil
}
