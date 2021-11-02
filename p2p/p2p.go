package p2p

import (
	"context"
	"fmt"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/utils"
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
	secInfo      *alts.SecInfo
}

// Run the kademlia network
func (s *p2p) Run(ctx context.Context) error {
	for {
		if err := s.run(ctx); err != nil {
			if utils.IsContextErr(err) {
				return err
			}
			log.WithContext(ctx).WithError(err).Error("failed to run kadmelia, retrying.")
		} else {
			return nil
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

	if err := s.dht.ConfigureBootstrapNodes(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("failed to get bootstap ip")

		return fmt.Errorf("unable to get p2p bootstrap ip: %s", err)
	}

	// join the kademlia network if bootstrap nodes is set
	if err := s.dht.Bootstrap(ctx); err != nil {
		// stop the node for kademlia network
		s.dht.Stop(ctx)
		// close the store of kademlia network
		s.store.Close(ctx)
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

// Store store data into the kademlia network
func (s *p2p) Store(ctx context.Context, data []byte) (string, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return "", errors.New("p2p service is not running")
	}

	return s.dht.Store(ctx, data)
}

// Retrieve retrive the data from the kademlia network
func (s *p2p) Retrieve(ctx context.Context, key string) ([]byte, error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	if !s.running {
		return nil, errors.New("p2p service is not running")
	}

	return s.dht.Retrieve(ctx, key)
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

// configure the distributed hash table for p2p service
func (s *p2p) configure(ctx context.Context) error {
	// new the local storage
	store, err := db.NewStore(ctx, s.config.DataDir)
	if err != nil {
		return errors.Errorf("new kademlia store: %w", err)
	}
	s.store = store

	/*
		// FIXME - use this code to enable secure connection
		transportCredentials := credentials.NewClientCreds(s.pastelClient, s.secInfo)

		// new a kademlia distributed hash table
		dht, err := kademlia.NewDHT(store, s.pastelClient, transportCredentials, &kademlia.Options{
			BootstrapNodes: []*kademlia.Node{},
			IP:             s.config.ListenAddress,
			Port:           s.config.Port,
		})
	*/

	// new a kademlia distributed hash table
	dht, err := kademlia.NewDHT(store, s.pastelClient, nil, &kademlia.Options{
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
func New(config *Config, pastelClient pastel.Client, secInfo *alts.SecInfo) P2P {
	return &p2p{
		config:       config,
		pastelClient: pastelClient,
		secInfo:      secInfo,
	}
}
