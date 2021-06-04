package metadb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/metadb/rqlite/auth"
	"github.com/pastelnetwork/gonode/metadb/rqlite/cluster"
	"github.com/pastelnetwork/gonode/metadb/rqlite/disco"
	httpd "github.com/pastelnetwork/gonode/metadb/rqlite/http"
	"github.com/pastelnetwork/gonode/metadb/rqlite/store"
	"github.com/pastelnetwork/gonode/metadb/rqlite/tcp"
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

		c := disco.New(ctx, s.config.DiscoveryURL)
		r, err := c.Register(s.config.DiscoveryID, apiAdv)
		if err != nil {
			return nil, errors.Errorf("discovery register: %w", err)
		}
		log.WithContext(ctx).Infof("discovery service responded with nodes: %v", r.Nodes)

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
		return errors.Errorf("parse RaftOpenTimeout: %w", err)
	}
	if _, err := dbStore.WaitForLeader(openTimeout); err != nil {
		if s.config.RaftWaitForLeader {
			return errors.Errorf("leader did not appear within timeout: %w", err)
		}
		log.WithContext(ctx).Infof("ignoring error while waiting for leader")
	}
	if openTimeout != 0 {
		if err := dbStore.WaitForApplied(openTimeout); err != nil {
			return errors.Errorf("store log not applied within timeout: %w", err)
		}
	} else {
		log.WithContext(ctx).Info("not waiting for logs to be applied")
	}

	return nil
}

// open the auth file, and returns a credential store instance
func (s *Service) credentialStore() (*auth.CredentialsStore, error) {
	if s.config.AuthFile == "" {
		return nil, nil
	}

	file, err := os.Open(s.config.AuthFile)
	if err != nil {
		return nil, errors.Errorf("open authentication file: %w", err)
	}

	store := auth.NewCredentialsStore()
	if store.Load(file); err != nil {
		return nil, errors.Errorf("store load: %w", err)
	}
	return store, nil
}

// start the http server
func (s *Service) startHTTPServer(ctx context.Context, dbStore *store.Store, cs *cluster.Service) error {
	// load the credential store
	cred, err := s.credentialStore()
	if err != nil {
		return errors.Errorf("load credentail store: %w", err)
	}

	var server *httpd.Service
	// create http server and load authentication information if required
	if cred != nil {
		server = httpd.New(ctx, s.config.HTTPAddress, dbStore, cs, cred)
	} else {
		server = httpd.New(ctx, s.config.HTTPAddress, dbStore, cs, nil)
	}
	server.CertFile = s.config.X509Cert
	server.KeyFile = s.config.X509Key
	server.TLS1011 = s.config.TLS1011

	// start the http server
	return server.Start()
}

// start a mux for rqlite node
func (s *Service) startNodeMux(ctx context.Context, ln net.Listener) (*tcp.Mux, error) {
	var adv net.Addr
	var err error
	if s.config.RaftAdvertiseAddress != "" {
		adv, err = net.ResolveTCPAddr("tcp", s.config.RaftAdvertiseAddress)
		if err != nil {
			return nil, errors.Errorf("resolve advertise address %s: %w", s.config.RaftAdvertiseAddress, err)
		}
	}

	var mux *tcp.Mux
	if s.config.NodeEncrypt {
		log.WithContext(ctx).Infof("enabling node-to-node encryption with cert: %s, key: %s", s.config.NodeX509Cert, s.config.NodeX509Key)
		mux, err = tcp.NewTLSMux(ctx, ln, adv, s.config.NodeX509Cert, s.config.NodeX509Key, s.config.NodeX509CACert)
	} else {
		mux, err = tcp.NewMux(ctx, ln, adv)
	}
	if err != nil {
		return nil, errors.Errorf("create node-to-node mux: %w", err)
	}
	mux.InsecureSkipVerify = s.config.NodeNoVerify

	go mux.Serve()

	return mux, nil
}

// start the cluster server
func (s *Service) startClusterService(ctx context.Context, tn cluster.Transport) (*cluster.Service, error) {
	c := cluster.New(ctx, tn)

	apiAddr := s.config.HTTPAddress
	if s.config.HTTPAdvertiseAddress != "" {
		apiAddr = s.config.HTTPAdvertiseAddress
	}
	c.SetAPIAddr(apiAddr)
	// conditions met for a HTTPS API
	c.EnableHTTPS(s.config.X509Cert != "" && s.config.X509Key != "")

	// open the cluster service
	if err := c.Open(); err != nil {
		return nil, err
	}
	return c, nil
}

