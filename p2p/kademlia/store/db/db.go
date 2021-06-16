package db

import (
	"bytes"
	"context"
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
	replicationPrefix   = "r:"
	expirationPrefix    = "e:"
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

	log.WithContext(ctx).Infof("data dir: %v", dataDir)
	// mkdir the data directory for badger
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, errors.Errorf("mkdir %q: %w", dataDir, err)
	}
	// init the badger options
	badgerOptions := badger.DefaultOptions(dataDir)
	badgerOptions = badgerOptions.WithCompression(options.ZSTD)
	badgerOptions = badgerOptions.WithLogger(log.DefaultLogger)

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

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Badger) Store(ctx context.Context, key []byte, value []byte, replication time.Time, expiration time.Time) error {
	if err := s.db.Update(func(txn *badger.Txn) error {
		// set the key/value to badger
		if err := txn.Set(key, value); err != nil {
			return errors.Errorf("transaction set key/value: %w", err)
		}

		// set the replication to badger
		rk := replicationPrefix + base58.Encode(key)
		rv := strconv.FormatInt(replication.Unix(), 10)
		if err := txn.Set([]byte(rk), []byte(rv)); err != nil {
			return errors.Errorf("transaction set replication: %w", err)
		}

		// set the expiration to badger
		ek := expirationPrefix + base58.Encode(key)
		ev := strconv.FormatInt(expiration.Unix(), 10)
		if err := txn.Set([]byte(ek), []byte(ev)); err != nil {
			return errors.Errorf("transaction set expiration: %w", err)
		}

		return nil
	}); err != nil {
		return errors.Errorf("badger update: %v, %w", base58.Encode(key), err)
	}

	log.WithContext(ctx).Infof("store key: %s, data: %v", base58.Encode(key), base58.Encode(value))
	return nil
}

// Retrieve the local key/value from store
func (s *Badger) Retrieve(ctx context.Context, key []byte) ([]byte, error) {
	var value []byte
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
			value = append([]byte{}, val...)
			return nil
		}); err != nil {
			return errors.Errorf("read item value: %w", err)
		}
		return nil
	}); err != nil {
		return nil, errors.Errorf("badger view: %w", err)
	}

	log.WithContext(ctx).Infof("retrieve key: %s, value: %v", base58.Encode(key), base58.Encode(value))
	return value, nil
}

// Delete a key/value pair from the Store
func (s *Badger) Delete(ctx context.Context, key []byte) {
	// delete a key from badger
	if err := s.db.Update(func(txn *badger.Txn) error {
		// delete the key/value
		if err := txn.Delete(key); err != nil {
			return err
		}

		// delete the replication
		rk := replicationPrefix + base58.Encode(key)
		if err := txn.Delete([]byte(rk)); err != nil {
			return err
		}

		// delete the expiration
		ek := expirationPrefix + base58.Encode(key)
		if err := txn.Delete([]byte(ek)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.WithContext(ctx).Errorf("badger delete: %v", err)
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
			key := it.Item().Key()
			if bytes.HasPrefix(key, []byte(replicationPrefix)) {
				continue
			}
			if bytes.HasPrefix(key, []byte(expirationPrefix)) {
				continue
			}
			keys = append(keys, key)
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).Errorf("badger iterate keys: %v", err)
	}
	return keys
}

func (s *Badger) findKeysWithPrefix(ctx context.Context, prefix string) [][]byte {
	var keys [][]byte

	if err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		keyPrefix := []byte(prefix)
		for it.Seek(keyPrefix); it.ValidForPrefix(keyPrefix); it.Next() {
			item := it.Item()

			key := item.Key()
			if err := item.Value(func(value []byte) error {
				t, err := strconv.ParseInt(string(value), 10, 64)
				if err != nil {
					return errors.Errorf("strconv parse: %w", err)
				}

				if time.Now().Unix() > t {
					trimed := strings.TrimPrefix(string(key), prefix)
					keys = append(keys, base58.Decode(trimed))
				}

				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).Errorf("badger iterate with prefix: %s, %s", prefix, err)
	}

	return keys
}

// KeysForExpiration returns the keys of all data to be removed from local storage
func (s *Badger) KeysForExpiration(ctx context.Context) [][]byte {
	return s.findKeysWithPrefix(ctx, expirationPrefix)
}

// KeysForReplication returns the keys of all data to be replicated across the network
func (s *Badger) KeysForReplication(ctx context.Context) [][]byte {
	return s.findKeysWithPrefix(ctx, replicationPrefix)
}

// Close the badger store
func (s *Badger) Close(ctx context.Context) {
	if s.done != nil {
		close(s.done)
	}
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.WithContext(ctx).Errorf("close badger: %s", err)
		}
	}
}
