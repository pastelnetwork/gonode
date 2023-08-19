package kademlia

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"sync/atomic"

	"github.com/btcsuite/btcutil/base58"
	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/net/credentials/alts"
	"github.com/pastelnetwork/gonode/common/storage"
	"github.com/pastelnetwork/gonode/common/storage/memory"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
	"golang.org/x/crypto/sha3"
)

var (
	defaultNetworkAddr                   = "0.0.0.0"
	defaultNetworkPort                   = 4445
	defaultRefreshTime                   = time.Second * 3600
	defaultPingTime                      = time.Second * 10
	defaultCleanupInterval               = time.Minute * 2
	defaultDisabledKeyExpirationInterval = time.Minute * 30
)

const maxIterations = 5

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
	ht             *HashTable       // the hashtable for routing
	options        *Options         // the options of DHT
	network        *Network         // the network of DHT
	store          Store            // the storage of DHT
	metaStore      MetaStore        // the meta storage of DHT
	done           chan struct{}    // distributed hash table is done
	cache          storage.KeyValue // store bad bootstrap addresses
	pastelClient   pastel.Client
	externalIP     string
	mtx            sync.Mutex
	authHelper     *AuthHelper
	ignorelist     *BanList
	replicationMtx sync.RWMutex
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
func NewDHT(ctx context.Context, store Store, metaStore MetaStore, pc pastel.Client, secInfo *alts.SecInfo, options *Options) (*DHT, error) {
	// validate the options, if it's invalid, set them to default value
	if options.IP == "" {
		options.IP = defaultNetworkAddr
	}
	if options.Port <= 0 {
		options.Port = defaultNetworkPort
	}

	s := &DHT{
		metaStore:      metaStore,
		store:          store,
		options:        options,
		pastelClient:   pc,
		done:           make(chan struct{}),
		cache:          memory.NewKeyValue(),
		ignorelist:     NewBanList(ctx),
		replicationMtx: sync.RWMutex{},
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
	network, err := NewNetwork(ctx, s, ht.self, nil, s.authHelper)
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
	go s.startDisabledKeysCleanupWorker(ctx)

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

func (s *DHT) retryStore(ctx context.Context, key []byte, data []byte, typ int) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 1 * time.Minute
	b.InitialInterval = 200 * time.Millisecond

	return backoff.Retry(backoff.Operation(func() error {
		return s.store.Store(ctx, key, data, typ, true)
	}), b)
}

// Store the data into the network
func (s *DHT) Store(ctx context.Context, data []byte, typ int) (string, error) {
	key := s.hashKey(data)

	retKey := base58.Encode(key)
	// store the key to local storage
	if err := s.retryStore(ctx, key, data, typ); err != nil {
		log.WithContext(ctx).WithError(err).Error("local data store failure after retries")
		return "", fmt.Errorf("retry store data to local storage: %v", err)
	}

	if _, err := s.iterate(ctx, IterateStore, key, data, typ); err != nil {
		log.WithContext(ctx).WithError(err).Error("iterate data store failure")
		return "", fmt.Errorf("iterative store data: %v", err)
	}

	return retKey, nil
}

// StoreBatch will store a batch of values with their SHA256 hash as the key
func (s *DHT) StoreBatch(ctx context.Context, values [][]byte, typ int) error {
	var (
		counter int32
	)

	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}

	log.WithContext(ctx).WithField("taskID", taskID).WithField("records", len(values)).Info("store db batch begin")
	if err := s.store.StoreBatch(ctx, values, typ, true); err != nil {
		return fmt.Errorf("store batch: %v", err)
	}
	log.WithContext(ctx).WithField("taskID", taskID).Info("store db batch done,store network batch begin")

	// Launch the workers

	var wg sync.WaitGroup

	for i := 0; i < len(values); i++ {
		// Wait until we have a free spot
		for {
			if atomic.LoadInt32(&counter) < 2000 {
				break
			}
			// Sleep for a bit to prevent 100% CPU usage
			time.Sleep(1 * time.Second)
		}

		wg.Add(1)
		go func(val []byte, typ int) {
			// Increment the counter
			atomic.AddInt32(&counter, 1)

			key := s.hashKey(val)
			_, err := s.iterate(ctx, IterateStore, key, val, typ)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("iterate data batch store failure")
			}

			// Decrement the counter when done
			atomic.AddInt32(&counter, -1)
			wg.Done()
		}(values[i], typ)
	}
	wg.Wait()

	log.WithContext(ctx).WithField("taskID", taskID).Info("store network batch workers done")

	return nil
}

