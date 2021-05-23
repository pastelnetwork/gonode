package memory

import (
	"sync"

	"github.com/pastelnetwork/gonode/common/storage"
)

type keyValue struct {
	sync.Mutex
	values map[string][]byte
}

// Get implements storage.KeyValue.Set().
func (db *keyValue) Get(key string) ([]byte, error) {
	db.Lock()
	defer db.Unlock()

	if value, ok := db.values[key]; ok {
		return value, nil
	}
	return nil, storage.ErrKeyValueNotFound
}

// Delete implements storage.KeyValue.Delete().
func (db *keyValue) Delete(key string) error {
	db.Lock()
	defer db.Unlock()

	delete(db.values, key)
	return nil
}

// Get implements storage.KeyValue.Set().
func (db *keyValue) Set(key string, value []byte) error {
	db.Lock()
	defer db.Unlock()

	db.values[key] = value
	return nil
}

// NewKeyValue return new instance key value storage
func NewKeyValue() storage.KeyValue {
	return &keyValue{
		values: make(map[string][]byte),
	}
}
