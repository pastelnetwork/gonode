package sqlite

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
	store, err := NewStore(ctx, storePath, time.Second*10, time.Second*60)
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

		err = store.Store(ctx, key, writeData)
		assert.Nil(t, err)

		readData, err := store.Retrieve(ctx, key)
		assert.Nil(t, err)
		assert.Equal(t, writeData, readData)
	}

	assert.Equal(t, 0, len(store.GetKeysForReplication(ctx)))
	time.Sleep(time.Second * 10)
	assert.Equal(t, numberOfKeys, len(store.GetKeysForReplication(ctx)))
	assert.Equal(t, 0, len(store.GetKeysForReplication(ctx)))
}