// Retrieve data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (s *DHT) Retrieve(ctx context.Context, key string, localOnly ...bool) ([]byte, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return nil, fmt.Errorf("invalid key: %v", key)
	}

	dbKey := hex.EncodeToString(decoded)
	if s.metaStore != nil {
		if err := s.metaStore.Retrieve(ctx, dbKey); err == nil {
			return nil, fmt.Errorf("key is disabled: %v", key)
		}
	}

	// retrieve the key/value from local storage
	value, err := s.store.Retrieve(ctx, decoded)
	if err == nil && len(value) > 0 {
		return value, nil
	}

	// if local only option is set, do not search just return error
	if len(localOnly) > 0 && localOnly[0] {
		return nil, fmt.Errorf("local-only failed to get properly: " + err.Error())
	}

	// if not found locally, iterative find value from kademlia network
	peerValue, err := s.iterate(ctx, IterateFindValue, decoded, nil, 0)
	if err != nil {
		return nil, errors.Errorf("retrieve from peer: %w", err)
	}
	if len(peerValue) > 0 {
		log.WithContext(ctx).WithField("key", dbKey).WithField("data len", len(peerValue)).Info("Not found locally, retrieved from other nodes")
	} else {
		log.WithContext(ctx).WithField("key", dbKey).Info("Not found locally, not found in other nodes")
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
	if s.store == nil {
		return nil, fmt.Errorf("store is nil")
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

// GetValueFromNode get values from node
func (s *DHT) GetValueFromNode(ctx context.Context, target []byte, n *Node) ([]byte, error) {
	messageType := FindValue
	data := &FindValueRequest{Target: target}

	request := s.newMessage(messageType, n, data)
	// send the request and receive the response
	cctx, ccancel := context.WithTimeout(ctx, time.Second*5)
	defer ccancel()

	response, err := s.network.Call(cctx, request, false)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Errorf("network call request %s failed", request.String())
		return nil, fmt.Errorf("network call request %s failed: %w", request.String(), err)
	}

	v, ok := response.Data.(*FindValueResponse)
	if ok && v.Status.Result == ResultOk && len(v.Value) > 0 {
		return v.Value, nil
	}

	return nil, fmt.Errorf("claim to have value but not found - %s - node: %s", response.String(), n.String())
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
				response, err := s.network.Call(ctx, request, false)
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
func (s *DHT) iterate(ctx context.Context, iterativeType int, target []byte, data []byte, typ int) ([]byte, error) {
	if iterativeType == IterateFindValue {
		return s.iterateFindValue(ctx, iterativeType, target)
	}

	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}

	sKey := hex.EncodeToString(target)

	igList := s.ignorelist.ToNodeList()
	// find the closest contacts for the target node from local route tables
	nl, _ := s.ht.closestContacts(Alpha, target, igList)
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
		hashedID, _ := utils.Sha3256hash(s.ht.self.ID)
		hashedTargetID, _ := utils.Sha3256hash(target)

		bucket := s.ht.bucketIndex(hashedID, hashedTargetID)
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
	maxIterations := 5 // Adjust the maximum iterations as needed

	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Debug("begin iteration")

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

				default:
					log.WithContext(ctx).WithField("type", response.MessageType).Error("unknown message type")
				}
			}

			// Stop search if no more nodes to contact
			if !searchRest && len(nl.Nodes) == 0 {
				log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("search stopped")
				return nil, nil
			}

			// Sort the nodes in the node list
			nl.Comparator = target
			nl.Sort()

			log.P2P().WithContext(ctx).Debugf("id: %v, iterate %d, sorted nodes: %v", base58.Encode(s.ht.self.ID), iterativeType, nl.String())

			switch iterativeType {
			case IterateFindNode:
				// If closestNode is unchanged
				if bytes.Equal(nl.Nodes[0].ID, closestNode.ID) || searchRest {
					if !searchRest {
						// Search all the remaining nodes
						searchRest = true
						continue
					}
					return nil, nil
				}
				// Update the closest node
				closestNode = nl.Nodes[0]

			case IterateStore:
				// Store the value to the node list
				if err := s.storeToAlphaNodes(ctx, nl, data, typ, taskID); err != nil {
					log.WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Error("could not store value to remaining network")
				}

				return nil, nil
			}

		}
	}
	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("finish iteration without results")
	return nil, nil
}

