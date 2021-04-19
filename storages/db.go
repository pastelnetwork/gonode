package storages

import "github.com/pastelnetwork/go-commons/errors"

var (
	// ErrKeyNotFound is returned when key isn't found.
	ErrKeyNotFound = errors.New("key not found")
)

// KeyValue represents database that uses a simple key-value method to store data.
type KeyValue interface {
	// Init initializes the database.
	Init() error

	// Get looks for key and returns corresponding Item.
	// If key is not found, ErrKeyNotFound is returned.
	Get(key string) (value []byte, err error)

	// Set adds a key-value pair to the database.
	Set(key string, value []byte) (err error)

	// Delete deletes a key.
	Delete(key string) error
}