// start the rqlite server, and try to join rqlite cluster if the join addresses is not empty
func (s *Service) startServer(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	// create internode network mux and configure.
	muxListener, err := net.Listen("tcp", s.config.RaftAddress)
	if err != nil {
		return errors.Errorf("listen on %s: %w", s.config.RaftAddress, err)
	}
	// close the mux listener
	defer func() {
		if err := muxListener.Close(); err != nil {
			log.WithContext(ctx).Errorf("close mux listener: %w", err)
		}
	}()

	mux, err := s.startNodeMux(ctx, muxListener)
	if err != nil {
		return errors.Errorf("start node mux: %w", err)
	}
	raftTn := mux.Listen(cluster.MuxRaftHeader)

	// create cluster service, so nodes can learn information about each other.
	// This can be started now since it doesn't require a functioning Store yet.
	cs, err := s.startClusterService(ctx, mux.Listen(cluster.MuxClusterHeader))
	if err != nil {
		return errors.Errorf("start create cluster service: %w", err)
	}

	// create and open the store
	dbConf := store.NewDBConfig(s.config.DNS, !s.config.OnDisk)
	db := store.New(ctx, raftTn, &store.Config{
		DBConf: dbConf,
		Dir:    s.config.DataDir,
		ID:     s.idOrRaftAddr(),
	})

	// set optional parameters on store
	db.SetRequestCompression(s.config.CompressionBatch, s.config.CompressionSize)
	db.RaftLogLevel = s.config.RaftLogLevel
	db.ShutdownOnRemove = s.config.RaftShutdownOnRemove
	db.SnapshotThreshold = s.config.RaftSnapThreshold
	db.SnapshotInterval, err = time.ParseDuration(s.config.RaftSnapInterval)
	if err != nil {
		return errors.Errorf("parse RaftSnapInterval: %w", err)
	}
	db.LeaderLeaseTimeout, err = time.ParseDuration(s.config.RaftLeaderLeaseTimeout)
	if err != nil {
		return errors.Errorf("parse RaftLeaderLeaseTimeout: %w", err)
	}
	db.HeartbeatTimeout, err = time.ParseDuration(s.config.RaftHeartbeatTimeout)
	if err != nil {
		return errors.Errorf("parse RaftHeartbeatTimeout: %w", err)
	}
	db.ElectionTimeout, err = time.ParseDuration(s.config.RaftElectionTimeout)
	if err != nil {
		return errors.Errorf("parse RaftElectionTimeout: %w", err)
	}
	db.ApplyTimeout, err = time.ParseDuration(s.config.RaftApplyTimeout)
	if err != nil {
		return errors.Errorf("parse RaftApplyTimeout: %w", err)
	}

	// a pre-existing node
	bootstrap := false
	isNew := store.IsNewNode(s.config.DataDir)
	if isNew {
		bootstrap = true // new node, it needs to bootstrap
	} else {
		log.WithContext(ctx).Infof("node is detected in: %v", s.config.DataDir)
	}

	// determine the join addresses
	joins, err := s.determineJoinAddresses(ctx)
	if err != nil {
		return errors.Errorf("determine join addresses: %w", err)
	}
	// supplying join addresses means bootstrapping a new cluster won't be required.
	if len(joins) > 0 {
		bootstrap = false
		log.WithContext(ctx).Info("join addresses specified, node is not bootstrap")
	} else {
		log.WithContext(ctx).Info("no join addresses")
	}
	// join address supplied, but we don't need them
	if !isNew && len(joins) > 0 {
		log.WithContext(ctx).Info("node is already member of cluster")
	}

	// open store
	if err := db.Open(bootstrap); err != nil {
		return errors.Errorf("open store: %w", err)
	}
	s.db = db

	// close the rqlite store
	defer func() {
		if err := db.Close(true); err != nil {
			log.WithContext(ctx).Errorf("close store: %w", err)
		}
	}()

	// execute any requested join operation
	if len(joins) > 0 && isNew {
		log.WithContext(ctx).Infof("join addresses are: %v", joins)
		advAddr := s.config.RaftAddress
		if s.config.RaftAdvertiseAddress != "" {
			advAddr = s.config.RaftAdvertiseAddress
		}

		// try to parse the join duration
		joinDuration, err := time.ParseDuration(s.config.JoinInterval)
		if err != nil {
			return errors.Errorf("parse JoinInterval: %w", err)
		}

		tlsConfig := tls.Config{InsecureSkipVerify: s.config.NoVerify}
		if s.config.X509CACert != "" {
			data, err := ioutil.ReadFile(s.config.X509CACert)
			if err != nil {
				return errors.Errorf("ioutil read: %w", err)
			}
			tlsConfig.RootCAs = x509.NewCertPool()
			if ok := tlsConfig.RootCAs.AppendCertsFromPEM(data); !ok {
				return errors.Errorf("parse root CA certificate in %v", s.config.X509CACert)
			}
		}

		// join rqlite cluster
		joinAddr, err := cluster.Join(
			ctx,
			s.config.JoinSourceIP,
			joins,
			db.ID(),
			advAddr,
			!s.config.RaftNotVoter,
			s.config.JoinAttempts,
			joinDuration,
			&tlsConfig,
		)
		if err != nil {
			return errors.Errorf("join cluster at %v: %w", joins, err)
		}
		log.WithContext(ctx).Infof("successfully joined cluster at %v", joinAddr)
	}

	// wait until the store is in full consensus
	if err := s.waitForConsensus(ctx, db); err != nil {
		return errors.Errorf("wait for consensus: %w", err)
	}
	log.WithContext(ctx).Info("store has reached consensus")

	// start the HTTP API server
	if err := s.startHTTPServer(ctx, db, cs); err != nil {
		return errors.Errorf("start http server: %w", err)
	}
	log.WithContext(ctx).Info("node is ready, block until context is done")

	// mark the rqlite node is ready
	s.ready <- struct{}{}

	// block until context is done
	<-ctx.Done()

	log.WithContext(ctx).Info("rqlite server is stopped")
	return nil
}
