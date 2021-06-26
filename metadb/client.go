package metadb

import "context"

// Client exposes the interfaces for metadb service
type Client interface {
	// Query execute a query, not support multple statements
	Query(ctx context.Context, statement string, level string) (*QueryResult, error)
	// Write execute a statement with parameters, not support multple statements
	Write(ctx context.Context, statement string, params ...interface{}) (*WriteResult, error)
}
