package kademlia

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)

	return db
}

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	err = migrate(db)
	assert.NoError(t, err)
}

func TestStore(t *testing.T) {
	db := setupDB(t)

	rows, err := store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
}

func TestRetrieve(t *testing.T) {
	db := setupDB(t)

	rows, err := store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	data, err := retrieve(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"))
	assert.NoError(t, err)
	assert.Equal(t, "data", string(data))
}

func TestGetAllKeysForReplication(t *testing.T) {
	db := setupDB(t)

	_, err := store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	keys, err := getAllKeysForReplication(db)
	assert.NoError(t, err)
	assert.Len(t, keys, 2)
}

func TestExpireKeys(t *testing.T) {
	db := setupDB(t)

	_, err := store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	err = expireKeys(db)
	assert.NoError(t, err)
}

func TestRemove(t *testing.T) {
	db := setupDB(t)

	_, err := store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	_, err = store(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)

	err = remove(db, []byte("cf23df2207d99a74fbe169e3eba035e633b65d94"))
	assert.NoError(t, err)

	keys, err := getAllKeysForReplication(db)
	assert.NoError(t, err)
	assert.Len(t, keys, 0)
}
