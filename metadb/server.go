package metadb

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/metadb/rqlite/cluster"
	httpd "github.com/pastelnetwork/gonode/metadb/rqlite/http"
	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
	"github.com/pastelnetwork/gonode/metadb/rqlite/tcp"
)

const (
	defaultJoinAttempts             = 5
	defaultJoinInterval             = 5 * time.Second
	defaultRaftHeartbeatTimeout     = time.Second
	defaultRaftElectionTimeout      = time.Second
	defaultRaftApplyTimeout         = 10 * time.Second
	defaultRaftOpenTimeout          = 60 * time.Second
	defaultRaftWaitForLeader        = true
	defaultRaftSnapThreshold        = 8192
	defaultRaftSnapInterval         = 30 * time.Second
	defaultRaftLeaderLeaseTimeout   = 0 * time.Second
	defaultCompressionSize          = 150
	defaultCompressionBatch         = 5
	defaultCheckLeaderInterval      = 30 * time.Second
	defaultCheckBlockCountInterval  = 30 * time.Second
	defaultJoinClusterRetryInterval = 10 * time.Second
)

// wait until the store is in full consensus
func (s *service) waitForConsensus(ctx context.Context, dbStore *store.Store) error {
	if _, err := dbStore.WaitForLeader(ctx, defaultRaftOpenTimeout); err != nil {
		if defaultRaftWaitForLeader {
			return errors.Errorf("leader did not appear within timeout: %w", err)
		}
		log.MetaDB().WithContext(ctx).Infof("ignoring error while waiting for leader")
	}
	if err := dbStore.WaitForApplied(ctx, defaultRaftOpenTimeout); err != nil {
		return errors.Errorf("store log not applied within timeout: %w", err)
	}

	return nil
}

// start the http server
func (s *service) startHTTPServer(ctx context.Context, dbStore *store.Store, cs *cluster.Service) error {
	httpAddr := fmt.Sprintf("%s:%d", s.config.ListenAddress, s.config.HTTPPort)

	//clstrDialer := tcp.NewDialer(cluster.MuxClusterHeader, nodeEncrypt, noNodeVerify)
	//clstrClient := cluster.NewClient(clstrDialer)
	//if err := clstrClient.SetLocal(raftAdv, clstr); err != nil {
	//	log.Fatalf("failed to set cluster client local parameters: %s", err.Error())
	//}

	// create http server
	server := httpd.New(ctx, httpAddr, dbStore, cs, nil)

	// start the http server
	return server.Start()
}

// start a mux for rqlite node
func (s *service) startNodeMux(ctx context.Context, ln net.Listener, raftAdvAddr string) (*tcp.Mux, error) {
	var adv net.Addr
	var err error
	if raftAdvAddr != "" {
		adv, err = net.ResolveTCPAddr("tcp", raftAdvAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve advertise address %s: %s", raftAdvAddr, err.Error())
		}
	}

	mux, err := tcp.NewMux(ctx, ln, adv)
	if err != nil {
		return nil, errors.Errorf("create node-to-node mux: %w", err)
	}

	go mux.Serve()

	return mux, nil
}

// start the cluster server
func (s *service) startClusterService(ctx context.Context, tn cluster.Transport, httpAdvAddr string) (*cluster.Service, error) {
	c := cluster.New(ctx, tn, s.db)

	// set the api address
	c.SetAPIAddr(httpAdvAddr)

	// open the cluster service
	if err := c.Open(); err != nil {
		return nil, err
	}
	return c, nil
}

// create and open the store of rqlite cluster
func (s *service) initStore(ctx context.Context, raftTn *tcp.Layer, externalAddress string) (*store.Store, error) {
	// create and open the store, which is on disk
	dbConf := store.NewDBConfig(false)
	db := store.New(ctx, raftTn, &store.Config{
		DBConf: dbConf,
		Dir:    s.config.DataDir,
		ID:     s.nodeID,
	})

	// set optional parameters on store
	db.SetRequestCompression(defaultCompressionBatch, defaultCompressionSize)
	db.ShutdownOnRemove = false
	db.SnapshotThreshold = defaultRaftSnapThreshold
	db.SnapshotInterval = defaultRaftSnapInterval
	db.LeaderLeaseTimeout = defaultRaftLeaderLeaseTimeout
	db.HeartbeatTimeout = defaultRaftHeartbeatTimeout
	db.ElectionTimeout = defaultRaftElectionTimeout
	db.ApplyTimeout = defaultRaftApplyTimeout

	// a pre-existing node
	bootstrap := false
	isNew := store.IsNewNode(s.config.DataDir)
	if isNew {
		bootstrap = true // new node, it needs to bootstrap
	} else {
		log.MetaDB().WithContext(ctx).Infof("node is detected in: %v", s.config.DataDir)
	}

	httpAdvAddr := fmt.Sprintf("%s:%d", externalAddress, s.config.HTTPPort)

	var joinIPAddresses []string
	if !s.config.IsLeader {
		if s.config.LeaderAddress != "" {
			joinIPAddresses = append(joinIPAddresses, s.config.LeaderAddress)
		} else {
			for _, ip := range s.nodeIPList {
				if httpAdvAddr == ip {
					continue
				}
				joinIPAddresses = append(joinIPAddresses, ip)
			}
		}
	}

	if len(joinIPAddresses) > 0 {
		bootstrap = false
		log.MetaDB().WithContext(ctx).Info("join addresses specified, node is not bootstrap")
	} else {
		log.MetaDB().WithContext(ctx).Info("no join addresses")
	}
	// join address supplied, but we don't need them
	if !isNew && len(joinIPAddresses) > 0 {
		log.MetaDB().WithContext(ctx).Info("node is already member of cluster")
	}

	// open store
	if err := db.Open(bootstrap); err != nil {
		return nil, errors.Errorf("open store: %w", err)
	}
	s.db = db

	// execute any requested join operation
	if len(joinIPAddresses) == 0 {
		log.MetaDB().WithContext(ctx).Warn("No nearby nodes found. Not joining cluster")
		return db, nil
	}

	raftAdvAddr := fmt.Sprintf("%s:%d", externalAddress, s.config.RaftPort)

	if err := s.joinCluster(ctx, joinIPAddresses, raftAdvAddr); err == nil {
		return db, nil
	}

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-time.After(defaultJoinClusterRetryInterval):
			err := s.joinCluster(ctx, joinIPAddresses, raftAdvAddr)
			if err == nil {
				break loop
			}

			log.MetaDB().WithContext(ctx).WithError(err).Errorf("metdadb join cluster failure, retrying in %v s", defaultJoinClusterRetryInterval)
		}
	}

	return db, nil
}

