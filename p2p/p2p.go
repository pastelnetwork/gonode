package p2p

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/store/db"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	logPrefix = "p2p"
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
}

// Run the kademlia network
func (s *p2p) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	// configure the kademlia dht for p2p service
	if err := s.configure(ctx); err != nil {
		return errors.Errorf("configure kademlia dht: %w", err)
	}

	// start the node for kademlia network
	if err := s.dht.Start(ctx); err != nil {
		return errors.Errorf("start the kademlia network: %w", err)
	}

	if err := s.dht.ConfigureBootstrapNodes(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get bootstap ip")

		return fmt.Errorf("unable to get p2p bootstrap ip: %s", err)
	}

	// join the kademlia network if bootstrap nodes is set
	if err := s.dht.Bootstrap(ctx); err != nil {
		return errors.Errorf("bootstrap the node: %w", err)
	}
	s.running = true

	log.WithContext(ctx).Info("p2p service is started")

	// block until context is done
	<-ctx.Done()

	// stop the node for kademlia network
	s.dht.Stop(ctx)

	// close the store of kademlia network
	s.store.Close(ctx)

	log.WithContext(ctx).Info("p2p service is stopped")
	return nil
}

// Store data into the kademlia network
func (s *p2p) Store(ctx context.Context, data []byte) (string, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return "", errors.New("p2p service is not running")
	}

	return s.dht.Store(ctx, data)
}

// Retrive the data from the kademlia network
func (s *p2p) Retrieve(ctx context.Context, key string) ([]byte, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return nil, errors.New("p2p service is not running")
	}

	return s.dht.Retrieve(ctx, key)
}

// configure the distributed hash table for p2p service
func (s *p2p) configure(ctx context.Context) error {
	// new the local storage
	store, err := db.NewStore(ctx, s.config.DataDir)
	if err != nil {
		return errors.Errorf("new kademlia store: %w", err)
	}
	s.store = store

	// new a kademlia distributed hash table
	dht, err := kademlia.NewDHT(store, s.pastelClient, &kademlia.Options{
		BootstrapNodes: []*kademlia.Node{},
		IP:             s.config.ListenAddress,
		Port:           s.config.Port,
	})
	if err != nil {
		return errors.Errorf("new kademlia dht: %w", err)
	}
	s.dht = dht

	return nil
}

// New returns a new p2p instance.
func New(config *Config, pastelClient pastel.Client) P2P {
	return &p2p{
		config:       config,
		pastelClient: pastelClient,
	}
}
