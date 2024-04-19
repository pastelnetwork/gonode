package rqstore

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGetToDoStoreSymbolDirs(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		setup       func(db *sqlx.DB) // Function to setup database conditions for each test
		expected    []SymbolDir
		expectError bool
	}{
		{
			name: "All matching entries",
			setup: func(db *sqlx.DB) {
				// Insert multiple entries, all should be returned
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir, is_first_batch_stored, is_completed) VALUES ('tx1', 'dir1', TRUE, FALSE)")
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir, is_first_batch_stored, is_completed) VALUES ('tx2', 'dir2', TRUE, FALSE)")
			},
			expected: []SymbolDir{
				{TXID: "tx1", Dir: "dir1"},
				{TXID: "tx2", Dir: "dir2"},
			},
			expectError: false,
		},
		{
			name: "Some not matching entries",
			setup: func(db *sqlx.DB) {
				// Insert both matching and non-matching entries
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir, is_first_batch_stored, is_completed) VALUES ('tx1', 'dir1', TRUE, FALSE)")
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir, is_first_batch_stored, is_completed) VALUES ('tx2', 'dir2', FALSE, TRUE)")
			},
			expected: []SymbolDir{
				{TXID: "tx1", Dir: "dir1"},
			},
			expectError: false,
		},
		{
			name: "No matching entries",
			setup: func(db *sqlx.DB) {
				// All entries do not match the condition
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir, is_first_batch_stored, is_completed) VALUES ('tx1', 'dir1', FALSE, TRUE)")
			},
			expected:    nil,
			expectError: false,
		},
		{
			name: "Empty database",
			setup: func(db *sqlx.DB) {
				// No setup needed, testing empty state
			},
			expected:    nil,
			expectError: false,
		},
		{
			name: "Database error",
			setup: func(db *sqlx.DB) {
				// Simulate a database error by dropping the table before querying
				db.Exec("DROP TABLE rq_symbols_dir")
			},
			expected:    nil,
			expectError: true,
		},
	}

	// Run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := SetupTestDB(t)
			defer db.Close()

			// Setup the database as per the test case
			tc.setup(db.db)

			// Execute the function to test
			result, err := db.GetToDoStoreSymbolDirs()

			// Check expected results
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestStoreSymbolDirectory(t *testing.T) {
	tests := []struct {
		name          string
		txid          string
		dir           string
		setup         func(*sqlx.DB)
		expectError   bool
		verifyResults func(*sqlx.DB, string, string) error // Additional function to verify the results
	}{
		{
			name:        "successful insert",
			txid:        "tx123",
			dir:         "dir123",
			expectError: false,
			verifyResults: func(db *sqlx.DB, txid string, dir string) error {
				var count int
				err := db.Get(&count, "SELECT COUNT(*) FROM rq_symbols_dir WHERE txid = ? AND dir = ?", txid, dir)
				if err != nil {
					return err
				}
				if count != 1 {
					return fmt.Errorf("expected 1 entry in database, got %d", count)
				}
				return nil
			},
		},
		{
			name:        "failed due to duplicate txid",
			txid:        "tx123",
			dir:         "dir456",
			expectError: true,
			setup: func(db *sqlx.DB) {
				db.Exec("INSERT INTO rq_symbols_dir (txid, dir) VALUES (?, ?)", "tx123", "dir123")
			},
		},
		{
			name:        "failed due to SQL error",
			txid:        "txXYZ",
			dir:         "dirXYZ",
			expectError: true,
			setup: func(db *sqlx.DB) {
				db.Exec("DROP TABLE rq_symbols_dir") // Drop the table to induce an SQL error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := SetupTestDB(t)
			defer store.Close()

			if tt.setup != nil {
				tt.setup(store.db)
			}

			err := store.StoreSymbolDirectory(tt.txid, tt.dir)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.verifyResults != nil {
					assert.NoError(t, tt.verifyResults(store.db, tt.txid, tt.dir))
				}
			}
		})
	}
}