func (s *DHT) handleResponses(ctx context.Context, responses <-chan *Message, nl *NodeList) (*NodeList, []byte) {
	for response := range responses {
		s.addNode(ctx, response.Sender)
		if response.MessageType == FindNode || response.MessageType == StoreData {
			v, ok := response.Data.(*FindNodeResponse)
			if ok && v.Status.Result == ResultOk && len(v.Closest) > 0 {
				nl.AddNodes(v.Closest)
			}
		} else if response.MessageType == FindValue {
			v, ok := response.Data.(*FindValueResponse)
			if ok {
				if v.Status.Result == ResultOk && len(v.Value) > 0 {
					log.P2P().WithContext(ctx).Info("iterate found value from network")
					return nl, v.Value
				} else if len(v.Closest) > 0 {
					nl.AddNodes(v.Closest)
				}
			}
		}
	}

	return nl, nil
}

func (s *DHT) iterateFindValue(ctx context.Context, iterativeType int, target []byte) ([]byte, error) {
	// If task_id is available, use it for logging - This helps with debugging
	val := ctx.Value(log.TaskIDKey)
	taskID := ""
	if val != nil {
		taskID = fmt.Sprintf("%v", val)
	}

	// for logging, helps with debugging
	sKey := hex.EncodeToString(target)

	// the nodes which are unreachable right now are stored in 'ignore list'- we want to avoid hitting them repeatedly
	igList := s.ignorelist.ToNodeList()

	// nl will have the closest nodes to the target value, it will ignore the nodes in igList
	nl, _ := s.ht.closestContacts(Alpha, target, igList)
	if len(igList) > 0 {
		log.P2P().WithContext(ctx).WithField("nodes", nl.String()).WithField("ignored", s.ignorelist.String()).Info("closest contacts")
	}

	// if no nodes are found, return - this is a corner case and should not happen in practice
	if nl.Len() == 0 {
		return nil, nil
	}

	searchRest := false
	// keep track of contacted nodes so that we don't hit them again
	contacted := make(map[string]bool)
	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Debug("begin iteration")

	var closestNode *Node
	var iterationCount int
	for iterationCount = 0; iterationCount < maxIterations; iterationCount++ {
		log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("nl", nl.Len()).Debugf("begin find value - target: %v , nodes: %v", sKey, nl.String())

		if nl.Len() == 0 {
			log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).WithField("iteration count", iterationCount).Error("nodes list length is 0")
			return nil, nil
		}

		// if the closest node is the same as the last iteration  and we don't want to search rest of nodes, we are done
		if !searchRest && (closestNode != nil && bytes.Equal(nl.Nodes[0].ID, closestNode.ID)) {
			log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).WithField("iteration count", iterationCount).Debug("closest node is the same as the last iteration")
			return nil, nil
		}

		closestNode = nl.Nodes[0]
		responses := s.doMultiWorkers(ctx, iterativeType, target, nl, contacted, searchRest)
		var value []byte
		nl, value = s.handleResponses(ctx, responses, nl)
		if len(value) > 0 {
			return value, nil
		}

		nl.Sort()

		log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Infof("iteration: %v, nodes: %v", iterationCount, nl.Len())
	}

	log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Info("finished iterations without results")
	return nil, nil
}

func (s *DHT) sendReplicateData(ctx context.Context, n *Node, request *ReplicateDataRequest) (*ReplicateDataResponse, error) {
	// new a request message
	reqMsg := s.newMessage(Replicate, n, request)

	rspMsg, err := s.network.Call(ctx, reqMsg, true)
	if err != nil {
		return nil, errors.Errorf("replicate network call: %w", err)
	}

	response, ok := rspMsg.Data.(*ReplicateDataResponse)
	if !ok {
		return nil, errors.New("invalid ReplicateDataResponse")
	}

	return response, nil
}

func (s *DHT) sendStoreData(ctx context.Context, n *Node, request *StoreDataRequest) (*StoreDataResponse, error) {
	// new a request message
	reqMsg := s.newMessage(StoreData, n, request)

	rspMsg, err := s.network.Call(ctx, reqMsg, false)
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
	hashedID, _ := utils.Sha3256hash(s.ht.self.ID)
	hashedIncomingID, _ := utils.Sha3256hash(node.ID)

	index := s.ht.bucketIndex(hashedID, hashedIncomingID)

	s.ht.mutex.Lock()
	defer s.ht.mutex.Unlock()
	// 1. if the node is existed, refresh the node to the end of bucket
	if s.ht.hasBucketNode(index, node.ID) {
		s.ht.refreshNode(node.ID)
		return nil
	}

	if err := s.updateReplicationNode(ctx, node.ID, node.IP, node.Port, true); err != nil {
		log.P2P().WithContext(ctx).WithField("node-id", string(node.ID)).WithField("node-ip", node.IP).WithError(err).Error("update replication node failed")
	}

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
		response, err := s.network.Call(ctx, request, false)
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
	nodeList, _ := s.ht.closestContacts(n, base58.Decode(key), ignores)

	return nodeList.Nodes
}

