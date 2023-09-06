package sqlite

import (
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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
	testValues := []string{"testValue1", "testValue2"}
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
		var res result
		err = db.Get(&res, "SELECT key, id, ip, port FROM replication_keys WHERE key = ?", testValue)
		if err != nil {
			t.Fatalf("Failed to retrieve data: %v", err)
		}

		assert.Equal(t, testValue, res.Key)
		assert.Equal(t, testID, res.ID)
		assert.Equal(t, testIP, res.IP)
		assert.Equal(t, testPort, res.Port)
	}
}

func TestGetBatchRepKeys(t *testing.T) {
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
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS replication_keys(
        key TEXT PRIMARY KEY,
		id TEXT NOT NULL,
        ip TEXT NOT NULL,
		port INTEGER NOT NULL,
		attempts INTEGER DEFAULT 0,
		updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Define test data
	testValues := []string{"testValue1", "testValue2"}
	testID := "testID"
	testIP := "127.0.0.1"
	testPort := 8080

	// Call StoreBatchRepKeys
	err = store.StoreBatchRepKeys(testValues, testID, testIP, testPort)
	if err != nil {
		t.Fatalf("Failed to store batch rep keys: %v", err)
	}

	// Call StoreBatchRepKeys
	list, err := store.GetAllToDoRepKeys(0, 12)
	if err != nil {
		t.Fatalf("Failed to get batch rep keys: %v", err)
	}

	assert.Equal(t, 2, len(list))

}
