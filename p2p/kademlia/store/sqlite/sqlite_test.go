package sqlite

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/cloud.go"
	"github.com/stretchr/testify/assert"
)

func TestStoreAndRetrieve(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "sqlite-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	store, err := NewStore(context.Background(), dbPath, nil, nil)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close(context.Background())

	testCases := []struct {
		key   []byte
		value []byte
	}{
		{generateRandomBytes(16), []byte("value1")},
		{generateRandomBytes(16), []byte("value2")},
		{generateRandomBytes(16), []byte("value3")},
	}

	// Store test cases
	for _, tc := range testCases {
		err := store.Store(context.Background(), tc.key, tc.value, 0, true)
		assert.NoError(t, err, "Store should not return an error")
	}
	time.Sleep(2 * time.Second)
	// Retrieve test cases
	for _, tc := range testCases {
		retrievedValue, err := store.Retrieve(context.Background(), tc.key)
		assert.NoError(t, err, "Retrieve should not return an error")
		assert.Equal(t, tc.value, retrievedValue, "Retrieved value should match the stored value")
	}
}

// Helper function to generate random bytes
func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func TestStore(t *testing.T) {
	cloud := cloud.NewRcloneStorage("test", "test")

	mst, err := NewMigrationMetaStore(context.Background(), ".", cloud)

	// override the tickers for testing
	mst.updateTicker = time.NewTicker(2 * time.Second)
	mst.insertTicker = time.NewTicker(2 * time.Second)

	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	store, err := NewStore(context.Background(), ".", cloud, mst)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	r1 := []byte("test-record-1")
	r2 := []byte("test-record-2")
	r3 := []byte("test-record-3")

	hashed, err := utils.Sha3256hash(r1)
	if err != nil {
		t.Fatalf("failed to hash record: %v", err)
	}

	r1Key := hex.EncodeToString(hashed)

	hashed, err = utils.Sha3256hash(r2)
	if err != nil {
		t.Fatalf("failed to hash record: %v", err)
	}

	r2Key := hex.EncodeToString(hashed)

	hashed, err = utils.Sha3256hash(r3)
	if err != nil {
		t.Fatalf("failed to hash record: %v", err)
	}

	r3Key := hex.EncodeToString(hashed)

	err = store.storeBatchRecord([][]byte{r1, r2, r3}, 0, true)
	if err != nil {
		t.Fatalf("failed to store record: %v", err)
	}

	time.Sleep(3 * time.Second)

	type record struct {
		Key           string    `db:"key"`
		LastAcccessed time.Time `db:"last_accessed"`
		AccessCount   int       `db:"access_count"`
		DataSize      int       `db:"data_size"`
	}

	var keys []record
	err = store.migrationStore.db.Select(&keys, "SELECT key,last_accessed,access_count,data_size FROM meta where key in (?, ?, ?)", r1Key, r2Key, r3Key)
	if err != nil {
		t.Fatalf("failed to retrieve record: %v", err)
	}

	if len(keys) != 3 {
		t.Fatalf("expected 3 records, got %d", len(keys))
	}

	time.Sleep(1 * time.Second)

	_, _, err = store.RetrieveBatchValues(context.Background(), []string{r1Key, r2Key, r3Key}, true)
	if err != nil {
		t.Fatalf("failed to retrieve record: %v", err)
	}

	time.Sleep(3 * time.Second)

	var nkeys []record
	err = store.migrationStore.db.Select(&nkeys, "SELECT  key,last_accessed,access_count,data_size FROM meta where key in (?, ?, ?)", r1Key, r2Key, r3Key)
	if err != nil {
		t.Fatalf("failed to retrieve record: %v", err)
	}

	if len(nkeys) != 3 {
		t.Fatalf("expected 3 records, got %d", len(nkeys))
	}

	for _, key := range nkeys {
		for _, k := range keys {
			if key.Key == k.Key {
				fmt.Println(key.LastAcccessed.String(), k.LastAcccessed.String())
				if !key.LastAcccessed.After(k.LastAcccessed) {
					t.Fatalf("last accessed time not updated")
				}

				if key.AccessCount != k.AccessCount+1 {
					t.Fatalf("access count not updated")
				}
			}
		}
	}

	os.Remove("data001.sqlite3")
	os.Remove("data001-migration-meta.sqlite3")
	os.Remove("data001.sqlite3-shm")
	os.Remove("data001.sqlite3-wal")

}
