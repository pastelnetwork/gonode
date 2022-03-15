package p2p

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/sqlite"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "p2p"
)

var (
	defaultReplicateInterval = time.Second * 3600
	defaultRepublishInterval = time.Second * 3600 * 24
)

// P2P represents the p2p service.
type P2P interface {
	Client

	// Run the node of distributed hash table
	Run(ctx context.Context) error
}

// p2p structure to implements interface
type p2p struct {
	store        kademlia.Store // the store for kademlia network
	dht          *kademlia.DHT  // the kademlia network
	config       *Config        // the service configuration
	running      bool           // if the kademlia network is ready
	pastelClient pastel.Client
	secInfo      *alts.SecInfo
}

// Run the kademlia network
func (s *p2p) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(5 * time.Second):
			if err := s.run(ctx); err != nil {
				if utils.IsContextErr(err) {
					return err
				}

				log.P2P().WithContext(ctx).WithError(err).Error("failed to run kadmelia, retrying.")
			} else {
				log.P2P().WithContext(ctx).Info("kadmelia started successfully")
				return nil
			}
		}
	}
}

// run the kademlia network
func (s *p2p) run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	// configure the kademlia dht for p2p service
	if err := s.configure(ctx); err != nil {
		return errors.Errorf("configure kademlia dht: %w", err)
	}

	// start the node for kademlia network
	if err := s.dht.Start(ctx); err != nil {
		return errors.Errorf("start the kademlia network: %w", err)
	}

	if err := s.dht.ConfigureBootstrapNodes(ctx, s.config.BootstrapIPs); err != nil {
		log.P2P().WithContext(ctx).WithError(err).Error("failed to get bootstap ip")
	}

	// join the kademlia network if bootstrap nodes is set
	if err := s.dht.Bootstrap(ctx, s.config.BootstrapIPs); err != nil {
		// stop the node for kademlia network
		s.dht.Stop(ctx)
		// close the store of kademlia network
		s.store.Close(ctx)
		return errors.Errorf("bootstrap the node: %w", err)
	}
	s.running = true

	//go s.store.InitCleanup(ctx, 5*time.Minute)

	log.P2P().WithContext(ctx).Info("p2p service is started")

	// block until context is done
	<-ctx.Done()

	// stop the node for kademlia network
	s.dht.Stop(ctx)

	// close the store of kademlia network
	s.store.Close(ctx)

	log.P2P().WithContext(ctx).Info("p2p service is stopped")
	return nil
}

// Store store data into the kademlia network
func (s *p2p) Store(ctx context.Context, data []byte) (string, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return "", errors.New("p2p service is not running")
	}

	return s.dht.Store(ctx, data)
}

// Retrieve retrive the data from the kademlia network
func (s *p2p) Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return nil, errors.New("p2p service is not running")
	}

	return s.dht.Retrieve(ctx, key, localOnly...)
}

// Delete delete key in local node
func (s *p2p) Delete(ctx context.Context, key string) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return errors.New("p2p service is not running")
	}

	return s.dht.Delete(ctx, key)
}

// Stats return status of p2p
func (s *p2p) Stats(ctx context.Context) (map[string]interface{}, error) {
	retStats := map[string]interface{}{}
	dhtStats, err := s.dht.Stats(ctx)
	if err != nil {
		return nil, err
	}

	retStats["dht"] = dhtStats
	retStats["config"] = s.config

	// get free space of current kademlia folder
	diskUse, err := utils.DiskUsage(s.config.DataDir)
	if err != nil {
		return nil, errors.Errorf("get disk info failed: %w", err)
	}

	retStats["disk-info"] = &diskUse
	return retStats, nil
}

// NClosestNodes returns a list of n closest masternode to a given string
func (s *p2p) NClosestNodes(ctx context.Context, n int, key string, ignores ...string) []string {
	var ret = make([]string, 0)
	var ignoreNodes = make([]*kademlia.Node, len(ignores))
	for idx, node := range ignores {
		ignoreNodes[idx] = &kademlia.Node{ID: []byte(node)}
	}
	nodes := s.dht.NClosestNodes(ctx, n, key, ignoreNodes...)
	for _, node := range nodes {
		ret = append(ret, string(node.ID))
	}
	return ret
}

// configure the distributed hash table for p2p service
func (s *p2p) configure(ctx context.Context) error {
	// new the local storage
	store, err := sqlite.NewStore(ctx, s.config.DataDir, defaultReplicateInterval, defaultRepublishInterval)
	if err != nil {
		return errors.Errorf("new kademlia store: %w", err)
	}
	s.store = store

	kadOpts := &kademlia.Options{
		BootstrapNodes: []*kademlia.Node{},
		IP:             s.config.ListenAddress,
		Port:           s.config.Port,
		ID:             []byte(s.config.ID),
		PeerAuth:       true, // Enable peer authentication
	}

	if s.config.BootstrapIPs != "" && s.config.ExternalIP != "" {
		kadOpts.ExternalIP = s.config.ExternalIP
	}

	// new a kademlia distributed hash table
	dht, err := kademlia.NewDHT(store, s.pastelClient, s.secInfo, kadOpts)

	if err != nil {
		return errors.Errorf("new kademlia dht: %w", err)
	}
	s.dht = dht

	return nil
}

// New returns a new p2p instance.
func New(config *Config, pastelClient pastel.Client, secInfo *alts.SecInfo) P2P {
	return &p2p{
		config:       config,
		pastelClient: pastelClient,
		secInfo:      secInfo,
	}
}
