package meta

import (
	"encoding/hex"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
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
