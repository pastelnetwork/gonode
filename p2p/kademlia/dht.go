package kademlia

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
	"github.com/pastelnetwork/gonode/pastel"
	"golang.org/x/crypto/sha3"
)

var (
	defaultNetworkAddr = "0.0.0.0"
	defaultNetworkPort = 4445
	defaultRefreshTime = time.Second * 3600
	defaultPingTime    = time.Second * 10
)

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
	ht                   *HashTable       // the hashtable for routing
	options              *Options         // the options of DHT
	network              *Network         // the network of DHT
	store                Store            // the storage of DHT
	done                 chan struct{}    // distributed hash table is done
	cache                storage.KeyValue // store bad bootstrap addresses
	pastelClient         pastel.Client
	externalIP           string
	mtx                  sync.Mutex
	authHelper           *AuthHelper
	ignorelist           *BanList
	nodeReplicationTimes map[string]domain.NodeReplicationInfo
	replicationMtx       sync.Mutex
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

	// Authentication is required or not
	PeerAuth bool

	ExternalIP string
}

// NewDHT returns a new DHT node
func NewDHT(ctx context.Context, store Store, pc pastel.Client, secInfo *alts.SecInfo, options *Options) (*DHT, error) {
	// validate the options, if it's invalid, set them to default value
	if options.IP == "" {
		options.IP = defaultNetworkAddr
	}
	if options.Port <= 0 {
		options.Port = defaultNetworkPort
	}
	info, err := store.GetAllReplicationInfo(ctx)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Errorf("get all replicationInfo failed")
	}

	replicationMap := make(map[string]domain.NodeReplicationInfo)
	for _, v := range info {
		replicationMap[string(v.ID)] = v
	}

	s := &DHT{
		store:                store,
		options:              options,
		pastelClient:         pc,
		done:                 make(chan struct{}),
		cache:                memory.NewKeyValue(),
		ignorelist:           NewBanList(ctx),
		nodeReplicationTimes: replicationMap,
		replicationMtx:       sync.Mutex{},
	}

	if options.ExternalIP != "" {
		s.externalIP = options.ExternalIP
	}

	if options.PeerAuth && options.ExternalIP != "" {
		s.authHelper = NewAuthHelper(pc, secInfo)
	}

	// new a hashtable with options
	ht, err := NewHashTable(options)
	if err != nil {
		return nil, fmt.Errorf("new hashtable: %v", err)
	}
	s.ht = ht

	// add bad boostrap addresss
	s.skipBadBootstrapAddrs()

	/*
		// FIXME - use this code to enable secure connection
		secureHelper := credentials.NewClientCreds(s.pastelClient, s.secInfo)
		network, err := NewNetwork(s, ht.self, secureHelper)
	*/

	// new network service for dht
	network, err := NewNetwork(s, ht.self, nil, s.authHelper)
	if err != nil {
		return nil, fmt.Errorf("new network: %v", err)
	}
	s.network = network

	return s, nil
}

func (s *DHT) getExternalIP() (string, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	// if listen IP is localhost - then return itself
	if s.ht.self.IP == "127.0.0.1" || s.ht.self.IP == "localhost" {
		return s.ht.self.IP, nil
	}

	if s.externalIP != "" {
		return s.externalIP, nil
	}

	externalIP, err := utils.GetExternalIPAddress()
	if err != nil {
		return "", fmt.Errorf("get external ip addr: %s", err)
	}

	s.externalIP = externalIP
	return externalIP, nil
}

// Start the distributed hash table
func (s *DHT) Start(ctx context.Context) error {
	// start the network
	if err := s.network.Start(ctx); err != nil {
		return fmt.Errorf("start network: %v", err)
	}

	go s.StartReplicationWorker(ctx)

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

func (s *DHT) retryStore(ctx context.Context, key []byte, data []byte) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 1 * time.Minute
	b.InitialInterval = 200 * time.Millisecond

	return backoff.Retry(backoff.Operation(func() error {
		return s.store.Store(ctx, key, data)
	}), b)
}

// Store the data into the network
func (s *DHT) Store(ctx context.Context, data []byte) (string, error) {
	key := s.hashKey(data)

	retKey := base58.Encode(key)
	// store the key to local storage
	if err := s.retryStore(ctx, key, data); err != nil {
		log.WithContext(ctx).WithError(err).Error("local data store failure after retries")
		return "", fmt.Errorf("retry store data to local storage: %v", err)
	}

	return retKey, nil
}

