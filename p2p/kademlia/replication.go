package kademlia

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
	"github.com/pastelnetwork/gonode/p2p/kademlia/domain"
)

var (
	// defaultReplicationTime is the default time for replication - in case lastReplicated is nil
	// it will be used as the lastReplicated time
	defaultReplicationTime = time.Hour * 24 * 30 * 6 // 6 months

	// defaultReplicationInterval is the default interval for replication.
	defaultReplicationInterval = time.Minute * 5

	// nodeShowUpDeadline is the time after which the node will considered to be permeant offline
	// we'll adjust the keys once the node is permeant offline
	nodeShowUpDeadline = time.Minute * 35

	// check for active & inactive nodes after this interval
	checkNodeActivityInterval = time.Minute * 3
)

// StartReplicationWorker starts replication
func (s *DHT) StartReplicationWorker(ctx context.Context) error {
	log.P2P().WithContext(ctx).Info("replication worker started")

	go s.checkNodeActivity(ctx)

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

	for nodeID, infoVar := range s.nodeReplicationTimes {
		info := infoVar

		logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", nodeID)
		if !info.Active {
			adjustNodeKeys := isNodeGoneAndShouldBeAdjusted(info.LastReplicatedAt, info.IsAdjusted)
			if adjustNodeKeys {
				if err := s.adjustNodeKeys(ctx, info); err != nil {
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

		from := time.Now().Add(-defaultReplicationTime)
		if info.LastReplicatedAt != nil {
			from = *info.LastReplicatedAt

			// TODO: REMOVE THIS BEFORE RELEASE
			//if time.Since(from) > defaultReplicationTime {
			//	from = time.Now().Add(-defaultReplicationTime)
			//}
		}

		logEntry.WithField("from", from.String()).Info("replication from")
		replicationKeys := s.store.GetKeysForReplication(ctx, from)
		now := time.Now()
		logEntry = logEntry.WithField("len-rep-keys", len(replicationKeys))
		logEntry.Info("replication keys")

		if len(replicationKeys) == 0 {
			continue
		}

		replicatedCount := 0
		for _, key := range replicationKeys {
			ignores := s.ignorelist.ToNodeList()
			nodeList := s.ht.closestContacts(Alpha, key, ignores)

			n := &Node{ID: []byte(nodeID), IP: info.IP, Port: info.Port}
			if !nodeList.Exists(n) {
				// the node is not supposed to hold this key as its not in 6 closest contacts
				continue
			}
			logEntry.WithField("key", hex.EncodeToString(key)).WithField("ip", info.IP).Info("this key is supposed to be hold by this node")

			value, typ, err := s.store.RetrieveWithType(ctx, key)
			if err != nil {
				log.P2P().WithContext(ctx).WithError(err).Error("replicate store retrieve failed")
				continue
			}

			if value != nil {
				request := &StoreDataRequest{Data: value, Type: typ}
				response, err := s.sendStoreData(ctx, n, request)
				if err != nil {
					logEntry.WithError(err).Error("replicate store data failed")
				} else if response.Status.Result != ResultOk {
					logEntry.WithError(errors.New(response.Status.ErrMsg)).Error("reply replicate store data failed")
				} else {
					replicatedCount++
					logEntry.Info("replicate store data success")
				}
			}
		}

		if err := s.updateLastReplicated(ctx, []byte(nodeID), now); err != nil {
			logEntry.Error("replicate update lastReplicated failed")
		} else {
			logEntry.WithField("node", info.IP).WithField("expected-rep-keys", len(replicationKeys)).
				WithField("keys-replicated", replicatedCount).Info("replicate update lastReplicated success")
		}

	}

	log.P2P().WithContext(ctx).Info("Replication done")
}

func (s *DHT) adjustNodeKeys(ctx context.Context, info domain.NodeReplicationInfo) error {
	logEntry := log.P2P().WithContext(ctx).WithField("rep-ip", info.IP).WithField("rep-id", info.ID)

	from := time.Now().Add(-defaultReplicationTime)
	logEntry.WithField("from", from.String()).Info("adjusting node keys")

	replicationKeys := s.store.GetKeysForReplication(ctx, from)
	logEntry = logEntry.WithField("len-rep-keys", len(replicationKeys))
	logEntry.Info("replication keys")

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
		logEntry.WithField("key", hex.EncodeToString(key)).WithField("ip", info.IP).Info("this key is supposed to be hold by this node")

		// get the value
		value, typ, err := s.store.RetrieveWithType(ctx, key)
		if err != nil {
			log.P2P().WithContext(ctx).WithError(err).Error("replicate store retrieve failed")
			continue
		}

		key := s.hashKey(value)

		// iterative store the data
		if _, err := s.iterate(ctx, IterateStore, key, value, typ); err != nil {
			log.WithContext(ctx).WithError(err).Error("replicate iterate data store failure")
			continue
		}
	}

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
