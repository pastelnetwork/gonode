package queries

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
const createTaskHistory string = `
  CREATE TABLE IF NOT EXISTS task_history (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  time DATETIME NOT NULL,
  task_id TEXT NOT NULL,
  status TEXT NOT NULL
  );`

const alterTaskHistory string = `ALTER TABLE task_history ADD COLUMN details TEXT;`

const createStorageChallengeMessages string = `
  CREATE TABLE IF NOT EXISTS storage_challenge_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  sender_signature BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createBroadcastChallengeMessages string = `
  CREATE TABLE IF NOT EXISTS broadcast_challenge_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  challenger TEXT NOT NULL,
  recipient TEXT NOT NULL,
  observers TEXT NOT NULL,
  data BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createStorageChallengeMessagesUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS storage_challenge_messages_unique ON storage_challenge_messages(challenge_id, message_type, sender_id);
`

const createSelfHealingChallenges string = `
  CREATE TABLE IF NOT EXISTS self_healing_challenges (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  merkleroot TEXT NOT NULL,
  file_hash TEXT NOT NULL,
  challenging_node TEXT NOT NULL,
  responding_node TEXT NOT NULL,
  verifying_node TEXT,
  reconstructed_file_hash BLOB,
  status TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL                                                  
  );`

const createPingHistory string = `
  CREATE TABLE IF NOT EXISTS ping_history (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  supernode_id TEXT UNIQUE NOT NULL,
  ip_address TEXT UNIQUE NOT NULL,
  total_pings INTEGER NOT NULL,
  total_successful_pings INTEGER NOT NULL,
  avg_ping_response_time FLOAT NOT NULL,
  is_online BOOLEAN NOT NULL,
  is_on_watchlist BOOLEAN NOT NULL,
  is_adjusted BOOLEAN NOT NULL,
  cumulative_response_time FLOAT NOT NULL,
  last_seen DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL                                                  
  );`

const createPingHistoryUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS ping_history_unique ON ping_history(supernode_id, ip_address);
`

const createSelfHealingGenerationMetrics string = `
  CREATE TABLE IF NOT EXISTS self_healing_generation_metrics (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  trigger_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  sender_signature BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createSelfHealingGenerationMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_generation_metrics_unique ON self_healing_generation_metrics(trigger_id);
`

const createSelfHealingExecutionMetrics string = `
  CREATE TABLE IF NOT EXISTS self_healing_execution_metrics (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  trigger_id TEXT NOT NULL,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  sender_signature BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createSelfHealingChallengeTickets string = `
  CREATE TABLE IF NOT EXISTS self_healing_challenge_events (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  trigger_id TEXT NOT NULL,
  ticket_id TEXT NOT NULL,
  challenge_id TEXT NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  is_processed BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
`

const createSelfHealingChallengeTicketsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_challenge_events_unique ON self_healing_challenge_events(trigger_id, ticket_id, challenge_id);
`

const createSelfHealingExecutionMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_execution_metrics_unique ON self_healing_execution_metrics(trigger_id, challenge_id, message_type);
`

const alterTablePingHistory = `ALTER TABLE ping_history
ADD COLUMN metrics_last_broadcast_at DATETIME NULL;`

const alterTablePingHistoryGenerationMetrics = `ALTER TABLE ping_history
ADD COLUMN generation_metrics_last_broadcast_at DATETIME NULL;`

const alterTablePingHistoryExecutionMetrics = `ALTER TABLE ping_history
ADD COLUMN execution_metrics_last_broadcast_at DATETIME NULL;`

const createStorageChallengeMetrics string = `
  CREATE TABLE IF NOT EXISTS storage_challenge_metrics (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createStorageChallengeMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS storage_challenge_metrics_unique ON storage_challenge_metrics(challenge_id, message_type, sender_id);
`

const createHealthCheckChallengeMessages string = `
  CREATE TABLE IF NOT EXISTS healthcheck_challenge_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  sender_signature BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createBroadcastHealthCheckChallengeMessages string = `
  CREATE TABLE IF NOT EXISTS broadcast_healthcheck_challenge_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  challenger TEXT NOT NULL,
  recipient TEXT NOT NULL,
  observers TEXT NOT NULL,
  data BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createHealthCheckChallengeMetrics string = `
  CREATE TABLE IF NOT EXISTS healthcheck_challenge_metrics (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
`
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

const createHealthCheckChallengeMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS healthcheck_challenge_metrics_unique ON healthcheck_challenge_metrics(challenge_id, message_type, sender_id);
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

const alterTablePingHistoryHealthCheckColumn = `ALTER TABLE ping_history
ADD COLUMN health_check_metrics_last_broadcast_at DATETIME NULL;`

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
	historyDBName = "history.db"
	emptyString   = ""
)

// SQLiteStore handles sqlite ops
type SQLiteStore struct {
	db *sqlx.DB
}

// CloseHistoryDB closes history database
func (s *SQLiteStore) CloseHistoryDB(ctx context.Context) {
	if err := s.db.Close(); err != nil {
		log.WithContext(ctx).WithError(err).Error("error closing history db")
	}
}

// OpenHistoryDB opens history DB
func OpenHistoryDB() (LocalStoreInterface, error) {
	dbFile := filepath.Join(configurer.DefaultPath(), historyDBName)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}

	if _, err := db.Exec(createTaskHistory); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createStorageChallengeMessages); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createStorageChallengeMessagesUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot execute migration: %w", err)
	}

	if _, err := db.Exec(createBroadcastChallengeMessages); err != nil {
		return nil, fmt.Errorf("cannot execute migration: %w", err)
	}

	if _, err := db.Exec(createSelfHealingChallenges); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createPingHistory); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createPingHistoryUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingGenerationMetrics); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingGenerationMetricsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingExecutionMetrics); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingExecutionMetricsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingChallengeTickets); err != nil {
		return nil, fmt.Errorf("cannot create createSelfHealingChallengeTickets: %w", err)
	}

	if _, err := db.Exec(createSelfHealingChallengeTicketsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create createSelfHealingChallengeTicketsUniqueIndex: %w", err)
	}

	if _, err := db.Exec(createStorageChallengeMetrics); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createStorageChallengeMetricsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createHealthCheckChallengeMessages); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createHealthCheckChallengeMetrics); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createHealthCheckChallengeMetricsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createBroadcastHealthCheckChallengeMessages); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
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

	_, _ = db.Exec(alterTaskHistory)

	_, _ = db.Exec(alterTablePingHistory)

	_, _ = db.Exec(alterTablePingHistoryGenerationMetrics)

	_, _ = db.Exec(alterTablePingHistoryExecutionMetrics)

	_, _ = db.Exec(alterTablePingHistoryHealthCheckColumn)

	pragmas := []string{
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA cache_size=-262144;",
		"PRAGMA busy_timeout=120000;",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return nil, fmt.Errorf("cannot set sqlite database parameter: %w", err)
		}
	}

	return &SQLiteStore{
		db: db,
	}, nil
}
