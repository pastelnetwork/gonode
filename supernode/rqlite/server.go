package rqlite

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/rqlite/cluster"
	"github.com/pastelnetwork/gonode/rqlite/disco"
	httpd "github.com/pastelnetwork/gonode/rqlite/http"
	"github.com/pastelnetwork/gonode/rqlite/store"
	"github.com/pastelnetwork/gonode/rqlite/tcp"
)

// the order is: node id, raft advertise address vs raft address
func (s *Service) idOrRaftAddr() string {
	if s.config.NodeID != "" {
		return s.config.NodeID
	}
	if s.config.RaftAdvertiseAddress == "" {
		return s.config.RaftAddress
	}
	return s.config.RaftAdvertiseAddress
}

// determine the join addresses
func (s *Service) determineJoinAddresses(ctx context.Context) ([]string, error) {
	apiAdv := s.config.HTTPAddress
	if s.config.HTTPAdvertiseAddress != "" {
		apiAdv = s.config.HTTPAdvertiseAddress
	}

	var addrs []string
	if s.config.JoinAddress != "" {
		// explicit join addresses are first priority.
		addrs = strings.Split(s.config.JoinAddress, ",")
	}

	if s.config.DiscoveryID != "" {
		log.WithContext(ctx).Infof("register with discovery service at %s with ID %s", s.config.DiscoveryURL, s.config.DiscoveryID)

		c := disco.New(s.config.DiscoveryURL)
		r, err := c.Register(s.config.DiscoveryID, apiAdv)
		if err != nil {
			return nil, err
		}
		log.WithContext(ctx).Infof("discovery service responded with nodes:", r.Nodes)

		for _, a := range r.Nodes {
			if a != apiAdv {
				addrs = append(addrs, a)
			}
		}
	}

	return addrs, nil
}

// wait until the store is in full consensus
func (s *Service) waitForConsensus(ctx context.Context, dbStore *store.Store) error {
	openTimeout, err := time.ParseDuration(s.config.RaftOpenTimeout)
	if err != nil {
		return fmt.Errorf("parse RaftOpenTimeout: %v", err)
	}
	if _, err := dbStore.WaitForLeader(openTimeout); err != nil {
		if s.config.RaftWaitForLeader {
			return fmt.Errorf("leader did not appear within timeout: %v", err)
		}
		log.WithContext(ctx).Infof("ignoring error while waiting for leader")
	}
	if openTimeout != 0 {
		if err := dbStore.WaitForApplied(openTimeout); err != nil {
			return fmt.Errorf("log was not fully applied within timeout: %s", err.Error())
		}
	} else {
		log.WithContext(ctx).Info("not waiting for logs to be applied")
	}
	return nil
}

// start the http server
func (s *Service) startHTTPServer(ctx context.Context, dbStore *store.Store) error {
	logger := log.DefaultLogger.WithField("prefix", "http")

	// new http server
	server := httpd.New(s.config.HTTPAddress, dbStore, nil, logger)

	// start the http server
	return server.Start()
}

