package db

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/dgraph-io/badger/v2"
	"github.com/dgraph-io/badger/v2/options"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	replicationPrefix = "r:"
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

	log.WithContext(ctx).Debugf("data dir: %v", dataDir)
	// mkdir the data directory for badger
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, errors.Errorf("mkdir %q: %w", dataDir, err)
	}
	// init the badger options
	badgerOptions := badger.DefaultOptions(dataDir)
	badgerOptions = badgerOptions.WithCompression(options.ZSTD)
	// discard the log for badger
	logger := log.NewLogger()
	logger.SetOutput(ioutil.Discard)
	badgerOptions = badgerOptions.WithLogger(logger)

	// open the badger
	db, err := badger.Open(badgerOptions)
	if err != nil {
		return nil, errors.Errorf("open badger: %w", err)
	}
	s.db = db

	return s, nil
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (s *Badger) Store(ctx context.Context, key []byte, value []byte, replication time.Time) error {
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

		return nil
	}); err != nil {
		return errors.Errorf("badger update: %v, %w", base58.Encode(key), err)
	}

	log.WithContext(ctx).Infof("store key: %s", base58.Encode(key))
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

	log.WithContext(ctx).Infof("retrieve key: %s", base58.Encode(key))
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
		return txn.Delete([]byte(rk))
	}); err != nil {
		log.WithContext(ctx).WithError(err).Error("badger delete failed")
		return
	}
	log.WithContext(ctx).Infof("delete key: %s", base58.Encode(key))
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
			keys = append(keys, key)
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).WithError(err).Error("badger iterate keys failed")
	}
	return keys
}

// ForEachKey process each key
func (s *Badger) ForEachKey(_ context.Context, handler func(key []byte)) error {
	return s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()

			if bytes.HasPrefix(key, []byte(replicationPrefix)) {
				continue
			}

			handler(key)
		}
		return nil
	})
}

// ForEach will loop through all items and processing them

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
		log.WithContext(ctx).WithError(err).Errorf("badger iterate failed with prefix: %s", prefix)
	}

	return keys
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
			log.WithContext(ctx).WithError(err).Errorf("close badger failed")
		}
	}
}

// Stats returns stats of store
func (s *Badger) Stats(ctx context.Context) (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	lsm, vlog := s.db.Size()
	stats["log_size"] = vlog
	stats["dir_size"] = lsm + vlog
	stats["record_count"] = len(s.Keys(ctx))
	return stats, nil
}
