package rqstore

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

const createRQSymbolsTable string = `
  CREATE TABLE IF NOT EXISTS rq_symbols (
  txid TEXT NOT NULL,
  hash TEXT NOT NULL,
  data BLOB NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

type Store interface {
	StoreSymbols(id string, symbols map[string][]byte) error
	LoadSymbols(id string, symbols map[string][]byte) (map[string][]byte, error)
}

// SQLiteRQStore is sqlite implementation of DD store and Score store
type SQLiteRQStore struct {
	db *sqlx.DB
}

// NewSQLiteRQStore is new sqlite store constructor
func NewSQLiteRQStore(file string) (*SQLiteRQStore, error) {
	db, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("cannot open rq-service database: %w", err)
	}

	// Create the rq_symbols table if it doesn't exist
	_, err = db.Exec(createRQSymbolsTable)
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

// StoreSymbols inserts all the symbols in the given table 'rq_symbols' through a batch insert using a transaction.
func (s *SQLiteRQStore) StoreSymbols(id string, symbols map[string][]byte) error {
	// Begin a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO rq_symbols (txid, hash, data) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback() // Roll back the transaction in case of an error
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for hash, data := range symbols {
		if _, err := stmt.Exec(id, hash, data); err != nil {
			tx.Rollback() // Roll back the transaction in case of an error
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LoadSymbols loads the values for the given keys (hashes) and populates the symbols map with the loaded values.
func (s *SQLiteRQStore) LoadSymbols(id string, symbols map[string][]byte) (map[string][]byte, error) {
	if len(symbols) == 0 {
		return symbols, nil
	}

	hashes := make([]string, 0, len(symbols))
	for hash := range symbols {
		hashes = append(hashes, hash)
	}

	// Construct the query with placeholders for each hash
	query := fmt.Sprintf("SELECT hash, data FROM rq_symbols WHERE txid = ? AND hash IN (%s)",
		strings.Repeat("?,", len(hashes)-1)+"?")

	// Prepare arguments for the query
	args := make([]interface{}, 0, len(hashes)+1)
	args = append(args, id) // Add txid as the first argument
	for _, hash := range hashes {
		args = append(args, hash)
	}

	rows, err := s.db.Queryx(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to load symbols: %w", err)
	}
	defer rows.Close()

	// Iterate over the result set
	loadedSymbols := make(map[string][]byte)
	for rows.Next() {
		var hash string
		var data []byte
		if err := rows.Scan(&hash, &data); err != nil {
			return nil, fmt.Errorf("failed to scan symbol: %w", err)
		}
		loadedSymbols[hash] = data
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over symbols: %w", err)
	}

	return loadedSymbols, nil
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
	_, err = db.Exec(createRQSymbolsTable)
	if err != nil {
		t.Fatalf("Could not create test table: %v", err)
	}

	return &SQLiteRQStore{db: db}
}
