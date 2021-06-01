package rqlite

import (
	"context"
	"errors"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/rqlite/command"
	"github.com/pastelnetwork/gonode/rqlite/db"
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

// Execute the multiple statements
func (s *Service) Execute(ctx context.Context, sqls []string) ([]*db.Result, error) {
	if len(sqls) == 0 {
		return nil, errors.New("sqls are empty")
	}

	// prepare the statements for command
	stmts := []*command.Statement{}
	for _, sql := range sqls {
		stmts = append(stmts, &command.Statement{
			Sql: sql,
		})
	}
	// prepare the command request
	er := &command.ExecuteRequest{
		Request: &command.Request{
			Transaction: false,
			Statements:  stmts,
		},
		Timings: false,
	}

	// execute the command by store
	results, err := s.db.Execute(er)
	if err != nil {
		log.WithContext(ctx).Errorf("store execute %s: %v", sqls, err)
		return nil, err
	}
	return results, err
}
