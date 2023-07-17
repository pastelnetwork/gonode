package sqlite

import (
	"context"
	"crypto/rand"

	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoreAndRetrieve(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "sqlite-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	store, err := NewStore(context.Background(), dbPath, time.Minute, time.Minute)
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
