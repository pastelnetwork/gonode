package memory

import (
	"sync"

	"github.com/pastelnetwork/go-commons/log"
	"github.com/pastelnetwork/walletnode/dao"
)

const logPrefix = "[memory]"

type keyValue struct {
	sync.Mutex
	values map[string][]byte
}

// Init implements dao.KeyValue.Init
func (db *keyValue) Init() error {
	return nil
}

// Get implements dao.KeyValue.Set
func (db *keyValue) Get(key []byte) ([]byte, error) {
	db.Lock()
	defer db.Unlock()

	if value, ok := db.values[string(key)]; ok {
		return value, nil
	}
	return nil, dao.ErrKeyNotFound
}

// Delete implements dao.KeyValue.Delete
func (db *keyValue) Delete(key []byte) error {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.values[string(key)]; ok {
		log.WithField("key", string(key)).Debugln(logPrefix, "Remove value")
		delete(db.values, string(key))
	}
	return nil
}

// Get implements dao.KeyValue.Set
func (db *keyValue) Set(key, value []byte) error {
	db.Lock()
	defer db.Unlock()

	log.WithField("key", string(key)).WithField("bytes", len(value)).Debugln(logPrefix, "Add value")
	db.values[string(key)] = value

	return nil
}

// NewKeyValue returns a new KeyValue instance
func NewKeyValue() dao.KeyValue {
	return &keyValue{
		values: make(map[string][]byte),
	}
}
