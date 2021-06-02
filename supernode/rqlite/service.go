package rqlite

import (
	"context"

	"github.com/pastelnetwork/gonode/rqlite/store"
)

const (
	logPrefix = "rqlite"
)

// Service represents the rqlite cluster
type Service struct {
	config *Config       // the service configuration
	db     *store.Store  // the store for accessing the rqlite cluster
	ready  chan struct{} // mark the rqlite node is started
}

// New returns a new service for rqlite cluster
func NewService(config *Config) *Service {
	return &Service{
		config: config,
		ready:  make(chan struct{}),
	}
}

// Run starts the rqlite server
func (s *Service) Run(ctx context.Context) error {
	return s.startServer(ctx)
}
