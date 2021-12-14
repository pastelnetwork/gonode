package db

import (
	"context"
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	storePath := "/tmp/testDb"
	// new the local storage
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	store, err := NewStore(ctx, storePath)
	assert.Nil(t, err)

	defer func() {
		store.Close(ctx)
		os.RemoveAll(storePath)
	}()

	numberOfKeys := 2000
	for i := 0; i < numberOfKeys; i++ {
		key := make([]byte, 32)
		rand.Read(key)
		writeData := make([]byte, 50*1024)
		rand.Read(writeData)

		err = store.Store(ctx, key, writeData, time.Now())
		assert.Nil(t, err)

		readData, err := store.Retrieve(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, writeData, readData)
	}

	assert.Equal(t, 1, len(store.Keys(ctx, 0, 1)))
	assert.Equal(t, 10, len(store.Keys(ctx, 0, 10)))
	assert.Equal(t, numberOfKeys, len(store.Keys(ctx, 0, -1)))
	assert.Equal(t, 0, len(store.Keys(ctx, numberOfKeys, 1)))
	assert.Equal(t, 0, len(store.Keys(ctx, numberOfKeys, -1)))
}

// Run this test to remove 20% symbolfiles for storage challenge test
// func TestRemoveKeysDB(t *testing.T) {
// 	home := os.Getenv("HOME")
// 	storePath := path.Join(home, ".pastel/p2p-data")
// 	// new the local storage
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()
// 	store, err := NewStore(ctx, storePath)
// 	assert.Nil(t, err)

// 	defer func() {
// 		store.Close(ctx)
// 	}()
// 	var keys = make([][]byte, 0)
// 	offset, limit, count := 0, 300, 0

// 	for ; count <= limit; offset += count - 1 {
// 		ks := store.Keys(ctx, offset, limit)
// 		keys = append(keys, ks...)
// 		count = len(ks)
// 	}

// 	for _, key := range keys {
// 		if rd.Float64() < 0.8 {
// 			// 80% keeps
// 			continue
// 		}
// 		// 20% delete
// 		store.Delete(ctx, key)
// 	}
// }
