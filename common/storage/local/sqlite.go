package local

import (
	"database/sql"
	"path/filepath"

	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/types"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/storage"
)

const create_task_history string = `
  CREATE TABLE IF NOT EXISTS task_history (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  time DATETIME NOT NULL,
  task_id TEXT NOT NULL,
  status TEXT NOT NULL
  );`

const historyDBName = "history.db"

type SQLiteStore struct {
	db *sql.DB
}

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

func OpenHistoryDB() (storage.LocalStoreInterface, error) {
	db, err := sql.Open("sqlite3", filepath.Join(configurer.DefaultPath(), historyDBName))
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create_task_history); err != nil {
		return nil, err
	}

	return &SQLiteStore{
		db: db,
	}, nil
}
