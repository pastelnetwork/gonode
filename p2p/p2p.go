package p2p

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/dao/mem"
)

const (
	logPrefix = "p2p"
)

// P2P represents the p2p service.
type P2P interface {
	Client
	Run(ctx context.Context) error
}

type p2p struct {
	dht    *kademlia.DHT
	config *Config
}

// Run starts the DHT service
func (service *p2p) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := service.configure(ctx); err != nil {
		return err
	}

	if err := service.dht.CreateSocket(); err != nil {
		return err
	}

	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		log.WithContext(ctx).Infof("Kademlia server listening on %q", service.dht.GetNetworkAddr())
		return service.dht.Listen(ctx)
	})
	group.Go(func() error {
		<-ctx.Done()

		log.WithContext(ctx).Infof("Shutting down Kademlia server at %q", service.dht.GetNetworkAddr())
		return service.dht.Disconnect()
	})

	if err := service.dht.Bootstrap(ctx); err != nil {
		return err
	}

	return group.Wait()
}

// Store stores data on the network. This will trigger an iterateStore message.
// The base58 encoded identifier will be returned if the store is successful.
func (service *p2p) Store(ctx context.Context, data []byte) (id string, err error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	return service.dht.Store(ctx, data)
}

// Get retrieves data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (service *p2p) Get(ctx context.Context, key string) (data []byte, found bool, err error) {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	return service.dht.Get(ctx, key)
}

// configure configures service DHT
func (service *p2p) configure(ctx context.Context) error {
	var bootstrapNodes []*kademlia.NetworkNode
	if service.config.BootstrapIP != "" || service.config.BootstrapPort != "" {
		bootstrapNode := kademlia.NewNetworkNode(service.config.BootstrapIP, service.config.BootstrapPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	dht, err := kademlia.NewDHT(ctx, &mem.Key{}, &kademlia.Options{
		BootstrapNodes: bootstrapNodes,
		IP:             service.config.ListenAddresses,
		Port:           service.config.Port,
		UseStun:        service.config.UseStun,
		DataSourceName: service.config.DataDir,
		MemoryDB:       service.config.MemoryDB,
	})
	if err != nil {
		return err
	}
	service.dht = dht

	return nil
}

// New returns a new p2p instance.
func New(config *Config) P2P {
	return &p2p{
		config: config,
	}
}
