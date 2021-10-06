package metadb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pastelnetwork/gonode/metadb/rqlite/store"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
)

// initLeadershipTransferTrigger is supposed to transfer leadership from current leader
// to another Top Ranked SN every N blocks so the system remains decentralized
func (s *service) initLeadershipTransferTrigger(ctx context.Context,
	leaderCheckInterval time.Duration, blockCountCheckInterval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(leaderCheckInterval):
			// this logic is only for leader, we don't want to exit for loop
			// as current node can become leader in future
			if !s.IsLeader() {
				log.WithContext(ctx).Debug("leadership transfer: non-leader node")
				break
			}
			log.WithContext(ctx).Info("leadership transfer: leader node! Initiaing scheduler")

			count, err := s.pastelClient.GetBlockCount(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("initLeaderElectionTrigger failed to call pastel client BlockCount")
				break
			}
			// store current block count so we may check block interval
			log.WithContext(ctx).Infof("leadership transfer: initial block count: %d", count)
			s.currentBlockCount = count

			// this will wait for configured block intervals to pass & then
			// execute transfer leadership logic
			s.leadershipTransfer(ctx, blockCountCheckInterval)
		}
	}
}

// leadershipTransfer will wait for configured block intervals to pass & then
func (s *service) leadershipTransfer(ctx context.Context, interval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			// we'd want to make sure that only leader executes the transfer logic
			// so checking it again
			if !s.IsLeader() {
				// so its no longer the leader, go back to initLeadershipTransferTrigger
				// and wait for itself to become leader
				log.WithContext(ctx).Info("leadership transfer: node is no longer a leader.")

				return
			}

			count, err := s.pastelClient.GetBlockCount(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("leadershipTransfer failed to call pastel client BlockCount")
				break
			}
			log.WithContext(ctx).Debugf("leadership transfer: current block count: %d", count)

			// check if configured block inteval has passed
			if count-s.currentBlockCount < s.config.BlockInterval {
				// No! so will try again after 5 seconds
				break
			}
			log.WithContext(ctx).Info("leadership transfer: block interval surpass, transferring leadership")

			nodes, err := s.pastelClient.MasterNodesTop(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("leadershipTransfer failed to call pastel client MasterNodesTop")
				break
			}

			// get server id & address of next leader node
			id, address, err := s.getNextLeaderAddressAndID(ctx, nodes)
			if err != nil {
				log.WithContext(ctx).WithError(err).Error("leadershipTransfer - getNextLeaderAddress failed")
				break
			}
			log.WithContext(ctx).Infof("leadership transfer: transferring to: %s", address)

			// transfer leadership to the next leader
			if err := s.db.TransferLeadership(id, address); err != nil {
				if errors.Is(err, store.ErrNotLeader) {
					log.WithContext(ctx).Warn("leadershipTransfer called by non-leader node")
					return
				}

				log.WithContext(ctx).WithError(err).Error("leadershipTransfer - failed to transfer Leadership")
			} else {
				log.WithContext(ctx).WithField("next-leader", address).
					Info("leadership transfer: leadership transfer successful.")

				// we only want to update the count if leadership was transferred successfully
				// otherwise it will wait again for the block interval to pass while still being leader
				s.currentBlockCount = count
			}
		}
	}
}

func (s *service) getNextLeaderAddressAndID(ctx context.Context, nodes pastel.MasterNodes) (nodeID string, nodeAddress string, err error) {
	if len(nodes) < 1 {
		return "", "", errors.New("no node found to transfer leadership to")
	}

	for _, node := range nodes {
		address := node.Address
		segments := strings.Split(address, ":")
		if len(segments) != 2 {
			log.WithContext(ctx).WithField("address", address).Error("leadershipTransfer: malformed address")
			continue
		}
		nodeAddress = fmt.Sprintf("%s:%d", segments[0], s.config.RaftPort)
		nodeID, err = s.db.GetServerID(nodeAddress)
		if err != nil {
			log.WithContext(ctx).WithField("address", address).Error("leadershipTransfer: Get ServerID failed")
			continue
		}

		if nodeID == "" {
			log.WithContext(ctx).WithField("address", address).Error("leadershipTransfer: unable to find server ID")
			continue
		}

		break
	}

	if nodeAddress == "" || nodeID == "" {
		return "", "", fmt.Errorf("no node found to have correct address - node count: %d", len(nodes))
	}

	return nodeID, nodeAddress, nil
}
