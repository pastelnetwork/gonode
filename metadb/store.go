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
	// AllowedSubStatement - the supported count of sub statements
	AllowedSubStatement = 1
)

// Write execute a statement, not support multple statements
func (s *service) Write(_ context.Context, statement string, params ...interface{}) (*WriteResult, error) {
	if len(statement) == 0 {
		return nil, errors.New("statement are empty")
	}
	if len(strings.Split(statement, ";")) > AllowedSubStatement {
		return nil, errors.New("only support one sub statement")
	}

	// validate the parameters
	parameters := []*command.Parameter{}
	for _, parameter := range params {
		var item *command.Parameter
		switch v := parameter.(type) {
		case bool:
			item = &command.Parameter{Value: &command.Parameter_B{B: v}}
		case int64:
			item = &command.Parameter{Value: &command.Parameter_I{I: v}}
		case float64:
			item = &command.Parameter{Value: &command.Parameter_D{D: v}}
		case string:
			item = &command.Parameter{Value: &command.Parameter_S{S: v}}
		case []uint8:
			item = &command.Parameter{Value: &command.Parameter_Y{Y: v}}
		default:
			return nil, errors.Errorf("parameter type is not supported: %T", v)
		}
		parameters = append(parameters, item)
	}

	// prepare the execute command request
	request := &command.ExecuteRequest{
		Request: &command.Request{
			Transaction: true,
			Statements: []*command.Statement{
				{
					Sql:        statement,
					Parameters: parameters,
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
	first := results[0]
	if first.Error != "" {
		return nil, errors.Errorf("execute statement: %s", first.Error)
	}

	return &WriteResult{
		Timing:       first.Time,
		RowsAffected: first.RowsAffected,
		LastInsertID: first.LastInsertID,
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
	if first.Error != "" {
		return nil, errors.Errorf("query statement: %s", first.Error)
	}

	return &QueryResult{
		columns:   first.Columns,
		types:     first.Types,
		values:    first.Values,
		rowNumber: -1,
		Timing:    first.Time,
	}, nil
}
