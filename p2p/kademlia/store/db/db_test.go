package db

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	// new the local storage
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	store, err := NewStore(ctx, "/tmp/testDb")
	assert.Nil(t, err)

	for i := 0; i <= 2000; i++ {
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
}
