package kademlia

import (
	"bytes"
	"context"
	"math"
	"sort"
	"time"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/errors"

	b58 "github.com/jbenet/go-base58"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/crypto"
	"github.com/pastelnetwork/gonode/p2p/kademlia/dao"
)

// DHT represents the state of the local node in the distributed hash table
type DHT struct {
	ht         *hashTable
	options    *Options
	networking networking
	store      dao.Key
}

// Options contains configuration options for the local node
type Options struct {
	ID []byte

	// The local IPv4 or IPv6 address
	IP string

	// The local port to listen for connections on
	Port int

	// Whether or not to use the STUN protocol to determine public IP and Port
	// May be necessary if the node is behind a NAT
	UseStun bool

	// Specifies the the host of the STUN server. If left empty will use the
	// default specified in go-stun.
	StunAddr string

	// A logger interface
	Logger log.Logger

	// The nodes being used to bootstrap the network. Without a bootstrap
	// node there is no way to connect to the network. NetworkNodes can be
	// initialized via dht.NewNetworkNode()
	BootstrapNodes []*NetworkNode

	// The time after which a key/value pair expires;
	// this is a time-to-live (TTL) from the original publication date
	TExpire time.Duration

	// Seconds after which an otherwise unaccessed bucket must be refreshed
	TRefresh time.Duration

	// The interval between Kademlia replication events, when a node is
	// required to publish its entire database
	TReplicate time.Duration

	// The time after which the original publisher must
	// republish a key/value pair. Currently not implemented.
	TRepublish time.Duration

	// The maximum time to wait for a response from a node before discarding
	// it from the bucket
	TPingMax time.Duration

	// The maximum time to wait for a response to any message
	TMsgTimeout time.Duration
}

// NewDHT initializes a new DHT node. A store and options struct must be
// provided.
func NewDHT(ctx context.Context, store dao.Key, options *Options) (*DHT, error) {
	dht := &DHT{}

	dht.options = options

	ht, err := newHashTable(options)
	if err != nil {
		return nil, err
	}

	dht.store = store
	dht.ht = ht
	dht.networking = &realNetworking{}

	if err = store.Init(ctx); err != nil {
		return nil, err
	}

	if options.TExpire == 0 {
		options.TExpire = time.Second * 86410
	}

	if options.TRefresh == 0 {
		options.TRefresh = time.Second * 3600
	}

	if options.TReplicate == 0 {
		options.TReplicate = time.Second * 3600
	}

	if options.TRepublish == 0 {
		options.TRepublish = time.Second * 86400
	}

	if options.TPingMax == 0 {
		options.TPingMax = time.Second * 1
	}

	if options.TMsgTimeout == 0 {
		options.TMsgTimeout = time.Second * 2
	}

	return dht, nil
}

func (dht *DHT) getExpirationTime(key []byte) time.Time {
	bucket := getBucketIndexFromDifferingBit(key, dht.ht.Self.ID)
	var total int
	for i := 0; i < bucket; i++ {
		total += dht.ht.getTotalNodesInBucket(i)
	}
	closer := dht.ht.getAllNodesInBucketCloserThan(bucket, key)
	score := total + len(closer)

	if score == 0 {
		score = 1
	}

	if score > k {
		return time.Now().Add(dht.options.TExpire)
	}

	day := dht.options.TExpire
	seconds := day.Nanoseconds() * int64(math.Exp(float64(k/score)))
	dur := time.Second * time.Duration(seconds)
	return time.Now().Add(dur)
}

// Store stores data on the network. This will trigger an iterateStore message.
// The base58 encoded identifier will be returned if the store is successful.
func (dht *DHT) Store(ctx context.Context, data []byte) (id string, err error) {
	key := crypto.GetKey(data)
	expiration := dht.getExpirationTime(key)
	replication := time.Now().Add(dht.options.TReplicate)
	if err = dht.store.Store(ctx, data, replication, expiration, true); err != nil {
		return "", err
	}

	if _, _, err = dht.iterate(ctx, iterateStore, key[:], data); err != nil {
		return "", err
	}
	str := b58.Encode(key)
	return str, nil
}

// Get retrieves data from the networking using key. Key is the base58 encoded
// identifier of the data.
func (dht *DHT) Get(ctx context.Context, key string) (data []byte, found bool, err error) {
	keyBytes := b58.Decode(key)
	value, exists := dht.store.Retrieve(ctx, keyBytes)

	// if len(keyBytes) != k {
	// 	return nil, false, errors.New("Invalid key")
	// }

	if !exists {
		var err error
		value, _, err = dht.iterate(ctx, iterateFindValue, keyBytes, nil)
		if err != nil {
			return nil, false, err
		}
		if value != nil {
			exists = true
		}
	}

	return value, exists, nil
}

