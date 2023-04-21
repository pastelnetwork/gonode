package sqlite

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	dataDir, err := ioutil.TempDir("", "sqlite")
	require.NoError(t, err)
	defer os.RemoveAll(dataDir)

	replicateInterval := 5 * time.Minute
	republishInterval := 24 * time.Hour

	store, err := NewStore(ctx, dataDir, replicateInterval, republishInterval)
	require.NoError(t, err)

	tests := []struct {
		name        string
		key         []byte
		value       []byte
		expectedErr error
	}{
		{"Valid key-value", []byte("key1"), []byte("value1"), nil},
		{"Empty key", []byte{}, []byte("value2"), errors.New("key cannot be empty")},
		{"Nil key", []byte("value4"), []byte{}, errors.New("value cannot be empty")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := store.Store(ctx, test.key, test.value)
			if test.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestRetrieve(t *testing.T) {
	ctx := context.Background()
	dataDir, err := ioutil.TempDir("", "sqlite")
	require.NoError(t, err)
	defer os.RemoveAll(dataDir)

	replicateInterval := 5 * time.Minute
	republishInterval := 24 * time.Hour

	store, err := NewStore(ctx, dataDir, replicateInterval, republishInterval)
	require.NoError(t, err)

	testKey := []byte("key1")
	testValue := []byte("value1")
	err = store.Store(ctx, testKey, testValue)
	require.NoError(t, err)

	tests := []struct {
		name        string
		key         []byte
		expected    []byte
		expectedErr error
	}{
		{"Valid key", testKey, testValue, nil},
		{"Non-existent key", []byte("nonexistent"), nil, fmt.Errorf("failed to get record by key %s: sql: no rows in result set", hex.EncodeToString([]byte("nonexistent")))},
		{"Empty key", []byte{}, nil, fmt.Errorf("failed to get record by key %s: sql: no rows in result set", "")},
		{"Nil key", nil, nil, fmt.Errorf("failed to get record by key %s: sql: no rows in result set", "")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := store.Retrieve(ctx, test.key)
			if test.expectedErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, value)
			} else {
				assert.EqualError(t, err, test.expectedErr.Error())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	dataDir, err := ioutil.TempDir("", "sqlite")
	require.NoError(t, err)
	defer os.RemoveAll(dataDir)

	replicateInterval := 5 * time.Minute
	republishInterval := 24 * time.Hour

	store, err := NewStore(ctx, dataDir, replicateInterval, republishInterval)
	require.NoError(t, err)

	testKey := []byte("key1")
	testValue := []byte("value1")
	err = store.Store(ctx, testKey, testValue)
	require.NoError(t, err)

	tests := []struct {
		name       string
		key        []byte
		shouldFail bool
	}{
		{"Valid key", testKey, false},
		{"Non-existent key", []byte("nonexistent"), false},
		{"Empty key", []byte{}, false},
		{"Nil key", nil, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.Delete(ctx, test.key)

			// Verify deletion
			_, err := store.Retrieve(ctx, test.key)
			assert.Error(t, err, "failed to get record by key %s: sql: no rows in result set", hex.EncodeToString(test.key))
		})
	}
}