// StoreBatch will store a batch of values with their SHA256 hash as the key
func (s *DHT) StoreBatch(ctx context.Context, values [][]byte) error {
	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}

	log.WithContext(ctx).WithField("taskID", taskID).Info("store db batch begin")
	if err := s.store.StoreBatch(ctx, values); err != nil {
		return fmt.Errorf("store batch: %v", err)
	}

	log.WithContext(ctx).WithField("taskID", taskID).Info("store db batch done")

	// Define a maximum number of goroutines
	/*const maxGoroutines = 1000
	const minSuccessRate = 0.6

	// Create a worker pool and results slice
	workers := make(chan struct{}, maxGoroutines)
	results := make([]error, len(values))

	// Create wait group
	var wg sync.WaitGroup
	wg.Add(len(values))

	// Launch the workers
	for i := 0; i < len(values); i++ {
		go func(i int, value []byte) {
			// Acquire a slot from the workers channel
			workers <- struct{}{}

			// Perform the operation and save the result
			key := s.hashKey(value)
			_, err := s.iterate(ctx, IterateStore, key, value)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("iterate data batch store failure")
			}
			results[i] = err

			// Release the slot when the goroutine finishes
			<-workers

			wg.Done()
		}(i, values[i]) // pass a copy of the index and value to the goroutine
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Count successes and failures
	var successCount int
	for i := 0; i < len(results); i++ {
		if results[i] == nil {
			successCount++
		}
	}

	log.WithContext(ctx).WithField("success", successCount).WithField("total", len(values)).Info("store batch network results")

	// Check if successful requests are less than the minimum success rate
	successRate := float64(successCount) / float64(len(values))
	if successRate < minSuccessRate {
		return fmt.Errorf("less than %.2f%% of requests were successful", minSuccessRate*100)
	}
	*/

	log.WithContext(ctx).WithField("taskID", taskID).Info("store network batch begin")
	// Launch the workers
	for i := 0; i < len(values); i++ {
		key := s.hashKey(values[i])
		_, err := s.iterate(ctx, IterateStore, key, values[i])
		if err != nil {
			log.WithContext(ctx).WithError(err).Error("iterate data batch store failure")
		}
	}
	log.WithContext(ctx).WithField("taskID", taskID).Info("store network batch done")

	return nil
}

// Retrieve data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (s *DHT) Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return nil, fmt.Errorf("invalid key: %v", key)
	}

	log.WithContext(ctx).WithField("key", key).Info("begin local data retrieval")
	// retrieve the key/value from local storage
	value, err := s.store.Retrieve(ctx, decoded)
	if err == nil && len(value) > 0 {
		return value, nil
	}
	log.WithContext(ctx).WithField("key", key).WithField("len", len(value)).Info("finish local data retrieval")

	// if local only option is set, do not search just return error
	if len(localOnly) > 0 && localOnly[0] {
		return nil, fmt.Errorf("local-only failed to get properly: " + err.Error())
	}
	log.WithContext(ctx).WithField("key", key).Info("Not found locally, searching in other nodes")

	// if not found locally, iterative find value from kademlia network
	peerValue, err := s.iterate(ctx, IterateFindValue, decoded, nil)
	if err != nil {
		return nil, errors.Errorf("retrieve from peer: %w", err)
	}
	if len(peerValue) > 0 {
		log.WithContext(ctx).WithField("key", key).WithField("data len", len(peerValue)).Info("Not found locally, retrieved from other nodes")
	} else {
		log.WithContext(ctx).WithField("key", key).Info("Not found locally, not found in other nodes")
	}

	return peerValue, nil
}

// Delete delete key in local node
func (s *DHT) Delete(ctx context.Context, key string) error {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return fmt.Errorf("invalid key: %v", key)
	}

	s.store.Delete(ctx, decoded)

	return nil
}

// Stats returns stats of DHT
func (s *DHT) Stats(ctx context.Context) (map[string]interface{}, error) {
	dbStats, err := s.store.Stats(ctx)
	if err != nil {
		return nil, err
	}

	dhtStats := map[string]interface{}{}
	dhtStats["self"] = s.ht.self
	dhtStats["peers_count"] = len(s.ht.nodes())
	dhtStats["peers"] = s.ht.nodes()
	dhtStats["database"] = dbStats

	return dhtStats, nil
}

