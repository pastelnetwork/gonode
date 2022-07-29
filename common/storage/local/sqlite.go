package local

import (
	"fmt"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/rqlite/go-sqlite3" //go-sqlite3

	"github.com/pastelnetwork/gonode/common/storage"
)

const createTaskHistory string = `
  CREATE TABLE IF NOT EXISTS task_history (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  time DATETIME NOT NULL,
  task_id TEXT NOT NULL,
  status TEXT NOT NULL
  );`

const historyDBName = "history.db"

// SQLiteStore handles sqlite ops
type SQLiteStore struct {
	db *sqlx.DB
}

// InsertTaskHistory inserts task history
func (s *SQLiteStore) InsertTaskHistory(history types.TaskHistory) (int, error) {
	res, err := s.db.Exec("INSERT INTO task_history VALUES(NULL,?,?,?);", history.CreatedAt,
		history.TaskID, history.Status)

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
	rows, err := s.db.Query("SELECT * FROM task_history WHERE task_id = ? LIMIT 100", taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []types.TaskHistory{}
	for rows.Next() {
		i := types.TaskHistory{}
		err = rows.Scan(&i.ID, &i.CreatedAt, &i.TaskID, &i.Status)
		if err != nil {
			return nil, err
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

	return &SQLiteStore{
		db: db,
	}, nil
}
