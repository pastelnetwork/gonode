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
func (s *service) WaitForStarting(ctx context.Context) error {
	select {
	case <-s.ready:
		return nil
	case <-ctx.Done():
		return errors.Errorf("context done: %w", ctx.Err())
	}
}

// LeaderAddr returns the address of the current leader. Returns a
// blank string if there is no leader.
func (s *service) LeaderAddress() string {
	address, _ := s.db.LeaderAddr()
	return address
}

// IsLeader let us know if this instance is leader or not
func (s *service) IsLeader() bool {
	return s.db.IsLeader()
}

// Stats return status of store
func (s *service) Stats(_ context.Context) (map[string]interface{}, error) {
	return s.db.Stats()
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
		LastInsertID: firsts.LastInsertId,
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

	vals := [][]interface{}{}
	for _, v := range first.Values {
		params := make([]interface{}, len(v.Parameters))

		for j, k := range v.Parameters {
			switch w := k.GetValue().(type) {
			case *command.Parameter_I:
				params[j] = w.I
			case *command.Parameter_D:
				params[j] = w.D
			case *command.Parameter_B:
				params[j] = w.B
			case *command.Parameter_Y:
				params[j] = w.Y
			case *command.Parameter_S:
				params[j] = w.S
			case nil:
				params[j] = nil
			default:
				return nil, errors.Errorf("query: unsupported type: %w", w)
			}

		}

		vals = append(vals, params)
	}

	return &QueryResult{
		columns:   first.Columns,
		types:     first.Types,
		values:    vals,
		rowNumber: -1,
		Timing:    first.Time,
	}, nil
}
