package selfhealing

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// Define a timeout duration
const timeoutDuration = 10 * time.Second

//FetchAndMaintainPingInfo fetch and maintains the ping info in db for every node
func (task *SHTask) FetchAndMaintainPingInfo(ctx context.Context) error {
	log.WithContext(ctx).Println("Self Healing Ping Nodes Worker invoked")

	nodesToPing, err := task.getNodesAddressesToConnect(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving nodes addresses")
	}

	if err := task.pingNodes(ctx, nodesToPing); err != nil {
		log.WithContext(ctx).WithError(err).Error("error pinging nodes")
	}

	return nil
}

func (task *SHTask) getNodesAddressesToConnect(ctx context.Context) ([]pastel.MasterNode, error) {
	supernodes, err := task.SuperNodeService.PastelClient.MasterNodesExtra(ctx)
	if err != nil {
		log.WithContext(ctx).WithField("method", "FetchAndMaintainPingInfo").WithError(err).Warn("could not get Supernode extra: ", err.Error())
		return nil, err
	}

	mapSupernodes := make(map[string]pastel.MasterNode)
	for _, mn := range supernodes {
		if mn.ExtAddress == "" || mn.ExtKey == "" {
			log.WithContext(ctx).WithField("method", "FetchAndMaintainPingInfo").
				WithField("node_id", mn.ExtKey).Warn("node address or node id is empty")

			continue
		}

		mapSupernodes[mn.ExtKey] = mn
	}

	return supernodes, nil
}

// pingNodes will ping the nodes and record their responses
func (task *SHTask) pingNodes(ctx context.Context, nodesToPing pastel.MasterNodes) error {
	if nodesToPing == nil {
		return errors.Errorf("no nodes found to connect for maintaining ping info")
	}

	req := &pb.PingRequest{
		SenderId: task.nodeID,
	}

	var wg sync.WaitGroup
	for _, node := range nodesToPing {
		node := node
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := task.ping(ctx, req, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithField("sn_address", node.ExtAddress).WithError(err).
					Error("error pinging sn")
				return

				//Ping Failed Response
				//a) if a node didn’t respond, set is_online = false, increment total_pings
				//a.1) Check if last_seen is before 20 minutes, if so, set is_on_watchlist field to true.
			}

			//ping success scenario
			// b) if a node responded, update last_seen to current time, increment total_pings and total_successful_pings and recalculate average_response_time (use total_successful_pings)
			//  ===> also set is_on_watchlist to false.

		}()
	}

	wg.Wait()

	return nil
}

// ping just pings the given node address
func (task *SHTask) ping(ctx context.Context, req *pb.PingRequest, supernodeAddr string) (*pb.PingResponse, error) {
	log.WithContext(ctx).Info("pinging supernode: " + supernodeAddr)

	// Create a context with timeout
	pingCtx, cancel := context.WithTimeout(ctx, timeoutDuration)
	defer cancel()

	// Connect over gRPC with timeout context
	nodeClientConn, err := task.nodeClient.Connect(pingCtx, supernodeAddr)
	if err != nil {
		err = fmt.Errorf("could not use node client to connect to: %s, error: %v", supernodeAddr, err)
		log.WithContext(ctx).Warn(err.Error())
		return nil, err
	}
	defer nodeClientConn.Close()

	selfHealingIF := nodeClientConn.SelfHealingChallenge()

	// Use the timeout context for the ping operation
	pingResponse, err := selfHealingIF.Ping(pingCtx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.WithContext(ctx).Warnf("Ping to %s timed out", supernodeAddr)
		} else {
			log.WithContext(ctx).Warn(err.Error())
		}
		return nil, err
	}

	return pingResponse, nil
}
