package kademlia

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	bootstrapRetryInterval = 10
	badAddrExpiryHours     = 12
)

// ConfigureBootstrapNodes connects with pastel client & gets p2p boostrap ip & port
func (s *DHT) ConfigureBootstrapNodes(ctx context.Context) error {
	get := func(ctx context.Context, f func(context.Context) (pastel.MasterNodes, error)) (string, error) {
		mns, err := f(ctx)
		if err != nil {
			return "", err
		}

		for _, mn := range mns {
			if mn.ExtP2P != "" {
				if _, err := s.cache.Get(mn.ExtP2P); err == nil {
					log.WithContext(ctx).WithField("addr", mn.ExtP2P).Info("configure: skip bad p2p boostrap addr")

					continue
				}

				return mn.ExtP2P, nil
			}
		}

		return "", nil
	}

	extP2p, err := get(ctx, s.pastelClient.MasterNodesTop)
	if err != nil {
		return fmt.Errorf("masternodesTop failed: %s", err)
	} else if extP2p == "" {
		extP2p, err = get(ctx, s.pastelClient.MasterNodesExtra)
		if err != nil {
			return fmt.Errorf("masternodesExtra failed: %s", err)
		} else if extP2p == "" {
			log.WithContext(ctx).Info("unable to fetch bootstrap ip. Missing extP2P")

			return nil
		}
	}

	addr := strings.Split(extP2p, ":")
	if len(addr) != 2 {
		return fmt.Errorf("invalid extP2P format: %s", extP2p)
	}

	ip := addr[0]
	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return fmt.Errorf("invalid extP2P port: %s", extP2p)
	}

	log.WithContext(ctx).WithField("bootstap_ip", ip).
		WithField("bootstrap_port", port).Info("adding p2p bootstap node")

	bootstrapNode := &Node{
		IP:   ip,
		Port: port,
	}
	s.options.BootstrapNodes = append(s.options.BootstrapNodes, bootstrapNode)

	return nil
}

// Bootstrap attempts to bootstrap the network using the BootstrapNodes provided
// to the Options struct
func (s *DHT) Bootstrap(ctx context.Context) error {
	if len(s.options.BootstrapNodes) == 0 {
		time.AfterFunc(bootstrapRetryInterval*time.Minute, func() {
			s.retryBootstrap(ctx)
		})

		return nil
	}

	var wg sync.WaitGroup
	for _, node := range s.options.BootstrapNodes {
		// sync the node id when it's empty
		if len(node.ID) == 0 {
			addr := fmt.Sprintf("%s:%v", node.IP, node.Port)
			if _, err := s.cache.Get(addr); err == nil {
				log.WithContext(ctx).WithField("addr", addr).Info("skip bad p2p boostrap addr")
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				// new a ping request message
				request := s.newMessage(Ping, node, nil)
				// new a context with timeout
				ctx, cancel := context.WithTimeout(ctx, defaultPingTime)
				defer cancel()

				// invoke the request and handle the response
				response, err := s.network.Call(ctx, request)
				if err != nil {
					s.cache.SetWithExpiry(addr, []byte("true"), badAddrExpiryHours*time.Hour)

					log.WithContext(ctx).WithError(err).Error("network call failed")
					return
				}
				log.WithContext(ctx).Debugf("ping response: %v", response.String())

				// add the node to the route table
				s.addNode(response.Sender)
			}()
		}
	}

	// wait until all are done
	wg.Wait()

	// if it has nodes in local route tables
	if s.ht.totalCount() > 0 {
		// iterative find node from the nodes
		if _, err := s.iterate(ctx, IterateFindNode, s.ht.self.ID, nil); err != nil {
			log.WithContext(ctx).WithError(err).Error("iterative find node failed")
			return err
		}
	} else {
		time.AfterFunc(bootstrapRetryInterval*time.Minute, func() {
			s.retryBootstrap(ctx)
		})
	}

	return nil
}

func (s *DHT) retryBootstrap(ctx context.Context) {
	if err := s.ConfigureBootstrapNodes(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("retry failed to get bootstap ip")
	}

	// join the kademlia network if bootstrap nodes is set
	if err := s.Bootstrap(ctx); err != nil {
		log.WithContext(ctx).WithError(err).Error("retry failed - bootstrap the node.")
	}
}
