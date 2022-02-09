package memory

import (
	"errors"
	"time"

	"github.com/pastelnetwork/gonode/common/storage"
	cache "github.com/patrickmn/go-cache"
)

type keyValue struct {
	store *cache.Cache
}

// Get implements storage.KeyValue.Set().
func (db *keyValue) Get(key string) ([]byte, error) {
	if value, ok := db.store.Get(key); ok {
		if bytes, ok := value.([]byte); ok {
			return bytes, nil
		}

		return nil, errors.New("unable to get bytes from value")
	}

	return nil, storage.ErrKeyValueNotFound
}

// Delete implements storage.KeyValue.Delete().
func (db *keyValue) Delete(key string) error {
	db.store.Delete(key)

	return nil
}

// Get implements storage.KeyValue.Set().
func (db *keyValue) Set(key string, value []byte) error {
	db.store.Set(key, value, cache.NoExpiration)
	return nil
}

// Get implements storage.KeyValue.Set().
func (db *keyValue) SetWithExpiry(key string, value []byte, expiry time.Duration) error {
	db.store.Set(key, value, expiry)
	return nil
}

// NewKeyValue returns a new key value storage instance
func NewKeyValue() storage.KeyValue {
	c := cache.New(-1*time.Minute, 12*time.Hour)
	return &keyValue{
		store: c,
	}
}
