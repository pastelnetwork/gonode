package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// Store implements an in-memory cache
type Store interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Reset(expiry time.Duration, purge time.Duration)
}

type store struct {
	db  *cache.Cache
	mtx *sync.Mutex

	// counters are used to keep track of how many times a key has been requested
	// it gives the flexibility to the clients to receive unique responses depending
	// upon the number of request againt the same key (command)
	storeCounter map[string]int
	getCounter   map[string]int
}

// New returns a new instance of store
func New(expiry time.Duration, purge time.Duration) Store {
	return &store{
		db:           cache.New(expiry, purge),
		storeCounter: make(map[string]int),
		getCounter:   make(map[string]int),
		mtx:          &sync.Mutex{},
	}
}

// Get retrieves from the store
// it uses get counter to keep track of how many times a key has ben requested
func (s *store) Get(key string) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	counter := 1
	if _, ok := s.getCounter[key]; ok {
		counter = s.getCounter[key]
	}

	getKey := key + "-" + fmt.Sprint(counter)
	data, found := s.db.Get(getKey)
	if found {
		s.getCounter[key] = counter + 1
		return []byte(data.(string)), nil
	}

	return nil, errors.New("not found")
}

/// Set stores in the store
// it uses storeCounter to keep track of how many times a key has ben stored
func (s *store) Set(key string, data []byte) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	counter := 1
	if _, ok := s.storeCounter[key]; ok {
		counter = s.storeCounter[key]
	}

	storeKey := key + "-" + fmt.Sprint(counter)

	s.storeCounter[key] = counter + 1
	s.db.Set(storeKey, string(data), cache.DefaultExpiration)

	return nil
}

// Reset clears up the store
func (s *store) Reset(expiry time.Duration, purge time.Duration) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.db = cache.New(expiry, purge)
	s.storeCounter = make(map[string]int)
	s.getCounter = make(map[string]int)
}
