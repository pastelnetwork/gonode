package sqlite

import (
	"context"
	"database/sql"
	"time"

	// add sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/pastelnetwork/gonode/common/errors"
)

// Key is a simple in-memory key/value store used for unit testing, and
// the CLI example
type Key struct {
	db *sql.DB
}

// GetAllKeysForReplication should return the keys of all data to be
// replicated across the network. Typically all data should be
// replicated every tReplicate seconds.
func (k *Key) GetAllKeysForReplication(ctx context.Context) ([][]byte, error) {
	return GetAllKeysForReplication(ctx, k.db)
}

// ExpireKeys should expire all key/values due for expiration.
func (k *Key) ExpireKeys(ctx context.Context) error {
	return ExpireKeys(ctx, k.db)
}

// Init initializes the Store
func (k *Key) Init(ctx context.Context) (closerr error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return errors.Errorf("failed opening new sqlite connection: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			closerr = errors.Errorf("failed closing db connection: %w", err)
		}
	}()

	if err := Migrate(ctx, k.db); err != nil {
		return err
	}

	k.db = db
	return closerr
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (k *Key) Store(ctx context.Context, data []byte, replication time.Time, expiration time.Time, _ bool) error {
	_, err := Store(ctx, k.db, data, replication, expiration)
	return err
}

// Retrieve will return the local key/value if it exists
func (k *Key) Retrieve(ctx context.Context, key []byte) (data []byte, found bool) {
	data, err := Retrieve(ctx, k.db, key)
	return data, err == nil
}

// Delete deletes a key/value pair from the Key
func (k *Key) Delete(ctx context.Context, key []byte) error {
	return Remove(ctx, k.db, key)
}
