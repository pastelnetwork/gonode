package kademlia

import (
	"context"
	"time"

	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/utils"
)

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
					nrt := s.getNodeReplicationTimesCopy()

					for nodeID, info := range nrt {
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
						_, err := s.network.Call(ctx, request, false)
						if err != nil && info.Active {
							log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).
								Error("failed to ping node, setting node to inactive")

							// add node to ignore list
							// we maintain this list to avoid pinging nodes that are not responding
							s.ignorelist.IncrementCount(node)

							// mark node as inactive in database
							info.Active = false
							info.UpdatedAt = time.Now()

							s.replicationMtx.Lock()
							s.nodeReplicationTimes[string(nodeID)] = info
							s.replicationMtx.Unlock()

							if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
								log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Error("failed to update replication info, node is inactive")
							}

						} else if err == nil {
							if !info.Active {
								log.P2P().WithContext(ctx).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Info("node found to be active again")

								// remove node from ignore list
								s.ignorelist.Delete(node)

								// mark node as active in database
								info.Active = true
								info.IsAdjusted = false
								info.UpdatedAt = time.Now()

								s.replicationMtx.Lock()
								s.nodeReplicationTimes[string(nodeID)] = info
								s.replicationMtx.Unlock()

								if err := s.store.UpdateReplicationInfo(ctx, info); err != nil {
									log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Error("failed to update replication info, node is active")
								}
							}

							s.replicationMtx.RLock()
							upInfo := s.nodeReplicationTimes[string(nodeID)]
							s.replicationMtx.RUnlock()

							now := time.Now()
							upInfo.LastSeen = &now

							s.replicationMtx.Lock()
							s.nodeReplicationTimes[string(nodeID)] = upInfo
							s.replicationMtx.Unlock()

							if err := s.store.UpdateLastSeen(ctx, string(nodeID)); err != nil {
								log.WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(nodeID)).Error("failed to update last seen")
							}
						}
					}
				}
			}()
		}
	}
}
