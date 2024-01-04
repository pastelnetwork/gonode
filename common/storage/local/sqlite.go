package local

import (
	"context"
	"database/sql"
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

const createSelfHealingMessages string = `
  CREATE TABLE IF NOT EXISTS self_healing_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  message_type INTEGER NOT NULL,
  data BLOB NOT NULL,
  sender_id TEXT NOT NULL,
  sender_signature BLOB NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);`

const createSelfHealingMessagesUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_messages_unique ON self_healing_messages(challenge_id, message_type, sender_id);
`

const createSelfHealingMetricsTable string = `
	CREATE TABLE IF NOT EXISTS self_healing_metrics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	challenge_id TEXT NOT NULL,
	sent_tickets_for_self_healing INTEGER,
	estimated_missing_keys INTEGER,
	tickets_in_progress INTEGER,
	tickets_required_self_healing INTEGER,
	successfully_self_healed_tickets INTEGER,
	successfully_verified_tickets INTEGER,
	created_at DATETIME,
	updated_at DATETIME
);
`

const createSelfHealingMetricsUniqueIndex string = `
CREATE UNIQUE INDEX IF NOT EXISTS self_healing_metrics_unique ON self_healing_metrics(challenge_id);
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

// InsertSelfHealingChallenge inserts self-healing challenge
func (s *SQLiteStore) InsertSelfHealingChallenge(challenge types.SelfHealingLogMessage) error {
	now := time.Now().UTC()
	const insertQuery = "INSERT INTO self_healing_messages(id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,$1,$2,$3,$4,$5,$6,$7);"

	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MessageType, challenge.Data, challenge.SenderID, challenge.SenderSignature, now, now)
	if err != nil {
		return err
	}

	return nil
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

// InsertSelfHealingMetrics inserts the self-healing metrics against challenge
func (s *SQLiteStore) InsertSelfHealingMetrics(metrics types.SelfHealingMetrics) error {
	now := time.Now().UTC()
	const insertQuery = `
        INSERT INTO self_healing_metrics(
            challenge_id,
            sent_tickets_for_self_healing,
            estimated_missing_keys,
            tickets_in_progress,
            tickets_required_self_healing,
            successfully_self_healed_tickets,
            successfully_verified_tickets,
            created_at,
            updated_at
        )
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT(challenge_id) 
		DO UPDATE SET
			sent_tickets_for_self_healing = excluded.sent_tickets_for_self_healing,
			estimated_missing_keys = excluded.estimated_missing_keys,
			tickets_in_progress = excluded.tickets_in_progress,
			tickets_required_self_healing = excluded.tickets_required_self_healing,
			successfully_self_healed_tickets = excluded.successfully_self_healed_tickets,
			successfully_verified_tickets = excluded.successfully_verified_tickets,
			updated_at = excluded.updated_at;
    `

	_, err := s.db.Exec(
		insertQuery,
		metrics.ChallengeID,
		metrics.SentTicketsForSelfHealing,
		metrics.EstimatedMissingKeys,
		metrics.TicketsInProgress,
		metrics.TicketsRequiredSelfHealing,
		metrics.SuccessfullySelfHealedTickets,
		metrics.SuccessfullyVerifiedTickets,
		now,
		now,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSelfHealingMetrics returns the sum of self-healing metrics of all challenges stored in DB
func (s *SQLiteStore) GetSelfHealingMetrics() (types.SelfHealingMetrics, error) {
	query := `
    SELECT
        SUM(sent_tickets_for_self_healing) AS sent_tickets_for_self_healing,
        SUM(estimated_missing_keys) AS estimated_missing_keys,
        SUM(tickets_in_progress) AS tickets_in_progress,
        SUM(tickets_required_self_healing) AS tickets_required_self_healing,
        SUM(successfully_self_healed_tickets) AS successfully_self_healed_tickets,
        SUM(successfully_verified_tickets) AS successfully_verified_tickets
    FROM self_healing_metrics;
    `

	var metrics types.SelfHealingMetrics
	var nullMetrics struct {
		SentTicketsForSelfHealing     sql.NullInt64
		EstimatedMissingKeys          sql.NullInt64
		TicketsInProgress             sql.NullInt64
		TicketsRequiredSelfHealing    sql.NullInt64
		SuccessfullySelfHealedTickets sql.NullInt64
		SuccessfullyVerifiedTickets   sql.NullInt64
	}

	err := s.db.QueryRow(query).Scan(
		&nullMetrics.SentTicketsForSelfHealing,
		&nullMetrics.EstimatedMissingKeys,
		&nullMetrics.TicketsInProgress,
		&nullMetrics.TicketsRequiredSelfHealing,
		&nullMetrics.SuccessfullySelfHealedTickets,
		&nullMetrics.SuccessfullyVerifiedTickets,
	)
	if err != nil {
		return metrics, err
	}

	// Assign values if they are not NULL
	if nullMetrics.SentTicketsForSelfHealing.Valid {
		metrics.SentTicketsForSelfHealing = int(nullMetrics.SentTicketsForSelfHealing.Int64)
	}
	if nullMetrics.EstimatedMissingKeys.Valid {
		metrics.EstimatedMissingKeys = int(nullMetrics.EstimatedMissingKeys.Int64)
	}
	if nullMetrics.TicketsInProgress.Valid {
		metrics.TicketsInProgress = int(nullMetrics.TicketsInProgress.Int64)
	}
	if nullMetrics.TicketsRequiredSelfHealing.Valid {
		metrics.TicketsRequiredSelfHealing = int(nullMetrics.TicketsRequiredSelfHealing.Int64)
	}
	if nullMetrics.SuccessfullySelfHealedTickets.Valid {
		metrics.SuccessfullySelfHealedTickets = int(nullMetrics.SuccessfullySelfHealedTickets.Int64)
	}
	if nullMetrics.SuccessfullyVerifiedTickets.Valid {
		metrics.SuccessfullyVerifiedTickets = int(nullMetrics.SuccessfullyVerifiedTickets.Int64)
	}

	return metrics, nil
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

	if _, err := db.Exec(createSelfHealingMessages); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingMessagesUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingMetricsTable); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createSelfHealingMetricsUniqueIndex); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	_, _ = db.Exec(alterTaskHistory)

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
