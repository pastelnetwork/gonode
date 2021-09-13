package kademlia

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/p2p/kademlia/helpers"
	"github.com/pastelnetwork/gonode/pastel"
	"golang.org/x/crypto/sha3"
)

var (
	defaultNetworkAddr   = "0.0.0.0"
	defaultNetworkPort   = 4445
	defaultRefreshTime   = time.Second * 3600
	defaultReplicateTime = time.Second * 3600
	defaultPingTime      = time.Second * 5
	defaultUpdateTime    = time.Minute * 1
)

type KeyType byte

const (
	Sha3Sum256HashLength         = 32
	KeyTypeData          KeyType = 0
	KeyTypeFingerprints  KeyType = 1
	KeyTypeThumbnails    KeyType = 2
)

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
	ht           *HashTable       // the hashtable for routing
	options      *Options         // the options of DHT
	network      *Network         // the network of DHT
	store        Store            // the storage of DHT
	done         chan struct{}    // distributed hash table is done
	cache        storage.KeyValue // store bad bootstrap addresses
	pastelClient pastel.Client
}

// Options contains configuration options for the local node
type Options struct {
	ID []byte

	// The local IPv4 or IPv6 address
	IP string

	// The local port to listen for connections
	Port int

	// The nodes being used to bootstrap the network. Without a bootstrap
	// node there is no way to connect to the network
	BootstrapNodes []*Node

	// PastelClient to retrieve p2p bootstrap addrs
	PastelClient pastel.Client
}

// NewDHT returns a new DHT node
func NewDHT(store Store, pc pastel.Client, options *Options) (*DHT, error) {
	// validate the options, if it's invalid, set them to default value
	if options.IP == "" {
		options.IP = defaultNetworkAddr
	}
	if options.Port <= 0 {
		options.Port = defaultNetworkPort
	}

	s := &DHT{
		store:        store,
		options:      options,
		pastelClient: pc,
		done:         make(chan struct{}),
		cache:        memory.NewKeyValue(),
	}
	// new a hashtable with options
	ht, err := NewHashTable(options)
	if err != nil {
		return nil, fmt.Errorf("new hashtable: %v", err)
	}
	s.ht = ht

	// new network service for dht
	network, err := NewNetwork(s, ht.self)
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

	// start a timer for doing update work
	helpers.StartTimer(ctx, "update work", s.done, defaultUpdateTime, func() error {
		// refresh
		for i := 0; i < B; i++ {
			if time.Since(s.ht.refreshTime(i)) > defaultRefreshTime {
				// refresh the bucket by iterative find node
				id := s.ht.randomIDFromBucket(K)
				if _, err := s.iterate(ctx, IterateFindNode, id, nil); err != nil {
					log.WithContext(ctx).WithError(err).Error("iterate find node failed")
				}
			}
		}

		// replication
		replicationKeys := s.store.KeysForReplication(ctx)
		for _, key := range replicationKeys {
			value, err := s.store.Retrieve(ctx, key)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("store retrieve failed")
				continue
			}
			if value != nil {
				// iteratve store the value
				if _, err := s.iterate(ctx, IterateStore, key, value); err != nil {
					log.WithContext(ctx).WithError(err).Error("iterate store data failed")
				}
			}
		}

		return nil
	})

	return nil
}

// Stop the distributed hash table
func (s *DHT) Stop(ctx context.Context) {
	if s.done != nil {
		close(s.done)
	}

	// stop the network
	s.network.Stop(ctx)
}

// a hash key for the data
func (s *DHT) hashKey(data []byte) []byte {
	sha := sha3.Sum256(data)
	return sha[:]
}

func (s *DHT) storedKey(keyType KeyType, key []byte) []byte {

	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(keyType))
	buffer.Write(key)

	return buffer.Bytes()
}

// Store the data into the network
func (s *DHT) Store(ctx context.Context, keyType KeyType, data []byte) (string, error) {
	key := s.hashKey(data)
	actualKey := s.storedKey(keyType, key)

	// replicate time for the key
	replication := time.Now().Add(defaultReplicateTime)

	// store the key to local storage
	if err := s.store.Store(ctx, actualKey, data, replication); err != nil {
		return "", fmt.Errorf("store data to local storage: %v", err)
	}

	// iterative store the data
	if _, err := s.iterate(ctx, IterateStore, actualKey, data); err != nil {
		return "", fmt.Errorf("iterative store data: %v", err)
	}

	return base58.Encode(key), nil
}

