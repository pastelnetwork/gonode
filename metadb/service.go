package metadb

import (
	"context"

	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
)

const (
	logPrefix = "metadb"
)

// Service represents the rqlite cluster
type Service struct {
	config *Config       // the service configuration
	db     *store.Store  // the store for accessing the rqlite cluster
	ready  chan struct{} // mark the rqlite node is started
}

// NewService returns a new service for rqlite cluster
func NewService(config *Config) *Service {
	return &Service{
		config: config,
		ready:  make(chan struct{}, 1),
	}
}

// Run starts the rqlite server
func (s *Service) Run(ctx context.Context) error {
	return s.startServer(ctx)
}
