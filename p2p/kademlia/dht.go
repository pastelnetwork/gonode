package kademlia

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
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
	"github.com/pastelnetwork/gonode/common/storage/rqstore"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/pastel"
)

var (
	defaultNetworkAddr                   = "0.0.0.0"
	defaultNetworkPort                   = 4445
	defaultRefreshTime                   = time.Second * 3600
	defaultPingTime                      = time.Second * 10
	defaultCleanupInterval               = time.Minute * 2
	defaultDisabledKeyExpirationInterval = time.Minute * 30
	defaultRedundantDataCleanupInterval  = 12 * time.Hour
	defaultDeleteDataInterval            = 11 * time.Hour
	delKeysCountThreshold                = 10
	lowSpaceThreshold                    = 50 // GB
	batchStoreSize                       = 2500
	storeSameSymbolsBatchConcurrency     = 1
	storeSymbolsBatchConcurrency         = 2.0
)

const maxIterations = 4

// DHT represents the state of the queries node in the distributed hash table
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
	rqstore        rqstore.Store
}

// Options contains configuration options for the queries node
type Options struct {
	ID []byte

	// The queries IPv4 or IPv6 address
	IP string

	// The queries port to listen for connections
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
func NewDHT(ctx context.Context, store Store, metaStore MetaStore, pc pastel.Client, secInfo *alts.SecInfo, options *Options, rqstore rqstore.Store) (*DHT, error) {
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
		rqstore:        rqstore,
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
	go s.startCleanupRedundantDataWorker(ctx)
	go s.startDeleteDataWorker(ctx)
	go s.startStoreSymbolsWorker(ctx)

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
	key, _ := utils.Sha3256hash(data)

	retKey := base58.Encode(key)
	// store the key to queries storage
	if err := s.retryStore(ctx, key, data, typ); err != nil {
		log.WithContext(ctx).WithError(err).Error("queries data store failure after retries")
		return "", fmt.Errorf("retry store data to queries storage: %v", err)
	}

	if _, err := s.iterate(ctx, IterateStore, key, data, typ); err != nil {
		log.WithContext(ctx).WithError(err).Error("iterate data store failure")
		return "", fmt.Errorf("iterative store data: %v", err)
	}

	return retKey, nil
}

// StoreBatch will store a batch of values with their SHA256 hash as the key
func (s *DHT) StoreBatch(ctx context.Context, values [][]byte, typ int, taskID string) error {
	log.WithContext(ctx).WithField("taskID", taskID).WithField("records", len(values)).Info("store db batch begin")
	if err := s.store.StoreBatch(ctx, values, typ, true); err != nil {
		return fmt.Errorf("store batch: %v", err)
	}
	log.WithContext(ctx).WithField("taskID", taskID).Info("store db batch done,store network batch begin")

	if err := s.IterateBatchStore(ctx, values, typ, taskID); err != nil {
		return fmt.Errorf("iterate batch store: %v", err)
	}

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

	// retrieve the key/value from queries storage
	value, err := s.store.Retrieve(ctx, decoded)
	if err == nil && len(value) > 0 {
		return value, nil
	}

	// if queries only option is set, do not search just return error
	if len(localOnly) > 0 && localOnly[0] {
		return nil, fmt.Errorf("queries-only failed to get properly: " + err.Error())
	}

	// if not found locally, iterative find value from kademlia network
	peerValue, err := s.iterate(ctx, IterateFindValue, decoded, nil, 0)
	if err != nil {
		return nil, errors.Errorf("retrieve from peer: %w", err)
	}
	if len(peerValue) > 0 {
		log.WithContext(ctx).WithField("key", dbKey).WithField("data len", len(peerValue)).Debug("Not found locally, retrieved from other nodes")
	} else {
		log.WithContext(ctx).WithField("key", dbKey).Debug("Not found locally, not found in other nodes")
	}

	return peerValue, nil
}

// Delete delete key in queries node
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
		log.P2P().WithContext(ctx).WithError(err).Debugf("network call request %s failed", request.String())
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
					log.P2P().WithContext(ctx).WithError(err).Debugf("network call request %s failed", request.String())
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

		// close the message channel
		close(responses)
	}()

	return responses
}

