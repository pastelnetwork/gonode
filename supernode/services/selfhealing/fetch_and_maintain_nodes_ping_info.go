package selfhealing

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pastelnetwork/gonode/common/storage/local"
	"github.com/pastelnetwork/gonode/common/types"
	"sync"
	"time"

	"github.com/pastelnetwork/gonode/common/errors"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/pastel"
	pb "github.com/pastelnetwork/gonode/proto/supernode"
)

// Define a timeout duration
const (
	timeoutDuration = 10 * time.Second
	UTCTimeLayout   = "2006-01-02T15:04:05Z"
)

// FetchAndMaintainPingInfo fetch and maintains the ping info in db for every node
func (task *SHTask) FetchAndMaintainPingInfo(ctx context.Context) error {
	log.WithContext(ctx).Infoln("Self Healing Ping Nodes Worker invoked")

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

			timeBeforePing := time.Now().UTC()
			res, err := task.ping(ctx, req, node.ExtAddress)
			if err != nil {
				log.WithContext(ctx).WithField("sn_address", node.ExtAddress).WithError(err).
					Error("error pinging sn")

				pi := types.PingInfo{
					SupernodeID:      node.ExtKey,
					IPAddress:        node.ExtAddress,
					IsOnline:         false,
					IsAdjusted:       false,
					LastResponseTime: 0.0,
				}

				if err := task.StorePingInfo(ctx, pi); err != nil {
					log.WithContext(ctx).WithError(err).Error("error storing ping info")
				}

				return
			}

			respondedAt, err := time.Parse(UTCTimeLayout, res.RespondedAt)
			if err != nil {
				log.WithContext(ctx).WithError(err)
			}

			pi := types.PingInfo{
				SupernodeID:      node.ExtKey,
				IPAddress:        node.ExtAddress,
				IsOnline:         res.IsOnline,
				IsAdjusted:       false,
				LastResponseTime: respondedAt.Sub(timeBeforePing).Seconds(),
				LastSeen:         sql.NullTime{Time: respondedAt, Valid: true},
			}

			if err := task.StorePingInfo(ctx, pi); err != nil {
				log.WithContext(ctx).WithError(err).Error("error storing ping info")
			}
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

// StorePingInfo stores the ping info to db
func (task *SHTask) StorePingInfo(ctx context.Context, info types.PingInfo) error {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return err
	}

	existedInfo, err := task.GetPingInfoFromDB(ctx, info.SupernodeID)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("error retrieving existed ping info")
		return err
	}

	inf := GetPingInfoToInsert(existedInfo, &info)

	if store != nil && inf != nil {
		defer store.CloseHistoryDB(ctx)

		err = store.UpsertPingHistory(*inf)
		if err != nil {
			log.WithContext(ctx).
				WithField("supernode_id", info.SupernodeID).
				Error("error storing ping history")

			return err
		}

	}

	return nil
}

// GetPingInfoFromDB get the ping info from db
func (task *SHTask) GetPingInfoFromDB(ctx context.Context, supernodeID string) (*types.PingInfo, error) {
	store, err := local.OpenHistoryDB()
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("Error Opening DB")
		return nil, err
	}

	var info *types.PingInfo

	if store != nil {
		defer store.CloseHistoryDB(ctx)

		info, err = store.GetPingInfoBySupernodeID(supernodeID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return &types.PingInfo{}, nil
			}

			log.WithContext(ctx).
				WithField("supernode_id", info.SupernodeID).
				Error("error retrieving ping history")

			return nil, err
		}
	}

	return info, nil
}

// GetPingInfoToInsert finalise the ping info after comparing the existed and new info that needs to be inserted/updated
func GetPingInfoToInsert(existedInfo, info *types.PingInfo) *types.PingInfo {
	if existedInfo == nil {
		return nil
	}

	info.TotalPings = existedInfo.TotalPings + 1

	if !existedInfo.LastSeen.Valid { //for the first row
		info.LastSeen = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	}

	if info.IsOnline {
		info.TotalSuccessfulPings = existedInfo.TotalSuccessfulPings + 1
	}

	if !info.IsOnline {
		info.TotalSuccessfulPings = existedInfo.TotalSuccessfulPings
		info.LastSeen = existedInfo.LastSeen
	}

	info.CumulativeResponseTime = existedInfo.CumulativeResponseTime + info.LastResponseTime

	var avgPingResponseTime float64
	if info.TotalSuccessfulPings != 0 {
		avgPingResponseTime = info.CumulativeResponseTime / float64(info.TotalSuccessfulPings)
	}
	info.AvgPingResponseTime = avgPingResponseTime

	return info
}
