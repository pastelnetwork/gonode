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

// StartReplicationWorker starts replication
func (s *DHT) StartReplicationWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("replication worker started")

	go s.checkNodeActivity(ctx)
	go s.StartBatchFetchAndStoreWorker(ctx)
	go s.StartFailedFetchAndStoreWorker(ctx)

	for {
		select {
		case <-time.After(defaultReplicationInterval):
			//log.WithContext(ctx).Info("replication worker disabled")
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

// StartFailedFetchAndStoreWorker starts replication
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
	// check if record exists
	ok, err := s.store.RecordExists(string(nodeID))
	if err != nil {
		return fmt.Errorf("err checking if replication info record exists: %w", err)
	}

	now := time.Now().UTC()
	info := domain.NodeReplicationInfo{
		UpdatedAt: time.Now().UTC(),
		Active:    isActive,
		IP:        ip,
		Port:      port,
		ID:        nodeID,
		LastSeen:  &now,
	}

	if ok {
		if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
			log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).WithField("ip", ip).Error("failed to update replication info")
			return err
		}
	} else {
		if err := s.store.AddReplicationInfo(ctx, info); err != nil {
			log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).WithField("ip", ip).Error("failed to add replication info")
			return err
		}
	}

	return nil
}

func (s *DHT) updateLastReplicated(ctx context.Context, nodeID []byte, timestamp time.Time) error {
	if err := s.store.UpdateLastReplicated(ctx, string(nodeID), timestamp); err != nil {
		log.P2P().WithContext(ctx).WithError(err).WithField("node_id", string(nodeID)).Error("failed to update replication info last replicated")
	}

	return nil
}