// new a message
func (s *DHT) newMessage(messageType int, receiver *Node, data interface{}) *Message {
	externalIP, _ := s.getExternalIP()
	sender := &Node{
		IP:   externalIP,
		ID:   s.ht.self.ID,
		Port: s.ht.self.Port,
	}
	return &Message{
		Sender:      sender,
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
		//var removedNodes []*Node

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

			log.P2P().WithContext(ctx).Debugf("start work %v for node: %s", iterativeType, node.String())

			wg.Add(1)
			// send and receive message concurrently
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
					log.P2P().WithContext(ctx).WithError(err).Errorf("network call request %s failed", request.String())
					// node is unreachable, remove the node
					//removedNodes = append(removedNodes, receiver)
					return
				}

				// send the response to message channel
				responses <- response
			}(node)
		}

		// wait until tasks are done
		wg.Wait()

		// delete the node which is unreachable
		/*go func() {
			for _, node := range removedNodes {
				nl.DelNode(node)
			}
		}()*/

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
	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}

	sKey := base58.Encode(target)

	igList := s.ignorelist.ToNodeList()
	// find the closest contacts for the target node from local route tables
	nl := s.ht.closestContacts(Alpha, target, igList)
	if len(igList) > 0 {
		log.P2P().WithContext(ctx).WithField("nodes", nl.String()).WithField("ignored", s.ignorelist.String()).Info("closest contacts")
	}
	// if no closer node, stop search
	if nl.Len() == 0 {
		return nil, nil
	}
	log.P2P().WithContext(ctx).WithField("task_id", taskID).Infof("type: %v, target: %v, nodes: %v", iterativeType, sKey, nl.String())

	// keep the closer node
	closestNode := nl.Nodes[0]
	// if it's a find node, reset the refresh timer
	if iterativeType == IterateFindNode {
		bucket := s.ht.bucketIndex(target, s.ht.self.ID)
		log.P2P().WithContext(ctx).Debugf("bucket for target: %v", sKey)

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

	// Set a timeout for the iteration process
	timeout := time.After(10 * time.Second) // Adjust the timeout duration as needed

	// Set a maximum number of iterations to prevent indefinite looping
	maxIterations := 100 // Adjust the maximum iterations as needed

	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("begin iteration")

	for i := 0; i < maxIterations; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("iterate cancelled: %w", ctx.Err())
		case <-timeout:
			log.P2P().WithContext(ctx).Debug("iteration timed out")
			return nil, nil
		default:
			// Do the requests concurrently
			responses := s.doMultiWorkers(ctx, iterativeType, target, nl, contacted, searchRest)

			// Handle the responses one by one
			for response := range responses {
				// Add the target node to the bucket
				s.addNode(ctx, response.Sender)

				switch response.MessageType {
				case FindNode, StoreData:
					v, ok := response.Data.(*FindNodeResponse)
					if ok && v.Status.Result == ResultOk {
						if len(v.Closest) > 0 {
							nl.AddNodes(v.Closest)
						}
					}
				case FindValue:
					v, ok := response.Data.(*FindValueResponse)
					if !ok {
						log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Error("invalid find value response")
					}

					if v.Status.Result == ResultOk {
						if len(v.Value) > 0 {
							log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).WithField("resp", len(v.Value)).Info("iterate found value from network")
							return v.Value, nil
						} else if len(v.Closest) > 0 {
							nl.AddNodes(v.Closest)
						}
					} else {
						log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).WithError(errors.New(v.Status.ErrMsg)).Info("iterate could not found value from node")
					}
				}
			}

			// Stop search if no more nodes to contact
			if !searchRest && len(nl.Nodes) == 0 {
				log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("search stopped")
				return nil, nil
			}

			// Sort the nodes in the node list
			nl.Sort()

			log.P2P().WithContext(ctx).Debugf("id: %v, iterate %d, sorted nodes: %v", base58.Encode(s.ht.self.ID), iterativeType, nl.String())

			// If closestNode is unchanged
			if bytes.Equal(nl.Nodes[0].ID, closestNode.ID) || searchRest {
				switch iterativeType {
				case IterateFindNode:
					if !searchRest {
						// Search all the remaining nodes
						searchRest = true
						continue
					}
					return nil, nil
				case IterateFindValue:
					log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Error("attempt to find value from remaining network")
					// Attempt to find the value from the remaining nodes
					for _, n := range nl.Nodes {
						request := &FindValueRequest{Target: target}
						reqMsg := s.newMessage(FindValue, n, request)
						// Send the request and receive the response
						rspMsg, err := s.network.Call(ctx, reqMsg)
						if err != nil {
							return nil, errors.Errorf("network call: %w", err)
						}

						v, ok := rspMsg.Data.(*FindValueResponse)
						if ok && v.Status.Result == ResultOk && len(v.Value) > 0 {
							return v.Value, nil
						}
					}

					log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Error("could not find value from remaining network")
					return nil, nil

				case IterateStore:
					// Store the value to the node list
					storeCount := 0
					for i := 0; i < len(nl.Nodes); i++ {
						logEntry := log.P2P().WithContext(ctx).WithField("node", nl.Nodes[i]).WithField("task_id", taskID)
						// Limit the count below K
						if i >= K {
							return nil, nil
						}
						n := nl.Nodes[i]

						if s.ignorelist.Banned(n) {
							logEntry.Info("ignore node as its continuous failed count is above threshold")
							continue
						}

						request := &StoreDataRequest{Data: data}
						response, err := s.sendStoreData(ctx, n, request)
						if err != nil {
							logEntry.WithError(err).Error("send store data failed")
						} else if response.Status.Result != ResultOk {
							logEntry.WithError(errors.New(response.Status.ErrMsg)).Error("reply store data failed")
						} else {
							storeCount++
						}
					}

					log.P2P().WithContext(ctx).WithField("store count", storeCount).WithField("task_id", taskID).Info("task data stored")

					return nil, nil
				}
			} else {
				// Update the closest node
				closestNode = nl.Nodes[0]
			}
		}
	}
	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("finish iteration without results")
	return nil, nil
}

