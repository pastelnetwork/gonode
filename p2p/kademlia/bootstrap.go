package kademlia

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/utils"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

const (
	bootstrapRetryInterval = 10
	badAddrExpiryHours     = 12
)

func (s *DHT) skipBadBootstrapAddrs() {
	skipAddress1 := fmt.Sprintf("%s:%d", "127.0.0.1", s.options.Port)
	skipAddress2 := fmt.Sprintf("%s:%d", "localhost", s.options.Port)
	s.cache.Set(skipAddress1, []byte("true"))
	s.cache.Set(skipAddress2, []byte("true"))
}

func (s *DHT) parseNode(extP2P string, selfAddr string) (*Node, error) {
	if extP2P == "" {
		return nil, errors.New("empty address")
	}

	if strings.Contains(extP2P, "0.0.0.0") {
		return nil, errors.New("invalid address")
	}

	if extP2P == selfAddr {
		return nil, errors.New("self address")
	}

	if _, err := s.cache.Get(extP2P); err == nil {
		return nil, errors.New("configure: skip bad p2p boostrap addr")
	}

	addr := strings.Split(extP2P, ":")
	if len(addr) != 2 {
		return nil, errors.New("wrong number of field")
	}
	ip := addr[0]

	if ip == "" {
		return nil, errors.New("empty ip")
	}

	port, err := strconv.Atoi(addr[1])
	if err != nil {
		return nil, errors.New("invalid port")
	}

	return &Node{
		IP:   ip,
		Port: port,
	}, nil
}

// ConfigureBootstrapNodes connects with pastel client & gets p2p boostrap ip & port
func (s *DHT) ConfigureBootstrapNodes(ctx context.Context) error {
	selfAddress, err := utils.GetExternalIPAddress()
	if err != nil {
		return fmt.Errorf("get external ip addr: %s", err)
	}
	selfAddress = fmt.Sprintf("%s:%d", selfAddress, s.options.Port)

	get := func(ctx context.Context, f func(context.Context) (pastel.MasterNodes, error)) ([]*Node, error) {
		mns, err := f(ctx)
		if err != nil {
			return []*Node{}, err
		}

		mapNodes := map[string]*Node{}
		for _, mn := range mns {
			node, err := s.parseNode(mn.ExtP2P, selfAddress)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("extP2P", mn.ExtP2P).Warn("Skip Bad Boostrap Address")
				continue
			}

			mapNodes[mn.ExtP2P] = node
		}

		nodes := []*Node{}
		for _, node := range mapNodes {
			nodes = append(nodes, node)
		}

		return nodes, nil
	}

	boostrapNodes, err := get(ctx, s.pastelClient.MasterNodesExtra)
	if err != nil {
		return fmt.Errorf("masternodesTop failed: %s", err)
	} else if len(boostrapNodes) == 0 {
		boostrapNodes, err = get(ctx, s.pastelClient.MasterNodesTop)
		if err != nil {
			return fmt.Errorf("masternodesExtra failed: %s", err)
		} else if len(boostrapNodes) == 0 {
			log.WithContext(ctx).Error("unable to fetch bootstrap ip. Missing extP2P")

			return nil
		}
	}

	for _, node := range boostrapNodes {
		log.WithContext(ctx).WithFields(log.Fields{
			"bootstap_ip":    node.IP,
			"bootstrap_port": node.Port,
		}).Info("adding p2p bootstap node")
	}

	s.options.BootstrapNodes = append(s.options.BootstrapNodes, boostrapNodes...)

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
				continue
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
				s.addNode(ctx, response.Sender)
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
