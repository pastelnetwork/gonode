package kademlia

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	b58 "github.com/jbenet/go-base58"
	"github.com/sirupsen/logrus"
)

const (
	defaultNetworkAddr   = "0.0.0.0"
	defaultNetworkPort   = 3000
	defaultExpireTime    = time.Second * 86410
	defaultRefreshTime   = time.Second * 3600
	defaultReplicateTime = time.Second * 3600
	defaultRepublishTime = time.Second * 86400
	defaultPingTime      = time.Second * 1
	defaultReceiveTime   = time.Second * 2
)

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
	ht      *HashTable // the hashtable for routing
	options *Options   // the options of DHT
	network *Network   // the network of DHT
	store   Store      // the storage of DHT
	mutex   sync.Mutex // mutex for update node
}

// Options contains configuration options for the local node
type Options struct {
	ID []byte

	// The local IPv4 or IPv6 address
	IP string

	// The local port to listen for connections
	Port int

	// The nodes being used to bootstrap the network. Without a bootstrap
	// node there is no way to connect to the network. NetworkNodes can be
	// initialized via s.NewNetworkNode()
	BootstrapNodes []*Node

	// The time after which a key/value pair expires;
	// this is a time-to-live (TTL) from the original publication date
	ExpireTime time.Duration

	// Seconds after which an otherwise unaccessed bucket must be refreshed
	RefreshTime time.Duration

	// The interval between Kademlia replication events, when a node is
	// required to publish its entire database
	ReplicateTime time.Duration

	// The time after which the original publisher must
	// republish a key/value pair. Currently not implemented.
	RepublishTime time.Duration

	// The maximum time to wait for a response from a node before discarding
	// it from the bucket
	PingTime time.Duration

	// The maximum time to wait for a response to any message
	ReceiveTime time.Duration
}

// NewDHT returns a new DHT node
func NewDHT(store Store, options *Options) (*DHT, error) {
	if store == nil {
		return nil, errors.New("store is nil")
	}

	// validate the options, if it's invalid, set them to default value
	if options.IP == "" {
		options.IP = defaultNetworkAddr
	}
	if options.Port <= 0 {
		options.Port = defaultNetworkPort
	}
	if options.ExpireTime == 0 {
		options.ExpireTime = defaultExpireTime
	}
	if options.RefreshTime == 0 {
		options.RefreshTime = defaultRefreshTime
	}
	if options.ReplicateTime == 0 {
		options.ReplicateTime = defaultReplicateTime
	}
	if options.RepublishTime == 0 {
		options.RepublishTime = defaultReplicateTime
	}
	if options.PingTime == 0 {
		options.PingTime = defaultPingTime
	}
	if options.ReceiveTime == 0 {
		options.ReceiveTime = defaultReceiveTime
	}

	s := &DHT{
		store:   store,
		options: options,
	}
	// new a hashtable with options
	ht, err := NewHashTable(options)
	if err != nil {
		return nil, fmt.Errorf("new hashtable: %v", err)
	}
	s.ht = ht

	// new network service for dht
	network, err := NewNetwork(s, ht.Self)
	if err != nil {
		return nil, fmt.Errorf("new network: %v", err)
	}
	s.network = network

	return s, nil
}

// Start the distributed hash table
func (s *DHT) Start(ctx context.Context) error {
	// start the network
	if err := s.network.Start(ctx); err != nil {
		return fmt.Errorf("start network: %v", err)
	}

	return nil
}

// Stop the distributed hash table
func (s *DHT) Stop(ctx context.Context) {
	// stop the network
	s.network.Stop(ctx)
}

// keyExpireTime returns the expire time of key
func (s *DHT) keyExpireTime(key []byte) time.Time {
	// calculate the bucket index with key and local node id
	bucket := BucketIndex(key, s.ht.Self.ID)

	var total int
	// total nodes before the bucket which key in
	for i := 0; i < bucket; i++ {
		total += s.ht.GetBucketNodes(i)
	}

	// find the nodes in bucket which is closer to key
	closers := s.ht.CloserNodes(bucket, key)
	score := total + len(closers)
	if score == 0 {
		score = 1
	}
	// if the nodes is more than maximum value, return the expire time
	if score > K {
		return time.Now().Add(s.options.ExpireTime)
	}

	expire := s.options.ExpireTime
	seconds := expire.Nanoseconds() * int64(math.Exp(float64(K/score)))
	return time.Now().Add(time.Second * time.Duration(seconds))
}

