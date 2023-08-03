package kademlia

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"encoding/hex"

	"github.com/cenkalti/backoff"
	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

var (
	// defaultReplicationInterval is the default interval for replication.
	defaultReplicationInterval = time.Minute * 10

	// nodeShowUpDeadline is the time after which the node will considered to be permeant offline
	// we'll adjust the keys once the node is permeant offline
	nodeShowUpDeadline = time.Minute * 35

	// check for active & inactive nodes after this interval
	checkNodeActivityInterval = time.Minute * 2

	defaultFetchAndStoreInterval = time.Minute * 10

	defaultBatchFetchAndStoreInterval = time.Minute * 5

	maxBackOff = 45 * time.Second
)

func (s *DHT) getNodeReplicationTimesCopy() map[string]domain.NodeReplicationInfo {
	s.replicationMtx.RLock()
	defer s.replicationMtx.RUnlock()

	copyInfo := make(map[string]domain.NodeReplicationInfo)
	for id, info := range s.nodeReplicationTimes {
		copyInfo[id] = info
	}

	return copyInfo
}

// StartReplicationWorker starts replication
func (s *DHT) StartReplicationWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("replication worker started")

	go s.checkNodeActivity(ctx)
	go s.StartBatchFetchAndStoreWorker(ctx)
	go s.StartFailedFetchAndStoreWorker(ctx)

	for {
		select {
		case <-time.After(defaultReplicationInterval):
			s.Replicate(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing replication worker")
			return nil
		}
	}
}

// StartBatchFetchAndStoreWorker starts replication
func (s *DHT) StartBatchFetchAndStoreWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("batch fetch and store worker started")

	for {
		select {
		case <-time.After(defaultBatchFetchAndStoreInterval):
			s.BatchFetchAndStore(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing batch fetch & store worker")
			return nil
		}
	}
}

// StartFetchAndStoreWorker starts replication
func (s *DHT) StartFailedFetchAndStoreWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("fetch and store worker started")

	for {
		select {
		case <-time.After(defaultFetchAndStoreInterval):
			s.BatchFetchAndStoreFailedKeys(ctx)
		case <-ctx.Done():
			log.P2P().WithContext(ctx).Error("closing fetch & store worker")
			return nil
		}
	}
}

func (s *DHT) updateReplicationNode(ctx context.Context, nodeID []byte, ip string, port int, isActive bool) error {
	s.replicationMtx.RLock()
	info, ok := s.nodeReplicationTimes[string(nodeID)]
	s.replicationMtx.RUnlock()

	if ok {
		info.Active = isActive
		info.UpdatedAt = time.Now()
		info.IP = ip
		info.Port = port
		info.ID = nodeID

		s.replicationMtx.Lock()
		s.nodeReplicationTimes[string(nodeID)] = info
		s.replicationMtx.Unlock()

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

		s.replicationMtx.Lock()
		s.nodeReplicationTimes[string(nodeID)] = info
		s.replicationMtx.Unlock()

		if err := s.store.AddReplicationInfo(ctx, info); err != nil {
			log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).WithField("ip", ip).Error("failed to add replication info")
			return err
		}
	}

	return nil
}

func (s *DHT) updateLastReplicated(ctx context.Context, nodeID []byte, timestamp time.Time) error {
	s.replicationMtx.RLock()
	info, ok := s.nodeReplicationTimes[string(nodeID)]
	s.replicationMtx.RUnlock()
	if !ok {
		return errors.New("node not found")
	}

	info.LastReplicatedAt = &timestamp

	s.replicationMtx.Lock()
	s.nodeReplicationTimes[string(nodeID)] = info
	s.replicationMtx.Unlock()

	if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
		log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).Error("failed to update replication info")
	}

	return nil
}

