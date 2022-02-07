package metadb

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "metadb"
)

// MetaDB represents the metadb service
type MetaDB interface {
	// Run starts the rqlite service
	Run(ctx context.Context) error
	// Query execute a query, not support multple statements
	// level can be 'none', 'weak', and 'strong'
	Query(ctx context.Context, statement string, level string) (*QueryResult, error)
	// Write execute a statement, not support multple statements
	Write(ctx context.Context, statement string) (*WriteResult, error)
	// WaitForStarting wait for the db to completely started, just run once after starting up the db
	WaitForStarting(ctx context.Context) error
	// LeaderAddr returns the address of the current leader. Returns a
	// blank string if there is no leader.
	LeaderAddress() string
	// IsLeader let us know if this instance is leader or not
	IsLeader() bool
	// Stats return status of MetaDB
	Stats(ctx context.Context) (map[string]interface{}, error)
}

type service struct {
	nodeID            string        // the node id for rqlite cluster
	config            *Config       // the service configuration
	db                *store.Store  // the store for accessing the rqlite cluster
	ready             chan struct{} // mark the rqlite node is started
	nodeIPList        []string
	pastelClient      pastel.Client
	currentBlockCount int32
}

// New returns a new service for metadb
func New(config *Config, nodeID string, pastelClient pastel.Client) MetaDB {
	return &service{
		nodeID:       nodeID,
		config:       config,
		ready:        make(chan struct{}, 5),
		pastelClient: pastelClient,
	}
}

// Run starts the rqlite server
func (s *service) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := s.startServer(ctx); err != nil {
				if utils.IsContextErr(err) {
					return err
				}

				//log.MetaDB().WithContext(ctx).WithError(err).Error("failed to run metadb, retrying.")
			} else {
				return nil
			}
		}
	}
}
