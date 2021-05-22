package p2p

import (
	"context"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia"
)

const (
	logPrefix = "p2p"
)

// Service represents the p2p service.
type Service struct {
	dht DHT
}

// DHT represents the methods by which a library consumer interacts with a DHT.
type DHT interface {
	Store(ctx context.Context, data []byte, key []byte) (id string, err error)
	Get(ctx context.Context, key string) (data []byte, found bool, err error)
	CreateSocket() error
	Listen(ctx context.Context) error
	GetNetworkAddr() string
	Bootstrap(ctx context.Context) error
	UseStun() bool
	Disconnect() error
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
		err := service.dht.Listen(ctx)
		panic(err)
	}()

	if err := service.dht.Bootstrap(ctx); err != nil {
		return err
	}

	return nil
}

// NewService returns a new Service instance.
func NewService(config *Config) (*Service, error) {
	var bootstrapNodes []*kademlia.NetworkNode
	if config.BootstrapIP != "" || config.BootstrapPort != "" {
		bootstrapNode := kademlia.NewNetworkNode(config.BootstrapIP, config.BootstrapPort)
		bootstrapNodes = append(bootstrapNodes, bootstrapNode)
	}

	dht, err := kademlia.NewDHT(&kademlia.MemoryStore{}, &kademlia.Options{
		BootstrapNodes: bootstrapNodes,
		IP:             config.IP,
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