// Replicate is called periodically by the replication worker to replicate the data across the network by refreshing the buckets
// it iterates over the node replication Info map and replicates the data to the nodes that are active
func (s *DHT) Replicate(ctx context.Context) {
	historicStart, err := s.store.GetOwnCreatedAt(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get own createdAt")
		historicStart = time.Now().Add(-24 * time.Hour * 180)
	}

	log.P2P().WithContext(ctx).WithField("historic-start", historicStart).Info("replicating data")

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
	nrt := s.getNodeReplicationTimesCopy()
	for nodeID, infoVar := range nrt {
		info := infoVar

		logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", string(nodeID))
		if !info.Active {
			adjustNodeKeys := isNodeGoneAndShouldBeAdjusted(info.LastSeen, info.IsAdjusted)
			if adjustNodeKeys {
				if err := s.adjustNodeKeys(ctx, infoVar.CreatedAt, info); err != nil {
					logEntry.WithError(err).Error("failed to adjust node keys")
				} else {
					info.IsAdjusted = true
					info.UpdatedAt = time.Now()

					// lock for write
					s.replicationMtx.Lock()
					s.nodeReplicationTimes[string(nodeID)] = info
					// unlock after write
					s.replicationMtx.Unlock()

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

		logEntry.WithField("from", from).WithField("to", to).Info("getting replication keys")
		replicationKeys := s.store.GetKeysForReplication(ctx, from, to)

		if len(replicationKeys) == 0 {
			// Now closestContactKeys contains all the keys that are in the closest contacts.
			if err := s.updateLastReplicated(ctx, []byte(nodeID), to); err != nil {
				logEntry.Error("replicate update lastReplicated failed")
			} else {
				logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("fetch-keys", 0).Debug("replicate update lastReplicated success")
			}

			continue
		}

		logEntry.WithField("len-rep-keys", len(replicationKeys)).Info("count of replication keys to be checked")
		// Preallocate a slice with a capacity equal to the number of keys.
		closestContactKeys := make([][]byte, 0, len(replicationKeys))
		for i := 0; i < len(replicationKeys); i++ {
			key := replicationKeys[i]
			ignores := s.ignorelist.ToNodeList()
			n := &Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}

			nodeList := s.ht.closestContactsWithInlcudingNode(Alpha, key, ignores, n)

			if nodeList.Exists(n) {
				// the node is supposed to hold this key as it's in the 6 closest contacts
				closestContactKeys = append(closestContactKeys, key)
			}
		}
		logEntry.WithField("len-rep-keys", len(closestContactKeys)).Info("closest contact keys count")

		if len(closestContactKeys) == 0 {
			if err := s.updateLastReplicated(ctx, []byte(nodeID), to); err != nil {
				logEntry.Error("replicate update lastReplicated failed")
			} else {
				logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("closest-contact-keys", 0).Info("replicate update lastReplicated success")
			}
		}

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

		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = maxBackOff

		err = backoff.RetryNotify(func() error {
			response, err := s.sendReplicateData(ctx, n, request)
			if err != nil {
				return err
			}

			if response.Status.Result != ResultOk {
				return errors.New(response.Status.ErrMsg)
			}

			return nil
		}, b, func(err error, duration time.Duration) {
			logEntry.WithError(err).WithField("duration", duration).Error("retrying send replicate data")
		})

		if err != nil {
			logEntry.WithError(err).Error("send replicate data failed after retries")
			continue
		}

		// Now closestContactKeys contains all the keys that are in the closest contacts.
		if err := s.updateLastReplicated(ctx, []byte(nodeID), to); err != nil {
			logEntry.Error("replicate update lastReplicated failed")
		} else {
			logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("expected-rep-keys", len(closestContactKeys)).Info("replicate update lastReplicated success")
		}
	}

	log.P2P().WithContext(ctx).Info("Replication done")
}

func (s *DHT) adjustNodeKeys(ctx context.Context, from time.Time, info domain.NodeReplicationInfo) error {
	replicationKeys := s.store.GetKeysForReplication(ctx, from, time.Now())
	logEntry := log.P2P().WithContext(ctx).WithField("offline-node-ip", info.IP).WithField("offline-node-id", string(info.ID)).
		WithField("total-rep-keys", len(replicationKeys))

	logEntry.WithField("from", from.String()).Info("begin adjusting node keys process for offline node")

	// prepare ignored nodes list but remove the node we are adjusting
	// because we want to find if this node was supposed to hold this key
	ignores := s.ignorelist.ToNodeList()
	var updatedIgnored []*Node
	for _, ignore := range ignores {
		if !bytes.Equal(ignore.ID, info.ID) {
			updatedIgnored = append(updatedIgnored, ignore)
		}
	}

	nodeKeysMap := make(map[string][]string)
	for _, key := range replicationKeys {

		offNode := &Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}

		// get closest contacts to the key
		nodeList := s.ht.closestContactsWithInlcudingNode(Alpha, key, updatedIgnored, offNode)

		// check if the node that is gone was supposed to hold the key
		if !nodeList.Exists(offNode) {
			// the node is not supposed to hold this key as its not in 6 closest contacts
			continue
		}

		for i := 0; i < len(nodeList.Nodes); i++ {
			if nodeList.Nodes[i].IP == info.IP {
				continue // because we do not want to send request to the server that's already offline
			}

			// If the node is supposed to hold the key, we map the node's info to the key
			nodeInfo := generateKeyFromNode(nodeList.Nodes[i])
			// convert key to string representation
			keyString := hex.EncodeToString(key)

			// append the key to the list of keys that the node is supposed to have
			nodeKeysMap[nodeInfo] = append(nodeKeysMap[nodeInfo], keyString)

		}
	}

	// iterate over the map and send the keys to the node
	// Loop over the map
	successCount := 0
	failureCount := 0

	for nodeInfoKey, keys := range nodeKeysMap {
		logEntry.WithField("adjust-to-node", nodeInfoKey).WithField("to-adjust-keys-len", len(keys)).Info("sending adjusted replication keys to node")
		// Retrieve the node object from the key
		node, err := getNodeFromKey(nodeInfoKey)
		if err != nil {
			logEntry.WithError(err).Error("Failed to parse node info from key")
			return fmt.Errorf("failed to parse node info from key: %w", err)
		}

		// Convert the hexadecimal string keys back to byte slices
		byteKeys := make([][]byte, len(keys))
		for i, key := range keys {
			byteKey, err := hex.DecodeString(key)
			if err != nil {
				logEntry.WithError(err).Errorf("Failed to decode key: %s", key)
				return fmt.Errorf("failed to decode key: %s: %w", key, err)
			}
			byteKeys[i] = byteKey
		}

		data, err := compressKeys(byteKeys)
		if err != nil {
			logEntry.WithField("adjust-rep-keys", len(byteKeys)).WithError(err).Error("unable to compress keys - adjust keys failed")
			return fmt.Errorf("unable to compress keys - adjust keys failed: %w", err)
		}

		// TODO: Check if data size is bigger than 32 MB
		request := &ReplicateDataRequest{
			Keys: data,
		}

		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = maxBackOff

		err = backoff.RetryNotify(func() error {
			response, err := s.sendReplicateData(ctx, node, request)
			if err != nil {
				return err
			}

			if response.Status.Result != ResultOk {
				return errors.New(response.Status.ErrMsg)
			}

			successCount++
			return nil
		}, b, func(err error, duration time.Duration) {
			logEntry.WithError(err).WithField("duration", duration).Error("retrying send replicate data")
		})

		if err != nil {
			logEntry.WithError(err).Error("send replicate data failed after retries")
			failureCount++
		}
	}

	totalCount := successCount + failureCount

	if totalCount > 0 { // Prevent division by zero
		successRate := float64(successCount) / float64(totalCount) * 100

		if successRate < 75 {
			// Success rate is less than 75%
			return fmt.Errorf("adjust keys success rate is less than 75%%: %v", successRate)
		}
	} else {
		return fmt.Errorf("adjust keys totalCount is 0")
	}

	info.IsAdjusted = true
	info.UpdatedAt = time.Now()

	if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
		return fmt.Errorf("replicate update isAdjusted failed: %v", err)
	}

	logEntry.Info("offline node was successfully adjusted")

	return nil
}

func isNodeGoneAndShouldBeAdjusted(lastSeen *time.Time, isAlreadyAdjusted bool) bool {
	if lastSeen == nil {
		log.Info("lastSeen is nil - aborting node adjustment")
		return false
	}

	return time.Since(*lastSeen) > nodeShowUpDeadline && !isAlreadyAdjusted
}