// Store the data into the network
func (s *DHT) Store(ctx context.Context, data []byte) ([]byte, error) {
	key := s.store.GetKey(data)

	// expire time for the key
	expiration := s.keyExpireTime(key)
	// replicate time for the key
	replication := time.Now().Add(s.options.ReplicateTime)

	// store the key to local storage
	if err := s.store.Store(ctx, key, data, replication, expiration); err != nil {
		return nil, fmt.Errorf("store data to local storage: %v", err)
	}

	// iterative store the data
	if _, err := s.iterate(ctx, IterateStore, key, data); err != nil {
		return nil, fmt.Errorf("iterative store data: %v", err)
	}

	return []byte(b58.Encode(key)), nil
}

// Get retrieves data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (s *DHT) Get(ctx context.Context, key string) ([]byte, error) {
	decoded := b58.Decode(key)
	if len(decoded) != K {
		return nil, fmt.Errorf("invalid key: %v", key)
	}

	// retrieve the key/value from local storage
	value := s.store.Retrieve(ctx, decoded)
	// if not found locally, iterative find value from kademlia network
	if value == nil {
		data, err := s.iterate(ctx, IterateFindValue, decoded, nil)
		if err != nil {
			return nil, fmt.Errorf("iterate find value: %v", err)
		}
		return data, nil
	}

	return value, nil
}

// Listen begins listening on the socket for incoming messages
// func (s *DHT) Listen() error {
// 	if !s.network.isInitialized() {
// 		return errors.New("socket not created")
// 	}
// 	go s.listen()
// 	go s.timers()
// 	return s.network.Accept()
// }

// Bootstrap attempts to bootstrap the network using the BootstrapNodes provided
// to the Options struct
func (s *DHT) Bootstrap(ctx context.Context) error {
	if len(s.options.BootstrapNodes) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	for _, node := range s.options.BootstrapNodes {
		// sync the node id when it's empty
		if len(node.ID) == 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// new a ping request message
				request := s.newMessage(Ping, node, nil)
				// new a context with timeout
				ctx, cancel := context.WithTimeout(ctx, s.options.PingTime)
				defer cancel()

				// invoke the request and handle the response
				response, err := s.network.Call(ctx, request)
				if err != nil {
					logrus.Errorf("network call: %v", err)
					return
				}
				logrus.Infof("ping response sequence: %v", response.Sender)

				// add the node to the route table
				s.addNode(response.Sender)
			}()
		}
	}

	// wait until all are done
	wg.Wait()

	// if it has nodes in local route tables
	if s.ht.TotalCount() > 0 {
		// iterative find node from the nodes
		if _, err := s.iterate(ctx, IterateFindNode, s.ht.Self.ID, nil); err != nil {
			logrus.Errorf("iterative find node: %v", err)
			return err
		}
	}

	return nil
}

// new a message
func (s *DHT) newMessage(messageType int, receiver *Node, data interface{}) *Message {
	return &Message{
		Sender:      s.ht.Self,
		Receiver:    receiver,
		MessageType: messageType,
		Data:        data,
	}
}

func (s *DHT) doMultiWorkers(ctx context.Context, iterativeType int, target []byte, nl *NodeList, contacted map[string]bool, haveRest bool) []*Message {
	var number int

	// nodes which will be removed
	removedNodes := []*Node{}
	// responses from remote node
	responses := []*Message{}

	var wg sync.WaitGroup
	// send the message to the first (closest) alpha nodes
	for _, node := range nl.Nodes {
		// only contact alpha nodes
		if number >= Alpha && !haveRest {
			break
		}
		// ignore the contacted node
		if contacted[string(node.ID)] {
			continue
		} else {
			contacted[string(node.ID)] = true
		}
		// update the running goroutines
		number++

		wg.Add(1)
		// send and recive message concurrently
		go func(receiver *Node) {
			defer wg.Done()

			var data interface{}
			var messageType int
			switch iterativeType {
			case IterateFindNode, IterateStore:
				messageType = FindNode
				data = &FindNodeRequest{Target: target}
			case IterateFindValue:
				messageType = FindValue
				data = &FindValueRequest{Target: target}
			}

			// new a request message
			request := s.newMessage(messageType, receiver, data)
			// send the request and receive the response
			response, err := s.network.Call(ctx, request)
			if err != nil {
				logrus.Errorf("network call: %v", err)

				// node is unreachable, remove the node
				removedNodes = append(removedNodes, receiver)
			}
			responses = append(responses, response)
		}(node)
	}
	// delete the node which is unreachable
	for _, node := range removedNodes {
		nl.DelNode(node)
	}

	// wait until tasks are done
	wg.Wait()

	return responses
}

