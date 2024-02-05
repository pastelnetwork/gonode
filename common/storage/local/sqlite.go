package local

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	json "github.com/json-iterator/go"
	_ "github.com/mattn/go-sqlite3" //go-sqlite3
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/types"
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

const createSelfHealingExecutionMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_execution_metrics_unique ON self_healing_execution_metrics(trigger_id, challenge_id, message_type);
`

const alterTablePingHistory = `ALTER TABLE ping_history
ADD COLUMN metrics_last_broadcast_at DATETIME NULL;`

const alterTablePingHistoryGenerationMetrics = `ALTER TABLE ping_history
ADD COLUMN generation_metrics_last_broadcast_at DATETIME NULL;`

const alterTablePingHistoryExecutionMetrics = `ALTER TABLE ping_history
ADD COLUMN execution_metrics_last_broadcast_at DATETIME NULL;`

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

// InsertTaskHistory inserts task history
func (s *SQLiteStore) InsertTaskHistory(history types.TaskHistory) (hID int, err error) {
	var stringifyDetails string
	if history.Details != nil {
		stringifyDetails = history.Details.Stringify()
	}

	const insertQuery = "INSERT INTO task_history(id, time, task_id, status, details) VALUES(NULL,?,?,?,?);"
	res, err := s.db.Exec(insertQuery, history.CreatedAt, history.TaskID, history.Status, stringifyDetails)

	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
}

// QueryTaskHistory gets task history by taskID
func (s *SQLiteStore) QueryTaskHistory(taskID string) (history []types.TaskHistory, err error) {
	const selectQuery = "SELECT * FROM task_history WHERE task_id = ? LIMIT 100"
	rows, err := s.db.Query(selectQuery, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []types.TaskHistory
	for rows.Next() {
		i := types.TaskHistory{}
		var details string
		err = rows.Scan(&i.ID, &i.CreatedAt, &i.TaskID, &i.Status, &details)
		if err != nil {
			return nil, err
		}

		if details != emptyString {
			err = json.Unmarshal([]byte(details), &i.Details)
			if err != nil {
				log.Info(details)
				log.WithError(err).Error(fmt.Sprintf("cannot unmarshal task history details: %s", details))
				i.Details = nil
			}
		}

		data = append(data, i)
	}

	return data, nil
}

// InsertStorageChallengeMessage inserts failed storage challenge to db
func (s *SQLiteStore) InsertStorageChallengeMessage(challenge types.StorageChallengeLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO storage_challenge_messages(id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MessageType, challenge.Data, challenge.Sender, challenge.SenderSignature, now, now)
	if err != nil {
		return err
	}

	return nil
}

// UpsertPingHistory inserts/update ping information into the ping_history table
func (s *SQLiteStore) UpsertPingHistory(pingInfo types.PingInfo) error {
	now := time.Now().UTC()

	const upsertQuery = `
		INSERT INTO ping_history (
			supernode_id, ip_address, total_pings, total_successful_pings, 
			avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(supernode_id, ip_address) 
		DO UPDATE SET
			total_pings = excluded.total_pings,
			total_successful_pings = excluded.total_successful_pings,
			avg_ping_response_time = excluded.avg_ping_response_time,
			is_online = excluded.is_online,
			is_on_watchlist = excluded.is_on_watchlist,
			is_adjusted = excluded.is_adjusted,
		    last_seen = excluded.last_seen,
		    cumulative_response_time = excluded.cumulative_response_time,
			updated_at = excluded.updated_at;`

	_, err := s.db.Exec(upsertQuery,
		pingInfo.SupernodeID, pingInfo.IPAddress, pingInfo.TotalPings,
		pingInfo.TotalSuccessfulPings, pingInfo.AvgPingResponseTime,
		pingInfo.IsOnline, pingInfo.IsOnWatchlist, pingInfo.IsAdjusted, pingInfo.LastSeen.Time, pingInfo.CumulativeResponseTime, now, now)
	if err != nil {
		return err
	}

	return nil
}

// GetPingInfoBySupernodeID retrieves a ping history record by supernode ID
func (s *SQLiteStore) GetPingInfoBySupernodeID(supernodeID string) (*types.PingInfo, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        WHERE supernode_id = ?;`

	var pingInfo types.PingInfo
	row := s.db.QueryRow(selectQuery, supernodeID)

	// Scan the row into the PingInfo struct
	err := row.Scan(
		&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
		&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
		&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
		&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &pingInfo, nil
}

// GetWatchlistPingInfo retrieves all the nodes that are on watchlist
func (s *SQLiteStore) GetWatchlistPingInfo() ([]types.PingInfo, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        WHERE is_on_watchlist = true AND is_adjusted = false;`

	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {
		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// UpdatePingInfo updates the ping info
func (s *SQLiteStore) UpdatePingInfo(supernodeID string, isOnWatchlist, isAdjusted bool) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET is_adjusted = ?, is_on_watchlist = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, isAdjusted, isOnWatchlist, supernodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMetricsBroadcastTimestamp updates the ping info metrics_last_broadcast_at
func (s *SQLiteStore) UpdateMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateGenerationMetricsBroadcastTimestamp updates the ping info generation_metrics_last_broadcast_at
func (s *SQLiteStore) UpdateGenerationMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET generation_metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().Add(-60*time.Minute).UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateExecutionMetricsBroadcastTimestamp updates the ping info execution_metrics_last_broadcast_at
func (s *SQLiteStore) UpdateExecutionMetricsBroadcastTimestamp(nodeID string) error {
	// Update query
	const updateQuery = `
UPDATE ping_history
SET execution_metrics_last_broadcast_at = ?
WHERE supernode_id = ?;`

	// Execute the update query
	_, err := s.db.Exec(updateQuery, time.Now().Add(-60*time.Minute).UTC(), nodeID)
	if err != nil {
		return err
	}

	return nil
}

// InsertBroadcastMessage inserts broadcast storage challenge msg to db
func (s *SQLiteStore) InsertBroadcastMessage(challenge types.BroadcastLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO broadcast_challenge_messages(id, challenge_id, data, challenger, recipient, observers, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.Data, challenge.Challenger, challenge.Recipient, challenge.Observers, now, now)
	if err != nil {
		return err
	}

	return nil
}

// QueryStorageChallengeMessage retrieves storage challenge message against challengeID and messageType
func (s *SQLiteStore) QueryStorageChallengeMessage(challengeID string, messageType int) (challengeMessage types.StorageChallengeLogMessage, err error) {
	const selectQuery = "SELECT * FROM storage_challenge_messages WHERE challenge_id=? AND message_type=?"
	err = s.db.QueryRow(selectQuery, challengeID, messageType).Scan(
		&challengeMessage.ID, &challengeMessage.ChallengeID, &challengeMessage.MessageType, &challengeMessage.Data,
		&challengeMessage.Sender, &challengeMessage.SenderSignature, &challengeMessage.CreatedAt, &challengeMessage.UpdatedAt)

	if err != nil {
		return challengeMessage, err
	}

	return challengeMessage, nil
}

// CleanupStorageChallenges cleans up challenges stored in DB for self-healing
func (s *SQLiteStore) CleanupStorageChallenges() (err error) {
	const delQuery = "DELETE FROM storage_challenge_messages"
	_, err = s.db.Exec(delQuery)
	return err
}

// QuerySelfHealingChallenges retrieves self-healing audit logs stored in DB for self-healing
func (s *SQLiteStore) QuerySelfHealingChallenges() (challenges []types.SelfHealingChallenge, err error) {
	const selectQuery = "SELECT * FROM self_healing_challenges"
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		challenge := types.SelfHealingChallenge{}
		err = rows.Scan(&challenge.ID, &challenge.ChallengeID, &challenge.MerkleRoot, &challenge.FileHash,
			&challenge.ChallengingNode, &challenge.RespondingNode, &challenge.VerifyingNode, &challenge.ReconstructedFileHash,
			&challenge.Status, &challenge.CreatedAt, &challenge.UpdatedAt)
		if err != nil {
			return nil, err
		}

		challenges = append(challenges, challenge)
	}

	return challenges, nil
}

// GetAllPingInfos retrieves all ping infos
func (s *SQLiteStore) GetAllPingInfos() (types.PingInfos, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time,
               created_at, updated_at
        FROM ping_history
        `
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {

		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// GetAllPingInfoForOnlineNodes retrieves all ping infos for nodes that are online
func (s *SQLiteStore) GetAllPingInfoForOnlineNodes() (types.PingInfos, error) {
	const selectQuery = `
        SELECT id, supernode_id, ip_address, total_pings, total_successful_pings,
               avg_ping_response_time, is_online, is_on_watchlist, is_adjusted, last_seen, cumulative_response_time, 
               metrics_last_broadcast_at, generation_metrics_last_broadcast_at, execution_metrics_last_broadcast_at,
               created_at, updated_at
        FROM ping_history
        WHERE is_online = true

        `
	rows, err := s.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pingInfos types.PingInfos
	for rows.Next() {

		var pingInfo types.PingInfo
		if err := rows.Scan(
			&pingInfo.ID, &pingInfo.SupernodeID, &pingInfo.IPAddress, &pingInfo.TotalPings,
			&pingInfo.TotalSuccessfulPings, &pingInfo.AvgPingResponseTime,
			&pingInfo.IsOnline, &pingInfo.IsOnWatchlist, &pingInfo.IsAdjusted, &pingInfo.LastSeen, &pingInfo.CumulativeResponseTime,
			&pingInfo.MetricsLastBroadcastAt, &pingInfo.GenerationMetricsLastBroadcastAt, &pingInfo.ExecutionMetricsLastBroadcastAt,
			&pingInfo.CreatedAt, &pingInfo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pingInfos = append(pingInfos, pingInfo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pingInfos, nil
}

// CleanupSelfHealingChallenges cleans up self-healing challenges stored in DB for inspection
func (s *SQLiteStore) CleanupSelfHealingChallenges() (err error) {
	const delQuery = "DELETE FROM self_healing_challenges"
	_, err = s.db.Exec(delQuery)
	return err
}

// OpenHistoryDB opens history DB
func OpenHistoryDB() (storage.LocalStoreInterface, error) {
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

	_, _ = db.Exec(alterTaskHistory)

	_, _ = db.Exec(alterTablePingHistory)

	_, _ = db.Exec(alterTablePingHistoryGenerationMetrics)

	_, _ = db.Exec(alterTablePingHistoryExecutionMetrics)

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
