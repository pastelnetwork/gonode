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
			if !utils.CheckInternetConnectivity() {
				log.WithContext(ctx).Info("no internet connectivity, not checking node activity")
			} else {
				repInfo, err := s.store.GetAllReplicationInfo(ctx)
				if err != nil {
					log.P2P().WithContext(ctx).WithError(err).Errorf("get all replicationInfo failed")
				}

				for _, info := range repInfo {
					// new a ping request message
					node := &Node{
						ID:   []byte(info.ID),
						IP:   info.IP,
						Port: info.Port,
					}

					request := s.newMessage(Ping, node, nil)

					// invoke the request and handle the response
					_, err := s.network.Call(ctx, request, false)
					if err != nil {
						log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(info.ID)).
							Debug("failed to ping node")
						if info.Active {
							log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(info.ID)).
								Error("setting node to inactive")

							// add node to ignore list
							// we maintain this list to avoid pinging nodes that are not responding
							s.ignorelist.IncrementCount(node)
							// remove from route table
							s.removeNode(ctx, node)

							// mark node as inactive in database
							if err := s.store.UpdateIsActive(ctx, string(info.ID), false, false); err != nil {
								log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(info.ID)).Error("failed to update replication info, node is inactive")
							}
						}

					} else if err == nil {
						// remove node from ignore list
						s.ignorelist.Delete(node)

						if !info.Active {
							log.P2P().WithContext(ctx).WithField("ip", info.IP).WithField("node_id", string(info.ID)).Info("node found to be active again")
							// add node adds in the route table
							s.addNode(ctx, node)
							if err := s.store.UpdateIsActive(ctx, string(info.ID), true, false); err != nil {
								log.P2P().WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(info.ID)).Error("failed to update replication info, node is inactive")
							}
						}

						if err := s.store.UpdateLastSeen(ctx, string(info.ID)); err != nil {
							log.WithContext(ctx).WithError(err).WithField("ip", info.IP).WithField("node_id", string(info.ID)).Error("failed to update last seen")
						}
					}
				}
			}
		}
	}
}
