package db

import (
	"context"
	"encoding/hex"
	"os"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

// Badger defines a key/value database for store
type Badger struct {
	db *badger.DB // the badger for key/value persistence
}

// NewStore returns a badger store
func NewStore(dataDir string) (*Badger, error) {
	// mkdir the data directory for badger
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, errors.Errorf("mkdir %q: %w", dataDir, err)
	}
	// init the badger options
	badgerOptions := badger.DefaultOptions(dataDir)
	badgerOptions = badgerOptions.WithCompression(options.ZSTD)

	// open the badger
	db, err := badger.Open(badgerOptions)
	if err != nil {
		return nil, errors.Errorf("open badger: %w", err)
	}

	return &Badger{db: db}, nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Badger) Store(ctx context.Context, key []byte, data []byte, replication time.Time, expiration time.Time) error {
	if err := s.db.Update(func(txn *badger.Txn) error {

		return nil
	}); err != nil {
		return errors.Errorf("badger update: %v, %w", hex.EncodeToString(key), err)
	}
	return nil
}

// Close the badger store
func (s *Badger) Close(ctx context.Context) {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.WithContext(ctx).Errorf("close badger: %w", err)
		}
	}
}
