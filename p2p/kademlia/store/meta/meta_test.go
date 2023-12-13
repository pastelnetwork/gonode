package meta

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreDisabledKey(t *testing.T) {
	// Create an in-memory SQLite database
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Create the disabled_keys table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS disabled_keys(
        key TEXT PRIMARY KEY,
        createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	// Prepare the store
	store := &Store{
		db: db,
		// worker and other store fields should be initialized as needed
	}

	// Define the test cases
	testCases := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     []byte("valid_key"),
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     []byte(""),
			wantErr: true, // Depending on how you want to handle empty keys
		},
		// Add more test cases as needed
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := store.storeDisabledKey(tc.key)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the key is stored
				var hkey string
				query := "SELECT key FROM disabled_keys WHERE key = ?"
				err := db.Get(&hkey, query, hex.EncodeToString(tc.key))
				require.NoError(t, err)
				assert.Equal(t, hex.EncodeToString(tc.key), hkey)
			}
		})
	}
}

// TestStore encapsulates testing functions for Store
func TestStore(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	require.NoError(t, err)

	// Prepare the store
	store := &Store{
		db: db,
		// worker and other store fields should be initialized as needed
	}

	err = store.migrateToDelKeys()
	require.NoError(t, err)

	testCases := []struct {
		name         string
		setup        func() error
		testFunc     func() error
		validateFunc func() error
	}{
		{
			name: "Insert New Key",
			setup: func() error {
				// Reset the database state
				_, err := db.Exec("DELETE FROM del_keys")
				if err != nil {
					return err
				}

				return nil // No setup required for this test case
			},
			testFunc: func() error {
				return store.BatchInsertDelKeys(context.Background(), domain.DelKeys{
					{Key: "key1", Count: 1, CreatedAt: time.Now()},
				})
			},
			validateFunc: func() error {
				keys, err := store.GetAllToDelKeys(1)
				if err != nil {
					return err
				}
				if len(keys) != 1 || keys[0].Key != "key1" || keys[0].Count != 1 {
					return fmt.Errorf("unexpected result: %+v", keys)
				}
				return nil
			},
		},
		{
			name: "Increment Count on Duplicate Insert",
			setup: func() error {
				// Reset the database state
				_, err := db.Exec("DELETE FROM del_keys")
				if err != nil {
					return err
				}
				// Pre-insert a key
				return store.BatchInsertDelKeys(context.Background(), domain.DelKeys{
					{Key: "key1", Count: 1, CreatedAt: time.Now()},
				})
			},
			testFunc: func() error {
				return store.BatchInsertDelKeys(context.Background(), domain.DelKeys{
					{Key: "key1", Count: 1, CreatedAt: time.Now()},
				})
			},
			validateFunc: func() error {
				keys, err := store.GetAllToDelKeys(1)
				if err != nil {
					return err
				}
				if len(keys) != 1 || keys[0].Key != "key1" || keys[0].Count != 2 {
					return fmt.Errorf("count did not increment as expected: %+v", keys)
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, tc.setup())
			require.NoError(t, tc.testFunc())
			require.NoError(t, tc.validateFunc())
		})
	}
}