// NumNodes returns the total number of nodes stored in the local routing table
func (dht *DHT) NumNodes() int {
	return dht.ht.totalNodes()
}

// GetSelfID returns the base58 encoded identifier of the local node
func (dht *DHT) GetSelfID() string {
	str := b58.Encode(dht.ht.Self.ID)
	return str
}

// GetNetworkAddr returns the publicly accessible IP and Port of the local
// node
func (dht *DHT) GetNetworkAddr() string {
	return dht.networking.getNetworkAddr()
}

// UseStun checks if DHT is using STUN.
func (dht *DHT) UseStun() bool {
	return dht.options.UseStun
}

// CreateSocket attempts to open a UDP socket on the port provided to options
func (dht *DHT) CreateSocket() error {
	ip := dht.options.IP
	port := dht.options.Port

	netMsgInit()
	dht.networking.init(dht.ht.Self)

	publicHost, publicPort, err := dht.networking.createSocket(ip, port, dht.options.UseStun, dht.options.StunAddr)
	if err != nil {
		return err
	}

	if dht.options.UseStun {
		dht.ht.setSelfAddr(publicHost, publicPort)
	}

	return nil
}

// Listen begins listening on the socket for incoming messages
func (dht *DHT) Listen(ctx context.Context) error {
	if !dht.networking.isInitialized() {
		return errors.New("socket not created")
	}

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		return dht.listen(ctx)
	})
	group.Go(func() error {
		return dht.timers(ctx)
	})
	group.Go(func() error {
		return dht.networking.listen(ctx)
	})
	return group.Wait()
}

// Bootstrap attempts to bootstrap the network using the BootstrapNodes provided
// to the Options struct. This will trigger an iterativeFindNode to the provided
// BootstrapNodes.
func (dht *DHT) Bootstrap(ctx context.Context) error {
	if len(dht.options.BootstrapNodes) == 0 {
		return nil
	}
	expectedResponses := []*expectedResponse{}

	for _, bn := range dht.options.BootstrapNodes {
		query := &message{}
		query.Sender = dht.ht.Self
		query.Receiver = bn
		query.Type = messageTypePing
		if bn.ID == nil {
			res, err := dht.networking.sendMessage(ctx, query, true, -1)
			if err != nil {
				continue
			}
			expectedResponses = append(expectedResponses, res)
		} else {
			node := newNode(bn)
			if err := dht.addNode(ctx, node); err != nil {
				return err
			}
		}
	}

	numExpectedResponses := len(expectedResponses)

	if numExpectedResponses > 0 {
		group, _ := errgroup.WithContext(ctx)

		for _, r := range expectedResponses {
			r := r
			group.Go(func() error {
				select {
				case result := <-r.ch:
					// If result is nil, channel was closed
					if result != nil {
						if err := dht.addNode(ctx, newNode(result.Sender)); err != nil {
							return err
						}
					}
					return nil
				case <-time.After(dht.options.TMsgTimeout):
					dht.networking.cancelResponse(r)
					return nil
				case <-ctx.Done():
					return nil
				}
			})
		}

		if err := group.Wait(); err != nil {
			return err
		}
	}

	if dht.NumNodes() > 0 {
		_, _, err := dht.iterate(ctx, iterateFindNode, dht.ht.Self.ID, nil)
		return err
	}

	return nil
}

// Disconnect will trigger a disconnect from the network. All underlying sockets
// will be closed.
func (dht *DHT) Disconnect() error {
	// TODO if .CreateSocket() is called, but .Listen() is never called, we
	// don't provide a way to close the socket
	return dht.networking.disconnect()
}