func (s *service) joinCluster(ctx context.Context, joinAddrs []string, raftAdvAddr string) (err error) {

	log.MetaDB().WithContext(ctx).Infof("join addresses are: %v", joinAddrs)

	// join rqlite cluster
	joinAddr, err := cluster.Join(ctx, "", joinAddrs, s.db.ID(), raftAdvAddr, !s.config.NoneVoter,
		defaultJoinAttempts, defaultJoinInterval, nil)
	if err != nil {
		return fmt.Errorf("join cluster at %v: %s", joinAddrs, err.Error())
	}

	log.MetaDB().WithContext(ctx).Infof("successfully joined cluster at %v", joinAddr)
	return nil
}

// start the rqlite server, and try to join rqlite cluster if the join addresses is not empty
func (s *service) startServer(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)
	if err := s.setDatabaseNodes(ctx); err != nil {
		return errors.Errorf("set database nodes failure: %w", err)
	}

	externalAddress := s.config.ListenAddress
	if externalAddress == "0.0.0.0" {
		var err error
		if externalAddress, err = utils.GetExternalIPAddress(); err != nil {
			return fmt.Errorf("cannot find own external address: %s. Replace 'listen_address: 0.0.0.0' in config for real IP", err.Error())
		}
	}

	raftListenAddr := fmt.Sprintf("%s:%d", s.config.ListenAddress, s.config.RaftPort)
	raftAdvAddr := fmt.Sprintf("%s:%d", externalAddress, s.config.RaftPort)
	httpAdvAddr := fmt.Sprintf("%s:%d", externalAddress, s.config.HTTPPort)

	// create internode network mux and configure.
	muxListener, err := net.Listen("tcp", raftListenAddr)
	if err != nil {
		return errors.Errorf("listen on %s: %w", raftListenAddr, err)
	}
	mux, err := s.startNodeMux(ctx, muxListener, raftAdvAddr)
	if err != nil {
		return errors.Errorf("start node mux: %w", err)
	}
	raftTn := mux.Listen(cluster.MuxRaftHeader)

	// create cluster service, so nodes can learn information about each other.
	// This can be started now since it doesn't require a functioning Store yet.
	cs, err := s.startClusterService(ctx, mux.Listen(cluster.MuxClusterHeader), httpAdvAddr)
	if err != nil {
		return errors.Errorf("start create cluster service: %w", err)
	}

	// create and open the store
	db, err := s.initStore(ctx, raftTn, externalAddress)
	if err != nil {
		return errors.Errorf("create and open store: %w", err)
	}
	s.db = db

	go s.initLeadershipTransferTrigger(ctx, defaultCheckLeaderInterval,
		defaultCheckBlockCountInterval)

	// wait until the store is in full consensus
	if err := s.waitForConsensus(ctx, db); err != nil {
		return errors.Errorf("wait for consensus: %w", err)
	}
	log.MetaDB().WithContext(ctx).Info("store has reached consensus")

	// start the HTTP API server
	if err := s.startHTTPServer(ctx, db, cs); err != nil {
		return errors.Errorf("start http server: %w", err)
	}
	// mark the rqlite node is ready
	select {
	case <-ctx.Done():
		return errors.Errorf("context done: %w", ctx.Err())
	case s.ready <- struct{}{}:
		// do nothing, continue
	}

	log.MetaDB().WithContext(ctx).Info("metadb service is started")

	// block until context is done
	<-ctx.Done()

	// close the rqlite store
	if err := db.Close(true); err != nil {
		log.MetaDB().WithContext(ctx).WithError(err).Errorf("close store failed")
	}

	// close the mux listener
	if err := muxListener.Close(); err != nil {
		log.MetaDB().WithContext(ctx).WithError(err).Errorf("Close mux listener failed")
	}

	log.MetaDB().WithContext(ctx).Info("metadb service is stopped")
	return nil
}

func (s *service) setDatabaseNodes(ctx context.Context) error {
	var nodeIPList []string
	nodeList, err := s.pastelClient.MasterNodesExtra(ctx)
	if err != nil {
		return err
	}

	for _, nodeInfo := range nodeList {
		address := nodeInfo.ExtAddress
		segments := strings.Split(address, ":")
		if len(segments) != 2 {
			return errors.Errorf("malformed db node address: %s", address)
		}
		nodeAddress := fmt.Sprintf("%s:%d", segments[0], s.config.HTTPPort)
		nodeIPList = append(nodeIPList, nodeAddress)
	}
	s.nodeIPList = nodeIPList

	return nil
}
