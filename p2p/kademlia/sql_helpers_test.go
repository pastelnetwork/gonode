package kademlia

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_migrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)
}

func Test_store(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	rows, err := store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}


func Test_retrieve(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	rows, err := store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	data, err := retrieve(db, []byte("sample key"))
	assert.NoError(t, err)
	assert.Equal(t, "data", string(data))
}

func Test_getAllKeysForReplication(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	keys, err := getAllKeysForReplication(db)
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
}

func Test_expireKeys(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	err = expireKeys(db)
	assert.NoError(t, err)
}

func Test_remove(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	err = remove(db, []byte("sample key"))
	assert.NoError(t, err)

	keys, err := getAllKeysForReplication(db)
	assert.NoError(t, err)
	assert.Len(t, keys, 0)
}