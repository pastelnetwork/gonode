package local

import (
	"encoding/json"
	"fmt"
	"path/filepath"

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

const (
	historyDBName = "history.db"
	emptyString   = ""
)

// SQLiteStore handles sqlite ops
type SQLiteStore struct {
	db *sqlx.DB
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

	if _, err := db.Exec(alterTaskHistory); err != nil {
		log.WithError(err).Error("error details column to task history:")
	}

	return &SQLiteStore{
		db: db,
	}, nil
}
