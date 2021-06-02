package rqlite

import (
	"context"
	"strings"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/rqlite/command"
	"github.com/pastelnetwork/gonode/rqlite/db"
)

const (
	// ReadLevelNono - with none, the node simply queries its local SQLite database, and does not even check if it is the leader
	ReadLevelNono = "nono"
	// ReadLevelWeek - weak instructs the node to check that it is the leader, before querying the local SQLite file
	ReadLevelWeek = "week"
	// ReadLevelStrong - to avoid even the issues associated with weak consistency, rqlite also offers strong
	ReadLevelStrong = "strong"
)

// Execute executes a slice of queries, each of which is not expected
// to return rows. If timings is true, then timing information will
// be return. If tx is true, then either all queries will be executed
// successfully or it will as though none executed.
func (s *Service) Execute(ctx context.Context, stmts []string, tx bool, timings bool) ([]*db.Result, error) {
	if len(stmts) == 0 {
		return nil, errors.New("execute statements are empty")
	}

	// prepare the statements for command
	statements := []*command.Statement{}
	for _, stmt := range stmts {
		statements = append(statements, &command.Statement{
			Sql: stmt,
		})
	}
	// prepare the execute command request
	request := &command.ExecuteRequest{
		Request: &command.Request{
			Transaction: tx,
			Statements:  statements,
		},
		Timings: timings,
	}

	// execute the command by store
	results, err := s.db.Execute(request)
	if err != nil {
		return nil, errors.Errorf("store execute statements: %v", err)
	}

	return results, nil
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

// Query executes a slice of queries, each of which returns rows. If
// timings is true, then timing information will be returned. If tx
// is true, then all queries will take place while a read transaction
// is held on the database. The level can be 'none', 'weak', and 'strong'
func (s *Service) Query(ctx context.Context, stmts []string, tx bool, timings bool, level string) ([]*db.Rows, error) {
	if len(stmts) == 0 {
		return nil, errors.New("query statements are empty")
	}

	// prepare the statements for command
	statements := []*command.Statement{}
	for _, stmt := range stmts {
		statements = append(statements, &command.Statement{
			Sql: stmt,
		})
	}

	// prepare the query command request
	request := &command.QueryRequest{
		Request: &command.Request{
			Transaction: tx,
			Statements:  statements,
		},
		Timings: timings,
		Level:   s.consistencyLevel(level),
	}

	// query the command by store
	results, err := s.db.Query(request)
	if err != nil {
		return nil, errors.Errorf("store query statements: %v", err)
	}

	return results, nil
}
