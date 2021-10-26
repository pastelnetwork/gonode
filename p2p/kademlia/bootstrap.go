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

// ConfigureBootstrapNodes connects with pastel client & gets p2p boostrap ip & port
func (s *DHT) ConfigureBootstrapNodes(ctx context.Context) error {
	selfAddress, err := utils.GetExternalIPAddress()
	if err != nil {
		return fmt.Errorf("get external ip addr: %s", err)
	}
	selfAddress = fmt.Sprintf("%s:%d", selfAddress, s.options.Port)

	get := func(ctx context.Context, f func(context.Context) (pastel.MasterNodes, error)) ([]string, error) {
		extP2PList := []string{}
		mns, err := f(ctx)
		if err != nil {
			return extP2PList, err
		}

		for _, mn := range mns {
			if mn.ExtP2P == "" {
				continue
			}

			if mn.ExtP2P == selfAddress {
				continue
			}

			if _, err := s.cache.Get(mn.ExtP2P); err == nil {
				log.WithContext(ctx).WithField("addr", mn.ExtP2P).Info("configure: skip bad p2p boostrap addr")
				continue
			}

			extP2PList = append(extP2PList, mn.ExtP2P)
		}

		return extP2PList, nil
	}

	extP2PList, err := get(ctx, s.pastelClient.MasterNodesTop)
	if err != nil {
		return fmt.Errorf("masternodesTop failed: %s", err)
	} else if len(extP2PList) == 0 {
		extP2PList, err = get(ctx, s.pastelClient.MasterNodesExtra)
		if err != nil {
			return fmt.Errorf("masternodesExtra failed: %s", err)
		} else if len(extP2PList) == 0 {
			log.WithContext(ctx).Error("unable to fetch bootstrap ip. Missing extP2P")

			return nil
		}
	}

	available := false
	for _, extP2P := range extP2PList {
		addr := strings.Split(extP2P, ":")
		if len(addr) != 2 {
			log.WithContext(ctx).WithField("extP2P", extP2P).Warn("invalid extP2P")
			continue
		}

		ip := addr[0]
		port, err := strconv.Atoi(addr[1])
		if err != nil {
			log.WithContext(ctx).WithField("extP2P", extP2P).Warn("invalid extP2P's port")
			continue
		}

		log.WithContext(ctx).WithField("bootstap_ip", ip).
			WithField("bootstrap_port", port).Info("adding p2p bootstap node")

		bootstrapNode := &Node{
			IP:   ip,
			Port: port,
		}

		s.options.BootstrapNodes = append(s.options.BootstrapNodes, bootstrapNode)
		available = true
	}

	if !available {
		return errors.New("not found any boostrap node")
	}

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