func (s *DHT) fetchAndAddLocalKeys(ctx context.Context, hexKeys []string, result *sync.Map, req int32) (count int32, err error) {
	batchSize := 5000

	// Process in batches
	for start := 0; start < len(hexKeys); start += batchSize {
		end := start + batchSize
		if end > len(hexKeys) {
			end = len(hexKeys)
		}

		batchHexKeys := hexKeys[start:end]

		log.WithFields(log.Fields{
			"batchSize": len(batchHexKeys),
			"totalKeys": len(hexKeys),
		}).Info("Processing batch of local keys")

		// Retrieve values for the current batch of local keys
		localValues, _, batchErr := s.store.RetrieveBatchValues(ctx, batchHexKeys)
		if batchErr != nil {
			log.WithField("error", batchErr).Error("Failed to retrieve batch values")
			err = fmt.Errorf("retrieve batch values (local): %v", batchErr)
			continue // Optionally continue with next batch or return depending on use case
		}

		// Populate the result map with the local values and count the found keys
		for i, val := range localValues {
			if len(val) > 0 {
				count++
				result.Store(batchHexKeys[i], val)
				if count >= req {
					return count, nil
				}
			}
		}
	}

	return count, err
}

// BatchRetrieve data from the networking using keys. Keys are the base58 encoded
func (s *DHT) BatchRetrieve(ctx context.Context, keys []string, required int32, txID string, localOnly ...bool) (result map[string][]byte, err error) {
	result = make(map[string][]byte) // the result of the batch retrieve - keys are b58 encoded (as received in request)
	var resMap sync.Map              // the result of the batch retrieve - keys are b58 encoded
	var foundLocalCount int32        // the number of values found so far

	hexKeys := make([]string, len(keys))                // the hex keys keys[i] = hex.EncodeToString(base58.Decode(keys[i]))
	globalClosestContacts := make(map[string]*NodeList) // This will store the global top 6 nodes for each symbol's hash
	hashes := make([][]byte, len(keys))                 // the hashes of the keys - hashes[i] = base58.Decode(keys[i])
	knownNodes := make(map[string]*Node)                // This will store the nodes we know about

	defer func() {
		// Transfer data from resMap to result
		resMap.Range(func(key, value interface{}) bool {
			hexKey := key.(string)
			valBytes := value.([]byte)

			k, err := hex.DecodeString(hexKey)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("key", hexKey).WithField("txid", txID).Error("failed to decode hex key in resMap.Range")
				return true
			}

			result[base58.Encode(k)] = valBytes

			return true
		})

		for key, value := range result {
			if len(value) == 0 {
				delete(result, key)
			}
		}
	}()

	// populate result map with required keys
	for _, key := range keys {
		result[key] = nil
	}

	self := &Node{ID: s.ht.self.ID, IP: s.externalIP, Port: s.ht.self.Port}
	self.SetHashedID()
	log.WithContext(ctx).WithField("self", self.String()).
		WithField("txid", txID).Info("batch retrieve")

	// populate hexKeys and hashes
	for i, key := range keys {
		decoded := base58.Decode(key)
		if len(decoded) != B/8 {
			return nil, fmt.Errorf("invalid key: %v", key)
		}
		hashes[i] = decoded
		hexKeys[i] = hex.EncodeToString(decoded)
	}
	log.WithContext(ctx).WithField("self", self.String()).WithField("txid", txID).Info("populated keys and hashes")

	// Add nodes from route table to known nodes map
	for _, node := range s.ht.nodes() {
		n := &Node{ID: node.ID, IP: node.IP, Port: node.Port}
		n.SetHashedID()
		knownNodes[string(node.ID)] = n

	}

	log.WithContext(ctx).WithField("txid", txID).Info("done looping over responses")

	// Calculate the local top 6 nodes for each value
	for i := range keys {
		// Calculate the local top 6 nodes for each value
		top6 := s.ht.closestContactsWithInlcudingNode(Alpha, hashes[i], s.ignorelist.ToNodeList(), nil)
		globalClosestContacts[keys[i]] = top6

		s.addKnownNodes(ctx, top6.Nodes, knownNodes)
	}

	log.WithContext(ctx).WithField("txid", txID).Info("closest contacts populated, fetching local keys now")

	// remove self from the map
	delete(knownNodes, string(self.ID))

	foundLocalCount, err = s.fetchAndAddLocalKeys(ctx, hexKeys, &resMap, required)
	if err != nil {
		return nil, fmt.Errorf("fetch and add local keys: %v", err)
	}
	log.WithContext(ctx).WithField("txid", txID).WithField("local-foundCount", foundLocalCount).Info("batch find values count")

	if foundLocalCount >= required {
		return result, nil
	}

	// We don't have enough values locally, so we need to fetch from the network
	batchSize := batchStoreSize
	var networkFound int32
	totalBatches := int(math.Ceil(float64(required) / float64(batchSize)))
	parallelBatches := int(math.Min(float64(totalBatches), storeSymbolsBatchConcurrency))

	semaphore := make(chan struct{}, parallelBatches)
	var wg sync.WaitGroup
	gctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.WithContext(ctx).WithField("txid", txID).WithField("parallel batches", parallelBatches).Info("begin iterate batch get values")
	// Process in batches
	for start := 0; start < len(keys); start += batchSize {
		end := start + batchSize
		if end > len(keys) {
			end = len(keys)
		}

		// Check for early termination
		if atomic.LoadInt32(&networkFound)+int32(foundLocalCount) >= int32(required) {
			break
		}

		wg.Add(1)
		semaphore <- struct{}{} // Acquire a semaphore slot before launching the goroutine

		go s.processBatch(gctx, keys[start:end], hexKeys[start:end], semaphore, &wg, globalClosestContacts, knownNodes, &resMap,
			required, foundLocalCount, &networkFound, cancel, txID)
	}

	log.WithContext(ctx).WithField("txid", txID).Info("called iterate batch get values - waiting for workers to finish")
	wg.Wait() // Wait for all goroutines to finish
	log.WithContext(ctx).WithField("txid", txID).Info("iterate batch get values workers done")

	return result, nil
}

