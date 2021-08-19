package metadb

import (
	"context"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/metadb/rqlite/command"
)

const (
	// ReadLevelNono - with none, the node simply queries its local SQLite database, and does not even check if it is the leader
	ReadLevelNono = "nono"
	// ReadLevelWeak - weak instructs the node to check that it is the leader, before querying the local SQLite file
	ReadLevelWeak = "weak"
	// ReadLevelStrong - to avoid even the issues associated with weak consistency, rqlite also offers strong
	ReadLevelStrong = "strong"
)

// WaitForStarting wait for the db to completely started, just run once after starting up the db
func (s *service) WaitForStarting() {
	<-s.ready
}

// LeaderAddr returns the address of the current leader. Returns a
// blank string if there is no leader.
func (s *service) LeaderAddress() string {
	address, _ := s.db.LeaderAddr()
	return address
}

// IsLeader let us know if this instance is leader or not
func (s *service) IsLeader() bool {
	address, _ := s.db.LeaderAddr()
	return address == s.db.Addr()
}

// EnableFKConstraints is used to enable foreign key constraint
func (s *service) EnableFKConstraints(e bool) error {
	return s.db.EnableFKConstraints(e)
}

// Write execute a statement, not support multple statements
func (s *service) Write(_ context.Context, statement string) (*WriteResult, error) {
	if len(statement) == 0 {
		return nil, errors.New("statement are empty")
	}
	if !s.IsLeader() {
		return nil, errors.New("not a leader")
	}

	// prepare the execute command request
	request := &command.ExecuteRequest{
		Request: &command.Request{
			Transaction: true,
			Statements: []*command.Statement{
				{
					Sql: statement,
				},
			},
		},
		Timings: true,
	}

	// execute the command by store
	results, err := s.db.Execute(request)
	if err != nil {
		return nil, errors.Errorf("execute statement: %w", err)
	}
	firsts := results[0]

	return &WriteResult{
		Error:        firsts.Error,
		Timing:       firsts.Time,
		RowsAffected: firsts.RowsAffected,
		LastInsertID: firsts.LastInsertID,
	}, nil
}

// for the read consistency level, Weak is probably sufficient for most applications, and is the default read consistency level
func (s *service) consistencyLevel(level string) command.QueryRequest_Level {
	switch strings.ToLower(level) {
	case ReadLevelNono:
		return command.QueryRequest_QUERY_REQUEST_LEVEL_NONE
	case ReadLevelWeak:
		return command.QueryRequest_QUERY_REQUEST_LEVEL_WEAK
	case ReadLevelStrong:
		return command.QueryRequest_QUERY_REQUEST_LEVEL_STRONG
	default:
		return command.QueryRequest_QUERY_REQUEST_LEVEL_WEAK
	}
}

// Query execute a query, not support multple statements
// level can be 'none', 'weak', and 'strong'
func (s *service) Query(_ context.Context, statement string, level string) (*QueryResult, error) {
	if len(statement) == 0 {
		return nil, errors.New("statement are empty")
	}

	// prepare the query command request
	request := &command.QueryRequest{
		Request: &command.Request{
			Transaction: true,
			Statements: []*command.Statement{
				{
					Sql: statement,
				},
			},
		},
		Timings: true,
		Level:   s.consistencyLevel(level),
	}

	// query the command by store
	rows, err := s.db.Query(request)
	if err != nil {
		return nil, errors.Errorf("query statement: %w", err)
	}
	first := rows[0]

	return &QueryResult{
		columns:   first.Columns,
		types:     first.Types,
		values:    first.Values,
		rowNumber: -1,
		Timing:    first.Time,
	}, nil
}
