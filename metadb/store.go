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
	// ReadLevelWeek - weak instructs the node to check that it is the leader, before querying the local SQLite file
	ReadLevelWeek = "week"
	// ReadLevelStrong - to avoid even the issues associated with weak consistency, rqlite also offers strong
	ReadLevelStrong = "strong"
	// AllowedSubStatement - the supported count of sub statements
	AllowedSubStatement = 1
)

// Write execute a statement, not support multple statements
func (s *Service) Write(_ context.Context, statement string) (*WriteResult, error) {
	if len(statement) == 0 {
		return nil, errors.New("statement are empty")
	}
	if len(strings.Split(statement, ";")) > AllowedSubStatement {
		return nil, errors.New("only support one sub statement")
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
func (s *Service) consistencyLevel(level string) command.QueryRequest_Level {
	switch strings.ToLower(level) {
	case "none":
		return command.QueryRequest_QUERY_REQUEST_LEVEL_NONE
	case "weak":
		return command.QueryRequest_QUERY_REQUEST_LEVEL_WEAK
	case "strong":
		return command.QueryRequest_QUERY_REQUEST_LEVEL_STRONG
	default:
		return command.QueryRequest_QUERY_REQUEST_LEVEL_WEAK
	}
}

// Query execute a query, not support multple statements
// level can be 'none', 'weak', and 'strong'
func (s *Service) Query(_ context.Context, statement string, level string) (*QueryResult, error) {
	if len(statement) == 0 {
		return nil, errors.New("statement are empty")
	}
	if len(strings.Split(statement, ";")) > AllowedSubStatement {
		return nil, errors.New("only support one sub statement")
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