func (s *DHT) processBatch(ctx context.Context, batchKeys []string, batchHexKeys []string, semaphore chan struct{}, wg *sync.WaitGroup,
	globalClosestContacts map[string]*NodeList, knownNodes map[string]*Node, resMap *sync.Map, required int32, foundLocalCount int32, networkFound *int32,
	cancel context.CancelFunc, txID string) {

	defer wg.Done()
	defer func() { <-semaphore }()

	for i := 0; i < maxIterations; i++ {
		// Early check if context is done to stop processing
		select {
		case <-ctx.Done():
			return
		default:
		}

		fetchMap := make(map[string][]int)
		for i, key := range batchKeys {
			fetchNodes := globalClosestContacts[key].Nodes
			for _, node := range fetchNodes {
				nodeID := string(node.ID)
				fetchMap[nodeID] = append(fetchMap[nodeID], i)
			}
		}

		log.WithContext(ctx).WithField("len(fetchMap)", len(fetchMap)).WithField("len(batchHexKeys)", len(batchHexKeys)).WithField("len(batchKeys)", len(batchKeys)).
			WithField("network-found", atomic.LoadInt32(networkFound)).WithField("txid", txID).Info("fetch map")

		// Iterate through the network to get the values for the current batch
		foundCount, newClosestContacts, batchErr := s.iterateBatchGetValues(ctx, knownNodes, batchKeys, batchHexKeys, fetchMap, resMap, required, foundLocalCount+atomic.LoadInt32(networkFound))
		if batchErr != nil {
			log.WithContext(ctx).WithError(batchErr).WithField("txid", txID).Error("iterate batch get values failed")
		}

		// Update the global counter for found values
		atomic.AddInt32(networkFound, int32(foundCount))

		// Check and propagate early termination
		if atomic.LoadInt32(networkFound)+int32(foundLocalCount) >= int32(required) {
			cancel() // Cancels the context, signaling other goroutines to stop
			break
		}

		// we now need to check if the nodes in the globalClosestContacts Map are still in the top 6
		// if not, we need to send calls to the newly found nodes to inquire about the top 6 nodes
		changed := false
		for key, nodesList := range newClosestContacts {
			if nodesList == nil || nodesList.Nodes == nil {
				continue
			}

			if globalClosestContacts[key] == nil || globalClosestContacts[key].Nodes == nil {
				log.WithContext(ctx).WithField("key", key).Warn("global contacts list doesn't have the key")
				continue
			}

			if !haveAllNodes(nodesList.Nodes, globalClosestContacts[key].Nodes) {
				log.WithContext(ctx).WithField("key", key).WithField("have", nodesList.String()).WithField("task-id", txID).WithField("got", globalClosestContacts[key].String()).Info("global closest contacts list changed in fetch!")
				changed = true
			}

			nodesList.AddNodes(globalClosestContacts[key].Nodes)
			nodesList.Sort()
			nodesList.TopN(Alpha)

			s.addKnownNodes(ctx, nodesList.Nodes, knownNodes)
			globalClosestContacts[key] = nodesList
		}

		if !changed {
			log.WithContext(ctx).WithField("iter", i).WithField("task-id", txID).Info("global closest contacts list did not change")
			break
		}

		if i == maxIterations-1 {
			log.WithContext(ctx).WithField("iter", i).WithField("task-id", txID).Warn("max iterations reached, still top 6 list was changed")
		}
	}
}

