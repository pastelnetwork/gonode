package db

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/jbenet/go-base58"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/helpers"
)

const (
	defaultBadgerGCTime = 5 * time.Minute
)

// Badger defines a key/value database for store
type Badger struct {
	db   *badger.DB    // the badger for key/value persistence
	done chan struct{} // the badger is done
}

// NewStore returns a badger store
func NewStore(ctx context.Context, dataDir string) (*Badger, error) {
	s := &Badger{
		done: make(chan struct{}),
	}

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
	s.db = db

	// start a timer for the badger GC
	helpers.StartTimer(ctx, "badger gc", s.done, defaultBadgerGCTime, func() error {
		return db.RunValueLogGC(0.7)
	})

	return s, nil
}

// serialize from the data, replication time and expiration time
func (s *Badger) serialize(data []byte, replication time.Time, expiration time.Time) []byte {
	serialized := fmt.Sprintf("%d:%d:%s", replication.Unix(), expiration.Unix(), base58.Encode(data))
	return []byte(serialized)
}

// deserialize to the data, replication time and expiration time
func (s *Badger) deserialize(data []byte) ([]byte, int64, int64, error) {
	parts := strings.Split(string(data), ":")
	if len(parts) != 3 {
		return nil, 0, 0, errors.Errorf("invalid data: %v", string(data))
	}
	replication, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, 0, 0, errors.Errorf("parse replication time: %w", err)
	}
	expiration, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, 0, 0, errors.Errorf("parse expiration time: %w", err)
	}

	return []byte(parts[2]), replication, expiration, nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Badger) Store(ctx context.Context, key []byte, data []byte, replication time.Time, expiration time.Time) error {
	value := s.serialize(data, replication, expiration)

	if err := s.db.Update(func(txn *badger.Txn) error {
		if err := txn.Set(key, value); err != nil {
			return errors.Errorf("transaction set entry: %w", err)
		}
		return nil
	}); err != nil {
		return errors.Errorf("badger update: %v, %w", hex.EncodeToString(key), err)
	}

	log.WithContext(ctx).Debugf("store key: %s, data: %v", hex.EncodeToString(key), hex.EncodeToString(data))
	return nil
}

// Retrieve the local key/value from store
func (s *Badger) Retrieve(ctx context.Context, key []byte) ([]byte, error) {
	var buf []byte
	// find key from badger
	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return errors.Errorf("read from badger: %w", err)
		}
		if err := item.Value(func(val []byte) error {
			buf = append([]byte{}, val...)
			return nil
		}); err != nil {
			return errors.Errorf("read item value: %w", err)
		}
		return nil
	}); err != nil {
		return nil, errors.Errorf("badger view: %w", err)
	}
	// return nil, when key not found
	if buf == nil {
		return nil, nil
	}

	// deserialize the bytes to data, replication and expiration
	data, _, _, err := s.deserialize(buf)
	if err != nil {
		return nil, errors.Errorf("deserialize value: %w", err)
	}

	log.WithContext(ctx).Debugf("retrieve key: %s, data: %v", hex.EncodeToString(key), hex.EncodeToString(data))
	return data, nil
}

// Delete a key/value pair from the Store
func (s *Badger) Delete(ctx context.Context, key []byte) {
	// delete a key from badger
	if err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	}); err != nil {
		log.WithContext(ctx).Errorf("badger delete: %w", err)
	}
}

// Keys returns all the keys from the Store
func (s *Badger) Keys(ctx context.Context) [][]byte {
	keys := [][]byte{}

	if err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			keys = append(keys, it.Item().Key())
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).Errorf("badger iterate keys: %w", err)
	}
	return keys
}

// ExpireKeys should expire all key/values due for expiration
func (s *Badger) ExpireKeys(ctx context.Context) {
}

// KeysForReplication should return the keys of all data to be
// replicated across the network
func (s *Badger) KeysForReplication(ctx context.Context) [][]byte {
	var keys [][]byte

	if err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()

			key := item.Key()
			if err := item.Value(func(value []byte) error {
				// deserialize the data, replication and expiration
				_, replication, _, err := s.deserialize(value)
				if err != nil {
					return err
				}
				if time.Now().Unix() > replication {
					keys = append(keys, key)
				}

				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).Errorf("badger iterate keys/values: %w", err)
	}

	return keys
}

// Close the badger store
func (s *Badger) Close(ctx context.Context) {
	if s.done != nil {
		close(s.done)
	}
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.WithContext(ctx).Errorf("close badger: %w", err)
		}
	}
}