// Retrieve data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (s *DHT) Retrieve(ctx context.Context, keyType KeyType, key string) ([]byte, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return nil, fmt.Errorf("invalid key: %v", key)
	}

	actualKey := s.storedKey(keyType, decoded)

	// retrieve the key/value from local storage
	value, err := s.store.Retrieve(ctx, actualKey)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("store retrive failed")
	}
	// if not found locally, iterative find value from kademlia network
	if value == nil {
		data, err := s.iterate(ctx, IterateFindValue, actualKey, nil)
		if err != nil {
			//return nil, fmt.Errorf("iterate find value: %v", err)
			// be compatible with old key
			return s.oldRetrieve(ctx, decoded)
		}
		return data, nil
	}

	return value, nil
}

// oldRetrieve use previous key name match, be compatible with previous version
func (s *DHT) oldRetrieve(ctx context.Context, key []byte) ([]byte, error) {
	// retrieve the key/value from local storage
	value, err := s.store.Retrieve(ctx, key)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("store retrive failed")
	}
	// if not found locally, iterative find value from kademlia network
	if value == nil {
		data, err := s.iterate(ctx, IterateFindValue, key, nil)
		if err != nil {
			return nil, fmt.Errorf("iterate find value: %v", err)
		}
		return data, nil
	}

	return value, nil
}

// Stats returns stats of DHT
func (s *DHT) Stats(ctx context.Context) (map[string]interface{}, error) {
	fakekey := make([]byte, 20)
	rand.Read(fakekey)
	fakeData := make([]byte, 20)
	rand.Read(fakeData)
	// store the key to local storage
	if err := s.store.Store(ctx, fakekey, fakeData, time.Now()); err != nil {
		log.WithContext(ctx).WithError(err).Error("dhtStoreFailed")
	}

	dbStats, err := s.store.Stats(ctx)
	if err != nil {
		return nil, err
	}

	dhtStats := map[string]interface{}{}
	dhtStats["self"] = s.ht.self
	dhtStats["peers_count"] = len(s.ht.nodes())
	dhtStats["peers"] = s.ht.nodes()
	dhtStats["database"] = dbStats

	// count key status
	detailsStats, err := s.dataDetailsStats(ctx)
	if err != nil {
		return nil, err
	}

	dhtStats["data_details"] = detailsStats

	return dhtStats, nil
}

// Stats returns stats of DHT
func (s *DHT) dataDetailsStats(ctx context.Context) (map[string]interface{}, error) {
	// get allkeys
	cntsMap := map[string]int{}
	increaseCnt := func(keyType string) {
		cnt := cntsMap[keyType]
		cnt = cnt + 1
		cntsMap[keyType] = cnt
	}

	// count number of key for each kind of namespace
	keyFunc := func(key []byte) {
		if len(key) == Sha3Sum256HashLength+1 {
			keyType := KeyType(key[0])
			switch keyType {
			case KeyTypeData:
				increaseCnt("data")
			case KeyTypeFingerprints:
				increaseCnt("fingerprints")
			case KeyTypeThumbnails:
				increaseCnt("thumbnails")
			default:
				increaseCnt("unknown")
			}
		} else {
			// old items
			increaseCnt("unknown")
		}
	}

	err := s.store.ForEachKey(ctx, keyFunc)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{}
	// merge cnt data to output
	for ns, count := range cntsMap {
		stats[ns+"_count"] = count
	}
	return stats, nil
}

// new a message
func (s *DHT) newMessage(messageType int, receiver *Node, data interface{}) *Message {
	return &Message{
		Sender:      s.ht.self,
		Receiver:    receiver,
		MessageType: messageType,
		Data:        data,
	}
}

