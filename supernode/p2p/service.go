package p2p

import (
	"context"
	"log"
)

type Service struct {
	dht DHT
}

type DHT interface {
	Store(data []byte, key []byte) (id string, err error)
	Get(key string) (data []byte, found bool, err error)
	CreateSocket() error
	Listen() error
	GetNetworkAddr() string
	Bootstrap() error
}

func (service *Service) Run(ctx context.Context) error {
	err := service.dht.CreateSocket()
	if err != nil {
		panic(err)
	}

	go func() {
		log.Println("Now listening on " + service.dht.GetNetworkAddr())
		err := service.dht.Listen()
		panic(err)
	}()

	return service.dht.Bootstrap()
}

// NewService returns a new Service instance.
func NewService(dht DHT) *Service {
	return &Service{
		dht: dht,
	}
}
