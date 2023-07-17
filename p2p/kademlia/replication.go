package kademlia

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"encoding/gob"
	"encoding/hex"

	"github.com/DataDog/zstd"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

var (
	// defaultReplicationInterval is the default interval for replication.
	defaultReplicationInterval = time.Minute * 5

	// nodeShowUpDeadline is the time after which the node will considered to be permeant offline
	// we'll adjust the keys once the node is permeant offline
	nodeShowUpDeadline = time.Minute * 35

	// check for active & inactive nodes after this interval
	checkNodeActivityInterval = time.Minute * 3

	defaultFetchAndStoreInterval = time.Minute * 5
)

// StartReplicationWorker starts replication
func (s *DHT) StartReplicationWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("replication worker started")

	go s.checkNodeActivity(ctx)

	for {
		select {
		case <-time.After(defaultReplicationInterval):
			s.Replicate(ctx)
			log.P2P().WithContext(ctx).Error("no running replication worker")
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing replication worker")
			return nil
		}
	}
}

// StartFetchAndStoreWorker starts replication
func (s *DHT) StartFetchAndStoreWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("fetch and store worker started")

	for {
		select {
		case <-time.After(defaultFetchAndStoreInterval):
			s.FetchAndStore(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing fetch & store worker")
			return nil
		}
	}
}

func (s *DHT) updateReplicationNode(ctx context.Context, nodeID []byte, ip string, port int, isActive bool) error {
	if info, ok := s.nodeReplicationTimes[string(nodeID)]; ok {
		info.Active = isActive
		info.UpdatedAt = time.Now()
		info.IP = ip
		info.Port = port
		info.ID = nodeID

		if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
			log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).WithField("ip", ip).Error("failed to add replication info")
			return err
		}

	} else {
		info := domain.NodeReplicationInfo{
			UpdatedAt: time.Now(),
			Active:    isActive,
			IP:        ip,
			Port:      port,
			ID:        nodeID,
		}

		s.nodeReplicationTimes[string(nodeID)] = info

		if err := s.store.AddReplicationInfo(ctx, info); err != nil {
			log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).WithField("ip", ip).Error("failed to add replication info")
			return err
		}
	}

	return nil
}

func (s *DHT) updateLastReplicated(ctx context.Context, nodeID []byte, timestamp time.Time) error {
	info, ok := s.nodeReplicationTimes[string(nodeID)]
	if !ok {
		return errors.New("node not found")
	}

	info.LastReplicatedAt = &timestamp

	s.nodeReplicationTimes[string(nodeID)] = info

	if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
		log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).Error("failed to update replication info")
	}

	return nil
}

