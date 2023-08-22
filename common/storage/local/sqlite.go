package local

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
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
	now := time.Now()
	const insertQuery = "INSERT INTO storage_challenge_messages(id, challenge_id, message_type, data, sender_id, sender_signature, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?,?);"
	_, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MessageType, challenge.Data, challenge.Sender, challenge.SenderSignature, now, now)
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
func (s *SQLiteStore) InsertSelfHealingChallenge(challenge types.SelfHealingChallenge) (hID int, err error) {
	now := time.Now()
	const insertQuery = "INSERT INTO self_healing_challenges(id, challenge_id, merkleroot, file_hash, challenging_node, responding_node, verifying_node, reconstructed_file_hash, status, created_at, updated_at) VALUES(NULL,$1,$2,$3,$4,$5,$6,$7,$8,$9,$9);"

	res, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.MerkleRoot, challenge.FileHash, challenge.ChallengingNode, challenge.RespondingNode, challenge.VerifyingNode, challenge.ReconstructedFileHash, challenge.Status, now)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
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

	if _, err := db.Exec(createSelfHealingChallenges); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	_, _ = db.Exec(alterTaskHistory)

	return &SQLiteStore{
		db: db,
	}, nil
}
