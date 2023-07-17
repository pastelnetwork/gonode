package sqlite

import (
	"encoding/hex"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreBatchRepKeys(t *testing.T) {
	// Set up a temporary in-memory SQLite database
	db, err := sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the store
	store := Store{
		db: db,
	}

	// Create the table
	_, err = db.Exec(`CREATE TABLE replication_keys (
		key TEXT NOT NULL,
		updatedAt DATETIME NOT NULL,
		ip TEXT NOT NULL,
		port INTEGER NOT NULL,
		id TEXT NOT NULL,
		PRIMARY KEY(key)
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Define test data
	testValues := [][]byte{[]byte("testValue1"), []byte("testValue2")}
	testID := "testID"
	testIP := "127.0.0.1"
	testPort := 8080

	// Call StoreBatchRepKeys
	err = store.StoreBatchRepKeys(testValues, testID, testIP, testPort)
	if err != nil {
		t.Fatalf("Failed to store batch rep keys: %v", err)
	}

	// Retrieve and verify the data
	type result struct {
		Key  string
		ID   string
		IP   string
		Port int
	}

	for _, testValue := range testValues {
		hkey := hex.EncodeToString(testValue)
		var res result
		err = db.Get(&res, "SELECT key, id, ip, port FROM replication_keys WHERE key = ?", hkey)
		if err != nil {
			t.Fatalf("Failed to retrieve data: %v", err)
		}

		assert.Equal(t, hkey, res.Key)
		assert.Equal(t, testID, res.ID)
		assert.Equal(t, testIP, res.IP)
		assert.Equal(t, testPort, res.Port)
	}
}
