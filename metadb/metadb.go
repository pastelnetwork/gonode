package metadb

import (
	"context"

	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
)

const (
	logPrefix = "metadb"
)

// MetaDB represents the metadb service
type MetaDB interface {
	Client

	// Run starts the rqlite service
	Run(ctx context.Context) error
	// AfterFunc updates the after function for rqlite service
	AfterFunc(fn func() error)
}

type service struct {
	nodeID    string        // the node id for rqlite cluster
	config    *Config       // the service configuration
	db        *store.Store  // the store for accessing the rqlite cluster
	ready     chan struct{} // mark the rqlite node is started
	afterFunc func() error  // call the function after the rqlite cluster is ready
}

// New returns a new service for metadb
func New(config *Config, nodeID string) MetaDB {
	return &service{
		nodeID: nodeID,
		config: config,
		ready:  make(chan struct{}, 5),
	}
}

// Run starts the rqlite server
func (s *service) Run(ctx context.Context) error {
	return s.startServer(ctx)
}

// AfterFunc updates the after function for rqlite service
func (s *service) AfterFunc(fn func() error) {
	s.afterFunc = fn
}
