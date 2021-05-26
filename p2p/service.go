package p2p

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
	"github.com/pastelnetwork/gonode/p2p/kademlia/dao/mem"
)

const (
	logPrefix = "p2p"
)

// Service represents the p2p service.
type Service struct {
	dht *kademlia.DHT
}

// Run starts the DHT service
func (service *Service) Run(ctx context.Context) error {
	ctx = log.ContextWithPrefix(ctx, logPrefix)

	err := service.dht.CreateSocket()
	if err != nil {
		panic(err)
	}

	go func() {
		log.WithContext(ctx).Infof("Server listening on %q", service.dht.GetNetworkAddr())
		if err := service.dht.Listen(ctx); err != nil {
			log.WithContext(ctx).Error(err)
			<-ctx.Done()
		}
	}()

	return service.dht.Bootstrap(ctx)
}

// New returns a new Service instance.
func New(config Config) (*Service, error) {
	var bootstrapNodes []*kademlia.NetworkNode
	if config.BootstrapIP != "" || config.BootstrapPort != "" {
		bootstrapNode := kademlia.NewNetworkNode(config.BootstrapIP, config.BootstrapPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	dht, err := kademlia.NewDHT(&mem.Key{}, &kademlia.Options{
		BootstrapNodes: bootstrapNodes,
		IP:             config.ListenAddress,
		Port:           config.Port,
		UseStun:        config.UseStun,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		dht: dht,
	}, nil
}
