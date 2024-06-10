package ticketstore

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pastelnetwork/gonode/common/configurer"
	"github.com/pastelnetwork/gonode/common/log"
	"path/filepath"
)

const createFilesTable string = `
CREATE TABLE IF NOT EXISTS files (
    file_id TEXT NOT NULL PRIMARY KEY,
    upload_timestamp DATETIME,
    path TEXT,
    index TEXT,
    base_file_id TEXT,
    task_id TEXT,
    reg_txid TEXT,
    activation_txid TEXT,
    req_burn_txn_amount FLOAT,
    burn_txn_id TEXT,
    req_amount FLOAT,
    is_concluded BOOLEAN,
    cascade_metadata_ticket_id TEXT,
    uuid_key TEXT,
    hash_of_original_big_file TEXT,
    name_of_original_big_file_with_ext TEXT,
    size_of_original_big_file FLOAT,
    data_type_of_original_big_file TEXT,
    start_block INTEGER,
    done_block INTEGER,
    FOREIGN KEY (base_file_id) REFERENCES files(file_id)
);
`

const createRegistrationAttemptsTable string = `
CREATE TABLE IF NOT EXISTS registration_attempts (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    file_id TEXT NOT NULL,
    reg_started_at DATETIME,
    processor_sns TEXT,
    finished_at DATETIME,
    is_successful BOOLEAN,
    error_message TEXT,
    FOREIGN KEY (file_id) REFERENCES files(file_id)
);
`

const createActivationAttemptsTable string = `
CREATE TABLE IF NOT EXISTS activation_attempts (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    file_id TEXT NOT NULL,
    activation_attempt_at DATETIME,
    is_successful BOOLEAN,
    error_message TEXT,
    FOREIGN KEY (file_id) REFERENCES files(file_id)
);
`

const (
	ticketDBName = "ticket.db"
)

// TicketStore handles sqlite ops for tickets
type TicketStore struct {
	db *sqlx.DB
}

// CloseTicketDB closes ticket database
func (s *TicketStore) CloseTicketDB(ctx context.Context) {
	if err := s.db.Close(); err != nil {
		log.WithContext(ctx).WithError(err).Error("error closing ticket db")
	}
}

// OpenTicketingDb opens ticket DB
func OpenTicketingDb() (TicketStorageInterface, error) {
	dbFile := filepath.Join(configurer.DefaultPath(), ticketDBName)
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite database: %w", err)
	}

	if _, err := db.Exec(createFilesTable); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createRegistrationAttemptsTable); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
	}

	if _, err := db.Exec(createActivationAttemptsTable); err != nil {
		return nil, fmt.Errorf("cannot create table(s): %w", err)
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

	return &TicketStore{
		db: db,
	}, nil
}