func (s *DHT) iterateBatchGetValues(ctx context.Context, nodes map[string]*Node, keys []string, hexKeys []string, fetchMap map[string][]int,
	resMap *sync.Map, req, alreadyFound int32) (int, map[string]*NodeList, error) {
	semaphore := make(chan struct{}, storeSameSymbolsBatchConcurrency) // Limit concurrency to 1
	closestContacts := make(map[string]*NodeList)
	var wg sync.WaitGroup
	contactsMap := make(map[string]map[string][]*Node)
	var firstErr error
	var mu sync.Mutex // To protect the firstErr
	foundCount := int32(0)

	gctx, cancel := context.WithCancel(ctx) // Create a cancellable context
	defer cancel()
	for nodeID, node := range nodes {
		if s.ignorelist.Banned(node) {
			log.WithContext(ctx).WithField("node", node.String()).Info("Ignore banned node in iterate batch get values")
			continue
		}

		contactsMap[nodeID] = make(map[string][]*Node)
		wg.Add(1)
		go func(node *Node, nodeID string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case <-gctx.Done():
				return
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			}

			indices := fetchMap[nodeID]
			requestKeys := make(map[string]KeyValWithClosest)
			for _, idx := range indices {
				if idx < len(hexKeys) {
					_, loaded := resMap.Load(hexKeys[idx]) //  check if key is already there in resMap
					if !loaded {
						requestKeys[hexKeys[idx]] = KeyValWithClosest{}
					}
				}
			}

			if len(requestKeys) == 0 {
				return
			}

			decompressedData, err := s.doBatchGetValuesCall(gctx, node, requestKeys)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			for k, v := range decompressedData {
				if len(v.Value) > 0 {
					_, loaded := resMap.LoadOrStore(k, v.Value)
					if !loaded {
						atomic.AddInt32(&foundCount, 1)
						if atomic.LoadInt32(&foundCount) >= int32(req-alreadyFound) {
							cancel() // Cancel context to stop other goroutines
							return
						}
					}
				} else {
					contactsMap[nodeID][k] = v.Closest
				}
			}
		}(node, nodeID)
	}

	wg.Wait()

	log.WithContext(ctx).WithField("found-count", atomic.LoadInt32(&foundCount)).Info("iterate batch get values done")

	if firstErr != nil {
		log.WithContext(ctx).WithError(firstErr).WithField("found-count", atomic.LoadInt32(&foundCount)).Error("encountered error in iterate batch get values")
	}

	for _, closestNodes := range contactsMap {
		for key, nodes := range closestNodes {
			comparator, err := hex.DecodeString(key)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("key", key).Error("failed to decode hex key in closestNodes.Range")
				return 0, nil, err
			}
			bkey := base58.Encode(comparator)

			if _, ok := closestContacts[bkey]; !ok {
				closestContacts[bkey] = &NodeList{Nodes: nodes, Comparator: comparator}
			} else {
				closestContacts[bkey].AddNodes(nodes)
			}
		}
	}

	for key, nodes := range closestContacts {
		nodes.Sort()
		nodes.TopN(Alpha)
		closestContacts[key] = nodes
	}

	return int(foundCount), closestContacts, firstErr
}