// Iterate does an iterative search through the kademlia network
// - IterativeStore - used to store new information in the kademlia network
// - IterativeFindNode - used to bootstrap the network
// - IterativeFindValue - used to find a value among the network given a key
func (s *DHT) iterate(ctx context.Context, iterativeType int, target []byte, data []byte) ([]byte, error) {
	// find the closest contacts for target node from local route tables
	nl := s.ht.ClosestContacts(Alpha, target, []*Node{})
	// no a closer node, stop search
	if nl.Len() == 0 {
		return nil, nil
	}

	// keep the closer node
	closestNode := nl.Nodes[0]
	// if it's find node, reset the refresh timer
	if iterativeType == IterateFindNode {
		bucket := BucketIndex(target, s.ht.Self.ID)
		// reset the refresh time for the bucket
		s.ht.ResetRefreshTime(bucket)
	}

	// According to the Kademlia white paper, after a round of FIND_NODE RPCs
	// fails to provide a node closer than closestNode, we should send a
	// FIND_NODE RPC to all remaining nodes in the node list that have not
	// yet been contacted.
	searchRest := false

	// keep track of nodes contacted
	var contacted = make(map[string]bool)
	for {
		// do the requests concurrently
		responses := s.doMultiWorkers(ctx, iterativeType, target, nl, contacted, searchRest)
		// handle the response one by one
		for _, response := range responses {
			// add the target node to the bucket
			s.addNode(response.Sender)

			switch response.MessageType {
			case FindNode, StoreData:
				v := response.Data.(*FindNodeResponse)
				if len(v.Closest) > 0 {
					nl.AddNodes(v.Closest)
				}
			case FindValue:
				v := response.Data.(*FindValueResponse)
				if v.Value != nil {
					return v.Value, nil
				}
				if len(v.Closest) > 0 {
					nl.AddNodes(v.Closest)
				}
			}
		}

		// stop search
		if !searchRest && len(nl.Nodes) == 0 {
			return nil, nil
		}

		// sort the nodes for node list
		sort.Sort(nl)

		// if closestNode is unchanged
		if bytes.Equal(nl.Nodes[0].ID, closestNode.ID) || searchRest {
			switch iterativeType {
			case FindNode:
				if !searchRest {
					// search all the rest nodes
					searchRest = true
					continue
				}
				return nil, nil
			case FindValue:
				return nil, nil
			case StoreData:
				// store the value to node list
				for i, n := range nl.Nodes {
					// limite the count below K
					if i >= K {
						return nil, nil
					}

					data := &StoreDataRequest{Data: data}
					// new a request message
					request := s.newMessage(StoreData, n, data)
					// send the request and receive the response
					if _, err := s.network.Call(ctx, request); err != nil {
						// <TODO> need to remove the node ?
						logrus.Errorf("network call: %v", err)
					}
				}
				return nil, nil
			}
		} else {
			// update the closest node
			closestNode = nl.Nodes[0]
		}
	}
}

// add a node into the appropriate k bucket
func (s *DHT) addNode(node *Node) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// the bucket index for the node
	index := BucketIndex(s.ht.Self.ID, node.ID)

	// 1. if the node is existed, refresh the node to the end of bucket
	if s.ht.NodeExists(index, node.ID) {
		s.ht.RefreshNode(node.ID)
		return
	}

	// 2. if the bucket is full, ping the first node
	bucket := s.ht.RouteTable[index]
	if len(bucket) == K {
		node := bucket[0]

		// new a ping request message
		request := s.newMessage(Ping, node, nil)
		// new a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), s.options.PingTime)
		defer cancel()

		// invoke the request and handle the response
		response, err := s.network.Call(ctx, request)
		if err != nil {
			// the node is down, remove the node from bucket
			bucket = append(bucket, node)
			bucket = bucket[1:]
		}
		logrus.Infof("ping response sequence: %v", response.Sender)

		// refresh the node to the end of bucket
		bucket = bucket[1:]
		bucket = append(bucket, node)
	} else {
		// 3. append the node to the end of the bucket
		bucket = append(bucket, node)
	}

	// need to reset the route table with the bucket
	s.ht.RouteTable[index] = bucket
}

// func (s *DHT) timers() {
// 	t := time.NewTicker(time.Second)
// 	for {
// 		select {
// 		case <-t.C:
// 			// Refresh
// 			for i := 0; i < B; i++ {
// 				if time.Since(s.ht.getRefreshTimeForBucket(i)) > s.options.RefreshTime {
// 					id := s.ht.getRandomIDFromBucket(K)
// 					s.iterate(IterateFindNode, id, nil)
// 				}
// 			}

// 			// Replication
// 			// keys := s.store.KeysForReplication()
// 			// for _, key := range keys {
// 			// 	value, _ := s.store.Retrieve(key)
// 			// 	s.iterate(IterateStore, key, value)
// 			// }

// 			// // Expiration
// 			// s.store.ExpireKeys()
// 		case <-s.network.getDisconnect():
// 			t.Stop()
// 			s.network.timersFin()
// 			return
// 		}
// 	}
// }