// Replicate is called periodically by the replication worker to replicate the data across the network by refreshing the buckets
// it iterates over the node replication Info map and replicates the data to the nodes that are active
func (s *DHT) Replicate(ctx context.Context) {
	s.replicationMtx.Lock()
	defer s.replicationMtx.Unlock()

	historicStart, err := s.store.GetOwnCreatedAt(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get own createdAt")
		historicStart = time.Now().Add(-24 * time.Hour * 180)
	}

	log.P2P().WithContext(ctx).Info("replicating data")

	for i := 0; i < B; i++ {
		if time.Since(s.ht.refreshTime(i)) > defaultRefreshTime {
			// refresh the bucket by iterative find node
			id := s.ht.randomIDFromBucket(K)
			if _, err := s.iterate(ctx, IterateFindNode, id, nil, 0); err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("replicate iterate find node failed")
			}
		}
	}

	to := time.Now()
	for nodeID, infoVar := range s.nodeReplicationTimes {
		info := infoVar

		logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", string(nodeID))
		if !info.Active {
			adjustNodeKeys := isNodeGoneAndShouldBeAdjusted(info.LastReplicatedAt, info.IsAdjusted)
			if adjustNodeKeys {
				if err := s.adjustNodeKeys(ctx, infoVar.CreatedAt, info); err != nil {
					logEntry.WithError(err).Error("failed to adjust node keys")
				} else {
					info.IsAdjusted = true
					info.UpdatedAt = time.Now()
					s.nodeReplicationTimes[string(nodeID)] = info

					if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
						logEntry.WithError(err).Error("failed to update replication info, set isAdjusted to true")
					} else {
						logEntry.Info("set isAdjusted to true")
					}
				}
			}

			logEntry.Info("replication node not active, skipping over it.")
			continue
		}

		from := historicStart
		if info.LastReplicatedAt != nil {
			from = *info.LastReplicatedAt
		}

		replicationKeys := s.store.GetKeysForReplication(ctx, from, to)
		logEntry.WithField("len-rep-keys", len(replicationKeys)).Info("count of replication keys to be checked")

		// Preallocate a slice with a capacity equal to the number of keys.
		closestContactKeys := make([][]byte, 0, len(replicationKeys))
		for i := 0; i < len(replicationKeys); i++ {
			key := replicationKeys[i]
			ignores := s.ignorelist.ToNodeList()
			nodeList := s.ht.closestContacts(Alpha, key, ignores)

			n := &Node{ID: []byte(nodeID), IP: info.IP, Port: info.Port}
			if nodeList.Exists(n) {
				// the node is supposed to hold this key as it's in the 6 closest contacts
				closestContactKeys = append(closestContactKeys, key)
			}
		}
		logEntry.WithField("len-rep-keys", len(closestContactKeys)).Info("closest contact keys count")

		if len(closestContactKeys) > 0 {
			data, err := compressKeys(closestContactKeys)
			if err != nil {
				logEntry.WithField("len-rep-keys", len(closestContactKeys)).WithError(err).Error("unable to compress keys - replication failed")
				continue
			}

			// TODO: Check if data size is bigger than 32 MB
			request := &ReplicateDataRequest{
				Keys: data,
			}

			n := &Node{ID: []byte(nodeID), IP: info.IP, Port: info.Port}
			response, err := s.sendReplicateData(ctx, n, request)
			if err != nil {
				logEntry.WithError(err).Error("send replicate data failed")
				continue
			} else if response.Status.Result != ResultOk {
				logEntry.WithError(errors.New(response.Status.ErrMsg)).Error("reply replicate data failed")
				continue
			}

			// Now closestContactKeys contains all the keys that are in the closest contacts.
			if err := s.updateLastReplicated(ctx, []byte(nodeID), to); err != nil {
				logEntry.Error("replicate update lastReplicated failed")
			} else {
				logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("expected-rep-keys", len(closestContactKeys)).Info("replicate update lastReplicated success")
			}
		}
	}

	log.P2P().WithContext(ctx).Info("Replication done")
}

func (s *DHT) adjustNodeKeys(ctx context.Context, from time.Time, info domain.NodeReplicationInfo) error {
	logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", string(info.ID))
	logEntry.WithField("from", from.String()).Info("adjusting node keys")

	replicationKeys := s.store.GetKeysForReplication(ctx, from, time.Now())
	logEntry = logEntry.WithField("len-rep-keys", len(replicationKeys))
	logEntry.Info("replication keys to be adjusted")

	/*storeMap := make(map[string][]Node)
	for _, key := range replicationKeys {

		// prepare ignored nodes list but remove the node we are adjusting
		// because we want to find if this node was supposed to hold this key
		ignores := s.ignorelist.ToNodeList()
		var updatedIgnored []*Node
		for _, ignore := range ignores {
			if string(ignore.ID) != string(info.ID) {
				updatedIgnored = append(updatedIgnored, ignore)
			}
		}

		// get closest contacts to the key
		nodeList := s.ht.closestContacts(Alpha, key, updatedIgnored)

		// check if the node that is gone was supposed to hold the key
		n := &Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}
		if !nodeList.Exists(n) {
			// the node is not supposed to hold this key as its not in 6 closest contacts
			continue
		}
		logEntry.WithField("key", hex.EncodeToString(key)).WithField("ip", info.IP).Info("adjust: this key is supposed to be hold by this node")

		storeMap[hex.EncodeToString(key)] = append(storeMap[hex.EncodeToString(key)], n)
	}
	*/
	info.IsAdjusted = true
	info.UpdatedAt = time.Now()

	if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
		return fmt.Errorf("replicate update isAdjusted failed: %v", err)
	}

	logEntry.Info("isAdjusted success")

	return nil
}

func isNodeGoneAndShouldBeAdjusted(lastReplicated *time.Time, isAlreadyAdjusted bool) bool {
	if lastReplicated == nil {
		return false
	}

	return time.Since(*lastReplicated) > nodeShowUpDeadline && !isAlreadyAdjusted
}