func (s *DHT) doBatchGetValuesCall(ctx context.Context, node *Node, requestKeys map[string]KeyValWithClosest) (map[string]KeyValWithClosest, error) {
	request := s.newMessage(BatchGetValues, node, &BatchGetValuesRequest{Data: requestKeys})
	response, err := s.network.Call(ctx, request, false)
	if err != nil {
		return nil, fmt.Errorf("network call request %s failed: %w", request.String(), err)
	}

	resp, ok := response.Data.(*BatchGetValuesResponse)
	if !ok {
		return nil, fmt.Errorf("invalid response type: %T", response.Data)
	}

	if resp.Status.Result != ResultOk {
		return nil, fmt.Errorf("response status: %v", resp.Status.ErrMsg)
	}

	return resp.Data, nil
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
	// find the closest contacts for the target node from queries route tables
	nl, _ := s.ht.closestContacts(Alpha, target, igList)
	if len(igList) > 0 {
		log.P2P().WithContext(ctx).WithField("nodes", nl.String()).WithField("ignored", s.ignorelist.String()).Info("closest contacts")
	}
	// if no closer node, stop search
	if nl.Len() == 0 {
		return nil, nil
	}
	log.P2P().WithContext(ctx).WithField("task_id", taskID).Debugf("type: %v, target: %v, nodes: %v", iterativeType, sKey, nl.String())

	// keep the closer node
	closestNode := nl.Nodes[0]
	// if it's a find node, reset the refresh timer
	if iterativeType == IterateFindNode {
		hashedTargetID, _ := utils.Sha3256hash(target)
		bucket := s.ht.bucketIndex(s.ht.self.HashedID, hashedTargetID)
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
					log.P2P().WithContext(ctx).Debug("iterate found value from network")
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

		log.P2P().WithContext(ctx).WithField("task_id", taskID).WithField("key", sKey).Debugf("iteration: %v, nodes: %v", iterationCount, nl.Len())
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
	node.SetHashedID()

	index := s.ht.bucketIndex(s.ht.self.HashedID, node.HashedID)

	if err := s.updateReplicationNode(ctx, node.ID, node.IP, node.Port, true); err != nil {
		log.P2P().WithContext(ctx).WithField("node-id", string(node.ID)).WithField("node-ip", node.IP).WithError(err).Error("update replication node failed")
	}

	if s.ht.hasBucketNode(index, node.ID) {
		s.ht.refreshNode(node.HashedID)
		return nil
	}

	log.WithContext(ctx).WithField("node", node.String()).Info("waiting for mutex")
	s.ht.mutex.Lock()
	log.WithContext(ctx).WithField("node", node.String()).Info("unlocked mutex")
	defer s.ht.mutex.Unlock()

	// 2. if the bucket is full, ping the first node
	bucket := s.ht.routeTable[index]
	if len(bucket) == K {
		first := bucket[0]

		var err error
		isBanned := s.ignorelist.Banned(first)
		if !isBanned {
			// new a ping request message
			request := s.newMessage(Ping, first, nil)
			// new a context with timeout
			ctx, cancel := context.WithTimeout(ctx, defaultPingTime)
			defer cancel()

			// invoke the request and handle the response
			_, err = s.network.Call(ctx, request, false)
			if err == nil {
				// refresh the node to the end of bucket
				bucket = bucket[1:]
				bucket = append(bucket, node)
				s.ht.routeTable[index] = bucket
				return nil
			}
		} else if isBanned || err != nil {
			s.ignorelist.IncrementCount(node)
			// the node is down, remove the node from bucket
			bucket = append(bucket, node)
			bucket = bucket[1:]

			// need to reset the route table with the bucket
			s.ht.routeTable[index] = bucket

			return first
		}

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

// NClosestNodesWithIncludingNodelist get n closest nodes to a key string with including node list
func (s *DHT) NClosestNodesWithIncludingNodelist(_ context.Context, n int, key string, ignores, includeNodeList []*Node) []*Node {
	list := s.ignorelist.ToNodeList()
	ignores = append(ignores, list...)

	for i := 0; i < len(includeNodeList); i++ {
		includeNodeList[i].SetHashedID()
	}

	nodeList := s.ht.closestContactsWithIncludingNodeList(n, base58.Decode(key), ignores, includeNodeList)

	return nodeList.Nodes
}

// LocalStore the data into the network
func (s *DHT) LocalStore(ctx context.Context, key string, data []byte) (string, error) {
	decoded := base58.Decode(key)
	if len(decoded) != B/8 {
		return "", fmt.Errorf("invalid key: %v", key)
	}

	// store the key to queries storage
	if err := s.retryStore(ctx, decoded, data, 0); err != nil {
		log.WithContext(ctx).WithError(err).Error("queries data store failure after retries")
		return "", fmt.Errorf("retry store data to queries storage: %v", err)
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
				WithField("len-total-nodes", nl.Len()).Debug("store data to alpha nodes success")
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
		log.WithContext(ctx).WithField("task_id", taskID).WithField("len-total-nodes", nl.Len()).WithField("skey", hex.EncodeToString(skey)).Debug("store data to alpha nodes success")
		return nil
	}
	log.WithContext(ctx).WithField("task_id", taskID).WithField("store-count", finalStoreCount).WithField("skey", hex.EncodeToString(skey)).Info("store data to alpha nodes failed")

	return fmt.Errorf("store data to alpha nodes failed, only %d nodes stored", finalStoreCount)
}

// remove node from appropriate k bucket
func (s *DHT) removeNode(ctx context.Context, node *Node) {
	// ensure this is not itself address
	if bytes.Equal(node.ID, s.ht.self.ID) {
		log.P2P().WithContext(ctx).Debug("trying to remove itself")
		return
	}
	node.SetHashedID()

	index := s.ht.bucketIndex(s.ht.self.HashedID, node.HashedID)

	if removed := s.ht.RemoveNode(index, node.ID); !removed {
		log.P2P().WithContext(ctx).Errorf("remove node %s not found in bucket %d", node.String(), index)
	} else {
		log.P2P().WithContext(ctx).Infof("removed node %s from bucket %d success", node.String(), index)
	}
}

func (s *DHT) addKnownNodes(ctx context.Context, nodes []*Node, knownNodes map[string]*Node) {
	for _, node := range nodes {
		if _, ok := knownNodes[string(node.ID)]; ok {
			continue
		}
		node.SetHashedID()
		knownNodes[string(node.ID)] = node

		s.addNode(ctx, node)
		bucket := s.ht.bucketIndex(s.ht.self.HashedID, node.HashedID)
		s.ht.resetRefreshTime(bucket)
	}
}

func (s *DHT) IterateBatchStore(ctx context.Context, values [][]byte, typ int, id string) error {
	globalClosestContacts := make(map[string]*NodeList)
	knownNodes := make(map[string]*Node)
	contacted := make(map[string]bool)
	hashes := make([][]byte, len(values))

	log.WithContext(ctx).WithField("task-id", id).WithField("keys", len(values)).Info("iterate batch store begin")
	for i := 0; i < len(values); i++ {
		target, _ := utils.Sha3256hash(values[i])
		hashes[i] = target
		top6 := s.ht.closestContactsWithInlcudingNode(Alpha, target, s.ignorelist.ToNodeList(), nil)

		globalClosestContacts[base58.Encode(target)] = top6
		s.addKnownNodes(ctx, top6.Nodes, knownNodes)
	}

	var changed bool
	var i int
	for {
		i++
		log.WithContext(ctx).WithField("task-id", id).WithField("iter", i).WithField("keys", len(values)).Info("iterate batch store begin")
		changed = false
		localClosestNodes := make(map[string]*NodeList)
		responses, atleastOneContacted := s.batchFindNode(ctx, hashes, knownNodes, contacted, id)

		if !atleastOneContacted {
			break
		}

		for response := range responses {
			if response.Error != nil {
				log.WithContext(ctx).WithError(response.Error).WithField("task-id", id).Error("batch find node failed on a node")
				continue
			}

			if response.Message == nil {
				continue
			}

			v, ok := response.Message.Data.(*BatchFindNodeResponse)
			if ok && v.Status.Result == ResultOk {
				for key, nodesList := range v.ClosestNodes {
					if nodesList != nil {
						nl, exists := localClosestNodes[key]
						if exists {
							nl.AddNodes(nodesList)
							localClosestNodes[key] = nl
						} else {
							localClosestNodes[key] = &NodeList{Nodes: nodesList, Comparator: base58.Decode(key)}
						}

						s.addKnownNodes(ctx, nodesList, knownNodes)
					}
				}
			}
		}

		// we now need to check if the nodes in the globalClosestContacts Map are still in the top 6
		// if yes, we can store the data to them
		// if not, we need to send calls to the newly found nodes to inquire about the top 6 nodes
		log.WithContext(ctx).WithField("task-id", id).WithField("iter", i).WithField("keys", len(values)).Info("check closest nodes")
		for key, nodesList := range localClosestNodes {
			if nodesList == nil {
				continue
			}

			nodesList.Comparator = base58.Decode(key)
			nodesList.Sort()
			nodesList.TopN(Alpha)
			s.addKnownNodes(ctx, nodesList.Nodes, knownNodes)

			if !haveAllNodes(nodesList.Nodes, globalClosestContacts[key].Nodes) {
				log.WithContext(ctx).WithField("key", key).WithField("have", nodesList.String()).WithField("task-id", id).WithField("got", globalClosestContacts[key].String()).Info("global closest contacts list changed!")
				changed = true
			}

			nodesList.AddNodes(globalClosestContacts[key].Nodes)
			nodesList.Sort()
			nodesList.TopN(Alpha)
			globalClosestContacts[key] = nodesList
		}
		log.WithContext(ctx).WithField("task-id", id).WithField("iter", i).WithField("keys", len(values)).Info("check closest nodes done")

		if !changed {
			log.WithContext(ctx).WithField("iter", i).WithField("task-id", id).Info("global closest contacts list did not change, we can now store the data")
			break
		}
	}

	// assume at this point, we have True\Golabl top 6 nodes for each symbol's hash stored in globalClosestContacts Map
	// we now need to store the data to these nodes

	storageMap := make(map[string][]int) // This will store the index of the data in the values array that needs to be stored to the node
	for i := 0; i < len(hashes); i++ {
		storageNodes := globalClosestContacts[base58.Encode(hashes[i])]
		for j := 0; j < len(storageNodes.Nodes); j++ {
			storageMap[string(storageNodes.Nodes[j].ID)] = append(storageMap[string(storageNodes.Nodes[j].ID)], i)
		}
	}

	requests := 0
	successful := 0

	storeResponses := s.batchStoreNetwork(ctx, values, knownNodes, storageMap, typ)
	for response := range storeResponses {
		requests++
		if response.Error != nil {
			sender := ""
			if response.Message != nil && response.Message.Sender != nil {
				sender = response.Message.Sender.String()
			}
			log.WithContext(ctx).WithField("node", sender).WithError(response.Error).Error("batch store failed on a node")
		}

		if response.Message == nil {
			continue
		}

		v, ok := response.Message.Data.(*StoreDataResponse)
		if ok && v.Status.Result == ResultOk {
			successful++
		} else {
			errMsg := "unknwon error"
			if v != nil {
				errMsg = v.Status.ErrMsg
			}

			log.WithContext(ctx).WithField("err", errMsg).WithField("task-id", id).Errorf("batch store to node %s failed", response.Message.Sender.String())
		}
	}

	if requests > 0 {
		successRate := float64(successful) / float64(requests) * 100
		if successRate >= 80 {
			log.WithContext(ctx).WithField("task-id", id).Infof("Successful store operations: %.2f%%", successRate)
			return nil
		} else {
			log.WithContext(ctx).WithField("task-id", id).Infof("Failed to achieve desired success rate, only: %.2f%%", successRate)
			return fmt.Errorf("failed to achieve desired success rate, only: %.2f%% successful", successRate)
		}
	}

	return fmt.Errorf("no store operations were performed")
}

func (s *DHT) batchStoreNetwork(ctx context.Context, values [][]byte, nodes map[string]*Node, storageMap map[string][]int, typ int) chan *MessageWithError {
	responses := make(chan *MessageWithError, len(nodes))
	semaphore := make(chan struct{}, 3) // Semaphore to limit concurrency to 3

	var wg sync.WaitGroup

	for key, node := range nodes {
		if s.ignorelist.Banned(node) {
			log.WithContext(ctx).WithField("node", node.String()).Info("Ignoring banned node in batch store network call")
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore
		go func(receiver *Node, key string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			select {
			case <-ctx.Done():
				responses <- &MessageWithError{Error: ctx.Err()}
				return
			default:
				keysToStore := storageMap[key]
				toStore := make([][]byte, len(keysToStore))
				totalBytes := 0
				for i, idx := range keysToStore {
					toStore[i] = values[idx]
					totalBytes += len(values[idx])
				}

				log.WithContext(ctx).WithField("keys", len(toStore)).WithField("size-before-compress", utils.BytesIntToMB(totalBytes)).Info("batch store to node")

				data := &BatchStoreDataRequest{Data: toStore, Type: typ}
				request := s.newMessage(BatchStoreData, receiver, data)
				response, err := s.network.Call(ctx, request, false)
				if err != nil {
					s.ignorelist.IncrementCount(receiver)
					log.P2P().WithContext(ctx).WithError(err).Debugf("network call batch store request %s failed", request.String())
					responses <- &MessageWithError{Error: err, Message: response}
					return
				}

				responses <- &MessageWithError{Message: response}
			}
		}(node, key)
	}

	wg.Wait()
	log.WithContext(ctx).Info("closing response channel")
	close(responses)
	log.WithContext(ctx).Info("closed response channel")

	return responses
}

func (s *DHT) batchFindNode(ctx context.Context, payload [][]byte, nodes map[string]*Node, contacted map[string]bool, txid string) (chan *MessageWithError, bool) {
	log.WithContext(ctx).WithField("task-id", txid).WithField("nodes-count", len(nodes)).Info("batch find node begin")

	responses := make(chan *MessageWithError, len(nodes))
	atleastOneContacted := false
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 20)

	for _, node := range nodes {
		if _, ok := contacted[string(node.ID)]; ok {
			continue
		}
		if s.ignorelist.Banned(node) {
			log.WithContext(ctx).WithField("node", node.String()).WithField("txid", txid).Info("Ignoring banned node in batch find call")
			continue
		}

		contacted[string(node.ID)] = true
		atleastOneContacted = true
		wg.Add(1)
		semaphore <- struct{}{}

		go func(receiver *Node) {
			defer wg.Done()
			defer func() { <-semaphore }()

			select {
			case <-ctx.Done():
				responses <- &MessageWithError{Error: ctx.Err()}
				return
			default:
				data := &BatchFindNodeRequest{HashedTarget: payload}
				request := s.newMessage(BatchFindNode, receiver, data)
				response, err := s.network.Call(ctx, request, false)
				if err != nil {
					s.ignorelist.IncrementCount(receiver)
					log.WithContext(ctx).WithError(err).WithField("node", receiver.String()).WithField("txid", txid).Warn("batch find node network call request failed")
					responses <- &MessageWithError{Error: err, Message: response}
					return
				}

				responses <- &MessageWithError{Message: response}
			}
		}(node)
	}
	wg.Wait()
	close(responses)
	log.WithContext(ctx).WithField("nodes-count", len(nodes)).WithField("len-resp", len(responses)).WithField("txid", txid).Info("batch find node done")

	return responses, atleastOneContacted
}