// Iterate does an iterative search through the network. This can be done
// for multiple reasons. These reasons include:
//     iterativeStore - Used to store new information in the network.
//     iterativeFindNode - Used to bootstrap the network.
//     iterativeFindValue - Used to find a value among the network given a key.
func (dht *DHT) iterate(ctx context.Context, t int, target []byte, data []byte) (value []byte, closest []*NetworkNode, err error) {
	sl := dht.ht.getClosestContacts(alpha, target, []*NetworkNode{})

	// We keep track of nodes contacted so far. We don't contact the same node
	// twice.
	var contacted = make(map[string]bool)

	// According to the Kademlia white paper, after a round of FIND_NODE RPCs
	// fails to provide a node closer than closestNode, we should send a
	// FIND_NODE RPC to all remaining nodes in the shortlist that have not
	// yet been contacted.
	queryRest := false

	// We keep a reference to the closestNode. If after performing a search
	// we do not find a closer node, we stop searching.
	if len(sl.Nodes) == 0 {
		return nil, nil, nil
	}

	closestNode := sl.Nodes[0]

	if t == iterateFindNode {
		bucket := getBucketIndexFromDifferingBit(target, dht.ht.Self.ID)
		dht.ht.resetRefreshTimeForBucket(bucket)
	}

	removeFromShortlist := []*NetworkNode{}

	for {
		expectedResponses := []*expectedResponse{}
		numExpectedResponses := 0

		// Next we send messages to the first (closest) alpha nodes in the
		// shortlist and wait for a response

		for i, node := range sl.Nodes {
			// Contact only alpha nodes
			if i >= alpha && !queryRest {
				break
			}

			// Don't contact nodes already contacted
			if contacted[string(node.ID)] {
				continue
			}

			contacted[string(node.ID)] = true
			query := &message{}
			query.Sender = dht.ht.Self
			query.Receiver = node

			switch t {
			case iterateFindNode:
				query.Type = messageTypeFindNode
				queryData := &queryDataFindNode{}
				queryData.Target = target
				query.Data = queryData
			case iterateFindValue:
				query.Type = messageTypeFindValue
				queryData := &queryDataFindValue{}
				queryData.Target = target
				query.Data = queryData
			case iterateStore:
				query.Type = messageTypeFindNode
				queryData := &queryDataFindNode{}
				queryData.Target = target
				query.Data = queryData
			default:
				return nil, nil, errors.New("unknown iterate type")
			}

			// Send the async queries and wait for a response
			res, err := dht.networking.sendMessage(ctx, query, true, -1)
			if err != nil {
				// Node was unreachable for some reason. We will have to remove
				// it from the shortlist, but we will keep it in our routing
				// table in hopes that it might come back online in the future.
				removeFromShortlist = append(removeFromShortlist, query.Receiver)
				continue
			}

			expectedResponses = append(expectedResponses, res)
		}

		for _, n := range removeFromShortlist {
			sl.RemoveNode(n)
		}

		numExpectedResponses = len(expectedResponses)

		resultChan := make(chan (*message))
		group, _ := errgroup.WithContext(ctx)
		for _, r := range expectedResponses {
			r := r
			group.Go(func() error {
				select {
				case result := <-r.ch:
					if result == nil {
						// Channel was closed
						return nil
					}
					if err := dht.addNode(ctx, newNode(result.Sender)); err != nil {
						return err
					}
					resultChan <- result
					return nil
				case <-time.After(dht.options.TMsgTimeout):
					dht.networking.cancelResponse(r)
					return nil
				}
			})
		}

		var results []*message
		if numExpectedResponses > 0 {
		Loop:
			for {
				select {
				case result := <-resultChan:
					if result != nil {
						results = append(results, result)
					} else {
						numExpectedResponses--
					}
					if len(results) == numExpectedResponses {
						close(resultChan)
						break Loop
					}
				case <-time.After(dht.options.TMsgTimeout):
					close(resultChan)
					break Loop
				}
			}

			for _, result := range results {
				if result.Error != nil {
					sl.RemoveNode(result.Receiver)
					continue
				}
				switch t {
				case iterateFindNode:
					responseData := result.Data.(*responseDataFindNode)
					sl.AppendUniqueNetworkNodes(responseData.Closest)
				case iterateFindValue:
					responseData := result.Data.(*responseDataFindValue)
					// TODO When an iterativeFindValue succeeds, the initiator must
					// store the key/value pair at the closest node seen which did
					// not return the value.
					if responseData.Value != nil {
						return responseData.Value, nil, nil
					}
					sl.AppendUniqueNetworkNodes(responseData.Closest)
				case iterateStore:
					responseData := result.Data.(*responseDataFindNode)
					sl.AppendUniqueNetworkNodes(responseData.Closest)
				}
			}
		}

		if !queryRest && len(sl.Nodes) == 0 {
			return nil, nil, group.Wait()
		}

		sort.Sort(sl)

		// If closestNode is unchanged then we are done
		if bytes.Equal(sl.Nodes[0].ID, closestNode.ID) || queryRest {
			// We are done
			switch t {
			case iterateFindNode:
				if !queryRest {
					queryRest = true
					continue
				}
				return nil, sl.Nodes, nil
			case iterateFindValue:
				return nil, sl.Nodes, nil
			case iterateStore:
				for i, n := range sl.Nodes {
					if i >= k {
						return nil, nil, group.Wait()
					}

					query := &message{}
					query.Receiver = n
					query.Sender = dht.ht.Self
					query.Type = messageTypeStore
					queryData := &queryDataStore{}
					queryData.Data = data
					query.Data = queryData
					_, err := dht.networking.sendMessage(ctx, query, false, -1)
					if err != nil {
						return nil, nil, err
					}
				}
				return nil, nil, group.Wait()
			}
		} else {
			closestNode = sl.Nodes[0]
		}
	}
}

