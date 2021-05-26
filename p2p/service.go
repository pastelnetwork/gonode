package p2p

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
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
		return err
	}

	go func() {
		log.WithContext(ctx).Infof("Server listening on %q", service.dht.GetNetworkAddr())
		if err := service.dht.Listen(ctx); err != nil {
			log.WithContext(ctx).Fatal(err)
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		defer errors.Recover(log.Fatal)

		<-ctx.Done()
		log.WithContext(ctx).Infof("Server is shutting down...")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := service.dht.Disconnect(ctx); err != nil {
			errCh <- errors.Errorf("gracefully shutdown the server: %w", err)
		}
		close(errCh)
	}()

	return service.dht.Bootstrap(ctx)
}

// Setup configures service DHT
func (service *Service) Setup(ctx context.Context, config Config) error {
	var bootstrapNodes []*kademlia.NetworkNode
	if config.BootstrapIP != "" || config.BootstrapPort != "" {
		bootstrapNode := kademlia.NewNetworkNode(config.BootstrapIP, config.BootstrapPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	dht, err := kademlia.NewDHT(ctx, &mem.Key{}, &kademlia.Options{
		BootstrapNodes: bootstrapNodes,
		IP:             config.ListenAddresses,
		Port:           config.Port,
		UseStun:        config.UseStun,
	})
	if err != nil {
		return err
	}

	service.dht = dht
	return nil
}

// New returns a new Service instance.
func New() *Service {
	return &Service{}
}