func (s *DHT) doMultiWorkers(ctx context.Context, iterativeType int, target []byte, nl *NodeList, contacted map[string]bool, haveRest bool) chan *Message {
	// responses from remote node
	responses := make(chan *Message, Alpha)

	go func() {
		// the nodes which are unreachable
		removedNodes := []*Node{}

		var wg sync.WaitGroup

		var number int
		// send the message to the first (closest) alpha nodes
		for _, node := range nl.Nodes {
			// only contact alpha nodes
			if number >= Alpha && !haveRest {
				break
			}
			// ignore the contacted node
			if contacted[string(node.ID)] {
				continue
			}
			contacted[string(node.ID)] = true

			// update the running goroutines
			number++

			log.WithContext(ctx).Debugf("start work %v for node: %s", iterativeType, node.String())

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
					log.WithContext(ctx).WithError(err).Errorf("network call request %s failed", request.String())
					// node is unreachable, remove the node
					removedNodes = append(removedNodes, receiver)
					return
				}

				// send the response to message channel
				responses <- response
			}(node)
		}

		// wait until tasks are done
		wg.Wait()

		// delete the node which is unreachable
		for _, node := range removedNodes {
			nl.DelNode(node)
		}

		// close the message channel
		close(responses)
	}()

	return responses
}

// Iterate does an iterative search through the kademlia network
// - IterativeStore - used to store new information in the kademlia network
// - IterativeFindNode - used to bootstrap the network
// - IterativeFindValue - used to find a value among the network given a key
func (s *DHT) iterate(ctx context.Context, iterativeType int, target []byte, data []byte) ([]byte, error) {
	// find the closest contacts for target node from local route tables
	nl := s.ht.closestContacts(Alpha, target, []*Node{})
	// no a closer node, stop search
	if nl.Len() == 0 {
		return nil, nil
	}
	log.WithContext(ctx).Infof("type: %v, target: %v, nodes: %v", iterativeType, base58.Encode(target), nl.String())

	// keep the closer node
	closestNode := nl.Nodes[0]
	// if it's find node, reset the refresh timer
	if iterativeType == IterateFindNode {
		bucket := s.ht.bucketIndex(target, s.ht.self.ID)
		log.WithContext(ctx).Debugf("bucket for target: %v", base58.Encode(target))

		// reset the refresh time for the bucket
		s.ht.resetRefreshTime(bucket)
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
		for response := range responses {
			log.WithContext(ctx).Debugf("response: %v", response.String())
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

		log.WithContext(ctx).Debugf("id: %v, iterate %d, sorted nodes: %v", base58.Encode(s.ht.self.ID), iterativeType, nl.String())

		// if closestNode is unchanged
		if bytes.Equal(nl.Nodes[0].ID, closestNode.ID) || searchRest {
			switch iterativeType {
			case IterateFindNode:
				if !searchRest {
					// search all the rest nodes
					searchRest = true
					continue
				}
				return nil, nil
			case IterateFindValue:
				return nil, nil
			case IterateStore:
				// store the value to node list
				for i, n := range nl.Nodes {
					// limite the count below K
					if i >= K {
						return nil, nil
					}

					data := &StoreDataRequest{Key: target, Data: data}
					// new a request message
					request := s.newMessage(StoreData, n, data)
					// send the request and receive the response
					if _, err := s.network.Call(ctx, request); err != nil {
						// <TODO> need to remove the node ?
						log.WithContext(ctx).WithError(err).Error("network call")
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

// add a node into the appropriate k bucket, return the removed node if it's full
func (s *DHT) addNode(node *Node) *Node {
	// the bucket index for the node
	index := s.ht.bucketIndex(s.ht.self.ID, node.ID)

	// 1. if the node is existed, refresh the node to the end of bucket
	if s.ht.hasBucketNode(index, node.ID) {
		s.ht.refreshNode(node.ID)
		return nil
	}

	s.ht.mutex.Lock()
	defer s.ht.mutex.Unlock()

	// 2. if the bucket is full, ping the first node
	bucket := s.ht.routeTable[index]
	if len(bucket) == K {
		first := bucket[0]

		// new a ping request message
		request := s.newMessage(Ping, first, nil)
		// new a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), defaultPingTime)
		defer cancel()

		// invoke the request and handle the response
		response, err := s.network.Call(ctx, request)
		if err != nil {
			// the node is down, remove the node from bucket
			bucket = append(bucket, node)
			bucket = bucket[1:]

			// need to reset the route table with the bucket
			s.ht.routeTable[index] = bucket

			log.WithContext(ctx).Debugf("bucket: %d, network call: %v: %v", index, request, err)
			return first
		}
		log.WithContext(ctx).Debugf("ping response: %v", response.String())

		// refresh the node to the end of bucket
		bucket = bucket[1:]
		bucket = append(bucket, node)
	} else {
		// 3. append the node to the end of the bucket
		bucket = append(bucket, node)
	}

	// need to update the route table with the bucket
	s.ht.routeTable[index] = bucket

	return nil
}
