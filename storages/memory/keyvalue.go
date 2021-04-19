package memory

import (
	"sync"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/walletnode/storages"
)

const logPrefix = "[memory]"

type keyValue struct {
	sync.Mutex
	values map[string][]byte
}

// Init implements storages.KeyValue.Init().
func (db *keyValue) Init() error {
	return nil
}

// Get implements storages.KeyValue.Set().
func (db *keyValue) Get(key string) ([]byte, error) {
	db.Lock()
	defer db.Unlock()

	if value, ok := db.values[key]; ok {
		return value, nil
	}
	return nil, storages.ErrKeyNotFound
}

// Delete implements storages.KeyValue.Delete().
func (db *keyValue) Delete(key string) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.values[key]; ok {
		log.WithField("key", key).Debugln(logPrefix, "Remove value")
		delete(db.values, key)
	}
	return nil
}

// Get implements storages.KeyValue.Set().
func (db *keyValue) Set(key string, value []byte) error {
	db.Lock()
	defer db.Unlock()

	log.WithField("key", key).WithField("bytes", len(value)).Debugln(logPrefix, "Add value")
	db.values[key] = value

	return nil
}

// NewKeyValue returns a new KeyValue instance.
func NewKeyValue() storages.KeyValue {
	return &keyValue{
		values: make(map[string][]byte),
	}
}