// LocalStore the data into the network
func (s *DHT) LocalStore(ctx context.Context, key string, data []byte) (string, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return "", fmt.Errorf("invalid key: %v", key)
	}

	// store the key to local storage
	if err := s.retryStore(ctx, decoded, data, 0); err != nil {
		log.WithContext(ctx).WithError(err).Error("local data store failure after retries")
		return "", fmt.Errorf("retry store data to local storage: %v", err)
	}

	return key, nil
}

func (s *DHT) storeToAlphaNodes(ctx context.Context, nl *NodeList, data []byte, typ int, taskID string) error {
	storeCount := int32(0)
	alphaCh := make(chan bool, Alpha)

	// Start by sending the first Alpha requests in parallel
	for i := 0; i < Alpha && i < nl.Len(); i++ {
		n := nl.Nodes[i]

		if s.ignorelist.Banned(n) {
			continue
		}

		go func(n *Node) {
			logEntry := log.P2P().WithContext(ctx).WithField("node", n).WithField("task_id", taskID)

			request := &StoreDataRequest{Data: data, Type: typ}
			response, err := s.sendStoreData(ctx, n, request)
			if err != nil {
				logEntry.WithError(err).Error("send store data failed")
				alphaCh <- false
			} else if response.Status.Result != ResultOk {
				logEntry.WithError(errors.New(response.Status.ErrMsg)).Error("reply store data failed")
				alphaCh <- false
			} else {
				atomic.AddInt32(&storeCount, 1)
				alphaCh <- true
			}
		}(n)
	}
	skey, _ := utils.Sha3256hash(data)

	// Collect results from parallel requests
	for i := 0; i < Alpha && i < len(nl.Nodes); i++ {
		<-alphaCh
		if atomic.LoadInt32(&storeCount) >= int32(Alpha) {
			nl.TopN(Alpha)
			log.WithContext(ctx).WithField("task_id", taskID).WithField("skey", hex.EncodeToString(skey)).WithField("closest 6 nodes", nl.String()).
				WithField("len-total-nodes", nl.Len()).Info("store data to alpha nodes success")
			return nil
		}
	}

	finalStoreCount := atomic.LoadInt32(&storeCount)
	// If storeCount is still < Alpha, send requests sequentially until it reaches Alpha
	for i := Alpha; i < len(nl.Nodes); i++ {
		if finalStoreCount >= int32(Alpha) {
			break
		}

		logEntry := log.P2P().WithContext(ctx).WithField("node", nl.Nodes[i]).WithField("task_id", taskID)
		n := nl.Nodes[i]

		if s.ignorelist.Banned(n) {
			logEntry.Info("ignore node as its continuous failed count is above threshold")
			continue
		}

		request := &StoreDataRequest{Data: data, Type: typ}
		response, err := s.sendStoreData(ctx, n, request)
		if err != nil {
			logEntry.WithError(err).Error("send store data failed")
		} else if response.Status.Result != ResultOk {
			logEntry.WithError(errors.New(response.Status.ErrMsg)).Error("reply store data failed")
		} else {
			finalStoreCount++
		}
	}

	if finalStoreCount >= int32(Alpha) {
		log.WithContext(ctx).WithField("task_id", taskID).WithField("len-total-nodes", nl.Len()).WithField("skey", hex.EncodeToString(skey)).Info("store data to alpha nodes success")
		return nil
	}
	log.WithContext(ctx).WithField("task_id", taskID).WithField("store-count", finalStoreCount).WithField("skey", hex.EncodeToString(skey)).Info("store data to alpha nodes failed")

	return fmt.Errorf("store data to alpha nodes failed, only %d nodes stored", finalStoreCount)
}

func (s *DHT) startDisabledKeysCleanupWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("disabled keys cleanup worker started")

	for {
		select {
		case <-time.After(defaultCleanupInterval):
			s.cleanupDisabledKeys(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing disabled keys cleanup worker")
			return nil
		}
	}
}

func (s *DHT) cleanupDisabledKeys(ctx context.Context) error {
	if s.metaStore == nil {
		return nil
	}

	from := time.Now().Add(-1 * defaultDisabledKeyExpirationInterval)
	disabledKeys, err := s.metaStore.GetDisabledKeys(from)
	if err != nil {
		return errors.Errorf("get disabled keys: %w", err)
	}

	for i := 0; i < len(disabledKeys); i++ {
		dec, err := hex.DecodeString(disabledKeys[i].Key)
		if err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("decode disabled key failed")
			continue
		}
		s.metaStore.Delete(ctx, dec)
	}

	return nil
}