// start the rqlite server, and try to join rqlite cluster if the join addresses is not empty
func (s *Service) startServer(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	// create internode network layer
	transport := tcp.NewTransport()
	if err := transport.Open(s.config.RaftAddress); err != nil {
		log.WithContext(ctx).Errorf("open internode network layer: %v", err)
		return err
	}

	// create and open the store
	dbConf := store.NewDBConfig(s.config.DNS, !s.config.OnDisk)
	db := store.New(transport, &store.StoreConfig{
		DBConf: dbConf,
		Dir:    s.config.DataDir,
		ID:     s.idOrRaftAddr(),
		Logger: log.DefaultLogger.WithField("prefix", "store"),
	})

	var err error
	// set optional parameters on store
	db.SetRequestCompression(s.config.CompressionBatch, s.config.CompressionSize)
	db.RaftLogLevel = s.config.RaftLogLevel
	db.ShutdownOnRemove = s.config.RaftShutdownOnRemove
	db.SnapshotThreshold = s.config.RaftSnapThreshold
	db.SnapshotInterval, err = time.ParseDuration(s.config.RaftSnapInterval)
	if err != nil {
		log.WithContext(ctx).Errorf("parse RaftSnapInterval: %v", err)
		return err
	}
	db.LeaderLeaseTimeout, err = time.ParseDuration(s.config.RaftLeaderLeaseTimeout)
	if err != nil {
		log.WithContext(ctx).Errorf("parse RaftLeaderLeaseTimeout: %v", err)
		return err
	}
	db.HeartbeatTimeout, err = time.ParseDuration(s.config.RaftHeartbeatTimeout)
	if err != nil {
		log.WithContext(ctx).Errorf("parse RaftHeartbeatTimeout: %v", err)
		return err
	}
	db.ElectionTimeout, err = time.ParseDuration(s.config.RaftElectionTimeout)
	if err != nil {
		log.WithContext(ctx).Errorf("parse RaftElectionTimeout: %v", err)
		return err
	}
	db.ApplyTimeout, err = time.ParseDuration(s.config.RaftApplyTimeout)
	if err != nil {
		log.WithContext(ctx).Errorf("parse RaftApplyTimeout: %v", err)
		return err
	}

	// a pre-existing node
	bootstrap := false
	isNew := store.IsNewNode(s.config.DataDir)
	if isNew {
		bootstrap = true // new node, it needs to bootstrap
	} else {
		log.WithContext(ctx).Infof("preexisting node detected in: %v", s.config.DataDir)
	}

	// determine the join addresses
	joins, err := s.determineJoinAddresses(ctx)
	if err != nil {
		log.WithContext(ctx).Errorf("determine join addresses: %v", err)
		return err
	}
	// supplying join addresses means bootstrapping a new cluster won't be required.
	if len(joins) > 0 {
		bootstrap = false
		log.WithContext(ctx).Info("join addresses specified, node is not bootstrapping")
	} else {
		log.WithContext(ctx).Info("no join addresses")
	}
	// join address supplied, but we don't need them
	if !isNew && len(joins) > 0 {
		log.WithContext(ctx).Info("node is already member of cluster, ignoring join addresses")
	}

	// open store
	if err := db.Open(bootstrap); err != nil {
		log.WithContext(ctx).Errorf("open store: %v", err)
		return err
	}
	s.db = db

	// prepare metadata for join command.
	apiAdv := s.config.HTTPAddress
	if s.config.HTTPAdvertiseAddress != "" {
		apiAdv = s.config.HTTPAdvertiseAddress
	}
	apiProto := "http"
	meta := map[string]string{
		"api_addr":  apiAdv,
		"api_proto": apiProto,
	}
	// execute any requested join operation
	if len(joins) > 0 && isNew {
		log.WithContext(ctx).Infof("join addresses are: %v", joins)
		advAddr := s.config.RaftAddress
		if s.config.RaftAdvertiseAddress != "" {
			advAddr = s.config.RaftAdvertiseAddress
		}

		joinDuration, err := time.ParseDuration(s.config.JoinInterval)
		if err != nil {
			log.WithContext(ctx).Errorf("parse JoinInterval: %v", err)
			return err
		}

		// join rqlite cluster
		if j, err := cluster.Join(
			s.config.JoinSourceIP,
			joins,
			db.ID(),
			advAddr,
			!s.config.RaftNotVoter,
			meta,
			s.config.JoinAttempts,
			joinDuration,
			nil,
		); err != nil {
			log.WithContext(ctx).Errorf("join cluster at %v: %v", joins, err)
			return err
		} else {
			log.WithContext(ctx).Infof("successfully joined cluster at %v", j)
		}
	}

	// wait until the store is in full consensus
	if err := s.waitForConsensus(ctx, db); err != nil {
		log.WithContext(ctx).Errorf("wait for consensus: %v", err)
		return err
	}

	// this may be a standalone server. In that case set its own metadata.
	if err := db.SetMetadata(meta); err != nil && err != store.ErrNotLeader {
		// Non-leader errors are OK, since metadata will then be set through
		// consensus as a result of a join. All other errors indicate a problem.
		log.WithContext(ctx).Errorf("set store metadata: %v", err)
		return err
	}

	// start the HTTP API server
	if err := s.startHTTPServer(ctx, db); err != nil {
		log.WithContext(ctx).Errorf("start http server: %v", err)
		return err
	}
	log.WithContext(ctx).Info("node is ready, block until context is done")

	// block until context is done
	<-ctx.Done()

	// close the rqlite server
	if err := db.Close(true); err != nil {
		log.WithContext(ctx).Errorf("close store: %v", err)
	}

	log.WithContext(ctx).Info("rqlite server is stopped")
	return nil
}
