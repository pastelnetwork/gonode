package p2p

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/p2p/kademlia/store/meta"

	"github.com/btcsuite/btcutil/base58"
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
	// B is the number of bits in a SHA256 hash
	B = 256
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
	metaStore    kademlia.MetaStore
	dht          *kademlia.DHT // the kademlia network
	config       *Config       // the service configuration
	running      bool          // if the kademlia network is ready
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

	// close the store of kademlia network
	defer s.store.Close(ctx)

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
		return errors.Errorf("bootstrap the node: %w", err)
	}
	s.running = true

	log.P2P().WithContext(ctx).Info("p2p service is started")

	// block until context is done
	<-ctx.Done()

	// stop the node for kademlia network
	s.dht.Stop(ctx)

	log.P2P().WithContext(ctx).Info("p2p service is stopped")
	return nil
}

// Store store data into the kademlia network
func (s *p2p) Store(ctx context.Context, data []byte, typ int) (string, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return "", errors.New("p2p service is not running")
	}

	return s.dht.Store(ctx, data, typ)
}

// StoreBatch will store a batch of values with their SHA256 hash as the key
func (s *p2p) StoreBatch(ctx context.Context, data [][]byte, typ int) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return errors.New("p2p service is not running")
	}

	return s.dht.StoreBatch(ctx, data, typ)
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

	log.WithContext(ctx).WithField("no_of_closest_nodes", n).WithField("file_hash", key).
		WithField("closest_nodes", ret).
		Info("closest nodes retrieved")
	return ret
}

// configure the distributed hash table for p2p service
func (s *p2p) configure(ctx context.Context) error {
	// new the local storage
	kadOpts := &kademlia.Options{
		BootstrapNodes: []*kademlia.Node{},
		IP:             s.config.ListenAddress,
		Port:           s.config.Port,
		ID:             []byte(s.config.ID),
		PeerAuth:       true, // Enable peer authentication
	}

	if len(kadOpts.ID) == 0 {
		errors.Errorf("node id is empty")
	}

	// We Set ExternalIP only for integration tests
	if s.config.BootstrapIPs != "" && s.config.ExternalIP != "" {
		kadOpts.ExternalIP = s.config.ExternalIP
		kadOpts.PeerAuth = false
	}

	// new a kademlia distributed hash table
	dht, err := kademlia.NewDHT(ctx, s.store, s.metaStore, s.pastelClient, s.secInfo, kadOpts)

	if err != nil {
		return errors.Errorf("new kademlia dht: %w", err)
	}
	s.dht = dht

	return nil
}

// New returns a new p2p instance.
func New(ctx context.Context, config *Config, pastelClient pastel.Client, secInfo *alts.SecInfo) (P2P, error) {
	store, err := sqlite.NewStore(ctx, config.DataDir, defaultReplicateInterval, defaultRepublishInterval)
	if err != nil {
		return nil, errors.Errorf("new kademlia store: %w", err)
	}

	meta, err := meta.NewStore(ctx, config.DataDir)
	if err != nil {
		return nil, errors.Errorf("new kademlia meta store: %w", err)
	}

	return &p2p{
		store:        store,
		metaStore:    meta,
		config:       config,
		pastelClient: pastelClient,
		secInfo:      secInfo,
	}, nil
}

// LocalStore store data into the kademlia network
func (s *p2p) LocalStore(ctx context.Context, key string, data []byte) (string, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return "", errors.New("p2p service is not running")
	}

	return s.dht.LocalStore(ctx, key, data)
}

// DisableKey adds key to disabled keys list - It takes in a B58 encoded SHA-256 hash
func (s *p2p) DisableKey(ctx context.Context, b58EncodedHash string) error {
	decoded := base58.Decode(b58EncodedHash)
	if len(decoded) != B/8 {
		return fmt.Errorf("invalid key: %v", b58EncodedHash)
	}

	return s.metaStore.Store(ctx, decoded)
}

// EnableKey removes key from disabled list - It takes in a B58 encoded SHA-256 hash
func (s *p2p) EnableKey(ctx context.Context, b58EncodedHash string) error {
	decoded := base58.Decode(b58EncodedHash)
	if len(decoded) != B/8 {
		return fmt.Errorf("invalid key: %v", b58EncodedHash)
	}

	s.metaStore.Delete(ctx, decoded)

	return nil
}

// GetLocalKeys returns a list of all keys stored locally
func (s *p2p) GetLocalKeys(ctx context.Context, from *time.Time, to time.Time) ([]string, error) {
	if from == nil {
		fromTime, err := s.store.GetOwnCreatedAt(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get own created at: %w", err)
		}

		from = &fromTime
	}

	keys, err := s.store.GetLocalKeys(*from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get local keys: %w", err)
	}

	retkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		str, err := hex.DecodeString(keys[i])
		if err != nil {
			log.WithContext(ctx).WithField("key", keys[i]).Error("replicate failed to hex decode key")
			continue
		}

		retkeys[i] = base58.Encode(str)
	}

	return retkeys, nil
}
