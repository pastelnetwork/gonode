package rqstore

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const createRQSymbolsDir string = `
  CREATE TABLE IF NOT EXISTS rq_symbols_dir (
  txid TEXT NOT NULL,
  dir TEXT NOT NULL,
  is_first_batch_stored BOOLEAN NOT NULL DEFAULT FALSE,
  is_completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (txid)
);`

type Store interface {
	DeleteSymbolsByTxID(txid string) error
	StoreSymbolDirectory(txid, dir string) error
	GetDirectoryByTxID(txid string) (string, error)
	GetToDoStoreSymbolDirs() ([]SymbolDir, error)
	SetIsCompleted(txid string) error
	UpdateIsFirstBatchStored(txid string) error
}

// SQLiteRQStore is sqlite implementation of DD store and Score store
type SQLiteRQStore struct {
	db *sqlx.DB
}

type SymbolDir struct {
	TXID string `db:"txid"`
	Dir  string `db:"dir"`
}

// NewSQLiteRQStore is new sqlite store constructor
func NewSQLiteRQStore(file string) (*SQLiteRQStore, error) {
	db, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("cannot open rq-service database: %w", err)
	}

	// Create the rq_symbols_dir table if it doesn't exist
	_, err = db.Exec(createRQSymbolsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create rq_symbols table: %w", err)
	}

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

	return &SQLiteRQStore{
		db: db,
	}, nil
}

func (s *SQLiteRQStore) Close() error {
	return s.db.Close()
}

func SetupTestDB(t *testing.T) *SQLiteRQStore {
	// Create an in-memory SQLite instance for testing
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not connect to in-memory sqlite instance: %v", err)
	}

	// Create the rq_symbols table
	_, err = db.Exec(createRQSymbolsDir)
	if err != nil {
		t.Fatalf("Could not create test table: %v", err)
	}

	return &SQLiteRQStore{db: db}
}

// DeleteSymbolsByTxID deletes all symbols entries associated with a specific txid.
func (s *SQLiteRQStore) DeleteSymbolsByTxID(txid string) error {
	// Prepare the delete statement
	stmt, err := s.db.Prepare("DELETE FROM rq_symbols WHERE txid = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	// Execute the delete statement
	_, err = stmt.Exec(txid)
	if err != nil {
		return fmt.Errorf("failed to execute delete statement: %w", err)
	}

	return nil
}

// StoreSymbolDirectory associates a txid with a directory path and stores it in the database.
func (s *SQLiteRQStore) StoreSymbolDirectory(txid, dir string) error {
	stmt, err := s.db.Prepare("INSERT INTO rq_symbols_dir (txid, dir) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement for directory: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(txid, dir); err != nil {
		return fmt.Errorf("failed to execute insert statement for directory: %w", err)
	}

	return nil
}

// GetDirectoryByTxID retrieves the directory path associated with a given txid.
func (s *SQLiteRQStore) GetDirectoryByTxID(txid string) (string, error) {
	var dir string
	err := s.db.Get(&dir, "SELECT dir FROM rq_symbols_dir WHERE txid = ?", txid)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve directory for txid %s: %w", txid, err)
	}

	return dir, nil
}

// GetToDoStoreSymbolDirs retrieves directories and their txids where is_first_batch_stored is true and is_completed is false
func (s *SQLiteRQStore) GetToDoStoreSymbolDirs() ([]SymbolDir, error) {
	var symbolDirs []SymbolDir
	query := "SELECT txid, dir FROM rq_symbols_dir WHERE is_first_batch_stored = TRUE AND is_completed = FALSE"
	err := s.db.Select(&symbolDirs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve symbol directories: %w", err)
	}
	return symbolDirs, nil
}

// UpdateIsFirstBatchStored sets is_first_batch_stored to true for a given txid
func (s *SQLiteRQStore) UpdateIsFirstBatchStored(txid string) error {
	stmt, err := s.db.Prepare("UPDATE rq_symbols_dir SET is_first_batch_stored = TRUE WHERE txid = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(txid); err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	return nil
}

// SetIsCompleted sets is_completed to true for a given txid
func (s *SQLiteRQStore) SetIsCompleted(txid string) error {
	stmt, err := s.db.Prepare("UPDATE rq_symbols_dir SET is_completed = TRUE WHERE txid = ?")
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(txid); err != nil {
		return fmt.Errorf("failed to execute update statement: %w", err)
	}

	return nil
}