// checkNodeActivity keeps track of active nodes - the idea here is to ping nodes periodically and mark them as inactive if they don't respond
func (s *DHT) checkNodeActivity(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(checkNodeActivityInterval): // Adjust the interval as needed
			func() {
				if !utils.CheckInternetConnectivity() {
					log.WithContext(ctx).Info("no internet connectivity, not checking node activity")
				} else {
					s.replicationMtx.Lock()
					defer s.replicationMtx.Unlock()

					for nodeID, info := range s.nodeReplicationTimes {
						// new a ping request message
						node := &Node{
							ID:   []byte(nodeID),
							IP:   info.IP,
							Port: info.Port,
						}

						request := s.newMessage(Ping, node, nil)
						// new a context with timeout
						ctx, cancel := context.WithTimeout(ctx, defaultPingTime)
						defer cancel()

						// invoke the request and handle the response
						_, err := s.network.Call(ctx, request)
						if err != nil && info.Active {
							log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).
								Error("failed to ping node, setting node to inactive")

							// add node to ignore list
							// we maintain this list to avoid pinging nodes that are not responding
							s.ignorelist.IncrementCount(node)

							// mark node as inactive in database
							info.Active = false
							info.UpdatedAt = time.Now()
							s.nodeReplicationTimes[string(nodeID)] = info

							if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
								log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Error("failed to update replication info, node is inactive")
							}

						} else if err == nil {
							log.P2P().WithContext(ctx).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Info("node is active")

							if !info.Active {
								log.P2P().WithContext(ctx).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Info("node found to be active again")

								// remove node from ignore list
								s.ignorelist.Delete(node)

								// mark node as active in database
								info.Active = true
								info.IsAdjusted = false
								info.UpdatedAt = time.Now()
								s.nodeReplicationTimes[string(nodeID)] = info

								if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
									log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Error("failed to update replication info, node is active")
								}
							}
						}
					}
				}
			}()
		}
	}
}

func compressKeys(keys [][]byte) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(keys); err != nil {
		return nil, fmt.Errorf("encode error: %w", err)
	}

	compressed, err := zstd.CompressLevel(nil, buf.Bytes(), 22)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	return compressed, nil
}

func decompressKeys(data []byte) ([][]byte, error) {
	decompressed, err := zstd.Decompress(nil, data)
	if err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(decompressed))

	var keys [][]byte
	if err := dec.Decode(&keys); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return keys, nil
}

// FetchAndStore fetches all keys from the local TODO replicate list, fetches value from respective nodes and stores them in the local store
func (s *DHT) FetchAndStore(ctx context.Context) error {
	keys, err := s.store.GetAllToDoRepKeys()
	if err != nil {
		return fmt.Errorf("get all keys error: %w", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(keys))        // Add count of all keys before spawning goroutines
	var successCounter int32 // Create a counter for successful operations

	for i := 0; i < len(keys); i++ {
		key := keys[i]

		go func(info domain.ToRepKey) {
			defer wg.Done()
			sKey := hex.EncodeToString(info.Key)
			n := Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}
			value, err := s.GetValueFromNode(ctx, info.Key, &n)
			if err != nil {
				log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).WithError(err).Error("fetch & store key failed")
				value, err = s.iterateFindValue(ctx, IterateFindValue, info.Key)
				if err != nil {
					log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).WithError(err).Error("iterate fetch for replication failed")
					return
				}
				log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).Info("iterate fetch for replication success")
			}

			if err := s.store.Store(ctx, info.Key, value, 0, false); err != nil {
				log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).WithError(err).Error("fetch & local store key failed")
				return
			}

			if err := s.store.DeleteRepKey(info.Key); err != nil {
				log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).WithError(err).Error("delete key from todo list failed")
				return
			}

			atomic.AddInt32(&successCounter, 1) // Increment the counter atomically

			log.WithContext(ctx).WithField("key", sKey).WithField("ip", info.IP).Info("fetch & store key success")
		}(key)
	}

	wg.Wait()

	log.WithContext(ctx).Infof("Successfully fetched & stored %d keys", atomic.LoadInt32(&successCounter)) // Log the final count

	return nil
}
