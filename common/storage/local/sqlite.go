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

const createFailedStorageChallenges string = `
  CREATE TABLE IF NOT EXISTS failed_storage_challenges (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  challenge_id TEXT NOT NULL,
  file_hash TEXT NOT NULL,
  status TEXT NOT NULL,
  responding_node TEXT NOT NULL,
  file_reconstructing_node TEXT,
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

// InsertFailedStorageChallenge inserts failed storage challenge to db
func (s *SQLiteStore) InsertFailedStorageChallenge(challenge types.FailedStorageChallenge) (hID int, err error) {
	now := time.Now()
	const insertQuery = "INSERT INTO failed_storage_challenges(id, challenge_id, file_hash, status, responding_node, created_at, updated_at) VALUES(NULL,?,?,?,?,?,?);"

	res, err := s.db.Exec(insertQuery, challenge.ChallengeID, challenge.FileHash, challenge.Status, challenge.RespondingNode, now, now)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}

	return int(id), nil
}

// QueryFailedStorageChallenges retrieves failed challenges stored in DB for self-healing
func (s *SQLiteStore) QueryFailedStorageChallenges() (challenges []types.FailedStorageChallenge, err error) {
	const selectQuery = "SELECT * FROM failed_storage_challenges WHERE status = ?"
	rows, err := s.db.Query(selectQuery, types.CreatedSelfHealingStatus)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		challenge := types.FailedStorageChallenge{}
		err = rows.Scan(&challenge.ID, &challenge.ChallengeID, &challenge.FileHash, &challenge.Status,
			&challenge.RespondingNode, challenge.CreatedAt, challenge.UpdatedAt)
		if err != nil {
			return nil, err
		}

		challenges = append(challenges, challenge)
	}

	return challenges, nil
}

// UpdateFailedStorageChallenge update the given failed storage challenge by ID
func (s *SQLiteStore) UpdateFailedStorageChallenge(challenge types.FailedStorageChallenge) (err error) {
	const updateQuery = "update failed_storage_challenges set status = ?, file_reconstructing_node = ?, updated_at= ? WHERE id = ?"
	_, err = s.db.Exec(updateQuery, challenge.Status, challenge.FileReconstructingNode, time.Now(), challenge.ID)
	if err != nil {
		return err
	}

	return nil
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

	if _, err := db.Exec(createFailedStorageChallenges); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	_, _ = db.Exec(alterTaskHistory)

	return &SQLiteStore{
		db: db,
	}, nil
}