func (s *DHT) sendStoreData(ctx context.Context, n *Node, request *StoreDataRequest) (*StoreDataResponse, error) {
	// new a request message
	reqMsg := s.newMessage(StoreData, n, request)

	/*b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Second
	b.InitialInterval = 200 * time.Millisecond

	var rspMsg *Message
	if err := backoff.Retry(backoff.Operation(func() error {
		var err error
		// send the request and receive the response
		rspMsg, err = s.network.Call(ctx, reqMsg)
		if err != nil {
			return errors.Errorf("network call: %w", err)
		}

		return nil
	}), b); err != nil {
		return nil, fmt.Errorf("retry send store data failed: %w", err)
	}*/

	rspMsg, err := s.network.Call(ctx, reqMsg)
	if err != nil {
		return nil, errors.Errorf("network call: %w", err)
	}

	response, ok := rspMsg.Data.(*StoreDataResponse)
	if !ok {
		return nil, errors.New("invalid StoreDataResponse")
	}

	return response, nil
}

// add a node into the appropriate k bucket, return the removed node if it's full
func (s *DHT) addNode(ctx context.Context, node *Node) *Node {
	// ensure this is not itself address
	if bytes.Equal(node.ID, s.ht.self.ID) {
		log.P2P().WithContext(ctx).Debug("trying to add itself")
		return nil
	}

	// the bucket index for the node
	index := s.ht.bucketIndex(s.ht.self.ID, node.ID)

	// 1. if the node is existed, refresh the node to the end of bucket
	if s.ht.hasBucketNode(index, node.ID) {
		s.ht.refreshNode(node.ID)
		return nil
	}

	if err := s.updateReplicationNode(ctx, node.ID, node.IP, node.Port, true); err != nil {
		log.P2P().WithContext(ctx).WithField("node-id", string(node.ID)).WithField("node-ip", node.IP).WithError(err).Error("update replication node failed")
	} else {
		log.P2P().WithContext(ctx).WithField("node-id", string(node.ID)).WithField("node-ip", node.IP).Info("adding new node")
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
		ctx, cancel := context.WithTimeout(ctx, defaultPingTime)
		defer cancel()

		// invoke the request and handle the response
		response, err := s.network.Call(ctx, request)
		if err != nil {
			// the node is down, remove the node from bucket
			bucket = append(bucket, node)
			bucket = bucket[1:]

			// need to reset the route table with the bucket
			s.ht.routeTable[index] = bucket

			log.P2P().WithContext(ctx).Debugf("bucket: %d, network call: %v: %v", index, request, err)
			return first
		}
		log.P2P().WithContext(ctx).Debugf("ping response: %v", response.String())

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

// NClosestNodes get n closest nodes to a key string
func (s *DHT) NClosestNodes(_ context.Context, n int, key string, ignores ...*Node) []*Node {
	list := s.ignorelist.ToNodeList()
	ignores = append(ignores, list...)
	nodeList := s.ht.closestContacts(n, base58.Decode(key), ignores)

	return nodeList.Nodes
}

// LocalStore the data into the network
func (s *DHT) LocalStore(ctx context.Context, key string, data []byte) (string, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return "", fmt.Errorf("invalid key: %v", key)
	}

	// store the key to local storage
	if err := s.retryStore(ctx, decoded, data); err != nil {
		log.WithContext(ctx).WithError(err).Error("local data store failure after retries")
		return "", fmt.Errorf("retry store data to local storage: %v", err)
	}

	return key, nil
}
