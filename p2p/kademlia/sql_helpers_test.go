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

	err = store(db, []byte("sample key"), []byte("data"), time.Now(), time.Now())
	assert.NoError(t, err)
}