// addNode adds a node into the appropriate k bucket
// we store these buckets in big-endian order so we look at the bits
// from right to left in order to find the appropriate bucket
func (dht *DHT) addNode(ctx context.Context, node *node) error {
	index := getBucketIndexFromDifferingBit(dht.ht.Self.ID, node.ID)

	// Make sure node doesn't already exist
	// If it does, mark it as seen
	if dht.ht.doesNodeExistInBucket(index, node.ID) {
		return dht.ht.markNodeAsSeen(ctx, node.ID)
	}

	dht.ht.mutex.Lock()
	defer dht.ht.mutex.Unlock()

	bucket := dht.ht.RoutingTable[index]

	if len(bucket) == k {
		// If the bucket is full we need to ping the first node to find out
		// if it responds back in a reasonable amount of time. If not -
		// we may remove it
		n := bucket[0].NetworkNode
		query := &message{}
		query.Receiver = n
		query.Sender = dht.ht.Self
		query.Type = messageTypePing
		res, err := dht.networking.sendMessage(ctx, query, true, -1)
		if err != nil {
			bucket = append(bucket, node)
			bucket = bucket[1:]
		} else {
			select {
			case <-res.ch:
				return nil
			case <-time.After(dht.options.TPingMax):
				bucket = bucket[1:]
				bucket = append(bucket, node)
			}
		}
	} else {
		bucket = append(bucket, node)
	}

	dht.ht.RoutingTable[index] = bucket

	return nil
}

func (dht *DHT) timers(ctx context.Context) error {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			// Refresh
			for i := 0; i < b; i++ {
				if time.Since(dht.ht.getRefreshTimeForBucket(i)) > dht.options.TRefresh {
					id := dht.ht.getRandomIDFromBucket(k)
					dht.iterate(ctx, iterateFindNode, id, nil)
				}
			}

			// Replication
			keys, err := dht.store.GetAllKeysForReplication(ctx)
			if err != nil {
				return err
			}

			for _, key := range keys {
				value, _ := dht.store.Retrieve(ctx, key)
				dht.iterate(ctx, iterateStore, key, value)
			}

			// Expiration
			dht.store.ExpireKeys(ctx)
		case <-dht.networking.getDisconnect():
			t.Stop()
			dht.networking.timersFin()
			return nil
		}
	}
}

func (dht *DHT) listen(ctx context.Context) error {
	for {
		select {
		case msg := <-dht.networking.getMessage():
			if msg == nil {
				// Disconnected
				dht.networking.messagesFin()
				return nil
			}
			switch msg.Type {
			case messageTypeFindNode:
				data := msg.Data.(*queryDataFindNode)
				if err := dht.addNode(ctx, newNode(msg.Sender)); err != nil {
					return err
				}
				closest := dht.ht.getClosestContacts(k, data.Target, []*NetworkNode{msg.Sender})
				response := &message{IsResponse: true}
				response.Sender = dht.ht.Self
				response.Receiver = msg.Sender
				response.Type = messageTypeFindNode
				responseData := &responseDataFindNode{}
				responseData.Closest = closest.Nodes
				response.Data = responseData
				dht.networking.sendMessage(ctx, response, false, msg.ID)
			case messageTypeFindValue:
				data := msg.Data.(*queryDataFindValue)
				if err := dht.addNode(ctx, newNode(msg.Sender)); err != nil {
					return err
				}
				value, exists := dht.store.Retrieve(ctx, data.Target)
				response := &message{IsResponse: true}
				response.ID = msg.ID
				response.Receiver = msg.Sender
				response.Sender = dht.ht.Self
				response.Type = messageTypeFindValue
				responseData := &responseDataFindValue{}
				if exists {
					responseData.Value = value
				} else {
					closest := dht.ht.getClosestContacts(k, data.Target, []*NetworkNode{msg.Sender})
					responseData.Closest = closest.Nodes
				}
				response.Data = responseData
				dht.networking.sendMessage(ctx, response, false, msg.ID)
			case messageTypeStore:
				data := msg.Data.(*queryDataStore)
				if err := dht.addNode(ctx, newNode(msg.Sender)); err != nil {
					return err
				}
				key := crypto.GetKey(data.Data)
				expiration := dht.getExpirationTime(key)
				replication := time.Now().Add(dht.options.TReplicate)
				dht.store.Store(ctx, data.Data, replication, expiration, false)
			case messageTypePing:
				response := &message{IsResponse: true}
				response.Sender = dht.ht.Self
				response.Receiver = msg.Sender
				response.Type = messageTypePing
				dht.networking.sendMessage(ctx, response, false, msg.ID)
			}
		case <-dht.networking.getDisconnect():
			dht.networking.messagesFin()
			return nil
		}
	}
}
