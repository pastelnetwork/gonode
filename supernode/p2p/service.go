package p2p

import (
	"context"
	"os"
	"os/signal"

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
	Store(data []byte, key []byte) (id string, err error)
	Get(key string) (data []byte, found bool, err error)
	CreateSocket() error
	Listen() error
	GetNetworkAddr() string
	Bootstrap() error
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
		err := service.dht.Listen()
		panic(err)
	}()

	if err := service.dht.Bootstrap(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			err := service.dht.Disconnect()
			if err != nil {
				panic(err)
			}
		}
	}()

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
		UseStun:    config.UseStun,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		dht: dht,
	}, nil
}