// Replicate is called periodically by the replication worker to replicate the data across the network by refreshing the buckets
// it iterates over the node replication Info map and replicates the data to the nodes that are active
func (s *DHT) Replicate(ctx context.Context) {
	historicStart, err := s.store.GetOwnCreatedAt(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("unable to get own createdAt")
		historicStart = time.Now().UTC().Add(-24 * time.Hour * 180)
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

	repInfo, err := s.store.GetAllReplicationInfo(ctx)
	if err != nil {
		log.P2P().WithContext(ctx).WithError(err).Errorf("get all replicationInfo failed")
		return
	}

	if len(repInfo) == 0 {
		log.P2P().WithContext(ctx).Info("no replication info found")
		return
	}

	from := historicStart
	if repInfo[0].LastReplicatedAt != nil {
		from = *repInfo[0].LastReplicatedAt
	}

	log.P2P().WithContext(ctx).WithField("from", from).Info("getting all possible replication keys")
	to := time.Now().UTC()
	replicationKeys := s.store.GetKeysForReplication(ctx, from, to)

	ignores := s.ignorelist.ToNodeList()
	closestContactsMap := make(map[string][][]byte)

	for i := 0; i < len(replicationKeys); i++ {
		self := &Node{ID: s.ht.self.ID, IP: s.externalIP, Port: s.ht.self.Port}
		decKey, _ := hex.DecodeString(replicationKeys[i].Key)
		closestContactsMap[replicationKeys[i].Key] = s.ht.closestContactsWithInlcudingNode(Alpha, decKey, ignores, self).NodeIDs()
	}

	for _, info := range repInfo {
		if !info.Active {
			s.checkAndAdjustNode(ctx, info, historicStart)
			continue
		}

		logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", string(info.ID))

		start := historicStart
		if info.LastReplicatedAt != nil {
			start = *info.LastReplicatedAt
		}

		idx := replicationKeys.FindFirstAfter(start)
		if idx == -1 {
			// Now closestContactKeys contains all the keys that are in the closest contacts.
			if err := s.updateLastReplicated(ctx, info.ID, to); err != nil {
				logEntry.Error("replicate update lastReplicated failed")
			} else {
				logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("fetch-keys", 0).Debug("no replication keys - replicate update lastReplicated success")
			}

			continue
		}
		countToSendKeys := len(replicationKeys) - idx
		logEntry.WithField("len-rep-keys", countToSendKeys).Info("count of replication keys to be checked")
		// Preallocate a slice with a capacity equal to the number of keys.
		closestContactKeys := make([]string, 0, countToSendKeys)

		for i := idx; i < len(replicationKeys); i++ {
			for j := 0; j < len(closestContactsMap[replicationKeys[i].Key]); j++ {
				if bytes.Equal(closestContactsMap[replicationKeys[i].Key][j], info.ID) {
					// the node is supposed to hold this key as it's in the 6 closest contacts
					closestContactKeys = append(closestContactKeys, replicationKeys[i].Key)
				}
			}
		}

		logEntry.WithField("len-rep-keys", len(closestContactKeys)).Info("closest contact keys count")

		if len(closestContactKeys) == 0 {
			if err := s.updateLastReplicated(ctx, info.ID, to); err != nil {
				logEntry.Error("replicate update lastReplicated failed")
			} else {
				logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("closest-contact-keys", 0).Info("no closest keys found - replicate update lastReplicated success")
			}

			continue
		}

		data, err := compressKeysStr(closestContactKeys)
		if err != nil {
			logEntry.WithField("len-rep-keys", len(closestContactKeys)).WithError(err).Error("unable to compress keys - replication failed")
			continue
		}

		// TODO: Check if data size is bigger than 32 MB
		request := &ReplicateDataRequest{
			Keys: data,
		}

		n := &Node{ID: info.ID, IP: info.IP, Port: info.Port}

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
		if err := s.updateLastReplicated(ctx, info.ID, to); err != nil {
			logEntry.Error("replicate update lastReplicated failed")
		} else {
			logEntry.WithField("node", info.IP).WithField("to", to.String()).WithField("expected-rep-keys", len(closestContactKeys)).Info("replicate update lastReplicated success")
		}
	}

	log.P2P().WithContext(ctx).Info("Replication done")
}

func (s *DHT) adjustNodeKeys(ctx context.Context, from time.Time, info domain.NodeReplicationInfo) error {
	replicationKeys := s.store.GetKeysForReplication(ctx, from, time.Now().UTC())
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
	for i := 0; i < len(replicationKeys); i++ {

		offNode := &Node{ID: []byte(info.ID), IP: info.IP, Port: info.Port}

		// get closest contacts to the key
		key, _ := hex.DecodeString(replicationKeys[i].Key)
		nodeList := s.ht.closestContactsWithInlcudingNode(Alpha+1, key, updatedIgnored, offNode) // +1 because we want to include the node we are adjusting
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

			// append the key to the list of keys that the node is supposed to have
			nodeKeysMap[nodeInfo] = append(nodeKeysMap[nodeInfo], replicationKeys[i].Key)

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

		data, err := compressKeysStr(keys)
		if err != nil {
			logEntry.WithField("adjust-rep-keys", len(keys)).WithError(err).Error("unable to compress keys - adjust keys failed")
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

	if err := s.store.UpdateIsAdjusted(ctx, string(info.ID), true); err != nil {
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

func (s *DHT) checkAndAdjustNode(ctx context.Context, info domain.NodeReplicationInfo, start time.Time) {
	logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", string(info.ID))
	adjustNodeKeys := isNodeGoneAndShouldBeAdjusted(info.LastSeen, info.IsAdjusted)
	if adjustNodeKeys {
		if err := s.adjustNodeKeys(ctx, start, info); err != nil {
			logEntry.WithError(err).Error("failed to adjust node keys")
		} else {
			info.IsAdjusted = true
			info.UpdatedAt = time.Now().UTC()

			if err := s.store.UpdateIsAdjusted(ctx, string(info.ID), true); err != nil {
				logEntry.WithError(err).Error("failed to update replication info, set isAdjusted to true")
			} else {
				logEntry.Info("set isAdjusted to true")
			}
		}
	}

	logEntry.Info("replication node not active, skipping over it.")
}
