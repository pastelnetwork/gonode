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
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
)

const (
	replicationPrefix = "r:"

	// cleanupdiscardRatio of 0.5 would indicate that a file will be rewritten if half the space can be discarded.
	// This results in a lifetime value log write amplification of 2 (1 from original write + 0.5 rewrite + 0.25 + 0.125 + ... = 2).
	//  Setting it to higher value would result in fewer space reclaims, while setting it to a lower value would result in more space
	//  reclaims at the cost of increased activity on the LSM tree. discardRatio must be in the range (0.0, 1.0)
	cleanupDiscardRatio = 0.35
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
	if err := os.MkdirAll(dataDir, 0750); err != nil {
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
func (s *Badger) Store(_ context.Context, key []byte, value []byte, replication time.Time) error {
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

	return nil
}

// Retrieve the local key/value from store
func (s *Badger) Retrieve(_ context.Context, key []byte) ([]byte, error) {
	var value []byte
	// find key from badger
	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
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

// Keys return a list of keys with given offset + limit
func (s *Badger) Keys(ctx context.Context, offset int, limit int) [][]byte {
	keys := [][]byte{}
	if offset < 0 {
		offset = 0
	}
	curIndex := 0

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
			if limit >= 0 && curIndex >= offset+limit {
				break
			}

			if curIndex >= offset {
				keys = append(keys, key)
			}
			curIndex = curIndex + 1
		}
		return nil
	}); err != nil {
		log.WithContext(ctx).WithError(err).Error("badger iterate keys failed")
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
	stats["record_count"] = len(s.Keys(ctx, 0, -1))
	return stats, nil
}

// InitCleanup garbage collects the badger db. This should always run in a separate go-routine
func (s *Badger) InitCleanup(ctx context.Context, cleanupInterval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(cleanupInterval):
			for {
				if err := s.Cleanup(cleanupDiscardRatio); err != nil {
					break
				}
			}
		}
	}
}

// Cleanup garbage collects the badger db. Nil err means somethin was cleaned up
func (s *Badger) Cleanup(discardRatio float64) error {
	return s.db.RunValueLogGC(discardRatio)
}
