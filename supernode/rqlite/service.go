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
	config *Config
	db     *store.Store
}

// New returns a new service for rqlite cluster
func NewService(config *Config) *Service {
	return &Service{
		config: config,
	}
}

// Run starts the rqlite server
func (s *Service) Run(ctx context.Context) error {
	return s.startServer(ctx)
}
